package ui

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
}

func NewTypingModel(width, height int) TypingModel {
	text := NewText(GetSampleText())
	text.SetCursorType(DefaultCursorType)
	return TypingModel{
		text:         text,
		width:        width,
		height:       height,
		timerRunning: false,
		cursorType:   DefaultCursorType,
		needsRefresh: true,
		lastKeyTime:  time.Now(),
	}
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
			DebugLog("Game: Refresh triggered by UI change")
			return m, nil
		}

		// refreshrate WARN:migh cause alot of issues!
		if m.timerRunning && time.Since(m.lastKeyTime) > 5*time.Second {
			return m, gameTickCommand()
		}

		return m, gameTickCommand()

	case tea.KeyMsg:
		m.lastKeyTime = time.Now()
		m.needsRefresh = true

		keyStr := msg.String()
		DebugLog("Game: Key pressed: %s", keyStr)

		if !m.timerRunning && keyStr != "tab" && keyStr != "esc" && keyStr != "ctrl+c" {
			m.timerRunning = true
			m.startTime = time.Now()
			DebugLog("Game: Timer started")
		}

		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			DebugLog("Game: Quitting game")
			return m, tea.Quit
		case tea.KeyTab:
			DebugLog("Game: Restarting game")
			return NewTypingModel(m.width, m.height), gameTickCommand()
		case tea.KeyBackspace:
			before := m.text.GetCursorPos()
			m.text.Backspace()
			after := m.text.GetCursorPos()
			DebugLog("Game: Backspace - cursor moved from %d to %d", before, after)
			return m, nil
		default:
			if len(keyStr) == 1 {
				before := m.text.GetCursorPos()
				m.text.Type([]rune(keyStr)[0])
				after := m.text.GetCursorPos()
				DebugLog("Game: Typed '%s' - cursor moved from %d to %d", keyStr, before, after)
			}
			return m, nil
		}

	case tea.WindowSizeMsg:
		DebugLog("Game: Window size changed: %dx%d", msg.Width, msg.Height)
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	return m, nil
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

	content := lipgloss.NewStyle().
		Width(m.width * 3 / 4).
		Align(lipgloss.Center).
		Render(
			"\n" +
				"GoTyper - Typing Practice " + TimerStyle.Render(m.formatElapsedTime()) + "\n\n" +
				m.text.Render() + "\n\n" +
				HelpStyle("Type the text above. Press ESC to quit, TAB to restart. "+
					"Using "+(func() string {
					if m.cursorType == BlockCursor {
						return "Block"
					}
					return "Underline"
				})()+" cursor. "+
					(func() string {
						modeInfo := fmt.Sprintf("%s mode", strings.Title(CurrentSettings.GameMode))
						if CurrentSettings.UseNumbers {
							return modeInfo + " with numbers"
						}
						return modeInfo + " without numbers"
					})()) +
				(func() string {
					if !m.timerRunning {
						return ""
					}

					total, correct, errors := m.text.Stats()
					accuracy := 0.0
					if total > 0 {
						accuracy = float64(correct) / float64(total) * 100
					}

					elapsedMinutes := time.Since(m.startTime).Minutes()
					wpm := 0.0
					if elapsedMinutes > 0 {
						wpm = float64(total*5) / elapsedMinutes / 5
					}

					return fmt.Sprintf("\n\nWPM: %.1f | Accuracy: %.1f%% | Words: %d | Correct: %d | Errors: %d",
						wpm, accuracy, total, correct, errors)
				})(),
		)

	result := lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		content)

	renderTime := time.Since(startTime)
	DebugLog("Game: View rendering completed in %s", renderTime)

	return result
}

func StartTypingGame(width, height int) tea.Model {
	DebugLog("Game: Starting typing game with dimensions: %dx%d", width, height)

	startTime := time.Now()
	model := NewTypingModel(width, height)
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

	p := tea.NewProgram(NewTypingModel(0, 0), tea.WithAltScreen())

	DebugLog("Game: Starting main program loop")
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running typing game: %v\n", err)
		DebugLog("Game: Error running typing game: %v", err)
		os.Exit(1)
	}
	DebugLog("Game: Typing game exited normally")
}
