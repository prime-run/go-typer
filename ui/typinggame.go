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
	expectedNextPos  int
	correctChars     int
	incorrectChars   int
	typedChars       []rune
}

func NewTypingModel(width, height int) TypingModel {
	ta := textarea.New()
	ta.Focus()
	ta.SetWidth(MaxWidth)
	ta.SetHeight(3)

	// Remove textarea styling and set actual target text as initial content
	ta.Prompt = ""
	ta.FocusedStyle.Base = lipgloss.NewStyle()
	ta.BlurredStyle.Base = lipgloss.NewStyle()
	ta.SetValue(SampleText)

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
		expectedNextPos:  0,
		correctChars:     0,
		incorrectChars:   0,
		typedChars:       []rune{},
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
			if len(m.typedChars) > 0 {
				m.typedChars = m.typedChars[:len(m.typedChars)-1]
				if m.expectedNextPos > 0 {
					m.expectedNextPos--
				}
				m.updateTextarea()
				return m, nil
			}
		case tea.KeySpace:
			m.typedChars = append(m.typedChars, ' ')

			if m.expectedNextPos < len(m.targetText) {
				if m.targetText[m.expectedNextPos] == ' ' {
					m.correctChars++
				} else {
					m.incorrectChars++
					for m.expectedNextPos < len(m.targetText) && m.targetText[m.expectedNextPos] != ' ' {
						m.expectedNextPos++
					}
				}

				if m.expectedNextPos < len(m.targetText) {
					m.expectedNextPos++
				}
			}

			m.updateTextarea()
			return m, nil
		default:
			if len(msg.String()) == 1 {
				char := []rune(msg.String())[0]
				m.typedChars = append(m.typedChars, char)

				if m.expectedNextPos < len(m.targetText) {
					if m.targetText[m.expectedNextPos] == byte(char) {
						m.correctChars++
					} else {
						m.incorrectChars++
						if m.targetText[m.expectedNextPos] == ' ' {
							m.expectedNextPos++
						}
					}

					if m.expectedNextPos < len(m.targetText) {
						m.expectedNextPos++
					}
				} else {
					m.incorrectChars++
				}

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
	m.userTyped = string(m.typedChars)
	m.lastInput = m.userTyped
	m.typingLog = append(m.typingLog, m.userTyped)

	return m, cmd
}

func (m *TypingModel) updateTextarea() tea.Cmd {
	m.userTyped = string(m.typedChars)

	if m.expectedNextPos >= len(m.targetText) {
		m.textarea.SetValue(m.userTyped)
		m.textarea.SetCursor(len(m.userTyped))
		return nil
	}

	remainingPortion := m.targetText[m.expectedNextPos:]
	displayText := m.userTyped + remainingPortion

	m.textarea.SetValue(displayText)
	m.textarea.SetCursor(len(m.userTyped))

	return nil
}

func (m TypingModel) CompareWithTarget() string {
	if len(m.typedChars) == 0 {
		return "Start typing..."
	}

	var result strings.Builder

	for i, char := range m.typedChars {
		if i < len(m.targetText) {
			if byte(char) == m.targetText[i] {
				result.WriteString(InputStyle.Render(string(char)))
			} else {
				result.WriteString(ErrorStyle.Render(string(char)))
			}
		} else {
			result.WriteString(ErrorStyle.Render(string(char)))
		}
	}

	return result.String()
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
	if len(m.typedChars) == 0 {
		return TextToTypeStyle.Render(m.targetText)
	}

	var result strings.Builder

	for i, char := range m.targetText {
		if i < len(m.typedChars) {
			typedChar := m.typedChars[i]
			if typedChar == char {
				result.WriteString(InputStyle.Render(string(char)))
			} else {
				result.WriteString(ErrorStyle.Render(string(char)))
			}
		} else {
			result.WriteString(string(char))
		}
	}

	if len(m.typedChars) > len(m.targetText) {
		extraChars := string(m.typedChars[len(m.targetText):])
		result.WriteString(ErrorStyle.Render(extraChars))
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
