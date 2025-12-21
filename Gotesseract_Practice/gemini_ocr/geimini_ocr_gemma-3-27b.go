package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

func main() {
	// --- CONFIGURATION ---
	imagePath := "rajpatra_imgs/img-1.png"
	modelName := "gemma-3-27b-it"

	// --- SETUP ---
	err := godotenv.Load()
	if err != nil {
		log.Println("Note: .env file not found, checking system environment variables...")
	}

	apiKey := os.Getenv("GEMMA_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("GEMINI_API_KEY")
	}
	if apiKey == "" {
		log.Fatal("Error: Neither GEMMA_API_KEY nor GEMINI_API_KEY is set")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	defer client.Close()

	model := client.GenerativeModel(modelName)
	model.SetTemperature(0)

	// --- LOAD IMAGE ---
	imgData, err := os.ReadFile(imagePath)
	if err != nil {
		log.Fatalf("Error reading image file at %s: %v", imagePath, err)
	}

	// --- DETERMINE FORMAT ---
	ext := strings.ToLower(filepath.Ext(imagePath))
	format := "png"
	if ext == ".jpg" || ext == ".jpeg" {
		format = "jpeg"
	} else if ext == ".webp" {
		format = "webp"
	}

	// --- DEFINE PROMPT ---
	prompt := []genai.Part{
		genai.ImageData(format, imgData),
		genai.Text(`
			You are an expert Nepali Transcription and formatting engine.
			
			Task: Extract the text from this image and format it as clean HTML.

			Formatting Rules:
			1. Use <h1>, <h2>, or <h3> tags for headings.
			2. Use <p> tags for paragraphs.
			3. Use <br> if there are specific line breaks that need preserving within a paragraph.
			4. If there are lists, use <ol> or <ul> tags.
			5. If text is bold in the image, use <strong> tags.
			
			Strict Output Rule:
			Output strictly valid HTML. Do not wrap it in markdown backticks.
		`),
	}

	// --- RUN INFERENCE ---
	fmt.Printf("Processing image with %s... (Please wait)\n", modelName)
	resp, err := model.GenerateContent(ctx, prompt...)
	if err != nil {
		if strings.Contains(err.Error(), "400") || strings.Contains(err.Error(), "multimodal") {
			log.Fatalf("Error: It appears '%s' might not support Image input (OCR). \nTry switching back to 'gemini-2.5-flash'.\nOriginal Error: %v", modelName, err)
		}
		log.Fatalf("API Error: %v", err)
	}

	// --- OUTPUT RESULT ---
	printResponse(resp)
}

func printResponse(resp *genai.GenerateContentResponse) {
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				if txt, ok := part.(genai.Text); ok {
					rawText := string(txt)

					// --- FIX: Clean Markdown Backticks ---
					cleanText := cleanOutput(rawText)

					fmt.Println("\n================ HTML OUTPUT ================")
					fmt.Println(cleanText)
					fmt.Println("=============================================")

					saveHTML(cleanText)
				}
			}
		}
	}
}

// cleanOutput removes ```html, ```, and trims whitespace
func cleanOutput(input string) string {
	// Remove opening tag (e.g., ```html or ```)
	cleaned := strings.ReplaceAll(input, "```html", "")
	cleaned = strings.ReplaceAll(cleaned, "```", "")

	// Trim leading/trailing whitespace/newlines left over
	return strings.TrimSpace(cleaned)
}

func saveHTML(content string) {
	f, err := os.Create("output.html")
	if err != nil {
		fmt.Println("Could not save HTML file:", err)
		return
	}
	defer f.Close()

	// Add a proper doctype wrapper so it opens nicely in a browser
	fullHTML := fmt.Sprintf(`<!DOCTYPE html>
<html lang="ne">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body { font-family: sans-serif; padding: 40px; max-width: 800px; margin: auto; }
        h1, h2, h3 { color: #2c3e50; }
        p { line-height: 1.6; font-size: 18px; }
    </style>
    <title>Nepali OCR Output</title>
</head>
<body>
%s
</body>
</html>`, content)

	f.WriteString(fullHTML)
	fmt.Println("Saved cleaned result to 'output.html'")
}
