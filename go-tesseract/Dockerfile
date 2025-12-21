# Multi-stage build for smaller image size
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git tesseract-ocr-dev leptonica-dev gcc g++ musl-dev

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the API binary
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o ocr-api ./cmd/ocr-api

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache \
    tesseract-ocr \
    tesseract-ocr-data-nep \
    ca-certificates \
    && rm -rf /var/cache/apk/*

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/ocr-api .

# Create uploads directory
RUN mkdir -p /app/uploads

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the API
CMD ["./ocr-api"]
