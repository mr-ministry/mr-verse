package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/mr-ministry/mr-verse/internal/ui"
)

func init() {
	// Load .env file if it exists
	envFile := ".env"
	if _, err := os.Stat(envFile); err == nil {
		err := godotenv.Load(envFile)
		if err != nil {
			log.Printf("Warning: Error loading .env file: %v", err)
		} else {
			log.Println("Environment variables loaded from .env file")
		}
	}
}

func main() {
	// Set up logging
	// setupLogging()
	
	// Create data directory if it doesn't exist
	ensureDataDirectory()
	
	// Run the application
	// log.Println("Starting Mr Verse application")
	ui.RunApp()
}

// setupLogging configures the application logging
// func setupLogging() {
// // Create logs directory if it doesn't exist
// logsDir := "logs"
// if err := os.MkdirAll(logsDir, 0755); err != nil {
// 		log.Printf("Warning: Could not create logs directory: %v", err)
// 		return
// 	}

// 	// Set up log file
// 	logFile := filepath.Join(logsDir, "mr-verse.log")
// 	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
// 	if err != nil {
// 		log.Printf("Warning: Could not open log file: %v", err)
// 		return
// 	}

// 	// Create a multi-writer to log to both file and stdout
// 	multiWriter := io.MultiWriter(os.Stdout, f)
	
// 	// // Configure the logger to use the multi-writer
// 	log.SetOutput(multiWriter)
// 	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	
// 	log.Println("Logging configured to write to both console and", logFile)
// }

// ensureDataDirectory creates the data directory if it doesn't exist
func ensureDataDirectory() {
	dataDir := "data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("Error: Could not create data directory: %v", err)
	}
	
	// Check if there are any Bible translation files
	files, err := filepath.Glob(filepath.Join(dataDir, "*.json"))
	if err != nil {
		log.Printf("Warning: Could not check for Bible translation files: %v", err)
	} else if len(files) == 0 {
		fmt.Println("Warning: No Bible translation files found in the data directory.")
		fmt.Println("Please add at least one Bible translation JSON file to the data directory.")
	}
}
