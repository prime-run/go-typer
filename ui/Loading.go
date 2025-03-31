package ui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

const (
	padding  = 2
	maxWidth = 80
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render
var textToTypeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#333333")).Padding(1).Width(maxWidth)
var inputStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
var errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))

const sampleText = "The quick brown fox jumps over the lazy dog. Programming is the process of creating a set of instructions that tell a computer how to perform a task. Programming can be done using a variety of computer programming languages, such as JavaScript, Python, and C++."

func ProgressBar(cmd *cobra.Command, args []string) {
	m := loadingModel{
		progress: progress.New(progress.WithDefaultGradient()),
		done:     false,
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Oh no!", err)
		os.Exit(1)
	}
}

type tickMsg time.Time

type loadingModel struct {
	progress progress.Model
	done     bool
}

func (m loadingModel) Init() tea.Cmd {
	return tickCmd()
}

func (m loadingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case tickMsg:
		if m.progress.Percent() == 1.0 {
			return newTypingModel(), nil
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

func (m loadingModel) View() string {
	pad := strings.Repeat(" ", padding)
	return "\n" +
		pad + "Loading text..." + "\n\n" +
		pad + m.progress.View() + "\n\n" +
		pad + helpStyle("Text is being prepared...")
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type typingModel struct {
	textInput  textinput.Model
	targetText string
	width      int
}

func newTypingModel() typingModel {
	ti := textinput.New()
	ti.Placeholder = "Start typing..."
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = maxWidth

	return typingModel{
		textInput:  ti,
		targetText: sampleText,
		width:      maxWidth,
	}
}

func (m typingModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m typingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			// Restart the game
			return newTypingModel(), nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		if m.width > maxWidth {
			m.width = maxWidth
		}
		m.textInput.Width = m.width
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m typingModel) View() string {
	pad := strings.Repeat(" ", padding)

	formattedText := textToTypeStyle.Render(m.targetText)

	input := m.textInput.View()

	instructions := helpStyle("Type the text above. Press ESC to quit, ENTER to restart.")

	return "\n" +
		pad + "GoTyper - Typing Practice" + "\n\n" +
		pad + formattedText + "\n\n" +
		pad + input + "\n\n" +
		pad + instructions
}

func GameSession(cmd *cobra.Command, args []string) {
	cmd.Println("Starting a new game session...")
	p := tea.NewProgram(newTypingModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
