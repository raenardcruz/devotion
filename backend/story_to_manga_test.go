package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDefaultGeminiImageModelSetting(t *testing.T) {
	defaults := getDefaultSettings()
	if defaults.GeminiImageModel != "gemini-nano-banana-2.1" {
		t.Fatalf("expected default GeminiImageModel to be 'gemini-nano-banana-2.1', got %q", defaults.GeminiImageModel)
	}
}

func TestAdminGeminiImageModelsHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/admin/gemini-image-models", nil)
	rec := httptest.NewRecorder()

	adminGeminiImageModelsHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var resp struct {
		Models []string `json:"models"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if len(resp.Models) == 0 || resp.Models[0] != "gemini-nano-banana-2.1" {
		t.Fatalf("expected first model to be 'gemini-nano-banana-2.1', got %v", resp.Models)
	}
}

func TestMangaConfigHandler(t *testing.T) {
	memorySettings = getDefaultSettings()

	// GET default config
	getReq := httptest.NewRequest(http.MethodGet, "/api/manga/config", nil)
	getRec := httptest.NewRecorder()
	mangaConfigHandler(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Fatalf("expected GET status 200, got %d", getRec.Code)
	}

	var getResp struct {
		ImageModel string   `json:"imageModel"`
		Models     []string `json:"models"`
	}
	if err := json.Unmarshal(getRec.Body.Bytes(), &getResp); err != nil {
		t.Fatalf("failed to decode GET response: %v", err)
	}
	if getResp.ImageModel != "gemini-nano-banana-2.1" {
		t.Fatalf("expected default imageModel 'gemini-nano-banana-2.1', got %q", getResp.ImageModel)
	}

	// POST updated image model (requires admin)
	postBody, _ := json.Marshal(map[string]string{
		"imageModel": "gemini-2.5-flash-image-preview",
	})
	postReq := httptest.NewRequest(http.MethodPost, "/api/manga/config", bytes.NewReader(postBody))
	postReq.Header.Set("X-Admin-Token", adminSessionToken)
	postRec := httptest.NewRecorder()
	mangaConfigHandler(postRec, postReq)

	if postRec.Code != http.StatusOK {
		t.Fatalf("expected POST status 200, got %d", postRec.Code)
	}
	if GetSettings().GeminiImageModel != "gemini-2.5-flash-image-preview" {
		t.Fatalf("expected saved setting to be 'gemini-2.5-flash-image-preview', got %q", GetSettings().GeminiImageModel)
	}

	// Restore default
	memorySettings.GeminiImageModel = "gemini-nano-banana-2.1"
}

func TestParseMangaJSON(t *testing.T) {
	rawWithFence := "```json\n{\"title\":\"Leo's Forgiveness\",\"characters\":[{\"name\":\"Leo\",\"physicalDescription\":\"10-year-old boy\",\"personality\":\"Kind\",\"role\":\"Protagonist\"}],\"setting\":{\"timePeriod\":\"Modern day\",\"location\":\"School and Parish\",\"mood\":\"Hopeful\"}}\n```"
	var analysis MangaStoryAnalysis
	if err := parseMangaJSON(rawWithFence, &analysis); err != nil {
		t.Fatalf("expected parseMangaJSON to succeed, got error: %v", err)
	}
	if analysis.Title != "Leo's Forgiveness" || len(analysis.Characters) != 1 || analysis.Characters[0].Name != "Leo" {
		t.Fatalf("unexpected parsed analysis: %+v", analysis)
	}
}

func TestBuildReadingsSummary(t *testing.T) {
	input := &MangaReadingsInput{
		FirstReading: &ReadingWithContext{
			Citation: "Sirach 27:30-28:7",
			Text:     "<p>Wrath and anger are hateful things, yet the sinner hugs them tight.</p>",
			Context:  "Forgive your neighbor's injustice; then when you pray, your own sins will be forgiven.",
		},
		Gospel: &ReadingWithContext{
			Citation: "Matthew 18:21-35",
			Text:     "Peter approached Jesus and asked him, Lord, if my brother sins against me, how often must I forgive?",
			Context:  "Jesus teaches that we must forgive seventy-seven times from the heart.",
		},
	}

	summary := buildReadingsSummary(context.Background(), "2026-11-15", input)
	if !strings.Contains(summary, "Sirach 27:30-28:7") || !strings.Contains(summary, "Matthew 18:21-35") {
		t.Fatalf("expected summary to include citations, got: %s", summary)
	}
	if strings.Contains(summary, "<p>") {
		t.Fatalf("expected HTML tags to be stripped, got: %s", summary)
	}
}

func TestMangaAnalyzeStoryWordLimit(t *testing.T) {
	longStory := strings.Repeat("word ", 505)
	payload, _ := json.Marshal(MangaAnalyzeStoryRequest{
		Story: longStory,
		Style: "manga",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/manga/analyze-story", bytes.NewReader(payload))
	req.Header.Set("X-Admin-Token", adminSessionToken)
	rec := httptest.NewRecorder()

	mangaAnalyzeStoryHandler(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request for >500 words, got %d", rec.Code)
	}
}

func TestParseMangaStoryBreakdownWithNarrative(t *testing.T) {
	raw := `{"panels":[{"panelNumber":0,"characters":["Leo"],"sceneDescription":"Leo stares at a torn drawing on his desk.","dialogue":"Why did he do that?","narrative":"All afternoon, anger burned quietly in Leo's chest as he stared at the crumpled sketch his classmate had ruined.","cameraAngle":"Close-up","visualMood":"Tense"},{"panelNumber":2,"characters":["Leo"],"sceneDescription":"Leo sits in the parish pew listening.","narrative":"At evening Mass, the words of Sirach cut through his resentment: forgive your neighbor's injustice, and God will forgive yours.","cameraAngle":"Medium shot","visualMood":"Reflective"}]}`

	var breakdown MangaStoryBreakdown
	if err := parseMangaJSON(raw, &breakdown); err != nil {
		t.Fatalf("expected parseMangaJSON to succeed, got error: %v", err)
	}
	normalizeStoryBreakdown(&breakdown)

	if len(breakdown.Panels) != 2 {
		t.Fatalf("expected 2 panels, got %d", len(breakdown.Panels))
	}
	if breakdown.Panels[0].PanelNumber != 1 {
		t.Fatalf("expected panelNumber 0 to normalize to 1, got %d", breakdown.Panels[0].PanelNumber)
	}
	if !strings.Contains(breakdown.Panels[0].Narrative, "All afternoon, anger burned quietly") {
		t.Fatalf("expected first panel narrative to be preserved, got %q", breakdown.Panels[0].Narrative)
	}
	if !strings.Contains(breakdown.Panels[1].Narrative, "words of Sirach") {
		t.Fatalf("expected second panel narrative to be preserved, got %q", breakdown.Panels[1].Narrative)
	}
}

func TestNormalizeStoryBreakdownFallbackNarrative(t *testing.T) {
	breakdown := MangaStoryBreakdown{
		Panels: []MangaPanelPlan{
			{
				PanelNumber:      0,
				SceneDescription: "  Leo extends his hand in reconciliation after Mass.  ",
				Narrative:        "   ",
			},
		},
	}

	normalizeStoryBreakdown(&breakdown)

	if breakdown.Panels[0].PanelNumber != 1 {
		t.Fatalf("expected PanelNumber 1, got %d", breakdown.Panels[0].PanelNumber)
	}
	if breakdown.Panels[0].Narrative != "Leo extends his hand in reconciliation after Mass." {
		t.Fatalf("expected fallback narrative from trimmed SceneDescription, got %q", breakdown.Panels[0].Narrative)
	}
}

func TestMangaGenerationRequiresAdmin(t *testing.T) {
	endpoints := []struct {
		name    string
		method  string
		target  string
		body    string
		handler http.HandlerFunc
	}{
		{
			name:    "POST /api/manga/config",
			method:  http.MethodPost,
			target:  "/api/manga/config",
			body:    `{"imageModel":"gemini-nano-banana-2.1"}`,
			handler: mangaConfigHandler,
		},
		{
			name:    "POST /api/manga/analyze-story",
			method:  http.MethodPost,
			target:  "/api/manga/analyze-story",
			body:    `{"story":"test"}`,
			handler: mangaAnalyzeStoryHandler,
		},
		{
			name:    "POST /api/manga/generate-character-refs",
			method:  http.MethodPost,
			target:  "/api/manga/generate-character-refs",
			body:    `{"characters":[]}`,
			handler: mangaGenerateCharacterRefsHandler,
		},
		{
			name:    "POST /api/manga/chunk-story",
			method:  http.MethodPost,
			target:  "/api/manga/chunk-story",
			body:    `{"story":"test"}`,
			handler: mangaChunkStoryHandler,
		},
		{
			name:    "POST /api/manga/generate-panel",
			method:  http.MethodPost,
			target:  "/api/manga/generate-panel",
			body:    `{"panel":{"panelNumber":1}}`,
			handler: mangaGeneratePanelHandler,
		},
		{
			name:    "POST /api/manga/story",
			method:  http.MethodPost,
			target:  "/api/manga/story",
			body:    `{"date":"2026-11-16"}`,
			handler: mangaStoryHandler,
		},
		{
			name:    "DELETE /api/manga/story",
			method:  http.MethodDelete,
			target:  "/api/manga/story?id=1",
			body:    "",
			handler: mangaStoryHandler,
		},
	}

	for _, tc := range endpoints {
		req := httptest.NewRequest(tc.method, tc.target, strings.NewReader(tc.body))
		rec := httptest.NewRecorder()
		tc.handler(rec, req)
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("%s without admin token expected 401 Unauthorized, got %d", tc.name, rec.Code)
		}
	}
}
