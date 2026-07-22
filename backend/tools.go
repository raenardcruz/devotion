package main

import (
	"encoding/json"
	"log"

	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

func get_tools() []tool.Tool {
	massReadingsTool, err := functiontool.New(functiontool.Config{
		Name:        "GetMassReadings",
		Description: "Retrieves the daily Catholic mass readings (First Reading, Psalm, Second Reading, Gospel) from the USCCB website.",
	}, get_mass_readings)
	if err != nil {
		log.Fatalf("Failed to create GetMassReadings tool: %v", err)
	}

	popeQuoteTool, err := functiontool.New(functiontool.Config{
		Name:        "GetPopeQuote",
		Description: "Retrieves 'The words of the Popes' from Vatican News for a given date.",
	}, get_pope_quote)
	if err != nil {
		log.Fatalf("Failed to create GetPopeQuote tool: %v", err)
	}

	bibleContextTool, err := functiontool.New(functiontool.Config{
		Name:        "GetBibleContext",
		Description: "Unpacks the history, language, and theology behind the daily readings to reveal the deeper meaning of the Word.",
	}, get_bible_context_tool)
	if err != nil {
		log.Fatalf("Failed to create GetBibleContext tool: %v", err)
	}

	submitDevotionTool, err := functiontool.New(functiontool.Config{
		Name:        "SubmitDevotion",
		Description: "Call this tool at the very end with the final structured devotion data. This ensures the output is correctly formatted.",
	}, submit_devotion)
	if err != nil {
		log.Fatalf("Failed to create SubmitDevotion tool: %v", err)
	}

	return []tool.Tool{massReadingsTool, popeQuoteTool, bibleContextTool, submitDevotionTool}
}

func submit_devotion(ctx tool.Context, data DevotionData) (map[string]any, error) {
	// We return the data as a map so it can be used directly or stored in state
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var res map[string]any
	if err := json.Unmarshal(b, &res); err != nil {
		return nil, err
	}
	return res, nil
}

type DevotionData struct {
	FirstReading      *ReadingWithContext `json:"first_reading,omitempty"`
	ResponsorialPsalm *ReadingWithContext `json:"responsorial_psalm,omitempty"`
	SecondReading     *ReadingWithContext `json:"second_reading,omitempty"`
	Gospel            *ReadingWithContext `json:"gospel,omitempty"`
	PopeQuote         string              `json:"pope_quote,omitempty"`
}

type ReadingWithContext struct {
	Citation string `json:"citation,omitempty"`
	Context  string `json:"context,omitempty"`
	Text     string `json:"text,omitempty"`
}
