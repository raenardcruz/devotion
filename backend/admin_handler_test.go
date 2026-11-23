package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestAdminModelOptionsManagement(t *testing.T) {
	original := memorySettings
	defer func() {
		memorySettings = original
	}()

	memorySettings = getDefaultSettings()

	payload := map[string]interface{}{
		"context_provider":         "gemini",
		"fact_checker_provider":    "gemini",
		"magisterium_llm_provider": "ollama",
		"gemini_model":             "gemini-4.5-pro",
		"gemini_live_model":        "gemini-4-live",
		"gemini_image_model":       "gemini-nano-banana-3",
		"ollama_model":             "qwen3:8b",
		"gemini_models":            []string{" gemini-4.5-pro ", "gemini-4-argon", "gemini-4.5-pro", ""},
		"gemini_live_models":       []string{"gemini-4-live", "gemini-3.8-live"},
		"gemini_image_models":      []string{"gemini-nano-banana-3", "gemini-nano-banana-2.1"},
		"ollama_models":            []string{"qwen3:8b", "gemma4:cloud"},
	}

	body, _ := json.Marshal(payload)
	postReq := httptest.NewRequest(http.MethodPost, "/api/admin/settings", bytes.NewReader(body))
	postRec := httptest.NewRecorder()
	adminSettingsHandler(postRec, postReq)

	if postRec.Code != http.StatusOK {
		t.Fatalf("expected POST /api/admin/settings 200, got %d: %s", postRec.Code, postRec.Body.String())
	}

	// Verify GET /api/admin/gemini-models
	gmReq := httptest.NewRequest(http.MethodGet, "/api/admin/gemini-models", nil)
	gmRec := httptest.NewRecorder()
	adminGeminiModelsHandler(gmRec, gmReq)
	var gmResp struct {
		Models []string `json:"models"`
	}
	if err := json.Unmarshal(gmRec.Body.Bytes(), &gmResp); err != nil {
		t.Fatalf("failed to decode gemini-models: %v", err)
	}
	expectedGemini := []string{"gemini-4.5-pro", "gemini-4-argon"}
	if !reflect.DeepEqual(gmResp.Models, expectedGemini) {
		t.Fatalf("expected gemini models %v, got %v", expectedGemini, gmResp.Models)
	}

	// Verify GET /api/admin/gemini-live-models
	glReq := httptest.NewRequest(http.MethodGet, "/api/admin/gemini-live-models", nil)
	glRec := httptest.NewRecorder()
	adminGeminiLiveModelsHandler(glRec, glReq)
	var glResp struct {
		Models []string `json:"models"`
	}
	if err := json.Unmarshal(glRec.Body.Bytes(), &glResp); err != nil {
		t.Fatalf("failed to decode gemini-live-models: %v", err)
	}
	expectedLive := []string{"gemini-4-live", "gemini-3.8-live"}
	if !reflect.DeepEqual(glResp.Models, expectedLive) {
		t.Fatalf("expected gemini live models %v, got %v", expectedLive, glResp.Models)
	}

	// Verify GET /api/admin/gemini-image-models and GET /api/manga/config
	giReq := httptest.NewRequest(http.MethodGet, "/api/admin/gemini-image-models", nil)
	giRec := httptest.NewRecorder()
	adminGeminiImageModelsHandler(giRec, giReq)
	var giResp struct {
		Models []string `json:"models"`
	}
	if err := json.Unmarshal(giRec.Body.Bytes(), &giResp); err != nil {
		t.Fatalf("failed to decode gemini-image-models: %v", err)
	}
	expectedImage := []string{"gemini-nano-banana-3", "gemini-nano-banana-2.1"}
	if !reflect.DeepEqual(giResp.Models, expectedImage) {
		t.Fatalf("expected gemini image models %v, got %v", expectedImage, giResp.Models)
	}

	mcReq := httptest.NewRequest(http.MethodGet, "/api/manga/config", nil)
	mcRec := httptest.NewRecorder()
	mangaConfigHandler(mcRec, mcReq)
	var mcResp struct {
		ImageModel string   `json:"imageModel"`
		Models     []string `json:"models"`
	}
	if err := json.Unmarshal(mcRec.Body.Bytes(), &mcResp); err != nil {
		t.Fatalf("failed to decode manga config: %v", err)
	}
	if mcResp.ImageModel != "gemini-nano-banana-3" {
		t.Fatalf("expected manga config imageModel 'gemini-nano-banana-3', got %q", mcResp.ImageModel)
	}
	if !reflect.DeepEqual(mcResp.Models, expectedImage) {
		t.Fatalf("expected manga config models %v, got %v", expectedImage, mcResp.Models)
	}
}

func TestNormalizeModelListFallback(t *testing.T) {
	got := normalizeModelList([]string{"", "   "}, DefaultGeminiModels)
	if !reflect.DeepEqual(got, DefaultGeminiModels) {
		t.Fatalf("expected fallback to DefaultGeminiModels %v, got %v", DefaultGeminiModels, got)
	}
}
