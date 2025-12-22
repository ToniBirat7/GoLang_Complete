package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ToniBirat7/tesseract_ocr_ne/pkg/ocr"
)

func main() {
	// Command-line flags
	imagePath := flag.String("image", "", "Path to image file (required)")
	outputPath := flag.String("output", "", "Path to output JSON file (optional, prints to stdout if not specified)")
	includeLines := flag.Bool("lines", false, "Include individual lines in output")
	minConfidence := flag.Float64("min-confidence", 0.0, "Minimum confidence threshold (0-100)")
	flag.Parse()

	// Validate input
	if *imagePath == "" {
		fmt.Fprintln(os.Stderr, "Error: -image flag is required")
		flag.Usage()
		os.Exit(1)
	}

	// Check if image exists
	if _, err := os.Stat(*imagePath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: Image file '%s' does not exist\n", *imagePath)
		os.Exit(1)
	}

	// Configure OCR
	config := ocr.DefaultConfig()
	config.IncludeLines = *includeLines
	config.MinConfidence = *minConfidence

	// Perform OCR
	result, err := ocr.ExtractFromImage(*imagePath, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Marshal to JSON
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}

	// Output result
	if *outputPath != "" {
		// Ensure output directory exists
		outputDir := filepath.Dir(*outputPath)
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
			os.Exit(1)
		}

		// Write to file
		if err := os.WriteFile(*outputPath, jsonData, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Results saved to: %s\n", *outputPath)
	} else {
		// Print to stdout
		fmt.Println(string(jsonData))
	}
}
