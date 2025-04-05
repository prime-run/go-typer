package ui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type WelcomeModel struct {
	width     int
	height    int
	step      int
	done      bool
	startTime time.Time
	lastTick  time.Time
}

// Gradient is now defined in the shared gradient.go file

func NewWelcomeModel() *WelcomeModel {
	return &WelcomeModel{
		step:      0,
		done:      false,
		startTime: time.Now(),
		lastTick:  time.Now(),
	}
}

func (m *WelcomeModel) Init() tea.Cmd {
	return InitGlobalTick()
}

func (m *WelcomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case GlobalTickMsg:
		// Handle the global tick
		var cmd tea.Cmd
		m.lastTick, _, cmd = HandleGlobalTick(m.lastTick, msg)
		return m, cmd

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC || msg.Type == tea.KeyEsc {
			return m, tea.Quit
		}
		// Any key press advances to next step
		m.step++
		if m.step >= 2 {
			m.done = true
			CurrentSettings.HasSeenWelcome = true
			SaveSettings()
			return m, tea.Quit
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	return m, nil
}

func (m *WelcomeModel) View() string {
	if m.done {
		return ""
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Margin(1, 0)

	textStyle := lipgloss.NewStyle().
		Width(60).
		Align(lipgloss.Center)

	var content string
	switch m.step {
	case 0:
		title := RenderGradientText("Welcome to Go Typer!", m.lastTick)
		description := RenderGradientText("A modern, feature-rich typing practice tool built with Go.", m.lastTick)
		content = titleStyle.Render(title) + "\n\n" +
			textStyle.Render(description) + "\n\n" +
			HintStyle("Press any key to continue...")

	case 1:
		title := RenderGradientText("Getting Started", m.lastTick)
		features := RenderGradientText("• Practice with different text lengths\n"+
			"• Choose between normal and simple modes\n"+
			"• Track your WPM and accuracy\n"+
			"• Customize your experience in settings", m.lastTick)
		content = titleStyle.Render(title) + "\n\n" +
			textStyle.Render(features) + "\n\n" +
			HintStyle("Press any key to start typing...")
	}

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		content)
}

func ShowWelcomeScreen() bool {
	// Initialize settings first
	InitSettings()

	// Only show welcome screen if user hasn't seen it before
	if CurrentSettings.HasSeenWelcome {
		return false
	}

	model := NewWelcomeModel()
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running welcome screen: %v\n", err)
	}

	return true
}
