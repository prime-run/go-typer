package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type gameTickMsg time.Time
type gameStateMsg struct {
	text         *Text
	startTime    time.Time
	timerRunning bool
}

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
}

func NewTypingModel(width, height int, text string) TypingModel {
	DebugLog("Game: Creating new typing model with text: %s", text)
	model := TypingModel{
		width:        width,
		height:       height,
		timerRunning: false,
		cursorType:   DefaultCursorType,
		needsRefresh: true,
		lastKeyTime:  time.Now(),
	}
	model.text = NewText(text)
	model.text.SetCursorType(DefaultCursorType)
	return model
}

func (m TypingModel) Init() tea.Cmd {
	DebugLog("Game: Init called")
	return tea.Batch(
		gameTickCommand(),
	)
}

func gameTickCommand() tea.Cmd {
	return tea.Tick(time.Millisecond*33, func(t time.Time) tea.Msg {
		return gameTickMsg(t)
	})
}

func (m TypingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case gameTickMsg:
		if m.needsRefresh {
			m.needsRefresh = false
			return m, nil
		}

		// Check for game completion
		if !m.gameComplete && m.text.GetCursorPos() == len(m.text.words)-1 {
			lastWord := m.text.words[m.text.GetCursorPos()]
			if lastWord.IsComplete() {
				return m.handleGameCompletion()
			}
		}

		// refreshrate WARN:migh cause alot of issues!
		if m.timerRunning && time.Since(m.lastKeyTime) > 5*time.Second {
			return m, gameTickCommand()
		}

		return m, gameTickCommand()

	case tea.KeyMsg:
		// If game is complete, ignore all typing input
		if m.gameComplete {
			return m, nil
		}

		m.lastKeyTime = time.Now()
		m.needsRefresh = true

		keyStr := msg.String()
		DebugLog("Game: Key pressed: %s", keyStr)

		if !m.timerRunning && keyStr != "tab" && keyStr != "esc" && keyStr != "ctrl+c" {
			m.timerRunning = true
			m.startTime = time.Now()
		}

		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyTab:
			// Restart with the same text
			return NewTypingModel(m.width, m.height, m.text.GetText()), gameTickCommand()
		case tea.KeyBackspace:
			m.text.Backspace()
			return m, nil
		default:
			if len(keyStr) == 1 {
				m.text.Type([]rune(keyStr)[0])

				// Check for completion after each keystroke
				if m.text.GetCursorPos() == len(m.text.words)-1 {
					lastWord := m.text.words[m.text.GetCursorPos()]
					if lastWord.IsComplete() {
						return m.handleGameCompletion()
					}
				}
			}
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	return m, nil
}

// handleGameCompletion handles the game completion logic and returns the end game model
func (m TypingModel) handleGameCompletion() (tea.Model, tea.Cmd) {
	m.timerRunning = false // Stop the timer
	total, correct, errors := m.text.Stats()
	accuracy := 0.0
	if total > 0 {
		accuracy = float64(correct) / float64(total) * 100
	}
	elapsedMinutes := time.Since(m.startTime).Minutes()
	wpm := 0.0
	if elapsedMinutes > 0 {
		wpm = float64(correct*5) / elapsedMinutes / 5
	}

	// Create and initialize the end game model
	endModel := NewEndGameModel(wpm, accuracy, total, correct, errors, m.text.GetText())
	endModel.width = m.width
	endModel.height = m.height
	return endModel, nil
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
	startTime := time.Now()
	DebugLog("Game: View rendering started")
	DebugLog("Game: Current text: %s", m.text.GetText())

	// Get the text content
	textContent := m.text.Render()
	DebugLog("Game: Rendered text: %s", textContent)

	// If game is complete, add a completion message
	if m.gameComplete {
		textContent = lipgloss.NewStyle().
			Foreground(GetColor("text_correct")).
			Render(textContent)
	}

	content := lipgloss.NewStyle().
		Width(m.width * 3 / 4).
		Align(lipgloss.Center).
		Render(
			"\n" +
				"GoTyper - Typing Practice " + TimerStyle.Render(m.formatElapsedTime()) + "\n\n" +
				textContent + "\n\n" +
				HintStyle("⚫ Type the text above. results would pop when you are done typing.\n⚫ Timer will start as soon as you press the first key.\n⚫ Paragraph's lenght, gameplay and alot more can be adjusted in settings.\n⚫ Press ESC to quit, TAB to reset current passage.") +
				"\n\n" +
				SettingsStyle("Current Settings:") +
				HelpStyle(" • "+
					(func() string {
						var settings []string

						// Cursor type
						cursorType := "Underline cursor"
						if m.cursorType == BlockCursor {
							cursorType = "Block cursor"
						}
						settings = append(settings, cursorType)

						// Game mode and numbers
						modeInfo := cases.Title(language.English).String(CurrentSettings.GameMode) + " mode"
						if CurrentSettings.UseNumbers {
							modeInfo += " with numbers"
						}
						settings = append(settings, modeInfo)

						// Text length
						lengthMap := map[string]string{
							TextLengthShort:    "Short passage (1 quote)",
							TextLengthMedium:   "Medium passage (2 quotes)",
							TextLengthLong:     "Long passage (3 quotes)",
							TextLengthVeryLong: "Very Long passage (5 quotes)",
						}
						settings = append(settings, lengthMap[CurrentSettings.TextLength])

						return strings.Join(settings, " • ")
					})()))

	result := lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		content)

	renderTime := time.Since(startTime)
	DebugLog("Game: View rendering completed in %s", renderTime)

	return result
}

func StartTypingGame(width, height int, text string) tea.Model {
	DebugLog("Game: Starting typing game with dimensions: %dx%d", width, height)

	startTime := time.Now()
	model := NewTypingModel(width, height, text)
	initTime := time.Since(startTime)

	DebugLog("Game: Model initialization completed in %s", initTime)
	DebugLog("Game: Using theme: %s, cursor: %s", CurrentSettings.ThemeName, CurrentSettings.CursorType)

	return model
}

func RunTypingGame() {
	DebugLog("Game: RunTypingGame started")

	DebugLog("Game: Running in terminal mode")

	DebugLog("Game: Running with settings - Theme: %s, Cursor: %s, Mode: %s, UseNumbers: %v",
		CurrentSettings.ThemeName, CurrentSettings.CursorType,
		CurrentSettings.GameMode, CurrentSettings.UseNumbers)

	// Start with the loading screen
	StartLoadingWithOptions(CurrentSettings.CursorType)
}
