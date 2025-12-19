package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Entry point for the pipeline
func main() {
	imagesDir := "../rajpatra_imgs"     // Change if your images are elsewhere
	outputRoot := "../pipeline_outputs" // All results go here

	// 1. List all image files in the imagesDir
	files, err := ioutil.ReadDir(imagesDir)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		// Accept common image extensions
		if !isImageFile(file.Name()) {
			continue
		}
		imagePath := filepath.Join(imagesDir, file.Name())
		imageBase := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
		imageOutDir := filepath.Join(outputRoot, imageBase)

		// Create output directory for this image
		if err := os.MkdirAll(imageOutDir, 0755); err != nil {
			panic(err)
		}

		// Run all three extraction levels
		fmt.Printf("Processing %s...\n", file.Name())
		runWordLevel(imagePath, imageOutDir)
		runSentenceLevel(imagePath, imageOutDir)
		runParagraphLevel(imagePath, imageOutDir)
	}

	fmt.Println("All images processed.")
}

func isImageFile(name string) bool {
	name = strings.ToLower(name)
	return strings.HasSuffix(name, ".png") || strings.HasSuffix(name, ".jpg") || strings.HasSuffix(name, ".jpeg")
}

// These functions will call the actual extraction logic (to be implemented)
func runWordLevel(imagePath, outDir string) {
	if err := ExtractWordLevel(imagePath, outDir); err != nil {
		fmt.Printf("[ERROR] Word-level extraction failed for %s: %v\n", imagePath, err)
	}
}

func runSentenceLevel(imagePath, outDir string) {
	if err := ExtractSentenceLevel(imagePath, outDir); err != nil {
		fmt.Printf("[ERROR] Sentence-level extraction failed for %s: %v\n", imagePath, err)
	}
}

func runParagraphLevel(imagePath, outDir string) {
	if err := ExtractParagraphLevel(imagePath, outDir); err != nil {
		fmt.Printf("[ERROR] Paragraph-level extraction failed for %s: %v\n", imagePath, err)
	}
}
