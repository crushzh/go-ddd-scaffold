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

# ==================== Frontend ====================

.PHONY: web
web:
	@echo "Building frontend..."
	@cd web && npm install && npm run build
	@echo "Frontend build complete"

# ==================== Development ====================

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

# ==================== Packaging ====================

# Generate .run self-extracting installer
define make_package
	@echo "Packaging $(APP_NAME)-$(1)..."
	@mkdir -p $(BUILD_DIR)/pkg-$(1)
	@cp $(BUILD_DIR)/linux-$(1)/$(APP_NAME) $(BUILD_DIR)/pkg-$(1)/
	@cp -r configs $(BUILD_DIR)/pkg-$(1)/
	@cp scripts/install.sh $(BUILD_DIR)/pkg-$(1)/
	@cp scripts/uninstall.sh $(BUILD_DIR)/pkg-$(1)/
	@cd $(BUILD_DIR) && tar czf pkg-$(1).tar.gz -C pkg-$(1) .
	@cat scripts/makeself-header.sh $(BUILD_DIR)/pkg-$(1).tar.gz > $(BUILD_DIR)/$(APP_NAME)-$(1)-$(VERSION).run
	@chmod +x $(BUILD_DIR)/$(APP_NAME)-$(1)-$(VERSION).run
	@rm -rf $(BUILD_DIR)/pkg-$(1) $(BUILD_DIR)/pkg-$(1).tar.gz
	@echo "Package: $(BUILD_DIR)/$(APP_NAME)-$(1)-$(VERSION).run"
endef

.PHONY: package-linux
package-linux: build-linux
	$(call make_package,amd64)

.PHONY: package-arm64
package-arm64: build-arm64
	$(call make_package,arm64)

.PHONY: package-arm32
package-arm32: build-arm32
	$(call make_package,arm32)

.PHONY: package-all
package-all: package-linux package-arm64 package-arm32
	@echo "All packages complete!"

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
	@echo "Development:"
	@echo "  run             Run in development mode"
	@echo "  build           Build for local platform"
	@echo "  web             Build frontend"
	@echo ""
	@echo "Cross Compile:"
	@echo "  build-linux     Build Linux amd64"
	@echo "  build-arm64     Build Linux arm64"
	@echo "  build-arm32     Build Linux arm32 (ARMv7)"
	@echo "  build-windows   Build Windows"
	@echo "  build-all       Build all platforms"
	@echo ""
	@echo "Packaging:"
	@echo "  package-linux   Package .run installer (amd64)"
	@echo "  package-arm64   Package .run installer (arm64)"
	@echo "  package-arm32   Package .run installer (arm32)"
	@echo "  package-all     Package all platforms"
	@echo ""
	@echo "Code Generator:"
	@echo "  gen             Generate DDD module (make gen name=order cn=Order)"
	@echo ""
	@echo "Documentation:"
	@echo "  docs            Generate Swagger docs"
	@echo ""
	@echo "Quality:"
	@echo "  lint            Run linter"
	@echo "  test            Run tests"
	@echo "  coverage        Test coverage"
	@echo "  fmt             Format code"
	@echo ""
	@echo "Misc:"
	@echo "  deps            Download dependencies"
	@echo "  clean           Clean build artifacts"
	@echo "  info            Show build configuration"
	@echo "  help            Show this help"
