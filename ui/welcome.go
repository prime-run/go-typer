package ui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"time"
)

type WelcomeModel struct {
	width     int
	height    int
	step      int
	done      bool
	startTime time.Time
	lastTick  time.Time
}

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
		var cmd tea.Cmd
		m.lastTick, _, cmd = HandleGlobalTick(m.lastTick, msg)
		return m, cmd

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC || msg.Type == tea.KeyEsc {
			return m, tea.Quit
		}
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
		Width(80).
		Align(lipgloss.Center)

	var content string
	switch m.step {
	case 0:
		title := RenderGradientText("Welcome to Go Typer!", m.lastTick)
		description := RenderGradientText("\nYou only see these messages once, so first, thank you for playing :)", m.lastTick)
		content = titleStyle.Render(title) + "\n\n" +
			textStyle.Render(description) + "\n\n" +
			HintStyle("\n\n\nPress any key to continue...")

	case 1:
		title := RenderGradientText("Please keep in mind:", m.lastTick)
		features := RenderGradientText("• This was a passion project, built in a few-days, so if you found it buggy, check for updates hopefully they are fixed :) \n"+
			"\n• I tried to mimic the monkey type styles such as <Space> jumps to  the beginning of next word and ...\n"+
			"\n• I'd be grateful to hear your feedback on github.\n https://github.com/prime-run/go-typer\n"+
			"\n• to see enjoy the TUI you pprobably need a nerd-font and a modern shell and that can't be fixed via docker :)\n"+
			"\nImma head out out now, I hope you enjoy it, gl hf", m.lastTick)
		content = titleStyle.Render(title) + "\n\n" +
			textStyle.Render(features) + "\n\n" +
			HintStyle("\n\nPress any key to start typing...")
	}

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		content)
}

func ShowWelcomeScreen() bool {
	InitSettings()

	// TODO:REMOVE BELOW COMMENT IN PROD
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
