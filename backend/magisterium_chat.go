package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type MagisteriumChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type MagisteriumChatRequest struct {
	Messages []MagisteriumChatMessage `json:"messages"`
}

type MagisteriumCitation struct {
	Title         string  `json:"title"`
	DocumentTitle string  `json:"document_title,omitempty"`
	Author        string  `json:"author,omitempty"`
	Ref           string  `json:"ref,omitempty"`
	Score         float64 `json:"score,omitempty"`
	URL           string  `json:"url,omitempty"`
	Text          string  `json:"text,omitempty"`
}

type MagisteriumUsage struct {
	Limit     string `json:"limit,omitempty"`
	Remaining string `json:"remaining,omitempty"`
	Reset     string `json:"reset,omitempty"`
}

type MagisteriumChatResponse struct {
	Response  string                `json:"response"`
	Citations []MagisteriumCitation `json:"citations,omitempty"`
	Usage     *MagisteriumUsage     `json:"usage,omitempty"`
	Error     string                `json:"error,omitempty"`
}

func magisteriumChatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req MagisteriumChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if len(req.Messages) == 0 {
		respondWithError(w, http.StatusBadRequest, "No messages provided")
		return
	}

	settings := GetSettings()
	apiKey := settings.MagisteriumAPIKey
	if apiKey == "" {
		apiKey = os.Getenv("MAGISTERIUM_API_KEY")
	}

	if apiKey == "" {
		respondWithError(w, http.StatusBadRequest, "Magisterium AI API Key is not configured in Admin Settings.")
		return
	}

	// Prepare query: if multiple conversation turns exist, format past turn context into query
	var queryBuilder strings.Builder
	if len(req.Messages) > 1 {
		queryBuilder.WriteString("Conversation Context:\n")
		// Include up to 6 previous turns for context
		startIdx := 0
		if len(req.Messages) > 7 {
			startIdx = len(req.Messages) - 7
		}
		for i := startIdx; i < len(req.Messages)-1; i++ {
			msg := req.Messages[i]
			roleLabel := "User"
			if msg.Role == "assistant" {
				roleLabel = "Assistant"
			}
			queryBuilder.WriteString(fmt.Sprintf("%s: %s\n", roleLabel, msg.Content))
		}
		queryBuilder.WriteString("\nFollow-up Question: ")
		queryBuilder.WriteString(req.Messages[len(req.Messages)-1].Content)
	} else {
		queryBuilder.WriteString(req.Messages[len(req.Messages)-1].Content)
	}

	fullQuery := queryBuilder.String()

	// Call Magisterium API endpoint
	magisteriumURL := os.Getenv("MAGISTERIUM_API_URL")
	if magisteriumURL == "" {
		magisteriumURL = "https://www.magisterium.com/api/v1/search"
	}

	payload := map[string]interface{}{
		"query":      fullQuery,
		"messages":   req.Messages,
		"numResults": 5,
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to marshal request payload")
		return
	}

	client := http.Client{Timeout: 30 * time.Second}
	httpReq, err := http.NewRequest(http.MethodPost, magisteriumURL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create HTTP request to Magisterium AI")
		return
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	resp, err := client.Do(httpReq)
	if err != nil {
		log.Printf("[magisteriumChatHandler] HTTP error calling Magisterium AI: %v", err)
		respondWithError(w, http.StatusBadGateway, "Failed to connect to Magisterium AI service")
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to read response from Magisterium AI")
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("[magisteriumChatHandler] Magisterium API error response (Status %d): %s", resp.StatusCode, string(respBody))
		var errData map[string]interface{}
		errMsg := fmt.Sprintf("Magisterium AI request failed with status %d", resp.StatusCode)
		if json.Unmarshal(respBody, &errData) == nil {
			if msg, ok := errData["message"].(string); ok {
				errMsg = msg
			} else if msg, ok := errData["error"].(string); ok {
				errMsg = msg
			}
		}
		respondWithError(w, resp.StatusCode, errMsg)
		return
	}

	log.Printf("[magisteriumChatHandler] Raw Magisterium AI Response: %s", string(respBody))

	var rawResp map[string]interface{}
	if err := json.Unmarshal(respBody, &rawResp); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to parse Magisterium AI response")
		return
	}

	var citations []MagisteriumCitation
	var answerText string

	// Extract results according to target schema { data: { results: [...] } }
	var resultsList []interface{}

	if dataObj, ok := rawResp["data"].(map[string]interface{}); ok {
		if resArr, ok := dataObj["results"].([]interface{}); ok {
			resultsList = resArr
		}
		if ans, ok := dataObj["answer"].(string); ok && strings.TrimSpace(ans) != "" {
			answerText = ans
		}
	}

	// Fallback to top-level results if data wrapper is omitted
	if len(resultsList) == 0 {
		if resArr, ok := rawResp["results"].([]interface{}); ok {
			resultsList = resArr
		}
	}

	if answerText == "" {
		if ans, ok := rawResp["answer"].(string); ok && strings.TrimSpace(ans) != "" {
			answerText = ans
		}
	}

	// Process schema items from data.results
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

	// If no explicit summary answer was provided, build formatted response synthesis from results
	if strings.TrimSpace(answerText) == "" && len(citations) > 0 {
		var parts []string
		for idx, c := range citations {
			header := fmt.Sprintf("**%d. %s**", idx+1, c.DocumentTitle)
			if c.Author != "" {
				header += fmt.Sprintf(" — *%s*", c.Author)
			}
			if c.Ref != "" {
				header += fmt.Sprintf(" (%s)", c.Ref)
			}
			parts = append(parts, fmt.Sprintf("%s\n%s", header, c.Text))
		}
		answerText = strings.Join(parts, "\n\n")
	}

	if strings.TrimSpace(answerText) == "" {
		answerText = string(respBody)
	}

	// Extract rate limit usage headers if provided by Magisterium API
	limitHeader := resp.Header.Get("X-RateLimit-Limit")
	if limitHeader == "" {
		limitHeader = resp.Header.Get("x-ratelimit-limit")
	}
	remHeader := resp.Header.Get("X-RateLimit-Remaining")
	if remHeader == "" {
		remHeader = resp.Header.Get("x-ratelimit-remaining")
	}
	resetHeader := resp.Header.Get("X-RateLimit-Reset")
	if resetHeader == "" {
		resetHeader = resp.Header.Get("x-ratelimit-reset")
	}

	var usageInfo *MagisteriumUsage
	if limitHeader != "" || remHeader != "" || resetHeader != "" {
		usageInfo = &MagisteriumUsage{
			Limit:     limitHeader,
			Remaining: remHeader,
			Reset:     resetHeader,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(MagisteriumChatResponse{
		Response:  answerText,
		Citations: citations,
		Usage:     usageInfo,
	})
}

type SavePublicConversationRequest struct {
	ID         string          `json:"id"`
	AuthorName string          `json:"author_name"`
	Title      string          `json:"title"`
	Messages   json.RawMessage `json:"messages"`
}

func publicConversationsSaveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SavePublicConversationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if strings.TrimSpace(req.AuthorName) == "" {
		respondWithError(w, http.StatusBadRequest, "Author name is required")
		return
	}

	if strings.TrimSpace(req.ID) == "" || len(req.Messages) == 0 {
		respondWithError(w, http.StatusBadRequest, "Conversation ID and messages are required")
		return
	}

	if err := SavePublicConversation(req.ID, req.AuthorName, req.Title, req.Messages); err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to publish conversation: %v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Conversation successfully published to public sanctuary database!",
	})
}

func publicConversationsListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	list, err := GetPublicConversations()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve public conversations")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"conversations": list,
	})
}
