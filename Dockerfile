# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-linkmode external -extldflags "-static"' -o app ./cmd/api

# Runtime stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache \
    sqlite \
    curl \
    iputils \
    bash \
    libxml2-utils \
    libxslt \
    && rm -rf /var/cache/apk/*

WORKDIR /app

# Create necessary directories
RUN mkdir -p /app/data /app/uploads /app/files /app/images /app/static /app/exports /app/backups

# Copy binary from builder
COPY --from=builder /build/app /app/app

# Set environment variables
ENV PORT=8080
ENV DB_PATH=/app/data/app.db
ENV UPLOAD_DIR=/app/uploads
ENV FILES_DIR=/app/files
ENV IMAGES_DIR=/app/images
ENV STATIC_DIR=/app/static
ENV EXPORTS_DIR=/app/exports
ENV BACKUPS_DIR=/app/backups

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Run the application
CMD ["/app/app"]
