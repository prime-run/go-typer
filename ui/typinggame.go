package ui

//TODO: live hilight placeholders

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
	userTyped        string
	placeholderText  string
	cursorPos        int
	showPlaceholder  bool
}

func NewTypingModel(width, height int) TypingModel {
	ta := textarea.New()
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
		userTyped:        "",
		placeholderText:  SampleText,
		cursorPos:        0,
		showPlaceholder:  true,
	}
}

func (m TypingModel) Init() tea.Cmd {
	cmd := m.updateTextarea()
	return tea.Batch(cmd, textarea.Blink)
}

func (m TypingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

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
		case tea.KeyBackspace:
			if len(m.userTyped) > 0 {
				m.userTyped = m.userTyped[:len(m.userTyped)-1]
				m.updateTextarea()
				return m, nil
			}
		case tea.KeySpace:
			m.userTyped += " "
			m.updateTextarea()
			m.advanceToNextWord()
			return m, nil
		default:
			if len(msg.String()) == 1 {
				m.userTyped += msg.String()
				m.updateTextarea()
				return m, nil
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.textarea.SetWidth(MaxWidth)
	}

	m.textarea, cmd = m.textarea.Update(msg)

	m.lastInput = m.userTyped
	m.typingLog = append(m.typingLog, m.userTyped)

	return m, cmd
}

func (m *TypingModel) updateTextarea() tea.Cmd {
	if len(m.userTyped) == 0 && m.showPlaceholder {
		placeholderText := strings.Repeat("·", len(m.targetText))
		m.textarea.SetValue(placeholderText)
		m.textarea.SetCursor(0)
		return nil
	}

	typed := m.userTyped

	targetWords := strings.Fields(m.targetText)
	typedWords := strings.Fields(typed)

	shift := 0

	for i, typedWord := range typedWords {
		if i < len(targetWords) {
			targetWord := targetWords[i]
			if len(typedWord) > len(targetWord) {
				shift += len(typedWord) - len(targetWord)
			}
		} else {
			shift += len(typedWord) + 1
		}
	}

	var placeholderText string

	if m.showPlaceholder && len(typed)+shift < len(m.targetText) {
		remainingChars := len(m.targetText) - (len(typed) + shift)
		placeholderChars := "·"
		placeholderText = strings.Repeat(placeholderChars, remainingChars)
	}

	m.textarea.SetValue(typed + placeholderText)
	m.textarea.SetCursor(len(typed))

	return nil
}

func (m *TypingModel) advanceToNextWord() {
	targetWords := strings.Fields(m.targetText)
	if m.currentWordIndex < len(targetWords)-1 {
		m.currentWordIndex++
	}
}

func (m TypingModel) CompareWithTarget() string {
	userInput := m.userTyped

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
			result.WriteString(ErrorStyle.Render(typedWord))
		}
	}

	return result.String()
}

func (m TypingModel) compareWord(typed, target string) string {
	var result strings.Builder

	for i, char := range typed {
		if i < len(target) {
			if string(char) == string(target[i]) {
				result.WriteString(string(char))
			} else {
				result.WriteString(ErrorStyle.Render(string(char)))
			}
		} else {
			result.WriteString(ErrorStyle.Render(string(char)))
		}
	}

	if typed == target {
		return InputStyle.Render(typed)
	}

	if len(typed) < len(target) && strings.HasSuffix(typed, " ") {
		return ErrorStyle.Render(result.String())
	}

	return result.String()
}

// func isPrefixOf(s1, s2 string) bool {
// 	if len(s1) > len(s2) {
// 		return false
// 	}
// 	return s2[:len(s1)] == s1
// }

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

	targetDisplay := m.renderTargetWithProgress()

	userTyped := m.CompareWithTarget()
	typingPreview := PreviewStyle.Render("Live Typing Analysis:\n" + userTyped)
	textareaView := m.textarea.View()
	instructions := HelpStyle("Type the text above. Press ESC to quit, TAB to restart.")

	content := "\n" +
		pad + "GoTyper - Typing Practice " + timerDisplay + "\n\n" +
		targetDisplay + "\n\n" +
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

func (m TypingModel) renderTargetWithProgress() string {
	if len(m.userTyped) == 0 {
		return TextToTypeStyle.Render(m.targetText)
	}

	targetWords := strings.Fields(m.targetText)
	typedWords := strings.Fields(m.userTyped)

	var result strings.Builder
	var currentPos int

	for i, targetWord := range targetWords {
		if i > 0 {
			if currentPos < len(m.userTyped) && currentPos < len(m.targetText) && m.targetText[currentPos] == ' ' {
				result.WriteString(InputStyle.Render(" "))
			} else {
				result.WriteString(" ")
			}
			currentPos++
		}

		if i < len(typedWords) {
			typedWord := typedWords[i]

			for j, char := range targetWord {
				if j < len(typedWord) {
					if j < len(typedWord) && j < len(targetWord) &&
						string(typedWord[j]) == string(targetWord[j]) {
						result.WriteString(InputStyle.Render(string(char)))
					} else {
						result.WriteString(ErrorStyle.Render(string(char)))
					}
				} else {
					result.WriteString(string(char))
				}
				currentPos++
			}

			if len(typedWord) > len(targetWord) {
				extraChars := typedWord[len(targetWord):]
				result.WriteString(ErrorStyle.Render(extraChars))
				currentPos += len(extraChars)
			}
		} else {
			result.WriteString(targetWord)
			currentPos += len(targetWord)
		}
	}

	return TextContainerStyle.Render(result.String())
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
