package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNormalizeMagazineReadingKey(t *testing.T) {
	cases := map[string]string{
		"cover":              "cover",
		"FIRST_READING":      "first_reading",
		"responsorial_psalm": "responsorial_psalm",
		"second_reading":     "second_reading",
		"gospel":             "gospel",
		"unknown":            "reading",
	}

	for input, expected := range cases {
		if got := normalizeMagazineReadingKey(input); got != expected {
			t.Errorf("normalizeMagazineReadingKey(%q) = %q, want %q", input, got, expected)
		}
	}
}

func TestBuildReadingIllustrationPrompt(t *testing.T) {
	prompt := buildReadingIllustrationPrompt("gospel", "Luke 19:1-10", "<p>Zacchaeus climbed a sycamore tree.</p>", "Jesus seeks out the lost.")
	if !strings.Contains(prompt, "Luke 19:1-10") {
		t.Fatalf("expected citation in prompt, got: %s", prompt)
	}
	if !strings.Contains(prompt, "Zacchaeus climbed a sycamore tree.") {
		t.Fatalf("expected stripped text in prompt, got: %s", prompt)
	}
	if !strings.Contains(prompt, "NO text") {
		t.Fatalf("expected no-text instruction in prompt")
	}
}

func TestMagazineImagesCacheAndListHandler(t *testing.T) {
	ctx := context.Background()
	setCachedMagazineImage(ctx, MagazineReadingImageResponse{
		Date:       "2026-11-15",
		ReadingKey: "gospel",
		Image:      "data:image/jpeg;base64,ZmFrZQ==",
	})

	req := httptest.NewRequest(http.MethodGet, "/api/magazine/images?date=2026-11-15", nil)
	rec := httptest.NewRecorder()
	magazineImagesListHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var payload struct {
		Date   string                                  `json:"date"`
		Images map[string]MagazineReadingImageResponse `json:"images"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if payload.Images["gospel"].Image == "" {
		t.Fatalf("expected cached gospel image in response")
	}
}

func TestMagazineReadingImageRequiresAdmin(t *testing.T) {
	body := bytes.NewBufferString(`{"date":"2026-11-16","readingKey":"cover","citation":"Luke 19:1-10"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/magazine/reading-image", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Access-Code", "GUEST-CODE")
	rec := httptest.NewRecorder()

	magazineReadingImageHandler(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized for non-admin, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestMagazineReadingImageAllowsAdminWithCachedImage(t *testing.T) {
	ctx := context.Background()
	setCachedMagazineImage(ctx, MagazineReadingImageResponse{
		Date:       "2026-11-16",
		ReadingKey: "cover",
		Image:      "data:image/jpeg;base64,ZmFrZQ==",
	})

	body := bytes.NewBufferString(`{"date":"2026-11-16","readingKey":"cover","citation":"Luke 19:1-10"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/magazine/reading-image", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Admin-Token", adminSessionToken)
	rec := httptest.NewRecorder()

	magazineReadingImageHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for admin, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp MagazineReadingImageResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if !resp.Cached || resp.Image == "" {
		t.Fatalf("expected cached cover image for admin, got %+v", resp)
	}
}

func TestMagazineImagesListFallbackToR2(t *testing.T) {
	resetMangaTestStores()
	defer resetMangaTestStores()

	memMagazineImagesMu.Lock()
	memMagazineImages = make(map[string]MagazineReadingImageResponse)
	memMagazineImagesMu.Unlock()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Query().Get("list-type") == "2" {
			w.Header().Set("Content-Type", "application/xml")
			_, _ = io.WriteString(w, `<?xml version="1.0" encoding="UTF-8"?>
<ListBucketResult>
  <KeyCount>2</KeyCount>
  <Contents><Key>magazine/2026-11-17/cover.png</Key></Contents>
  <Contents><Key>magazine/2026-11-17/gospel.jpg</Key></Contents>
</ListBucketResult>`)
			return
		}
		if r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/magazine/2026-11-17/cover.png") {
			w.Header().Set("Content-Type", "image/png")
			_, _ = w.Write([]byte("png-bytes"))
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()

	memorySettings.R2Endpoint = srv.URL
	memorySettings.R2Bucket = "test-bucket"
	memorySettings.R2AccessKeyID = "test-key"
	memorySettings.R2SecretAccessKey = "test-secret"
	memorySettings.R2PublicURL = ""

	// Guest fetches previously generated images for date
	req := httptest.NewRequest(http.MethodGet, "/api/magazine/images?date=2026-11-17", nil)
	rec := httptest.NewRecorder()
	magazineImagesListHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", rec.Code, rec.Body.String())
	}

	var payload struct {
		Date   string                                  `json:"date"`
		Images map[string]MagazineReadingImageResponse `json:"images"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if payload.Images["cover"].R2Key != "magazine/2026-11-17/cover.png" {
		t.Fatalf("expected cover R2Key from fallback, got %+v", payload.Images["cover"])
	}
	if payload.Images["gospel"].R2Key != "magazine/2026-11-17/gospel.jpg" {
		t.Fatalf("expected gospel R2Key from fallback, got %+v", payload.Images["gospel"])
	}

	// Verify /api/manga/asset serves magazine/ keys for guests
	assetReq := httptest.NewRequest(http.MethodGet, "/api/manga/asset?key=magazine/2026-11-17/cover.png", nil)
	assetRec := httptest.NewRecorder()
	mangaAssetHandler(assetRec, assetReq)
	if assetRec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for magazine asset, got %d: %s", assetRec.Code, assetRec.Body.String())
	}
}
