# Makefile for ssh-thing

BINARY_NAME=ssh-thing
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

.PHONY: all build build-all clean help darwin-arm64 linux-amd64 windows-amd64

# Default target
all: build

# Help command
help:
	@echo "Usage:"
	@echo "  make build         - Build for current platform"
	@echo "  make build-all     - Build for all platforms (macOS ARM, Linux, Windows)"
	@echo "  make darwin-arm64  - Build for macOS ARM64"
	@echo "  make linux-amd64   - Build for Linux AMD64"
	@echo "  make windows-amd64 - Build for Windows AMD64"
	@echo "  make clean         - Remove build artifacts"

# Build for current platform
build:
	go build $(LDFLAGS) -o ./bin/$(BINARY_NAME)

# Build for all target platforms
build-all: darwin-arm64 linux-amd64 windows-amd64

# Build for macOS ARM64 (Apple Silicon)
darwin-arm64:
	@echo "Building for macOS ARM64..."
	@mkdir -p ./bin/darwin-arm64
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o ./bin/darwin-arm64/$(BINARY_NAME)
	@echo "Done! Executable: ./bin/darwin-arm64/$(BINARY_NAME)"

# Build for Linux AMD64
linux-amd64:
	@echo "Building for Linux AMD64..."
	@mkdir -p ./bin/linux-amd64
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o ./bin/linux-amd64/$(BINARY_NAME)
	@echo "Done! Executable: ./bin/linux-amd64/$(BINARY_NAME)"

# Build for Windows AMD64
windows-amd64:
	@echo "Building for Windows AMD64..."
	@mkdir -p ./bin/windows-amd64
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o ./bin/windows-amd64/$(BINARY_NAME).exe
	@echo "Done! Executable: ./bin/windows-amd64/$(BINARY_NAME).exe"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf ./bin
	@echo "Done!"
