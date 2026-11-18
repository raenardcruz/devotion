package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"google.golang.org/genai"
)

// AIUsageRecord represents a single logged AI API call and its computed cost.
type AIUsageRecord struct {
	ID               int64     `json:"id"`
	Provider         string    `json:"provider"`
	Model            string    `json:"model"`
	Feature          string    `json:"feature"`
	Operation        string    `json:"operation"`
	InputTokens      int       `json:"input_tokens"`
	OutputTokens     int       `json:"output_tokens"`
	ThoughtsTokens   int       `json:"thoughts_tokens"`
	CachedTokens     int       `json:"cached_tokens"`
	TotalTokens      int       `json:"total_tokens"`
	ImageCount       int       `json:"image_count"`
	EstimatedCostUSD float64   `json:"estimated_cost_usd"`
	CreatedAt        time.Time `json:"created_at"`
}

// AIUsageSummary aggregates totals for the selected period.
type AIUsageSummary struct {
	TotalRequests       int     `json:"total_requests"`
	TotalInputTokens    int64   `json:"total_input_tokens"`
	TotalOutputTokens   int64   `json:"total_output_tokens"`
	TotalThoughtsTokens int64   `json:"total_thoughts_tokens"`
	TotalCachedTokens   int64   `json:"total_cached_tokens"`
	TotalTokens         int64   `json:"total_tokens"`
	TotalImages         int64   `json:"total_images"`
	TotalCostUSD        float64 `json:"total_cost_usd"`
}

// AIUsageProviderBreakdown summarizes usage grouped by provider and model.
type AIUsageProviderBreakdown struct {
	Provider         string  `json:"provider"`
	Model            string  `json:"model"`
	Requests         int     `json:"requests"`
	InputTokens      int64   `json:"input_tokens"`
	OutputTokens     int64   `json:"output_tokens"`
	ThoughtsTokens   int64   `json:"thoughts_tokens"`
	CachedTokens     int64   `json:"cached_tokens"`
	TotalTokens      int64   `json:"total_tokens"`
	ImageCount       int64   `json:"image_count"`
	EstimatedCostUSD float64 `json:"estimated_cost_usd"`
}

// AIUsageFeatureBreakdown summarizes usage grouped by app feature.
type AIUsageFeatureBreakdown struct {
	Feature          string  `json:"feature"`
	Requests         int     `json:"requests"`
	TotalTokens      int64   `json:"total_tokens"`
	ImageCount       int64   `json:"image_count"`
	EstimatedCostUSD float64 `json:"estimated_cost_usd"`
}

// AIUsageTimelinePoint represents a time bucket (hour/day/month) for charts.
type AIUsageTimelinePoint struct {
	Bucket           string  `json:"bucket"`
	Requests         int     `json:"requests"`
	TotalTokens      int64   `json:"total_tokens"`
	EstimatedCostUSD float64 `json:"estimated_cost_usd"`
}

// AIUsageReport is the full response returned by GET /api/admin/ai-usage.
type AIUsageReport struct {
	Period     string                     `json:"period"`
	StartDate  string                     `json:"start_date"`
	EndDate    string                     `json:"end_date"`
	Summary    AIUsageSummary             `json:"summary"`
	ByProvider []AIUsageProviderBreakdown `json:"by_provider"`
	ByFeature  []AIUsageFeatureBreakdown  `json:"by_feature"`
	Timeline   []AIUsageTimelinePoint     `json:"timeline"`
	RecentLogs []AIUsageRecord            `json:"recent_logs"`
}

// AIModelPricing represents configurable per-model pricing rates (USD per 1M tokens / per image).
type AIModelPricing struct {
	ID               int64     `json:"id"`
	Provider         string    `json:"provider"`
	Model            string    `json:"model"`
	InputPerMillion  float64   `json:"input_per_million"`
	OutputPerMillion float64   `json:"output_per_million"`
	CachedPerMillion float64   `json:"cached_per_million"`
	ImageCostEach    float64   `json:"image_cost_each"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type modelPricing struct {
	InputPerMillion  float64
	OutputPerMillion float64
	CachedPerMillion float64
	ImageCostEach    float64
}

var (
	memAIUsageMu     sync.RWMutex
	memAIUsageLogs   []AIUsageRecord
	memAIUsageNextID int64 = 1

	memModelPricingMu     sync.RWMutex
	memModelPricing       []AIModelPricing
	memModelPricingNextID int64 = 1
)

func defaultModelPricingList() []AIModelPricing {
	return []AIModelPricing{
		{Provider: "gemini", Model: "gemini-3.1-flash-lite", InputPerMillion: 0.075, OutputPerMillion: 0.30, CachedPerMillion: 0.01875, ImageCostEach: 0},
		{Provider: "gemini", Model: "gemini-2.5-flash", InputPerMillion: 0.15, OutputPerMillion: 0.60, CachedPerMillion: 0.0375, ImageCostEach: 0},
		{Provider: "gemini", Model: "gemini-2.5-pro", InputPerMillion: 1.25, OutputPerMillion: 10.00, CachedPerMillion: 0.3125, ImageCostEach: 0},
		{Provider: "gemini", Model: "gemini-3.8-live", InputPerMillion: 0.50, OutputPerMillion: 2.00, CachedPerMillion: 0.125, ImageCostEach: 0},
		{Provider: "gemini", Model: "gemini-nano-banana-2.1", InputPerMillion: 0.15, OutputPerMillion: 0.60, CachedPerMillion: 0.0375, ImageCostEach: 0.039},
		{Provider: "ollama", Model: "gemma4:cloud", InputPerMillion: 0, OutputPerMillion: 0, CachedPerMillion: 0, ImageCostEach: 0},
	}
}

func lookupCustomModelPricing(provider, model string) (modelPricing, bool) {
	p := strings.ToLower(strings.TrimSpace(provider))
	m := strings.ToLower(strings.TrimSpace(model))
	if m == "" {
		return modelPricing{}, false
	}

	rules, err := GetAllModelPricing()
	if err != nil {
		return modelPricing{}, false
	}
	for _, r := range rules {
		if strings.ToLower(strings.TrimSpace(r.Provider)) == p && strings.ToLower(strings.TrimSpace(r.Model)) == m {
			return modelPricing{
				InputPerMillion:  r.InputPerMillion,
				OutputPerMillion: r.OutputPerMillion,
				CachedPerMillion: r.CachedPerMillion,
				ImageCostEach:    r.ImageCostEach,
			}, true
		}
	}
	return modelPricing{}, false
}

// getModelPricing returns estimated USD pricing per 1M tokens (or per image) for a provider/model.
// Checks configured ai_model_pricing rules first, then falls back to built-in defaults.
func getModelPricing(provider, model, operation string) modelPricing {
	if custom, ok := lookupCustomModelPricing(provider, model); ok {
		return custom
	}

	p := strings.ToLower(strings.TrimSpace(provider))
	m := strings.ToLower(strings.TrimSpace(model))
	op := strings.ToLower(strings.TrimSpace(operation))

	if p == "ollama" || p == "magisterium" {
		return modelPricing{}
	}

	// Image generation models (Gemini Flash Image / Nano Banana)
	if op == "generate_image" || strings.Contains(m, "image") || strings.Contains(m, "nano-banana") || strings.Contains(m, "imagen") {
		return modelPricing{
			InputPerMillion:  0.15,
			OutputPerMillion: 0.60,
			CachedPerMillion: 0.0375,
			ImageCostEach:    0.039,
		}
	}

	// Gemini Live (audio/multimodal)
	if strings.Contains(m, "live") || op == "live_session" {
		return modelPricing{
			InputPerMillion:  0.50,
			OutputPerMillion: 2.00,
			CachedPerMillion: 0.125,
		}
	}

	// Gemini Pro models
	if strings.Contains(m, "pro") {
		return modelPricing{
			InputPerMillion:  1.25,
			OutputPerMillion: 10.00,
			CachedPerMillion: 0.3125,
		}
	}

	// Gemini Flash-Lite / 2.0 Flash
	if strings.Contains(m, "flash-lite") || strings.Contains(m, "2.0-flash") {
		return modelPricing{
			InputPerMillion:  0.075,
			OutputPerMillion: 0.30,
			CachedPerMillion: 0.01875,
		}
	}

	// Default Gemini Flash / Argon / general Gemini text model
	return modelPricing{
		InputPerMillion:  0.15,
		OutputPerMillion: 0.60,
		CachedPerMillion: 0.0375,
	}
}

// EstimateAICostUSD computes estimated USD cost from token counts and image count.
func EstimateAICostUSD(provider, model, operation string, inputTokens, outputTokens, thoughtsTokens, cachedTokens, imageCount int) float64 {
	pricing := getModelPricing(provider, model, operation)
	if pricing.InputPerMillion == 0 && pricing.OutputPerMillion == 0 && pricing.CachedPerMillion == 0 && pricing.ImageCostEach == 0 {
		return 0
	}

	uncachedInput := inputTokens - cachedTokens
	if uncachedInput < 0 {
		uncachedInput = 0
	}
	cachedRate := pricing.CachedPerMillion
	if cachedRate == 0 && pricing.InputPerMillion > 0 {
		cachedRate = pricing.InputPerMillion * 0.25
	}

	inputCost := (float64(uncachedInput) / 1_000_000.0) * pricing.InputPerMillion
	cachedCost := (float64(cachedTokens) / 1_000_000.0) * cachedRate
	outputCost := (float64(outputTokens+thoughtsTokens) / 1_000_000.0) * pricing.OutputPerMillion
	imageCost := float64(imageCount) * pricing.ImageCostEach

	return inputCost + cachedCost + outputCost + imageCost
}

// tokenCountsFromGenAI extracts token counts from genai.GenerateContentResponseUsageMetadata.
func tokenCountsFromGenAI(meta *genai.GenerateContentResponseUsageMetadata) (input, output, thoughts, cached, total int) {
	if meta == nil {
		return 0, 0, 0, 0, 0
	}
	input = int(meta.PromptTokenCount)
	output = int(meta.CandidatesTokenCount)
	thoughts = int(meta.ThoughtsTokenCount)
	cached = int(meta.CachedContentTokenCount)
	total = int(meta.TotalTokenCount)
	if total == 0 {
		total = input + output + thoughts
	}
	return
}

// tokenCountsFromLive extracts token counts from genai.UsageMetadata (Live API).
func tokenCountsFromLive(meta *genai.UsageMetadata) (input, output, thoughts, cached, total int) {
	if meta == nil {
		return 0, 0, 0, 0, 0
	}
	input = int(meta.PromptTokenCount)
	output = int(meta.ResponseTokenCount)
	thoughts = int(meta.ThoughtsTokenCount)
	cached = int(meta.CachedContentTokenCount)
	total = int(meta.TotalTokenCount)
	if total == 0 {
		total = input + output + thoughts
	}
	return
}

// RecordAIUsage computes cost if not already set and persists the usage record.
func RecordAIUsage(record AIUsageRecord) {
	if record.CreatedAt.IsZero() {
		record.CreatedAt = time.Now()
	}
	if record.TotalTokens == 0 {
		record.TotalTokens = record.InputTokens + record.OutputTokens + record.ThoughtsTokens
	}
	if record.EstimatedCostUSD == 0 {
		record.EstimatedCostUSD = EstimateAICostUSD(
			record.Provider,
			record.Model,
			record.Operation,
			record.InputTokens,
			record.OutputTokens,
			record.ThoughtsTokens,
			record.CachedTokens,
			record.ImageCount,
		)
	}

	if db == nil {
		memAIUsageMu.Lock()
		record.ID = memAIUsageNextID
		memAIUsageNextID++
		memAIUsageLogs = append(memAIUsageLogs, record)
		memAIUsageMu.Unlock()
		return
	}

	query := `
	INSERT INTO ai_usage_logs (
		provider, model, feature, operation,
		input_tokens, output_tokens, thoughts_tokens, cached_tokens, total_tokens,
		image_count, estimated_cost_usd, created_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	if _, err := db.Exec(
		query,
		record.Provider,
		record.Model,
		record.Feature,
		record.Operation,
		record.InputTokens,
		record.OutputTokens,
		record.ThoughtsTokens,
		record.CachedTokens,
		record.TotalTokens,
		record.ImageCount,
		record.EstimatedCostUSD,
		record.CreatedAt,
	); err != nil {
		log.Printf("[RecordAIUsage] Failed to insert usage log: %v", err)
	}
}

// resolvePeriodBounds computes [start, end) time bounds and normalized period name.
func resolvePeriodBounds(period, anchorDate string) (string, time.Time, time.Time, error) {
	p := strings.ToLower(strings.TrimSpace(period))
	if p == "" {
		p = "month"
	}

	now := time.Now()
	loc := now.Location()
	anchor := now
	if strings.TrimSpace(anchorDate) != "" {
		parsed, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(anchorDate), loc)
		if err != nil {
			return "", time.Time{}, time.Time{}, fmt.Errorf("invalid date format (expected YYYY-MM-DD): %w", err)
		}
		anchor = parsed
	}

	dayStart := time.Date(anchor.Year(), anchor.Month(), anchor.Day(), 0, 0, 0, 0, loc)

	switch p {
	case "day":
		return "day", dayStart, dayStart.AddDate(0, 0, 1), nil
	case "week":
		// ISO week starting Monday
		weekday := int(dayStart.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		weekStart := dayStart.AddDate(0, 0, -(weekday - 1))
		return "week", weekStart, weekStart.AddDate(0, 0, 7), nil
	case "month":
		monthStart := time.Date(anchor.Year(), anchor.Month(), 1, 0, 0, 0, 0, loc)
		return "month", monthStart, monthStart.AddDate(0, 1, 0), nil
	case "year":
		yearStart := time.Date(anchor.Year(), 1, 1, 0, 0, 0, 0, loc)
		return "year", yearStart, yearStart.AddDate(1, 0, 0), nil
	default:
		return "", time.Time{}, time.Time{}, fmt.Errorf("invalid period %q (must be day, week, month, or year)", period)
	}
}

func timelineBucketKey(t time.Time, period string) string {
	switch period {
	case "day":
		return t.Format("15:00")
	case "year":
		return t.Format("2006-01")
	default:
		return t.Format("2006-01-02")
	}
}

// QueryAIUsageReport returns aggregated AI usage and cost for the requested period.
func QueryAIUsageReport(period, anchorDate string) (AIUsageReport, error) {
	normPeriod, start, end, err := resolvePeriodBounds(period, anchorDate)
	if err != nil {
		return AIUsageReport{}, err
	}

	var records []AIUsageRecord
	if db == nil {
		memAIUsageMu.RLock()
		for _, rec := range memAIUsageLogs {
			if !rec.CreatedAt.Before(start) && rec.CreatedAt.Before(end) {
				records = append(records, rec)
			}
		}
		memAIUsageMu.RUnlock()
		sort.Slice(records, func(i, j int) bool {
			return records[i].CreatedAt.After(records[j].CreatedAt)
		})
	} else {
		rows, qErr := db.Query(`
			SELECT id, provider, model, feature, operation,
			       input_tokens, output_tokens, thoughts_tokens, cached_tokens, total_tokens,
			       image_count, estimated_cost_usd, created_at
			FROM ai_usage_logs
			WHERE created_at >= $1 AND created_at < $2
			ORDER BY created_at DESC
		`, start, end)
		if qErr != nil {
			return AIUsageReport{}, qErr
		}
		defer rows.Close()

		for rows.Next() {
			var rec AIUsageRecord
			if scanErr := rows.Scan(
				&rec.ID,
				&rec.Provider,
				&rec.Model,
				&rec.Feature,
				&rec.Operation,
				&rec.InputTokens,
				&rec.OutputTokens,
				&rec.ThoughtsTokens,
				&rec.CachedTokens,
				&rec.TotalTokens,
				&rec.ImageCount,
				&rec.EstimatedCostUSD,
				&rec.CreatedAt,
			); scanErr != nil {
				return AIUsageReport{}, scanErr
			}
			records = append(records, rec)
		}
		if rowsErr := rows.Err(); rowsErr != nil {
			return AIUsageReport{}, rowsErr
		}
	}

	report := AIUsageReport{
		Period:     normPeriod,
		StartDate:  start.Format("2006-01-02"),
		EndDate:    end.AddDate(0, 0, -1).Format("2006-01-02"),
		ByProvider: []AIUsageProviderBreakdown{},
		ByFeature:  []AIUsageFeatureBreakdown{},
		Timeline:   []AIUsageTimelinePoint{},
		RecentLogs: []AIUsageRecord{},
	}

	providerMap := make(map[string]*AIUsageProviderBreakdown)
	featureMap := make(map[string]*AIUsageFeatureBreakdown)
	timelineMap := make(map[string]*AIUsageTimelinePoint)

	for i, rec := range records {
		report.Summary.TotalRequests++
		report.Summary.TotalInputTokens += int64(rec.InputTokens)
		report.Summary.TotalOutputTokens += int64(rec.OutputTokens)
		report.Summary.TotalThoughtsTokens += int64(rec.ThoughtsTokens)
		report.Summary.TotalCachedTokens += int64(rec.CachedTokens)
		report.Summary.TotalTokens += int64(rec.TotalTokens)
		report.Summary.TotalImages += int64(rec.ImageCount)
		report.Summary.TotalCostUSD += rec.EstimatedCostUSD

		pmKey := rec.Provider + "::" + rec.Model
		pb, ok := providerMap[pmKey]
		if !ok {
			pb = &AIUsageProviderBreakdown{
				Provider: rec.Provider,
				Model:    rec.Model,
			}
			providerMap[pmKey] = pb
		}
		pb.Requests++
		pb.InputTokens += int64(rec.InputTokens)
		pb.OutputTokens += int64(rec.OutputTokens)
		pb.ThoughtsTokens += int64(rec.ThoughtsTokens)
		pb.CachedTokens += int64(rec.CachedTokens)
		pb.TotalTokens += int64(rec.TotalTokens)
		pb.ImageCount += int64(rec.ImageCount)
		pb.EstimatedCostUSD += rec.EstimatedCostUSD

		featName := rec.Feature
		if featName == "" {
			featName = "general"
		}
		fb, ok := featureMap[featName]
		if !ok {
			fb = &AIUsageFeatureBreakdown{Feature: featName}
			featureMap[featName] = fb
		}
		fb.Requests++
		fb.TotalTokens += int64(rec.TotalTokens)
		fb.ImageCount += int64(rec.ImageCount)
		fb.EstimatedCostUSD += rec.EstimatedCostUSD

		bucket := timelineBucketKey(rec.CreatedAt.In(start.Location()), normPeriod)
		tp, ok := timelineMap[bucket]
		if !ok {
			tp = &AIUsageTimelinePoint{Bucket: bucket}
			timelineMap[bucket] = tp
		}
		tp.Requests++
		tp.TotalTokens += int64(rec.TotalTokens)
		tp.EstimatedCostUSD += rec.EstimatedCostUSD

		if i < 100 {
			report.RecentLogs = append(report.RecentLogs, rec)
		}
	}

	for _, pb := range providerMap {
		report.ByProvider = append(report.ByProvider, *pb)
	}
	sort.Slice(report.ByProvider, func(i, j int) bool {
		if report.ByProvider[i].EstimatedCostUSD != report.ByProvider[j].EstimatedCostUSD {
			return report.ByProvider[i].EstimatedCostUSD > report.ByProvider[j].EstimatedCostUSD
		}
		return report.ByProvider[i].Requests > report.ByProvider[j].Requests
	})

	for _, fb := range featureMap {
		report.ByFeature = append(report.ByFeature, *fb)
	}
	sort.Slice(report.ByFeature, func(i, j int) bool {
		if report.ByFeature[i].EstimatedCostUSD != report.ByFeature[j].EstimatedCostUSD {
			return report.ByFeature[i].EstimatedCostUSD > report.ByFeature[j].EstimatedCostUSD
		}
		return report.ByFeature[i].Requests > report.ByFeature[j].Requests
	})

	for _, tp := range timelineMap {
		report.Timeline = append(report.Timeline, *tp)
	}
	sort.Slice(report.Timeline, func(i, j int) bool {
		return report.Timeline[i].Bucket < report.Timeline[j].Bucket
	})

	return report, nil
}

// PurgeAIUsageBefore deletes all records strictly older than the given YYYY-MM-DD date.
func PurgeAIUsageBefore(beforeDate string) (int64, error) {
	trimmed := strings.TrimSpace(beforeDate)
	if trimmed == "" {
		return 0, fmt.Errorf("before_date is required (YYYY-MM-DD)")
	}
	cutoff, err := time.ParseInLocation("2006-01-02", trimmed, time.Now().Location())
	if err != nil {
		return 0, fmt.Errorf("invalid before_date format (expected YYYY-MM-DD): %w", err)
	}

	if db == nil {
		memAIUsageMu.Lock()
		defer memAIUsageMu.Unlock()
		kept := make([]AIUsageRecord, 0, len(memAIUsageLogs))
		var deleted int64
		for _, rec := range memAIUsageLogs {
			if rec.CreatedAt.Before(cutoff) {
				deleted++
			} else {
				kept = append(kept, rec)
			}
		}
		memAIUsageLogs = kept
		return deleted, nil
	}

	res, execErr := db.Exec(`DELETE FROM ai_usage_logs WHERE created_at < $1`, cutoff)
	if execErr != nil {
		return 0, execErr
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affected, nil
}

func adminAIUsageHandler(w http.ResponseWriter, r *http.Request) {
	if !isAdminRequest(r) {
		respondWithError(w, http.StatusUnauthorized, "Admin authentication required")
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	period := r.URL.Query().Get("period")
	dateParam := r.URL.Query().Get("date")

	report, err := QueryAIUsageReport(period, dateParam)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

func adminAIUsageCleanupHandler(w http.ResponseWriter, r *http.Request) {
	if !isAdminRequest(r) {
		respondWithError(w, http.StatusUnauthorized, "Admin authentication required")
		return
	}
	if r.Method != http.MethodPost && r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	beforeDate := strings.TrimSpace(r.URL.Query().Get("before"))
	if beforeDate == "" && r.Body != nil {
		var body struct {
			BeforeDate string `json:"before_date"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err == nil {
			beforeDate = strings.TrimSpace(body.BeforeDate)
		}
	}

	deleted, err := PurgeAIUsageBefore(beforeDate)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":       true,
		"before_date":   beforeDate,
		"deleted_count": deleted,
	})
}

// GetAllModelPricing returns configured model pricing rules (or seeded defaults if none exist yet).
func GetAllModelPricing() ([]AIModelPricing, error) {
	if db == nil {
		memModelPricingMu.RLock()
		if len(memModelPricing) == 0 {
			memModelPricingMu.RUnlock()
			defaults := defaultModelPricingList()
			memModelPricingMu.Lock()
			if len(memModelPricing) == 0 {
				now := time.Now()
				for _, d := range defaults {
					d.ID = memModelPricingNextID
					d.UpdatedAt = now
					memModelPricingNextID++
					memModelPricing = append(memModelPricing, d)
				}
			}
			out := append([]AIModelPricing(nil), memModelPricing...)
			memModelPricingMu.Unlock()
			return out, nil
		}
		out := append([]AIModelPricing(nil), memModelPricing...)
		memModelPricingMu.RUnlock()
		return out, nil
	}

	rows, err := db.Query(`
		SELECT id, provider, model, input_per_million, output_per_million, cached_per_million, image_cost_each, updated_at
		FROM ai_model_pricing
		ORDER BY provider ASC, model ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []AIModelPricing
	for rows.Next() {
		var p AIModelPricing
		if err := rows.Scan(&p.ID, &p.Provider, &p.Model, &p.InputPerMillion, &p.OutputPerMillion, &p.CachedPerMillion, &p.ImageCostEach, &p.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Seed defaults on first read if table is empty
	if len(list) == 0 {
		defaults := defaultModelPricingList()
		if err := SaveModelPricingBatch(defaults); err == nil {
			return GetAllModelPricing()
		}
		return defaults, nil
	}
	return list, nil
}

// SaveModelPricingBatch upserts one or more model pricing rules.
func SaveModelPricingBatch(rules []AIModelPricing) error {
	now := time.Now()
	if db == nil {
		memModelPricingMu.Lock()
		defer memModelPricingMu.Unlock()
		for _, r := range rules {
			prov := strings.ToLower(strings.TrimSpace(r.Provider))
			if prov == "" {
				prov = "gemini"
			}
			mod := strings.TrimSpace(r.Model)
			if mod == "" {
				continue
			}
			updated := false
			for i := range memModelPricing {
				if strings.EqualFold(memModelPricing[i].Provider, prov) && strings.EqualFold(memModelPricing[i].Model, mod) {
					memModelPricing[i].InputPerMillion = r.InputPerMillion
					memModelPricing[i].OutputPerMillion = r.OutputPerMillion
					memModelPricing[i].CachedPerMillion = r.CachedPerMillion
					memModelPricing[i].ImageCostEach = r.ImageCostEach
					memModelPricing[i].UpdatedAt = now
					updated = true
					break
				}
			}
			if !updated {
				memModelPricing = append(memModelPricing, AIModelPricing{
					ID:               memModelPricingNextID,
					Provider:         prov,
					Model:            mod,
					InputPerMillion:  r.InputPerMillion,
					OutputPerMillion: r.OutputPerMillion,
					CachedPerMillion: r.CachedPerMillion,
					ImageCostEach:    r.ImageCostEach,
					UpdatedAt:        now,
				})
				memModelPricingNextID++
			}
		}
		return nil
	}

	for _, r := range rules {
		prov := strings.ToLower(strings.TrimSpace(r.Provider))
		if prov == "" {
			prov = "gemini"
		}
		mod := strings.TrimSpace(r.Model)
		if mod == "" {
			continue
		}
		_, err := db.Exec(`
			INSERT INTO ai_model_pricing (provider, model, input_per_million, output_per_million, cached_per_million, image_cost_each, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT (provider, model) DO UPDATE SET
				input_per_million = EXCLUDED.input_per_million,
				output_per_million = EXCLUDED.output_per_million,
				cached_per_million = EXCLUDED.cached_per_million,
				image_cost_each = EXCLUDED.image_cost_each,
				updated_at = EXCLUDED.updated_at
		`, prov, mod, r.InputPerMillion, r.OutputPerMillion, r.CachedPerMillion, r.ImageCostEach, now)
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteModelPricing removes a pricing rule by provider and model.
func DeleteModelPricing(provider, model string) error {
	prov := strings.ToLower(strings.TrimSpace(provider))
	mod := strings.TrimSpace(model)
	if mod == "" {
		return fmt.Errorf("model is required")
	}

	if db == nil {
		memModelPricingMu.Lock()
		defer memModelPricingMu.Unlock()
		kept := make([]AIModelPricing, 0, len(memModelPricing))
		for _, item := range memModelPricing {
			if strings.EqualFold(item.Provider, prov) && strings.EqualFold(item.Model, mod) {
				continue
			}
			kept = append(kept, item)
		}
		memModelPricing = kept
		return nil
	}

	_, err := db.Exec(`DELETE FROM ai_model_pricing WHERE LOWER(provider) = LOWER($1) AND LOWER(model) = LOWER($2)`, prov, mod)
	return err
}

func adminAIPricingHandler(w http.ResponseWriter, r *http.Request) {
	if !isAdminRequest(r) {
		respondWithError(w, http.StatusUnauthorized, "Admin authentication required")
		return
	}

	switch r.Method {
	case http.MethodGet:
		pricing, err := GetAllModelPricing()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to load model pricing: %v", err))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"pricing": pricing,
		})

	case http.MethodPost, http.MethodPut:
		var payload struct {
			Pricing []AIModelPricing `json:"pricing"`
			Rule    *AIModelPricing  `json:"rule,omitempty"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid pricing payload")
			return
		}
		rules := payload.Pricing
		if payload.Rule != nil {
			rules = append(rules, *payload.Rule)
		}
		if len(rules) == 0 {
			respondWithError(w, http.StatusBadRequest, "No pricing rules provided")
			return
		}
		if err := SaveModelPricingBatch(rules); err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to save pricing: %v", err))
			return
		}
		updated, err := GetAllModelPricing()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Saved but failed to reload pricing: %v", err))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"pricing": updated,
		})

	case http.MethodDelete:
		provider := r.URL.Query().Get("provider")
		model := r.URL.Query().Get("model")
		if err := DeleteModelPricing(provider, model); err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		updated, _ := GetAllModelPricing()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"pricing": updated,
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
