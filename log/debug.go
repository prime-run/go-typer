package devlog

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/prime-run/go-typer/utils"
)

var (
	DebugEnabled = false   // Enable debug logging
	debugFile    *os.File  // File handle for the debug log
	lastFlush    time.Time // Last time the log was flushed
)

// InitDebugLog initializes the debug log file. It creates a new file in the
// application's config directory and writes initial debug information.
// It also logs the current OS, architecture, number of CPUs, and goroutines.
func InitLog() {
	if !DebugEnabled {
		return
	}

	configDir, err := utils.GetAppConfigDir()
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
	Log("Debug logging initialized")

	Log("OS: %s, Arch: %s, NumCPU: %d, NumGoroutine: %d",
		runtime.GOOS, runtime.GOARCH, runtime.NumCPU(), runtime.NumGoroutine())
}

// DebugLog writes a debug message to the debug log file. It includes a timestamp and the number of goroutines.
// The message is formatted using fmt.Sprintf, and the log is flushed every second.
func Log(format string, args ...interface{}) {
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

// CloseDebugLog closes the debug log file if it is open.
func CloseLog() {
	if debugFile != nil {
		Log("Closing debug log")
		debugFile.Close()
		debugFile = nil
	}
}
