package ocr

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/otiai10/gosseract/v2"
)

// Regex patterns pre-compiled for performance
var (
	// Clean non-Devanagari text, keeping numbers and punctuation
	regexGibberish = regexp.MustCompile(`[^\x{0900}-\x{097F}0-9\s.,?!।\-()]+`)
	// Ensure at least one Devanagari character exists in a line
	regexHasDevanagari = regexp.MustCompile(`[\x{0900}-\x{097F}]`)
	// Reduce multiple spaces
	regexMultiSpace = regexp.MustCompile(`\s+`)
)

// ExtractFromImage performs OCR on an image file and returns structured results
// This is the main exported function for library users
func ExtractFromImage(imagePath string, config *OCRConfig) (*OCRResult, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// Extract lines from image
	lines, avgConf, err := extractSentenceLevel(imagePath, config.Language)
	if err != nil {
		return nil, fmt.Errorf("failed to extract text: %w", err)
	}

	// Process and clean lines
	var cleanedLines []ExtractedLine
	var validTexts []string

	for _, line := range lines {
		// Skip low confidence lines if threshold is set
		if config.MinConfidence > 0 && line.Confidence < config.MinConfidence {
			continue
		}

		cleaned := line.Text
		if config.CleanDevanagari {
			cleaned = cleanDevanagariText(line.Text)
		}

		// Skip empty or invalid lines
		if cleaned == "" || !regexHasDevanagari.MatchString(cleaned) {
			continue
		}

		cleanedLines = append(cleanedLines, ExtractedLine{
			Text:       cleaned,
			Confidence: line.Confidence,
		})
		validTexts = append(validTexts, cleaned)
	}

	result := &OCRResult{
		Text:              strings.Join(validTexts, "\n"),
		AverageConfidence: avgConf,
		LineCount:         len(cleanedLines),
	}

	if config.IncludeLines {
		result.Lines = cleanedLines
	}

	return result, nil
}

// extractSentenceLevel performs OCR at the line/sentence level
// Returns: Slice of lines, Average Confidence, Error
func extractSentenceLevel(imagePath string, language string) ([]ExtractedLine, float64, error) {
	client := gosseract.NewClient()
	defer client.Close()

	if err := client.SetImage(imagePath); err != nil {
		return nil, 0, fmt.Errorf("failed to set image: %w", err)
	}

	if err := client.SetLanguage(language); err != nil {
		return nil, 0, fmt.Errorf("failed to set language: %w", err)
	}

	// Get Text Lines (Sentence Level)
	boundingBoxes, err := client.GetBoundingBoxes(gosseract.RIL_TEXTLINE)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get bounding boxes: %w", err)
	}

	var results []ExtractedLine
	var totalConfidence float64
	var lineCount int

	for _, box := range boundingBoxes {
		cleanText := strings.TrimSpace(box.Word)
		if cleanText != "" {
			results = append(results, ExtractedLine{
				Text:       cleanText,
				Confidence: box.Confidence,
			})
			totalConfidence += box.Confidence
			lineCount++
		}
	}

	avgConfidence := 0.0
	if lineCount > 0 {
		avgConfidence = totalConfidence / float64(lineCount)
	}

	return results, avgConfidence, nil
}

// cleanDevanagariText removes non-Devanagari gibberish and normalizes spacing
func cleanDevanagariText(text string) string {
	// 1. Remove non-Nepali gibberish (English noise, random symbols)
	cleaned := regexGibberish.ReplaceAllString(text, " ")

	// 2. Fix multiple spaces
	cleaned = regexMultiSpace.ReplaceAllString(cleaned, " ")

	// 3. Trim whitespace
	cleaned = strings.TrimSpace(cleaned)

	return cleaned
}
