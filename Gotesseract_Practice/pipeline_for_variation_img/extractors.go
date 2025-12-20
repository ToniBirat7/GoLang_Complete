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

// ExtractWordLevel performs OCR at the word level and returns Average Confidence
func ExtractWordLevel(imagePath, outputDir string) (float64, error) {
	outputFolder := filepath.Join(outputDir, "extraction_results_word_level")
	outputInfoFile := "words_data.txt"
	outputImageFile := "mapped_image.png"

	client := gosseract.NewClient()
	defer client.Close()
	client.SetImage(imagePath)
	client.SetLanguage("nep")

	// Get Word Boxes
	boundingBoxes, err := client.GetBoundingBoxes(gosseract.RIL_WORD)
	if err != nil {
		return 0, err
	}

	if err := os.MkdirAll(outputFolder, 0755); err != nil {
		return 0, err
	}

	txtFile, err := os.Create(filepath.Join(outputFolder, outputInfoFile))
	if err != nil {
		return 0, err
	}
	defer txtFile.Close()

	srcFile, err := os.Open(imagePath)
	if err != nil {
		return 0, err
	}
	defer srcFile.Close()
	srcImg, _, err := image.Decode(srcFile)
	if err != nil {
		return 0, err
	}

	bounds := srcImg.Bounds()
	rgbaImg := image.NewRGBA(bounds)
	draw.Draw(rgbaImg, bounds, srcImg, bounds.Min, draw.Src)
	boxColor := color.RGBA{255, 0, 0, 255} // RED for Words

	var totalConfidence float64
	var wordCount int

	for i, box := range boundingBoxes {
		cleanWord := strings.TrimSpace(box.Word)
		if cleanWord != "" {
			line := fmt.Sprintf("[%d] || %s || %.2f\n", i, cleanWord, box.Confidence)
			txtFile.WriteString(line)

			// Accumulate confidence stats
			totalConfidence += box.Confidence
			wordCount++
		}

		rect := box.Box
		drawBox(rgbaImg, rect, boxColor)
	}

	// Calculate Average Confidence
	avgConfidence := 0.0
	if wordCount > 0 {
		avgConfidence = totalConfidence / float64(wordCount)
	}

	err = saveImage(rgbaImg, filepath.Join(outputFolder, outputImageFile))
	return avgConfidence, err
}

// ExtractSentenceLevel performs OCR at the line/sentence level
func ExtractSentenceLevel(imagePath, outputDir string) error {
	outputFolder := filepath.Join(outputDir, "extraction_results_sentence_level")
	outputInfoFile := "sentences_data.txt"
	outputImageFile := "mapped_sentences.png"

	client := gosseract.NewClient()
	defer client.Close()
	client.SetImage(imagePath)
	client.SetLanguage("nep")

	boundingBoxes, err := client.GetBoundingBoxes(gosseract.RIL_TEXTLINE)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(outputFolder, 0755); err != nil {
		return err
	}

	txtFile, err := os.Create(filepath.Join(outputFolder, outputInfoFile))
	if err != nil {
		return err
	}
	defer txtFile.Close()

	srcFile, err := os.Open(imagePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	srcImg, _, err := image.Decode(srcFile)
	if err != nil {
		return err
	}

	bounds := srcImg.Bounds()
	rgbaImg := image.NewRGBA(bounds)
	draw.Draw(rgbaImg, bounds, srcImg, bounds.Min, draw.Src)
	boxColor := color.RGBA{0, 0, 255, 255} // BLUE for Sentences

	for i, box := range boundingBoxes {
		cleanText := strings.TrimSpace(box.Word)
		if cleanText != "" {
			line := fmt.Sprintf("[%d] || %s || %.2f\n", i, cleanText, box.Confidence)
			txtFile.WriteString(line)
		}
		drawBox(rgbaImg, box.Box, boxColor)
	}

	return saveImage(rgbaImg, filepath.Join(outputFolder, outputImageFile))
}

// ExtractParagraphLevel performs OCR and manually groups lines into paragraphs
func ExtractParagraphLevel(imagePath, outputDir string) error {
	outputFolder := filepath.Join(outputDir, "extraction_results_paragraph_level")
	outputInfoFile := "custom_paragraphs_data.txt"
	outputImageFile := "mapped_custom_paragraphs.png"
	// Sensitivity: If gap > 60% of line height, it's a new paragraph
	gapThresholdRatio := 0.60

	type CustomParagraph struct {
		Box        image.Rectangle
		Text       string
		Confidence float64
		LineCount  int
	}

	client := gosseract.NewClient()
	defer client.Close()
	client.SetImage(imagePath)
	client.SetLanguage("nep")

	lineBoxes, err := client.GetBoundingBoxes(gosseract.RIL_TEXTLINE)
	if err != nil {
		return err
	}

	var paragraphs []CustomParagraph

	if len(lineBoxes) > 0 {
		currentPara := CustomParagraph{
			Box:        lineBoxes[0].Box,
			Text:       lineBoxes[0].Word,
			Confidence: lineBoxes[0].Confidence,
			LineCount:  1,
		}

		for i := 1; i < len(lineBoxes); i++ {
			prevBox := lineBoxes[i-1].Box
			currBox := lineBoxes[i].Box
			lineHeight := currBox.Max.Y - currBox.Min.Y
			verticalGap := currBox.Min.Y - prevBox.Max.Y

			if float64(verticalGap) > float64(lineHeight)*gapThresholdRatio {
				paragraphs = append(paragraphs, currentPara)
				currentPara = CustomParagraph{
					Box:        currBox,
					Text:       lineBoxes[i].Word,
					Confidence: lineBoxes[i].Confidence,
					LineCount:  1,
				}
			} else {
				currentPara.Box = currentPara.Box.Union(currBox)
				currentPara.Text += " " + lineBoxes[i].Word
				totalConf := (currentPara.Confidence * float64(currentPara.LineCount)) + lineBoxes[i].Confidence
				currentPara.LineCount++
				currentPara.Confidence = totalConf / float64(currentPara.LineCount)
			}
		}
		paragraphs = append(paragraphs, currentPara)
	}

	if err := os.MkdirAll(outputFolder, 0755); err != nil {
		return err
	}

	txtFile, err := os.Create(filepath.Join(outputFolder, outputInfoFile))
	if err != nil {
		return err
	}
	defer txtFile.Close()

	srcFile, err := os.Open(imagePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	srcImg, _, err := image.Decode(srcFile)
	if err != nil {
		return err
	}

	bounds := srcImg.Bounds()
	rgbaImg := image.NewRGBA(bounds)
	draw.Draw(rgbaImg, bounds, srcImg, bounds.Min, draw.Src)
	boxColor := color.RGBA{0, 128, 0, 255} // GREEN for Paragraphs

	for i, para := range paragraphs {
		cleanText := strings.ReplaceAll(para.Text, "\n", " ")
		cleanText = strings.TrimSpace(cleanText)

		if cleanText != "" {
			line := fmt.Sprintf("[%d] || %s || %.2f\n", i, cleanText, para.Confidence)
			txtFile.WriteString(line)
		}
		drawBox(rgbaImg, para.Box, boxColor)
	}

	return saveImage(rgbaImg, filepath.Join(outputFolder, outputImageFile))
}

func drawBox(img *image.RGBA, rect image.Rectangle, col color.RGBA) {
	bounds := img.Bounds()
	minX, minY := int(math.Max(0, float64(rect.Min.X))), int(math.Max(0, float64(rect.Min.Y)))
	maxX, maxY := int(math.Min(float64(bounds.Max.X), float64(rect.Max.X))), int(math.Min(float64(bounds.Max.Y), float64(rect.Max.Y)))

	for i := 0; i < 2; i++ {
		for x := minX; x < maxX; x++ {
			if minY+i < maxY {
				img.Set(x, minY+i, col)
			}
			if maxY-1-i > minY {
				img.Set(x, maxY-1-i, col)
			}
		}
		for y := minY; y < maxY; y++ {
			if minX+i < maxX {
				img.Set(minX+i, y, col)
			}
			if maxX-1-i > minX {
				img.Set(maxX-1-i, y, col)
			}
		}
	}
}

func saveImage(img image.Image, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}
