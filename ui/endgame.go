package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
}

func NewEndGameModel(wpm, accuracy float64, words, correct, errors int, text string) EndGameModel {
	return EndGameModel{
		selectedItem: 0,
		wpm:          wpm,
		accuracy:     accuracy,
		words:        words,
		correct:      correct,
		errors:       errors,
		text:         text,
	}
}

func (m EndGameModel) Init() tea.Cmd {
	return nil
}

func (m EndGameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			m.selectedItem--
			if m.selectedItem < 0 {
				m.selectedItem = 1
			}

		case "down", "j":
			m.selectedItem++
			if m.selectedItem > 1 {
				m.selectedItem = 0
			}

		case "enter", " ":
			switch m.selectedItem {
			case 0: // Play with Same Text
				return NewTypingModel(m.width, m.height, m.text), nil
			case 1: // Play with New Text
				StartLoadingWithOptions(CurrentSettings.CursorType)
				return m, tea.Quit
			}

		case "esc":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m EndGameModel) View() string {
	// Title with border
	title := EndGameTitleStyle.Render("Game Complete!")

	// Stats with different colors and styles
	stats := fmt.Sprintf("%s %.1f | %s %.1f%% | %s %d | %s %d | %s %d",
		EndGameWpmStyle.Render("WPM:"), m.wpm,
		EndGameAccuracyStyle.Render("Accuracy:"), m.accuracy,
		EndGameWordsStyle.Render("Words:"), m.words,
		EndGameCorrectStyle.Render("Correct:"), m.correct,
		EndGameErrorsStyle.Render("Errors:"), m.errors)

	// Stats container with border
	statsBox := EndGameStatsBoxStyle.Render(stats)

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
				statsBox + "\n\n" +
				menu + "\n\n" +
				HelpStyle("Use arrow keys to navigate, enter to select, esc to quit"),
		)

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		content)
}
