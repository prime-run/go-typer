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
type textFetchedMsg string

type LoadingModel struct {
	progress progress.Model
	done     bool
	width    int
	height   int
	text     string
}

func NewLoadingModel() LoadingModel {
	return LoadingModel{
		progress: progress.New(
			progress.WithGradient("#00ADD8", "#00FFFF"), // Go blue to Cyan
		),
		done: false,
	}
}

func (m LoadingModel) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(),
		fetchTextCmd(),
	)
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
		if m.done {
			return m, nil
		}
		cmd := m.progress.IncrPercent(0.1)
		return m, tea.Batch(tickCmd(), cmd)

	case textFetchedMsg:
		m.done = true
		m.text = string(msg)
		DebugLog("Loading: Fetched text: %s", m.text)
		return StartTypingGame(m.width, m.height, m.text), nil

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	case StartGameMsg:
		selectedCursorType := BlockCursor
		if msg.cursorType == "underline" {
			selectedCursorType = UnderlineCursor
		}
		DefaultCursorType = selectedCursorType

		if msg.theme != "" {
			if err := LoadTheme(msg.theme); err == nil {
				UpdateStyles()
			}
		}

		return m, nil

	default:
		return m, nil
	}
}

func (m LoadingModel) View() string {
	pad := strings.Repeat(" ", Padding)

	content := "\n" +
		pad + "Loading text..." + "\n\n" +
		pad + m.progress.View() + "\n\n" +
		pad + HelpStyle("Fetching text from server...")

	if m.width > 0 {
		return lipgloss.Place(m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			content)
	}

	return content
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func fetchTextCmd() tea.Cmd {
	return func() tea.Msg {
		textCount := map[string]int{
			TextLengthShort:    1,
			TextLengthMedium:   2,
			TextLengthLong:     3,
			TextLengthVeryLong: 5,
		}

		count := textCount[CurrentSettings.TextLength]
		var texts []string

		for i := 0; i < count; i++ {
			text := GetRandomText()
			texts = append(texts, text)
		}

		// Join all texts with a space
		finalText := strings.Join(texts, " ")
		return textFetchedMsg(finalText)
	}
}

func StartLoading(cmd *cobra.Command, args []string) {
	StartLoadingWithOptions("block")
}

func StartLoadingWithOptions(cursorTypeStr string) {
	selectedCursorType := BlockCursor
	if cursorTypeStr == "underline" {
		selectedCursorType = UnderlineCursor
	}

	DefaultCursorType = selectedCursorType

	model := NewLoadingModel()
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Oh no!", err)
		os.Exit(1)
	}
}

func ReloadTheme(filePath string) error {
	err := LoadTheme(filePath)
	if err != nil {
		return err
	}
	UpdateStyles()
	return nil
}
