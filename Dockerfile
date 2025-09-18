# Multi-stage Dockerfile for OpenSerial
FROM golang:1.24.1-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the applications
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s -X main.version=${VERSION:-dev}" \
    -o openserial \
    ./cmd/openserial

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s -X main.version=${VERSION:-dev}" \
    -o tcpbridge \
    ./cmd/tcpbridge

# Final stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S openserial && \
    adduser -u 1001 -S openserial -G openserial

# Copy binaries from builder stage
COPY --from=builder /app/openserial /usr/local/bin/openserial
COPY --from=builder /app/tcpbridge /usr/local/bin/tcpbridge

# Copy default configuration
COPY --from=builder /app/configs/config.yaml /etc/openserial/config.yaml

# Create directories and set permissions
RUN mkdir -p /var/log/openserial /var/log/tcpbridge && \
    chown -R openserial:openserial /var/log/openserial /var/log/tcpbridge && \
    chmod +x /usr/local/bin/openserial /usr/local/bin/tcpbridge

# Switch to non-root user
USER openserial

# Expose default port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Default command
CMD ["openserial", "-config", "/etc/openserial/config.yaml"]
