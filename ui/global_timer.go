package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// GlobalTickMsg is sent by the global timer to all screens
type GlobalTickMsg time.Time

// GetRefreshInterval returns the current refresh interval based on settings
func GetRefreshInterval() time.Duration {
	// Default to 10 FPS if settings not initialized
	fps := 10
	if CurrentSettings.RefreshRate > 0 {
		fps = CurrentSettings.RefreshRate
	}

	// Calculate milliseconds per frame
	return time.Duration(1000/fps) * time.Millisecond
}

// GlobalTickCmd returns a command that sends a tick every interval
func GlobalTickCmd(interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(t time.Time) tea.Msg {
		return GlobalTickMsg(t)
	})
}

// InitGlobalTick starts the global tick for any model
func InitGlobalTick() tea.Cmd {
	return GlobalTickCmd(GetRefreshInterval())
}

// HandleGlobalTick helps models respond to global ticks
// Returns whether the model should be updated and any commands to execute
func HandleGlobalTick(lastTick time.Time, msg GlobalTickMsg) (time.Time, bool, tea.Cmd) {
	newTick := time.Time(msg)

	// Always schedule the next tick with current refresh rate
	cmd := GlobalTickCmd(GetRefreshInterval())

	// Return the new tick time, whether to update, and the next tick command
	return newTick, true, cmd
}
