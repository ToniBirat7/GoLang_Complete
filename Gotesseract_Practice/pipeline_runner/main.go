package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Regex patterns pre-compiled for performance
var (
	// FIXED REGEX: Removed unnecessary backslashes.
	// We keep \s (whitespace) and \- (hyphen).
	// The rest (.,?!।()) are treated as literals inside [].
	regexGibberish = regexp.MustCompile(`[^\x{0900}-\x{097F}0-9\s.,?!।\-()]+`)

	// Helper to check if a string contains at least one Devanagari character.
	regexHasDevanagari = regexp.MustCompile(`[\x{0900}-\x{097F}]`)

	// Reduce multiple spaces to one
	regexMultiSpace = regexp.MustCompile(`\s+`)
)

func main() {
	// --- CONFIGURATION ---
	imagesDir := "../rajpatra_imgs"
	outputRoot := "../rajpatra_outputs"

	// Check if input directory exists
	if _, err := os.Stat(imagesDir); os.IsNotExist(err) {
		fmt.Printf("Error: Input directory '%s' does not exist.\n", imagesDir)
		return
	}

	files, err := os.ReadDir(imagesDir)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Found %d files. Starting processing...\n", len(files))

	for _, file := range files {
		if file.IsDir() || !isImageFile(file.Name()) {
			continue
		}

		imagePath := filepath.Join(imagesDir, file.Name())
		imageBase := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
		imageOutDir := filepath.Join(outputRoot, imageBase)

		fmt.Printf(">>> Processing: %s\n", file.Name())

		// 1. Run Extractions (Calls functions in extractors.go)
		// We ignore errors here to allow continuing to next stages/images
		if err := ExtractWordLevel(imagePath, imageOutDir); err != nil {
			fmt.Printf("   [Error] Word Level: %v\n", err)
		}
		if err := ExtractSentenceLevel(imagePath, imageOutDir); err != nil {
			fmt.Printf("   [Error] Sentence Level: %v\n", err)
		}
		if err := ExtractParagraphLevel(imagePath, imageOutDir); err != nil {
			fmt.Printf("   [Error] Paragraph Level: %v\n", err)
		}

		// 2. Post-Process: Clean and Reconstruct
		fmt.Println("   ... Cleaning and Formatting Corpus data ...")

		// A. Process Words -> Final Document (Joined with spaces, empty lines on sentence end)
		processLevel(imageOutDir, "extraction_results_word_level", "words_data.txt", "final_clean_words.txt", "word")

		// B. Process Sentences -> Final Document (One sentence per line)
		processLevel(imageOutDir, "extraction_results_sentence_level", "sentences_data.txt", "final_clean_sentences.txt", "sentence")

		// C. Process Paragraphs -> Final Document (Double Newline separated)
		processLevel(imageOutDir, "extraction_results_paragraph_level", "custom_paragraphs_data.txt", "final_clean_paragraphs.txt", "paragraph")
	}

	fmt.Println("\nAll Done! Check the output folder.")
}

func isImageFile(name string) bool {
	name = strings.ToLower(name)
	return strings.HasSuffix(name, ".png") || strings.HasSuffix(name, ".jpg") || strings.HasSuffix(name, ".jpeg")
}

// processLevel reads the OCR data file, cleans it, and writes the NLP-ready file
func processLevel(baseDir, subFolder, inputFile, outputFile, levelType string) {
	inputPath := filepath.Join(baseDir, subFolder, inputFile)
	outputPath := filepath.Join(baseDir, subFolder, outputFile)

	// Open raw OCR data file
	file, err := os.Open(inputPath)
	if err != nil {
		// File might not exist if extraction failed or image was empty
		return
	}
	defer file.Close()

	var validTokens []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		rawLine := scanner.Text()

		// Parse our format: "[Index] || Text Content || Confidence"
		parts := strings.Split(rawLine, "||")
		if len(parts) < 3 {
			continue
		}

		textPart := parts[1] // The middle part is the text

		// --- STEP 1: Cleaning ---
		cleaned := cleanText(textPart)

		// --- STEP 2: Filter Noise ---
		if cleaned == "" {
			continue
		}

		// If the line has NO Devanagari characters (only numbers or punctuation),
		// we skip it to ensure the corpus is high quality language data.
		if !regexHasDevanagari.MatchString(cleaned) {
			continue
		}

		validTokens = append(validTokens, cleaned)
	}

	// --- STEP 3: Reconstruction ---
	finalContent := ""

	switch levelType {
	case "word":
		// Join words. If a word ends with punctuation, add an EMPTY LINE (double newline)
		var builder strings.Builder
		for _, word := range validTokens {
			builder.WriteString(word)

			// Check for sentence endings (Purna Viram, Question, Exclamation)
			if strings.HasSuffix(word, "।") || strings.HasSuffix(word, "?") || strings.HasSuffix(word, "!") {
				builder.WriteString("\n\n") // Empty line separator
			} else {
				builder.WriteString(" ")
			}
		}
		finalContent = builder.String()

	case "sentence":
		// One sentence per line
		finalContent = strings.Join(validTokens, "\n")

	case "paragraph":
		// Paragraphs separated by double newlines to preserve structure
		finalContent = strings.Join(validTokens, "\n\n")
	}

	// Write Final File
	err = os.WriteFile(outputPath, []byte(finalContent), 0644)
	if err != nil {
		fmt.Printf("Failed to write %s: %v\n", outputFile, err)
	}
}

// cleanText applies regex to remove non-Devanagari noise
func cleanText(input string) string {
	// 1. Remove non-allowed characters (English letters, weird symbols)
	// We replace them with space to avoid merging words accidentally
	s := regexGibberish.ReplaceAllString(input, " ")

	// 2. Normalize spaces (turn "   " into " ")
	s = regexMultiSpace.ReplaceAllString(s, " ")

	// 3. Trim edges
	return strings.TrimSpace(s)
}
