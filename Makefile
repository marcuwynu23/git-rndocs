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
	rm -rf $(BUILD_DIR) dist coverage.out coverage.html
	@echo "Cleaned"

# ------------------------
# Release / Installer targets
# ------------------------

BINARY_NAME ?= git-rndocs
NSIS_VERSION ?= 0.0.0
DEB_VERSION ?= 0.0.0

dist/$(BINARY_NAME)-windows-amd64.exe:
	GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) -o $@ .

dist/$(BINARY_NAME)-linux-amd64:
	GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -o $@ .

dist/$(BINARY_NAME)-darwin-amd64:
	GOOS=darwin GOARCH=amd64 $(GO) build $(GOFLAGS) -o $@ .

dist/$(BINARY_NAME)-darwin-arm64:
	GOOS=darwin GOARCH=arm64 $(GO) build $(GOFLAGS) -o $@ .

dist/$(BINARY_NAME)-setup.exe: dist/$(BINARY_NAME)-windows-amd64.exe
	@mkdir -p dist
	@echo "Creating NSIS installer for version $(NSIS_VERSION)..."
	@makensis -DVERSION=$(NSIS_VERSION) -DBINARY=dist/$(BINARY_NAME)-windows-amd64.exe -DOUTFILE="$@" installer.nsi 2>/dev/null || \
	 echo "NSIS not installed — copying binary as fallback"; cp $< $@

installer-nsis: dist/$(BINARY_NAME)-windows-amd64.exe
	@mkdir -p dist
	@echo "Creating NSIS installer for version $(NSIS_VERSION)..."
	@if command -v makensis >/dev/null 2>&1; then \
		makensis -DVERSION=$(NSIS_VERSION) \
			-DBINARY=dist/$(BINARY_NAME)-windows-amd64.exe \
			-DOUTFILE="dist/$(BINARY_NAME)-setup.exe" \
			installer.nsi; \
	else \
		echo "NSIS not installed — skipping installer"; \
		cp dist/$(BINARY_NAME)-windows-amd64.exe dist/$(BINARY_NAME)-setup.exe; \
	fi

deb: dist/$(BINARY_NAME)-linux-amd64
	@echo "Building .deb package for version $(DEB_VERSION)..."
	@if command -v fpm >/dev/null 2>&1; then \
		fpm -s dir -t deb \
			-n $(BINARY_NAME) \
			-v $(DEB_VERSION) \
			--description "Professional release notes from Git history" \
			--license "Apache-2.0" \
			--url "https://github.com/marcuwynu23/git-rndocs" \
			--maintainer "marcuwynu23" \
			dist/$(BINARY_NAME)-linux-amd64=/usr/local/bin/$(BINARY_NAME); \
	else \
		echo "fpm not installed — skipping .deb"; \
	fi
