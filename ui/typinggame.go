package ui

//
//TODO: switch to a textare isntead of inpput for the input!
//
import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TypingModel struct {
	textInput  textinput.Model
	targetText string
	width      int
	height     int
}

func NewTypingModel(width, height int) TypingModel {
	ti := textinput.New()
	ti.Placeholder = "Start typing..."
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = MaxWidth

	return TypingModel{
		textInput:  ti,
		targetText: SampleText,
		width:      width,
		height:     height,
	}
}

func (m TypingModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m TypingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			return NewTypingModel(m.width, m.height), nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.textInput.Width = MaxWidth
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m TypingModel) View() string {
	pad := strings.Repeat(" ", Padding)

	formattedText := TextToTypeStyle.Render(m.targetText)
	input := m.textInput.View()
	instructions := HelpStyle("Type the text above. Press ESC to quit, ENTER to restart.")

	content := "\n" +
		pad + "GoTyper - Typing Practice" + "\n\n" +
		pad + formattedText + "\n\n" +
		pad + input + "\n\n" +
		pad + instructions

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
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
