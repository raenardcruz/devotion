package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type Readings struct {
	Type     string `json:"type"`
	Citation string `json:"citation"`
}

type GetMassReadingsResponse struct {
	Readings []Readings `json:"readings"`
}

func get_mass_readings(dateStr string) (GetMassReadingsResponse, error) {
	var targetDateFormatted string // YYYY-MM-DD
	var targetDateMMDDYY string    // MMDDYY

	if dateStr != "" {
		t, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return GetMassReadingsResponse{}, fmt.Errorf("invalid reading date %q: use YYYY-MM-DD", dateStr)
		}
		targetDateFormatted = t.Format("2006-01-02")
		targetDateMMDDYY = t.Format("010206")
	} else {
		now := time.Now()
		targetDateFormatted = now.Format("2006-01-02")
		targetDateMMDDYY = now.Format("010206")
	}

	c := colly.NewCollector(
		colly.AllowedDomains("bible.usccb.org"),
	)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c.WithTransport(tr)
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36"

	// Step 1: Attempt to find reading path from calendar table (td[data-date="YYYY-MM-DD"])
	var calendarPath string
	calCollector := c.Clone()
	calCollector.OnHTML("td[data-date]", func(e *colly.HTMLElement) {
		if e.Attr("data-date") == targetDateFormatted {
			link := e.ChildAttr("a", "href")
			if link != "" {
				calendarPath = link
			}
		}
	})

	log.Printf("Fetching readings calendar from https://bible.usccb.org/readings/calendar...")
	if err := calCollector.Visit("https://bible.usccb.org/readings/calendar"); err != nil {
		log.Printf("[get_mass_readings] Warning fetching calendar: %v", err)
	}

	var urlsToTry []string
	if calendarPath != "" {
		if !strings.HasPrefix(calendarPath, "http") {
			if !strings.HasPrefix(calendarPath, "/") {
				calendarPath = "/" + calendarPath
			}
			calendarPath = "https://bible.usccb.org" + calendarPath
		}
		log.Printf("Found calendar reading URL for date %s: %s", targetDateFormatted, calendarPath)
		urlsToTry = append(urlsToTry, calendarPath)
	}

	// Step 2: Fallbacks to date-based URLs (.cfm and standard date)
	urlsToTry = append(urlsToTry,
		fmt.Sprintf("https://bible.usccb.org/bible/readings/%s.cfm", targetDateMMDDYY),
		fmt.Sprintf("https://bible.usccb.org/bible/readings/%s", targetDateMMDDYY),
	)

	var readings []Readings

	// Target .container blocks containing .name and .address matching catholic-mass-readings & USCCB layout
	c.OnHTML(".container", func(e *colly.HTMLElement) {
		typeText := strings.TrimSpace(e.ChildText(".name"))
		if typeText == "" {
			typeText = strings.TrimSpace(e.ChildText("h3.name"))
		}

		citation := strings.TrimSpace(e.ChildText(".address a"))
		if citation == "" {
			citation = strings.TrimSpace(e.ChildText(".address"))
		}
		citation = cleanScrapedText(citation)

		if typeText != "" && citation != "" {
			readings = append(readings, Readings{
				Type:     typeText,
				Citation: citation,
			})
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("[get_mass_readings] Error requesting %s: %v", r.Request.URL, err)
	})

	var lastErr error
	for _, url := range urlsToTry {
		log.Printf("Fetching readings from %s...", url)
		err := c.Visit(url)
		if err == nil && len(readings) > 0 {
			return GetMassReadingsResponse{Readings: readings}, nil
		}
		if err != nil {
			lastErr = err
		}
	}

	if len(readings) > 0 {
		return GetMassReadingsResponse{Readings: readings}, nil
	}

	if lastErr != nil {
		return GetMassReadingsResponse{}, fmt.Errorf("failed to visit USCCB: %w", lastErr)
	}

	return GetMassReadingsResponse{Readings: readings}, nil
}

func cleanScrapedText(s string) string {
	s = strings.ReplaceAll(s, "\u00a0", " ")
	s = strings.ReplaceAll(s, "\n", " ")
	words := strings.Fields(s)
	return strings.Join(words, " ")
}
