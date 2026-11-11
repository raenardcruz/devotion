package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// MagisteriumRAGResponse represents the structured response returned by the Magisterium RAG Agent.
type MagisteriumRAGResponse struct {
	Query            string                `json:"query"`
	TotalCitations   int                   `json:"total_citations"`
	Citations        []MagisteriumCitation `json:"citations"`
	FormattedContext string                `json:"formatted_context"`
	Usage            *MagisteriumUsage     `json:"usage,omitempty"`
	Error            string                `json:"error,omitempty"`
}

// MagisteriumRAGRequest represents the input request to the Magisterium RAG Agent.
type MagisteriumRAGRequest struct {
	Query      string `json:"query"`
	NumResults int    `json:"num_results,omitempty"`
}

// ExecuteMagisteriumRAG queries the Magisterium AI API with a search query and returns structured citations.
func ExecuteMagisteriumRAG(ctx context.Context, query string, numResults int, apiKey string) ([]MagisteriumCitation, *MagisteriumUsage, error) {
	if apiKey == "" {
		settings := GetSettings()
		apiKey = settings.MagisteriumAPIKey
		if apiKey == "" {
			apiKey = os.Getenv("MAGISTERIUM_API_KEY")
		}
	}

	if apiKey == "" {
		return nil, nil, fmt.Errorf("Magisterium AI API key is not configured in Admin Settings or environment")
	}

	trimmedQuery := strings.TrimSpace(query)
	if trimmedQuery == "" {
		return nil, nil, fmt.Errorf("search query cannot be empty")
	}

	if numResults <= 0 {
		numResults = 5
	}

	magisteriumURL := os.Getenv("MAGISTERIUM_API_URL")
	if magisteriumURL == "" {
		magisteriumURL = "https://www.magisterium.com/api/v1/search"
	}

	payload := map[string]interface{}{
		"query":      trimmedQuery,
		"numResults": numResults,
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal Magisterium request payload: %w", err)
	}

	client := http.Client{Timeout: 30 * time.Second}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, magisteriumURL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create HTTP request to Magisterium AI: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	resp, err := client.Do(httpReq)
	if err != nil {
		log.Printf("[ExecuteMagisteriumRAG] HTTP error: %v", err)
		return nil, nil, fmt.Errorf("failed to connect to Magisterium AI service: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read response from Magisterium AI: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("[ExecuteMagisteriumRAG] Error response (Status %d): %s", resp.StatusCode, string(respBody))
		var errData map[string]interface{}
		errMsg := fmt.Sprintf("Magisterium AI request failed with status %d", resp.StatusCode)
		if json.Unmarshal(respBody, &errData) == nil {
			if msg, ok := errData["message"].(string); ok {
				errMsg = msg
			} else if msg, ok := errData["error"].(string); ok {
				errMsg = msg
			}
		}
		return nil, nil, fmt.Errorf("%s", errMsg)
	}

	var rawResp map[string]interface{}
	if err := json.Unmarshal(respBody, &rawResp); err != nil {
		return nil, nil, fmt.Errorf("failed to parse Magisterium AI response: %w", err)
	}

	var resultsList []interface{}
	if dataObj, ok := rawResp["data"].(map[string]interface{}); ok {
		if resArr, ok := dataObj["results"].([]interface{}); ok {
			resultsList = resArr
		}
	}

	if len(resultsList) == 0 {
		if resArr, ok := rawResp["results"].([]interface{}); ok {
			resultsList = resArr
		}
	}

	var citations []MagisteriumCitation
	for _, item := range resultsList {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		author, _ := itemMap["author"].(string)
		docTitle, _ := itemMap["document_title"].(string)
		refStr, _ := itemMap["ref"].(string)
		textVal, _ := itemMap["text"].(string)
		titleVal, _ := itemMap["title"].(string)
		urlVal, _ := itemMap["url"].(string)
		scoreVal, _ := itemMap["score"].(float64)

		displayTitle := titleVal
		if displayTitle == "" {
			displayTitle = docTitle
		}
		if displayTitle == "" {
			displayTitle = refStr
		}

		citations = append(citations, MagisteriumCitation{
			Title:         displayTitle,
			DocumentTitle: docTitle,
			Author:        author,
			Ref:           refStr,
			Score:         scoreVal,
			URL:           urlVal,
			Text:          textVal,
		})
	}

	limitHeader := resp.Header.Get("X-RateLimit-Limit")
	remHeader := resp.Header.Get("X-RateLimit-Remaining")
	resetHeader := resp.Header.Get("X-RateLimit-Reset")

	var usageInfo *MagisteriumUsage
	if limitHeader != "" || remHeader != "" || resetHeader != "" {
		usageInfo = &MagisteriumUsage{
			Limit:     limitHeader,
			Remaining: remHeader,
			Reset:     resetHeader,
		}
	}

	log.Printf("[ExecuteMagisteriumRAG] Successfully retrieved %d citations for query: %q", len(citations), trimmedQuery)
	return citations, usageInfo, nil
}

// FormatCitationsAsContext compiles citations into clean, structured text for LLM grounding.
func FormatCitationsAsContext(citations []MagisteriumCitation) string {
	if len(citations) == 0 {
		return "No citations found."
	}

	var sb strings.Builder
	for i, c := range citations {
		docName := c.DocumentTitle
		if docName == "" {
			docName = c.Title
		}
		sb.WriteString(fmt.Sprintf("[%d] %s", i+1, docName))
		if c.Author != "" {
			sb.WriteString(fmt.Sprintf(" (Author: %s)", c.Author))
		}
		if c.Ref != "" {
			sb.WriteString(fmt.Sprintf(" [Reference: %s]", c.Ref))
		}
		if c.Score > 0 {
			sb.WriteString(fmt.Sprintf(" [Match: %.1f%%]", c.Score*100))
		}
		sb.WriteString("\n")
		if c.Text != "" {
			sb.WriteString(fmt.Sprintf("Excerpt: %s\n\n", c.Text))
		}
	}
	return strings.TrimSpace(sb.String())
}

// magisteriumRAGHandler exposes the Magisterium RAG Agent as an HTTP REST endpoint.
func magisteriumRAGHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req MagisteriumRAGRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if strings.TrimSpace(req.Query) == "" {
		respondWithError(w, http.StatusBadRequest, "Search query is required")
		return
	}

	settings := GetSettings()
	citations, usage, err := ExecuteMagisteriumRAG(r.Context(), req.Query, req.NumResults, settings.MagisteriumAPIKey)
	if err != nil {
		log.Printf("[magisteriumRAGHandler] RAG execution error: %v", err)
		respondWithError(w, http.StatusBadGateway, fmt.Sprintf("Magisterium RAG error: %v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(MagisteriumRAGResponse{
		Query:            req.Query,
		TotalCitations:   len(citations),
		Citations:        citations,
		FormattedContext: FormatCitationsAsContext(citations),
		Usage:            usage,
	})
}
