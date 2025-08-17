# Makefile for Advanced Password Manager

# Variables
BINARY_NAME=password-manager
BUILD_DIR=build
MAIN_PATH=cmd/main.go

# Go commands
GO=go
GOBUILD=$(GO) build
GOTEST=$(GO) test
GOCLEAN=$(GO) clean
GOGET=$(GO) get
GOMOD=$(GO) mod

# Build flags
LDFLAGS=-ldflags "-X main.version=$(shell git describe --tags --always --dirty 2>/dev/null || echo 'dev')"

# Default target
.PHONY: all
all: clean build

# Build the application
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Build for multiple platforms
.PHONY: build-all
build-all: clean
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	
	# Windows
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	
	# Linux
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PATH)
	
	# macOS
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	
	@echo "Multi-platform build complete!"

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	$(GOMOD) download
	$(GOGET) golang.org/x/crypto@latest
	$(GOGET) golang.org/x/term@latest
	$(GOGET) github.com/mattn/go-sqlite3@latest
	@echo "Dependencies installed!"

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...
	@echo "Tests complete!"

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests with race detection
.PHONY: test-race
test-race:
	@echo "Running tests with race detection..."
	$(GOTEST) -race -v ./...

# Run benchmarks
.PHONY: benchmark
benchmark:
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...
	@echo "Code formatted!"

# Vet code
.PHONY: vet
vet:
	@echo "Vetting code..."
	$(GO) vet ./...
	@echo "Code vetted!"

# Lint code (requires golangci-lint)
.PHONY: lint
lint:
	@echo "Linting code..."
	golangci-lint run
	@echo "Linting complete!"

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "Clean complete!"

# Install the application
.PHONY: install
install: build
	@echo "Installing $(BINARY_NAME)..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	@echo "Installation complete!"

# Uninstall the application
.PHONY: uninstall
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	@rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "Uninstallation complete!"

# Run the application
.PHONY: run
run: build
	@echo "Running $(BINARY_NAME)..."
	@$(BUILD_DIR)/$(BINARY_NAME) help

# Generate password example
.PHONY: example-gen
example-gen: build
	@echo "Generating example password..."
	@$(BUILD_DIR)/$(BINARY_NAME) generate --length 20 --uppercase --lowercase --numbers --symbols

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all              - Clean and build"
	@echo "  build            - Build the application"
	@echo "  build-all        - Build for multiple platforms"
	@echo "  deps             - Install dependencies"
	@echo "  test             - Run tests"
	@echo "  test-coverage    - Run tests with coverage"
	@echo "  test-race        - Run tests with race detection"
	@echo "  benchmark        - Run benchmarks"
	@echo "  fmt              - Format code"
	@echo "  vet              - Vet code"
	@echo "  lint             - Lint code (requires golangci-lint)"
	@echo "  clean            - Clean build artifacts"
	@echo "  install          - Install the application"
	@echo "  uninstall        - Uninstall the application"
	@echo "  run              - Run the application"
	@echo "  example-gen      - Generate example password"
	@echo "  help             - Show this help message"

# Development helpers
.PHONY: dev-setup
dev-setup: deps fmt vet test
	@echo "Development environment setup complete!"

# Quick development cycle
.PHONY: dev
dev: fmt vet test build
	@echo "Development cycle complete!"

# Show project info
.PHONY: info
info:
	@echo "Advanced Password Manager"
	@echo "========================"
	@echo "Version: $(shell git describe --tags --always --dirty 2>/dev/null || echo 'dev')"
	@echo "Go version: $(shell go version)"
	@echo "Build directory: $(BUILD_DIR)"
	@echo "Main file: $(MAIN_PATH)"
