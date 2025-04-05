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
}

// Gradient colors for animation
var gradientColors = []string{
	"#00ADD8", // Go blue
	"#15B5DB",
	"#2ABEDE",
	"#3FC6E1",
	"#54CFE4",
	"#69D7E7",
	"#7EE0EA",
	"#93E8ED",
	"#A8F1F0",
	"#BDF9F3",
	"#D2FFF6",
	"#E7FFF9",
	"#FCFFFC",
	"#E7FFF9",
	"#D2FFF6",
	"#BDF9F3",
	"#A8F1F0",
	"#93E8ED",
	"#7EE0EA",
	"#69D7E7",
	"#54CFE4",
	"#3FC6E1",
	"#2ABEDE",
	"#15B5DB",
}

func getGradientIndex(startTime time.Time) int {
	elapsed := time.Since(startTime).Milliseconds()
	return int(elapsed/30) % len(gradientColors)
}

func renderGradientText(text string, startTime time.Time) string {
	var result string
	colorIndex := getGradientIndex(startTime)

	for _, char := range text {
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(gradientColors[colorIndex]))
		result += style.Render(string(char))
		colorIndex = (colorIndex + 1) % len(gradientColors)
	}
	return result
}

func NewWelcomeModel() WelcomeModel {
	return WelcomeModel{
		step:      0,
		done:      false,
		startTime: time.Now(),
	}
}

func (m WelcomeModel) Init() tea.Cmd {
	return tea.Tick(time.Millisecond*16, func(t time.Time) tea.Msg {
		return t
	})
}

func (m WelcomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
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
		return m, tea.Tick(time.Millisecond*16, func(t time.Time) tea.Msg {
			return t
		})

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case time.Time:
		return m, tea.Tick(time.Millisecond*16, func(t time.Time) tea.Msg {
			return t
		})
	}

	return m, nil
}

func (m WelcomeModel) View() string {
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
		title := renderGradientText("Welcome to Go Typer!", m.startTime)
		description := renderGradientText("A modern, feature-rich typing practice tool built with Go.", m.startTime)
		content = titleStyle.Render(title) + "\n\n" +
			textStyle.Render(description) + "\n\n" +
			HintStyle("Press any key to continue...")

	case 1:
		title := renderGradientText("Getting Started", m.startTime)
		features := renderGradientText("• Practice with different text lengths\n"+
			"• Choose between normal and simple modes\n"+
			"• Track your WPM and accuracy\n"+
			"• Customize your experience in settings", m.startTime)
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

	// During development, always show welcome screen
	// TODO: Uncomment this check when welcome screen is finalized
	// if CurrentSettings.HasSeenWelcome {
	// 	return false
	// }

	model := NewWelcomeModel()
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running welcome screen: %v\n", err)
	}

	return true
}
