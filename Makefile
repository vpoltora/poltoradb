# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Binary names
PARSER_BINARY=parser
PARSER_SOURCE=./cmd/poltoradb

# Build directory
BUILD_DIR=bin

.PHONY: all build parser clean test deps help

# Default target
all: build

# Build all binaries
build: parser

# Build parser binary
parser:
	@echo "Building parser..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(PARSER_BINARY) $(PARSER_SOURCE)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run parser
run-parser: parser
	@echo "Running parser..."
	./$(BUILD_DIR)/$(PARSER_BINARY)

# Development build with debug info
dev-build:
	@echo "Building with debug info..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -gcflags="all=-N -l" -o $(BUILD_DIR)/$(PARSER_BINARY) $(PARSER_SOURCE)

# Help target
help:
	@echo "Available targets:"
	@echo "  all          - Build all binaries (default)"
	@echo "  build        - Build all binaries"
	@echo "  parser       - Build parser binary"
	@echo "  clean        - Clean build artifacts"
	@echo "  test         - Run tests"
	@echo "  run-parser   - Build and run parser"
	@echo "  dev-build    - Build with debug information"
	@echo "  help         - Show this help message"
