.PHONY: all clean test test-go test-integration build build-windows build-windows-arm64 build-stubs build-testapp smoke-wine run-wine fmt lint vet check dev-setup help

# Build configuration
BIN_DIR=bin
BIN=$(BIN_DIR)/gjg-launcher.exe
STUBS_DIR=testdata/jre/bin
JAVA_TEST_DIR=testdata/java

# Version information
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

# Default target
all: test build

# Clean all build artifacts
clean:
	rm -rf $(BIN_DIR)/ $(STUBS_DIR)/
	$(MAKE) -C $(JAVA_TEST_DIR) clean
	go clean -testcache

# Development setup - install dependencies and verify environment
dev-setup:
	@echo "Setting up development environment..."
	@echo "Checking Go version..."
	@go version
	@echo "Checking Java version..."
	@java -version 2>&1 | head -n 1
	@javac -version 2>&1
	@if command -v wine64 >/dev/null 2>&1; then \
		echo "Wine version:"; \
		wine64 --version; \
	else \
		echo "Warning: Wine not found. Wine-based testing will not work."; \
	fi
	@echo "Development environment ready!"

# Test targets
test: test-go test-integration

test-go:
	@echo "Running Go unit tests..."
	go test -v ./internal/...

test-integration: build-stubs build-testapp smoke-wine

# Build targets
build: build-windows

build-windows:
	@echo "Building launcher for Windows amd64..."
	@mkdir -p $(BIN_DIR)
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BIN) ./cmd/launcher

build-windows-arm64:
	@echo "Building launcher for Windows arm64..."
	@mkdir -p $(BIN_DIR)
	GOOS=windows GOARCH=arm64 go build $(LDFLAGS) -o $(BIN_DIR)/gjg-launcher-arm64.exe ./cmd/launcher

# Build test stubs (Windows Java executables for testing)
build-stubs:
	@echo "Building test stubs..."
	@mkdir -p $(STUBS_DIR)
	GOOS=windows GOARCH=amd64 go build -o $(STUBS_DIR)/java.exe ./testdata/stubs
	cp $(STUBS_DIR)/java.exe $(STUBS_DIR)/javaw.exe

# Build test Java application
build-testapp:
	@echo "Building test Java application..."
	$(MAKE) -C $(JAVA_TEST_DIR) jar

# Wine-based testing
smoke-wine: build-windows build-stubs
	@echo "Running Wine smoke tests..."
	bash scripts/wine-smoke.sh

run-wine: build-windows
	bash scripts/wine-run.sh $(BIN) --gjg-dry-run

# Code quality
fmt:
	@echo "Formatting Go code..."
	go fmt ./...

lint:
	@echo "Running Go lint..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found, running go vet instead"; \
		$(MAKE) vet; \
	fi

vet:
	@echo "Running go vet..."
	go vet ./...

# Combined check target
check: fmt vet test

# CI target - comprehensive testing for continuous integration
ci: dev-setup check test-integration
	@echo "All CI checks passed!"

# Help target
help:
	@echo "Available targets:"
	@echo "  all              - Run tests and build (default)"
	@echo "  clean            - Clean all build artifacts"
	@echo "  dev-setup        - Set up development environment"
	@echo "  test             - Run all tests (Go unit + integration)"
	@echo "  test-go          - Run Go unit tests only"
	@echo "  test-integration - Run integration tests (includes smoke-wine)"
	@echo "  build            - Build Windows launcher (amd64)"
	@echo "  build-windows    - Build Windows launcher (amd64)"
	@echo "  build-windows-arm64 - Build Windows launcher (arm64)"
	@echo "  build-stubs      - Build test stub executables"
	@echo "  build-testapp    - Build test Java application"
	@echo "  smoke-wine       - Run Wine-based smoke tests"
	@echo "  run-wine         - Quick Wine test run"
	@echo "  fmt              - Format Go code"
	@echo "  lint             - Run linter (golangci-lint or go vet)"
	@echo "  vet              - Run go vet"
	@echo "  check            - Run fmt, vet, and tests"
	@echo "  ci               - Full CI pipeline (dev-setup, check, integration tests)"
	@echo "  help             - Show this help message"

