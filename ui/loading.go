package ui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

type tickMsg time.Time

type LoadingModel struct {
	progress progress.Model
	done     bool
	width    int
	height   int
}

func NewLoadingModel() LoadingModel {
	return LoadingModel{
		progress: progress.New(progress.WithDefaultGradient()),
		done:     false,
	}
}

func (m LoadingModel) Init() tea.Cmd {
	return tickCmd()
}

func (m LoadingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.progress.Width = msg.Width - Padding*2 - 4
		if m.progress.Width > MaxWidth {
			m.progress.Width = MaxWidth
		}
		return m, nil

	case tickMsg:
		if m.progress.Percent() == 1.0 {
			return StartTypingGame(m.width, m.height), nil
		}

		cmd := m.progress.IncrPercent(0.25)
		return m, tea.Batch(tickCmd(), cmd)

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	default:
		return m, nil
	}
}

func (m LoadingModel) View() string {
	pad := strings.Repeat(" ", Padding)

	content := "\n" +
		pad + "Loading text..." + "\n\n" +
		pad + m.progress.View() + "\n\n" +
		pad + HelpStyle("Text is being prepared...")

	if m.width > 0 {
		return lipgloss.Place(m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			content)
	}

	return content
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func StartLoading(cmd *cobra.Command, args []string) {
	if _, err := tea.NewProgram(NewLoadingModel(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Oh no!", err)
		os.Exit(1)
	}
}
