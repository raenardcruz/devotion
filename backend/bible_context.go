package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"google.golang.org/adk/model"
	"google.golang.org/adk/tool"
	"google.golang.org/genai"

	"github.com/raenardcruz/ollama-adk-wrapper"
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
		// We continue even if passage text fetching fails, using just the citation
	}

	prompt := fmt.Sprintf(`As a Catholic Bible scholar, provide a detailed and insightful context for the passage %s. Focus on the historical setting, literal, moral, allegorical and anagogical senses of the scripture. Relate this passage to the Catechism of the Catholic Church (CCC) or Catholic doctrine if possible. Additionally, explore the typology (e.g., how an Old Testament passage prefigures New Testament realities, or how a New Testament passage fulfills Old Testament foreshadowing). Where relevant, mention early Church writers or Doctors of the Church (such as Saint Justin Martyr, Saint Ignatius, Saint Irenaeus, or Saint Augustine) and other saints. Use clear, accessible English. Structure your response with clean formatting and typography (using Markdown bold, italics, or lists where helpful to make the details highly readable). Limit your response to a maximum of 15 sentences.

CRITICAL CITATION RULES:

DEFAULT TO GENERAL DOCTRINAL REFERENCES: Because exact paragraph numbers are highly prone to hallucination, do not attempt to guess specific CCC paragraph numbers. Instead, name the specific Catholic doctrine, section, or theme being referenced (e.g., "The Catechism's teaching on Grace," or "The CCC section on the Sacrament of Baptism"). Only include an exact paragraph number if you are executing a verified internal data match; otherwise, omit the number entirely and use the doctrinal name.

VERIFY TEXTUAL ALIGNMENT: Before outputting a specific CCC paragraph number or Bible verse, verify that the text you are writing explicitly matches the actual content of that number/verse. If there is a mismatch between the description and the number, delete the number and keep only the general text description.

STRICT BIBLE CITATION: Any verse cited must exactly match the text and citation of that verse. If unsure of the exact verse number, refer to the chapter generally.`, args.Citation)
	if passageText != "" {
		prompt = fmt.Sprintf(`As a Catholic Bible scholar, provide a detailed and insightful context for the passage %s.

Passage Text:
%s

Focus on the historical setting, literal, moral, allegorical and anagogical senses of the scripture. Relate this passage to the Catechism of the Catholic Church (CCC) or Catholic doctrine if possible. Additionally, explore the typology (e.g., how an Old Testament passage prefigures New Testament realities, or how a New Testament passage fulfills Old Testament foreshadowing). Where relevant, mention early Church writers or Doctors of the Church (such as Saint Justin Martyr, Saint Ignatius, Saint Irenaeus, or Saint Augustine) and other saints. Use clear, accessible English. Structure your response with clean formatting and typography (using Markdown bold, italics, or lists where helpful to make the details highly readable). Limit your response to a maximum of 15 sentences.

CRITICAL CITATION RULES:

DEFAULT TO GENERAL DOCTRINAL REFERENCES: Because exact paragraph numbers are highly prone to hallucination, do not attempt to guess specific CCC paragraph numbers. Instead, name the specific Catholic doctrine, section, or theme being referenced (e.g., "The Catechism's teaching on Grace," or "The CCC section on the Sacrament of Baptism"). Only include an exact paragraph number if you are executing a verified internal data match; otherwise, omit the number entirely and use the doctrinal name.

VERIFY TEXTUAL ALIGNMENT: Before outputting a specific CCC paragraph number or Bible verse, verify that the text you are writing explicitly matches the actual content of that number/verse. If there is a mismatch between the description and the number, delete the number and keep only the general text description.

STRICT BIBLE CITATION: Any verse cited must exactly match the text and citation of that verse. If unsure of the exact verse number, refer to the chapter generally.`, args.Citation, passageText)
	}

	ollamaURL := os.Getenv("OLLAMA_URL")
	if ollamaURL == "" {
		ollamaURL = "http://localhost:11434"
	}
	m := ollama.NewModel("gemma4:cloud", ollamaURL)

	req := &model.LLMRequest{
		Contents: []*genai.Content{
			genai.NewContentFromText(prompt, "user"),
		},
	}
	log.Printf("[get_bible_context] Generating bible context for %s with ollama...\n", args.Citation)
	var contextText strings.Builder
	for resp, err := range m.GenerateContent(ctx, req, false) {
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

	log.Printf("[get_bible_context] Generated bible context for %s with ollama...\n%s", args.Citation, contextText.String())
	return BibleContext{
		Citation: args.Citation,
		Context:  contextText.String(),
	}, nil
}
