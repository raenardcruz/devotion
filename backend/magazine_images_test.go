package main

import (
	"context"
	"encoding/json"
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
