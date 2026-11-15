package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type MangaStoryRecord struct {
	ID                  int                       `json:"id"`
	Date                string                    `json:"date"`
	Style               string                    `json:"style"`
	Title               string                    `json:"title"`
	MoralLesson         string                    `json:"moralLesson"`
	Story               string                    `json:"story"`
	ImageModel          string                    `json:"imageModel"`
	StoryAnalysis       MangaStoryAnalysis        `json:"storyAnalysis"`
	CharacterReferences []MangaCharacterReference `json:"characterReferences"`
	StoryBreakdown      MangaStoryBreakdown       `json:"storyBreakdown"`
	GeneratedPanels     []MangaGeneratedPanel     `json:"generatedPanels"`
	R2Prefix            string                    `json:"r2Prefix"`
	R2Keys              []string                  `json:"r2Keys"`
	R2Stored            bool                      `json:"r2Stored"`
	CreatedAt           string                    `json:"createdAt"`
	UpdatedAt           string                    `json:"updatedAt"`
}

type MangaStorySummary struct {
	ID             int      `json:"id"`
	Date           string   `json:"date"`
	Style          string   `json:"style"`
	Title          string   `json:"title"`
	MoralLesson    string   `json:"moralLesson"`
	ImageModel     string   `json:"imageModel"`
	PanelCount     int      `json:"panelCount"`
	CharacterCount int      `json:"characterCount"`
	R2Prefix       string   `json:"r2Prefix"`
	R2Keys         []string `json:"r2Keys"`
	R2Stored       bool     `json:"r2Stored"`
	PreviewImage   string   `json:"previewImage,omitempty"`
	CreatedAt      string   `json:"createdAt"`
	UpdatedAt      string   `json:"updatedAt"`
}

type MangaSaveStoryRequest struct {
	Date                string                    `json:"date"`
	Style               string                    `json:"style"`
	ImageModel          string                    `json:"imageModel"`
	GeneratedStory      string                    `json:"generatedStory"`
	MoralLesson         string                    `json:"moralLesson"`
	StoryAnalysis       MangaStoryAnalysis        `json:"storyAnalysis"`
	CharacterReferences []MangaCharacterReference `json:"characterReferences"`
	StoryBreakdown      MangaStoryBreakdown       `json:"storyBreakdown"`
	GeneratedPanels     []MangaGeneratedPanel     `json:"generatedPanels"`
}

var (
	memMangaMu      sync.Mutex
	memMangaStories []MangaStoryRecord
	memNextMangaID  = 1
)

func imageExtensionFromMIME(mimeType string) string {
	lower := strings.ToLower(strings.TrimSpace(mimeType))
	switch {
	case strings.Contains(lower, "png"):
		return "png"
	case strings.Contains(lower, "webp"):
		return "webp"
	default:
		return "jpg"
	}
}

func normalizeMangaStyle(style string) string {
	s := strings.ToLower(strings.TrimSpace(style))
	if s != "comic" {
		return "manga"
	}
	return s
}

func hydrateRecordAssetURLs(rec *MangaStoryRecord) {
	if rec == nil {
		return
	}
	r2Cfg := ResolveR2Config(nil)
	if !r2Cfg.IsConfigured() {
		return
	}
	for i := range rec.CharacterReferences {
		if rec.CharacterReferences[i].R2Key != "" {
			rec.CharacterReferences[i].Image = r2Cfg.PublicAssetURL(rec.CharacterReferences[i].R2Key)
		}
	}
	for i := range rec.GeneratedPanels {
		if rec.GeneratedPanels[i].R2Key != "" {
			rec.GeneratedPanels[i].Image = r2Cfg.PublicAssetURL(rec.GeneratedPanels[i].R2Key)
		}
	}
}

func SaveMangaStoryRecord(ctx context.Context, req MangaSaveStoryRequest) (*MangaStoryRecord, error) {
	date := strings.TrimSpace(req.Date)
	if date == "" {
		return nil, fmt.Errorf("date is required")
	}
	style := normalizeMangaStyle(req.Style)
	imageModel := resolveGeminiImageModel(req.ImageModel)
	storyText := strings.TrimSpace(req.GeneratedStory)
	if storyText == "" {
		storyText = strings.TrimSpace(req.StoryAnalysis.Story)
	}
	moralLesson := strings.TrimSpace(req.MoralLesson)
	if moralLesson == "" {
		moralLesson = strings.TrimSpace(req.StoryAnalysis.MoralLesson)
	}
	title := strings.TrimSpace(req.StoryAnalysis.Title)
	if title == "" {
		title = fmt.Sprintf("Daily Story (%s)", date)
	}

	normalizeStoryBreakdown(&req.StoryBreakdown)

	existing, _ := GetMangaStoryByDateAndStyle(date, style)
	var oldR2Keys []string
	if existing != nil {
		oldR2Keys = existing.R2Keys
	}

	r2Cfg := ResolveR2Config(nil)
	prefix := fmt.Sprintf("stories/%s/%s", date, style)
	stamp := fmt.Sprintf("%d", time.Now().UnixMilli())

	charRefs := make([]MangaCharacterReference, len(req.CharacterReferences))
	copy(charRefs, req.CharacterReferences)

	panels := make([]MangaGeneratedPanel, len(req.GeneratedPanels))
	copy(panels, req.GeneratedPanels)

	var newR2Keys []string
	r2Stored := false

	if r2Cfg.IsConfigured() {
		for i := range charRefs {
			imgVal := strings.TrimSpace(charRefs[i].Image)
			if strings.HasPrefix(imgVal, "data:") {
				rawBytes, mimeType, err := decodeDataURL(imgVal)
				if err != nil {
					return nil, fmt.Errorf("failed to decode character %d image: %w", i, err)
				}
				ext := imageExtensionFromMIME(mimeType)
				key := fmt.Sprintf("%s/characters/char_%d_%s.%s", prefix, i, stamp, ext)
				if err := r2Cfg.PutObject(ctx, key, rawBytes, mimeType); err != nil {
					return nil, fmt.Errorf("failed to upload character %d image to R2: %w", i, err)
				}
				charRefs[i].R2Key = key
				charRefs[i].Image = r2Cfg.PublicAssetURL(key)
				newR2Keys = append(newR2Keys, key)
			} else if charRefs[i].R2Key != "" {
				newR2Keys = append(newR2Keys, charRefs[i].R2Key)
				charRefs[i].Image = r2Cfg.PublicAssetURL(charRefs[i].R2Key)
			}
		}

		for i := range panels {
			imgVal := strings.TrimSpace(panels[i].Image)
			if strings.HasPrefix(imgVal, "data:") {
				rawBytes, mimeType, err := decodeDataURL(imgVal)
				if err != nil {
					return nil, fmt.Errorf("failed to decode panel %d image: %w", panels[i].PanelNumber, err)
				}
				ext := imageExtensionFromMIME(mimeType)
				key := fmt.Sprintf("%s/panels/panel_%d_%s.%s", prefix, panels[i].PanelNumber, stamp, ext)
				if err := r2Cfg.PutObject(ctx, key, rawBytes, mimeType); err != nil {
					return nil, fmt.Errorf("failed to upload panel %d image to R2: %w", panels[i].PanelNumber, err)
				}
				panels[i].R2Key = key
				panels[i].Image = r2Cfg.PublicAssetURL(key)
				newR2Keys = append(newR2Keys, key)
			} else if panels[i].R2Key != "" {
				newR2Keys = append(newR2Keys, panels[i].R2Key)
				panels[i].Image = r2Cfg.PublicAssetURL(panels[i].R2Key)
			}
		}

		manifestKey := fmt.Sprintf("%s/manifest.json", prefix)
		manifestPayload := map[string]interface{}{
			"date":                date,
			"style":               style,
			"title":               title,
			"moralLesson":         moralLesson,
			"story":               storyText,
			"imageModel":          imageModel,
			"storyAnalysis":       req.StoryAnalysis,
			"characterReferences": charRefs,
			"storyBreakdown":      req.StoryBreakdown,
			"generatedPanels":     panels,
			"updatedAt":           time.Now().UTC().Format(time.RFC3339),
		}
		if manifestBytes, err := json.MarshalIndent(manifestPayload, "", "  "); err == nil {
			if err := r2Cfg.PutObject(ctx, manifestKey, manifestBytes, "application/json"); err != nil {
				return nil, fmt.Errorf("failed to upload manifest.json to R2: %w", err)
			}
			newR2Keys = append(newR2Keys, manifestKey)
		}

		r2Stored = len(newR2Keys) > 0

		// Clean up obsolete R2 keys from an earlier generation
		newKeySet := make(map[string]bool, len(newR2Keys))
		for _, k := range newR2Keys {
			newKeySet[k] = true
		}
		var staleKeys []string
		for _, oldKey := range oldR2Keys {
			if oldKey != "" && !newKeySet[oldKey] {
				staleKeys = append(staleKeys, oldKey)
			}
		}
		if len(staleKeys) > 0 {
			if err := r2Cfg.DeleteObjects(ctx, staleKeys); err != nil {
				log.Printf("[SaveMangaStoryRecord] Warning cleaning stale R2 keys: %v", err)
			}
		}
	}

	nowStr := time.Now().UTC().Format(time.RFC3339)
	rec := &MangaStoryRecord{
		Date:                date,
		Style:               style,
		Title:               title,
		MoralLesson:         moralLesson,
		Story:               storyText,
		ImageModel:          imageModel,
		StoryAnalysis:       req.StoryAnalysis,
		CharacterReferences: charRefs,
		StoryBreakdown:      req.StoryBreakdown,
		GeneratedPanels:     panels,
		R2Prefix:            prefix,
		R2Keys:              uniqueStrings(newR2Keys),
		R2Stored:            r2Stored,
		CreatedAt:           nowStr,
		UpdatedAt:           nowStr,
	}

	if dbAvailable() {
		analysisJSON, _ := json.Marshal(rec.StoryAnalysis)
		charRefsJSON, _ := json.Marshal(rec.CharacterReferences)
		breakdownJSON, _ := json.Marshal(rec.StoryBreakdown)
		panelsJSON, _ := json.Marshal(rec.GeneratedPanels)
		r2KeysJSON, _ := json.Marshal(rec.R2Keys)

		query := `
		INSERT INTO generated_manga_stories (
			story_date, style, title, moral_lesson, story_text, image_model,
			analysis_json, character_refs_json, breakdown_json, panels_json,
			r2_prefix, r2_keys_json, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT (story_date, style) DO UPDATE SET
			title = EXCLUDED.title,
			moral_lesson = EXCLUDED.moral_lesson,
			story_text = EXCLUDED.story_text,
			image_model = EXCLUDED.image_model,
			analysis_json = EXCLUDED.analysis_json,
			character_refs_json = EXCLUDED.character_refs_json,
			breakdown_json = EXCLUDED.breakdown_json,
			panels_json = EXCLUDED.panels_json,
			r2_prefix = EXCLUDED.r2_prefix,
			r2_keys_json = EXCLUDED.r2_keys_json,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id, created_at, updated_at;`

		var createdAt, updatedAt time.Time
		err := db.QueryRow(
			query,
			rec.Date,
			rec.Style,
			rec.Title,
			rec.MoralLesson,
			rec.Story,
			rec.ImageModel,
			string(analysisJSON),
			string(charRefsJSON),
			string(breakdownJSON),
			string(panelsJSON),
			rec.R2Prefix,
			string(r2KeysJSON),
		).Scan(&rec.ID, &createdAt, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed saving generated manga story to database: %w", err)
		}
		rec.CreatedAt = createdAt.Format(time.RFC3339)
		rec.UpdatedAt = updatedAt.Format(time.RFC3339)
		return rec, nil
	}

	// In-memory fallback
	memMangaMu.Lock()
	defer memMangaMu.Unlock()
	for i := range memMangaStories {
		if memMangaStories[i].Date == rec.Date && memMangaStories[i].Style == rec.Style {
			rec.ID = memMangaStories[i].ID
			rec.CreatedAt = memMangaStories[i].CreatedAt
			memMangaStories[i] = *rec
			return rec, nil
		}
	}
	rec.ID = memNextMangaID
	memNextMangaID++
	memMangaStories = append([]MangaStoryRecord{*rec}, memMangaStories...)
	return rec, nil
}

func scanMangaStoryRow(scanner interface {
	Scan(dest ...interface{}) error
}) (*MangaStoryRecord, error) {
	var rec MangaStoryRecord
	var analysisRaw, charRefsRaw, breakdownRaw, panelsRaw, r2KeysRaw string
	var createdAt, updatedAt time.Time

	err := scanner.Scan(
		&rec.ID,
		&rec.Date,
		&rec.Style,
		&rec.Title,
		&rec.MoralLesson,
		&rec.Story,
		&rec.ImageModel,
		&analysisRaw,
		&charRefsRaw,
		&breakdownRaw,
		&panelsRaw,
		&rec.R2Prefix,
		&r2KeysRaw,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	_ = json.Unmarshal([]byte(analysisRaw), &rec.StoryAnalysis)
	_ = json.Unmarshal([]byte(charRefsRaw), &rec.CharacterReferences)
	_ = json.Unmarshal([]byte(breakdownRaw), &rec.StoryBreakdown)
	_ = json.Unmarshal([]byte(panelsRaw), &rec.GeneratedPanels)
	_ = json.Unmarshal([]byte(r2KeysRaw), &rec.R2Keys)
	rec.R2Stored = len(rec.R2Keys) > 0
	rec.CreatedAt = createdAt.Format(time.RFC3339)
	rec.UpdatedAt = updatedAt.Format(time.RFC3339)
	hydrateRecordAssetURLs(&rec)
	return &rec, nil
}

func GetMangaStoryByDateAndStyle(date, style string) (*MangaStoryRecord, error) {
	cleanDate := strings.TrimSpace(date)
	cleanStyle := strings.ToLower(strings.TrimSpace(style))

	if dbAvailable() {
		var row *sql.Row
		if cleanStyle != "" {
			row = db.QueryRow(`
				SELECT id, story_date, style, title, moral_lesson, story_text, image_model,
				       analysis_json, character_refs_json, breakdown_json, panels_json,
				       r2_prefix, r2_keys_json, created_at, updated_at
				FROM generated_manga_stories
				WHERE story_date = $1 AND style = $2
				LIMIT 1`, cleanDate, normalizeMangaStyle(cleanStyle))
		} else {
			row = db.QueryRow(`
				SELECT id, story_date, style, title, moral_lesson, story_text, image_model,
				       analysis_json, character_refs_json, breakdown_json, panels_json,
				       r2_prefix, r2_keys_json, created_at, updated_at
				FROM generated_manga_stories
				WHERE story_date = $1
				ORDER BY updated_at DESC
				LIMIT 1`, cleanDate)
		}
		rec, err := scanMangaStoryRow(row)
		if err == sql.ErrNoRows {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}
		return rec, nil
	}

	memMangaMu.Lock()
	defer memMangaMu.Unlock()
	for i := range memMangaStories {
		if memMangaStories[i].Date == cleanDate && (cleanStyle == "" || memMangaStories[i].Style == normalizeMangaStyle(cleanStyle)) {
			copied := memMangaStories[i]
			hydrateRecordAssetURLs(&copied)
			return &copied, nil
		}
	}
	return nil, nil
}

func GetMangaStoryByID(id int) (*MangaStoryRecord, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid story id")
	}

	if dbAvailable() {
		row := db.QueryRow(`
			SELECT id, story_date, style, title, moral_lesson, story_text, image_model,
			       analysis_json, character_refs_json, breakdown_json, panels_json,
			       r2_prefix, r2_keys_json, created_at, updated_at
			FROM generated_manga_stories
			WHERE id = $1
			LIMIT 1`, id)
		rec, err := scanMangaStoryRow(row)
		if err == sql.ErrNoRows {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}
		return rec, nil
	}

	memMangaMu.Lock()
	defer memMangaMu.Unlock()
	for i := range memMangaStories {
		if memMangaStories[i].ID == id {
			copied := memMangaStories[i]
			hydrateRecordAssetURLs(&copied)
			return &copied, nil
		}
	}
	return nil, nil
}

func recordToSummary(rec *MangaStoryRecord) MangaStorySummary {
	preview := ""
	for _, p := range rec.GeneratedPanels {
		if p.Image != "" && !strings.HasPrefix(p.Image, "data:") {
			preview = p.Image
			break
		}
	}
	if preview == "" && len(rec.GeneratedPanels) > 0 {
		preview = rec.GeneratedPanels[0].Image
	}

	panelCount := len(rec.StoryBreakdown.Panels)
	if len(rec.GeneratedPanels) > panelCount {
		panelCount = len(rec.GeneratedPanels)
	}

	return MangaStorySummary{
		ID:             rec.ID,
		Date:           rec.Date,
		Style:          rec.Style,
		Title:          rec.Title,
		MoralLesson:    rec.MoralLesson,
		ImageModel:     rec.ImageModel,
		PanelCount:     panelCount,
		CharacterCount: len(rec.CharacterReferences),
		R2Prefix:       rec.R2Prefix,
		R2Keys:         rec.R2Keys,
		R2Stored:       rec.R2Stored,
		PreviewImage:   preview,
		CreatedAt:      rec.CreatedAt,
		UpdatedAt:      rec.UpdatedAt,
	}
}

func ListMangaStories() ([]MangaStorySummary, error) {
	if dbAvailable() {
		rows, err := db.Query(`
			SELECT id, story_date, style, title, moral_lesson, story_text, image_model,
			       analysis_json, character_refs_json, breakdown_json, panels_json,
			       r2_prefix, r2_keys_json, created_at, updated_at
			FROM generated_manga_stories
			ORDER BY story_date DESC, updated_at DESC`)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var list []MangaStorySummary
		for rows.Next() {
			rec, scanErr := scanMangaStoryRow(rows)
			if scanErr != nil {
				continue
			}
			list = append(list, recordToSummary(rec))
		}
		return list, nil
	}

	memMangaMu.Lock()
	defer memMangaMu.Unlock()
	list := make([]MangaStorySummary, 0, len(memMangaStories))
	for i := range memMangaStories {
		copied := memMangaStories[i]
		hydrateRecordAssetURLs(&copied)
		list = append(list, recordToSummary(&copied))
	}
	return list, nil
}

func DeleteMangaStory(ctx context.Context, id int, date, style string) error {
	var target *MangaStoryRecord
	var err error
	if id > 0 {
		target, err = GetMangaStoryByID(id)
	} else {
		target, err = GetMangaStoryByDateAndStyle(date, style)
	}
	if err != nil {
		return err
	}
	if target == nil {
		return fmt.Errorf("generated manga story not found")
	}

	// Delete objects from R2 if any keys are recorded and R2 is configured
	if len(target.R2Keys) > 0 {
		r2Cfg := ResolveR2Config(nil)
		if r2Cfg.IsConfigured() {
			if err := r2Cfg.DeleteObjects(ctx, target.R2Keys); err != nil {
				return fmt.Errorf("failed deleting objects from R2: %w", err)
			}
		}
	}

	if dbAvailable() {
		_, err := db.Exec(`DELETE FROM generated_manga_stories WHERE id = $1`, target.ID)
		if err != nil {
			return fmt.Errorf("failed deleting story from database: %w", err)
		}
		return nil
	}

	memMangaMu.Lock()
	defer memMangaMu.Unlock()
	filtered := make([]MangaStoryRecord, 0, len(memMangaStories))
	for _, item := range memMangaStories {
		if item.ID != target.ID {
			filtered = append(filtered, item)
		}
	}
	memMangaStories = filtered
	return nil
}

// === HTTP Handlers ===

func mangaStoryHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		idStr := strings.TrimSpace(r.URL.Query().Get("id"))
		date := strings.TrimSpace(r.URL.Query().Get("date"))
		style := strings.TrimSpace(r.URL.Query().Get("style"))

		var rec *MangaStoryRecord
		var err error
		if idStr != "" {
			id, convErr := strconv.Atoi(idStr)
			if convErr != nil {
				respondWithError(w, http.StatusBadRequest, "Invalid id parameter")
				return
			}
			rec, err = GetMangaStoryByID(id)
		} else {
			if date == "" {
				respondWithError(w, http.StatusBadRequest, "date or id parameter is required")
				return
			}
			rec, err = GetMangaStoryByDateAndStyle(date, style)
		}

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to load generated manga story: %v", err))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if rec == nil {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"found":   false,
			})
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"found":   true,
			"story":   rec,
		})

	case http.MethodPost, http.MethodPut:
		if !isAdminRequest(r) {
			respondWithError(w, http.StatusUnauthorized, "Admin authentication required")
			return
		}
		var req MangaSaveStoryRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		rec, err := SaveMangaStoryRecord(r.Context(), req)
		if err != nil {
			log.Printf("[mangaStoryHandler] Save error: %v", err)
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"story":   rec,
		})

	case http.MethodDelete:
		if !isAdminRequest(r) {
			respondWithError(w, http.StatusUnauthorized, "Admin authentication required")
			return
		}
		idStr := strings.TrimSpace(r.URL.Query().Get("id"))
		date := strings.TrimSpace(r.URL.Query().Get("date"))
		style := strings.TrimSpace(r.URL.Query().Get("style"))
		id := 0
		if idStr != "" {
			parsedID, convErr := strconv.Atoi(idStr)
			if convErr != nil {
				respondWithError(w, http.StatusBadRequest, "Invalid id parameter")
				return
			}
			id = parsedID
		}
		if id == 0 && date == "" {
			respondWithError(w, http.StatusBadRequest, "id or date parameter is required")
			return
		}
		if err := DeleteMangaStory(r.Context(), id, date, style); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func mangaAssetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	key := strings.TrimSpace(r.URL.Query().Get("key"))
	if key == "" || !strings.HasPrefix(key, "stories/") || strings.Contains(key, "..") {
		respondWithError(w, http.StatusBadRequest, "Invalid asset key")
		return
	}

	r2Cfg := ResolveR2Config(nil)
	if !r2Cfg.IsConfigured() {
		respondWithError(w, http.StatusServiceUnavailable, "R2 storage is not configured")
		return
	}

	data, contentType, err := r2Cfg.GetObject(r.Context(), key)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondWithError(w, http.StatusNotFound, "Asset not found")
			return
		}
		respondWithError(w, http.StatusBadGateway, fmt.Sprintf("Failed to load asset from R2: %v", err))
		return
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func adminTestR2Handler(w http.ResponseWriter, r *http.Request) {
	if !isAdminRequest(r) {
		respondWithError(w, http.StatusUnauthorized, "Admin authentication required")
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var override R2Config
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&override)
	}
	cfg := ResolveR2Config(&override)
	count, err := cfg.TestConnection(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":     true,
		"message":     fmt.Sprintf("Connected to R2 bucket %q successfully.", cfg.Bucket),
		"bucket":      cfg.Bucket,
		"endpoint":    cfg.ResolvedEndpoint(),
		"sampleCount": count,
	})
}

func adminMangaHandler(w http.ResponseWriter, r *http.Request) {
	if !isAdminRequest(r) {
		respondWithError(w, http.StatusUnauthorized, "Admin authentication required")
		return
	}

	switch r.Method {
	case http.MethodGet:
		stories, err := ListMangaStories()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to list generated manga: %v", err))
			return
		}
		if stories == nil {
			stories = []MangaStorySummary{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":      true,
			"r2Configured": ResolveR2Config(nil).IsConfigured(),
			"stories":      stories,
		})

	case http.MethodDelete:
		if !isAdminRequest(r) {
			respondWithError(w, http.StatusUnauthorized, "Admin authentication required")
			return
		}
		idStr := strings.TrimSpace(r.URL.Query().Get("id"))
		date := strings.TrimSpace(r.URL.Query().Get("date"))
		style := strings.TrimSpace(r.URL.Query().Get("style"))
		id := 0
		if idStr != "" {
			parsedID, convErr := strconv.Atoi(idStr)
			if convErr != nil {
				respondWithError(w, http.StatusBadRequest, "Invalid id parameter")
				return
			}
			id = parsedID
		}
		if id == 0 && date == "" {
			respondWithError(w, http.StatusBadRequest, "id or date parameter is required")
			return
		}
		if err := DeleteMangaStory(r.Context(), id, date, style); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Deleted generated manga from database and R2.",
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func adminRegenerateMangaHandler(w http.ResponseWriter, r *http.Request) {
	if !isAdminRequest(r) {
		respondWithError(w, http.StatusUnauthorized, "Admin authentication required")
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ID         int    `json:"id"`
		Date       string `json:"date"`
		Style      string `json:"style"`
		ImageModel string `json:"imageModel"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	date := strings.TrimSpace(req.Date)
	style := strings.TrimSpace(req.Style)
	imageModel := strings.TrimSpace(req.ImageModel)

	if req.ID > 0 {
		existing, err := GetMangaStoryByID(req.ID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to load story record: %v", err))
			return
		}
		if existing == nil {
			respondWithError(w, http.StatusNotFound, "Generated manga story not found")
			return
		}
		if date == "" {
			date = existing.Date
		}
		if style == "" {
			style = existing.Style
		}
		if imageModel == "" {
			imageModel = existing.ImageModel
		}
	}

	if date == "" {
		respondWithError(w, http.StatusBadRequest, "id or date is required to regenerate")
		return
	}
	style = normalizeMangaStyle(style)
	imageModel = resolveGeminiImageModel(imageModel)

	ctx := r.Context()

	// Step 1: Story & character/setting analysis
	analysis, _, err := generateStoryFromReadingsInternal(ctx, MangaAnalyzeStoryRequest{
		Date:  date,
		Style: style,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Step 1 (story) failed: %v", err))
		return
	}

	// Step 2: Character reference sheets
	charRefs, usedModel, err := generateCharacterRefsInternal(ctx, MangaGenerateCharacterRefsRequest{
		Characters: analysis.Characters,
		Setting:    analysis.Setting,
		Style:      style,
		ImageModel: imageModel,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Step 2 (character sheets) failed: %v", err))
		return
	}

	// Step 3: Panel breakdown + narratives
	breakdown, err := chunkStoryInternal(ctx, MangaChunkStoryRequest{
		Story:      analysis.Story,
		Characters: analysis.Characters,
		Setting:    analysis.Setting,
		Style:      style,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Step 3 (panel plan) failed: %v", err))
		return
	}

	// Step 4: Panel illustrations
	var generatedPanels []MangaGeneratedPanel
	for _, panelPlan := range breakdown.Panels {
		panelImg, modelUsed, panelErr := generatePanelInternal(ctx, MangaGeneratePanelRequest{
			Panel:               panelPlan,
			CharacterReferences: charRefs,
			Setting:             analysis.Setting,
			Style:               style,
			ImageModel:          usedModel,
		})
		if panelErr != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Step 4 (panel %d) failed: %v", panelPlan.PanelNumber, panelErr))
			return
		}
		usedModel = modelUsed
		generatedPanels = append(generatedPanels, *panelImg)
	}

	// Save & upload to R2 + DB (automatically cleans up stale R2 keys)
	saved, err := SaveMangaStoryRecord(ctx, MangaSaveStoryRequest{
		Date:                date,
		Style:               style,
		ImageModel:          usedModel,
		GeneratedStory:      analysis.Story,
		MoralLesson:         analysis.MoralLesson,
		StoryAnalysis:       *analysis,
		CharacterReferences: charRefs,
		StoryBreakdown:      *breakdown,
		GeneratedPanels:     generatedPanels,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed saving regenerated manga: %v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"story":   saved,
		"summary": recordToSummary(saved),
	})
}
