package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type BibleResponse struct {
	Data BibleContent `json:"data"`
}

type BibleContent struct {
	Content   string `json:"content"`
	Copyright string `json:"copyright"`
}

func get_bible_passage(passage string) (string, string, error) {
	log.Printf("[get_bible_passage] Getting bible passage for %s", passage)
	cacheKey := fmt.Sprintf("bible_passage:%s", passage)
	if redisClient != nil {
		val, err := redisClient.Get(context.Background(), cacheKey).Result()
		if err == nil {
			log.Printf("[get_bible_passage] Cache hit for %s", passage)
			var cached BibleContent
			if err := json.Unmarshal([]byte(val), &cached); err == nil {
				return cached.Content, cached.Copyright, nil
			}
			log.Printf("[get_bible_passage] Warning: failed to unmarshal cached bible passage: %v", err)
		}
	}

	passage_text := ""
	copy_right := ""
	passage_verses := passage_sanitize(passage)
	last_space_idx := strings.LastIndex(passage_verses, " ")
	if last_space_idx == -1 {
		return "", "", fmt.Errorf("invalid passage format: %s", passage)
	}
	passage_book := passage_verses[:last_space_idx]
	passage_coords := passage_verses[last_space_idx+1:]

	passage_book_short := get_bible_book_short(passage_book)
	if passage_book_short == "" {
		return "", "", fmt.Errorf("unknown bible book: %s", passage_book)
	}

	split_section := strings.Split(passage_coords, ",")
	current_chapter := ""

	for _, section := range split_section {
		rv := strings.Split(section, "-")
		if len(rv) == 0 {
			continue
		}

		// Parse start reference
		start_part := rv[0]
		start_sv := strings.Split(start_part, ":")
		start_verse := ""
		if len(start_sv) == 2 {
			current_chapter = start_sv[0]
			start_verse = start_sv[1]
		} else {
			if current_chapter == "" {
				return "", "", fmt.Errorf("invalid passage format: missing chapter for %s", section)
			}
			start_verse = start_sv[0]
		}

		transformed_passage := fmt.Sprintf("%s.%s.%s", passage_book_short, current_chapter, strip_verse_suffix(start_verse))

		if len(rv) > 1 {
			// Parse end reference
			end_part := rv[1]
			end_sv := strings.Split(end_part, ":")
			end_verse := ""
			end_chapter := current_chapter

			if len(end_sv) == 2 {
				end_chapter = end_sv[0]
				end_verse = end_sv[1]
				current_chapter = end_chapter
			} else {
				end_verse = end_sv[0]
			}

			transformed_passage += fmt.Sprintf("-%s.%s.%s", passage_book_short, end_chapter, strip_verse_suffix(end_verse))
		}

		resp, err := get_bible_passage_text(transformed_passage)
		if err != nil {
			return "", "", err
		}
		passage_text += resp.Content
		copy_right = resp.Copyright
	}
	if redisClient != nil && passage_text != "" {
		cached := BibleContent{
			Content:   passage_text,
			Copyright: copy_right,
		}
		val, err := json.Marshal(cached)
		if err == nil {
			err = redisClient.Set(context.Background(), cacheKey, string(val), 7*24*time.Hour).Err()
			if err != nil {
				log.Printf("[get_bible_passage] Error caching to Redis: %v", err)
			} else {
				log.Printf("[get_bible_passage] Cached bible passage for %s", passage)
			}
		}
	}
	return passage_text, copy_right, nil
}

func passage_sanitize(passage string) string {
	passage = strings.ReplaceAll(passage, "—", "-")
	passage = strings.ReplaceAll(passage, ";", ",")
	passage = strings.ReplaceAll(passage, " and ", ",")
	passage = strings.ReplaceAll(passage, ", ", ",")
	passage = strings.ReplaceAll(passage, ": ", ":")
	passage = strings.ReplaceAll(passage, " :", ":")
	passage = strings.ReplaceAll(passage, " : ", ":")
	passage = strings.ReplaceAll(passage, "- ", "-")
	passage = strings.ReplaceAll(passage, " -", "-")
	passage = strings.ReplaceAll(passage, " - ", "-")
	return passage
}

func get_bible_passage_text(passage string) (BibleContent, error) {
	settings := GetSettings()
	api_key := settings.BibleAPIKey
	if api_key == "" {
		return BibleContent{}, fmt.Errorf("BIBLE_API_KEY not set in database settings")
	}

	bible_id := os.Getenv("BIBLE_ID")
	if bible_id == "" {
		return BibleContent{}, fmt.Errorf("BIBLE_ID not set")
	}

	url := fmt.Sprintf("https://rest.api.bible/v1/bibles/%s/passages/%s", bible_id, passage)
	log.Printf("[get_bible_passage_text] URL: %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return BibleContent{}, err
	}

	req.Header.Set("api-key", api_key)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return BibleContent{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return BibleContent{}, err
	}

	var bibleResponse BibleResponse
	if err := json.Unmarshal(body, &bibleResponse); err != nil {
		return BibleContent{}, err
	}

	return bibleResponse.Data, nil
}

func get_bible_book_short(passage_book string) string {
	cleanBook := strings.TrimSpace(passage_book)
	norm := strings.ToLower(cleanBook)

	// Psalm aliases
	if norm == "ps" || norm == "psalm" || norm == "psalms" || norm == "psa" {
		return "PSA"
	}

	for _, book := range bible_books.Data {
		if strings.EqualFold(book.Name, cleanBook) ||
			strings.EqualFold(book.NameLong, cleanBook) ||
			strings.EqualFold(book.Abbreviation, cleanBook) ||
			strings.EqualFold(book.ID, cleanBook) {
			return book.Abbreviation
		}
	}

	// Additional common book aliases
	switch norm {
	case "matt", "mt", "matthew":
		return "MAT"
	case "mk", "mark":
		return "MRK"
	case "lk", "luke":
		return "LUK"
	case "jn", "john":
		return "JHN"
	case "act", "acts":
		return "ACT"
	case "rom", "romans":
		return "ROM"
	case "1 cor", "1 corinthians", "1cor":
		return "1CO"
	case "2 cor", "2 corinthians", "2cor":
		return "2CO"
	case "gal", "galatians":
		return "GAL"
	case "eph", "ephesians":
		return "EPH"
	case "phil", "philippians":
		return "PHP"
	case "col", "colossians":
		return "COL"
	case "gen", "genesis":
		return "GEN"
	case "ex", "exod", "exodus":
		return "EXO"
	}

	return ""
}


func strip_verse_suffix(v string) string {
	re := regexp.MustCompile(`[a-zA-Z]+$`)
	return re.ReplaceAllString(v, "")
}
