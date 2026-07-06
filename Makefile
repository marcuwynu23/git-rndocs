.PHONY: all build clean test lint cover install help

BINARY_NAME=git-rndocs
BUILD_DIR=./build
GO=go
GOFLAGS=-ldflags="-s -w"
GOTEST=$(GO) test

help:
	@echo "Usage:"
	@echo "  make build       - Build the binary"
	@echo "  make install     - Install the binary (go install)"
	@echo "  make test        - Run all tests"
	@echo "  make lint        - Run golangci-lint"
	@echo "  make cover       - Run tests with coverage"
	@echo "  make clean       - Remove build artifacts"
	@echo "  make all         - Build + test + lint"

all: clean lint test build

build:
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "Built $(BUILD_DIR)/$(BINARY_NAME)"

install:
	$(GO) install $(GOFLAGS) .
	@echo "Installed $(BINARY_NAME)"

test:
	$(GOTEST) -v -race ./...

lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed, skipping"; \
	fi

cover:
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	$(GO) tool cover -func=coverage.out
	@echo "Coverage report: coverage.html"

clean:
	rm -rf $(BUILD_DIR) coverage.out coverage.html
	@echo "Cleaned"
