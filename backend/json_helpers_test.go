package main

import (
	"testing"
)

func TestExtractDevotionData(t *testing.T) {
	// Sample data provided by user (reconstructed to be valid)
	sampleJSON := `{
  "first_reading": {
    "citation": "James 1:12-18",
    "text": "Blessed is he...",
    "context": "Written by James..."
  },
  "responsorial_psalm": {
    "citation": "Psalm 94:12-13a, 14-15, 18-19",
    "text": "Blessed the man...",
    "context": "As a Catholic Bible scholar..."
  },
  "second_reading": null,
  "gospel": {
    "citation": "Mark 8:14-21",
    "text": "The disciples...",
    "context": "Jesus warns..."
  }
}`

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Clean JSON",
			input:   sampleJSON,
			wantErr: false,
		},
		{
			name:    "Markdown Code Block",
			input:   "Here is the JSON:\n```json\n" + sampleJSON + "\n```",
			wantErr: false,
		},
		{
			name:    "Markdown Code Block with previous text",
			input:   "Sure, here is the devotion data:\n\n```json\n" + sampleJSON + "\n```\nHope this helps!",
			wantErr: false,
		},
		{
			name:    "Multiple JSON objects (should find the DevotionData)",
			input:   `{"other": "stuff"} ` + sampleJSON + ` {"more": "stuff"}`,
			wantErr: false,
		},
		{
			name:    "Broken formatted JSON (fail)",
			input:   "```json\n{ invalid json ...",
			wantErr: true,
		},
		{
			name:    "Brace in String (Robustness Check)",
			input:   `{"first_reading": {"context": "He said: {Behold!} and we rejoiced."}}`,
			wantErr: false,
		},
		{
			name:    "Unbalanced Brace in String (Robustness Check)",
			input:   `{"first_reading": {"context": "He said: {Behold! and } we rejoiced."}}`,
			wantErr: false,
		},
		{
			name:    "Irrelevant JSON (Validation Check)",
			input:   `{"some": "other", "data": "here"}`,
			wantErr: true, // Should fail because it doesn't match DevotionData fields
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractDevotionData(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractDevotionData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("ExtractDevotionData() returned nil success")
			}
		})
	}
}
func TestPopulatePassageTexts(t *testing.T) {
	// Note: This test requires BIBLE_API_KEY and BIBLE_ID to be set in environment
	// if it were to actually fetch data. For unit testing, we'd ideally mock get_bible_passage.

	data := &DevotionData{
		FirstReading: &ReadingWithContext{
			Citation: "James 1:12-18",
		},
		Gospel: &ReadingWithContext{
			Citation: "Mark 8:14-21",
		},
	}

	err := PopulatePassageTexts(data)
	if err != nil {
		t.Errorf("PopulatePassageTexts() error = %v", err)
	}

	// We can't easily assert the content without API keys, but we can check if it ran without crashing
	// and if the struct members are still there.
	if data.FirstReading.Citation != "James 1:12-18" {
		t.Errorf("Citation changed unexpectedly")
	}
}
