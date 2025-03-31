package ui

//
//TODO: live hilight placeholders
//
import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TypingModel struct {
	textarea   textarea.Model
	targetText string
	width      int
	height     int
	typingLog  []string
}

func NewTypingModel(width, height int) TypingModel {
	ta := textarea.New()
	ta.Placeholder = "Start typing here..."
	ta.Focus()
	ta.SetWidth(MaxWidth)
	ta.SetHeight(3)

	return TypingModel{
		textarea:   ta,
		targetText: SampleText,
		width:      width,
		height:     height,
		typingLog:  []string{},
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
		case tea.KeyTab: //reset or tabshift? do we need tabshift in typing game ?!
			return NewTypingModel(m.width, m.height), nil
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
			if userWord == targetWords[i] {
				// green word
				result.WriteString(InputStyle.Render(userWord))
			} else {
				result.WriteString(ErrorStyle.Render(userWord))
			}
		} else {
			// extra stuff
			result.WriteString(ErrorStyle.Render(userWord))
		}
	}

	return result.String()
}

func (m TypingModel) View() string {
	padStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	pad := padStyle.Render(strings.Repeat(" ", Padding))

	formattedText := TextToTypeStyle.Render(m.targetText)

	userTyped := m.CompareWithTarget()
	previewStyle := lipgloss.NewStyle().
		Padding(1).
		Margin(10, 0, 0, 0).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#000000")).
		Width(MaxWidth)

	typingPreview := previewStyle.Render("---DEBUG ONLY---:\n" + userTyped)

	textareaView := m.textarea.View()
	instructions := HelpStyle("Type the text above. Press ESC to quit, TAB to restart.")

	content := "\n" +
		//FIX: the \n count should be derried from actual text (calculated based on terminal width! oh god i already can feel "fuck microsoft powershell" vibes)
		pad + "GoTyper - Typing Practice" + "\n\n" +
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
