.PHONY: build run clean test

# Application name
APP_NAME := mr-verse

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GORUN := $(GOCMD) run
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod

# Main package path
MAIN_PATH := ./cmd/main.go
BIN_PATH := ./bin

# Build the application
build:
	$(GOBUILD) -o $(APP_NAME) $(MAIN_PATH)

# Run the application
run:
	$(GORUN) $(MAIN_PATH)

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(APP_NAME)

# Run tests
test:
	$(GOTEST) -v ./...

# Update dependencies
deps:
	$(GOMOD) tidy

# Build for Windows
build-windows:
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(APP_NAME).exe $(MAIN_PATH)

# Build for Windows GUI (no terminal)
build-windows-gui:
	$(GOBUILD) -o $(BIN_PATH)/$(APP_NAME).exe -ldflags -H=windowsgui $(MAIN_PATH)

# Build for macOS
build-macos:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(APP_NAME)-macos $(MAIN_PATH)

# Build for Linux
build-linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(APP_NAME)-linux $(MAIN_PATH)

# Build for all platforms
build-all: build-windows-gui build-macos build-linux

# Default target
all: build 
