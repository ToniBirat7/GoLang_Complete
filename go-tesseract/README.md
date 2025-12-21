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

## Installation & Usage (Docker Recommended)

The easiest way to use this package is via Docker, which includes all dependencies (Tesseract, Nepali language data, etc.).

### 1. Build the Docker Image
```bash
docker compose build
```

### 2. Run the HTTP API
Start the API server on port 8080:
```bash
docker compose up -d
```

Test it:
```bash
curl http://localhost:8080/health
```

### 3. Run the CLI Tool
Use the helper script `ocr-cli.bat` to run the CLI tool without installing anything locally.

**Usage:**
```powershell
.\ocr-cli.bat test_img/img.png
```

**Or using Docker directly:**
```bash
docker run --rm -v ${PWD}:/app/data -t go-tesseract-ocr-api ./ocr-cli -image /app/data/YOUR_IMAGE.png
```

## API Response Structure

```json
{
  "text": "नेपाली पाठ यहाँ छ\n\nअर्को लाइन",
  "average_confidence": 95.5,
  "line_count": 2
}
```

## API Response Structure

```json
{
  "text": "नेपाली पाठ यहाँ छ\n\nअर्को लाइन",
  "average_confidence": 95.5,
  "line_count": 2
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

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
