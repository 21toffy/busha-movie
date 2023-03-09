package logger

import "log"

func Info(msg string) {
	log.Printf("[INFO] %s\n", msg)
}

// package logger

// import (
// 	"io/ioutil"
// 	"log"
// 	"os"
// 	"path/filepath"

// 	"github.com/spf13/viper"
// )

// // InitLogger initializes the logger for the application
// func InitLogger() {
// 	logFilePath := viper.GetString("log_file_path")

// 	// Create the log directory if it does not exist
// 	logDir := filepath.Dir(logFilePath)
// 	err := os.MkdirAll(logDir, 0755)
// 	if err != nil {
// 		log.Fatalf("Failed to create log directory: %s", err)
// 	}

// 	// Open the log file for writing
// 	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
// 	if err != nil {
// 		log.Fatalf("Failed to open log file: %s", err)
// 	}

// 	// If the environment is not production, suppress logging to stdout
// 	if viper.GetString("environment") != "production" {
// 		log.SetOutput(ioutil.Discard)
// 	} else {
// 		log.SetOutput(logFile)
// 	}

// 	// Set the log prefix
// 	log.SetPrefix("[MYAPP] ")

// 	// Set the log flags
// 	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
// }
