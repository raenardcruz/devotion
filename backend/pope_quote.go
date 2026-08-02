package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type GetPopeQuoteResponse struct {
	Text string `json:"text"`
}

func get_pope_quote(dateStr string) (GetPopeQuoteResponse, error) {
	var t time.Time
	var err error
	if dateStr != "" {
		t, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			log.Printf("[get_pope_quote] Error parsing date %s, falling back to today: %v", dateStr, err)
			t = time.Now()
		}
	} else {
		t = time.Now()
	}

	// URL format: https://www.vaticannews.va/en/word-of-the-day/2026/04/14.html
	url := fmt.Sprintf("https://www.vaticannews.va/en/word-of-the-day/%d/%02d/%02d.html", t.Year(), t.Month(), t.Day())

	c := colly.NewCollector(
		colly.AllowedDomains("www.vaticannews.va"),
	)

	var popeQuote string

	// Based on inspection, we look for section--evidence with h2 "The words of the Popes"
	c.OnHTML("section.section--evidence", func(e *colly.HTMLElement) {
		title := e.ChildText("h2")
		if strings.Contains(strings.ToLower(title), "the words of the popes") {
			popeQuote = e.ChildText(".section__content")
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("[get_pope_quote] Error requesting %s: %v", r.Request.URL, err)
	})

	log.Printf("Fetching Pope quote from %s...", url)
	err = c.Visit(url)
	if err != nil {
		return GetPopeQuoteResponse{}, fmt.Errorf("failed to visit Vatican News: %w", err)
	}

	return GetPopeQuoteResponse{Text: strings.TrimSpace(popeQuote)}, nil
}
