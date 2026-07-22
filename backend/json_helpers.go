package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

func ExtractDevotionData(content string) (*DevotionData, error) {
	content = strings.ReplaceAll(content, "\\'", "'")

	// Try to find all potential JSON objects
	// We'll look for anything between { and }
	var lastErr error
	for i := 0; i < len(content); i++ {
		if content[i] == '{' {
			for j := len(content) - 1; j > i; j-- {
				if content[j] == '}' {
					jsonStr := content[i : j+1]
					var data DevotionData
					if err := json.Unmarshal([]byte(jsonStr), &data); err == nil {
						if validatesDevotionData(&data) {
							return &data, nil
						}
					} else {
						lastErr = err
					}
				}
			}
		}
	}

	if lastErr != nil {
		return nil, lastErr
	}
	return nil, fmt.Errorf("no valid devotion data found in content")
}

func validatesDevotionData(data *DevotionData) bool {
	if data.FirstReading == nil && data.Gospel == nil && data.ResponsorialPsalm == nil && data.SecondReading == nil && data.PopeQuote == "" {
		return false
	}
	return true
}

func PopulatePassageTexts(data *DevotionData) error {
	if data.FirstReading != nil && data.FirstReading.Citation != "" {
		text, _, err := get_bible_passage(data.FirstReading.Citation)
		if err == nil {
			data.FirstReading.Text = text
		}
	}
	if data.ResponsorialPsalm != nil && data.ResponsorialPsalm.Citation != "" {
		text, _, err := get_bible_passage(data.ResponsorialPsalm.Citation)
		if err == nil {
			data.ResponsorialPsalm.Text = text
		}
	}
	if data.SecondReading != nil && data.SecondReading.Citation != "" {
		text, _, err := get_bible_passage(data.SecondReading.Citation)
		if err == nil {
			data.SecondReading.Text = text
		}
	}
	if data.Gospel != nil && data.Gospel.Citation != "" {
		text, _, err := get_bible_passage(data.Gospel.Citation)
		if err == nil {
			data.Gospel.Text = text
		}
	}
	return nil
}
