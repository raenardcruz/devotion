package main

import (
	"strings"
	"testing"
)

func TestGetMassReadings(t *testing.T) {
	resp, err := get_mass_readings("2026-08-02")
	if err != nil {
		t.Fatalf("get_mass_readings error: %v", err)
	}
	t.Logf("Readings count: %d", len(resp.Readings))
	for i, r := range resp.Readings {
		t.Logf("Reading %d: Type=%q, Citation=%q", i, r.Type, r.Citation)
	}
}

func TestCategorizeReadings(t *testing.T) {
	// Simulate USCCB output format
	readingsResp := GetMassReadingsResponse{
		Readings: []Readings{
			{Type: "Reading 1", Citation: "Isaiah 55:1-3"},
			{Type: "Responsorial Psalm", Citation: "Psalm 145:8-9"},
			{Type: "Reading 2", Citation: "Romans 8:35"},
			{Type: "Gospel", Citation: "Matthew 14:13-21"},
		},
	}

	devotion := &DevotionData{}
	for _, r := range readingsResp.Readings {
		tStr := strings.ToLower(r.Type)
		rw := &ReadingWithContext{Citation: r.Citation}
		switch {
		case strings.Contains(tStr, "second reading") || strings.Contains(tStr, "reading 2") || strings.Contains(tStr, "reading ii"):
			devotion.SecondReading = rw
		case strings.Contains(tStr, "first reading") || strings.Contains(tStr, "reading 1") || strings.Contains(tStr, "reading i"):
			devotion.FirstReading = rw
		case strings.Contains(tStr, "psalm"):
			devotion.ResponsorialPsalm = rw
		case strings.Contains(tStr, "gospel"):
			devotion.Gospel = rw
		}
	}

	if devotion.FirstReading == nil || devotion.FirstReading.Citation != "Isaiah 55:1-3" {
		t.Errorf("Expected First Reading 'Isaiah 55:1-3', got %v", devotion.FirstReading)
	}
	if devotion.SecondReading == nil || devotion.SecondReading.Citation != "Romans 8:35" {
		t.Errorf("Expected Second Reading 'Romans 8:35', got %v", devotion.SecondReading)
	}
	if devotion.ResponsorialPsalm == nil || devotion.ResponsorialPsalm.Citation != "Psalm 145:8-9" {
		t.Errorf("Expected Psalm 'Psalm 145:8-9', got %v", devotion.ResponsorialPsalm)
	}
	if devotion.Gospel == nil || devotion.Gospel.Citation != "Matthew 14:13-21" {
		t.Errorf("Expected Gospel 'Matthew 14:13-21', got %v", devotion.Gospel)
	}
}

