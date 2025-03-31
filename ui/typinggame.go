package ui

//
//TODO: live hilight placeholders
//
import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TypingModel struct {
	textarea     textarea.Model
	targetText   string
	width        int
	height       int
	typingLog    []string
	startTime    time.Time
	timerRunning bool
}

func NewTypingModel(width, height int) TypingModel {
	ta := textarea.New()
	ta.Placeholder = "Start typing here..."
	ta.Focus()
	ta.SetWidth(MaxWidth)
	ta.SetHeight(3)

	return TypingModel{
		textarea:     ta,
		targetText:   SampleText,
		width:        width,
		height:       height,
		typingLog:    []string{},
		timerRunning: false,
	}
}

func (m TypingModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m TypingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyTab:
			return NewTypingModel(m.width, m.height), nil
		}

		// Start the timer on first keystroke
		if !m.timerRunning && msg.String() != "tab" && msg.String() != "esc" && msg.String() != "ctrl+c" {
			m.timerRunning = true
			m.startTime = time.Now()
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.textarea.SetWidth(MaxWidth)
	}

	m.textarea, cmd = m.textarea.Update(msg)
	m.typingLog = append(m.typingLog, m.textarea.Value())

	return m, cmd
}

// compareWithTarget compares word by word NOTE:might be useful in some gamemod! but not for the main part
//
//	returns a formatted string with errors highlighted in red
func (m TypingModel) CompareWithTarget() string {
	userInput := m.textarea.Value()

	if userInput == "" {
		return "Start typing..."
	}

	targetWords := strings.Fields(m.targetText)
	userWords := strings.Fields(userInput)

	var result strings.Builder

	for i, userWord := range userWords {
		if i > 0 {
			result.WriteString(" ")
		}

		if i < len(targetWords) {
			targetWord := targetWords[i]

			if userWord == targetWord {
				// All good with CW (cw = current word)
				result.WriteString(InputStyle.Render(userWord))
			} else if isPrefixOf(userWord, targetWord) {
				result.WriteString(userWord)
			} else {
				// complete and incorrect, or incomplete and wrong already
				result.WriteString(ErrorStyle.Render(userWord))
			}
		} else {
			// player just smashed his/her face on the keyboard
			result.WriteString(ErrorStyle.Render(userWord))
		}
	}

	return result.String()
}

func isPrefixOf(s1, s2 string) bool {
	if len(s1) > len(s2) {
		return false
	}
	return s2[:len(s1)] == s1
}

func (m TypingModel) formatElapsedTime() string {
	if !m.timerRunning {
		return "00:00"
	}

	elapsed := time.Since(m.startTime)
	minutes := int(elapsed.Minutes())
	seconds := int(elapsed.Seconds()) % 60

	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

func (m TypingModel) View() string {
	padStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	pad := padStyle.Render(strings.Repeat(" ", Padding))

	// TODO: when the timer needs more features, ove to ui/time.go
	timerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFDB58")).
		Bold(true).
		Padding(0, 1)
	timerDisplay := timerStyle.Render(m.formatElapsedTime())

	formattedText := TextToTypeStyle.Render(m.targetText)

	userTyped := m.CompareWithTarget()
	previewStyle := lipgloss.NewStyle().
		Padding(1).
		Margin(8, 0, 0, 0).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7F9ABE")).
		Width(MaxWidth)

	typingPreview := previewStyle.Render("---DEBUG window ---:\n LIVE text eval \n" + userTyped)

	textareaView := m.textarea.View()
	instructions := HelpStyle("Type the text above. Press ESC to quit, TAB to restart.")

	content := "\n" +
		pad + "GoTyper - Typing Practice " + timerDisplay + "\n\n" +
		formattedText + "\n\n" +
		textareaView + "\n\n" +
		instructions + "\n\n" +
		typingPreview

	if m.width > 0 {
		return lipgloss.Place(m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			content)
	}

	return content
}

func StartTypingGame(width, height int) tea.Model {
	return NewTypingModel(width, height)
}

func RunTypingGame() {
	p := tea.NewProgram(NewTypingModel(0, 0), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Some thing went wrong: %v\n", err)
		os.Exit(1)
	}
}
