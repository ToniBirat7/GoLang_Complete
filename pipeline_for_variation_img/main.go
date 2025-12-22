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
	// Clean non-Devanagari text, keeping numbers and punctuation
	regexGibberish = regexp.MustCompile(`[^\x{0900}-\x{097F}0-9\s.,?!।\-()]+`)
	// Ensure at least one Devanagari character exists
	regexHasDevanagari = regexp.MustCompile(`[\x{0900}-\x{097F}]`)
	// Reduce multiple spaces
	regexMultiSpace = regexp.MustCompile(`\s+`)
)

func main() {
	// --- CONFIGURATION ---
	imagesDir := "../variation_imgs"
	outputRoot := "../variation_outputs"

	// Check if input directory exists
	if _, err := os.Stat(imagesDir); os.IsNotExist(err) {
		fmt.Printf("Error: Input directory '%s' does not exist.\n", imagesDir)
		return
	}

	// 1. Read the Sub-Folders inside variation_imgs
	subFolders, err := os.ReadDir(imagesDir)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Found %d entries in root. Checking for sub-folders...\n", len(subFolders))

	for _, subDir := range subFolders {
		if !subDir.IsDir() {
			continue // Skip loose files in the root folder
		}

		subFolderName := subDir.Name()
		subFolderPath := filepath.Join(imagesDir, subFolderName)

		// Create a matching Output Subfolder
		outputSubFolder := filepath.Join(outputRoot, subFolderName)
		if err := os.MkdirAll(outputSubFolder, 0755); err != nil {
			fmt.Printf("Error creating output dir %s: %v\n", outputSubFolder, err)
			continue
		}

		// Create the Confidence Summary File for this folder
		confSummaryPath := filepath.Join(outputSubFolder, "confidence_summary.csv")
		confFile, err := os.Create(confSummaryPath)
		if err != nil {
			fmt.Printf("Error creating confidence file: %v\n", err)
			continue
		}
		// Write Header
		confFile.WriteString("Image Name,Average Word Confidence\n")

		// Read Images inside this Sub-Folder
		imageFiles, err := os.ReadDir(subFolderPath)
		if err != nil {
			fmt.Printf("Error reading subfolder %s: %v\n", subFolderName, err)
			confFile.Close()
			continue
		}

		fmt.Printf("\n>>> Processing Folder: %s (%d files)\n", subFolderName, len(imageFiles))

		for _, file := range imageFiles {
			if file.IsDir() || !isImageFile(file.Name()) {
				continue
			}

			imagePath := filepath.Join(subFolderPath, file.Name())
			imageBase := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))

			// Output directory specifically for this image
			// Structure: variation_outputs / FolderName / ImageName / ...results...
			imageOutDir := filepath.Join(outputSubFolder, imageBase)

			fmt.Printf("   -> Processing Image: %s\n", file.Name())

			// 1. Run Word Level Extraction (Returns Average Confidence)
			avgConf, err := ExtractWordLevel(imagePath, imageOutDir)
			if err != nil {
				fmt.Printf("      [Error] Word Level: %v\n", err)
			} else {
				// Write to the CSV file
				confLine := fmt.Sprintf("%s,%.4f\n", file.Name(), avgConf)
				confFile.WriteString(confLine)
			}

			// 2. Run Sentence Level
			if err := ExtractSentenceLevel(imagePath, imageOutDir); err != nil {
				fmt.Printf("      [Error] Sentence Level: %v\n", err)
			}

			// 3. Run Paragraph Level
			if err := ExtractParagraphLevel(imagePath, imageOutDir); err != nil {
				fmt.Printf("      [Error] Paragraph Level: %v\n", err)
			}

			// 4. Post-Process: Clean and Reconstruct
			// A. Process Words -> Final Document
			processLevel(imageOutDir, "extraction_results_word_level", "words_data.txt", "final_clean_words.txt", "word")

			// B. Process Sentences -> Final Document
			processLevel(imageOutDir, "extraction_results_sentence_level", "sentences_data.txt", "final_clean_sentences.txt", "sentence")

			// C. Process Paragraphs -> Final Document
			processLevel(imageOutDir, "extraction_results_paragraph_level", "custom_paragraphs_data.txt", "final_clean_paragraphs.txt", "paragraph")
		}

		confFile.Close() // Close the summary file for this folder
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

	file, err := os.Open(inputPath)
	if err != nil {
		return
	}
	defer file.Close()

	var validTokens []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		rawLine := scanner.Text()
		parts := strings.Split(rawLine, "||")
		if len(parts) < 3 {
			continue
		}

		textPart := parts[1]
		cleaned := cleanText(textPart)

		if cleaned == "" {
			continue
		}

		if !regexHasDevanagari.MatchString(cleaned) {
			continue
		}

		validTokens = append(validTokens, cleaned)
	}

	finalContent := ""

	switch levelType {
	case "word":
		var builder strings.Builder
		for _, word := range validTokens {
			builder.WriteString(word)
			if strings.HasSuffix(word, "।") || strings.HasSuffix(word, "?") || strings.HasSuffix(word, "!") {
				builder.WriteString("\n\n")
			} else {
				builder.WriteString(" ")
			}
		}
		finalContent = builder.String()

	case "sentence":
		finalContent = strings.Join(validTokens, "\n")

	case "paragraph":
		finalContent = strings.Join(validTokens, "\n\n")
	}

	os.WriteFile(outputPath, []byte(finalContent), 0644)
}

func cleanText(input string) string {
	s := regexGibberish.ReplaceAllString(input, " ")
	s = regexMultiSpace.ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}
