package main

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
