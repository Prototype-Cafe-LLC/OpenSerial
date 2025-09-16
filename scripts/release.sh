#!/bin/bash

# OpenSerial Release Script
# This script automates the release process for OpenSerial

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
BINARY_NAME="openserial"
VERSION=""
DRY_RUN=false
SKIP_TESTS=false
SKIP_DOCKER=false

# Functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

show_help() {
    cat << EOF
OpenSerial Release Script

Usage: $0 [OPTIONS]

Options:
    -v, --version VERSION    Set the version number (required)
    -d, --dry-run           Perform a dry run without creating actual release
    -s, --skip-tests        Skip running tests
    -k, --skip-docker       Skip Docker build
    -h, --help              Show this help message

Examples:
    $0 --version 1.0.0
    $0 --version 1.0.0 --dry-run
    $0 --version 1.0.0 --skip-tests

EOF
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -v|--version)
            VERSION="$2"
            shift 2
            ;;
        -d|--dry-run)
            DRY_RUN=true
            shift
            ;;
        -s|--skip-tests)
            SKIP_TESTS=true
            shift
            ;;
        -k|--skip-docker)
            SKIP_DOCKER=true
            shift
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            log_error "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

# Validate version
if [[ -z "$VERSION" ]]; then
    log_error "Version is required. Use -v or --version option."
    show_help
    exit 1
fi

# Validate version format (semantic versioning)
if [[ ! "$VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    log_error "Invalid version format. Use semantic versioning (e.g., 1.0.0)"
    exit 1
fi

log_info "Starting release process for version $VERSION"

# Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    log_error "Not in a git repository"
    exit 1
fi

# Check if working directory is clean
if [[ -n $(git status --porcelain) ]]; then
    log_error "Working directory is not clean. Please commit or stash changes."
    exit 1
fi

# Check if tag already exists
if git rev-parse "$VERSION" >/dev/null 2>&1; then
    log_error "Tag $VERSION already exists"
    exit 1
fi

# Update version in Makefile
log_info "Updating version in Makefile..."
sed -i.bak "s/VERSION=.*/VERSION=$VERSION/" Makefile
rm Makefile.bak

# Update version in go.mod if needed
if grep -q "version" go.mod; then
    log_info "Updating version in go.mod..."
    sed -i.bak "s/version.*/version $VERSION/" go.mod
    rm go.mod.bak
fi

# Run tests unless skipped
if [[ "$SKIP_TESTS" == "false" ]]; then
    log_info "Running tests..."
    if ! make test; then
        log_error "Tests failed"
        exit 1
    fi
    log_success "Tests passed"
else
    log_warning "Skipping tests"
fi

# Build for all platforms
log_info "Building for all platforms..."
if ! make build-all; then
    log_error "Build failed"
    exit 1
fi
log_success "Build completed"

# Generate checksums
log_info "Generating checksums..."
make release-checksums
log_success "Checksums generated"

# Build Docker image unless skipped
if [[ "$SKIP_DOCKER" == "false" ]]; then
    log_info "Building Docker image..."
    if ! make docker-build-multi; then
        log_error "Docker build failed"
        exit 1
    fi
    log_success "Docker image built"
else
    log_warning "Skipping Docker build"
fi

# Prepare release assets
log_info "Preparing release assets..."
make release-prepare
log_success "Release assets prepared"

# Show release summary
log_info "Release summary:"
echo "  Version: $VERSION"
echo "  Binary name: $BINARY_NAME"
echo "  Platforms: linux/amd64, linux/arm64, windows/amd64, windows/arm64, darwin/amd64, darwin/arm64"
echo "  Build directory: build/"
echo "  Release directory: dist/"

if [[ "$DRY_RUN" == "true" ]]; then
    log_warning "Dry run completed. No git operations performed."
    log_info "To create the actual release, run:"
    echo "  git add ."
    echo "  git commit -m \"Release $VERSION\""
    echo "  git tag -a $VERSION -m \"Release $VERSION\""
    echo "  git push origin main"
    echo "  git push origin $VERSION"
else
    # Commit changes
    log_info "Committing changes..."
    git add .
    git commit -m "Release $VERSION"
    
    # Create tag
    log_info "Creating tag $VERSION..."
    git tag -a "$VERSION" -m "Release $VERSION"
    
    log_success "Release $VERSION created successfully!"
    log_info "To push to remote:"
    echo "  git push origin main"
    echo "  git push origin $VERSION"
fi

log_success "Release process completed!"
