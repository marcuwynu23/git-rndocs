.PHONY: all build clean test lint cover install help link unlink installer-nsis deb

# OS detection
ifeq ($(OS),Windows_NT)
  BINARY_EXT := .exe
  IS_WINDOWS := 1
else
  BINARY_EXT :=
  IS_WINDOWS := 0
endif

BINARY_BASE=git-rndocs
BINARY_NAME=$(BINARY_BASE)$(BINARY_EXT)
BUILD_DIR=./build
GO=go
GOFLAGS=-ldflags="-s -w"
GOTEST=$(GO) test

help:
	@echo "Usage:"
	@echo "  make build       - Build the binary"
	@echo "  make install     - Install the binary (go install)"
	@echo "  make link        - Symlink binary to C:/Bin/tools/"
	@echo "  make unlink      - Remove symlink from C:/Bin/tools/"
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
	$(GOTEST) -v ./...

lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed, skipping"; \
	fi

cover:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	$(GO) tool cover -func=coverage.out
	@echo "Coverage report: coverage.html"

clean:
	rm -rf $(BUILD_DIR) dist coverage.out coverage.html
	@echo "Cleaned"

# ------------------------
# Symlink (Windows: C:\Bin\tools)
# ------------------------

LINK_DIR ?= C:/Bin/tools
LINK_TARGET := $(LINK_DIR)/$(BINARY_NAME)

link: build
ifeq ($(IS_WINDOWS),1)
	@powershell -Command "if (-not (Test-Path '$(LINK_DIR)')) { New-Item -ItemType Directory -Path '$(LINK_DIR)' -Force | Out-Null }"
	@echo Creating symlink at $(LINK_TARGET)
	@powershell -Command "New-Item -ItemType SymbolicLink -Path '$(LINK_TARGET)' -Target '$(abspath $(BUILD_DIR)/$(BINARY_NAME))' -Force" 2>&1
else
	@mkdir -p $(LINK_DIR)
	@echo Creating symlink at $(LINK_TARGET)
	@ln -sf $(PWD)/$(BUILD_DIR)/$(BINARY_NAME) $(LINK_TARGET)
endif

unlink:
ifeq ($(IS_WINDOWS),1)
	@powershell -Command "if (Test-Path '$(LINK_TARGET)') { Remove-Item -Path '$(LINK_TARGET)' -Force; Write-Host 'Removed symlink at $(LINK_TARGET)' } else { Write-Host 'Symlink not found at $(LINK_TARGET)' }"
else
	@if [ -f $(LINK_TARGET) ]; then \
		echo Removing symlink $(LINK_TARGET); \
		rm $(LINK_TARGET); \
	else \
		echo Symlink not found at $(LINK_TARGET); \
	fi
endif

# ------------------------
# Release / Installer targets
# ------------------------

NSIS_VERSION ?= 0.0.0
DEB_VERSION ?= 0.0.0

dist/$(BINARY_BASE)-windows-amd64.exe:
	GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) -o $@ .

dist/$(BINARY_BASE)-linux-amd64:
	GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -o $@ .

dist/$(BINARY_BASE)-darwin-amd64:
	GOOS=darwin GOARCH=amd64 $(GO) build $(GOFLAGS) -o $@ .

dist/$(BINARY_BASE)-darwin-arm64:
	GOOS=darwin GOARCH=arm64 $(GO) build $(GOFLAGS) -o $@ .

dist/$(BINARY_BASE)-setup.exe: dist/$(BINARY_BASE)-windows-amd64.exe
	@mkdir -p dist
	@echo "Creating NSIS installer for version $(NSIS_VERSION)..."
	@makensis -DVERSION=$(NSIS_VERSION) -DBINARY=dist/$(BINARY_BASE)-windows-amd64.exe -DOUTFILE="$@" installers/installer.nsi 2>/dev/null || \
	 echo "NSIS not installed — copying binary as fallback"; cp $< $@

installer-nsis: dist/$(BINARY_BASE)-windows-amd64.exe
	@mkdir -p dist
	@echo "Creating NSIS installer for version $(NSIS_VERSION)..."
	@if command -v makensis >/dev/null 2>&1; then \
		makensis -DVERSION=$(NSIS_VERSION) \
			-DBINARY=dist/$(BINARY_BASE)-windows-amd64.exe \
			-DOUTFILE="dist/$(BINARY_BASE)-setup.exe" \
			installers/installer.nsi; \
	else \
		echo "NSIS not installed — skipping installer"; \
		cp dist/$(BINARY_BASE)-windows-amd64.exe dist/$(BINARY_BASE)-setup.exe; \
	fi

deb: dist/$(BINARY_BASE)-linux-amd64
	@echo "Building .deb package for version $(DEB_VERSION)..."
	@if command -v fpm >/dev/null 2>&1; then \
		fpm -s dir -t deb \
			-n $(BINARY_BASE) \
			-v $(DEB_VERSION) \
			--description "Professional release notes from Git history" \
			--license "Apache-2.0" \
			--url "https://github.com/marcuwynu23/git-rndocs" \
			--maintainer "marcuwynu23" \
			dist/$(BINARY_BASE)-linux-amd64=/usr/local/bin/$(BINARY_BASE); \
	else \
		echo "fpm not installed — skipping .deb"; \
	fi
