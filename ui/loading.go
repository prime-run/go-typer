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

// Msg types for loading screen
type textFetchedMsg string

type LoadingModel struct {
	spinner     *Spinner
	width       int
	height      int
	progress    float64
	text        string
	lastTick    time.Time
	progressBar progress.Model
}

func NewLoadingModel() *LoadingModel {
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(60),
		progress.WithoutPercentage(),
	)

	return &LoadingModel{
		spinner:     NewSpinner(),
		progress:    0.0,
		lastTick:    time.Now(),
		progressBar: p,
	}
}

func (m *LoadingModel) Init() tea.Cmd {
	return tea.Batch(
		InitGlobalTick(),
		fetchTextCmd(),
	)
}

func (m *LoadingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case GlobalTickMsg:
		// Update spinner and progress
		m.spinner.Update()
		m.progress += 0.03
		if m.progress > 1.0 {
			m.progress = 0.1
		}

		// Handle the global tick
		var cmd tea.Cmd
		m.lastTick, _, cmd = HandleGlobalTick(m.lastTick, msg)
		return m, cmd

	case textFetchedMsg:
		// Text has been fetched, start the game
		m.text = string(msg)
		DebugLog("Loading: Fetched text: %s", m.text)
		return StartTypingGame(m.width, m.height, m.text), nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	return m, nil
}

func (m *LoadingModel) View() string {
	// Make the spinner color transition based on progress
	spinnerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00ADD8"))

	// Pad string helper
	pad := strings.Repeat(" ", 25)

	// Use the bubbles progress bar instead of custom one
	progressBar := m.progressBar.ViewAs(m.progress)

	spinnerDisplay := spinnerStyle.Render(m.spinner.View())

	content := lipgloss.NewStyle().
		Width(m.width * 3 / 4).
		Align(lipgloss.Center).
		Render(
			"\n\n" +
				pad + spinnerDisplay + "\n\n" +
				pad + "Loading text..." + "\n\n" +
				pad + progressBar + "\n\n" +
				pad + HelpStyle("Fetching text from server..."))

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		content)
}

// renderProgressBar renders a simple progress bar
func renderProgressBar(progress float64, width int) string {
	// Ensure progress is between 0 and 1
	if progress < 0 {
		progress = 0
	} else if progress > 1 {
		progress = 1
	}

	// Calculate filled width
	filled := int(progress * float64(width))
	if filled > width {
		filled = width
	}

	// Create the bar
	bar := "["
	for i := 0; i < width; i++ {
		if i < filled {
			bar += "="
		} else {
			bar += " "
		}
	}
	bar += "]"

	return bar
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

		texts := make([]string, 0, count)

		estimatedTotalLen := count * 200

		for i := 0; i < count; i++ {
			text := GetRandomText()
			texts = append(texts, text)
		}

		var finalTextBuilder strings.Builder
		finalTextBuilder.Grow(estimatedTotalLen + count)

		for i, text := range texts {
			finalTextBuilder.WriteString(text)
			if i < len(texts)-1 {
				finalTextBuilder.WriteRune(' ')
			}
		}

		return textFetchedMsg(finalTextBuilder.String())
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
