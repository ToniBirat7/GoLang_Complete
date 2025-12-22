package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ToniBirat7/tesseract_ocr_ne/pkg/ocr"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

const (
	maxUploadSize = 10 * 1024 * 1024 // 10MB
	uploadDir     = "./uploads"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
	Time    string `json:"time"`
}

func main() {
	// Create upload directory
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName:               "Go Tesseract OCR API",
		BodyLimit:             maxUploadSize,
		DisableStartupMessage: false,
		ErrorHandler:          customErrorHandler,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Routes
	app.Get("/", handleRoot)
	app.Get("/health", handleHealth)
	app.Post("/ocr/extract", handleOCRExtract)

	// Start server
	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// handleRoot returns API information
func handleRoot(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"name":        "Go Tesseract OCR API",
		"version":     "1.0.0",
		"description": "Nepali text extraction API using Tesseract OCR",
		"endpoints": fiber.Map{
			"health": "GET /health",
			"ocr":    "POST /ocr/extract",
		},
	})
}

// handleHealth returns health status
func handleHealth(c *fiber.Ctx) error {
	return c.JSON(HealthResponse{
		Status:  "healthy",
		Version: "1.0.0",
		Time:    time.Now().Format(time.RFC3339),
	})
}

// handleOCRExtract processes uploaded image and returns OCR results
func handleOCRExtract(c *fiber.Ctx) error {
	// Parse multipart form
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "bad_request",
			Message: "No image file provided. Use 'image' field in multipart/form-data",
		})
	}

	// Validate file size
	if file.Size > maxUploadSize {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "file_too_large",
			Message: fmt.Sprintf("File size exceeds maximum limit of %d MB", maxUploadSize/(1024*1024)),
		})
	}

	// Validate file extension
	ext := filepath.Ext(file.Filename)
	if !isValidImageExtension(ext) {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "invalid_file_type",
			Message: "Only .png, .jpg, .jpeg image files are supported",
		})
	}

	// Save uploaded file temporarily
	timestamp := time.Now().UnixNano()
	tempFilePath := filepath.Join(uploadDir, fmt.Sprintf("%d%s", timestamp, ext))
	if err := c.SaveFile(file, tempFilePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "upload_failed",
			Message: "Failed to save uploaded file",
		})
	}
	defer os.Remove(tempFilePath) // Clean up after processing

	// Parse optional parameters
	includeLines := c.FormValue("include_lines") == "true"

	// Configure OCR
	config := ocr.DefaultConfig()
	config.IncludeLines = includeLines

	// Perform OCR
	result, err := ocr.ExtractFromImage(tempFilePath, config)
	if err != nil {
		log.Printf("OCR extraction error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "ocr_failed",
			Message: "Failed to extract text from image",
		})
	}

	// Return JSON result
	return c.JSON(result)
}

// isValidImageExtension checks if file extension is valid
func isValidImageExtension(ext string) bool {
	validExts := map[string]bool{
		".png":  true,
		".jpg":  true,
		".jpeg": true,
		".PNG":  true,
		".JPG":  true,
		".JPEG": true,
	}
	return validExts[ext]
}

// customErrorHandler handles Fiber errors
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(ErrorResponse{
		Error:   "internal_error",
		Message: err.Error(),
	})
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
