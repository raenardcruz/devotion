package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func resetMangaTestStores() {
	memorySettings = getDefaultSettings()
	memMangaMu.Lock()
	memMangaStories = nil
	memNextMangaID = 1
	memMangaMu.Unlock()
}

func TestAdminSettingsMaskingIncludesR2Secret(t *testing.T) {
	resetMangaTestStores()
	defer resetMangaTestStores()

	// Save settings with R2 secret
	payload := Settings{
		ContextProvider:   "ollama",
		R2AccountID:       "acct-123",
		R2Bucket:          "devotion-manga",
		R2AccessKeyID:     "key-abc",
		R2SecretAccessKey: "super-secret-r2-value",
		R2PublicURL:       "https://cdn.example.com",
	}
	body, _ := json.Marshal(payload)
	postReq := httptest.NewRequest(http.MethodPost, "/api/admin/settings", bytes.NewReader(body))
	postRec := httptest.NewRecorder()
	adminSettingsHandler(postRec, postReq)

	if postRec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", postRec.Code, postRec.Body.String())
	}

	// GET should mask r2_secret_access_key
	getReq := httptest.NewRequest(http.MethodGet, "/api/admin/settings", nil)
	getRec := httptest.NewRecorder()
	adminSettingsHandler(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on GET, got %d", getRec.Code)
	}

	var got Settings
	if err := json.Unmarshal(getRec.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to decode GET settings: %v", err)
	}
	if !isMaskedKey(got.R2SecretAccessKey) {
		t.Fatalf("expected masked R2SecretAccessKey, got %q", got.R2SecretAccessKey)
	}
	if got.R2Bucket != "devotion-manga" || got.R2AccountID != "acct-123" {
		t.Fatalf("unexpected R2 fields on GET: %+v", got)
	}

	// Re-POSTing the masked secret should preserve original secret in memorySettings
	got.R2Bucket = "devotion-manga-v2"
	reBody, _ := json.Marshal(got)
	reReq := httptest.NewRequest(http.MethodPost, "/api/admin/settings", bytes.NewReader(reBody))
	reRec := httptest.NewRecorder()
	adminSettingsHandler(reRec, reReq)

	if reRec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on re-POST, got %d", reRec.Code)
	}
	saved := GetSettings()
	if saved.R2SecretAccessKey != "super-secret-r2-value" {
		t.Fatalf("expected original R2SecretAccessKey to be preserved, got %q", saved.R2SecretAccessKey)
	}
	if saved.R2Bucket != "devotion-manga-v2" {
		t.Fatalf("expected updated R2Bucket, got %q", saved.R2Bucket)
	}
}

func TestR2StorageAndMangaLifecycle(t *testing.T) {
	resetMangaTestStores()
	defer resetMangaTestStores()

	var mu sync.Mutex
	objects := make(map[string][]byte)
	contentTypes := make(map[string]string)

	mockR2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.Header.Get("Authorization"), "AWS4-HMAC-SHA256 ") {
			http.Error(w, "missing sigv4 header", http.StatusForbidden)
			return
		}
		if r.Header.Get("x-amz-date") == "" || r.Header.Get("x-amz-content-sha256") == "" {
			http.Error(w, "missing amz headers", http.StatusForbidden)
			return
		}

		// ListObjectsV2 test connection
		if r.Method == http.MethodGet && r.URL.Query().Get("list-type") == "2" {
			mu.Lock()
			count := len(objects)
			mu.Unlock()
			w.Header().Set("Content-Type", "application/xml")
			_, _ = w.Write([]byte("<ListBucketResult><KeyCount>" + string(rune('0'+count)) + "</KeyCount></ListBucketResult>"))
			return
		}

		// Expect path: /<bucket>/<key>
		trimmed := strings.TrimPrefix(r.URL.Path, "/test-bucket/")
		if trimmed == r.URL.Path || trimmed == "" {
			http.Error(w, "bad bucket path", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodPut:
			data, _ := io.ReadAll(r.Body)
			mu.Lock()
			objects[trimmed] = data
			contentTypes[trimmed] = r.Header.Get("Content-Type")
			mu.Unlock()
			w.WriteHeader(http.StatusOK)

		case http.MethodGet:
			mu.Lock()
			data, ok := objects[trimmed]
			ct := contentTypes[trimmed]
			mu.Unlock()
			if !ok {
				http.NotFound(w, r)
				return
			}
			w.Header().Set("Content-Type", ct)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(data)

		case http.MethodDelete:
			mu.Lock()
			delete(objects, trimmed)
			delete(contentTypes, trimmed)
			mu.Unlock()
			w.WriteHeader(http.StatusNoContent)

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}))
	defer mockR2.Close()

	memorySettings.R2Endpoint = mockR2.URL
	memorySettings.R2Bucket = "test-bucket"
	memorySettings.R2AccessKeyID = "test-access-key"
	memorySettings.R2SecretAccessKey = "test-secret-key"

	// 1. Test connection endpoint
	testReq := httptest.NewRequest(http.MethodPost, "/api/admin/r2/test", strings.NewReader(`{}`))
	testReq.Header.Set("X-Admin-Token", adminSessionToken)
	testRec := httptest.NewRecorder()
	adminTestR2Handler(testRec, testReq)
	if testRec.Code != http.StatusOK {
		t.Fatalf("expected /api/admin/r2/test 200, got %d: %s", testRec.Code, testRec.Body.String())
	}

	// 2. Save generated manga story (uploads images + manifest.json to R2)
	samplePNG := "data:image/png;base64," + base64.StdEncoding.EncodeToString([]byte("fake-png-bytes"))
	savePayload := MangaSaveStoryRequest{
		Date:           "2026-11-16",
		Style:          "manga",
		ImageModel:     "gemini-nano-banana-2.1",
		GeneratedStory: "Leo learned to forgive after hearing Sirach.",
		MoralLesson:    "Forgive from the heart.",
		StoryAnalysis: MangaStoryAnalysis{
			Title: "Leo's Forgiveness",
			Characters: []MangaCharacter{
				{Name: "Leo", PhysicalDescription: "10-year-old boy", Personality: "Kind", Role: "Protagonist"},
			},
			Setting: MangaSetting{TimePeriod: "Modern day", Location: "Parish", Mood: "Hopeful"},
		},
		CharacterReferences: []MangaCharacterReference{
			{Name: "Leo", Image: samplePNG, Description: "10-year-old boy"},
		},
		StoryBreakdown: MangaStoryBreakdown{
			Panels: []MangaPanelPlan{
				{
					PanelNumber:      1,
					Characters:       []string{"Leo"},
					SceneDescription: "Leo sits in the pew.",
					Narrative:        "In the quiet parish church, Leo listened carefully to the first reading.",
					CameraAngle:      "Medium shot",
					VisualMood:       "Reflective",
				},
			},
		},
		GeneratedPanels: []MangaGeneratedPanel{
			{PanelNumber: 1, Image: samplePNG},
		},
	}

	rec, err := SaveMangaStoryRecord(context.Background(), savePayload)
	if err != nil {
		t.Fatalf("SaveMangaStoryRecord failed: %v", err)
	}
	if !rec.R2Stored || len(rec.R2Keys) != 3 {
		t.Fatalf("expected 3 R2 keys (char, panel, manifest), got R2Stored=%v keys=%v", rec.R2Stored, rec.R2Keys)
	}
	if !strings.HasPrefix(rec.GeneratedPanels[0].Image, "/api/manga/asset?key=") {
		t.Fatalf("expected proxy asset URL when R2PublicURL is empty, got %q", rec.GeneratedPanels[0].Image)
	}

	// 3. Asset handler serves the uploaded panel image & blocks traversal
	assetReq := httptest.NewRequest(http.MethodGet, rec.GeneratedPanels[0].Image, nil)
	assetRec := httptest.NewRecorder()
	mangaAssetHandler(assetRec, assetReq)
	if assetRec.Code != http.StatusOK || assetRec.Body.String() != "fake-png-bytes" {
		t.Fatalf("expected asset handler to serve bytes, got %d %q", assetRec.Code, assetRec.Body.String())
	}

	badAssetReq := httptest.NewRequest(http.MethodGet, "/api/manga/asset?key=../secret", nil)
	badAssetRec := httptest.NewRecorder()
	mangaAssetHandler(badAssetRec, badAssetReq)
	if badAssetRec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid asset key, got %d", badAssetRec.Code)
	}

	// 4. Admin list endpoint returns the story
	listReq := httptest.NewRequest(http.MethodGet, "/api/admin/manga", nil)
	listReq.Header.Set("X-Admin-Token", adminSessionToken)
	listRec := httptest.NewRecorder()
	adminMangaHandler(listRec, listReq)
	if listRec.Code != http.StatusOK {
		t.Fatalf("expected /api/admin/manga 200, got %d", listRec.Code)
	}
	var listResp struct {
		Success bool                `json:"success"`
		Stories []MangaStorySummary `json:"stories"`
	}
	if err := json.Unmarshal(listRec.Body.Bytes(), &listResp); err != nil {
		t.Fatalf("failed decoding list response: %v", err)
	}
	if len(listResp.Stories) != 1 || listResp.Stories[0].Title != "Leo's Forgiveness" {
		t.Fatalf("unexpected list stories: %+v", listResp.Stories)
	}

	// 5. Admin delete removes both DB entry and all R2 objects
	delReq := httptest.NewRequest(http.MethodDelete, "/api/admin/manga?id=1", nil)
	delReq.Header.Set("X-Admin-Token", adminSessionToken)
	delRec := httptest.NewRecorder()
	adminMangaHandler(delRec, delReq)
	if delRec.Code != http.StatusOK {
		t.Fatalf("expected DELETE /api/admin/manga 200, got %d: %s", delRec.Code, delRec.Body.String())
	}

	mu.Lock()
	remainingR2 := len(objects)
	mu.Unlock()
	if remainingR2 != 0 {
		t.Fatalf("expected 0 remaining R2 objects after delete, got %d", remainingR2)
	}

	storiesAfter, _ := ListMangaStories()
	if len(storiesAfter) != 0 {
		t.Fatalf("expected 0 stories in DB/memory after delete, got %d", len(storiesAfter))
	}
}
