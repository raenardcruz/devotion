package main

import (
	"bytes"
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func resetMemAIUsageForTest() {
	memAIUsageMu.Lock()
	memAIUsageLogs = nil
	memAIUsageNextID = 1
	memAIUsageMu.Unlock()
}

func TestEstimateAICostUSD(t *testing.T) {
	// Gemini Flash-Lite: $0.075 / 1M input, $0.30 / 1M output
	cost := EstimateAICostUSD("gemini", "gemini-3.1-flash-lite", "generate", 1_000_000, 1_000_000, 0, 0, 0)
	if math.Abs(cost-0.375) > 1e-6 {
		t.Fatalf("expected 0.375 for 1M in + 1M out on flash-lite, got %f", cost)
	}

	// Gemini Pro: $1.25 / 1M input, $10.00 / 1M output
	proCost := EstimateAICostUSD("gemini", "gemini-2.5-pro", "generate", 1_000_000, 1_000_000, 0, 0, 0)
	if math.Abs(proCost-11.25) > 1e-6 {
		t.Fatalf("expected 11.25 for 1M in + 1M out on pro, got %f", proCost)
	}

	// Gemini Image: $0.039 per image
	imgCost := EstimateAICostUSD("gemini", "gemini-nano-banana-2.1", "generate_image", 0, 0, 0, 0, 2)
	if math.Abs(imgCost-0.078) > 1e-6 {
		t.Fatalf("expected 0.078 for 2 images, got %f", imgCost)
	}

	// Ollama & Magisterium: $0
	if ollamaCost := EstimateAICostUSD("ollama", "gemma4:cloud", "generate", 50000, 10000, 0, 0, 0); ollamaCost != 0 {
		t.Fatalf("expected 0 for ollama, got %f", ollamaCost)
	}
	if magCost := EstimateAICostUSD("magisterium", "magisterium-search-v1", "search", 0, 0, 0, 0, 0); magCost != 0 {
		t.Fatalf("expected 0 for magisterium, got %f", magCost)
	}
}

func TestQueryAIUsageReportFiltersAndPurge(t *testing.T) {
	origDB := db
	db = nil
	defer func() { db = origDB }()
	resetMemAIUsageForTest()

	loc := time.Now().Location()
	// Anchor: Wednesday 2026-11-18
	RecordAIUsage(AIUsageRecord{
		Provider:     "gemini",
		Model:        "gemini-3.1-flash-lite",
		Feature:      "bible_context",
		Operation:    "generate",
		InputTokens:  10000,
		OutputTokens: 2000,
		CreatedAt:    time.Date(2026, 11, 18, 10, 0, 0, 0, loc), // same day
	})
	RecordAIUsage(AIUsageRecord{
		Provider:     "gemini",
		Model:        "gemini-nano-banana-2.1",
		Feature:      "image_generation",
		Operation:    "generate_image",
		ImageCount:   1,
		CreatedAt:    time.Date(2026, 11, 16, 14, 0, 0, 0, loc), // same week (Monday)
	})
	RecordAIUsage(AIUsageRecord{
		Provider:     "ollama",
		Model:        "gemma4:cloud",
		Feature:      "magisterium_chat",
		Operation:    "summarize",
		InputTokens:  5000,
		OutputTokens: 1000,
		CreatedAt:    time.Date(2026, 11, 2, 9, 0, 0, 0, loc), // same month, earlier week
	})
	RecordAIUsage(AIUsageRecord{
		Provider:     "gemini",
		Model:        "gemini-2.5-pro",
		Feature:      "story_to_manga",
		Operation:    "structured_json",
		InputTokens:  4000,
		OutputTokens: 1000,
		CreatedAt:    time.Date(2026, 5, 10, 12, 0, 0, 0, loc), // same year, earlier month
	})
	RecordAIUsage(AIUsageRecord{
		Provider:     "gemini",
		Model:        "gemini-3.1-flash-lite",
		Feature:      "bible_context",
		Operation:    "generate",
		InputTokens:  2000,
		OutputTokens: 500,
		CreatedAt:    time.Date(2025, 12, 20, 8, 0, 0, 0, loc), // previous year
	})

	dayReport, err := QueryAIUsageReport("day", "2026-11-18")
	if err != nil {
		t.Fatalf("day report error: %v", err)
	}
	if dayReport.Summary.TotalRequests != 1 {
		t.Fatalf("expected 1 request for day, got %d", dayReport.Summary.TotalRequests)
	}

	weekReport, err := QueryAIUsageReport("week", "2026-11-18")
	if err != nil {
		t.Fatalf("week report error: %v", err)
	}
	if weekReport.Summary.TotalRequests != 2 {
		t.Fatalf("expected 2 requests for week, got %d", weekReport.Summary.TotalRequests)
	}

	monthReport, err := QueryAIUsageReport("month", "2026-11-18")
	if err != nil {
		t.Fatalf("month report error: %v", err)
	}
	if monthReport.Summary.TotalRequests != 3 {
		t.Fatalf("expected 3 requests for month, got %d", monthReport.Summary.TotalRequests)
	}

	yearReport, err := QueryAIUsageReport("year", "2026-11-18")
	if err != nil {
		t.Fatalf("year report error: %v", err)
	}
	if yearReport.Summary.TotalRequests != 4 {
		t.Fatalf("expected 4 requests for year, got %d", yearReport.Summary.TotalRequests)
	}

	// Purge older than 2026-11-01 -> should delete the 2026-05-10 and 2025-12-20 records (2 records)
	deleted, err := PurgeAIUsageBefore("2026-11-01")
	if err != nil {
		t.Fatalf("purge error: %v", err)
	}
	if deleted != 2 {
		t.Fatalf("expected 2 deleted records, got %d", deleted)
	}

	afterYearReport, err := QueryAIUsageReport("year", "2026-11-18")
	if err != nil {
		t.Fatalf("after purge year report error: %v", err)
	}
	if afterYearReport.Summary.TotalRequests != 3 {
		t.Fatalf("expected 3 remaining requests in 2026 after purge, got %d", afterYearReport.Summary.TotalRequests)
	}
}

func TestAdminAIUsageHandlers(t *testing.T) {
	origDB := db
	db = nil
	defer func() { db = origDB }()
	resetMemAIUsageForTest()

	RecordAIUsage(AIUsageRecord{
		Provider:     "gemini",
		Model:        "gemini-3.1-flash-lite",
		Feature:      "bible_context",
		Operation:    "generate",
		InputTokens:  1000,
		OutputTokens: 500,
		CreatedAt:    time.Now(),
	})

	// Unauthorized GET
	unauthReq := httptest.NewRequest(http.MethodGet, "/api/admin/ai-usage?period=month", nil)
	unauthRec := httptest.NewRecorder()
	adminAIUsageHandler(unauthRec, unauthReq)
	if unauthRec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 without admin token, got %d", unauthRec.Code)
	}

	// Authorized GET
	req := httptest.NewRequest(http.MethodGet, "/api/admin/ai-usage?period=month", nil)
	req.Header.Set("X-Admin-Token", adminSessionToken)
	rec := httptest.NewRecorder()
	adminAIUsageHandler(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", rec.Code, rec.Body.String())
	}

	var report AIUsageReport
	if err := json.Unmarshal(rec.Body.Bytes(), &report); err != nil {
		t.Fatalf("failed to decode report: %v", err)
	}
	if report.Summary.TotalRequests != 1 || report.Summary.TotalCostUSD <= 0 {
		t.Fatalf("unexpected report summary: %+v", report.Summary)
	}

	// Unauthorized cleanup
	unauthCleanup := httptest.NewRequest(http.MethodPost, "/api/admin/ai-usage/cleanup", bytes.NewBufferString(`{"before_date":"2030-01-01"}`))
	unauthCleanupRec := httptest.NewRecorder()
	adminAIUsageCleanupHandler(unauthCleanupRec, unauthCleanup)
	if unauthCleanupRec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for cleanup without admin token, got %d", unauthCleanupRec.Code)
	}

	// Authorized cleanup
	cleanupReq := httptest.NewRequest(http.MethodPost, "/api/admin/ai-usage/cleanup", bytes.NewBufferString(`{"before_date":"2030-01-01"}`))
	cleanupReq.Header.Set("X-Admin-Token", adminSessionToken)
	cleanupRec := httptest.NewRecorder()
	adminAIUsageCleanupHandler(cleanupRec, cleanupReq)
	if cleanupRec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for cleanup, got %d: %s", cleanupRec.Code, cleanupRec.Body.String())
	}
	var cleanupResp map[string]interface{}
	if err := json.Unmarshal(cleanupRec.Body.Bytes(), &cleanupResp); err != nil {
		t.Fatalf("failed to decode cleanup response: %v", err)
	}
	if cleanupResp["deleted_count"] != float64(1) {
		t.Fatalf("expected deleted_count=1, got %v", cleanupResp["deleted_count"])
	}
}

func TestAdminAIPricingCRUDAndCostOverride(t *testing.T) {
	origDB := db
	db = nil
	defer func() { db = origDB }()

	memModelPricingMu.Lock()
	memModelPricing = nil
	memModelPricingNextID = 1
	memModelPricingMu.Unlock()

	// Unauthorized GET
	unauthReq := httptest.NewRequest(http.MethodGet, "/api/admin/ai-pricing", nil)
	unauthRec := httptest.NewRecorder()
	adminAIPricingHandler(unauthRec, unauthReq)
	if unauthRec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 without admin token, got %d", unauthRec.Code)
	}

	// Authorized GET (returns seeded defaults)
	getReq := httptest.NewRequest(http.MethodGet, "/api/admin/ai-pricing", nil)
	getReq.Header.Set("X-Admin-Token", adminSessionToken)
	getRec := httptest.NewRecorder()
	adminAIPricingHandler(getRec, getReq)
	if getRec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", getRec.Code, getRec.Body.String())
	}

	// POST custom rate for a model
	payload := `{"rule":{"provider":"gemini","model":"gemini-custom-test","input_per_million":2.0,"output_per_million":8.0,"cached_per_million":0.5,"image_cost_each":0}}`
	postReq := httptest.NewRequest(http.MethodPost, "/api/admin/ai-pricing", bytes.NewBufferString(payload))
	postReq.Header.Set("X-Admin-Token", adminSessionToken)
	postRec := httptest.NewRecorder()
	adminAIPricingHandler(postRec, postReq)
	if postRec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on POST, got %d: %s", postRec.Code, postRec.Body.String())
	}

	// Verify EstimateAICostUSD uses custom rate: 500k uncached input ($1.00) + 500k cached ($0.25) + 500k output ($4.00) = $5.25
	customCost := EstimateAICostUSD("gemini", "gemini-custom-test", "generate", 1_000_000, 500_000, 0, 500_000, 0)
	if math.Abs(customCost-5.25) > 1e-6 {
		t.Fatalf("expected custom cost 5.25, got %f", customCost)
	}

	// DELETE custom rate
	delReq := httptest.NewRequest(http.MethodDelete, "/api/admin/ai-pricing?provider=gemini&model=gemini-custom-test", nil)
	delReq.Header.Set("X-Admin-Token", adminSessionToken)
	delRec := httptest.NewRecorder()
	adminAIPricingHandler(delRec, delReq)
	if delRec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on DELETE, got %d: %s", delRec.Code, delRec.Body.String())
	}
}
