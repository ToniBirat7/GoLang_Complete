package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"

	// Import these to ensure Go can decode the input image formats
	_ "image/jpeg"
	_ "image/png"

	"github.com/otiai10/gosseract/v2"
)

func main() {
	// --- CONFIGURATION ---
	inputImagePath := "../output/bar-13.png"
	outputFolder := "extraction_results_word_level"
	outputInfoFile := "words_data.txt"
	outputImageFile := "mapped_image.png"

	// 1. Initialize Gosseract Client
	client := gosseract.NewClient()
	defer client.Close()

	client.SetImage(inputImagePath)
	client.SetLanguage("nep")

	// 2. Get Bounding Boxes at WORD level
	boundingBoxes, err := client.GetBoundingBoxes(gosseract.RIL_WORD)
	if err != nil {
		panic(err)
	}

	// 3. Create the output directory
	if err := os.MkdirAll(outputFolder, 0755); err != nil {
		panic(err)
	}

	// 4. Create the text file
	txtFile, err := os.Create(filepath.Join(outputFolder, outputInfoFile))
	if err != nil {
		panic(err)
	}
	defer txtFile.Close()

	// 5. Open the original image
	srcFile, err := os.Open(inputImagePath)
	if err != nil {
		panic(err)
	}
	defer srcFile.Close()

	srcImg, _, err := image.Decode(srcFile)
	if err != nil {
		panic(err)
	}

	// 6. Convert to RGBA (Mutable Image)
	// Standard image.Decode might return read-only types (like YCbCr).
	// We create a new RGBA image and draw the original onto it so we can edit pixels.
	bounds := srcImg.Bounds()
	rgbaImg := image.NewRGBA(bounds)
	draw.Draw(rgbaImg, bounds, srcImg, bounds.Min, draw.Src)

	// Define the color for the bounding box (Red)
	boxColor := color.RGBA{255, 0, 0, 255}

	fmt.Println("Processing extracted words and drawing boxes...")

	for i, box := range boundingBoxes {
		// --- A. Save Text Data ---
		line := fmt.Sprintf("[%d] Word: %s | Confidence: %.2f\n", i, box.Word, box.Confidence)
		txtFile.WriteString(line)

		// --- B. Draw Bounding Box on the Single Image ---
		rect := box.Box

		// Draw Top and Bottom lines
		for x := rect.Min.X; x < rect.Max.X; x++ {
			rgbaImg.Set(x, rect.Min.Y, boxColor)   // Top line
			rgbaImg.Set(x, rect.Max.Y-1, boxColor) // Bottom line

			// Optional: Make lines thicker (2px)
			rgbaImg.Set(x, rect.Min.Y+1, boxColor)
			rgbaImg.Set(x, rect.Max.Y-2, boxColor)
		}

		// Draw Left and Right lines
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			rgbaImg.Set(rect.Min.X, y, boxColor)   // Left line
			rgbaImg.Set(rect.Max.X-1, y, boxColor) // Right line

			// Optional: Make lines thicker (2px)
			rgbaImg.Set(rect.Min.X+1, y, boxColor)
			rgbaImg.Set(rect.Max.X-2, y, boxColor)
		}
	}

	// 7. Save the final annotated image
	finalImgPath := filepath.Join(outputFolder, outputImageFile)
	f, err := os.Create(finalImgPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := png.Encode(f, rgbaImg); err != nil {
		panic(err)
	}

	fmt.Printf("Done! Text saved to '%s' and image saved to '%s'\n", outputInfoFile, outputImageFile)
}
