package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/raenardcruz/ollama-adk-wrapper"
	"google.golang.org/adk/model"
	geminiModel "google.golang.org/adk/model/gemini"
	"google.golang.org/genai"
)

type BibleContext struct {
	Citation string `json:"citation"`
	Context  string `json:"context"`
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
			Citation: args.Citation,
			Context:  rawContext,
		}, nil
	}

	log.Printf("[get_bible_context] Generated draft bible context for %s successfully. Running fact-check & correction...\n", args.Citation)

	factCheckedContext, err := FactCheckBibleContext(ctx, args.Citation, rawContext, passageText)
	if err != nil {
		log.Printf("[get_bible_context] Warning: fact check failed for %s: %v. Returning raw context.", args.Citation, err)
		factCheckedContext = rawContext
	}

	return BibleContext{
		Citation: args.Citation,
		Context:  factCheckedContext,
	}, nil
}

// FactCheckBibleContext fact-checks generated scripture/theological context against scripture text and Catholic doctrine guidelines, correcting any inaccuracies while strictly maintaining original format.
func FactCheckBibleContext(ctx context.Context, citation string, rawContext string, passageText string) (string, error) {
	if strings.TrimSpace(rawContext) == "" {
		return rawContext, nil
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

	if settings.FactCheckerProvider == "ollama" {
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
				return rawContext, fmt.Errorf("error generating fact-check correction with ollama: %w", err)
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
			return rawContext, nil
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
			return rawContext, fmt.Errorf("failed to create gemini model for fact check: %w", err)
		}
		seq := gModel.GenerateContent(ctx, req, false)

		for resp, err := range seq {
			if err != nil {
				return rawContext, fmt.Errorf("error generating fact-check correction: %w", err)
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
		return rawContext, nil
	}
	return res, nil
}
