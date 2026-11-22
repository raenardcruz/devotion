package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// A2UICatalogSpec returns the shared A2UI JSONL protocol schema, component catalog,
// and curated sacred/scholarly icon pack reference used across all A2UI prompts.
func A2UICatalogSpec() string {
	return `A2UI (Agent-to-User Interface) JSONL Protocol Specification:
Output ONLY valid JSONL (one compact JSON object per line, NO markdown code fences, NO extra prose before or after JSON).

Protocol rules:
1. Every component in a "surfaceUpdate" has {"id": "<unique-id>", "component": {"<Type>": { ... }}}.
2. Emit "beginRendering" naming the root component ID (surfaceId: "response", root: "root").
3. For progressive rendering, you may emit multiple "surfaceUpdate" lines (e.g., first line defining "root" and header/lead components, second line with {"beginRendering":{"surfaceId":"response","root":"root"}}, followed by additional "surfaceUpdate" lines for each section or card) OR a single "surfaceUpdate" followed by "beginRendering".
4. You may optionally emit "dataModelUpdate" when binding via "path", or use "literalString" directly.
5. Every referenced child ID MUST be defined in a "surfaceUpdate" component. Never leave dangling IDs.

Available Component Catalog:
- Column: {"children": {"explicitList": ["id1", "id2"]}, "gap": "sm|md|lg"}
- Row: {"children": {"explicitList": ["id1", "id2"]}, "gap": "sm|md"}
- Grid: {"columns": 2, "children": {"explicitList": ["c1", "c2"]}} (responsive 2- or 3-column grid for comparisons, dual lists, or 4 senses)
- Card: {"title": {"literalString": "Optional Title"}, "subtitle": {"literalString": "Optional Subtitle"}, "icon": "<icon-name>", "variant": "default|highlight|quote|callout|subtle", "child": "body-id"} (also supports "children": {"explicitList": [...]})
- Quote: {"text": {"literalString": "Quoted Scripture, Catechism, Papal, or Saint text"}, "attribution": {"literalString": "St. Augustine"}, "source": {"literalString": "Confessions, Book I"}, "icon": "quote"}
- Accordion: {"title": {"literalString": "Expandable Detail Title"}, "subtitle": {"literalString": "Optional hint"}, "icon": "<icon-name>", "defaultOpen": false, "child": "body-id"}
- Tabs: {"tabItems": [{"title": {"literalString": "Tab 1"}, "icon": "<icon-name>", "child": "tab1-id"}]}
- KeyValue: {"items": [{"label": {"literalString": "Setting"}, "value": {"literalString": "First-century Galilee"}, "icon": "globe"}]}
- List: {"children": {"explicitList": ["item1", "item2"]}, "ordered": false}
- Text: {"text": {"literalString": "Inline markdown (**bold**, *italic*, inline [1] citations)"}, "usageHint": "h1|h2|h3|lead|body|caption|quote", "icon": "<optional-icon-for-headings>"}
- Badge: {"label": {"literalString": "Short Tag"}, "icon": "<optional-icon>", "variant": "primary|gold|neutral"}
- Icon: {"name": "<icon-name>", "size": "sm|md|lg", "tone": "primary|gold|neutral"}
- Divider: {"variant": "subtle|ornament"}
- Link: {"label": {"literalString": "Read Source"}, "url": "https://..."}

Curated A2UI Icon Pack (use ONLY these exact names in "icon" or Icon "name"):
scripture, cross, magisterium, catechism, scroll, dove, chalice, lamp, compass, crown, bridge, quill, pray, heart, scales, quote, bookmark, checklist, info, globe, shield, hourglass, olive, star`
}

// A2UIResponseInstruction returns response-adaptive A2UI formatting instructions for Magisterium AI chat.
// It instructs the model to first analyze the shape of its answer and select the best visual components.
func A2UIResponseInstruction() string {
	return `Before emitting A2UI, ANALYZE the structure of your answer and choose the layout that makes it clearest and most pleasing to read:
- Short direct answer: "Text" (h2) + concise "Text" (lead/body) + optional "Row" of "Badge" source tags + a "Quote" if citing a key magisterial passage.
- Comparative or multi-category answer (e.g., Corporal vs. Spiritual Works of Mercy, Mortal vs. Venial Sin, Faith & Reason): use "Tabs" or a 2-column "Grid" of "Card"s with distinct icons so the reader can compare categories effortlessly.
- Doctrinal synthesis or multi-point teaching: open with an unboxed "Text" (h2) and "Text" (lead) synthesis (avoid wrapping every paragraph in heavy cards), followed by a structured "List" or focused "Card"s, a reverent "Quote" with attribution for official texts, and a "Card" (variant: "callout", icon: "pray" or "bookmark") for practical takeaway.
- Deep theological or historical topic: include a compact "KeyValue" overview or an "Accordion" (defaultOpen: false) for secondary nuances, historical background, or patristic details so the main answer stays scannable.
- Use icons purposefully from the curated Icon Pack to give visual rhythm, and preserve inline citation markers like [1], [2] inside Text strings (they become interactive citation chips).

` + A2UICatalogSpec() + `

Example progressive JSONL output:
{"surfaceUpdate":{"surfaceId":"response","components":[{"id":"root","component":{"Column":{"children":{"explicitList":["title","tags","lead","quote-1","points-card","takeaway"]}}}},{"id":"title","component":{"Text":{"text":{"literalString":"Sanctifying Grace in Catholic Teaching"},"usageHint":"h2","icon":"dove"}}},{"id":"tags","component":{"Row":{"children":{"explicitList":["b1","b2"]}}}},{"id":"b1","component":{"Badge":{"label":{"literalString":"Catechism"},"icon":"catechism","variant":"primary"}}},{"id":"b2","component":{"Badge":{"label":{"literalString":"Grace & Justification"},"icon":"star","variant":"gold"}}},{"id":"lead","component":{"Text":{"text":{"literalString":"Sanctifying grace is a **habitual, supernatural gift** infused into the soul that makes us share in the divine life [1]."},"usageHint":"lead"}}}]}}
{"beginRendering":{"surfaceId":"response","root":"root"}}
{"surfaceUpdate":{"surfaceId":"response","components":[{"id":"quote-1","component":{"Quote":{"text":{"literalString":"Grace is favor, the free and undeserved help that God gives us to respond to his call."},"attribution":{"literalString":"Catechism of the Catholic Church"},"source":{"literalString":"Grace and Justification [1]"},"icon":"quote"}}},{"id":"points-card","component":{"Card":{"title":{"literalString":"Key Doctrinal Dimensions"},"icon":"magisterium","variant":"subtle","child":"points-list"}}},{"id":"points-list","component":{"List":{"ordered":false,"children":{"explicitList":["p1","p2"]}}}},{"id":"p1","component":{"Text":{"text":{"literalString":"**Participation in Divine Life**: It heals human nature and elevates us to adoptive sonship [1]."},"usageHint":"body"}}},{"id":"p2","component":{"Text":{"text":{"literalString":"**Sacramental Reception**: Ordinarily received in Baptism and restored through Penance [2]."},"usageHint":"body"}}},{"id":"takeaway","component":{"Card":{"title":{"literalString":"Spiritual Reflection"},"icon":"pray","variant":"callout","child":"takeaway-text"}}},{"id":"takeaway-text","component":{"Text":{"text":{"literalString":"Cultivate a state of grace through prayer, the sacraments, and charity toward neighbor [2]."},"usageHint":"body"}}}]}}`
}

// A2UIBibleContextInstruction returns A2UI formatting instructions tailored for
// Historical & Theological Scripture Context responses, aligned with Catholic hermeneutics
// and anti-hallucination citation rules.
func A2UIBibleContextInstruction() string {
	return `FORMAT OVERRIDE — A2UI JSONL OUTPUT:
Present the Scripture commentary using the A2UI (Agent-to-User Interface) JSONL protocol instead of raw multi-line Markdown headings.
Keep all doctrinal, historical, ethical, sentence-limit, and CRITICAL CITATION RULES intact (do NOT invent or guess exact CCC paragraph numbers—use general doctrinal names in badges and text unless explicitly verified).

Analyze the passage commentary and select the cleanest layout:
- Header: "Text" (usageHint: "h2", icon: "scripture") with a concise thematic title + a "Row" of 2–4 thematic "Badge" chips (e.g., doctrinal theme, typology motif, or Church Father—never unverified CCC numbers).
- Historical & Literal Setting: a crisp "Card" (variant: "callout", icon: "scroll") or compact "KeyValue" strip + concise prose setting the scene.
- Spiritual Senses & Typology: when covering multiple senses (Literal, Allegorical, Moral, Anagogical) or OT↔NT typology, organize them cleanly using "Tabs" or a 2-column "Grid" of subtle "Card"s with fitting icons ("scroll" for literal/historical, "lamp" for allegorical, "compass" for moral, "crown" for anagogical, "bridge" for typology).
- Church Tradition & Doctrine: use a "Quote" (icon: "quill") if citing a verified Father/Doctor/Papal teaching, or a "Card" (variant: "subtle", icon: "magisterium") connecting to Catholic doctrine and human dignity (Magnifica Humanitas).
- Practical Meditation: close with a warm "Card" (variant: "highlight", icon: "pray") pointing the reader toward prayer, Sacred Tradition, and parish life.
- Keep the layout airy and scannable—do not duplicate the same sentence across multiple cards.

` + A2UICatalogSpec() + `

Example JSONL output for Bible Context:
{"surfaceUpdate":{"surfaceId":"response","components":[{"id":"root","component":{"Column":{"children":{"explicitList":["header","badges","setting-card","senses-grid","tradition-card","meditation-card"]}}}},{"id":"header","component":{"Text":{"text":{"literalString":"The Good Shepherd & Covenant Fidelity"},"usageHint":"h2","icon":"scripture"}}},{"id":"badges","component":{"Row":{"children":{"explicitList":["b1","b2","b3"]}}}},{"id":"b1","component":{"Badge":{"label":{"literalString":"Old Testament Typology"},"icon":"bridge","variant":"primary"}}},{"id":"b2","component":{"Badge":{"label":{"literalString":"Catechism on Christ the Shepherd"},"icon":"catechism","variant":"gold"}}},{"id":"b3","component":{"Badge":{"label":{"literalString":"St. Augustine"},"icon":"quill","variant":"neutral"}}},{"id":"setting-card","component":{"Card":{"title":{"literalString":"Historical & Literal Setting"},"icon":"scroll","variant":"callout","child":"setting-text"}}},{"id":"setting-text","component":{"Text":{"text":{"literalString":"Spoken in Jerusalem during the Feast of Dedication, Jesus contrasts hired hands with the true shepherd who lays down his life for the flock."},"usageHint":"body"}}},{"id":"senses-grid","component":{"Grid":{"columns":2,"children":{"explicitList":["typology-card","moral-card"]}}}},{"id":"typology-card","component":{"Card":{"title":{"literalString":"Typology & Allegorical Sense"},"icon":"bridge","variant":"subtle","child":"typology-text"}}},{"id":"typology-text","component":{"Text":{"text":{"literalString":"Fulfills the prophetic promise in Ezekiel that God himself would shepherd Israel, revealing Christ's self-giving love for the Church."},"usageHint":"body"}}},{"id":"moral-card","component":{"Card":{"title":{"literalString":"Moral & Anagogical Sense"},"icon":"compass","variant":"subtle","child":"moral-text"}}},{"id":"moral-text","component":{"Text":{"text":{"literalString":"Calls believers to uphold human dignity through sacrificial care while looking forward to the eternal unity of one flock."},"usageHint":"body"}}},{"id":"tradition-card","component":{"Card":{"title":{"literalString":"Tradition & Church Fathers"},"icon":"magisterium","variant":"default","child":"tradition-text"}}},{"id":"tradition-text","component":{"Text":{"text":{"literalString":"St. Augustine reflects that the sheep recognize the Shepherd's voice through humility and fidelity to truth, resonating with Catholic Social Teaching on solidarity."},"usageHint":"body"}}},{"id":"meditation-card","component":{"Card":{"title":{"literalString":"Practical Meditation"},"icon":"pray","variant":"highlight","child":"meditation-text"}}},{"id":"meditation-text","component":{"Text":{"text":{"literalString":"Pause in prayer today to listen for Christ's voice in Scripture and bring that mercy into your local parish community."},"usageHint":"body"}}}]}}
{"beginRendering":{"surfaceId":"response","root":"root"}}`
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

func isA2UIMessageObject(obj map[string]any) bool {
	if _, ok := obj["surfaceUpdate"]; ok {
		return true
	}
	if _, ok := obj["dataModelUpdate"]; ok {
		return true
	}
	if _, ok := obj["beginRendering"]; ok {
		return true
	}
	if _, ok := obj["deleteSurface"]; ok {
		return true
	}
	return false
}

// extractTopLevelJSONObjects scans cleaned text for balanced top-level JSON objects,
// supporting both single-line JSONL and pretty-printed multi-line JSON objects.
func extractTopLevelJSONObjects(cleaned string) []string {
	var objs []string
	depth := 0
	inString := false
	escaped := false
	start := -1

	for i := 0; i < len(cleaned); i++ {
		ch := cleaned[i]
		if inString {
			if escaped {
				escaped = false
				continue
			}
			if ch == '\\' {
				escaped = true
				continue
			}
			if ch == '"' {
				inString = false
			}
			continue
		}

		switch ch {
		case '"':
			inString = true
		case '{':
			if depth == 0 {
				start = i
			}
			depth++
		case '}':
			if depth > 0 {
				depth--
				if depth == 0 && start >= 0 {
					objs = append(objs, cleaned[start:i+1])
					start = -1
				}
			}
		}
	}
	return objs
}

// ParseA2UIMessages parses a JSONL string (or multi-line / array JSON) into A2UI message objects.
// Returns the parsed messages and true if at least one valid A2UI message
// (surfaceUpdate, dataModelUpdate, beginRendering, deleteSurface) is found.
func ParseA2UIMessages(raw string) ([]map[string]any, bool) {
	cleaned := StripA2UICodeFences(raw)
	if cleaned == "" {
		return nil, false
	}

	var messages []map[string]any
	hasA2UIKey := false

	// Fast path: JSON array of messages
	if strings.HasPrefix(cleaned, "[") {
		var arr []map[string]any
		if err := json.Unmarshal([]byte(cleaned), &arr); err == nil {
			for _, obj := range arr {
				if isA2UIMessageObject(obj) {
					hasA2UIKey = true
					messages = append(messages, obj)
				}
			}
			if hasA2UIKey && len(messages) > 0 {
				return messages, true
			}
		}
	}

	// Balanced object scanner handles both compact JSONL and pretty-printed multi-line JSON
	for _, rawObj := range extractTopLevelJSONObjects(cleaned) {
		var obj map[string]any
		if err := json.Unmarshal([]byte(rawObj), &obj); err != nil {
			continue
		}
		if isA2UIMessageObject(obj) {
			hasA2UIKey = true
			messages = append(messages, obj)
		}
	}

	return messages, hasA2UIKey && len(messages) > 0
}

// IsA2UIPayload reports whether raw contains valid A2UI JSONL messages.
func IsA2UIPayload(raw string) bool {
	_, ok := ParseA2UIMessages(raw)
	return ok
}

// NormalizeA2UIPayload strips code fences and re-emits valid A2UI messages as clean
// single-line JSONL. If raw is not valid A2UI, it returns StripA2UICodeFences(raw).
func NormalizeA2UIPayload(raw string) string {
	msgs, ok := ParseA2UIMessages(raw)
	if !ok {
		return StripA2UICodeFences(raw)
	}
	lines := make([]string, 0, len(msgs))
	for _, m := range msgs {
		if b, err := json.Marshal(m); err == nil {
			lines = append(lines, string(b))
		}
	}
	if len(lines) == 0 {
		return StripA2UICodeFences(raw)
	}
	return strings.Join(lines, "\n")
}

func resolveA2UIStringValue(val any, dataModel map[string]any) string {
	switch v := val.(type) {
	case string:
		return v
	case map[string]any:
		if lit, ok := v["literalString"].(string); ok {
			return lit
		}
		if p, ok := v["path"].(string); ok {
			parts := strings.Split(strings.TrimLeft(p, "/"), "/")
			var cur any = dataModel
			for _, part := range parts {
				if part == "" {
					continue
				}
				m, ok := cur.(map[string]any)
				if !ok {
					cur = nil
					break
				}
				cur = m[part]
			}
			if cur != nil {
				return fmt.Sprint(cur)
			}
		}
	}
	return ""
}

func applyA2UIDataEntries(target map[string]any, entries []any) {
	for _, rawEntry := range entries {
		entry, ok := rawEntry.(map[string]any)
		if !ok {
			continue
		}
		key, _ := entry["key"].(string)
		if key == "" {
			continue
		}
		if v, ok := entry["valueString"]; ok {
			target[key] = v
		} else if v, ok := entry["valueNumber"]; ok {
			target[key] = v
		} else if v, ok := entry["valueBoolean"]; ok {
			target[key] = v
		} else if nestedList, ok := entry["valueMap"].([]any); ok {
			nested := make(map[string]any)
			applyA2UIDataEntries(nested, nestedList)
			target[key] = nested
		}
	}
}

func extractChildIDsFromProps(props map[string]any) []string {
	var ids []string
	seen := make(map[string]bool)
	add := func(id string) {
		if id != "" && !seen[id] {
			seen[id] = true
			ids = append(ids, id)
		}
	}

	if child, ok := props["child"].(string); ok {
		add(child)
	}
	switch ch := props["children"].(type) {
	case map[string]any:
		if explicit, ok := ch["explicitList"].([]any); ok {
			for _, item := range explicit {
				if s, ok := item.(string); ok {
					add(s)
				}
			}
		}
	case []any:
		for _, item := range ch {
			if s, ok := item.(string); ok {
				add(s)
			}
		}
	}
	return ids
}

// ExtractA2UIPlainText walks a valid A2UI JSONL payload and returns readable plain text.
// If raw is not an A2UI payload, it returns raw unchanged.
func ExtractA2UIPlainText(raw string) string {
	msgs, ok := ParseA2UIMessages(raw)
	if !ok {
		return raw
	}

	components := make(map[string]map[string]any)
	order := make([]string, 0)
	dataModel := make(map[string]any)
	rootID := ""

	for _, msg := range msgs {
		if su, ok := msg["surfaceUpdate"].(map[string]any); ok {
			if comps, ok := su["components"].([]any); ok {
				for _, rawComp := range comps {
					comp, ok := rawComp.(map[string]any)
					if !ok {
						continue
					}
					id, _ := comp["id"].(string)
					compBody, _ := comp["component"].(map[string]any)
					if id != "" && compBody != nil {
						if _, exists := components[id]; !exists {
							order = append(order, id)
						}
						components[id] = compBody
					}
				}
			}
		} else if dmu, ok := msg["dataModelUpdate"].(map[string]any); ok {
			if contents, ok := dmu["contents"].([]any); ok {
				applyA2UIDataEntries(dataModel, contents)
			}
		} else if br, ok := msg["beginRendering"].(map[string]any); ok {
			if r, ok := br["root"].(string); ok && r != "" {
				rootID = r
			}
		}
	}

	if len(components) == 0 {
		return raw
	}
	if rootID == "" {
		if _, hasRoot := components["root"]; hasRoot {
			rootID = "root"
		} else if len(order) > 0 {
			rootID = order[0]
		}
	}

	visited := make(map[string]bool)
	var parts []string

	var visit func(id string)
	visit = func(id string) {
		if id == "" || visited[id] {
			return
		}
		visited[id] = true
		compBody, exists := components[id]
		if !exists {
			return
		}
		for compType, rawProps := range compBody {
			props, _ := rawProps.(map[string]any)
			if props == nil {
				continue
			}

			if title := resolveA2UIStringValue(props["title"], dataModel); title != "" {
				parts = append(parts, title)
			}
			if subtitle := resolveA2UIStringValue(props["subtitle"], dataModel); subtitle != "" {
				parts = append(parts, subtitle)
			}

			switch compType {
			case "Text":
				if txt := resolveA2UIStringValue(props["text"], dataModel); txt != "" {
					parts = append(parts, txt)
				}
			case "Quote":
				txt := resolveA2UIStringValue(props["text"], dataModel)
				if txt == "" {
					txt = resolveA2UIStringValue(props["quote"], dataModel)
				}
				attr := resolveA2UIStringValue(props["attribution"], dataModel)
				src := resolveA2UIStringValue(props["source"], dataModel)
				if txt != "" {
					line := fmt.Sprintf("\"%s\"", strings.Trim(txt, `"`))
					if attr != "" && src != "" {
						line += fmt.Sprintf(" — %s (%s)", attr, src)
					} else if attr != "" {
						line += fmt.Sprintf(" — %s", attr)
					} else if src != "" {
						line += fmt.Sprintf(" — %s", src)
					}
					parts = append(parts, line)
				}
			case "Badge":
				lbl := resolveA2UIStringValue(props["label"], dataModel)
				if lbl == "" {
					lbl = resolveA2UIStringValue(props["text"], dataModel)
				}
				if lbl != "" {
					parts = append(parts, fmt.Sprintf("[%s]", lbl))
				}
			case "KeyValue":
				if items, ok := props["items"].([]any); ok {
					for _, rawItem := range items {
						item, ok := rawItem.(map[string]any)
						if !ok {
							continue
						}
						lbl := resolveA2UIStringValue(item["label"], dataModel)
						val := resolveA2UIStringValue(item["value"], dataModel)
						if lbl != "" && val != "" {
							parts = append(parts, fmt.Sprintf("%s: %s", lbl, val))
						} else if val != "" {
							parts = append(parts, val)
						}
					}
				}
			case "Tabs":
				if tabs, ok := props["tabItems"].([]any); ok {
					for _, rawTab := range tabs {
						tab, ok := rawTab.(map[string]any)
						if !ok {
							continue
						}
						if tTitle := resolveA2UIStringValue(tab["title"], dataModel); tTitle != "" {
							parts = append(parts, tTitle)
						}
						if childID, ok := tab["child"].(string); ok {
							visit(childID)
						}
					}
				}
			case "Link":
				lbl := resolveA2UIStringValue(props["label"], dataModel)
				u, _ := props["url"].(string)
				if lbl != "" && u != "" {
					parts = append(parts, fmt.Sprintf("%s (%s)", lbl, u))
				} else if u != "" {
					parts = append(parts, u)
				}
			}

			for _, childID := range extractChildIDsFromProps(props) {
				visit(childID)
			}
			break
		}
	}

	visit(rootID)
	for _, id := range order {
		if !visited[id] {
			visit(id)
		}
	}

	res := strings.TrimSpace(strings.Join(parts, "\n\n"))
	if res == "" {
		return raw
	}
	return res
}

// BuildCitationsFallbackA2UI creates a valid A2UI JSONL response directly from retrieved
// Magisterium citations when LLM synthesis is bypassed or fails.
func BuildCitationsFallbackA2UI(citations []MagisteriumCitation) string {
	if len(citations) == 0 {
		return ""
	}

	rootChildren := []string{"header"}
	components := []map[string]any{
		{
			"id": "header",
			"component": map[string]any{
				"Text": map[string]any{
					"text":      map[string]any{"literalString": "Retrieved Magisterial Sources"},
					"usageHint": "h2",
					"icon":      "magisterium",
				},
			},
		},
	}

	for idx, c := range citations {
		cardID := fmt.Sprintf("cite-card-%d", idx+1)
		colID := fmt.Sprintf("cite-col-%d", idx+1)
		rootChildren = append(rootChildren, cardID)

		title := c.DocumentTitle
		if title == "" {
			title = c.Title
		}
		if title == "" {
			title = fmt.Sprintf("Magisterial Source %d", idx+1)
		}

		cardChildIDs := []string{}

		if c.Author != "" || c.Ref != "" {
			metaRowID := fmt.Sprintf("cite-meta-%d", idx+1)
			metaChildren := []string{}
			if c.Author != "" {
				authorID := fmt.Sprintf("cite-author-%d", idx+1)
				metaChildren = append(metaChildren, authorID)
				components = append(components, map[string]any{
					"id": authorID,
					"component": map[string]any{
						"Badge": map[string]any{
							"label":   map[string]any{"literalString": c.Author},
							"icon":    "quill",
							"variant": "gold",
						},
					},
				})
			}
			if c.Ref != "" {
				refID := fmt.Sprintf("cite-ref-%d", idx+1)
				metaChildren = append(metaChildren, refID)
				components = append(components, map[string]any{
					"id": refID,
					"component": map[string]any{
						"Badge": map[string]any{
							"label":   map[string]any{"literalString": c.Ref},
							"icon":    "catechism",
							"variant": "neutral",
						},
					},
				})
			}
			components = append(components, map[string]any{
				"id": metaRowID,
				"component": map[string]any{
					"Row": map[string]any{
						"children": map[string]any{"explicitList": metaChildren},
					},
				},
			})
			cardChildIDs = append(cardChildIDs, metaRowID)
		}

		if strings.TrimSpace(c.Text) != "" {
			quoteID := fmt.Sprintf("cite-quote-%d", idx+1)
			components = append(components, map[string]any{
				"id": quoteID,
				"component": map[string]any{
					"Quote": map[string]any{
						"text": map[string]any{"literalString": strings.TrimSpace(c.Text)},
						"icon": "quote",
					},
				},
			})
			cardChildIDs = append(cardChildIDs, quoteID)
		}

		if strings.TrimSpace(c.URL) != "" {
			linkID := fmt.Sprintf("cite-link-%d", idx+1)
			components = append(components, map[string]any{
				"id": linkID,
				"component": map[string]any{
					"Link": map[string]any{
						"label": map[string]any{"literalString": "Read Original Document"},
						"url":   strings.TrimSpace(c.URL),
					},
				},
			})
			cardChildIDs = append(cardChildIDs, linkID)
		}

		components = append(components, map[string]any{
			"id": colID,
			"component": map[string]any{
				"Column": map[string]any{
					"children": map[string]any{"explicitList": cardChildIDs},
				},
			},
		})

		components = append(components, map[string]any{
			"id": cardID,
			"component": map[string]any{
				"Card": map[string]any{
					"title":   map[string]any{"literalString": fmt.Sprintf("%d. %s", idx+1, title)},
					"icon":    "catechism",
					"variant": "subtle",
					"child":   colID,
				},
			},
		})
	}

	rootComp := map[string]any{
		"id": "root",
		"component": map[string]any{
			"Column": map[string]any{
				"children": map[string]any{"explicitList": rootChildren},
			},
		},
	}

	allComponents := append([]map[string]any{rootComp}, components...)

	surfaceBytes, err := json.Marshal(map[string]any{
		"surfaceUpdate": map[string]any{
			"surfaceId":  "response",
			"components": allComponents,
		},
	})
	if err != nil {
		return ""
	}
	beginBytes, err := json.Marshal(map[string]any{
		"beginRendering": map[string]any{
			"surfaceId": "response",
			"root":      "root",
		},
	})
	if err != nil {
		return ""
	}
	return string(surfaceBytes) + "\n" + string(beginBytes)
}
