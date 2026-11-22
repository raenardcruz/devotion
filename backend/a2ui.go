package main

import (
	"encoding/json"
	"strings"
)

// A2UIResponseInstruction returns system/formatting instructions that teach the LLM
// to format its response using the A2UI (Agent-to-User Interface) JSONL protocol.
func A2UIResponseInstruction() string {
	return `Format your ENTIRE response using the A2UI (Agent-to-User Interface) JSONL protocol.
Output ONLY valid JSONL (one JSON object per line, NO markdown code fences, NO extra prose outside JSON).

Protocol rules:
1. Emit a "surfaceUpdate" message defining a flat list of components by ID, then a "beginRendering" message naming the root ID.
2. You may optionally emit "dataModelUpdate" if binding values via "path", or use "literalString" directly.
3. Every component in "components" has {"id": "<unique-id>", "component": {"<Type>": { ... }}}.
4. Available component types:
   - Column: {"children": {"explicitList": ["id1", "id2"]}}
   - Row: {"children": {"explicitList": ["id1", "id2"]}}
   - Card: {"child": "body-id", "title": {"literalString": "Optional Card Title"}, "variant": "default|highlight|quote|callout"}
   - List: {"children": {"explicitList": ["item1", "item2"]}, "ordered": false}
   - Divider: {}
   - Tabs: {"tabItems": [{"title": {"literalString": "Tab 1"}, "child": "tab1-id"}]}
   - Text: {"text": {"literalString": "Markdown-formatted text (bold, italic, inline [1] citations)"}, "usageHint": "h1|h2|h3|body|caption|quote"}
   - Badge: {"label": {"literalString": "CCC 1718"}, "variant": "primary|neutral"}
   - Link: {"label": {"literalString": "Read Source"}, "url": "https://..."}

Example JSONL output:
{"surfaceUpdate":{"surfaceId":"response","components":[{"id":"root","component":{"Column":{"children":{"explicitList":["title","intro","key-card","points"]}}}},{"id":"title","component":{"Text":{"text":{"literalString":"Summary Title"},"usageHint":"h2"}}},{"id":"intro","component":{"Text":{"text":{"literalString":"Clear synthesis paragraph with **bold** emphasis and inline citations [1]."},"usageHint":"body"}}},{"id":"key-card","component":{"Card":{"title":{"literalString":"Key Doctrine"},"variant":"highlight","child":"key-text"}}},{"id":"key-text","component":{"Text":{"text":{"literalString":"\"Quoted magisterial or scriptural text.\""},"usageHint":"quote"}}},{"id":"points","component":{"List":{"ordered":false,"children":{"explicitList":["p1","p2"]}}}},{"id":"p1","component":{"Text":{"text":{"literalString":"**First Point**: Explanation [1]."},"usageHint":"body"}}},{"id":"p2","component":{"Text":{"text":{"literalString":"**Second Point**: Explanation [2]."},"usageHint":"body"}}}]}}
{"beginRendering":{"surfaceId":"response","root":"root"}}`
}

// A2UIBibleContextInstruction returns A2UI formatting instructions tailored for
// Historical & Theological Scripture Context responses.
func A2UIBibleContextInstruction() string {
	return `Format the Bible Context response using the A2UI (Agent-to-User Interface) JSONL protocol.
Output ONLY valid JSONL (one JSON object per line, NO markdown code fences, NO extra prose outside JSON).
Organize the commentary into clean A2UI components:
- A root "Column" containing structured "Card" sections (e.g., Historical & Cultural Setting, Theological Meaning, Catechism & Church Fathers Connections, Practical Meditation), "Text" nodes (with usageHint "h2", "h3", "body", or "quote"), "List" items, and "Badge" chips for Scripture or Catechism references (e.g., "CCC 1718").
- End with {"beginRendering":{"surfaceId":"response","root":"root"}}.`
}

// StripA2UICodeFences removes accidental markdown code fences (```json, ```jsonl, ```)
// around an A2UI JSONL payload.
func StripA2UICodeFences(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if !strings.HasPrefix(trimmed, "```") {
		return trimmed
	}
	lines := strings.Split(trimmed, "\n")
	if len(lines) >= 2 && strings.HasPrefix(strings.TrimSpace(lines[0]), "```") {
		lines = lines[1:]
	}
	if len(lines) >= 1 && strings.TrimSpace(lines[len(lines)-1]) == "```" {
		lines = lines[:len(lines)-1]
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}

// ParseA2UIMessages parses a JSONL string into A2UI message objects.
// Returns the parsed messages and true if at least one valid A2UI message
// (surfaceUpdate, dataModelUpdate, beginRendering, deleteSurface) is found.
func ParseA2UIMessages(raw string) ([]map[string]any, bool) {
	cleaned := StripA2UICodeFences(raw)
	if cleaned == "" {
		return nil, false
	}

	var messages []map[string]any
	hasA2UIKey := false

	for _, line := range strings.Split(cleaned, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || !strings.HasPrefix(line, "{") {
			continue
		}
		var obj map[string]any
		if err := json.Unmarshal([]byte(line), &obj); err != nil {
			continue
		}
		if _, ok := obj["surfaceUpdate"]; ok {
			hasA2UIKey = true
		} else if _, ok := obj["dataModelUpdate"]; ok {
			hasA2UIKey = true
		} else if _, ok := obj["beginRendering"]; ok {
			hasA2UIKey = true
		} else if _, ok := obj["deleteSurface"]; ok {
			hasA2UIKey = true
		}
		messages = append(messages, obj)
	}

	return messages, hasA2UIKey && len(messages) > 0
}

// IsA2UIPayload reports whether raw contains valid A2UI JSONL messages.
func IsA2UIPayload(raw string) bool {
	_, ok := ParseA2UIMessages(raw)
	return ok
}
