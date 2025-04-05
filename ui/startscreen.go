package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const logoArt = `
   ________    _________                     
  / ____/ /   / ____/   |  ____  ___  _____  
 / / __/ /   / __/ / /| | / __ \/ _ \/ ___/  
/ /_/ / /___/ /___/ ___ |/ /_/ /  __/ /      
\____/_____/_____/_/  |_/ .___/\___/_/       
                       /_/                    
 _______  ____________  ____    __            
/_  __\ \/ /  _/ __ \ \/ / /   / /            
 / /   \  // // /_/ /\  / /   / /             
/ /    / // // ____/ / / /___/ /___           
\____/_/___/_/_/     \/_____/_____/           
`

const (
	MenuMain int = iota
	MenuSettings
)

type menuItem struct {
	title     string
	action    func(*StartScreenModel) tea.Cmd
	disabled  bool
	backColor string
}

type StartScreenModel struct {
	menuState       int
	selectedItem    int
	mainMenuItems   []menuItem
	settingsItems   []menuItem
	width           int
	height          int
	cursorType      string
	selectedTheme   string
	initialTheme    string
	availableThemes []string
	themeChanged    bool
	gameMode        string
	useNumbers      bool
}

func NewStartScreenModel() *StartScreenModel {
	themes := ListAvailableThemes()

	model := &StartScreenModel{
		menuState:       MenuMain,
		selectedItem:    0,
		cursorType:      CurrentSettings.CursorType,
		selectedTheme:   CurrentSettings.ThemeName,
		initialTheme:    CurrentSettings.ThemeName,
		availableThemes: themes,
		themeChanged:    false,
		gameMode:        CurrentSettings.GameMode,
		useNumbers:      CurrentSettings.UseNumbers,
	}

	model.mainMenuItems = []menuItem{
		{title: "Start Typing", action: startGame},
		{title: "Multiplayer Typeracer", action: nil, disabled: true, backColor: "#555555"},
		{title: "Settings", action: openSettings},
		{title: "Statistics", action: openStats, disabled: true, backColor: "#555555"},
		{title: "Quit", action: quitGame},
	}

	model.settingsItems = []menuItem{
		{title: "Theme", action: cycleTheme},
		{title: "Cursor Style", action: cycleCursor},
		{title: "Game Mode", action: cycleGameMode},
		{title: "Use Numbers", action: toggleNumbers},
		{title: "Back", action: saveAndGoBack},
	}

	if model.mainMenuItems[model.selectedItem].disabled {
		for i, item := range model.mainMenuItems {
			if !item.disabled {
				model.selectedItem = i
				break
			}
		}
	}

	return model
}

func (m *StartScreenModel) Init() tea.Cmd {
	return nil
}

func (m *StartScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			previousItem := m.selectedItem
			for {
				m.selectedItem--
				if m.selectedItem < 0 {
					if m.menuState == MenuMain {
						m.selectedItem = len(m.mainMenuItems) - 1
					} else {
						m.selectedItem = len(m.settingsItems) - 1
					}
				}

				if m.menuState == MenuMain {
					if !m.mainMenuItems[m.selectedItem].disabled || m.selectedItem == previousItem {
						break
					}
				} else {
					if !m.settingsItems[m.selectedItem].disabled || m.selectedItem == previousItem {
						break
					}
				}
			}

		case "down", "j":
			previousItem := m.selectedItem
			for {
				m.selectedItem++

				if m.menuState == MenuMain {
					if m.selectedItem >= len(m.mainMenuItems) {
						m.selectedItem = 0
					}

					if !m.mainMenuItems[m.selectedItem].disabled || m.selectedItem == previousItem {
						break
					}
				} else {
					if m.selectedItem >= len(m.settingsItems) {
						m.selectedItem = 0
					}

					if !m.settingsItems[m.selectedItem].disabled || m.selectedItem == previousItem {
						break
					}
				}
			}

		case "enter", " ":
			var items []menuItem
			if m.menuState == MenuMain {
				items = m.mainMenuItems
			} else {
				items = m.settingsItems
			}

			if m.selectedItem < len(items) && !items[m.selectedItem].disabled {
				return m, items[m.selectedItem].action(m)
			}

		case "backspace", "-", "esc":
			if m.menuState != MenuMain {
				if m.themeChanged {
					LoadTheme(m.initialTheme)
					UpdateStyles()
					m.selectedTheme = m.initialTheme
					m.themeChanged = false
				}

				m.cursorType = CurrentSettings.CursorType

				m.menuState = MenuMain
				m.selectedItem = 0

				for m.mainMenuItems[m.selectedItem].disabled {
					m.selectedItem++
					if m.selectedItem >= len(m.mainMenuItems) {
						m.selectedItem = 0
					}
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m *StartScreenModel) View() string {
	var menuContent string

	logoStyle := lipgloss.NewStyle().
		Foreground(GetColor("border")).
		Bold(true)

	logo := logoStyle.Render(logoArt)

	if m.menuState == MenuMain {
		menuContent = m.renderMainMenu()
	} else if m.menuState == MenuSettings {
		menuContent = m.renderSettingsMenu()
	}

	footer := "\n" + HelpStyle("↑/↓: Navigate • Enter: Select • Esc: Back • q: Quit")

	content := fmt.Sprintf("%s\n%s\n%s", logo, menuContent, footer)

	if m.width > 0 && m.height > 0 {
		return lipgloss.Place(m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			content)
	}

	return content
}

func (m *StartScreenModel) renderMainMenu() string {
	var sb strings.Builder

	titleStyle := lipgloss.NewStyle().
		Foreground(GetColor("timer")).
		Bold(true).
		Margin(1, 0, 2, 0)

	sb.WriteString(titleStyle.Render("Main Menu"))
	sb.WriteString("\n\n")

	for i, item := range m.mainMenuItems {
		var s lipgloss.Style

		if i == m.selectedItem {
			s = lipgloss.NewStyle().
				Foreground(GetColor("cursor_bg")).
				Bold(true).
				Padding(0, 4).
				Underline(true)
		} else if item.disabled {
			c := GetColor("text_dim")
			if item.backColor != "" {
				c = lipgloss.Color(item.backColor)
			}
			s = lipgloss.NewStyle().
				Foreground(c).
				Padding(0, 4)
		} else {
			s = lipgloss.NewStyle().
				Foreground(GetColor("text_preview")).
				Padding(0, 4)
		}

		sb.WriteString(s.Render(item.title))
		if item.disabled {
			sb.WriteString(" " + HelpStyle("(coming soon)"))
		}
		sb.WriteString("\n\n")
	}

	return sb.String()
}

func (m *StartScreenModel) renderSettingsMenu() string {
	var sb strings.Builder

	titleStyle := lipgloss.NewStyle().
		Foreground(GetColor("timer")).
		Bold(true).
		Margin(1, 0, 2, 0)

	sb.WriteString(titleStyle.Render("Settings"))
	sb.WriteString("\n\n")

	for i, item := range m.settingsItems {
		var s lipgloss.Style

		if i == m.selectedItem {
			s = lipgloss.NewStyle().
				Foreground(GetColor("cursor_bg")).
				Bold(true).
				Padding(0, 2).
				Underline(true)
		} else {
			s = lipgloss.NewStyle().
				Foreground(GetColor("text_preview")).
				Padding(0, 2)
		}

		menuText := item.title

		if i == 0 {
			menuText = fmt.Sprintf("%-15s: %s", item.title, m.selectedTheme)
		} else if i == 1 {
			menuText = fmt.Sprintf("%-15s: %s", item.title, m.cursorType)
		} else if i == 2 {
			menuText = fmt.Sprintf("%-15s: %s", item.title, m.gameMode)
		} else if i == 3 {
			menuText = fmt.Sprintf("%-15s: %v", item.title, m.useNumbers)
		}

		sb.WriteString(s.Render(menuText))
		sb.WriteString("\n")

		if i == m.selectedItem {
			sb.WriteString("\n")

			if i == 0 {
				exampleBox := renderThemeExample(m.selectedTheme)
				sb.WriteString(exampleBox)
			} else if i == 1 {
				exampleBox := renderCursorExample(m.cursorType)
				sb.WriteString(exampleBox)
			} else if i == 2 {
				exampleBox := renderGameModeExample(m.gameMode)
				sb.WriteString(exampleBox)
			} else if i == 3 {
				exampleBox := renderUseNumbersExample(m.useNumbers)
				sb.WriteString(exampleBox)
			}

			sb.WriteString("\n")
		}

		sb.WriteString("\n")
	}

	return sb.String()
}

func renderThemeExample(themeName string) string {
	exampleStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(GetColor("border")).
		Padding(1, 2).
		Margin(0, 0, 0, 4)

	var example strings.Builder

	example.WriteString("Text Preview: ")
	example.WriteString(DimStyle.Render("dim "))
	example.WriteString(InputStyle.Render("correct "))
	example.WriteString(ErrorStyle.Render("error"))
	example.WriteString("\n\n")

	example.WriteString("Word with errors: ")
	example.WriteString(InputStyle.Render("co"))
	example.WriteString(ErrorStyle.Render("d"))
	example.WriteString(PartialErrorStyle.Render("e"))
	example.WriteString(DimStyle.Render("r"))

	return exampleStyle.Render(example.String())
}

func renderCursorExample(cursorType string) string {
	exampleStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(GetColor("border")).
		Padding(1, 2).
		Margin(0, 0, 0, 4)

	var example strings.Builder

	var cursor *Cursor
	if cursorType == "block" {
		cursor = NewCursor(BlockCursor)
	} else {
		cursor = NewCursor(UnderlineCursor)
	}

	example.WriteString("Cursor appearance: ")
	example.WriteString(cursor.Render('A'))
	example.WriteString(InputStyle.Render("BC"))

	return exampleStyle.Render(example.String())
}

func renderGameModeExample(gameMode string) string {
	exampleStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(GetColor("border")).
		Padding(1, 2).
		Margin(0, 0, 0, 4)

	var example strings.Builder

	example.WriteString("Game Mode: ")

	if gameMode == "normal" {
		example.WriteString("Normal (With Punctuation)")
		example.WriteString("\n\n")
		example.WriteString("Example: ")
		example.WriteString(TextToTypeStyle.Render("The quick brown fox jumps."))
	} else {
		example.WriteString("Simple (No Punctuation)")
		example.WriteString("\n\n")
		example.WriteString("Example: ")
		example.WriteString(TextToTypeStyle.Render("the quick brown fox jumps"))
	}

	return exampleStyle.Render(example.String())
}

// NOTE: render a use numbers example
func renderUseNumbersExample(useNumbers bool) string {
	exampleStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(GetColor("border")).
		Padding(1, 2).
		Margin(0, 0, 0, 4)

	var example strings.Builder

	example.WriteString("Use Numbers: ")

	if useNumbers {
		example.WriteString("Yes")
		example.WriteString("\n\n")
		example.WriteString("Example: ")
		example.WriteString(TextToTypeStyle.Render("quick brown fox jumps over 5 lazy dogs"))
	} else {
		example.WriteString("No")
		example.WriteString("\n\n")
		example.WriteString("Example: ")
		example.WriteString(TextToTypeStyle.Render("quick brown fox jumps over lazy dogs"))
	}

	return exampleStyle.Render(example.String())
}

func startGame(m *StartScreenModel) tea.Cmd {
	return tea.Quit
}

func openSettings(m *StartScreenModel) tea.Cmd {
	m.initialTheme = m.selectedTheme
	m.themeChanged = false

	m.menuState = MenuSettings
	m.selectedItem = 0
	return nil
}

func openStats(m *StartScreenModel) tea.Cmd {
	return nil
}

func quitGame(m *StartScreenModel) tea.Cmd {
	return tea.Quit
}

func saveAndGoBack(m *StartScreenModel) tea.Cmd {
	m.initialTheme = m.selectedTheme
	m.themeChanged = false

	UpdateSettings(UserSettings{
		ThemeName:  m.selectedTheme,
		CursorType: m.cursorType,
		GameMode:   m.gameMode,
		UseNumbers: m.useNumbers,
	})

	m.menuState = MenuMain
	m.selectedItem = 0
	return nil
}

func cycleTheme(m *StartScreenModel) tea.Cmd {
	currentIndex := -1
	for i, theme := range m.availableThemes {
		if theme == m.selectedTheme {
			currentIndex = i
			break
		}
	}

	currentIndex = (currentIndex + 1) % len(m.availableThemes)
	m.selectedTheme = m.availableThemes[currentIndex]
	m.themeChanged = true

	LoadTheme(m.selectedTheme)
	UpdateStyles()

	return nil
}

func cycleCursor(m *StartScreenModel) tea.Cmd {
	if m.cursorType == "block" {
		m.cursorType = "underline"
	} else {
		m.cursorType = "block"
	}

	return nil
}

func cycleGameMode(m *StartScreenModel) tea.Cmd {
	if m.gameMode == "normal" {
		m.gameMode = "simple"
	} else {
		m.gameMode = "normal"
	}

	return nil
}

func toggleNumbers(m *StartScreenModel) tea.Cmd {
	m.useNumbers = !m.useNumbers
	return nil
}

type StartGameMsg struct {
	cursorType string
	theme      string
}

func RunStartScreen() {
	p := tea.NewProgram(NewStartScreenModel(), tea.WithAltScreen())

	model, err := p.Run()
	if err != nil {
		fmt.Println("Error running start screen:", err)
		return
	}

	if m, ok := model.(*StartScreenModel); ok {
		UpdateSettings(UserSettings{
			ThemeName:  m.selectedTheme,
			CursorType: m.cursorType,
			GameMode:   m.gameMode,
			UseNumbers: m.useNumbers,
		})

		if m.menuState == MenuMain && m.selectedItem < len(m.mainMenuItems) {
			item := m.mainMenuItems[m.selectedItem]

			if item.title == "Start Typing" {
				StartLoadingWithOptions(m.cursorType)
			}
		}
	}
}
