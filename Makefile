SHELL := /bin/bash

BINARY = gosh
BINARY_PATH = ./tests/tmp/$(BINARY)
SRC_PATH = ./cmd/gosh/main.go

test: 
	@echo "🔨 Building test $(BINARY)..."
	@mkdir -p bin
	@go build -o $(BINARY_PATH) $(SRC_PATH)
	@echo "✅ Build complete!"

	@echo "🚀 Running tests..."
	@go test ./tests -v

	@echo "🗑️  Cleaning up..."
	@rm -rf $(BINARY_PATH)
	@echo "✅ Cleanup complete!"