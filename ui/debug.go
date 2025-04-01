package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var (
	// Enable or disable debug logging
	DebugEnabled = false
	debugFile    *os.File
	lastFlush    time.Time
)

// InitDebugLog initializes debug logging to a file
func InitDebugLog() {
	if !DebugEnabled {
		return
	}

	// Get config directory
	configDir, err := GetConfigDir()
	if err != nil {
		fmt.Printf("Error getting config directory: %v\n", err)
		return
	}

	// Create debug log file
	logPath := filepath.Join(configDir, "debug.log")
	debugFile, err = os.Create(logPath)
	if err != nil {
		fmt.Printf("Error creating debug log file: %v\n", err)
		return
	}

	lastFlush = time.Now()
	DebugLog("Debug logging initialized")

	// Log some basic system info
	DebugLog("OS: %s, Arch: %s, NumCPU: %d, NumGoroutine: %d",
		runtime.GOOS, runtime.GOARCH, runtime.NumCPU(), runtime.NumGoroutine())
}

// DebugLog logs a message to the debug file
func DebugLog(format string, args ...interface{}) {
	if !DebugEnabled || debugFile == nil {
		return
	}

	// Add timestamp and format the message
	now := time.Now()
	timestamp := now.Format("15:04:05.000")
	message := fmt.Sprintf(format, args...)

	// Include goroutine count periodically
	numGoroutines := runtime.NumGoroutine()

	// Write to the debug file
	fmt.Fprintf(debugFile, "[%s] [G:%d] %s\n", timestamp, numGoroutines, message)

	// Flush the file periodically
	if now.Sub(lastFlush) > time.Second {
		debugFile.Sync()
		lastFlush = now
	}
}

// CloseDebugLog closes the debug log file
func CloseDebugLog() {
	if debugFile != nil {
		DebugLog("Closing debug log")
		debugFile.Close()
		debugFile = nil
	}
}
