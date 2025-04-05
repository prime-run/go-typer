package ui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
	"time"
)

type EndGameModel struct {
	selectedItem int
	width        int
	height       int
	wpm          float64
	accuracy     float64
	words        int
	correct      int
	errors       int
	text         string
	startTime    time.Time
	lastTick     time.Time
}

func NewEndGameModel(wpm, accuracy float64, words, correct, errors int, text string) *EndGameModel {
	return &EndGameModel{
		selectedItem: 0,
		wpm:          wpm,
		accuracy:     accuracy,
		words:        words,
		correct:      correct,
		errors:       errors,
		text:         text,
		startTime:    time.Now(),
		lastTick:     time.Now(),
	}
}

func (m *EndGameModel) Init() tea.Cmd {
	return InitGlobalTick()
}

func (m *EndGameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case GlobalTickMsg:
		var cmd tea.Cmd
		m.lastTick, _, cmd = HandleGlobalTick(m.lastTick, msg)
		return m, cmd

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			m.selectedItem--
			if m.selectedItem < 0 {
				m.selectedItem = 1
			}
			return m, nil

		case "down", "j":
			m.selectedItem++
			if m.selectedItem > 1 {
				m.selectedItem = 0
			}
			return m, nil

		case "enter", " ":
			switch m.selectedItem {
			case 0:
				return NewTypingModel(m.width, m.height, m.text), InitGlobalTick()
			case 1:
				StartLoadingWithOptions(CurrentSettings.CursorType)
				return m, tea.Quit
			}

		case "esc":
			return m, tea.Quit
		}

		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	return m, nil
}

func (m *EndGameModel) View() string {

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(GetColor("text_correct")).
		Render("Game Complete!")

	wpmStyle := lipgloss.NewStyle().Foreground(GetColor("timer"))
	accuracyStyle := lipgloss.NewStyle().Foreground(GetColor("text_correct"))
	wordsStyle := lipgloss.NewStyle().Foreground(GetColor("text_preview"))
	correctStyle := lipgloss.NewStyle().Foreground(GetColor("text_correct"))
	errorsStyle := lipgloss.NewStyle().Foreground(GetColor("text_error"))

	wpmText := RenderGradientOverlay(fmt.Sprintf("WPM: %.1f", m.wpm), wpmStyle, m.lastTick)
	accuracyText := RenderGradientOverlay(fmt.Sprintf("Accuracy: %.1f%%", m.accuracy), accuracyStyle, m.lastTick)
	wordsText := RenderGradientOverlay(fmt.Sprintf("Words: %d", m.words), wordsStyle, m.lastTick)
	correctText := RenderGradientOverlay(fmt.Sprintf("Correct: %d", m.correct), correctStyle, m.lastTick)
	errorsText := RenderGradientOverlay(fmt.Sprintf("Errors: %d", m.errors), errorsStyle, m.lastTick)

	stats := fmt.Sprintf("%s   %s   %s   %s   %s",
		wpmText, accuracyText, wordsText, correctText, errorsText)

	options := []string{
		"Play with Same Text",
		"Play with New Text",
	}

	var menuItems []string
	for i, option := range options {
		cursor := " "
		style := EndGameOptionStyle
		if m.selectedItem == i {
			cursor = ">"
			style = EndGameSelectedOptionStyle
		}
		menuItems = append(menuItems, style.Render(fmt.Sprintf("%s %s", cursor, option)))
	}

	menu := strings.Join(menuItems, "\n")

	content := lipgloss.NewStyle().
		Width(m.width * 3 / 4).
		Align(lipgloss.Center).
		Render(
			"\n" +
				title + "\n\n" +
				stats + "\n\n" +
				menu + "\n\n" +
				HelpStyle("Use arrow keys to navigate, enter to select, esc to quit"),
		)

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		content)
}
