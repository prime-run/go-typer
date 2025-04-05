package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const logoArt = `
   █████████     ███████       ███████████                                        
  ███░░░░░███  ███░░░░░███    ░█░░░███░░░█                                        
 ███     ░░░  ███     ░░███   ░   ░███  ░  █████ ████ ████████   ██████  ████████ 
░███         ░███      ░███       ░███    ░░███ ░███ ░░███░░███ ███░░███░░███░░███
░███    █████░███      ░███       ░███     ░███ ░███  ░███ ░███░███████  ░███ ░░░ 
░░███  ░░███ ░░███     ███        ░███     ░███ ░███  ░███ ░███░███░░░   ░███     
 ░░█████████  ░░░███████░         █████    ░░███████  ░███████ ░░██████  █████    
  ░░░░░░░░░     ░░░░░░░          ░░░░░      ░░░░░███  ░███░░░   ░░░░░░  ░░░░░     
                                            ███ ░███  ░███                        
                                           ░░██████   █████                       
                                            ░░░░░░   ░░░░░                        
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

	if m.menuState == MenuMain {
		// Create gradient colors from light blue to dark blue
		colors := []string{
			"#87CEEB", // Sky blue
			"#4682B4", // Steel blue
			"#1E90FF", // Dodger blue
			"#0000CD", // Medium blue
			"#000080", // Navy blue
		}

		// Apply gradient to each line
		var logo strings.Builder
		lines := strings.Split(logoArt, "\n")
		for i, line := range lines {
			if line == "" {
				logo.WriteString("\n")
				continue
			}
			colorIndex := i % len(colors)
			style := lipgloss.NewStyle().Foreground(lipgloss.Color(colors[colorIndex]))
			logo.WriteString(style.Render(line))
			logo.WriteString("\n")
		}

		menuContent = fmt.Sprintf("%s\n%s", logo.String(), m.renderMainMenu())
	} else if m.menuState == MenuSettings {
		menuContent = m.renderSettingsMenu()
	}

	footer := "\n" + HelpStyle("↑/↓: Navigate • Enter: Select • Esc: Back • q: Quit")

	content := fmt.Sprintf("%s\n%s", menuContent, footer)

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
		Margin(1, 0, 2, 0).
		PaddingLeft(2)

	sb.WriteString(titleStyle.Render("Settings"))
	sb.WriteString("\n\n")

	// Create left column for settings items
	var settingsList []string
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

		settingsList = append(settingsList, s.Render(menuText))
	}

	// Create right column for the example
	var exampleBox string
	if m.selectedItem < len(m.settingsItems) {
		switch m.selectedItem {
		case 0:
			exampleBox = renderThemeExample(m.selectedTheme)
		case 1:
			exampleBox = renderCursorExample(m.cursorType)
		case 2:
			exampleBox = renderGameModeExample(m.gameMode)
		case 3:
			exampleBox = renderUseNumbersExample(m.useNumbers)
		}
	}

	// Calculate column widths
	leftWidth := 30                       // Fixed width for settings column
	rightWidth := m.width - leftWidth - 4 // Remaining width for example box

	// Style the example box
	exampleStyle := lipgloss.NewStyle().
		Padding(1, 2).
		Width(rightWidth)

	// Create columns
	leftColumn := lipgloss.NewStyle().
		Width(leftWidth).
		Render(lipgloss.JoinVertical(lipgloss.Left, settingsList...))

	rightColumn := exampleStyle.Render(exampleBox)

	// Join columns
	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftColumn,
		"  ",
		rightColumn,
	)

	sb.WriteString(content)
	sb.WriteString("\n\n")

	return sb.String()
}

func renderThemeExample(theme string) string {
	var sb strings.Builder

	// Get the current theme colors
	colors := map[string]lipgloss.Color{
		"Help Text":    GetColor("help_text"),
		"Timer":        GetColor("timer"),
		"Border":       GetColor("border"),
		"Text Dim":     GetColor("text_dim"),
		"Text Preview": GetColor("text_preview"),
		"Text Correct": GetColor("text_correct"),
		"Text Error":   GetColor("text_error"),
		"Text Partial": GetColor("text_partial_error"),
		"Cursor FG":    GetColor("cursor_fg"),
		"Cursor BG":    GetColor("cursor_bg"),
		"Cursor Under": GetColor("cursor_underline"),
		"Padding":      GetColor("padding"),
	}

	// Create color preview boxes
	for name, color := range colors {
		style := lipgloss.NewStyle().
			Foreground(color).
			Padding(0, 1)

		// Create a color box
		colorBox := lipgloss.NewStyle().
			Background(color).
			Padding(0, 2).
			Render("  ")

		// Combine color name and box
		sb.WriteString(fmt.Sprintf("%-15s %s %s\n", name, colorBox, style.Render(string(color))))
	}

	return sb.String()
}

func renderCursorExample(cursorType string) string {
	var example strings.Builder

	// Title with timer color
	titleStyle := lipgloss.NewStyle().Foreground(GetColor("timer")).Bold(true)
	example.WriteString(titleStyle.Render("Cursor Style: "))

	// Cursor type with text_preview color
	cursorTypeStyle := lipgloss.NewStyle().Foreground(GetColor("text_preview"))
	example.WriteString(cursorTypeStyle.Render(cursorType))
	example.WriteString("\n\n")

	// Example section
	example.WriteString(titleStyle.Render("Example Text:\n"))

	// Show example with appropriate styling
	dimStyle := lipgloss.NewStyle().Foreground(GetColor("text_dim"))
	example.WriteString(dimStyle.Render("quick "))

	// Render "brown" with cursor on 'r'
	example.WriteString(dimStyle.Render("b"))
	if cursorType == "block" {
		cursorStyle := lipgloss.NewStyle().
			Foreground(GetColor("cursor_fg")).
			Background(GetColor("cursor_bg"))
		example.WriteString(cursorStyle.Render("r"))
	} else {
		cursorStyle := lipgloss.NewStyle().
			Foreground(GetColor("cursor_underline")).
			Underline(true)
		example.WriteString(cursorStyle.Render("r"))
	}
	example.WriteString(dimStyle.Render("own"))

	return example.String()
}

func renderGameModeExample(gameMode string) string {
	var example strings.Builder

	// Title with timer color
	titleStyle := lipgloss.NewStyle().Foreground(GetColor("timer")).Bold(true)
	example.WriteString(titleStyle.Render("Game Mode: "))

	// Mode name with text_preview color
	modeStyle := lipgloss.NewStyle().Foreground(GetColor("text_preview"))
	if gameMode == "normal" {
		example.WriteString(modeStyle.Render("Normal (With Punctuation)"))
	} else {
		example.WriteString(modeStyle.Render("Simple (No Punctuation)"))
	}
	example.WriteString("\n\n")

	// Example section
	example.WriteString(titleStyle.Render("Example:\n"))

	// Show example with appropriate styling
	if gameMode == "normal" {
		example.WriteString(TextToTypeStyle.Render("The quick brown fox jumps."))
	} else {
		example.WriteString(TextToTypeStyle.Render("the quick brown fox jumps"))
	}

	return example.String()
}

func renderUseNumbersExample(useNumbers bool) string {
	var example strings.Builder

	// Title with timer color
	titleStyle := lipgloss.NewStyle().Foreground(GetColor("timer")).Bold(true)
	example.WriteString(titleStyle.Render("Use Numbers: "))

	// Yes/No with text_preview color
	valueStyle := lipgloss.NewStyle().Foreground(GetColor("text_preview"))
	if useNumbers {
		example.WriteString(valueStyle.Render("Yes"))
	} else {
		example.WriteString(valueStyle.Render("No"))
	}
	example.WriteString("\n\n")

	// Example section
	example.WriteString(titleStyle.Render("Example:\n"))

	// Show example with appropriate styling
	if useNumbers {
		example.WriteString(TextToTypeStyle.Render("quick brown fox jumps over 5 lazy dogs"))
	} else {
		example.WriteString(TextToTypeStyle.Render("quick brown fox jumps over lazy dogs"))
	}

	return example.String()
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
