# OpenSerial Release Guide

This document explains how to create releases for OpenSerial with automated cross-platform builds.

## Release Process Overview

OpenSerial uses GitHub Actions for automated builds and releases. The process supports:

- **Cross-platform builds**: Linux (amd64/arm64), Windows (amd64/arm64), macOS (amd64/arm64)
- **Docker images**: Multi-architecture Docker images for Linux
- **Automated releases**: GitHub releases with proper asset naming and checksums
- **Security scanning**: Automated vulnerability scanning with Trivy

## Quick Start

### Creating a Release

1. **Update version** in `Makefile`:

   ```bash
   # Edit VERSION=1.0.0 in Makefile
   ```

2. **Create and push a tag**:

   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

3. **GitHub Actions will automatically**:
   - Build for all platforms
   - Create Docker images
   - Generate release assets
   - Create GitHub release

### Using the Release Script

For a more automated approach, use the provided release script:

```bash
# Dry run to see what will happen
./scripts/release.sh --version 1.0.0 --dry-run

# Create actual release
./scripts/release.sh --version 1.0.0

# Skip tests (if you're confident)
./scripts/release.sh --version 1.0.0 --skip-tests
```

## Release Assets

Each release includes:

### Binary Executables

- `openserial-linux-amd64` - Linux x86_64
- `openserial-linux-arm64` - Linux ARM64
- `openserial-windows-amd64.exe` - Windows x86_64
- `openserial-windows-arm64.exe` - Windows ARM64
- `openserial-darwin-amd64` - macOS Intel
- `openserial-darwin-arm64` - macOS Apple Silicon

**Note for macOS users**: After downloading, you'll need to make the binary executable:

```bash
chmod +x openserial-darwin-amd64
# or
chmod +x openserial-darwin-arm64
```

### Checksums

- `checksums.txt` - Combined SHA256 checksums for all binaries
- Individual `.sha256` files for each binary

### Docker Images

- `openserial/openserial:latest` - Latest version
- `openserial/openserial:v1.0.0` - Tagged version
- Multi-architecture support (linux/amd64, linux/arm64)

## GitHub Actions Workflows

### CI Workflow (`.github/workflows/ci.yml`)

Runs on every push and pull request:

- Tests (unit tests, coverage)
- Linting (golangci-lint)
- Build verification
- Docker build test
- Security scanning

### Release Workflow (`.github/workflows/build-and-release.yml`)

Runs on tag creation:

- Cross-platform builds
- Docker image creation and push
- Release asset preparation
- GitHub release creation
- Security scanning

## Docker Support

### Running with Docker

```bash
# Build and run locally
make docker-build
make docker-run

# Or use docker-compose
make docker-compose-up
```

### Docker Compose

The `docker-compose.yml` file provides:

- OpenSerial service with serial device access
- Optional web interface
- Volume mounts for configuration and logs
- Network configuration

### Multi-architecture Docker Images

Docker images are automatically built for:

- `linux/amd64` - Standard x86_64 Linux
- `linux/arm64` - ARM64 Linux (Raspberry Pi, ARM servers)

## Manual Build Process

If you need to build manually:

```bash
# Build for all platforms
make build-all

# Generate checksums
make release-checksums

# Build Docker images
make docker-build-multi

# Prepare release assets
make release-prepare
```

## Configuration

### Required Secrets

For Docker Hub publishing, set these GitHub secrets:

- `DOCKER_USERNAME` - Docker Hub username
- `DOCKER_PASSWORD` - Docker Hub password/token

### Environment Variables

The release process uses these environment variables:

- `GO_VERSION` - Go version (default: 1.24.1)
- `BINARY_NAME` - Binary name (default: openserial)

## Troubleshooting

### Common Issues

1. **Build failures**: Check Go version compatibility
2. **Docker build failures**: Ensure Docker Buildx is available
3. **Release failures**: Verify GitHub token permissions
4. **Checksum mismatches**: Rebuild and regenerate checksums

### Debug Commands

```bash
# Check current version
make version

# Run all CI checks locally
make ci-test
make ci-build
make ci-docker

# Clean and rebuild
make clean
make build-all
```

## Release Checklist

Before creating a release:

- [ ] Update version in `Makefile`
- [ ] Update `CHANGELOG.md` (if applicable)
- [ ] Run tests: `make test`
- [ ] Run linting: `make lint`
- [ ] Test build: `make build-all`
- [ ] Test Docker: `make docker-build`
- [ ] Commit all changes
- [ ] Create and push tag
- [ ] Verify GitHub Actions run successfully
- [ ] Check release assets in GitHub

## Version Management

### Semantic Versioning

OpenSerial follows semantic versioning (MAJOR.MINOR.PATCH):

- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

### Version Bumping

Use the Makefile target to bump versions:

```bash
# Bump patch version (1.0.0 -> 1.0.1)
make bump-version
```

## Security

### Vulnerability Scanning

Every release includes:

- Trivy security scanning
- Dependency vulnerability checks
- Container image scanning
- Results uploaded to GitHub Security tab

### Best Practices

- Keep dependencies updated
- Review security scan results
- Use minimal base images
- Follow security guidelines in Dockerfile

## Support

For release-related issues:

1. Check GitHub Actions logs
2. Review this documentation
3. Check the test plan in `TEST_PLAN.md`
4. Create an issue in the repository
