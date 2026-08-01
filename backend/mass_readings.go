package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gocolly/colly"
	"google.golang.org/adk/tool"
)

type Readings struct {
	Type     string `json:"type"`
	Citation string `json:"citation"`
}

type GetMassReadingsArgs struct {
	Date string `json:"date"` // Format: YYYY-MM-DD
}

type GetMassReadingsResponse struct {
	Readings []Readings `json:"readings"`
}

func get_mass_readings(ctx tool.Context, args GetMassReadingsArgs) (GetMassReadingsResponse, error) {
	var targetDate string
	if args.Date != "" {
		t, err := time.Parse("2006-01-02", args.Date)
		if err != nil {
			log.Printf("[get_mass_readings] Error parsing date %s, falling back to today: %v", args.Date, err)
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

	var readings []Readings

	c.OnHTML(".innerblock", func(e *colly.HTMLElement) {
		reading := Readings{
			Type:     e.ChildText("h3.name"),
			Citation: e.ChildText(".address a"),
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
		return GetMassReadingsResponse{}, fmt.Errorf("failed to visit USCCB: %w", err)
	}

	return GetMassReadingsResponse{Readings: readings}, nil
}
