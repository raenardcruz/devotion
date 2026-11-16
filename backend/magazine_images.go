package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type MagazineReadingImageRequest struct {
	Date            string `json:"date"`
	ReadingKey      string `json:"readingKey"`
	Citation        string `json:"citation"`
	Text            string `json:"text,omitempty"`
	Context         string `json:"context,omitempty"`
	ForceRegenerate bool   `json:"forceRegenerate,omitempty"`
}

type MagazineReadingImageResponse struct {
	Date       string `json:"date"`
	ReadingKey string `json:"readingKey"`
	Image      string `json:"image"`
	R2Key      string `json:"r2Key,omitempty"`
	Cached     bool   `json:"cached"`
}

var (
	memMagazineImagesMu sync.RWMutex
	memMagazineImages   = make(map[string]MagazineReadingImageResponse)
)

func normalizeMagazineReadingKey(raw string) string {
	key := strings.ToLower(strings.TrimSpace(raw))
	switch key {
	case "cover", "first_reading", "responsorial_psalm", "second_reading", "gospel", "pope_quote":
		return key
	default:
		return "reading"
	}
}

func aspectRatioForReadingKey(readingKey string) string {
	switch normalizeMagazineReadingKey(readingKey) {
	case "cover":
		return "4:3"
	case "responsorial_psalm":
		return "16:9"
	default:
		return "16:9"
	}
}

func magazineImageCacheKey(date, readingKey string) string {
	return fmt.Sprintf("magazine:image:%s:%s", strings.TrimSpace(date), normalizeMagazineReadingKey(readingKey))
}

func getCachedMagazineImage(ctx context.Context, date, readingKey string) (MagazineReadingImageResponse, bool) {
	key := magazineImageCacheKey(date, readingKey)

	memMagazineImagesMu.RLock()
	if cached, ok := memMagazineImages[key]; ok && cached.Image != "" {
		memMagazineImagesMu.RUnlock()
		cached.Cached = true
		return cached, true
	}
	memMagazineImagesMu.RUnlock()

	if redisClient != nil {
		val, err := redisClient.Get(ctx, key).Result()
		if err == nil && strings.TrimSpace(val) != "" {
			var cached MagazineReadingImageResponse
			if unmErr := json.Unmarshal([]byte(val), &cached); unmErr == nil && cached.Image != "" {
				memMagazineImagesMu.Lock()
				memMagazineImages[key] = cached
				memMagazineImagesMu.Unlock()
				cached.Cached = true
				return cached, true
			}
		}
	}

	return MagazineReadingImageResponse{}, false
}

func setCachedMagazineImage(ctx context.Context, resp MagazineReadingImageResponse) {
	key := magazineImageCacheKey(resp.Date, resp.ReadingKey)

	memMagazineImagesMu.Lock()
	memMagazineImages[key] = resp
	memMagazineImagesMu.Unlock()

	if redisClient != nil {
		if data, err := json.Marshal(resp); err == nil {
			_ = redisClient.Set(ctx, key, string(data), 7*24*time.Hour).Err()
		}
	}
}

func truncateForPrompt(input string, maxLen int) string {
	cleaned := strings.Join(strings.Fields(stripHTMLTags(input)), " ")
	if len(cleaned) <= maxLen {
		return cleaned
	}
	return cleaned[:maxLen] + "..."
}

func buildReadingIllustrationPrompt(readingKey, citation, text, commentary string) string {
	key := normalizeMagazineReadingKey(readingKey)
	snippet := truncateForPrompt(text, 450)
	contextSnippet := truncateForPrompt(commentary, 300)

	subjectLine := fmt.Sprintf("Sacred Scripture passage: %s.", citation)
	if key == "cover" {
		subjectLine = fmt.Sprintf("Magazine cover centerpiece inspired by today's Gospel (%s).", citation)
	} else if key == "responsorial_psalm" {
		subjectLine = fmt.Sprintf("Poetic liturgical scene for Responsorial Psalm (%s).", citation)
	}

	var details strings.Builder
	details.WriteString(subjectLine)
	if snippet != "" {
		details.WriteString("\nPassage excerpt: ")
		details.WriteString(snippet)
	}
	if contextSnippet != "" {
		details.WriteString("\nTheological theme: ")
		details.WriteString(contextSnippet)
	}

	return fmt.Sprintf(`Create a museum-grade editorial sacred art illustration for a Catholic liturgical print magazine.
%s

Visual Style & Direction:
- Classical sacred fine-art editorial illustration blending warm Renaissance oil-wash, illuminated manuscript warmth, and subtle chiaroscuro lighting.
- Color palette: warm parchment cream, deep liturgical burgundy, antique gold leaf accents, olive/cypress greens, and rich charcoal ink shadows.
- Reverent, dignified, contemplative composition suitable for a high-end print magazine.
- Strictly NO text, NO typography, NO titles, NO Latin/English words, NO watermarks, NO borders, and NO modern UI elements inside the image.`, details.String())
}

func magazineReadingImageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req MagazineReadingImageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	date := strings.TrimSpace(req.Date)
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	readingKey := normalizeMagazineReadingKey(req.ReadingKey)
	citation := strings.TrimSpace(req.Citation)
	if citation == "" && req.Text == "" {
		respondWithError(w, http.StatusBadRequest, "citation or text is required")
		return
	}

	ctx := r.Context()

	if !req.ForceRegenerate {
		if cached, ok := getCachedMagazineImage(ctx, date, readingKey); ok {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(cached)
			return
		}
	}

	prompt := buildReadingIllustrationPrompt(readingKey, citation, req.Text, req.Context)
	aspectRatio := aspectRatioForReadingKey(readingKey)
	requestedModel := resolveGeminiImageModel("")

	dataURL, _, err := generateGeminiImage(ctx, requestedModel, prompt, nil, aspectRatio)
	if err != nil {
		log.Printf("[magazineReadingImageHandler] Image generation failed for %s/%s: %v", date, readingKey, err)
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to generate reading illustration: %v", err))
		return
	}

	resp := MagazineReadingImageResponse{
		Date:       date,
		ReadingKey: readingKey,
		Image:      dataURL,
		Cached:     false,
	}

	// Optionally persist to Cloudflare R2 if configured
	r2Cfg := ResolveR2Config(nil)
	if r2Cfg.IsConfigured() {
		if imgBytes, mimeType, decErr := decodeDataURL(dataURL); decErr == nil && len(imgBytes) > 0 {
			ext := imageExtensionFromMIME(mimeType)
			objectKey := fmt.Sprintf("magazine/%s/%s.%s", date, readingKey, ext)
			if putErr := r2Cfg.PutObject(ctx, objectKey, imgBytes, mimeType); putErr == nil {
				resp.R2Key = objectKey
				if pubURL := r2Cfg.PublicAssetURL(objectKey); pubURL != "" {
					resp.Image = pubURL
				}
			} else {
				log.Printf("[magazineReadingImageHandler] Warning: failed to upload to R2: %v", putErr)
			}
		}
	}

	setCachedMagazineImage(ctx, resp)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func magazineImagesListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	date := strings.TrimSpace(r.URL.Query().Get("date"))
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	ctx := r.Context()
	keys := []string{"cover", "first_reading", "responsorial_psalm", "second_reading", "gospel"}
	images := make(map[string]MagazineReadingImageResponse)

	for _, k := range keys {
		if cached, ok := getCachedMagazineImage(ctx, date, k); ok {
			images[k] = cached
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"date":   date,
		"images": images,
	})
}
