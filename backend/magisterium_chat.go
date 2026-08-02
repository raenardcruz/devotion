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

	"github.com/raenardcruz/ollama-adk-wrapper"
	"google.golang.org/adk/model"
	geminiModel "google.golang.org/adk/model/gemini"
	"google.golang.org/genai"
)

type MagisteriumChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type MagisteriumChatRequest struct {
	Messages []MagisteriumChatMessage `json:"messages"`
	Mode     string                   `json:"mode,omitempty"` // "magisterium" or "llm_summary"
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

	// If user requested LLM summary mode (or if answerText is empty in llm_summary mode), summarize citations with configured LLM (Ollama or Gemini)
	if req.Mode == "llm_summary" && len(citations) > 0 {
		llmProvider := settings.MagisteriumLLMProvider
		if llmProvider == "" {
			llmProvider = "ollama"
		}
		userQuestion := req.Messages[len(req.Messages)-1].Content
		summary, err := summarizeCitationsWithLLM(r.Context(), userQuestion, citations, llmProvider, settings)
		if err != nil {
			log.Printf("[magisteriumChatHandler] Warning: custom LLM summary failed: %v. Falling back to default answer/citations.", err)
		} else if strings.TrimSpace(summary) != "" {
			answerText = summary
		}
	} else if strings.TrimSpace(answerText) == "" && len(citations) > 0 {
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

// summarizeCitationsWithLLM feeds Magisterium search citations into the selected LLM provider (Ollama or Gemini) to generate a synthesis summary.
func summarizeCitationsWithLLM(ctx context.Context, userQuestion string, citations []MagisteriumCitation, provider string, settings Settings) (string, error) {
	var citationText strings.Builder
	for i, c := range citations {
		citationText.WriteString(fmt.Sprintf("[%d] %s", i+1, c.DocumentTitle))
		if c.Author != "" {
			citationText.WriteString(fmt.Sprintf(" (Author: %s)", c.Author))
		}
		if c.Ref != "" {
			citationText.WriteString(fmt.Sprintf(" [Ref: %s]", c.Ref))
		}
		citationText.WriteString("\n")
		if c.Text != "" {
			citationText.WriteString(fmt.Sprintf("Excerpt: %s\n\n", c.Text))
		}
	}

	prompt := fmt.Sprintf(`You are a knowledgeable and faithful Catholic assistant.
The user asked: "%s"

Below are official search citations retrieved from Catholic Magisterial documents (Catechism, Encyclicals, Scripture, Council documents):

%s

Instructions:
1. Provide a clear, cohesive, and comprehensive summary answering the user's question using ONLY the provided citations above.
2. Maintain fidelity to Catholic doctrine. Refer back to the cited documents or authors inline (e.g. [1], [2] or naming the document) when summarizing.
3. Structure your response with clean formatting (using markdown bolding and bullet points where helpful).`, userQuestion, citationText.String())

	req := &model.LLMRequest{
		Contents: []*genai.Content{
			genai.NewContentFromText(prompt, "user"),
		},
	}

	var summaryBuilder strings.Builder

	if provider == "gemini" {
		apiKey := settings.GeminiAPIKey
		if apiKey == "" {
			return "", fmt.Errorf("Gemini API key not configured in Admin Settings")
		}
		selectedModel := settings.GeminiModel
		if selectedModel == "" {
			selectedModel = "gemini-3.1-flash-lite"
		}
		log.Printf("[summarizeCitationsWithLLM] Summarizing Magisterium citations with Gemini (%s)...", selectedModel)
		gModel, err := geminiModel.NewModel(ctx, selectedModel, &genai.ClientConfig{
			APIKey: apiKey,
		})
		if err != nil {
			return "", fmt.Errorf("failed to initialize Gemini model: %w", err)
		}
		seq := gModel.GenerateContent(ctx, req, false)
		for resp, err := range seq {
			if err != nil {
				return "", fmt.Errorf("error generating summary with Gemini: %w", err)
			}
			if resp.Content != nil {
				for _, part := range resp.Content.Parts {
					if part.Text != "" {
						summaryBuilder.WriteString(part.Text)
					}
				}
			}
		}
	} else {
		// Default to Ollama
		ollamaURL := os.Getenv("OLLAMA_URL")
		if ollamaURL == "" {
			ollamaURL = "http://localhost:11434"
		}
		selectedModel := settings.OllamaModel
		if selectedModel == "" {
			selectedModel = "gemma4:cloud"
		}
		log.Printf("[summarizeCitationsWithLLM] Summarizing Magisterium citations with Ollama (%s)...", selectedModel)
		oModel := ollama.NewModel(selectedModel, ollamaURL)
		seq := oModel.GenerateContent(ctx, req, false)
		for resp, err := range seq {
			if err != nil {
				return "", fmt.Errorf("error generating summary with Ollama: %w", err)
			}
			if resp.Content != nil {
				for _, part := range resp.Content.Parts {
					if part.Text != "" {
						summaryBuilder.WriteString(part.Text)
					}
				}
			}
		}
	}

	return strings.TrimSpace(summaryBuilder.String()), nil
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
