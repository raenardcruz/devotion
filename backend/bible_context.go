package main

import (
	"context"
	"fmt"
	"iter"
	"log"
	"os"
	"strings"

	"github.com/raenardcruz/ollama-adk-wrapper"
	"google.golang.org/adk/model"
	geminiModel "google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/genai"
)

type BibleContext struct {
	Citation string `json:"citation"`
	Context  string `json:"context"`
}

type GetBibleContextArgs struct {
	Citation string `json:"citation"`
}

func get_bible_context_tool(ctx tool.Context, args GetBibleContextArgs) (BibleContext, error) {
	return get_bible_context(ctx, args)
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

	provider := strings.ToLower(settings.ContextProvider)

	var seq iter.Seq2[*model.LLMResponse, error]

	if provider == "gemini" {
		apiKey := settings.GeminiAPIKey
		if apiKey == "" {
			return BibleContext{}, fmt.Errorf("gemini API key is not configured in database settings")
		}
		selectedModelName := settings.GeminiModel
		if selectedModelName == "" {
			selectedModelName = "gemini-3.1-flash-lite"
		}

		log.Printf("[get_bible_context] Generating bible context for %s with gemini (%s)...\n", args.Citation, selectedModelName)
		gModel, err := geminiModel.NewModel(ctx, selectedModelName, &genai.ClientConfig{
			APIKey: apiKey,
		})
		if err != nil {
			return BibleContext{}, fmt.Errorf("failed to create gemini model (%s): %w", selectedModelName, err)
		}
		seq = gModel.GenerateContent(ctx, req, false)
	} else {
		// Default to Ollama
		ollamaURL := os.Getenv("OLLAMA_URL")
		if ollamaURL == "" {
			ollamaURL = "http://localhost:11434"
		}
		selectedModelName := settings.OllamaModel
		if selectedModelName == "" {
			selectedModelName = "gemma4:cloud"
		}

		log.Printf("[get_bible_context] Generating bible context for %s with ollama (%s)...\n", args.Citation, selectedModelName)
		oModel := ollama.NewModel(selectedModelName, ollamaURL)
		seq = oModel.GenerateContent(ctx, req, false)
	}

	var contextText strings.Builder
	for resp, err := range seq {
		if err != nil {
			return BibleContext{}, fmt.Errorf("error generating bible context with %s: %w", provider, err)
		}
		if resp.Content != nil {
			for _, part := range resp.Content.Parts {
				if part.Text != "" {
					contextText.WriteString(part.Text)
				}
			}
		}
	}

	log.Printf("[get_bible_context] Generated bible context for %s successfully with %s\n", args.Citation, provider)
	return BibleContext{
		Citation: args.Citation,
		Context:  contextText.String(),
	}, nil
}
