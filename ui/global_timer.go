package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type GlobalTickMsg time.Time

// InitGlobalTick initializes the global tick command with the current refresh interval.
func InitGlobalTick() tea.Cmd {
	return GlobalTickCmd(GetRefreshInterval())
}

// HandleGlobalTick processes the global tick message and returns the new tick time,
func HandleGlobalTick(lastTick time.Time, msg GlobalTickMsg) (time.Time, bool, tea.Cmd) {
	newTick := time.Time(msg)

	cmd := GlobalTickCmd(GetRefreshInterval())

	return newTick, true, cmd
}

// GlobalTickCmd creates a command that ticks at the specified interval.
func GlobalTickCmd(interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(t time.Time) tea.Msg {
		return GlobalTickMsg(t)
	})
}

// GetRefreshInterval calculates the refresh interval based on the current settings.
// If the refresh rate is set to 0, it defaults to 10 FPS (100ms).
func GetRefreshInterval() time.Duration {
	fps := 10
	if CurrentSettings.RefreshRate > 0 {
		fps = CurrentSettings.RefreshRate
	}

	return time.Duration(1000/fps) * time.Millisecond
}
