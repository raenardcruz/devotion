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
	"regexp"
	"strings"
	"time"

	"github.com/raenardcruz/ollama-adk-wrapper"
	"google.golang.org/adk/model"
	geminiModel "google.golang.org/adk/model/gemini"
	"google.golang.org/genai"
)

type BibleContext struct {
	Citation  string                `json:"citation"`
	Context   string                `json:"context"`
	Citations []MagisteriumCitation `json:"citations,omitempty"`
}

type GetBibleContextArgs struct {
	Citation string `json:"citation"`
}

func get_bible_context(ctx context.Context, args GetBibleContextArgs) (BibleContext, error) {
	if strings.Contains(args.Citation, "Psalms") {
		return BibleContext{
			Citation: args.Citation,
			Context:  "",
		}, nil
	}

	passageText, _, err := get_bible_passage(args.Citation)
	if err != nil {
		log.Printf("[get_bible_context] Warning: could not fetch passage text for %s: %v", args.Citation, err)
	}

	settings := GetSettings()
	instructionTemplate := settings.ContextInstruction
	if instructionTemplate == "" {
		instructionTemplate = defaultInstruction
	}

	// Substitute placeholders in instruction
	prompt := strings.ReplaceAll(instructionTemplate, "{{citation}}", args.Citation)
	prompt = strings.ReplaceAll(prompt, "%s", args.Citation)
	if passageText != "" {
		passageFormatted := fmt.Sprintf("Passage Text:\n%s", passageText)
		prompt = strings.ReplaceAll(prompt, "{{passage_text}}", passageFormatted)
	} else {
		prompt = strings.ReplaceAll(prompt, "{{passage_text}}", "")
	}

	req := &model.LLMRequest{
		Contents: []*genai.Content{
			genai.NewContentFromText(prompt, "user"),
		},
	}

	var contextText strings.Builder

	if settings.ContextProvider == "gemini" {
		apiKey := settings.GeminiAPIKey
		if apiKey == "" {
			return BibleContext{}, fmt.Errorf("Gemini API key not configured in settings")
		}
		selectedModelName := settings.GeminiModel
		if selectedModelName == "" {
			selectedModelName = "gemini-3.1-flash-lite"
		}
		log.Printf("[get_bible_context] Generating draft bible context for %s with Gemini (%s)...\n", args.Citation, selectedModelName)
		gModel, err := geminiModel.NewModel(ctx, selectedModelName, &genai.ClientConfig{
			APIKey: apiKey,
		})
		if err != nil {
			return BibleContext{}, fmt.Errorf("failed to create gemini model: %w", err)
		}
		seq := gModel.GenerateContent(ctx, req, false)
		for resp, err := range seq {
			if err != nil {
				return BibleContext{}, fmt.Errorf("error generating bible context with gemini: %w", err)
			}
			if resp.Content != nil {
				for _, part := range resp.Content.Parts {
					if part.Text != "" {
						contextText.WriteString(part.Text)
					}
				}
			}
		}
	} else {
		// Context generation defaults to Ollama
		ollamaURL := os.Getenv("OLLAMA_URL")
		if ollamaURL == "" {
			ollamaURL = "http://localhost:11434"
		}
		selectedModelName := settings.OllamaModel
		if selectedModelName == "" {
			selectedModelName = "gemma4:cloud"
		}

		log.Printf("[get_bible_context] Generating draft bible context for %s with Ollama (%s)...\n", args.Citation, selectedModelName)
		oModel := ollama.NewModel(selectedModelName, ollamaURL)
		seq := oModel.GenerateContent(ctx, req, false)

		for resp, err := range seq {
			if err != nil {
				return BibleContext{}, fmt.Errorf("error generating bible context with ollama: %w", err)
			}
			if resp.Content != nil {
				for _, part := range resp.Content.Parts {
					if part.Text != "" {
						contextText.WriteString(part.Text)
					}
				}
			}
		}
	}

	rawContext := contextText.String()

	if !settings.EnableFactChecker {
		log.Printf("[get_bible_context] Fact checker is disabled in settings. Returning raw draft context for %s.\n", args.Citation)
		return BibleContext{
			Citation:  args.Citation,
			Context:   rawContext,
			Citations: extractCitationsFromContext(rawContext),
		}, nil
	}

	log.Printf("[get_bible_context] Generated draft bible context for %s successfully. Running fact-check & correction...\n", args.Citation)

	factCheckedContext, citations, err := FactCheckBibleContext(ctx, args.Citation, rawContext, passageText)
	if err != nil {
		log.Printf("[get_bible_context] Warning: fact check failed for %s: %v. Returning raw context.", args.Citation, err)
		factCheckedContext = rawContext
	}

	if len(citations) == 0 {
		citations = extractCitationsFromContext(factCheckedContext)
	}

	return BibleContext{
		Citation:  args.Citation,
		Context:   factCheckedContext,
		Citations: citations,
	}, nil
}

// FactCheckBibleContext fact-checks generated scripture/theological context against scripture text and Catholic doctrine guidelines, correcting any inaccuracies while strictly maintaining original format.
func FactCheckBibleContext(ctx context.Context, citation string, rawContext string, passageText string) (string, []MagisteriumCitation, error) {
	if strings.TrimSpace(rawContext) == "" {
		return rawContext, nil, nil
	}

	settings := GetSettings()
	factCheckPrompt := fmt.Sprintf(`You are a meticulous Catholic Bible scholar and editor.
Your task is to fact-check and correct the provided Bible Context for citation "%s".

Passage Text (Ground Truth):
%s

Draft Bible Context to Fact-Check:
%s

Fact-Checking Instructions:
1. Verify all historical, geographical, linguistic, and theological statements against Catholic Scripture and Church teaching.
2. Correct any factual inaccuracies, misattributed quotes, hallucinated Catechism (CCC) paragraph numbers, or misquoted scripture verses.
3. If an unverified exact CCC paragraph number is present, remove the specific number and refer to the general Catholic doctrine instead.
4. CRITICAL: Strictly preserve the exact tone, structure, paragraph layout, headings, and Markdown formatting (bold, italics, lists) of the draft. Do NOT alter the format or structure—only correct factual or citation errors.
5. Return ONLY the final corrected context text with no added commentary or metadata.`, citation, passageText, rawContext)

	req := &model.LLMRequest{
		Contents: []*genai.Content{
			genai.NewContentFromText(factCheckPrompt, "user"),
		},
	}

	var correctedText strings.Builder

	magisteriumApiKey := settings.MagisteriumAPIKey
	if magisteriumApiKey == "" {
		magisteriumApiKey = os.Getenv("MAGISTERIUM_API_KEY")
	}

	if settings.FactCheckerProvider == "magisterium" || magisteriumApiKey != "" {
		if magisteriumApiKey == "" {
			log.Printf("[FactCheckBibleContext] Warning: Magisterium API key not configured. Returning raw context.")
			return rawContext, extractCitationsFromContext(rawContext), nil
		}

		// Extract targeted text: search specifically for doctrine, CCC references, saints, and early Church writers from the raw context
		targetedQuery := fmt.Sprintf("Catholic Church teaching, Catechism of the Catholic Church (CCC), Saints, and Doctors of the Church referenced for passage %s: %s", citation, rawContext)
		if len(targetedQuery) > 800 {
			targetedQuery = targetedQuery[:800]
		}

		log.Printf("[FactCheckBibleContext] Querying Magisterium AI for targeted doctrine/CCC/saints verification on %s...", citation)

		magisteriumURL := os.Getenv("MAGISTERIUM_API_URL")
		if magisteriumURL == "" {
			magisteriumURL = "https://www.magisterium.com/api/v1/search"
		}

		payload := map[string]interface{}{
			"query":      targetedQuery,
			"numResults": 5,
		}

		jsonBytes, err := json.Marshal(payload)
		if err != nil {
			log.Printf("[FactCheckBibleContext] Error marshaling Magisterium payload: %v", err)
			return rawContext, extractCitationsFromContext(rawContext), nil
		}

		client := http.Client{Timeout: 20 * time.Second}
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, magisteriumURL, bytes.NewBuffer(jsonBytes))
		if err != nil {
			log.Printf("[FactCheckBibleContext] Error creating Magisterium request: %v", err)
			return rawContext, extractCitationsFromContext(rawContext), nil
		}

		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", magisteriumApiKey))

		resp, err := client.Do(httpReq)
		if err != nil {
			log.Printf("[FactCheckBibleContext] Error calling Magisterium AI: %v", err)
			return rawContext, extractCitationsFromContext(rawContext), nil
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil || resp.StatusCode != http.StatusOK {
			log.Printf("[FactCheckBibleContext] Magisterium API response error (Status %d): %s", resp.StatusCode, string(respBody))
			return rawContext, extractCitationsFromContext(rawContext), nil
		}

		var rawResp map[string]interface{}
		var retrievedCitations []string
		var retrievedCitationObjs []MagisteriumCitation
		if err := json.Unmarshal(respBody, &rawResp); err == nil {
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

			for _, item := range resultsList {
				if itemMap, ok := item.(map[string]interface{}); ok {
					docTitle, _ := itemMap["document_title"].(string)
					author, _ := itemMap["author"].(string)
					textVal, _ := itemMap["text"].(string)
					refVal, _ := itemMap["ref"].(string)
					titleVal, _ := itemMap["title"].(string)
					displayTitle := docTitle
					if displayTitle == "" {
						displayTitle = titleVal
					}
					retrievedCitations = append(retrievedCitations, fmt.Sprintf("- Document: %s | Author: %s | Ref: %s\nExcerpt: %s", displayTitle, author, refVal, textVal))
					retrievedCitationObjs = append(retrievedCitationObjs, MagisteriumCitation{
						Title:         displayTitle,
						DocumentTitle: displayTitle,
						Author:        author,
						Ref:           refVal,
						Text:          textVal,
					})
				}
			}
		}

		if len(retrievedCitations) == 0 {
			log.Printf("[FactCheckBibleContext] No specific Magisterium citations retrieved for verification. Returning raw context.")
			return rawContext, extractCitationsFromContext(rawContext), nil
		}

		// Use retrieved Magisterium teaching excerpts to refine and verify the specific doctrine/CCC/saints in the draft context
		magisteriumEvidence := strings.Join(retrievedCitations, "\n\n")
		verificationPrompt := fmt.Sprintf(`You are a Catholic theological editor. Verify and refine the specific Catholic teachings, Catechism (CCC) references, and Saints/Doctors of the Church mentioned in the draft context for "%s".

Authentic Magisterial Sources (Ground Truth Excerpts):
%s

Draft Context:
%s

Instructions:
1. Cross-reference the Catholic teachings, CCC references, and Saints mentioned in the draft context against the Magisterial sources above.
2. Correct any misquotations, wrong attributions, or hallucinated paragraph numbers.
3. Keep the original structure, headings, tone, and formatting of the draft context completely intact.
4. Return ONLY the final corrected context with no additional commentary.`, citation, magisteriumEvidence, rawContext)

		reqVerification := &model.LLMRequest{
			Contents: []*genai.Content{
				genai.NewContentFromText(verificationPrompt, "user"),
			},
		}

		// Use configured LLM (Ollama or Gemini) to apply the Magisterium evidence correction
		if settings.MagisteriumLLMProvider == "gemini" && settings.GeminiAPIKey != "" {
			gModel, err := geminiModel.NewModel(ctx, settings.GeminiModel, &genai.ClientConfig{APIKey: settings.GeminiAPIKey})
			if err == nil {
				seq := gModel.GenerateContent(ctx, reqVerification, false)
				for r, e := range seq {
					if e == nil && r.Content != nil {
						for _, p := range r.Content.Parts {
							correctedText.WriteString(p.Text)
						}
					}
				}
			}
		} else {
			ollamaURL := os.Getenv("OLLAMA_URL")
			if ollamaURL == "" {
				ollamaURL = "http://localhost:11434"
			}
			oModel := ollama.NewModel(settings.OllamaModel, ollamaURL)
			seq := oModel.GenerateContent(ctx, reqVerification, false)
			for r, e := range seq {
				if e == nil && r.Content != nil {
					for _, p := range r.Content.Parts {
						correctedText.WriteString(p.Text)
					}
				}
			}
		}

		finalContext := strings.TrimSpace(correctedText.String())
		if finalContext == "" {
			finalContext = rawContext
		}
		return finalContext, retrievedCitationObjs, nil

	} else if settings.FactCheckerProvider == "ollama" {
		ollamaURL := os.Getenv("OLLAMA_URL")
		if ollamaURL == "" {
			ollamaURL = "http://localhost:11434"
		}
		selectedModelName := settings.OllamaModel
		if selectedModelName == "" {
			selectedModelName = "gemma4:cloud"
		}
		log.Printf("[FactCheckBibleContext] Fact checking context for %s using Ollama (%s)...\n", citation, selectedModelName)
		oModel := ollama.NewModel(selectedModelName, ollamaURL)
		seq := oModel.GenerateContent(ctx, req, false)
		for resp, err := range seq {
			if err != nil {
				return rawContext, extractCitationsFromContext(rawContext), fmt.Errorf("error generating fact-check correction with ollama: %w", err)
			}
			if resp.Content != nil {
				for _, part := range resp.Content.Parts {
					if part.Text != "" {
						correctedText.WriteString(part.Text)
					}
				}
			}
		}
	} else {
		// Fact checking default: Gemini
		apiKey := settings.GeminiAPIKey
		if apiKey == "" {
			log.Printf("[FactCheckBibleContext] Warning: Gemini API key not configured in settings. Returning raw context.")
			return rawContext, extractCitationsFromContext(rawContext), nil
		}
		selectedModelName := settings.GeminiModel
		if selectedModelName == "" {
			selectedModelName = "gemini-3.1-flash-lite"
		}

		log.Printf("[FactCheckBibleContext] Fact checking context for %s using Gemini (%s)...\n", citation, selectedModelName)
		gModel, err := geminiModel.NewModel(ctx, selectedModelName, &genai.ClientConfig{
			APIKey: apiKey,
		})
		if err != nil {
			return rawContext, extractCitationsFromContext(rawContext), fmt.Errorf("failed to create gemini model for fact check: %w", err)
		}
		seq := gModel.GenerateContent(ctx, req, false)

		for resp, err := range seq {
			if err != nil {
				return rawContext, extractCitationsFromContext(rawContext), fmt.Errorf("error generating fact-check correction: %w", err)
			}
			if resp.Content != nil {
				for _, part := range resp.Content.Parts {
					if part.Text != "" {
						correctedText.WriteString(part.Text)
					}
				}
			}
		}
	}

	res := strings.TrimSpace(correctedText.String())
	if res == "" {
		res = rawContext
	}
	return res, extractCitationsFromContext(res), nil
}

func extractCitationsFromContext(contextText string) []MagisteriumCitation {
	var citations []MagisteriumCitation
	if contextText == "" {
		return citations
	}

	// Extract CCC references like CCC 1718 or Catechism 1718
	cccRegex := regexp.MustCompile(`(?:Catechism(?: of the Catholic Church)?|CCC)\s*(?:§§?|paragraph|#)?\s*(\d+(?:-\d+)?)`)
	matches := cccRegex.FindAllStringSubmatch(contextText, -1)
	seen := make(map[string]bool)
	for _, m := range matches {
		if len(m) > 1 && !seen[m[1]] {
			seen[m[1]] = true
			citations = append(citations, MagisteriumCitation{
				Title:         "Catechism of the Catholic Church",
				DocumentTitle: "Catechism of the Catholic Church",
				Ref:           fmt.Sprintf("CCC %s", m[1]),
				Text:          fmt.Sprintf("Sacred Catholic Doctrine on Catechism of the Catholic Church paragraph %s.", m[1]),
			})
		}
	}
	return citations
}
