package ui

//
//TODO: switch to a textare isntead of inpput for the input!
//
import (
	"fmt"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
	"strings"
)

type TypingModel struct {
	textarea   textarea.Model
	targetText string
	width      int
	height     int
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
			// reset the game
			return NewTypingModel(m.width, m.height), nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.textarea.SetWidth(MaxWidth)
	}

	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

func (m TypingModel) View() string {
	pad := strings.Repeat(" ", Padding)

	formattedText := TextToTypeStyle.Render(m.targetText)

	// WARN:DEBUG feature
	userTyped := m.textarea.Value()
	previewStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7F9ABE")).
		Padding(1).
		Width(MaxWidth)

	typingPreview := previewStyle.Render("debug window, live text:\n" + userTyped)

	//TODO: rendering seems magic! lookup the view method , dont trust it
	textareaView := m.textarea.View()
	instructions := HelpStyle("Type the text above. Press ESC to quit, TAB to restart.")

	content := "\n" +
		//TODO: the \n count should be derried from actual text (calculated based on terminal width! oh god i already can feel "fuck microsoft powershell" vibes)
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
