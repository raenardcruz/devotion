package main

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestGetMassReadingsRejectsInvalidDate(t *testing.T) {
	_, err := get_mass_readings("08/02/2026")
	if err == nil {
		t.Fatal("expected an invalid date to return an error")
	}
	if !strings.Contains(err.Error(), "use YYYY-MM-DD") {
		t.Errorf("expected date format guidance, got %v", err)
	}
}

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

func TestExtractMassReadingsDoesNotJoinAcrossContainer(t *testing.T) {
	html := `
	<html><body>
	<div class="container">
		<div class="page-header">
			<h2 class="name">Wednesday of the Thirty-third Week in Ordinary Time</h2>
		</div>
		<div class="innerblock">
			<div class="content-header">
				<h3 class="name">Reading 1</h3>
				<div class="address"><a href="/bible/revelation/4">Revelation 4:1-11</a></div>
			</div>
			<div class="content-body"><p>After this I had a vision...</p></div>
		</div>
		<div class="innerblock">
			<div class="content-header">
				<h3 class="name">Responsorial Psalm</h3>
				<div class="address"><a href="/bible/psalms/150">Psalm 150:1b-2, 3-4, 5-6</a></div>
			</div>
		</div>
		<div class="innerblock">
			<div class="content-header">
				<h3 class="name">Alleluia</h3>
				<div class="address">See <a href="/bible/john/15">John 15:16</a></div>
			</div>
		</div>
		<div class="innerblock">
			<div class="content-header">
				<h3 class="name">Gospel</h3>
				<div class="address"><a href="/bible/luke/19">Luke 19:11-28</a></div>
			</div>
		</div>
	</div>
	</body></html>`

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatalf("failed to parse HTML: %v", err)
	}

	got := extractMassReadings(doc.Selection)
	want := []Readings{
		{Type: "Reading 1", Citation: "Revelation 4:1-11"},
		{Type: "Responsorial Psalm", Citation: "Psalm 150:1b-2, 3-4, 5-6"},
		{Type: "Alleluia", Citation: "See John 15:16"},
		{Type: "Gospel", Citation: "Luke 19:11-28"},
	}

	if len(got) != len(want) {
		t.Fatalf("expected %d readings, got %d: %+v", len(want), len(got), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("reading[%d] = %+v, want %+v", i, got[i], want[i])
		}
	}
}

func TestExtractMassReadingsFlatSiblingsInContainer(t *testing.T) {
	html := `
	<html><body>
	<div class="container">
		<h3 class="name">Reading 1</h3>
		<div class="address"><a href="/bible/isaiah/55">Isaiah 55:1-3</a></div>
		<h3 class="name">Responsorial Psalm</h3>
		<div class="address"><a href="/bible/psalms/145">Psalm 145:8-9, 15-16, 17-18</a></div>
		<h3 class="name">Reading 2</h3>
		<div class="address"><a href="/bible/romans/8">Romans 8:35, 37-39</a></div>
		<h3 class="name">Gospel</h3>
		<div class="address"><a href="/bible/matthew/14">Matthew 14:13-21</a></div>
	</div>
	</body></html>`

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatalf("failed to parse HTML: %v", err)
	}

	got := extractMassReadings(doc.Selection)
	want := []Readings{
		{Type: "Reading 1", Citation: "Isaiah 55:1-3"},
		{Type: "Responsorial Psalm", Citation: "Psalm 145:8-9, 15-16, 17-18"},
		{Type: "Reading 2", Citation: "Romans 8:35, 37-39"},
		{Type: "Gospel", Citation: "Matthew 14:13-21"},
	}

	if len(got) != len(want) {
		t.Fatalf("expected %d readings, got %d: %+v", len(want), len(got), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("reading[%d] = %+v, want %+v", i, got[i], want[i])
		}
	}
}
