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

	// Use the model you found in your list
	modelName := "gemini-2.5-flash"

	// --- SETUP ---
	err := godotenv.Load()
	if err != nil {
		log.Println("Note: .env file not found, checking system environment variables...")
	}

	apiKey := os.Getenv("FLASH_API_KEY")
	if apiKey == "" {
		log.Fatal("Error: FLASH_API_KEY environment variable is not set")
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

	// --- FIX: DETERMINE FORMAT (NOT MIME TYPE) ---
	// The library expects "png", "jpeg", or "webp".
	// It automatically adds "image/" before it.
	ext := strings.ToLower(filepath.Ext(imagePath))
	format := "png" // Default

	if ext == ".jpg" || ext == ".jpeg" {
		format = "jpeg"
	} else if ext == ".webp" {
		format = "webp"
	} else if ext == ".png" {
		format = "png"
	}

	// --- DEFINE PROMPT ---
	prompt := []genai.Part{
		// FIX: Pass 'format' (e.g. "png"), NOT "image/png"
		genai.ImageData(format, imgData),
		genai.Text(`
			You are an expert Nepali OCR system.
			Task: Transcribe the handwritten Nepali text in this image.

			The format of the document should be maintained in the output. 

			Use HTML tags to wrap the extracted text. Use proper heading and paragraphs to keep the format intact i.e. exactly similar to the provided document.

			Strict Rules:
			1. Output ONLY the Nepali text. No introductions, no markdown, no English.
			2. Maintain the original structure (newlines, paragraphs).
			3. Preserve all punctuation (|, ?, !, etc.) and bullet points (१., क., etc.).
			4. If a word is cut off or messy, use context to infer the correct Nepali word.
			5. Use <h1> for bold text according to the provided document. Also maintain the spaces and line breaks to separate the paragraphs.
			6. The heading tags also should follow hierarchy based on the boldness of the font. 
		`),
	}

	// --- RUN INFERENCE ---
	fmt.Printf("Processing image with %s... (Please wait)\n", modelName)
	resp, err := model.GenerateContent(ctx, prompt...)
	if err != nil {
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
					fmt.Println("\n================ OCR OUTPUT ================")
					fmt.Println(string(txt))
					fmt.Println("============================================")
				}
			}
		}
	}
}
