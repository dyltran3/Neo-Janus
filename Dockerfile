# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY 3_janus_core/go.mod 3_janus_core/go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY 3_janus_core/ .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server/

# Runtime stage
FROM alpine:3.18

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server /app/

# Copy config
COPY config.yaml ./

# Create necessary directories
RUN mkdir -p ./data/logs ./data/vaccine

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:8080/health || exit 1

# Expose port
EXPOSE 8080

# Run server
CMD ["./server"]
