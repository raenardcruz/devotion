package main

import "testing"

func TestParseA2UIMessages(t *testing.T) {
	jsonl := `{"surfaceUpdate":{"surfaceId":"response","components":[{"id":"root","component":{"Column":{"children":{"explicitList":["h1","body"]}}}},{"id":"h1","component":{"Text":{"text":{"literalString":"Grace"},"usageHint":"h2"}}},{"id":"body","component":{"Text":{"text":{"literalString":"Sanctifying grace is a gift from God."},"usageHint":"body"}}}]}}
{"beginRendering":{"surfaceId":"response","root":"root"}}`

	msgs, ok := ParseA2UIMessages(jsonl)
	if !ok {
		t.Fatalf("expected valid A2UI JSONL to parse")
	}
	if len(msgs) != 2 {
		t.Fatalf("expected 2 A2UI messages, got %d", len(msgs))
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
