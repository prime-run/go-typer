package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var (
	DebugEnabled = false
	debugFile    *os.File
	lastFlush    time.Time
)

func InitDebugLog() {
	if !DebugEnabled {
		return
	}

	configDir, err := GetConfigDir()
	if err != nil {
		fmt.Printf("Error getting config directory: %v\n", err)
		return
	}

	logPath := filepath.Join(configDir, "debug.log")
	debugFile, err = os.Create(logPath)
	if err != nil {
		fmt.Printf("Error creating debug log file: %v\n", err)
		return
	}

	lastFlush = time.Now()
	DebugLog("Debug logging initialized")

	DebugLog("OS: %s, Arch: %s, NumCPU: %d, NumGoroutine: %d",
		runtime.GOOS, runtime.GOARCH, runtime.NumCPU(), runtime.NumGoroutine())
}

func DebugLog(format string, args ...interface{}) {
	if !DebugEnabled || debugFile == nil {
		return
	}

	now := time.Now()
	timestamp := now.Format("15:04:05.000")
	message := fmt.Sprintf(format, args...)

	numGoroutines := runtime.NumGoroutine()

	fmt.Fprintf(debugFile, "[%s] [G:%d] %s\n", timestamp, numGoroutines, message)

	if now.Sub(lastFlush) > time.Second {
		debugFile.Sync()
		lastFlush = now
	}
}

func CloseDebugLog() {
	if debugFile != nil {
		DebugLog("Closing debug log")
		debugFile.Close()
		debugFile = nil
	}
}
