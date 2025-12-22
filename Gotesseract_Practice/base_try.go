package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/otiai10/gosseract/v2"
)

type OCRWord struct {
	Text       string
	Confidence float64
	X          int
	Y          int
	Width      int
	Height     int
}

func main() {
	client := gosseract.NewClient()
	defer client.Close()

	client.SetImage("1.png")
	client.SetLanguage("nep")
	client.SetPageSegMode(gosseract.PSM_SINGLE_BLOCK)

	client.SetConfigFile("tsv")

	tsv, err := client.Text()
	if err != nil {
		panic(err)
	}

	// SAVE RAW TSV
	if err := os.WriteFile("output.tsv", []byte(tsv), 0644); err != nil {
		panic(err)
	}

	words := parseTSV(tsv)

	for _, w := range words {
		fmt.Printf(
			"Word: %-15s | Confidence: %6.2f | Box: (%d,%d,%d,%d)\n",
			w.Text,
			w.Confidence,
			w.X, w.Y, w.Width, w.Height,
		)
	}

	saveWords(words, "output_words.txt")
}

func parseTSV(tsv string) []OCRWord {
	lines := strings.Split(tsv, "\n")
	var results []OCRWord

	for i, line := range lines {
		if i == 0 || strings.TrimSpace(line) == "" {
			continue
		}

		fields := strings.Split(line, "\t")
		if len(fields) < 12 {
			continue
		}

		level, _ := strconv.Atoi(fields[0])

		// Level 5 = WORD
		if level != 5 {
			continue
		}

		x, _ := strconv.Atoi(fields[6])
		y, _ := strconv.Atoi(fields[7])
		w, _ := strconv.Atoi(fields[8])
		h, _ := strconv.Atoi(fields[9])

		confidence, _ := strconv.ParseFloat(fields[10], 64)
		text := strings.TrimSpace(fields[11])

		if text == "" {
			continue
		}

		results = append(results, OCRWord{
			Text:       text,
			Confidence: confidence,
			X:          x,
			Y:          y,
			Width:      w,
			Height:     h,
		})
	}

	return results
}

func saveWords(words []OCRWord, filename string) {
	var sb strings.Builder

	for _, w := range words {
		sb.WriteString(fmt.Sprintf("%s\t%.2f\n", w.Text, w.Confidence))
	}

	if err := os.WriteFile(filename, []byte(sb.String()), 0644); err != nil {
		panic(err)
	}
}
