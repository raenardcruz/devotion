package main

import (
	"strings"
	"testing"
)

func TestParseA2UIMessages(t *testing.T) {
	jsonl := `{"surfaceUpdate":{"surfaceId":"response","components":[{"id":"root","component":{"Column":{"children":{"explicitList":["h1","body"]}}}},{"id":"h1","component":{"Text":{"text":{"literalString":"Grace"},"usageHint":"h2","icon":"dove"}}},{"id":"body","component":{"Text":{"text":{"literalString":"Sanctifying grace is a gift from God."},"usageHint":"body"}}}]}}
{"beginRendering":{"surfaceId":"response","root":"root"}}`

	msgs, ok := ParseA2UIMessages(jsonl)
	if !ok {
		t.Fatalf("expected valid A2UI JSONL to parse")
	}
	if len(msgs) != 2 {
		t.Fatalf("expected 2 A2UI messages, got %d", len(msgs))
	}
}

func TestParseA2UIMessages_PrettyPrintedMultiLineAndArray(t *testing.T) {
	pretty := `{
  "surfaceUpdate": {
    "surfaceId": "response",
    "components": [
      {
        "id": "root",
        "component": {
          "Text": { "text": { "literalString": "Pretty Printed" } }
        }
      }
    ]
  }
}
{
  "beginRendering": {
    "surfaceId": "response",
    "root": "root"
  }
}`
	msgs, ok := ParseA2UIMessages(pretty)
	if !ok || len(msgs) != 2 {
		t.Fatalf("expected pretty-printed multi-line JSON to parse into 2 messages, got ok=%v len=%d", ok, len(msgs))
	}

	normalized := NormalizeA2UIPayload(pretty)
	if len(strings.Split(normalized, "\n")) != 2 {
		t.Fatalf("expected normalized A2UI to have 2 compact JSONL lines, got:\n%s", normalized)
	}

	arr := `[{"surfaceUpdate":{"surfaceId":"response","components":[{"id":"root","component":{"Text":{"text":{"literalString":"Array Mode"}}}}]}},{"beginRendering":{"surfaceId":"response","root":"root"}}]`
	arrMsgs, arrOk := ParseA2UIMessages(arr)
	if !arrOk || len(arrMsgs) != 2 {
		t.Fatalf("expected JSON array of A2UI messages to parse, got ok=%v len=%d", arrOk, len(arrMsgs))
	}
}

func TestParseA2UIMessages_WithCodeFences(t *testing.T) {
	wrapped := "```jsonl\n" +
		`{"surfaceUpdate":{"surfaceId":"response","components":[{"id":"root","component":{"Text":{"text":{"literalString":"Hello"}}}}]}}` + "\n" +
		`{"beginRendering":{"surfaceId":"response","root":"root"}}` + "\n" +
		"```"

	if !IsA2UIPayload(wrapped) {
		t.Fatalf("expected wrapped A2UI JSONL to be recognized after fence stripping")
	}
}

func TestParseA2UIMessages_PlainMarkdownFallback(t *testing.T) {
	markdown := "## Heading\n\nThis is a normal markdown response."
	if IsA2UIPayload(markdown) {
		t.Fatalf("expected plain markdown not to be treated as A2UI payload")
	}
}

func TestExtractA2UIPlainText(t *testing.T) {
	jsonl := `{"surfaceUpdate":{"surfaceId":"response","components":[{"id":"root","component":{"Column":{"children":{"explicitList":["title","kv","q1"]}}}},{"id":"title","component":{"Text":{"text":{"literalString":"The Good Shepherd"},"usageHint":"h2","icon":"scripture"}}},{"id":"kv","component":{"KeyValue":{"items":[{"label":{"literalString":"Setting"},"value":{"literalString":"Jerusalem"}}]}}}, {"id":"q1","component":{"Quote":{"text":{"literalString":"I am the good shepherd."},"attribution":{"literalString":"John 10:11"}}}}]}}
{"beginRendering":{"surfaceId":"response","root":"root"}}`

	plain := ExtractA2UIPlainText(jsonl)
	for _, want := range []string{"The Good Shepherd", "Setting: Jerusalem", "I am the good shepherd.", "John 10:11"} {
		if !strings.Contains(plain, want) {
			t.Fatalf("expected plain text to contain %q, got:\n%s", want, plain)
		}
	}
	if strings.Contains(plain, "surfaceUpdate") {
		t.Fatalf("expected JSON boilerplate to be stripped from plain text, got:\n%s", plain)
	}
}

func TestBuildCitationsFallbackA2UI(t *testing.T) {
	citations := []MagisteriumCitation{
		{
			DocumentTitle: "Catechism of the Catholic Church",
			Author:        "Catholic Church",
			Ref:           "CCC 1996",
			Text:          "Our justification comes from the grace of God.",
			URL:           "https://www.vatican.va/archive/ENG0015/_P6Z.HTM",
		},
	}

	payload := BuildCitationsFallbackA2UI(citations)
	if !IsA2UIPayload(payload) {
		t.Fatalf("expected fallback citations payload to be valid A2UI JSONL, got:\n%s", payload)
	}
	plain := ExtractA2UIPlainText(payload)
	if !strings.Contains(plain, "Catechism of the Catholic Church") || !strings.Contains(plain, "Our justification comes from the grace of God.") {
		t.Fatalf("expected fallback plain text to include citation content, got:\n%s", plain)
	}
}

func TestA2UIBibleContextInstructionIncludesSchemaAndIcons(t *testing.T) {
	instr := A2UIBibleContextInstruction()
	for _, token := range []string{"surfaceUpdate", "beginRendering", "explicitList", "literalString", "Curated A2UI Icon Pack", "scripture", "bridge", "pray"} {
		if !strings.Contains(instr, token) {
			t.Fatalf("expected A2UIBibleContextInstruction to contain %q", token)
		}
	}
}
