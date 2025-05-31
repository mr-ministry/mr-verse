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

LINUX_PATH := linux
WINDOWS_GUI_PATH := windows
WINDOWS_DEBUG_PATH := windows-debug
MACOS_PATH := macos

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

# GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BIN_PATH)/$(WINDOWS_DEBUG_PATH)/$(APP_NAME).exe -gcflags "-g -dwarf" $(MAIN_PATH)
# $(GOBUILD) -o $(BIN_PATH)/$(WINDOWS_GUI_PATH)/$(APP_NAME).exe -ldflags -H=windowsgui $(MAIN_PATH)

# Build for Windows
build-windows:
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BIN_PATH)/$(WINDOWS_DEBUG_PATH)/$(APP_NAME).exe -gcflags="all=-N -l" $(MAIN_PATH)

# Build for Windows GUI (no terminal)
build-windows-gui:
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BIN_PATH)/$(WINDOWS_GUI_PATH)/$(APP_NAME).exe -ldflags="-H=windowsgui -s -w" $(MAIN_PATH)

# Build for macOS
build-macos:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BIN_PATH)/$(MACOS_PATH)/$(APP_NAME)-macos $(MAIN_PATH)

# Build for Linux
build-linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BIN_PATH)/$(LINUX_PATH)/$(APP_NAME)-linux $(MAIN_PATH)

# Clean build artifacts
clean-all:
	rm -f $(BIN_PATH)/$(BIN_PATH)/$(WINDOWS_GUI_PATH)/$(APP_NAME).exe
	rm -f $(BIN_PATH)/$(BIN_PATH)/$(WINDOWS_DEBUG_PATH)/$(APP_NAME).exe
	rm -f $(BIN_PATH)/$(MACOS_PATH)/$(APP_NAME)-macos
	rm -f $(BIN_PATH)/$(LINUX_PATH)/$(APP_NAME)-linux
	#
# Build for all platforms
build-all: clean-all build-windows build-windows-gui build-macos build-linux

# Linux compaible build
build-all-linux: clean-all build-linux

# Default target
all: build 
