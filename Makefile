SHELL := /bin/bash

BINARY = gosh
TEST_TMP_PATH = ./tests/tmp
BINARY_PATH = $(TEST_TMP_PATH)/$(BINARY)
SRC_PATH = ./cmd/gosh/main.go

# 🏗️ Build the release binary
build:
	@echo "🔨 Building $(BINARY)..."
	@mkdir -p bin
	@go build -o bin/$(BINARY) $(SRC_PATH)
	@echo "✅ Build complete! Binary is in ./bin"

# 🧪 Run tests
test: 
	@echo "🔨 Building test $(BINARY)..."
	@mkdir -p bin
	@go build -o $(BINARY_PATH) $(SRC_PATH)
	@echo "✅ Build complete!"

	@echo "🚀 Running tests..."
	@go test ./tests -v

	@echo "🗑️  Cleaning up..."
	@rm -rf $(TEST_TMP_PATH)
	@echo "✅ Cleanup complete!"

# 🔍 Run linter
lint:
	@echo "🔍 Running golangci-lint..."
	@golangci-lint run ./...
