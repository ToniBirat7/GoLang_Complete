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
**Text Cleaning** - Automatic removal of OCR artifacts  

## Installation

### As a Library

```bash
go get github.com/ToniBirat7/tesseract_ocr_ne
```

### Prerequisites

- Go 1.23 or higher
- Tesseract OCR installed
- Nepali language data for Tesseract

#### Install Tesseract

**Ubuntu/Debian:**
```bash
sudo apt-get update
sudo apt-get install tesseract-ocr tesseract-ocr-nep libtesseract-dev libleptonica-dev
```

**macOS:**
```bash
brew install tesseract tesseract-lang
```

**Windows:**
Download and install from [Tesseract GitHub](https://github.com/UB-Mannheim/tesseract/wiki)

## Usage

### 1. As a Library

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/ToniBirat7/tesseract_ocr_ne/pkg/ocr"
)

func main() {
    // Use default configuration
    config := ocr.DefaultConfig()
    
    // Or customize
    config.IncludeLines = true
    config.MinConfidence = 50.0
    
    // Extract text from image
    result, err := ocr.ExtractFromImage("path/to/image.png", config)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Text: %s\n", result.Text)
    fmt.Printf("Confidence: %.2f%%\n", result.AverageConfidence)
    fmt.Printf("Lines: %d\n", result.LineCount)
}
```

### 2. CLI Tool

Build the CLI:
```bash
go build -o ocr-cli ./cmd/ocr-cli
```

Basic usage:
```bash
# Print to stdout
./ocr-cli -image path/to/image.png

# Save to file
./ocr-cli -image path/to/image.png -output result.json

# Include individual lines
./ocr-cli -image path/to/image.png -lines

# Set minimum confidence threshold
./ocr-cli -image path/to/image.png -min-confidence 60
```

**CLI Flags:**
- `-image` (required) - Path to image file
- `-output` (optional) - Path to output JSON file
- `-lines` (optional) - Include individual lines in output
- `-min-confidence` (optional) - Minimum confidence threshold (0-100)

### 3. HTTP API

#### Run Locally

```bash
go run ./cmd/ocr-api
```

The API will start on `http://localhost:8080`

#### Endpoints

**GET /** - API information
```bash
curl http://localhost:8080/
```

**GET /health** - Health check
```bash
curl http://localhost:8080/health
```

**POST /ocr/extract** - Extract text from image
```bash
curl -X POST http://localhost:8080/ocr/extract \
  -F "image=@path/to/image.png" \
  -F "include_lines=true"
```

**Request Parameters:**
- `image` (required) - Image file (multipart/form-data)
- `include_lines` (optional) - Set to "true" to include individual lines

**Response Example:**
```json
{
  "text": "नेपाली पाठ यहाँ छ",
  "average_confidence": 95.5,
  "line_count": 3,
  "lines": [
    {
      "text": "नेपाली पाठ",
      "confidence": 96.2
    },
    {
      "text": "यहाँ छ",
      "confidence": 94.8
    }
  ]
}
```

### 4. Docker Deployment

#### Using Docker Compose (Recommended)

```bash
# Build and start
docker-compose up -d

# View logs
docker-compose logs -f

# Stop
docker-compose down
```

#### Using Docker directly

```bash
# Build image
docker build -t go-tesseract-ocr .

# Run container
docker run -d -p 8080:8080 --name ocr-api go-tesseract-ocr

# Test
curl http://localhost:8080/health
```

## API Response Structure

```go
type OCRResult struct {
    Text              string          `json:"text"`
    AverageConfidence float64         `json:"average_confidence"`
    LineCount         int             `json:"line_count"`
    Lines             []ExtractedLine `json:"lines,omitempty"`
}

type ExtractedLine struct {
    Text       string  `json:"text"`
    Confidence float64 `json:"confidence"`
}
```

## Configuration

```go
type OCRConfig struct {
    Language          string  // Default: "nep"
    IncludeLines      bool    // Default: false
    CleanDevanagari   bool    // Default: true
    MinConfidence     float64 // Default: 0.0
}
```

## Development

### Project Structure

```
go-tesseract/
├── cmd/
│   ├── ocr-api/          # HTTP API server
│   │   └── main.go
│   └── ocr-cli/          # CLI tool
│       └── main.go
├── pkg/
│   └── ocr/              # Core library (importable)
│       ├── ocr.go
│       └── types.go
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── README.md
```

### Build All Binaries

```bash
# Build CLI
go build -o bin/ocr-cli ./cmd/ocr-cli

# Build API
go build -o bin/ocr-api ./cmd/ocr-api
```

### Run Tests

```bash
go test ./...
```

### Install Dependencies

```bash
go mod download
go mod tidy
```

## Environment Variables

- `PORT` - API server port (default: 8080)

## Supported Image Formats

- PNG (.png)
- JPEG (.jpg, .jpeg)

## Error Handling

The API returns structured error responses:

```json
{
  "error": "error_code",
  "message": "Human-readable error message"
}
```

**Common Error Codes:**
- `bad_request` - Invalid request parameters
- `file_too_large` - File exceeds 10MB limit
- `invalid_file_type` - Unsupported file format
- `ocr_failed` - OCR processing error

## Performance

- Maximum file size: 10MB
- Supported concurrent requests: Based on system resources
- Average processing time: 1-3 seconds per image (depends on image size and complexity)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details

## Author

**Toni Birat**  
GitHub: [@ToniBirat7](https://github.com/ToniBirat7)

## Acknowledgments

- [Tesseract OCR](https://github.com/tesseract-ocr/tesseract)
- [gosseract](https://github.com/otiai10/gosseract)
- [Fiber](https://github.com/gofiber/fiber)

## Support

For issues and questions, please open an issue on GitHub.
