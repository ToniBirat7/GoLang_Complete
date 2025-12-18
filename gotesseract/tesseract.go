package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg" // Added JPEG support just in case
	_ "image/png"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/otiai10/gosseract/v2"
)

type Block struct {
	Text        string
	X, Y, W, H  int
	StrokeWidth float64
}

type Article struct {
	Headline string
	Body     string
}

func main() {
	imagePath := "nep4.png" // Ensure this matches your filename exactly

	// 1. Load Image
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatal("Could not open file: ", err)
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		log.Fatal("Could not decode image: ", err)
	}
	fmt.Printf("Image loaded successfully. Format: %s, Bounds: %v\n", format, img.Bounds())

	// 2. Tesseract Setup
	client := gosseract.NewClient()
	defer client.Close()
	client.SetLanguage("nep", "eng")
	client.SetImage(imagePath)
	client.SetPageSegMode(gosseract.PSM_AUTO)

	// 3. Get Boxes
	fmt.Println("Extracting bounding boxes...")
	boxes, err := client.GetBoundingBoxes(gosseract.RIL_TEXTLINE)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d raw text lines.\n", len(boxes))

	var cleanBoxes []Block
	maxX := 0
	nepaliRegex := regexp.MustCompile(`[\p{Devanagari}]`)

	// 4. Analyze Blocks
	fmt.Println("--- START ANALYSIS ---")
	for _, b := range boxes {
		if b.Box.Max.X > maxX {
			maxX = b.Box.Max.X
		}

		txt := strings.TrimSpace(b.Word)

		// SKIP LOGIC: Top Menu
		if b.Box.Min.Y < 80 {
			continue
		}
		// SKIP LOGIC: Empty
		if txt == "" {
			continue
		}

		// DEBUG: Print what Tesseract sees before filtering
		// fmt.Printf("Raw: %s\n", txt)

		// SKIP LOGIC: Must have Nepali
		if !nepaliRegex.MatchString(txt) {
			continue
		}

		// SKIP LOGIC: UI Noise
		if isNoise(txt) {
			continue
		}

		// Calculate Metrics
		width := b.Box.Max.X - b.Box.Min.X
		height := b.Box.Max.Y - b.Box.Min.Y
		sw := calculateStrokeWidth(img, b.Box.Min.X, b.Box.Min.Y, width, height)

		// DEBUG LOG: See the values!
		// If SW is 0.00, the pixel logic isn't finding black pixels.
		fmt.Printf("Text: %-20s... | Height: %d | StrokeWidth: %.2f\n",
			string([]rune(txt)[:min(len([]rune(txt)), 15)]), height, sw)

		cleanBoxes = append(cleanBoxes, Block{
			Text: txt, X: b.Box.Min.X, Y: b.Box.Min.Y, W: width, H: height, StrokeWidth: sw,
		})
	}
	fmt.Println("--- END ANALYSIS ---")

	// 5. Zone Logic
	pageCenter := maxX / 2
	bufferZone := 50
	var fullWidthRows, leftColRows, rightColRows []Block

	for _, b := range cleanBoxes {
		if b.X < (pageCenter-bufferZone) && (b.X+b.W) > (pageCenter+bufferZone) {
			fullWidthRows = append(fullWidthRows, b)
		} else if (b.X + b.W) < (pageCenter + bufferZone) {
			leftColRows = append(leftColRows, b)
		} else {
			rightColRows = append(rightColRows, b)
		}
	}

	// 6. Group Articles
	var allArticles []Article
	fmt.Printf("Processing Zones: Full(%d), Left(%d), Right(%d)\n",
		len(fullWidthRows), len(leftColRows), len(rightColRows))

	allArticles = append(allArticles, processZone(fullWidthRows)...)
	allArticles = append(allArticles, processZone(leftColRows)...)
	allArticles = append(allArticles, processZone(rightColRows)...)

	// 7. Write Output
	writeOutput(allArticles)
}

func calculateStrokeWidth(img image.Image, x, y, w, h int) float64 {
	bounds := img.Bounds()
	// Safety Check
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	if x+w > bounds.Max.X {
		w = bounds.Max.X - x
	}
	if y+h > bounds.Max.Y {
		h = bounds.Max.Y - y
	}

	totalStrokeLen := 0
	strokeCount := 0

	// Scan middle 50%
	startY := y + (h / 4)
	endY := y + (h * 3 / 4)

	for currY := startY; currY < endY; currY++ { // Scan EVERY line (more accurate)
		isBlack := false
		currentRun := 0

		for currX := x; currX < x+w; currX++ {
			if isDarkPixel(img.At(currX, currY)) {
				if !isBlack {
					isBlack = true
					currentRun = 1
				} else {
					currentRun++
				}
			} else {
				if isBlack {
					isBlack = false
					// RELAXED FILTER: Accept strokes from 1px to 25px
					if currentRun >= 1 && currentRun < 25 {
						totalStrokeLen += currentRun
						strokeCount++
					}
				}
			}
		}
	}

	if strokeCount == 0 {
		return 0
	}
	return float64(totalStrokeLen) / float64(strokeCount)
}

func isDarkPixel(c color.Color) bool {
	r, g, b, _ := c.RGBA()
	// Standard Luminance Formula
	// RGBA corresponds to 0-65535.
	// Dark grey text is usually < 50% luminance.
	y := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
	return y < 40000 // Relaxed threshold (Higher number = lighter greys accepted)
}

func processZone(boxes []Block) []Article {
	sort.Slice(boxes, func(i, j int) bool { return boxes[i].Y < boxes[j].Y })

	var articles []Article
	var currentHeadline string
	var currentBodyBuilder strings.Builder
	var lastY int

	for _, b := range boxes {
		// LOGIC: Headline if Height > 22 OR Bold > 2.2
		isHeadline := b.H > 22 || b.StrokeWidth > 2.2

		// Gap Check
		if lastY > 0 && (b.Y-lastY) > 120 {
			if currentHeadline != "" || currentBodyBuilder.Len() > 0 {
				articles = append(articles, Article{currentHeadline, currentBodyBuilder.String()})
			}
			currentHeadline = ""
			currentBodyBuilder.Reset()
		}

		if isHeadline {
			if currentHeadline != "" || currentBodyBuilder.Len() > 0 {
				articles = append(articles, Article{currentHeadline, currentBodyBuilder.String()})
			}
			currentHeadline = b.Text
			currentBodyBuilder.Reset()
		} else {
			// If we have a headline, add to it.
			// FALLBACK: If we don't have a headline, treat this first text as a headline
			// if it's the very first item in a block.
			if currentHeadline == "" && currentBodyBuilder.Len() == 0 {
				currentHeadline = b.Text // Treat orphan text as headline to ensure capture
			} else {
				currentBodyBuilder.WriteString(b.Text + " ")
			}
		}
		lastY = b.Y + b.H
	}

	if currentHeadline != "" || currentBodyBuilder.Len() > 0 {
		articles = append(articles, Article{currentHeadline, currentBodyBuilder.String()})
	}
	return articles
}

func isNoise(text string) bool {
	noise := []string{"सेयर", "संग्रह", "Login", "कमेन्ट", "साझेदारी", "मिनेट", "अगाडि"}
	for _, n := range noise {
		if strings.Contains(text, n) && len([]rune(text)) < 20 {
			return true
		}
	}
	return false
}

func writeOutput(articles []Article) {
	f, _ := os.Create("final_debug_output.txt")
	defer f.Close()

	count := 0
	for _, art := range articles {
		head := strings.TrimSpace(art.Headline)
		body := strings.TrimSpace(art.Body)

		// RELAXED FILTER: Write even if only headline or only body exists
		if head == "" && body == "" {
			continue
		}

		out := fmt.Sprintf("HEADLINE: %s\nCONTENT: %s\n%s\n", head, body, strings.Repeat("-", 30))
		f.WriteString(out)
		count++
	}
	fmt.Printf("Successfully wrote %d articles to final_debug_output.txt\n", count)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
