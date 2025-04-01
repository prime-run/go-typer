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
	textarea         textarea.Model
	targetText       string
	width            int
	height           int
	typingLog        []string
	startTime        time.Time
	timerRunning     bool
	currentWordIndex int
	lastInput        string
}

func NewTypingModel(width, height int) TypingModel {
	ta := textarea.New()
	ta.Placeholder = "Start typing here..."
	ta.Focus()
	ta.SetWidth(MaxWidth)
	ta.SetHeight(3)

	return TypingModel{
		textarea:         ta,
		targetText:       SampleText,
		width:            width,
		height:           height,
		typingLog:        []string{},
		timerRunning:     false,
		currentWordIndex: 0,
		lastInput:        "",
	}
}

func (m TypingModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m TypingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Start the timer on first keystroke
	if !m.timerRunning {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() != "tab" && msg.String() != "esc" && msg.String() != "ctrl+c" {
				m.timerRunning = true
				m.startTime = time.Now()
			}
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyTab:
			return NewTypingModel(m.width, m.height), nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.textarea.SetWidth(MaxWidth)
	}

	// Update the textarea first to capture all keystroke input
	m.textarea, cmd = m.textarea.Update(msg)

	// Check if a space was just added
	currentInput := m.textarea.Value()
	if len(currentInput) > len(m.lastInput) &&
		strings.HasSuffix(currentInput, " ") &&
		!strings.HasSuffix(m.lastInput, " ") {
		// Space was just pressed, handle word advancement
		m.advanceToNextWord()
	}

	m.lastInput = currentInput
	m.typingLog = append(m.typingLog, currentInput)

	return m, cmd
}

// advanceToNextWord moves to the next word in the target text
func (m *TypingModel) advanceToNextWord() {
	targetWords := strings.Fields(m.targetText)
	if m.currentWordIndex < len(targetWords)-1 {
		m.currentWordIndex++
	}
}

// CompareWithTarget compares character by character and applies MonkeyType-style validation
func (m TypingModel) CompareWithTarget() string {
	userInput := m.textarea.Value()

	if userInput == "" {
		return "Start typing..."
	}

	targetWords := strings.Fields(m.targetText)
	typedContent := userInput
	typedWords := strings.Split(typedContent, " ")

	var result strings.Builder

	for i, typedWord := range typedWords {
		if i > 0 {
			result.WriteString(" ")
		}

		if typedWord == "" {
			continue
		}

		if i < len(targetWords) {
			targetWord := targetWords[i]
			result.WriteString(m.compareWord(typedWord, targetWord))
		} else {
			// Extra words beyond target text
			result.WriteString(ErrorStyle.Render(typedWord))
		}
	}

	return result.String()
}

// compareWord compares single words character by character
func (m TypingModel) compareWord(typed, target string) string {
	var result strings.Builder

	// Process character by character
	for i, char := range typed {
		if i < len(target) {
			// Character within target word range
			if string(char) == string(target[i]) {
				// Correct character
				result.WriteString(string(char))
			} else {
				// Incorrect character
				result.WriteString(ErrorStyle.Render(string(char)))
			}
		} else {
			// Extra characters beyond target word
			result.WriteString(ErrorStyle.Render(string(char)))
		}
	}

	// If the word is completely and correctly typed
	if typed == target {
		return InputStyle.Render(typed)
	}

	// If the word is undertyped (they pressed space before finishing)
	if len(typed) < len(target) && strings.HasSuffix(typed, " ") {
		return ErrorStyle.Render(result.String())
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
	pad := PadStyle.Render(strings.Repeat(" ", Padding))
	timerDisplay := TimerStyle.Render(m.formatElapsedTime())
	formattedText := TextToTypeStyle.Render(m.targetText)
	userTyped := m.CompareWithTarget()
	typingPreview := PreviewStyle.Render("Live Typing Analysis:\n" + userTyped)
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
