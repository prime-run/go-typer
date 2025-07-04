package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type GlobalTickMsg time.Time

func InitGlobalTick() tea.Cmd {
	return GlobalTickCmd(GetRefreshInterval())
}

func HandleGlobalTick(lastTick time.Time, msg GlobalTickMsg) (time.Time, bool, tea.Cmd) {
	newTick := time.Time(msg)

	cmd := GlobalTickCmd(GetRefreshInterval())

	return newTick, true, cmd
}

func GlobalTickCmd(interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(t time.Time) tea.Msg {
		return GlobalTickMsg(t)
	})
}

// if set to > 0,  defaults to 10 (100ms).
func GetRefreshInterval() time.Duration {
	fps := 10
	if CurrentSettings.RefreshRate > 0 {
		fps = CurrentSettings.RefreshRate
	}

	return time.Duration(1000/fps) * time.Millisecond
}
