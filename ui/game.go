package ui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	devlog "github.com/prime-run/go-typer/log"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type TypingModel struct {
	text         *Text
	width        int
	height       int
	startTime    time.Time
	timerRunning bool
	cursorType   CursorType
	lastKeyTime  time.Time
	needsRefresh bool
	gameComplete bool
	lastTick     time.Time
}

func NewTypingModel(width, height int, text string) *TypingModel {
	devlog.Log("Game: Creating new typing model with text: %s", text)
	model := &TypingModel{
		width:        width,
		height:       height,
		timerRunning: false,
		cursorType:   DefaultCursorType,
		needsRefresh: true,
		lastKeyTime:  time.Now(),
		lastTick:     time.Now(),
	}
	model.text = NewText(text)
	model.text.SetCursorType(DefaultCursorType)
	return model
}

func (m *TypingModel) Init() tea.Cmd {
	devlog.Log("Game: Init called")
	return InitGlobalTick()
}

func (m *TypingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case GlobalTickMsg:

		var cmd tea.Cmd
		m.lastTick, _, cmd = HandleGlobalTick(m.lastTick, msg)

		if !m.gameComplete && m.text.GetCursorPos() == len(m.text.words)-1 {
			lastWord := m.text.words[m.text.GetCursorPos()]
			if lastWord.IsComplete() {
				return m.handleGameCompletion()
			}
		}

		return m, cmd

	case tea.KeyMsg:

		if m.gameComplete {
			return m, nil
		}

		m.lastKeyTime = time.Now()

		keyStr := msg.String()
		devlog.Log("Game: Key pressed: %s", keyStr)

		if !m.timerRunning && keyStr != "tab" && keyStr != "esc" && keyStr != "ctrl+c" {
			m.timerRunning = true
			m.startTime = time.Now()
		}

		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyTab:

			newModel := NewTypingModel(m.width, m.height, m.text.GetText())
			return newModel, InitGlobalTick()
		case tea.KeyBackspace:
			m.text.Backspace()
		default:
			if len(keyStr) == 1 {
				m.text.Type([]rune(keyStr)[0])

				if m.text.GetCursorPos() == len(m.text.words)-1 {
					lastWord := m.text.words[m.text.GetCursorPos()]
					if lastWord.IsComplete() {
						return m.handleGameCompletion()
					}
				}
			}
		}

		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	return m, nil
}

func (m *TypingModel) handleGameCompletion() (tea.Model, tea.Cmd) {
	total, correct, errors := m.text.Stats()
	accuracy := 0.0
	if total > 0 {
		accuracy = float64(correct) / float64(total) * 100
	}

	elapsedMinutes := m.lastTick.Sub(m.startTime).Minutes()

	wpm := 0.0
	if elapsedMinutes > 0 {
		wpm = float64(correct*5) / elapsedMinutes / 5
	}

	endModel := NewEndGameModel(wpm, accuracy, total, correct, errors, m.text.GetText())
	endModel.width = m.width
	endModel.height = m.height
	return endModel, InitGlobalTick()
}

func (m *TypingModel) formatElapsedTime() string {
	if !m.timerRunning {
		return "00:00"
	}

	elapsed := m.lastTick.Sub(m.startTime)
	minutes := int(elapsed.Minutes())
	seconds := int(elapsed.Seconds()) % 60

	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

func (m *TypingModel) View() string {
	startTime := time.Now()
	devlog.Log("Game: View rendering started")
	devlog.Log("Game: Current text: %s", m.text.GetText())

	textContent := m.text.Render()
	devlog.Log("Game: Rendered text: %s", textContent)

	if m.gameComplete {
		textContent = lipgloss.NewStyle().
			Foreground(GetColor("text_correct")).
			Render(textContent)
	}

	cursorType := "Underline cursor"
	if m.cursorType == BlockCursor {
		cursorType = "Block cursor"
	}

	modeInfo := cases.Title(language.English).String(CurrentSettings.GameMode) + " mode"
	if CurrentSettings.UseNumbers {
		modeInfo += " with numbers"
	}

	lengthMap := map[string]string{
		TextLengthShort:    "Short passage (1 quote)",
		TextLengthMedium:   "Medium passage (2 quotes)",
		TextLengthLong:     "Long passage (3 quotes)",
		TextLengthVeryLong: "Very Long passage (5 quotes)",
	}

	// FIX:? Render the complete view in one go
	//LOL Bug or feature! i really don't konow what to call it!
	content := lipgloss.NewStyle().
		Width(m.width * 3 / 4).
		Align(lipgloss.Center).
		Render(fmt.Sprintf(
			"\nGoTyper - Typing Practice %s\n\n%s\n\n%s\n\n%s\n%s",
			TimerStyle.Render(m.formatElapsedTime()),
			textContent,
			HintStyle("◾ Type the text above. results would pop when you are done typing.\n◾ Timer will start as soon as you press the first key.\n◾ Paragraph's lenght, gameplay and alot more can be adjusted in settings.\n◾ Press ESC to quit, TAB to reset current passage."),
			SettingsStyle("Current Settings:"),
			HelpStyle(fmt.Sprintf(" • %s • %s • %s", cursorType, modeInfo, lengthMap[CurrentSettings.TextLength])),
		))

	result := lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		content)

	renderTime := time.Since(startTime)
	devlog.Log("Game: View rendering completed in %s", renderTime)

	return result
}

func StartTypingGame(width, height int, text string) tea.Model {
	devlog.Log("Game: Starting typing game with dimensions: %dx%d", width, height)

	startTime := time.Now()
	model := NewTypingModel(width, height, text)
	initTime := time.Since(startTime)

	devlog.Log("Game: Model initialization completed in %s", initTime)
	devlog.Log("Game: Using theme: %s, cursor: %s", CurrentSettings.ThemeName, CurrentSettings.CursorType)

	return model
}

func RunTypingGame() {
	devlog.Log("Game: RunTypingGame started")

	devlog.Log("Game: Running in terminal mode")

	devlog.Log("Game: Running with settings - Theme: %s, Cursor: %s, Mode: %s, UseNumbers: %v",
		CurrentSettings.ThemeName, CurrentSettings.CursorType,
		CurrentSettings.GameMode, CurrentSettings.UseNumbers)

	StartLoadingWithOptions(CurrentSettings.CursorType)
}
