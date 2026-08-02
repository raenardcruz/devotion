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
	var targetDate string
	if dateStr != "" {
		t, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			log.Printf("[get_mass_readings] Error parsing date %s, falling back to today: %v", dateStr, err)
			targetDate = time.Now().Format("010206")
		} else {
			targetDate = t.Format("010206")
		}
	} else {
		targetDate = time.Now().Format("010206")
	}
	url := fmt.Sprintf("https://bible.usccb.org/bible/readings/%s", targetDate)

	c := colly.NewCollector(
		colly.AllowedDomains("bible.usccb.org"),
	)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	c.WithTransport(tr)
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36"

	var readings []Readings

	c.OnHTML(".innerblock", func(e *colly.HTMLElement) {
		typeText := strings.TrimSpace(e.ChildText("h3.name"))
		citation := strings.TrimSpace(e.ChildText(".address a"))
		if citation == "" {
			citation = strings.TrimSpace(e.ChildText(".address"))
		}
		citation = strings.ReplaceAll(citation, "\n", " ")
		citation = strings.TrimSpace(citation)

		reading := Readings{
			Type:     typeText,
			Citation: citation,
		}

		if reading.Type != "" {
			readings = append(readings, reading)
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("[get_mass_readings] Error requesting %s: %v", r.Request.URL, err)
	})

	log.Printf("Fetching readings from %s...", url)
	err := c.Visit(url)
	if err != nil {
		if !strings.HasSuffix(url, ".cfm") {
			log.Printf("[get_mass_readings] Error visiting %s: %v. Retrying with .cfm...", url, err)
			url = url + ".cfm"
			log.Printf("Fetching readings from %s...", url)
			err = c.Visit(url)
		}
	}
	if err != nil {
		return GetMassReadingsResponse{}, fmt.Errorf("failed to visit USCCB: %w", err)
	}

	return GetMassReadingsResponse{Readings: readings}, nil
}
