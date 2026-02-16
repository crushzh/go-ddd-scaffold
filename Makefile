# Go Scaffold - DDD Template Makefile
# Cross-platform build, code generator, Swagger docs

# ==================== Variables ====================
APP_NAME := myapp
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "0.1.0")
BUILD_TIME := $(shell date +%Y-%m-%dT%H:%M:%S)
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

GO := go
GOFLAGS := -trimpath
LDFLAGS := -s -w \
	-X 'main.Version=$(VERSION)' \
	-X 'main.BuildTime=$(BUILD_TIME)' \
	-X 'main.GitCommit=$(GIT_COMMIT)'

BUILD_DIR := build

# ==================== Cross-compilation toolchains ====================
ifeq ($(OS),Windows_NT)
    DETECTED_OS := Windows
    CC_LINUX_AMD64 ?= x86_64-linux-gnu-gcc
    CC_LINUX_ARM64 ?= aarch64-linux-gnu-gcc
    CC_LINUX_ARM32 ?= arm-linux-gnueabihf-gcc
    CC_WINDOWS ?= gcc
else
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Darwin)
        DETECTED_OS := macOS
        CC_LINUX_AMD64 ?= x86_64-unknown-linux-gnu-gcc
        CC_LINUX_ARM64 ?= aarch64-unknown-linux-gnu-gcc
        CC_LINUX_ARM32 ?= arm-unknown-linux-gnueabihf-gcc
        CC_WINDOWS ?= x86_64-w64-mingw32-gcc
    else
        DETECTED_OS := Linux
        CC_LINUX_AMD64 ?= gcc
        CC_LINUX_ARM64 ?= aarch64-linux-gnu-gcc
        CC_LINUX_ARM32 ?= arm-linux-gnueabihf-gcc
        CC_WINDOWS ?= x86_64-w64-mingw32-gcc
    endif
endif

.PHONY: all
all: build

.PHONY: run
run:
	$(GO) run ./cmd/server/ -c configs/config.yaml

.PHONY: build
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME) ./cmd/server/
	@echo "Build complete: $(BUILD_DIR)/$(APP_NAME)"

.PHONY: build-linux
build-linux:
	@mkdir -p $(BUILD_DIR)/linux-amd64
	CGO_ENABLED=1 CC=$(CC_LINUX_AMD64) GOOS=linux GOARCH=amd64 \
		$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" \
		-o $(BUILD_DIR)/linux-amd64/$(APP_NAME) ./cmd/server/
	@cp -r configs $(BUILD_DIR)/linux-amd64/

.PHONY: build-arm64
build-arm64:
	@mkdir -p $(BUILD_DIR)/linux-arm64
	CGO_ENABLED=1 CC=$(CC_LINUX_ARM64) GOOS=linux GOARCH=arm64 \
		$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" \
		-o $(BUILD_DIR)/linux-arm64/$(APP_NAME) ./cmd/server/
	@cp -r configs $(BUILD_DIR)/linux-arm64/

.PHONY: build-arm32
build-arm32:
	@mkdir -p $(BUILD_DIR)/linux-arm32
	CGO_ENABLED=1 CC=$(CC_LINUX_ARM32) GOOS=linux GOARCH=arm GOARM=7 \
		$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" \
		-o $(BUILD_DIR)/linux-arm32/$(APP_NAME) ./cmd/server/
	@cp -r configs $(BUILD_DIR)/linux-arm32/

.PHONY: build-windows
build-windows:
	@mkdir -p $(BUILD_DIR)/windows-amd64
	CGO_ENABLED=1 CC=$(CC_WINDOWS) GOOS=windows GOARCH=amd64 \
		$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" \
		-o $(BUILD_DIR)/windows-amd64/$(APP_NAME).exe ./cmd/server/
	@cp -r configs $(BUILD_DIR)/windows-amd64/

.PHONY: build-all
build-all: build-linux build-arm64 build-arm32 build-windows
	@echo "All platforms build complete!"

# ==================== Code Generator ====================

.PHONY: gen
gen:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make gen name=order cn=Order"; \
		exit 1; \
	fi
	$(GO) run ./cmd/gen/ -name $(name) -cn "$(cn)"

# ==================== Documentation ====================

.PHONY: swag-install
swag-install:
	go install github.com/swaggo/swag/cmd/swag@latest

.PHONY: docs
docs:
	swag init -d ./cmd/server,./internal -g main.go -o docs/swagger --parseDependency --parseInternal

# ==================== Code Quality ====================

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: test
test:
	$(GO) test -v -race -cover ./...

.PHONY: coverage
coverage:
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

.PHONY: fmt
fmt:
	$(GO) fmt ./...

# ==================== Tools ====================

.PHONY: deps
deps:
	$(GO) mod download
	$(GO) mod tidy

.PHONY: clean
clean:
	@rm -rf $(BUILD_DIR) coverage.out coverage.html

.PHONY: info
info:
	@echo "App: $(APP_NAME) | Version: $(VERSION) | OS: $(DETECTED_OS)"

.PHONY: help
help:
	@echo "$(APP_NAME) Build System (DDD)"
	@echo ""
	@echo "  run             Run in development mode"
	@echo "  build           Build for local platform"
	@echo "  build-all       Build for all platforms"
	@echo "  gen             Generate DDD module (make gen name=order cn=Order)"
	@echo "  docs            Generate Swagger docs"
	@echo "  test            Run tests"
	@echo "  lint            Run linter"
	@echo "  deps            Download dependencies"
	@echo "  clean           Clean build artifacts"
	@echo "  help            Show this help"
