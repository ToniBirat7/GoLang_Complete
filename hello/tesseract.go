package main

import (
	"fmt"
	"log"
	"os"

	"github.com/otiai10/gosseract/v2"
)

func main() {
	client := gosseract.NewClient()
	defer client.Close()

	// Nepali + English
	client.SetLanguage("nep", "eng")

	// Set the image
	client.SetImage("nep4.png")

	// Optional: use paragraph mode for better accuracy
	client.SetPageSegMode(gosseract.PSM_SINGLE_BLOCK)

	text, err := client.Text()
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("output.txt", []byte(text), 0644)

	fmt.Println("OCR Output:")
	fmt.Println(text)
}
