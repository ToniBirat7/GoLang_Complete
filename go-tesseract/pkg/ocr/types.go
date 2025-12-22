package ocr

// ExtractedLine represents a single line of OCR text with its confidence
type ExtractedLine struct {
	Text       string  `json:"text"`
	Confidence float64 `json:"confidence"`
}

// OCRResult represents the complete OCR extraction result
type OCRResult struct {
	Text              string          `json:"text"`
	AverageConfidence float64         `json:"average_confidence"`
	LineCount         int             `json:"line_count"`
	Lines             []ExtractedLine `json:"lines,omitempty"`
}

// OCRConfig holds configuration for OCR processing
type OCRConfig struct {
	Language        string
	IncludeLines    bool
	CleanDevanagari bool
	MinConfidence   float64
}

// DefaultConfig returns the default OCR configuration for Nepali text
func DefaultConfig() *OCRConfig {
	return &OCRConfig{
		Language:        "nep",
		IncludeLines:    false,
		CleanDevanagari: true,
		MinConfidence:   0.0,
	}
}