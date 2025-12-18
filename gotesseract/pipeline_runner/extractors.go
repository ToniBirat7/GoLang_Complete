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

// Word-level extraction
func ExtractWordLevel(imagePath, outputDir string) error {
	outputFolder := filepath.Join(outputDir, "extraction_results_word_level")
	outputInfoFile := "words_data.txt"
	outputImageFile := "mapped_image.png"

	client := gosseract.NewClient()
	defer client.Close()
	client.SetImage(imagePath)
	client.SetLanguage("nep")

	boundingBoxes, err := client.GetBoundingBoxes(gosseract.RIL_WORD)
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
	boxColor := color.RGBA{255, 0, 0, 255}
	for i, box := range boundingBoxes {
		line := fmt.Sprintf("[%d] Word: %s | Confidence: %.2f\n", i, box.Word, box.Confidence)
		txtFile.WriteString(line)
		rect := box.Box
		for x := rect.Min.X; x < rect.Max.X; x++ {
			rgbaImg.Set(x, rect.Min.Y, boxColor)
			rgbaImg.Set(x, rect.Max.Y-1, boxColor)
			rgbaImg.Set(x, rect.Min.Y+1, boxColor)
			rgbaImg.Set(x, rect.Max.Y-2, boxColor)
		}
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			rgbaImg.Set(rect.Min.X, y, boxColor)
			rgbaImg.Set(rect.Max.X-1, y, boxColor)
			rgbaImg.Set(rect.Min.X+1, y, boxColor)
			rgbaImg.Set(rect.Max.X-2, y, boxColor)
		}
	}
	finalImgPath := filepath.Join(outputFolder, outputImageFile)
	f, err := os.Create(finalImgPath)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := png.Encode(f, rgbaImg); err != nil {
		return err
	}
	return nil
}

// Sentence-level extraction
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
	boxColor := color.RGBA{0, 0, 255, 255}
	for i, box := range boundingBoxes {
		line := fmt.Sprintf("[%d] Text: %s | Confidence: %.2f\n", i, box.Word, box.Confidence)
		txtFile.WriteString(line)
		rect := box.Box
		for x := rect.Min.X; x < rect.Max.X; x++ {
			rgbaImg.Set(x, rect.Min.Y, boxColor)
			rgbaImg.Set(x, rect.Max.Y-1, boxColor)
			rgbaImg.Set(x, rect.Min.Y+1, boxColor)
			rgbaImg.Set(x, rect.Max.Y-2, boxColor)
		}
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			rgbaImg.Set(rect.Min.X, y, boxColor)
			rgbaImg.Set(rect.Max.X-1, y, boxColor)
			rgbaImg.Set(rect.Min.X+1, y, boxColor)
			rgbaImg.Set(rect.Max.X-2, y, boxColor)
		}
	}
	finalImgPath := filepath.Join(outputFolder, outputImageFile)
	f, err := os.Create(finalImgPath)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := png.Encode(f, rgbaImg); err != nil {
		return err
	}
	return nil
}

// Paragraph-level extraction (custom grouping)
func ExtractParagraphLevel(imagePath, outputDir string) error {
	outputFolder := filepath.Join(outputDir, "extraction_results_custom_paras")
	outputInfoFile := "custom_paragraphs_data.txt"
	outputImageFile := "mapped_custom_paragraphs.png"
	gapThresholdRatio := 0.6

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
				if currBox.Min.X < currentPara.Box.Min.X {
					currentPara.Box.Min.X = currBox.Min.X
				}
				if currBox.Min.Y < currentPara.Box.Min.Y {
					currentPara.Box.Min.Y = currBox.Min.Y
				}
				if currBox.Max.X > currentPara.Box.Max.X {
					currentPara.Box.Max.X = currBox.Max.X
				}
				if currBox.Max.Y > currentPara.Box.Max.Y {
					currentPara.Box.Max.Y = currBox.Max.Y
				}
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
	boxColor := color.RGBA{128, 0, 128, 255}
	for i, para := range paragraphs {
		cleanText := strings.ReplaceAll(para.Text, "\n", " ")
		line := fmt.Sprintf("[%d] Para (Lines: %d): %s | Confidence: %.2f\n", i, para.LineCount, cleanText, para.Confidence)
		txtFile.WriteString(line)
		rect := para.Box
		minX, minY := int(math.Max(0, float64(rect.Min.X))), int(math.Max(0, float64(rect.Min.Y)))
		maxX, maxY := int(math.Min(float64(bounds.Max.X), float64(rect.Max.X))), int(math.Min(float64(bounds.Max.Y), float64(rect.Max.Y)))
		for x := minX; x < maxX; x++ {
			rgbaImg.Set(x, minY, boxColor)
			rgbaImg.Set(x, maxY-1, boxColor)
			rgbaImg.Set(x, minY+1, boxColor)
			rgbaImg.Set(x, maxY-2, boxColor)
		}
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
		return err
	}
	defer f.Close()
	if err := png.Encode(f, rgbaImg); err != nil {
		return err
	}
	return nil
}
