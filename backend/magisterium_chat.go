package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

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
	Stream   bool                     `json:"stream,omitempty"`
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

	isStreaming := req.Stream || strings.Contains(r.Header.Get("Accept"), "text/event-stream")
	var flusher http.Flusher
	if isStreaming {
		var ok bool
		flusher, ok = w.(http.Flusher)
		if !ok {
			isStreaming = false
		}
	}

	if apiKey == "" {
		if isStreaming {
			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Connection", "keep-alive")
			w.Header().Set("X-Accel-Buffering", "no")
			data, _ := json.Marshal(map[string]any{
				"type":  "error",
				"error": "Magisterium AI API Key is not configured in Admin Settings.",
			})
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
			return
		}
		respondWithError(w, http.StatusBadRequest, "Magisterium AI API Key is not configured in Admin Settings.")
		return
	}

	if isStreaming {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("X-Accel-Buffering", "no")
		flusher.Flush()
	}

	sendSSE := func(payload any) error {
		if !isStreaming || flusher == nil {
			return nil
		}
		data, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, "data: %s\n\n", data); err != nil {
			return err
		}
		flusher.Flush()
		return nil
	}

	llmProvider := settings.MagisteriumLLMProvider
	if llmProvider == "" {
		llmProvider = "ollama"
	}

	// Prepare query: if multiple conversation turns exist, generate a standalone search prompt via LLM using conversation history (excluding documents)
	var searchQuery string
	latestQuestion := req.Messages[len(req.Messages)-1].Content

	if len(req.Messages) > 1 {
		if isStreaming {
			_ = sendSSE(map[string]any{
				"type":    "status",
				"message": "Understanding conversation context...",
			})
		}
		log.Printf("[magisteriumChatHandler] Follow-up detected (%d messages). Generating standalone query via LLM (%s)...", len(req.Messages), llmProvider)
		generatedQuery, err := generateStandaloneQueryWithLLM(r.Context(), req.Messages, llmProvider, settings)
		if err != nil {
			log.Printf("[magisteriumChatHandler] Warning: failed to generate standalone query via LLM: %v. Falling back to latest user question.", err)
			searchQuery = latestQuestion
		} else if strings.TrimSpace(generatedQuery) != "" {
			searchQuery = generatedQuery
			log.Printf("[magisteriumChatHandler] Generated standalone query for Magisterium API: %s", searchQuery)
		} else {
			searchQuery = latestQuestion
		}
	} else {
		searchQuery = latestQuestion
	}

	if isStreaming {
		_ = sendSSE(map[string]any{
			"type":    "status",
			"message": "Searching Catholic Church documents...",
		})
	}

	// Call Magisterium API endpoint via Magisterium RAG Agent
	citations, usageInfo, err := ExecuteMagisteriumRAG(r.Context(), searchQuery, 5, apiKey)
	if err != nil {
		log.Printf("[magisteriumChatHandler] Magisterium RAG error: %v", err)
		if isStreaming {
			_ = sendSSE(map[string]any{
				"type":  "error",
				"error": fmt.Sprintf("Failed to query Magisterium AI: %v", err),
			})
			return
		}
		respondWithError(w, http.StatusBadGateway, fmt.Sprintf("Failed to query Magisterium AI: %v", err))
		return
	}

	if isStreaming {
		_ = sendSSE(map[string]any{
			"type":      "citations",
			"citations": citations,
			"usage":     usageInfo,
		})
	}

	var answerText string

	// Default to LLM summary mode (or if explicitly requested as "llm_summary"), summarize citations with configured LLM (Ollama or Gemini)
	if (req.Mode == "" || req.Mode == "llm_summary") && len(citations) > 0 {
		if isStreaming {
			_ = sendSSE(map[string]any{
				"type":    "status",
				"message": "Synthesizing Catholic response...",
			})
			summary, err := streamCitationsSummaryWithLLM(r.Context(), latestQuestion, citations, llmProvider, settings, func(delta string) error {
				return sendSSE(map[string]any{
					"type":  "chunk",
					"delta": delta,
				})
			})
			if err != nil {
				log.Printf("[magisteriumChatHandler] Warning: streaming custom LLM summary failed: %v", err)
				if r.Context().Err() != nil {
					return
				}
				// Fallback to formatted citations
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
				_ = sendSSE(map[string]any{
					"type":  "chunk",
					"delta": answerText,
				})
				_ = sendSSE(map[string]any{
					"type":     "done",
					"response": answerText,
				})
				fmt.Fprintf(w, "data: [DONE]\n\n")
				flusher.Flush()
				return
			}
			answerText = summary
			_ = sendSSE(map[string]any{
				"type":     "done",
				"response": answerText,
			})
			fmt.Fprintf(w, "data: [DONE]\n\n")
			flusher.Flush()
			return
		}

		summary, err := summarizeCitationsWithLLM(r.Context(), latestQuestion, citations, llmProvider, settings)
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
		answerText = "No response available from Magisterium AI."
	}

	if isStreaming {
		_ = sendSSE(map[string]any{
			"type":  "chunk",
			"delta": answerText,
		})
		_ = sendSSE(map[string]any{
			"type":     "done",
			"response": answerText,
		})
		fmt.Fprintf(w, "data: [DONE]\n\n")
		flusher.Flush()
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(MagisteriumChatResponse{
		Response:  answerText,
		Citations: citations,
		Usage:     usageInfo,
	})
}

// streamCitationsSummaryWithLLM feeds Magisterium search citations into the selected LLM provider (Ollama or Gemini)
// to generate a synthesis summary, streaming chunks via onChunk if provided.
func streamCitationsSummaryWithLLM(
	ctx context.Context,
	userQuestion string,
	citations []MagisteriumCitation,
	provider string,
	settings Settings,
	onChunk func(delta string) error,
) (string, error) {
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
	streamEnabled := onChunk != nil

	if provider == "gemini" {
		apiKey := settings.GeminiAPIKey
		if apiKey == "" {
			return "", fmt.Errorf("Gemini API key not configured in Admin Settings")
		}
		selectedModel := settings.GeminiModel
		if selectedModel == "" {
			selectedModel = "gemini-3.1-flash-lite"
		}
		log.Printf("[streamCitationsSummaryWithLLM] Summarizing Magisterium citations with Gemini (%s) [streaming=%t]...", selectedModel, streamEnabled)
		gModel, err := geminiModel.NewModel(ctx, selectedModel, &genai.ClientConfig{
			APIKey: apiKey,
		})
		if err != nil {
			return "", fmt.Errorf("failed to initialize Gemini model: %w", err)
		}
		seq := gModel.GenerateContent(ctx, req, streamEnabled)
		for resp, err := range seq {
			if err != nil {
				return "", fmt.Errorf("error generating summary with Gemini: %w", err)
			}
			if resp.Content != nil {
				for _, part := range resp.Content.Parts {
					if part.Text != "" {
						summaryBuilder.WriteString(part.Text)
						if onChunk != nil {
							if err := onChunk(part.Text); err != nil {
								return summaryBuilder.String(), err
							}
						}
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
		log.Printf("[streamCitationsSummaryWithLLM] Summarizing Magisterium citations with Ollama (%s) [streaming=%t]...", selectedModel, streamEnabled)
		oModel := ollama.NewModel(selectedModel, ollamaURL)
		seq := oModel.GenerateContent(ctx, req, streamEnabled)
		for resp, err := range seq {
			if err != nil {
				return "", fmt.Errorf("error generating summary with Ollama: %w", err)
			}
			if resp.Content != nil {
				for _, part := range resp.Content.Parts {
					if part.Text != "" {
						summaryBuilder.WriteString(part.Text)
						if onChunk != nil {
							if err := onChunk(part.Text); err != nil {
								return summaryBuilder.String(), err
							}
						}
					}
				}
			}
		}
	}

	return strings.TrimSpace(summaryBuilder.String()), nil
}

// summarizeCitationsWithLLM feeds Magisterium search citations into the selected LLM provider (Ollama or Gemini) to generate a synthesis summary.
func summarizeCitationsWithLLM(ctx context.Context, userQuestion string, citations []MagisteriumCitation, provider string, settings Settings) (string, error) {
	return streamCitationsSummaryWithLLM(ctx, userQuestion, citations, provider, settings, nil)
}

// generateStandaloneQueryWithLLM uses the configured LLM provider (Ollama or Gemini) to process conversation history
// and reformulate a follow-up user message into a concise standalone search query for Magisterium AI.
func generateStandaloneQueryWithLLM(ctx context.Context, messages []MagisteriumChatMessage, provider string, settings Settings) (string, error) {
	if len(messages) == 0 {
		return "", nil
	}

	var chatHistory strings.Builder
	// Include recent turns up to 6 for history context, stripping any document attachments
	startIdx := 0
	if len(messages) > 7 {
		startIdx = len(messages) - 7
	}
	for i := startIdx; i < len(messages)-1; i++ {
		msg := messages[i]
		roleLabel := "User"
		if msg.Role == "assistant" {
			roleLabel = "Assistant"
		}
		// Truncate message text if exceptionally long to conserve prompt budget
		content := msg.Content
		if len(content) > 1000 {
			content = content[:1000] + "... [truncated]"
		}
		chatHistory.WriteString(fmt.Sprintf("%s: %s\n\n", roleLabel, content))
	}

	latestQuestion := messages[len(messages)-1].Content

	prompt := fmt.Sprintf(`Given the following conversation history and a follow-up question, reformulate the follow-up question into a clear, standalone search query to look up Catholic magisterial documents.
Do NOT answer the question. Do NOT include any explanations, document text, preamble, or commentary. Output ONLY the standalone search query string.

Conversation History:
%s

Follow-up Question: %s

Standalone Search Query:`, chatHistory.String(), latestQuestion)

	req := &model.LLMRequest{
		Contents: []*genai.Content{
			genai.NewContentFromText(prompt, "user"),
		},
	}

	var resultBuilder strings.Builder

	if provider == "gemini" {
		apiKey := settings.GeminiAPIKey
		if apiKey == "" {
			return "", fmt.Errorf("Gemini API key not configured in Admin Settings")
		}
		selectedModel := settings.GeminiModel
		if selectedModel == "" {
			selectedModel = "gemini-3.1-flash-lite"
		}
		gModel, err := geminiModel.NewModel(ctx, selectedModel, &genai.ClientConfig{
			APIKey: apiKey,
		})
		if err != nil {
			return "", fmt.Errorf("failed to initialize Gemini model for query generation: %w", err)
		}
		seq := gModel.GenerateContent(ctx, req, false)
		for resp, err := range seq {
			if err != nil {
				return "", fmt.Errorf("error generating query with Gemini: %w", err)
			}
			if resp.Content != nil {
				for _, part := range resp.Content.Parts {
					if part.Text != "" {
						resultBuilder.WriteString(part.Text)
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
		oModel := ollama.NewModel(selectedModel, ollamaURL)
		seq := oModel.GenerateContent(ctx, req, false)
		for resp, err := range seq {
			if err != nil {
				return "", fmt.Errorf("error generating query with Ollama: %w", err)
			}
			if resp.Content != nil {
				for _, part := range resp.Content.Parts {
					if part.Text != "" {
						resultBuilder.WriteString(part.Text)
					}
				}
			}
		}
	}

	query := strings.TrimSpace(resultBuilder.String())
	// Strip quotes if the LLM wrapped the query in quotation marks
	query = strings.Trim(query, `"'` + "`")
	return query, nil
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
	if r.Method == http.MethodDelete {
		publicConversationsDeleteHandler(w, r)
		return
	}
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

type DeletePublicConversationRequest struct {
	ID string `json:"id"`
}

func publicConversationsDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if id == "" && r.Body != nil {
		var req DeletePublicConversationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err == nil {
			id = strings.TrimSpace(req.ID)
		}
	}

	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Conversation ID is required")
		return
	}

	if err := DeletePublicConversation(id); err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to delete public conversation: %v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Conversation successfully removed from public sanctuary database!",
	})
}
