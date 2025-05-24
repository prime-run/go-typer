package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	devlog "github.com/prime-run/go-typer/log"
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
	width           int
	height          int
	menuState       int
	selectedItem    int
	mainMenuItems   []menuItem
	settingsItems   []menuItem
	cursorType      string
	selectedTheme   string
	initialTheme    string
	availableThemes []string
	themeChanged    bool
	gameMode        string
	useNumbers      bool
	textLength      string
	refreshRate     int
	startTime       time.Time
	lastTick        time.Time
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
		textLength:      CurrentSettings.TextLength,
		refreshRate:     CurrentSettings.RefreshRate,
		startTime:       time.Now(),
		lastTick:        time.Now(),
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
		{title: "Text Length", action: cycleTextLength},
		{title: "Refresh Rate", action: cycleRefreshRate},
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
	return InitGlobalTick()
}

func (m *StartScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		return m, nil
	}

	return m, nil
}

func (m *StartScreenModel) View() string {
	var menuContent string

	if m.menuState == MenuMain {
		menuContent = fmt.Sprintf("%s\n%s",
			renderAnimatedAscii(logoArt, m.lastTick),
			m.renderMainMenu())
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
		Margin(1, 0, 2, 0)

	sb.WriteString(titleStyle.Render("Settings"))
	sb.WriteString("\n\n")

	var exampleContent string

	if m.selectedItem == 0 {
		exampleContent = renderThemeExample(m.selectedTheme)
	} else if m.selectedItem == 1 {
		exampleContent = renderCursorExample(m.cursorType)
	} else if m.selectedItem == 2 {
		exampleContent = renderGameModeExample(m.gameMode)
	} else if m.selectedItem == 3 {
		exampleContent = renderUseNumbersExample(m.useNumbers)
	} else if m.selectedItem == 4 {
		exampleContent = renderTextLengthExample(m.textLength)
	} else if m.selectedItem == 5 {
		exampleContent = renderRefreshRateExample(m.refreshRate, m.lastTick)
	}

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
		} else if i == 4 {
			menuText = fmt.Sprintf("%-15s: %s", item.title, m.textLength)
		} else if i == 5 {
			menuText = fmt.Sprintf("%-15s: %d FPS", item.title, m.refreshRate)
		}

		settingsList = append(settingsList, s.Render(menuText))
	}

	var exampleBox string
	if m.selectedItem < len(m.settingsItems) {
		switch m.selectedItem {
		case 0:
			exampleBox = exampleContent
		case 1:
			exampleBox = renderCursorExample(m.cursorType)
		case 2:
			exampleBox = renderGameModeExample(m.gameMode)
		case 3:
			exampleBox = renderUseNumbersExample(m.useNumbers)
		case 4:
			exampleBox = renderTextLengthExample(m.textLength)
		case 5:
			exampleBox = renderRefreshRateExample(m.refreshRate, m.lastTick)
		}
	}

	leftWidth := 30
	rightWidth := m.width - leftWidth - 4

	exampleStyle := lipgloss.NewStyle().
		Padding(1, 2).
		Width(rightWidth)

	// create columns
	leftColumn := lipgloss.NewStyle().
		Width(leftWidth).
		Render(lipgloss.JoinVertical(lipgloss.Left, settingsList...))

	rightColumn := exampleStyle.Render(exampleBox)

	// join them
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

	colorOrder := []string{
		"Help Text",
		"Timer",
		"Border",
		"Text Dim",
		"Text Preview",
		"Text Correct",
		"Text Error",
		"Text Partial",
		"Cursor FG",
		"Cursor BG",
		"Cursor Under",
		"Padding",
	}

	for _, name := range colorOrder {
		color := GetColor(strings.ToLower(strings.ReplaceAll(name, " ", "_")))
		style := lipgloss.NewStyle().
			Foreground(color).
			Padding(0, 1)

		colorBox := lipgloss.NewStyle().
			Background(color).
			Padding(0, 2).
			Render("  ")

		sb.WriteString(fmt.Sprintf("%-15s %s %s\n", name, colorBox, style.Render(string(color))))
	}

	return sb.String()
}

func renderCursorExample(cursorType string) string {
	var example strings.Builder

	titleStyle := lipgloss.NewStyle().Foreground(GetColor("timer")).Bold(true)
	example.WriteString(titleStyle.Render("Cursor Style: "))

	cursorTypeStyle := lipgloss.NewStyle().Foreground(GetColor("text_preview"))
	example.WriteString(cursorTypeStyle.Render(cursorType))
	example.WriteString("\n\n")

	example.WriteString(titleStyle.Render("Example Text:\n"))

	dimStyle := lipgloss.NewStyle().Foreground(GetColor("text_dim"))
	example.WriteString(dimStyle.Render("quick "))

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

	titleStyle := lipgloss.NewStyle().Foreground(GetColor("timer")).Bold(true)
	example.WriteString(titleStyle.Render("Game Mode: "))

	modeStyle := lipgloss.NewStyle().Foreground(GetColor("text_preview"))
	if gameMode == "normal" {
		example.WriteString(modeStyle.Render("Normal (With Punctuation)"))
	} else {
		example.WriteString(modeStyle.Render("Simple (No Punctuation)"))
	}
	example.WriteString("\n\n")

	example.WriteString(titleStyle.Render("Example:\n"))

	if gameMode == "normal" {
		example.WriteString(TextToTypeStyle.Render("The quick brown fox jumps."))
	} else {
		example.WriteString(TextToTypeStyle.Render("the quick brown fox jumps"))
	}

	return example.String()
}

func renderUseNumbersExample(useNumbers bool) string {
	var example strings.Builder

	titleStyle := lipgloss.NewStyle().Foreground(GetColor("timer")).Bold(true)
	example.WriteString(titleStyle.Render("Use Numbers: "))

	valueStyle := lipgloss.NewStyle().Foreground(GetColor("text_preview"))
	if useNumbers {
		example.WriteString(valueStyle.Render("Yes"))
	} else {
		example.WriteString(valueStyle.Render("No"))
	}
	example.WriteString("\n\n")

	example.WriteString(titleStyle.Render("Example:\n"))

	if useNumbers {
		example.WriteString(TextToTypeStyle.Render("quick brown fox jumps over 5 lazy dogs"))
	} else {
		example.WriteString(TextToTypeStyle.Render("quick brown fox jumps over lazy dogs"))
	}

	return example.String()
}

func renderTextLengthExample(length string) string {
	var example strings.Builder

	titleStyle := lipgloss.NewStyle().Foreground(GetColor("timer")).Bold(true)
	example.WriteString(titleStyle.Render("Text Length: "))

	valueStyle := lipgloss.NewStyle().Foreground(GetColor("text_preview"))
	example.WriteString(valueStyle.Render(length))
	example.WriteString("\n\n")

	example.WriteString(titleStyle.Render("Quotes to fetch:\n"))

	textCount := map[string]int{
		TextLengthShort:    1,
		TextLengthMedium:   2,
		TextLengthLong:     3,
		TextLengthVeryLong: 5,
	}

	count := textCount[length]
	example.WriteString(fmt.Sprintf("\nWill fetch and combine %d quote(s)", count))
	example.WriteString("\nEstimated word count: ")

	wordCount := count * 30
	example.WriteString(valueStyle.Render(fmt.Sprintf("%d words", wordCount)))

	return example.String()
}

func renderRefreshRateExample(rate int, tickTime time.Time) string {
	var sb strings.Builder

	titleStyle := lipgloss.NewStyle().
		Foreground(GetColor("timer")).
		Bold(true)

	sb.WriteString(titleStyle.Render("Refresh Rate: "))

	valueStyle := lipgloss.NewStyle().
		Foreground(GetColor("text_preview"))
	sb.WriteString(valueStyle.Render(fmt.Sprintf("%d FPS", rate)))
	sb.WriteString("\n\n")

	descStyle := lipgloss.NewStyle().
		Foreground(GetColor("text_dim"))

	sb.WriteString(descStyle.Render(
		fmt.Sprintf("Updates %d times per second (%.1f ms per frame)",
			rate, 1000.0/float64(rate))))
	sb.WriteString("\n\n")

	helpStyle := lipgloss.NewStyle().
		Foreground(GetColor("help_text"))

	sb.WriteString(helpStyle.Render(
		"Higher values give smoother animations\n" +
			"Lower values use less CPU/battery"))
	sb.WriteString("\n\n")

	var spinner string
	frames := []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}
	index := int(tickTime.UnixNano()/int64(time.Second/time.Duration(rate))) % len(frames)
	spinner = frames[index]

	spinnerStyle := lipgloss.NewStyle().
		Foreground(GetColor("text_correct"))

	sb.WriteString(spinnerStyle.Render(spinner + " "))
	sb.WriteString(valueStyle.Render(fmt.Sprintf("Animation at %d FPS", rate)))

	return sb.String()
}

func renderAnimatedAscii(logoArt string, tickTime time.Time) string {
	var result strings.Builder
	colors := []string{
		"#87CEEB", // Sky blue
		"#4682B4", // Steel blue
		"#1E90FF", // Dodger blue
		"#0000CD", // Medium blue
		"#000080", // Navy blue
	}

	startIndex := int(tickTime.UnixNano()/int64(100*time.Millisecond)) % len(colors)

	lines := strings.Split(logoArt, "\n")
	for i, line := range lines {
		if line == "" {
			result.WriteString("\n")
			continue
		}
		colorIndex := (startIndex + i) % len(colors)
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(colors[colorIndex]))
		result.WriteString(style.Render(line))
		result.WriteString("\n")
	}

	return result.String()
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
	settings := UserSettings{
		ThemeName:      m.selectedTheme,
		CursorType:     m.cursorType,
		GameMode:       m.gameMode,
		UseNumbers:     m.useNumbers,
		TextLength:     m.textLength,
		RefreshRate:    m.refreshRate,
		HasSeenWelcome: CurrentSettings.HasSeenWelcome,
	}

	if err := UpdateSettings(settings); err != nil {
		devlog.Log("Settings: Error updating settings: %v", err)
	}

	m.menuState = MenuMain
	m.selectedItem = 0

	for m.mainMenuItems[m.selectedItem].disabled {
		m.selectedItem++
		if m.selectedItem >= len(m.mainMenuItems) {
			m.selectedItem = 0
		}
	}

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

func cycleTextLength(m *StartScreenModel) tea.Cmd {
	lengths := []string{TextLengthShort, TextLengthMedium, TextLengthLong, TextLengthVeryLong}
	var currentIndex int

	for i, length := range lengths {
		if length == m.textLength {
			currentIndex = i
			break
		}
	}

	currentIndex = (currentIndex + 1) % len(lengths)
	m.textLength = lengths[currentIndex]

	return nil
}

func cycleRefreshRate(m *StartScreenModel) tea.Cmd {
	rates := []int{1, 5, 10, 15, 30, 60}

	currentIndex := -1
	for i, r := range rates {
		if r == m.refreshRate {
			currentIndex = i
			break
		}
	}

	currentIndex = (currentIndex + 1) % len(rates)
	m.refreshRate = rates[currentIndex]

	return nil
}

type StartGameMsg struct {
	cursorType string
	theme      string
}

func RunStartScreen() {
	ShowWelcomeScreen()

	p := tea.NewProgram(NewStartScreenModel(), tea.WithAltScreen())

	model, err := p.Run()
	if err != nil {
		fmt.Printf("Error running start screen: %v\n", err)
		return
	}

	if m, ok := model.(*StartScreenModel); ok {
		UpdateSettings(UserSettings{
			ThemeName:      m.selectedTheme,
			CursorType:     m.cursorType,
			GameMode:       m.gameMode,
			UseNumbers:     m.useNumbers,
			TextLength:     m.textLength,
			RefreshRate:    m.refreshRate,
			HasSeenWelcome: CurrentSettings.HasSeenWelcome,
		})

		if m.menuState == MenuMain && m.selectedItem < len(m.mainMenuItems) {
			item := m.mainMenuItems[m.selectedItem]

			if item.title == "Start Typing" {
				StartLoadingWithOptions(m.cursorType)
			}
		}
	}
}
