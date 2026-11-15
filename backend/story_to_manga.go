package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/raenardcruz/ollama-adk-wrapper"
	"google.golang.org/adk/model"
	"google.golang.org/genai"
)

const DefaultGeminiImageModel = "gemini-nano-banana-2.1"

var SupportedGeminiImageModels = []string{
	"gemini-nano-banana-2.1",
	"gemini-2.5-flash-image-preview",
	"gemini-2.0-flash-preview-image-generation",
	"imagen-3.0-generate-002",
}

// === Domain Types (ported from victorhuangwq/story-to-manga) ===

type MangaCharacter struct {
	Name                string `json:"name"`
	PhysicalDescription string `json:"physicalDescription"`
	Personality         string `json:"personality"`
	Role                string `json:"role"`
}

type MangaSetting struct {
	TimePeriod string `json:"timePeriod"`
	Location   string `json:"location"`
	Mood       string `json:"mood"`
}

type MangaStoryAnalysis struct {
	Title       string           `json:"title"`
	Story       string           `json:"story,omitempty"`
	MoralLesson string           `json:"moralLesson,omitempty"`
	Characters  []MangaCharacter `json:"characters"`
	Setting     MangaSetting     `json:"setting"`
}

type MangaCharacterReference struct {
	Name        string `json:"name"`
	Image       string `json:"image"` // base64 data URL
	Description string `json:"description"`
}

type MangaPanelPlan struct {
	PanelNumber      int      `json:"panelNumber"`
	Characters       []string `json:"characters"`
	SceneDescription string   `json:"sceneDescription"`
	Dialogue         string   `json:"dialogue,omitempty"`
	CameraAngle      string   `json:"cameraAngle"`
	VisualMood       string   `json:"visualMood"`
}

type MangaStoryBreakdown struct {
	Panels []MangaPanelPlan `json:"panels"`
}

type MangaGeneratedPanel struct {
	PanelNumber int    `json:"panelNumber"`
	Image       string `json:"image"` // base64 data URL
}

type MangaUploadedReference struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name"`
	Image    string `json:"image"`
	FileName string `json:"fileName,omitempty"`
}

// === Request / Response Payloads ===

type MangaReadingsInput struct {
	FirstReading      *ReadingWithContext `json:"first_reading,omitempty"`
	ResponsorialPsalm *ReadingWithContext `json:"responsorial_psalm,omitempty"`
	SecondReading     *ReadingWithContext `json:"second_reading,omitempty"`
	Gospel            *ReadingWithContext `json:"gospel,omitempty"`
	PopeQuote         string              `json:"pope_quote,omitempty"`
}

type MangaAnalyzeStoryRequest struct {
	Date            string              `json:"date,omitempty"`
	Readings        *MangaReadingsInput `json:"readings,omitempty"`
	Story           string              `json:"story,omitempty"`
	Style           string              `json:"style"`
	ProtagonistHint string              `json:"protagonistHint,omitempty"`
}

type MangaGenerateCharacterRefsRequest struct {
	Characters                  []MangaCharacter         `json:"characters"`
	Setting                     MangaSetting             `json:"setting"`
	Style                       string                   `json:"style"`
	ImageModel                  string                   `json:"imageModel,omitempty"`
	UploadedCharacterReferences []MangaUploadedReference `json:"uploadedCharacterReferences,omitempty"`
}

type MangaChunkStoryRequest struct {
	Story      string           `json:"story"`
	Characters []MangaCharacter `json:"characters"`
	Setting    MangaSetting     `json:"setting"`
	Style      string           `json:"style"`
	NoDialogue bool             `json:"noDialogue,omitempty"`
}

type MangaGeneratePanelRequest struct {
	Panel                     MangaPanelPlan            `json:"panel"`
	CharacterReferences       []MangaCharacterReference `json:"characterReferences"`
	Setting                   MangaSetting              `json:"setting"`
	Style                     string                    `json:"style"`
	NoDialogue                bool                      `json:"noDialogue,omitempty"`
	ImageModel                string                    `json:"imageModel,omitempty"`
	UploadedSettingReferences []MangaUploadedReference  `json:"uploadedSettingReferences,omitempty"`
}

// === Handlers ===

func mangaConfigHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		settings := GetSettings()
		imageModel := strings.TrimSpace(settings.GeminiImageModel)
		if imageModel == "" {
			imageModel = DefaultGeminiImageModel
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"imageModel": imageModel,
			"models":     SupportedGeminiImageModels,
		})

	case http.MethodPost, http.MethodPut:
		var req struct {
			ImageModel string `json:"imageModel"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		modelName := strings.TrimSpace(req.ImageModel)
		if modelName == "" {
			respondWithError(w, http.StatusBadRequest, "imageModel is required")
			return
		}
		settings := GetSettings()
		settings.GeminiImageModel = modelName
		if err := SaveSettings(settings); err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to save image model setting: %v", err))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":    true,
			"imageModel": settings.GeminiImageModel,
			"models":     SupportedGeminiImageModels,
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func mangaAnalyzeStoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req MangaAnalyzeStoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Style == "" {
		req.Style = "manga"
	}
	if req.Style != "manga" && req.Style != "comic" {
		respondWithError(w, http.StatusBadRequest, "Style must be 'manga' or 'comic'")
		return
	}

	ctx := r.Context()
	storyText := strings.TrimSpace(req.Story)

	// If a custom story is already provided, validate word count (<= 500 words, matching story-to-manga)
	if storyText != "" {
		words := strings.Fields(storyText)
		wordCount := len(words)
		if wordCount > 500 {
			respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Story too long. Maximum 500 words, got %d words.", wordCount))
			return
		}

		prompt := fmt.Sprintf(`Analyze this story and extract the main characters with their detailed characteristics:

Story: "%s"

Style: %s

Please provide:
1. A title for this story (create a catchy, appropriate title if one isn't explicitly mentioned)

2. A list of main characters (1-4 maximum, choose based on story complexity) with:
   - Name
   - Physical description (age, build, hair, clothing, distinctive features)
   - Personality traits
   - Role in the story

3. Setting description (time period, location, mood)
`, storyText, req.Style)

		schema := &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"title":       {Type: genai.TypeString},
				"moralLesson": {Type: genai.TypeString},
				"characters": {
					Type: genai.TypeArray,
					Items: &genai.Schema{
						Type: genai.TypeObject,
						Properties: map[string]*genai.Schema{
							"name":                {Type: genai.TypeString},
							"physicalDescription": {Type: genai.TypeString},
							"personality":         {Type: genai.TypeString},
							"role":                {Type: genai.TypeString},
						},
						PropertyOrdering: []string{"name", "physicalDescription", "personality", "role"},
						Required:         []string{"name", "physicalDescription", "personality", "role"},
					},
				},
				"setting": {
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"timePeriod": {Type: genai.TypeString},
						"location":   {Type: genai.TypeString},
						"mood":       {Type: genai.TypeString},
					},
					PropertyOrdering: []string{"timePeriod", "location", "mood"},
					Required:         []string{"timePeriod", "location", "mood"},
				},
			},
			PropertyOrdering: []string{"title", "moralLesson", "characters", "setting"},
			Required:         []string{"title", "characters", "setting"},
		}

		rawJSON, err := generateStructuredJSON(ctx, prompt, schema)
		if err != nil {
			log.Printf("[mangaAnalyzeStoryHandler] Error analyzing story: %v", err)
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to analyze story: %v", err))
			return
		}

		var analysis MangaStoryAnalysis
		if err := parseMangaJSON(rawJSON, &analysis); err != nil {
			log.Printf("[mangaAnalyzeStoryHandler] Error parsing story analysis JSON: %v (raw: %.300s)", err, rawJSON)
			respondWithError(w, http.StatusInternalServerError, "Failed to parse story analysis")
			return
		}

		analysis.Story = storyText

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":     true,
			"story":       storyText,
			"moralLesson": analysis.MoralLesson,
			"analysis":    analysis,
			"wordCount":   wordCount,
		})
		return
	}

	// Otherwise: 1-Click Custom Story Generation from the day's readings and lessons!
	readingsSummary := buildReadingsSummary(ctx, req.Date, req.Readings)
	if strings.TrimSpace(readingsSummary) == "" {
		respondWithError(w, http.StatusBadRequest, "Readings or valid date are required to generate a story from the day's readings")
		return
	}

	protagonistDirection := strings.TrimSpace(req.ProtagonistHint)
	if protagonistDirection == "" {
		protagonistDirection = "A young boy has a relatable everyday problem or conflict. He then listens to the Mass readings and homily/reflection for that day, takes the lesson from those readings to heart, applies that lesson in action to his situation, and reaches a meaningful resolution."
	}

	prompt := fmt.Sprintf(`You are a master Catholic storyteller and visual comic writer creating an inspiring short story for a %s adaptation based on today's Daily Mass Readings and their spiritual lessons.

Today's Liturgical Readings & Contextual Lessons:
%s

Narrative Arc Requirement:
%s

Follow this exact 4-part story arc in 220 to 350 words (strictly under 450 words):
1. THE PROBLEM: Introduce a relatable protagonist (e.g., a young boy with a specific name) facing a concrete, everyday personal struggle, temptation, misunderstanding, or dilemma at home, school, or in his neighborhood that mirrors the core theme of today's readings.
2. LISTENING TO THE READINGS: The boy attends Mass or family devotion and listens attentively to the readings for that day (explicitly reference the Scripture passage/lesson from today's readings in the story).
3. APPLYING THE LESSON: Illuminated by the lesson of the readings, he makes a courageous, concrete choice to practice that virtue or teaching right where he struggled.
4. RESOLUTION: Show the heartwarming resolution as he applies the lesson, bringing peace, reconciliation, and deeper faith.

Also extract the Story-to-Manga analysis:
- "title": A catchy, memorable title for this story
- "story": The complete 220-350 word narrative text with vivid scenes and natural dialogue
- "moralLesson": A concise 1-2 sentence takeaway summarizing how the day's Scripture reading resolves the problem
- "characters": 1 to 3 main characters from the story with detailed visual descriptions (age, build, hair, clothing, distinctive features so the image generator can draw consistent character reference sheets), personality traits, and role
- "setting": The timePeriod (e.g. "Modern day"), location, and mood`, req.Style, readingsSummary, protagonistDirection)

	schema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"title":       {Type: genai.TypeString},
			"story":       {Type: genai.TypeString},
			"moralLesson": {Type: genai.TypeString},
			"characters": {
				Type: genai.TypeArray,
				Items: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"name":                {Type: genai.TypeString},
						"physicalDescription": {Type: genai.TypeString},
						"personality":         {Type: genai.TypeString},
						"role":                {Type: genai.TypeString},
					},
					PropertyOrdering: []string{"name", "physicalDescription", "personality", "role"},
					Required:         []string{"name", "physicalDescription", "personality", "role"},
				},
			},
			"setting": {
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"timePeriod": {Type: genai.TypeString},
					"location":   {Type: genai.TypeString},
					"mood":       {Type: genai.TypeString},
				},
				PropertyOrdering: []string{"timePeriod", "location", "mood"},
				Required:         []string{"timePeriod", "location", "mood"},
			},
		},
		PropertyOrdering: []string{"title", "story", "moralLesson", "characters", "setting"},
		Required:         []string{"title", "story", "moralLesson", "characters", "setting"},
	}

	rawJSON, err := generateStructuredJSON(ctx, prompt, schema)
	if err != nil {
		log.Printf("[mangaAnalyzeStoryHandler] Error generating story from readings: %v", err)
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to generate story from readings: %v", err))
		return
	}

	var analysis MangaStoryAnalysis
	if err := parseMangaJSON(rawJSON, &analysis); err != nil {
		log.Printf("[mangaAnalyzeStoryHandler] Error parsing generated story JSON: %v (raw: %.300s)", err, rawJSON)
		respondWithError(w, http.StatusInternalServerError, "Failed to parse generated story and character analysis")
		return
	}

	wordCount := len(strings.Fields(analysis.Story))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":     true,
		"story":       analysis.Story,
		"moralLesson": analysis.MoralLesson,
		"analysis":    analysis,
		"wordCount":   wordCount,
	})
}

func mangaGenerateCharacterRefsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req MangaGenerateCharacterRefsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if len(req.Characters) == 0 || req.Style == "" {
		respondWithError(w, http.StatusBadRequest, "Characters, setting, and style are required")
		return
	}

	ctx := r.Context()
	imageModel := resolveGeminiImageModel(req.ImageModel)

	stylePrefix := "American comic book style, colorful superhero art with bold colors and clean line art"
	if req.Style == "manga" {
		stylePrefix = "Japanese manga style, black and white, detailed character design with clean line art and screentones, English text only"
	}

	var characterReferences []MangaCharacterReference

	for _, character := range req.Characters {
		var matchingUploads []MangaUploadedReference
		for _, ref := range req.UploadedCharacterReferences {
			if strings.Contains(strings.ToLower(ref.Name), strings.ToLower(character.Name)) ||
				strings.Contains(strings.ToLower(character.Name), strings.ToLower(ref.Name)) {
				matchingUploads = append(matchingUploads, ref)
			}
		}

		prompt := fmt.Sprintf(`Character reference sheet in %s. 

Full body character design showing front view of %s:
- Physical appearance: %s
- Personality: %s
- Role: %s
- Setting context: %s, %s
`, stylePrefix, character.Name, character.PhysicalDescription, character.Personality, character.Role, req.Setting.TimePeriod, req.Setting.Location)

		if len(matchingUploads) > 0 {
			prompt += fmt.Sprintf("\nIMPORTANT: Use the provided reference images as inspiration for this character's design. The reference images show visual elements that should be incorporated while adapting them to the %s aesthetic. Maintain the essence and key visual features shown in the references.\n", stylePrefix)
		} else if len(req.UploadedCharacterReferences) > 0 {
			prompt += "\nNote: Reference images are provided, but use them as general style inspiration for this character design.\n"
		}

		prompt += "\nThe character should be drawn in a neutral pose against a plain background, showing their full design clearly for reference purposes. This is a character reference sheet that will be used to maintain consistency across multiple comic panels.\n"

		var refImages []string
		for _, upload := range req.UploadedCharacterReferences {
			if upload.Image != "" {
				refImages = append(refImages, upload.Image)
			}
		}

		dataURL, usedModel, err := generateGeminiImage(ctx, imageModel, prompt, refImages)
		if err != nil {
			log.Printf("[mangaGenerateCharacterRefsHandler] Failed to generate character reference for %s: %v", character.Name, err)
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to generate reference for %s: %v", character.Name, err))
			return
		}
		imageModel = usedModel

		characterReferences = append(characterReferences, MangaCharacterReference{
			Name:        character.Name,
			Image:       dataURL,
			Description: character.PhysicalDescription,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":             true,
		"characterReferences": characterReferences,
		"imageModel":          imageModel,
	})
}

func mangaChunkStoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req MangaChunkStoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if strings.TrimSpace(req.Story) == "" || len(req.Characters) == 0 || req.Style == "" {
		respondWithError(w, http.StatusBadRequest, "Story, characters, setting, and style are required")
		return
	}

	var charNames []string
	for _, c := range req.Characters {
		charNames = append(charNames, c.Name)
	}
	characterNames := strings.Join(charNames, ", ")

	layoutGuidance := `American comic panel guidelines:
- Rectangular panels with consistent borders
- Wide establishing shots
- Medium shots for dialogue
- Close-ups for dramatic moments
- Clean, structured compositions
- Bold, clear visual storytelling`
	if req.Style == "manga" {
		layoutGuidance = `Manga panel guidelines:
- Dynamic panel shapes and sizes
- Vertical emphasis for dramatic moments
- Action lines and motion blur for movement
- Close-ups for emotional beats
- Wide shots for establishing scenes
- Dramatic angles and perspectives`
	}

	noDialogueNotice := ""
	actionDetail := "\n- Dialogue (if any)"
	if req.NoDialogue {
		noDialogueNotice = "IMPORTANT: This is NO DIALOGUE MODE - do not include any dialogue, speech bubbles, or spoken words. Focus purely on visual storytelling through actions, expressions, and scene composition."
		actionDetail = " (emphasize visual storytelling without dialogue)"
	}

	prompt := fmt.Sprintf(`Break down this story into individual comic panels with detailed descriptions.

Story: "%s"
Characters: %s
Setting: %s, %s, %s
Style: %s

%s

Create 4-6 panels to tell this story effectively from the opening problem, listening to the day's Scripture readings, applying the lesson, to the final resolution.

%s

For each panel, describe:
- Characters present
- Action/scene description%s
- Camera angle (close-up, medium shot, wide shot, etc.)
- Visual mood/atmosphere

Return as a flat array of panels with sequential panel numbers.`,
		req.Story,
		characterNames,
		req.Setting.Location,
		req.Setting.TimePeriod,
		req.Setting.Mood,
		req.Style,
		layoutGuidance,
		noDialogueNotice,
		actionDetail,
	)

	schema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"panels": {
				Type: genai.TypeArray,
				Items: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"panelNumber":      {Type: genai.TypeInteger},
						"characters":       {Type: genai.TypeArray, Items: &genai.Schema{Type: genai.TypeString}},
						"sceneDescription": {Type: genai.TypeString},
						"dialogue":         {Type: genai.TypeString},
						"cameraAngle":      {Type: genai.TypeString},
						"visualMood":       {Type: genai.TypeString},
					},
					PropertyOrdering: []string{
						"panelNumber",
						"characters",
						"sceneDescription",
						"dialogue",
						"cameraAngle",
						"visualMood",
					},
					Required: []string{"panelNumber", "characters", "sceneDescription", "cameraAngle", "visualMood"},
				},
			},
		},
		PropertyOrdering: []string{"panels"},
		Required:         []string{"panels"},
	}

	ctx := r.Context()
	rawJSON, err := generateStructuredJSON(ctx, prompt, schema)
	if err != nil {
		log.Printf("[mangaChunkStoryHandler] Error chunking story: %v", err)
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to chunk story: %v", err))
		return
	}

	var storyBreakdown MangaStoryBreakdown
	if err := parseMangaJSON(rawJSON, &storyBreakdown); err != nil {
		log.Printf("[mangaChunkStoryHandler] Error parsing story breakdown JSON: %v (raw: %.300s)", err, rawJSON)
		respondWithError(w, http.StatusInternalServerError, "Failed to parse story breakdown")
		return
	}

	for i := range storyBreakdown.Panels {
		if storyBreakdown.Panels[i].PanelNumber <= 0 {
			storyBreakdown.Panels[i].PanelNumber = i + 1
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":        true,
		"storyBreakdown": storyBreakdown,
	})
}

func mangaGeneratePanelHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req MangaGeneratePanelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Panel.PanelNumber <= 0 || req.Style == "" {
		respondWithError(w, http.StatusBadRequest, "Panel, character references, setting, and style are required")
		return
	}

	ctx := r.Context()
	imageModel := resolveGeminiImageModel(req.ImageModel)

	stylePrefix := "American comic book style, full color, clean line art"
	if req.Style == "manga" {
		stylePrefix = "Japanese manga visual style (black and white with screentones), but with English text"
	}

	var charDescriptions []string
	for _, charName := range req.Panel.Characters {
		matched := false
		for _, ref := range req.CharacterReferences {
			if strings.EqualFold(ref.Name, charName) && ref.Image != "" {
				matched = true
				break
			}
		}
		if matched {
			charDescriptions = append(charDescriptions, fmt.Sprintf("%s (matching the character design shown in reference image)", charName))
		} else {
			charDescriptions = append(charDescriptions, charName)
		}
	}
	charactersInPanel := strings.Join(charDescriptions, " and ")
	if charactersInPanel == "" {
		charactersInPanel = "the main character"
	}

	dialoguePart := "No dialogue."
	if req.NoDialogue {
		dialoguePart = "NO DIALOGUE MODE - Focus on pure visual storytelling."
	} else if strings.TrimSpace(req.Panel.Dialogue) != "" {
		dialoguePart = fmt.Sprintf(`Dialogue: "%s"`, req.Panel.Dialogue)
	}

	prompt := fmt.Sprintf(`Create a single comic panel in %s.

Setting: %s, %s, mood: %s

Panel Details:
Panel %d: %s shot of %s. Scene: %s. %s. Mood: %s.

IMPORTANT: Use the character reference images provided to maintain visual consistency. Each character should match their appearance from the reference images exactly.
`,
		stylePrefix,
		req.Setting.Location,
		req.Setting.TimePeriod,
		req.Setting.Mood,
		req.Panel.PanelNumber,
		req.Panel.CameraAngle,
		charactersInPanel,
		req.Panel.SceneDescription,
		dialoguePart,
		req.Panel.VisualMood,
	)

	if len(req.UploadedSettingReferences) > 0 {
		prompt += fmt.Sprintf(`
IMPORTANT: Use the provided setting/environment reference images to guide the visual style, atmosphere, and environmental details of this panel. Incorporate the visual elements, lighting, and mood shown in the setting references while adapting them to the %s aesthetic.
`, stylePrefix)
	}

	if req.NoDialogue {
		prompt += `
The panel should include:
- Clear panel border
- NO speech bubbles or dialogue text - pure visual storytelling only
- Emphasize character expressions, body language, and visual details
- Environmental storytelling without any text or speech
- Sound effects where appropriate
- Consistent character designs matching the references

Generate a single comic panel image with proper framing and composition.
`
	} else {
		prompt += `
The panel should include:
- Clear panel border
- Speech bubbles with dialogue text (if any) - IMPORTANT: If dialogue includes character attribution like "Character: 'text'", only put the spoken text in the speech bubble, NOT the character name
- Thought bubbles if needed
- Sound effects where appropriate
- Consistent character designs matching the references

Generate a single comic panel image with proper framing and composition.
`
	}

	var refImages []string
	for _, charRef := range req.CharacterReferences {
		if charRef.Image != "" {
			refImages = append(refImages, charRef.Image)
		}
	}
	for _, settingRef := range req.UploadedSettingReferences {
		if settingRef.Image != "" {
			refImages = append(refImages, settingRef.Image)
		}
	}

	dataURL, usedModel, err := generateGeminiImage(ctx, imageModel, prompt, refImages)
	if err != nil {
		log.Printf("[mangaGeneratePanelHandler] Failed to generate panel %d: %v", req.Panel.PanelNumber, err)
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to generate panel %d: %v", req.Panel.PanelNumber, err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"generatedPanel": MangaGeneratedPanel{
			PanelNumber: req.Panel.PanelNumber,
			Image:       dataURL,
		},
		"imageModel": usedModel,
	})
}

// === Helpers ===

func resolveGeminiImageModel(requestedModel string) string {
	if m := strings.TrimSpace(requestedModel); m != "" {
		return m
	}
	settings := GetSettings()
	if m := strings.TrimSpace(settings.GeminiImageModel); m != "" {
		return m
	}
	if m := strings.TrimSpace(os.Getenv("GEMINI_IMAGE_MODEL")); m != "" {
		return m
	}
	return DefaultGeminiImageModel
}

func resolveGeminiAPIKey() string {
	settings := GetSettings()
	if k := strings.TrimSpace(settings.GeminiAPIKey); k != "" {
		return k
	}
	if k := strings.TrimSpace(os.Getenv("GEMINI_API_KEY")); k != "" {
		return k
	}
	return strings.TrimSpace(os.Getenv("GOOGLE_AI_API_KEY"))
}

func stripHTMLTags(input string) string {
	detailsRe := regexp.MustCompile(`(?is)<details[\s\S]*?</details>`)
	cleaned := detailsRe.ReplaceAllString(input, "")
	tagRe := regexp.MustCompile(`<[^>]*>`)
	cleaned = tagRe.ReplaceAllString(cleaned, " ")
	return strings.Join(strings.Fields(cleaned), " ")
}

func buildReadingsSummary(ctx context.Context, date string, input *MangaReadingsInput) string {
	var first, psalm, second, gospel *ReadingWithContext
	var popeQuote string

	if input != nil && (input.FirstReading != nil || input.Gospel != nil || input.ResponsorialPsalm != nil || input.SecondReading != nil) {
		first = input.FirstReading
		psalm = input.ResponsorialPsalm
		second = input.SecondReading
		gospel = input.Gospel
		popeQuote = input.PopeQuote
	} else {
		targetDate := strings.TrimSpace(date)
		if targetDate == "" {
			targetDate = time.Now().Format("2006-01-02")
		}
		devotionData, err := fetchDevotionData(ctx, targetDate)
		if err != nil {
			log.Printf("[buildReadingsSummary] Warning: could not fetch devotion data for %s: %v", targetDate, err)
			return ""
		}
		first = devotionData.FirstReading
		psalm = devotionData.ResponsorialPsalm
		second = devotionData.SecondReading
		gospel = devotionData.Gospel
		popeQuote = devotionData.PopeQuote
	}

	var sb strings.Builder
	appendSection := func(label string, r *ReadingWithContext) {
		if r == nil || (r.Citation == "" && r.Text == "" && r.Context == "") {
			return
		}
		sb.WriteString(fmt.Sprintf("### %s: %s\n", label, r.Citation))
		if txt := stripHTMLTags(r.Text); txt != "" {
			if len(txt) > 900 {
				txt = txt[:900] + "..."
			}
			sb.WriteString(fmt.Sprintf("Scripture Text: %s\n", txt))
		}
		if ctxText := stripHTMLTags(r.Context); ctxText != "" {
			if len(ctxText) > 900 {
				ctxText = ctxText[:900] + "..."
			}
			sb.WriteString(fmt.Sprintf("Lesson / Contextual Meditation: %s\n", ctxText))
		}
		sb.WriteString("\n")
	}

	appendSection("First Reading", first)
	appendSection("Responsorial Psalm", psalm)
	appendSection("Second Reading", second)
	appendSection("Gospel", gospel)

	if pq := stripHTMLTags(popeQuote); pq != "" {
		if len(pq) > 600 {
			pq = pq[:600] + "..."
		}
		sb.WriteString(fmt.Sprintf("### Words of the Popes:\n%s\n", pq))
	}

	return strings.TrimSpace(sb.String())
}

func generateStructuredJSON(ctx context.Context, prompt string, schema *genai.Schema) (string, error) {
	settings := GetSettings()
	apiKey := resolveGeminiAPIKey()

	if apiKey != "" {
		textModel := strings.TrimSpace(settings.GeminiModel)
		if textModel == "" {
			textModel = "gemini-2.5-flash"
		}

		client, err := genai.NewClient(ctx, &genai.ClientConfig{
			APIKey: apiKey,
		})
		if err == nil {
			candidateModels := []string{textModel, "gemini-2.5-flash", "gemini-2.0-flash"}
			var lastErr error
			for _, m := range uniqueStrings(candidateModels) {
				cfg := &genai.GenerateContentConfig{
					ResponseMIMEType: "application/json",
					ResponseSchema:   schema,
				}
				resp, genErr := client.Models.GenerateContent(ctx, m, []*genai.Content{
					genai.NewContentFromText(prompt, "user"),
				}, cfg)
				if genErr == nil && resp != nil {
					text := resp.Text()
					if strings.TrimSpace(text) != "" {
						return text, nil
					}
				}
				if genErr != nil {
					lastErr = genErr
					log.Printf("[generateStructuredJSON] Gemini model %s returned error: %v", m, genErr)
				}
			}
			if lastErr != nil && settings.ContextProvider == "gemini" {
				return "", lastErr
			}
		}
	}

	// Fallback to Ollama if Gemini key is not set or failed
	ollamaURL := os.Getenv("OLLAMA_URL")
	if ollamaURL == "" {
		ollamaURL = "http://localhost:11434"
	}
	ollamaModel := strings.TrimSpace(settings.OllamaModel)
	if ollamaModel == "" {
		ollamaModel = "gemma4:cloud"
	}

	jsonPrompt := prompt + "\n\nCRITICAL: Return ONLY valid JSON matching the requested fields with no markdown code fences or extra commentary."
	req := &model.LLMRequest{
		Contents: []*genai.Content{
			genai.NewContentFromText(jsonPrompt, "user"),
		},
	}

	oModel := ollama.NewModel(ollamaModel, ollamaURL)
	seq := oModel.GenerateContent(ctx, req, false)
	var out strings.Builder
	for resp, err := range seq {
		if err != nil {
			return "", fmt.Errorf("error generating structured JSON with ollama (%s): %w", ollamaModel, err)
		}
		if resp.Content != nil {
			for _, part := range resp.Content.Parts {
				if part.Text != "" {
					out.WriteString(part.Text)
				}
			}
		}
	}
	return out.String(), nil
}

func decodeDataURL(dataURL string) ([]byte, string, error) {
	trimmed := strings.TrimSpace(dataURL)
	mimeType := "image/jpeg"
	base64Part := trimmed

	if strings.HasPrefix(trimmed, "data:") {
		commaIdx := strings.Index(trimmed, ",")
		if commaIdx == -1 {
			return nil, "", fmt.Errorf("invalid data URL")
		}
		header := trimmed[5:commaIdx]
		base64Part = trimmed[commaIdx+1:]
		if semiIdx := strings.Index(header, ";"); semiIdx != -1 {
			if header[:semiIdx] != "" {
				mimeType = header[:semiIdx]
			}
		}
	}

	rawBytes, err := base64.StdEncoding.DecodeString(base64Part)
	if err != nil {
		return nil, "", err
	}
	return rawBytes, mimeType, nil
}

func generateGeminiImage(ctx context.Context, requestedModel string, prompt string, referenceDataURLs []string) (string, string, error) {
	apiKey := resolveGeminiAPIKey()
	if apiKey == "" {
		return "", requestedModel, fmt.Errorf("Gemini API key is not configured in Admin Settings")
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return "", requestedModel, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	var parts []*genai.Part
	parts = append(parts, genai.NewPartFromText(prompt))

	for _, refURL := range referenceDataURLs {
		imgBytes, mimeType, decErr := decodeDataURL(refURL)
		if decErr == nil && len(imgBytes) > 0 {
			parts = append(parts, genai.NewPartFromBytes(imgBytes, mimeType))
		}
	}

	contents := []*genai.Content{
		{
			Role:  "user",
			Parts: parts,
		},
	}

	// Always try the user's selected/configured model first (e.g. gemini-nano-banana-2.1),
	// followed by public Gemini Flash Image aliases if the upstream endpoint returns 404/unsupported.
	modelsToTry := uniqueStrings([]string{
		requestedModel,
		"gemini-2.5-flash-image-preview",
		"gemini-2.0-flash-preview-image-generation",
	})

	var lastErr error
	for _, m := range modelsToTry {
		log.Printf("[generateGeminiImage] Calling Gemini image generation with model %q (refs=%d)", m, len(referenceDataURLs))

		var cfg *genai.GenerateContentConfig
		if strings.Contains(m, "2.0-flash") {
			cfg = &genai.GenerateContentConfig{
				ResponseModalities: []string{"TEXT", "IMAGE"},
			}
		}

		resp, genErr := client.Models.GenerateContent(ctx, m, contents, cfg)
		if genErr != nil {
			lastErr = genErr
			log.Printf("[generateGeminiImage] Model %q error: %v", m, genErr)
			continue
		}

		if resp != nil && len(resp.Candidates) > 0 {
			cand := resp.Candidates[0]
			if string(cand.FinishReason) == "PROHIBITED_CONTENT" {
				return "", m, fmt.Errorf("PROHIBITED_CONTENT: content was blocked by Gemini safety filters")
			}
			if cand.Content != nil {
				for _, part := range cand.Content.Parts {
					if part.InlineData != nil && len(part.InlineData.Data) > 0 {
						mimeType := part.InlineData.MIMEType
						if mimeType == "" {
							mimeType = "image/jpeg"
						}
						b64 := base64.StdEncoding.EncodeToString(part.InlineData.Data)
						return fmt.Sprintf("data:%s;base64,%s", mimeType, b64), requestedModel, nil
					}
				}
			}
		}
		lastErr = fmt.Errorf("no image data received in response parts from model %s", m)
	}

	return "", requestedModel, lastErr
}

// parseMangaJSON cleans markdown code blocks (ported from victorhuangwq/story-to-manga/src/lib/json-parser.ts)
func parseMangaJSON(text string, target interface{}) error {
	cleaned := strings.TrimSpace(text)
	fenceRe := regexp.MustCompile("(?s)^```(?:json)?\\s*(.*?)\\s*```$")
	if matches := fenceRe.FindStringSubmatch(cleaned); len(matches) > 1 {
		cleaned = strings.TrimSpace(matches[1])
	}

	if err := json.Unmarshal([]byte(cleaned), target); err == nil {
		return nil
	}

	// Fallback: find outermost JSON object
	start := strings.Index(cleaned, "{")
	end := strings.LastIndex(cleaned, "}")
	if start != -1 && end > start {
		return json.Unmarshal([]byte(cleaned[start:end+1]), target)
	}
	return fmt.Errorf("invalid JSON response")
}

func uniqueStrings(items []string) []string {
	seen := make(map[string]bool)
	var out []string
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" && !seen[trimmed] {
			seen[trimmed] = true
			out = append(out, trimmed)
		}
	}
	return out
}
