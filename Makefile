# SubdomainX Makefile

.PHONY: build test clean install run help

# Build variables
BINARY_NAME=subdomainx
BUILD_DIR=build
MAIN_PATH=./cmd/subdomainx

# Default target
all: build

# Build the application
build:
	@echo "Building SubdomainX..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Build for multiple platforms
build-all: clean
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	
	# Linux
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	GOOS=linux GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PATH)
	
	# macOS
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	
	# Windows
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	
	@echo "Multi-platform build complete!"

# Install the application globally
install:
	@echo "Installing SubdomainX..."
	go install $(MAIN_PATH)
	@echo "Installation complete!"

# Test go install from GitHub (requires pushing to GitHub first)
install-remote:
	@echo "Installing SubdomainX from GitHub..."
	go install github.com/itszeeshan/subdomainx/cmd/subdomainx@latest
	@echo "SubdomainX installed from GitHub!"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...
	@echo "Tests complete!"

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "Clean complete!"

# Run the application
run: build
	@echo "Running SubdomainX..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# Run with custom domains file
run-example: build
	@echo "Running SubdomainX with example domains..."
	@echo "example.com" > domains.txt
	@echo "test.org" >> domains.txt
	./$(BUILD_DIR)/$(BINARY_NAME) --wildcard domains.txt
	@echo "Check output/ directory for results"

# Check tool availability
check-tools: build
	@echo "Checking tool availability..."
	./$(BUILD_DIR)/$(BINARY_NAME) --check-tools

# Show installation instructions
install-tools: build
	@echo "Showing installation instructions..."
	./$(BUILD_DIR)/$(BINARY_NAME) --install-tools

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "Code formatting complete!"

# Lint code
lint:
	@echo "Linting code..."
	golangci-lint run
	@echo "Linting complete!"

# Update dependencies
deps:
	@echo "Updating dependencies..."
	go mod tidy
	go mod download
	@echo "Dependencies updated!"

# Create release
release: clean build-all
	@echo "Creating release..."
	@cd $(BUILD_DIR) && tar -czf ../subdomainx-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64
	@cd $(BUILD_DIR) && tar -czf ../subdomainx-linux-arm64.tar.gz $(BINARY_NAME)-linux-arm64
	@cd $(BUILD_DIR) && tar -czf ../subdomainx-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64
	@cd $(BUILD_DIR) && tar -czf ../subdomainx-darwin-arm64.tar.gz $(BINARY_NAME)-darwin-arm64
	@cd $(BUILD_DIR) && zip ../subdomainx-windows-amd64.zip $(BINARY_NAME)-windows-amd64.exe
	@echo "Release packages created!"

# Show help
help:
	@echo "SubdomainX Makefile Commands:"
	@echo ""
	@echo "  build          - Build the application"
	@echo "  build-all      - Build for multiple platforms (Linux, macOS, Windows)"
	@echo "  install        - Install the application globally"
	@echo "  install-remote - Install from GitHub (go install)"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  clean          - Clean build artifacts"
	@echo "  run            - Build and run the application"
	@echo "  run-example    - Run with example domains"
	@echo "  check-tools    - Check tool availability"
	@echo "  install-tools  - Show installation instructions"
	@echo "  fmt            - Format code"
	@echo "  lint           - Lint code (requires golangci-lint)"
	@echo "  deps           - Update dependencies"
	@echo "  release        - Create release packages"
	@echo "  help           - Show this help message"
