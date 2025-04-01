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
	initialTheme    string // Store the initial theme before settings changes
	availableThemes []string
	themeChanged    bool   // Track if theme was changed in settings
	gameMode        string // Game mode (normal or simple)
	useNumbers      bool   // Whether to include numbers in the text
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

	// Create main menu items
	model.mainMenuItems = []menuItem{
		{title: "Start Typing", action: startGame},
		{title: "Multiplayer Typeracer", action: nil, disabled: true, backColor: "#555555"},
		{title: "Settings", action: openSettings},
		{title: "Statistics", action: openStats, disabled: true, backColor: "#555555"},
		{title: "Quit", action: quitGame},
	}

	// Create settings menu items
	model.settingsItems = []menuItem{
		{title: "Theme", action: cycleTheme},
		{title: "Cursor Style", action: cycleCursor},
		{title: "Game Mode", action: cycleGameMode},
		{title: "Use Numbers", action: toggleNumbers},
		{title: "Back", action: saveAndGoBack},
	}

	// Make sure initial selection is on an enabled item
	if model.mainMenuItems[model.selectedItem].disabled {
		// Find first enabled item
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
			// Move up, skip disabled items
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

				// Check if this item is enabled
				if m.menuState == MenuMain {
					if !m.mainMenuItems[m.selectedItem].disabled || m.selectedItem == previousItem {
						break // Found an enabled item or looped through all items
					}
				} else {
					if !m.settingsItems[m.selectedItem].disabled || m.selectedItem == previousItem {
						break // Found an enabled item or looped through all items
					}
				}
			}

		case "down", "j":
			// Move down, skip disabled items
			previousItem := m.selectedItem
			for {
				m.selectedItem++

				if m.menuState == MenuMain {
					if m.selectedItem >= len(m.mainMenuItems) {
						m.selectedItem = 0
					}

					if !m.mainMenuItems[m.selectedItem].disabled || m.selectedItem == previousItem {
						break // Found an enabled item or looped through all items
					}
				} else {
					if m.selectedItem >= len(m.settingsItems) {
						m.selectedItem = 0
					}

					if !m.settingsItems[m.selectedItem].disabled || m.selectedItem == previousItem {
						break // Found an enabled item or looped through all items
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
				// Restore original theme if we're exiting settings without saving
				if m.themeChanged {
					LoadTheme(m.initialTheme)
					UpdateStyles()
					m.selectedTheme = m.initialTheme
					m.themeChanged = false
				}

				// Revert to initial cursor type as well
				m.cursorType = CurrentSettings.CursorType

				m.menuState = MenuMain
				m.selectedItem = 0

				// Skip disabled items
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
	// Determine which menu to display
	var menuContent string

	// Format the logo
	logoStyle := lipgloss.NewStyle().
		Foreground(GetColor("border")).
		Bold(true)

	logo := logoStyle.Render(logoArt)

	// Main content area
	if m.menuState == MenuMain {
		menuContent = m.renderMainMenu()
	} else if m.menuState == MenuSettings {
		menuContent = m.renderSettingsMenu()
	}

	footer := "\n" + HelpStyle("↑/↓: Navigate • Enter: Select • Esc: Back • q: Quit")

	// Combine all elements
	content := fmt.Sprintf("%s\n%s\n%s", logo, menuContent, footer)

	// Center everything on screen
	if m.width > 0 && m.height > 0 {
		return lipgloss.Place(m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			content)
	}

	return content
}

func (m *StartScreenModel) renderMainMenu() string {
	var sb strings.Builder

	// Title
	titleStyle := lipgloss.NewStyle().
		Foreground(GetColor("timer")).
		Bold(true).
		Margin(1, 0, 2, 0)

	sb.WriteString(titleStyle.Render("Main Menu"))
	sb.WriteString("\n\n")

	// Menu items
	for i, item := range m.mainMenuItems {
		var s lipgloss.Style

		if i == m.selectedItem {
			// Selected item - use underline with accent color
			s = lipgloss.NewStyle().
				Foreground(GetColor("cursor_bg")).
				Bold(true).
				Padding(0, 4).
				Underline(true)
		} else if item.disabled {
			// Disabled item
			c := GetColor("text_dim")
			if item.backColor != "" {
				c = lipgloss.Color(item.backColor)
			}
			s = lipgloss.NewStyle().
				Foreground(c).
				Padding(0, 4)
		} else {
			// Normal item
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

	// Title
	titleStyle := lipgloss.NewStyle().
		Foreground(GetColor("timer")).
		Bold(true).
		Margin(1, 0, 2, 0)

	sb.WriteString(titleStyle.Render("Settings"))
	sb.WriteString("\n\n")

	// Menu items
	for i, item := range m.settingsItems {
		var s lipgloss.Style

		if i == m.selectedItem {
			// Selected item - use underline with accent color
			s = lipgloss.NewStyle().
				Foreground(GetColor("cursor_bg")).
				Bold(true).
				Padding(0, 2).
				Underline(true)
		} else {
			// Normal item
			s = lipgloss.NewStyle().
				Foreground(GetColor("text_preview")).
				Padding(0, 2)
		}

		menuText := item.title

		// Show current value for settings
		if i == 0 { // Theme setting
			menuText = fmt.Sprintf("%-15s: %s", item.title, m.selectedTheme)
		} else if i == 1 { // Cursor setting
			menuText = fmt.Sprintf("%-15s: %s", item.title, m.cursorType)
		} else if i == 2 { // Game Mode setting
			menuText = fmt.Sprintf("%-15s: %s", item.title, m.gameMode)
		} else if i == 3 { // Use Numbers setting
			menuText = fmt.Sprintf("%-15s: %v", item.title, m.useNumbers)
		}

		sb.WriteString(s.Render(menuText))
		sb.WriteString("\n")

		// Show example for selected item
		if i == m.selectedItem {
			sb.WriteString("\n")

			// Show examples for the different settings
			if i == 0 { // Theme example
				exampleBox := renderThemeExample(m.selectedTheme)
				sb.WriteString(exampleBox)
			} else if i == 1 { // Cursor example
				exampleBox := renderCursorExample(m.cursorType)
				sb.WriteString(exampleBox)
			} else if i == 2 { // Game Mode example
				exampleBox := renderGameModeExample(m.gameMode)
				sb.WriteString(exampleBox)
			} else if i == 3 { // Use Numbers example
				exampleBox := renderUseNumbersExample(m.useNumbers)
				sb.WriteString(exampleBox)
			}

			sb.WriteString("\n")
		}

		sb.WriteString("\n")
	}

	return sb.String()
}

// Render a theme example showing colors for typing
func renderThemeExample(themeName string) string {
	exampleStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(GetColor("border")).
		Padding(1, 2).
		Margin(0, 0, 0, 4)

	var example strings.Builder

	// Show a mini example of the theme's colors
	example.WriteString("Text Preview: ")
	example.WriteString(DimStyle.Render("dim "))
	example.WriteString(InputStyle.Render("correct "))
	example.WriteString(ErrorStyle.Render("error"))
	example.WriteString("\n\n")

	// Show a word with errors
	example.WriteString("Word with errors: ")
	example.WriteString(InputStyle.Render("co"))
	example.WriteString(ErrorStyle.Render("d"))
	example.WriteString(PartialErrorStyle.Render("e"))
	example.WriteString(DimStyle.Render("r"))

	return exampleStyle.Render(example.String())
}

// Render a cursor example
func renderCursorExample(cursorType string) string {
	exampleStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(GetColor("border")).
		Padding(1, 2).
		Margin(0, 0, 0, 4)

	var example strings.Builder

	// Create a temporary cursor
	var cursor *Cursor
	if cursorType == "block" {
		cursor = NewCursor(BlockCursor)
	} else {
		cursor = NewCursor(UnderlineCursor)
	}

	// Show the cursor
	example.WriteString("Cursor appearance: ")
	example.WriteString(cursor.Render('A'))
	example.WriteString(InputStyle.Render("BC"))

	return exampleStyle.Render(example.String())
}

// Render a game mode example
func renderGameModeExample(gameMode string) string {
	exampleStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(GetColor("border")).
		Padding(1, 2).
		Margin(0, 0, 0, 4)

	var example strings.Builder

	// Show game mode description
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

// Render a use numbers example
func renderUseNumbersExample(useNumbers bool) string {
	exampleStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(GetColor("border")).
		Padding(1, 2).
		Margin(0, 0, 0, 4)

	var example strings.Builder

	// Show numbers setting description
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

// Action functions for menu items
func startGame(m *StartScreenModel) tea.Cmd {
	return tea.Quit
}

func openSettings(m *StartScreenModel) tea.Cmd {
	// Store initial theme when entering settings
	m.initialTheme = m.selectedTheme
	m.themeChanged = false

	m.menuState = MenuSettings
	m.selectedItem = 0
	return nil
}

func openStats(m *StartScreenModel) tea.Cmd {
	// This is disabled
	return nil
}

func quitGame(m *StartScreenModel) tea.Cmd {
	return tea.Quit
}

func saveAndGoBack(m *StartScreenModel) tea.Cmd {
	// When accepting changes, commit them
	m.initialTheme = m.selectedTheme
	m.themeChanged = false

	// Save settings to disk
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
	// Find current theme in the list
	currentIndex := -1
	for i, theme := range m.availableThemes {
		if theme == m.selectedTheme {
			currentIndex = i
			break
		}
	}

	// Move to next theme
	currentIndex = (currentIndex + 1) % len(m.availableThemes)
	m.selectedTheme = m.availableThemes[currentIndex]
	m.themeChanged = true

	// Update theme preview
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

// Message to start the game
type StartGameMsg struct {
	cursorType string
	theme      string
}

// RunStartScreen runs the start screen and returns when the user selects an option
func RunStartScreen() {
	// Initialize and run the program
	p := tea.NewProgram(NewStartScreenModel(), tea.WithAltScreen())

	// Run the program
	model, err := p.Run()
	if err != nil {
		fmt.Println("Error running start screen:", err)
		return
	}

	// Check if we need to start the game
	if m, ok := model.(*StartScreenModel); ok {
		// Apply and save the selected settings
		UpdateSettings(UserSettings{
			ThemeName:  m.selectedTheme,
			CursorType: m.cursorType,
			GameMode:   m.gameMode,
			UseNumbers: m.useNumbers,
		})

		// Check which item was selected
		if m.menuState == MenuMain && m.selectedItem < len(m.mainMenuItems) {
			item := m.mainMenuItems[m.selectedItem]

			// If Start Typing was selected
			if item.title == "Start Typing" {
				// Start game with the selected settings
				StartLoadingWithOptions(m.cursorType)
			}
		}
	}
}
