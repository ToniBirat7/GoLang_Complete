# Go Tesseract OCR

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A production-ready Go package for Nepali text extraction using Tesseract OCR. Provides both a library package, CLI tool, and HTTP API.

## Features

**Library Package** - Import and use in your Go projects  

**CLI Tool** - Command-line interface for batch processing

**HTTP API** - RESTful API using Go Fiber  

**Docker Support** - Containerized deployment  

**JSON Output** - Structured results with confidence scores  

**Nepali Language** - Optimized for Devanagari script  

## **Prerequisites**

1. ``Tesseract OCR`` installed with Nepali language data.  

2. 
   - Installation instructions: `https://tesseract-ocr.github.io/tessdoc/Installation.html`  
   - Nepali language data: `nep.traineddata`

3. `Go` installed (version 1.23 or higher recommended). 
 
   - Download: `https://golang.org/dl/`

4. Install `WSL` (for Windows users)  

   - Instructions: `https://docs.microsoft.com/en-us/windows/wsl/install`

## Installation

### As a Library

```bash
go get github.com/ToniBirat7/tesseract_ocr_ne
go mod tidy
```

```bash
# Inside your Go project

package main

import (
  "fmt"
  "github.com/ToniBirat7/tesseract_ocr_ne/pkg/ocr"
)

func main() {
  cfg := ocr.DefaultConfig()
  cfg.IncludeLines = true
  res, err := ocr.ExtractFromImage("path/to/image.png", cfg)
  if err != nil { panic(err) }
  fmt.Println(res.Text)
}
```

## Installation & Usage (Docker Recommended)

The easiest way to use this package is via Docker, which includes all dependencies (Tesseract, Nepali language data, etc.).

First, clone the repository:

```bash
git clone git@github.com:ToniBirat7/tesseract_ocr_ne.git .
```

### 1. Build the Docker Image
```bash
docker compose build
```

### 2. Run the HTTP API
Start the API server on port 8080:
```bash
docker compose up -d
```

Get OCR:
```bash
curl -s -X POST \
  -F "image=@test_img/img.png" \
  -F "include_lines=true" \
  http://localhost:8080/ocr/extract
```

### 3. Use the CLI Tool

```bash
go run ./cmd/ocr-cli -image test_img/img.png
```

## API Response Structure

```json
{
  "text": "नेपाली पाठ यहाँ छ\n\nअर्को लाइन",
  "average_confidence": 95.5,
  "line_count": 2
}
```

## Performance

- Maximum file size: 10MB

- Supported concurrent requests: Based on system resources

- Average processing time: 1-3 seconds per image (depends on image size and complexity)

## License

MIT License - see LICENSE file for details

## Author

**Toni Birat**  
GitHub: [@ToniBirat7](https://github.com/ToniBirat7)

## Acknowledgments

- [Tesseract OCR](https://github.com/tesseract-ocr/tesseract)

- [gosseract](https://github.com/otiai10/gosseract)

- [Fiber](https://github.com/gofiber/fiber)
