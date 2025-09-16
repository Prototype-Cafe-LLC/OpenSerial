# OpenSerial Makefile

# Variables
BINARY_NAME=openserial
VERSION=0.1.0
BUILD_DIR=build
GO_FILES=$(shell find . -name "*.go" -not -path "./vendor/*")

# Build flags
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/openserial

# Build for multiple platforms
.PHONY: build-all
build-all: build-linux build-windows build-darwin

.PHONY: build-linux
build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)/linux
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/linux/$(BINARY_NAME) ./cmd/openserial
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/linux/$(BINARY_NAME)-arm64 ./cmd/openserial

.PHONY: build-windows
build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)/windows
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/windows/$(BINARY_NAME).exe ./cmd/openserial
	GOOS=windows GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/windows/$(BINARY_NAME)-arm64.exe ./cmd/openserial

.PHONY: build-darwin
build-darwin:
	@echo "Building for macOS..."
	@mkdir -p $(BUILD_DIR)/darwin
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/darwin/$(BINARY_NAME) ./cmd/openserial
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/darwin/$(BINARY_NAME)-arm64 ./cmd/openserial

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run the application
.PHONY: run
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME) -config configs/config.yaml

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
.PHONY: lint
lint:
	@echo "Linting code..."
	golangci-lint run

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Install development tools
.PHONY: install-tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Docker targets
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	docker build -t openserial/openserial:$(VERSION) .

.PHONY: docker-build-multi
docker-build-multi:
	@echo "Building multi-architecture Docker image..."
	docker buildx build --platform linux/amd64,linux/arm64 -t openserial/openserial:$(VERSION) .

.PHONY: docker-run
docker-run: docker-build
	@echo "Running OpenSerial in Docker..."
	docker run --rm -p 8080:8080 --device=/dev/ttyUSB0 openserial/openserial:$(VERSION)

.PHONY: docker-compose-up
docker-compose-up:
	@echo "Starting OpenSerial with Docker Compose..."
	docker-compose up -d

.PHONY: docker-compose-down
docker-compose-down:
	@echo "Stopping OpenSerial with Docker Compose..."
	docker-compose down

# Release targets
.PHONY: release-prepare
release-prepare: clean test build-all
	@echo "Preparing release assets..."
	@mkdir -p dist
	@cp -r $(BUILD_DIR)/* dist/
	@echo "Release assets prepared in dist/"

.PHONY: release-checksums
release-checksums:
	@echo "Generating checksums..."
	@find $(BUILD_DIR) -type f -name "$(BINARY_NAME)*" -exec sha256sum {} \; > $(BUILD_DIR)/checksums.txt
	@echo "Checksums generated: $(BUILD_DIR)/checksums.txt"

# CI/CD targets
.PHONY: ci-test
ci-test: deps fmt lint test-coverage

.PHONY: ci-build
ci-build: deps build-all release-checksums

.PHONY: ci-docker
ci-docker: deps docker-build-multi

# Version management
.PHONY: version
version:
	@echo "Current version: $(VERSION)"

.PHONY: bump-version
bump-version:
	@echo "Bumping version..."
	@new_version=$$(echo $(VERSION) | awk -F. '{$NF = $NF + 1;} 1' | sed 's/ /./g'); \
	sed -i.bak "s/VERSION=$(VERSION)/VERSION=$$new_version/" Makefile; \
	rm Makefile.bak; \
	echo "Version bumped to: $$new_version"

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build              Build the binary"
	@echo "  build-all          Build for all platforms"
	@echo "  test               Run tests"
	@echo "  test-coverage      Run tests with coverage"
	@echo "  run                Build and run the application"
	@echo "  clean              Clean build artifacts"
	@echo "  fmt                Format code"
	@echo "  lint               Lint code"
	@echo "  deps               Install dependencies"
	@echo "  install-tools      Install development tools"
	@echo ""
	@echo "Docker targets:"
	@echo "  docker-build       Build Docker image"
	@echo "  docker-build-multi Build multi-architecture Docker image"
	@echo "  docker-run         Run OpenSerial in Docker"
	@echo "  docker-compose-up  Start with Docker Compose"
	@echo "  docker-compose-down Stop Docker Compose"
	@echo ""
	@echo "Release targets:"
	@echo "  release-prepare    Prepare release assets"
	@echo "  release-checksums  Generate checksums"
	@echo ""
	@echo "CI/CD targets:"
	@echo "  ci-test            Run all tests (deps, fmt, lint, test-coverage)"
	@echo "  ci-build           Build for all platforms with checksums"
	@echo "  ci-docker          Build multi-architecture Docker image"
	@echo ""
	@echo "Version management:"
	@echo "  version            Show current version"
	@echo "  bump-version       Bump patch version"
	@echo ""
	@echo "  help               Show this help"
