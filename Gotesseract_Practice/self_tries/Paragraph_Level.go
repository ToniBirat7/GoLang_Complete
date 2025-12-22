package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"strings"

	_ "image/jpeg"
	_ "image/png"

	"github.com/otiai10/gosseract/v2"
)

// Custom struct to hold our manually merged paragraph data
type CustomParagraph struct {
	Box        image.Rectangle
	Text       string
	Confidence float64
	LineCount  int
}

func main() {
	// --- CONFIGURATION ---
	inputImagePath := "../output/bar-06.png"
	outputFolder := "extraction_results_custom_paras"
	outputInfoFile := "custom_paragraphs_data.txt"
	outputImageFile := "mapped_custom_paragraphs.png"

	// SENSITIVITY SETTING:
	// If the vertical gap between lines is > (LineHeight * 0.60),
	// we consider it a new paragraph. Decrease this to split paragraphs more aggressively.
	gapThresholdRatio := 0.6

	// 1. Initialize Gosseract Client
	client := gosseract.NewClient()
	defer client.Close()

	client.SetImage(inputImagePath)
	client.SetLanguage("nep")

	// 2. Get Bounding Boxes at LINE level (We will group them manually)
	lineBoxes, err := client.GetBoundingBoxes(gosseract.RIL_TEXTLINE)
	if err != nil {
		panic(err)
	}

	// 3. Logic to Group Lines into Paragraphs
	var paragraphs []CustomParagraph

	if len(lineBoxes) > 0 {
		// Initialize the first paragraph with the first line
		currentPara := CustomParagraph{
			Box:        lineBoxes[0].Box,
			Text:       lineBoxes[0].Word,
			Confidence: lineBoxes[0].Confidence,
			LineCount:  1,
		}

		for i := 1; i < len(lineBoxes); i++ {
			prevBox := lineBoxes[i-1].Box
			currBox := lineBoxes[i].Box

			// Calculate the average height of the current line
			lineHeight := currBox.Max.Y - currBox.Min.Y

			// Calculate the vertical gap between this line and the previous one
			verticalGap := currBox.Min.Y - prevBox.Max.Y

			// DECISION: Is this a new paragraph?
			if float64(verticalGap) > float64(lineHeight)*gapThresholdRatio {
				// YES: Save the old paragraph and start a new one
				paragraphs = append(paragraphs, currentPara)

				currentPara = CustomParagraph{
					Box:        currBox,
					Text:       lineBoxes[i].Word,
					Confidence: lineBoxes[i].Confidence,
					LineCount:  1,
				}
			} else {
				// NO: Merge this line into the current paragraph

				// 1. Expand the rectangle to include the new line
				// Union logic: Min X/Y is the smaller of the two, Max X/Y is the larger
				if currBox.Min.X < currentPara.Box.Min.X {
					currentPara.Box.Min.X = currBox.Min.X
				}
				if currBox.Min.Y < currentPara.Box.Min.Y {
					currentPara.Box.Min.Y = currBox.Min.Y
				} // Should handle itself naturally
				if currBox.Max.X > currentPara.Box.Max.X {
					currentPara.Box.Max.X = currBox.Max.X
				}
				if currBox.Max.Y > currentPara.Box.Max.Y {
					currentPara.Box.Max.Y = currBox.Max.Y
				}

				// 2. Append text
				currentPara.Text += " " + lineBoxes[i].Word

				// 3. Average the confidence
				totalConf := (currentPara.Confidence * float64(currentPara.LineCount)) + lineBoxes[i].Confidence
				currentPara.LineCount++
				currentPara.Confidence = totalConf / float64(currentPara.LineCount)
			}
		}
		// Don't forget to append the very last paragraph being built
		paragraphs = append(paragraphs, currentPara)
	}

	// 4. Create Output Files
	if err := os.MkdirAll(outputFolder, 0755); err != nil {
		panic(err)
	}

	txtFile, err := os.Create(filepath.Join(outputFolder, outputInfoFile))
	if err != nil {
		panic(err)
	}
	defer txtFile.Close()

	srcFile, err := os.Open(inputImagePath)
	if err != nil {
		panic(err)
	}
	defer srcFile.Close()

	srcImg, _, err := image.Decode(srcFile)
	if err != nil {
		panic(err)
	}

	// Convert to RGBA
	bounds := srcImg.Bounds()
	rgbaImg := image.NewRGBA(bounds)
	draw.Draw(rgbaImg, bounds, srcImg, bounds.Min, draw.Src)

	// Color: Purple for Custom Paragraphs
	boxColor := color.RGBA{128, 0, 128, 255}

	fmt.Printf("Detected %d paragraphs manually.\n", len(paragraphs))

	// 5. Draw and Save
	for i, para := range paragraphs {
		// A. Write Text Info
		// Clean up newlines in text for single line display in txt file
		cleanText := strings.ReplaceAll(para.Text, "\n", " ")
		line := fmt.Sprintf("[%d] Para (Lines: %d): %s | Confidence: %.2f\n", i, para.LineCount, cleanText, para.Confidence)
		txtFile.WriteString(line)

		// B. Draw Box
		rect := para.Box

		// Ensure bounds don't crash (though they shouldn't if source is same)
		minX, minY := int(math.Max(0, float64(rect.Min.X))), int(math.Max(0, float64(rect.Min.Y)))
		maxX, maxY := int(math.Min(float64(bounds.Max.X), float64(rect.Max.X))), int(math.Min(float64(bounds.Max.Y), float64(rect.Max.Y)))

		// Draw Horizontal lines
		for x := minX; x < maxX; x++ {
			rgbaImg.Set(x, minY, boxColor)
			rgbaImg.Set(x, maxY-1, boxColor)
			rgbaImg.Set(x, minY+1, boxColor)
			rgbaImg.Set(x, maxY-2, boxColor)
		}

		// Draw Vertical lines
		for y := minY; y < maxY; y++ {
			rgbaImg.Set(minX, y, boxColor)
			rgbaImg.Set(maxX-1, y, boxColor)
			rgbaImg.Set(minX+1, y, boxColor)
			rgbaImg.Set(maxX-2, y, boxColor)
		}
	}

	finalImgPath := filepath.Join(outputFolder, outputImageFile)
	f, err := os.Create(finalImgPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := png.Encode(f, rgbaImg); err != nil {
		panic(err)
	}

	fmt.Printf("Done! Saved to '%s' folder.\n", outputFolder)
}
