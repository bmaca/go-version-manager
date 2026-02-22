BINARY=goversion
MODULE=github.com/bmaca/go-toolchain-manager
GO ?= go
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
BUILD_DIR=build
DIST_DIR=dist
PLATFORMS=linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

# Default target
.PHONY: all
all: build

# Build for current platform
.PHONY: build
build:
	@echo "Building $(BINARY) for $(GOOS)/$(GOARCH)..."
	@mkdir -p $(BUILD_DIR)
	@$(GO) build -o $(BUILD_DIR)/$(BINARY)-$(GOOS)-$(GOARCH) ./cmd/
	@chmod +x $(BUILD_DIR)/$(BINARY)-$(GOOS)-$(GOARCH)

.PHONY: static-linux
static-linux:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build \
		-a -ldflags '-extldflags "-static"' \
		-o $(BUILD_DIR)/$(BINARY)-static ./cmd/

# Cross-compile for all platforms
.PHONY: cross
cross:
	@rm -rf $(DIST_DIR)
	@mkdir -p $(DIST_DIR)
	$(foreach platform,$(PLATFORMS), \
		CGO_ENABLED=0 GOOS=$(firstword $(subst /, ,$(platform))) \
		GOARCH=$(lastword $(subst /, ,$(platform))) \
		$(GO) build -o $(DIST_DIR)/$(BINARY)-$(platform) ./cmd/ && \
		echo "$(BINARY)-$(platform) built"; \
	)
	@echo "Built all binaries in $(DIST_DIR)/"

# Clean build artifacts
.PHONY: clean
clean:
	@rm -rf $(BUILD_DIR) $(DIST_DIR)
	@echo "Cleaned build artifacts"

# Run tests
.PHONY: test
test:
	@$(GO) test -v ./...

# Format code
.PHONY: fmt
fmt:
	@$(GO) fmt ./...

# Lint (requires golangci-lint)
.PHONY: lint
lint:
	@golangci-lint run

# Install dependencies
.PHONY: tidy
tidy:
	@$(GO) mod tidy

# Generate help
.PHONY: help
help: ## Show this help
	@egrep -h '\s##\s' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Show build info
.PHONY: info
info: ## Show Go environment info
	@$(GO) version
	@echo "GOOS: $(GOOS)"
	@echo "GOARCH: $(GOARCH)"
	@echo "Module: $(MODULE)"
