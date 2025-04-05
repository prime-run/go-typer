package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

type GlobalTickMsg time.Time

func GetRefreshInterval() time.Duration {
	fps := 10
	if CurrentSettings.RefreshRate > 0 {
		fps = CurrentSettings.RefreshRate
	}

	return time.Duration(1000/fps) * time.Millisecond
}

func GlobalTickCmd(interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(t time.Time) tea.Msg {
		return GlobalTickMsg(t)
	})
}

func InitGlobalTick() tea.Cmd {
	return GlobalTickCmd(GetRefreshInterval())
}

func HandleGlobalTick(lastTick time.Time, msg GlobalTickMsg) (time.Time, bool, tea.Cmd) {
	newTick := time.Time(msg)

	cmd := GlobalTickCmd(GetRefreshInterval())

	return newTick, true, cmd
}
