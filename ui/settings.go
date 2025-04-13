package ui

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
	"path/filepath"
)

type UserSettings struct {
	ThemeName      string `json:"theme"`
	CursorType     string `json:"cursor_type"`
	GameMode       string `json:"game_mode"`
	UseNumbers     bool   `json:"use_numbers"`
	TextLength     string `json:"text_length"`
	HasSeenWelcome bool   `json:"has_seen_welcome"`
	RefreshRate    int    `json:"refresh_rate"` // NOTE:in frames per second not tick
}

const (
	GameModeNormal = "normal"
	GameModeSimple = "simple"

	TextLengthShort    = "short"     // 1
	TextLengthMedium   = "medium"    // 2
	TextLengthLong     = "long"      // 3
	TextLengthVeryLong = "very long" // 5
)

var DefaultSettings = UserSettings{
	ThemeName:      "default",
	CursorType:     "block",
	GameMode:       GameModeNormal,
	UseNumbers:     true,
	TextLength:     TextLengthShort,
	HasSeenWelcome: false,
	RefreshRate:    10,
}

var CurrentSettings UserSettings

func GetConfigDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user config directory: %w", err)
	}

	appConfigDir := filepath.Join(configDir, "go-typer")

	if err := os.MkdirAll(appConfigDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	return appConfigDir, nil
}

func GetSettingsFilePath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "settings.json"), nil
}

func LoadSettings() error {
	CurrentSettings = DefaultSettings

	settingsPath, err := GetSettingsFilePath()
	if err != nil {
		return fmt.Errorf("failed to get settings file path: %w", err)
	}

	data, err := os.ReadFile(settingsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return SaveSettings()
		}
		return fmt.Errorf("error reading settings file: %w", err)
	}

	if err := json.Unmarshal(data, &CurrentSettings); err != nil {
		return fmt.Errorf("error parsing settings file: %w", err)
	}

	ApplySettings()

	return nil
}

func SaveSettings() error {
	settingsPath, err := GetSettingsFilePath()
	if err != nil {
		return fmt.Errorf("failed to get settings file path: %w", err)
	}

	data, err := json.MarshalIndent(CurrentSettings, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling settings: %w", err)
	}

	if err := os.WriteFile(settingsPath, data, 0644); err != nil {
		return fmt.Errorf("error writing settings file: %w", err)
	}

	return nil
}

func UpdateSettings(settings UserSettings) error {
	if settings.ThemeName != "" {
		CurrentSettings.ThemeName = settings.ThemeName
	}

	if settings.CursorType != "" {
		CurrentSettings.CursorType = settings.CursorType
	}

	if settings.GameMode != "" {
		CurrentSettings.GameMode = settings.GameMode
	}

	if settings.UseNumbers != CurrentSettings.UseNumbers {
		CurrentSettings.UseNumbers = settings.UseNumbers
	}

	if settings.TextLength != "" {
		CurrentSettings.TextLength = settings.TextLength
	}

	if settings.RefreshRate > 0 {
		CurrentSettings.RefreshRate = settings.RefreshRate
	}

	ApplySettings()

	return SaveSettings()
}

func ApplySettings() {
	if CurrentSettings.ThemeName != "" {
		LoadTheme(CurrentSettings.ThemeName)
		UpdateStyles()
	}

	if CurrentSettings.CursorType == "underline" {
		DefaultCursorType = UnderlineCursor
	} else {
		DefaultCursorType = BlockCursor
	}
}

func InitSettings() {
	if err := LoadSettings(); err != nil {
		fmt.Printf("Warning: Could not load settings: %v\n", err)
		fmt.Println("Using default settings")

		CurrentSettings = DefaultSettings
		ApplySettings()
	}
}

type SettingsItem struct {
	title    string
	options  []string
	details  string
	selected int
	expanded bool
	key      string
}

func (i SettingsItem) Title() string {
	arrow := "→"
	if i.expanded {
		arrow = "↓"
	}
	return fmt.Sprintf("%s %s", arrow, i.title)
}

func (i SettingsItem) Description() string {
	return fmt.Sprintf("%s: %s", i.options[i.selected], i.details)
}

func (i SettingsItem) FilterValue() string { return i.title }

type SettingsModel struct {
	list          list.Model
	height, width int
	settings      UserSettings
}

func createSettingsItems(settings UserSettings) []list.Item {
	themeOptions := []string{"default", "dark", "light"}
	themeSelected := 0
	for i, opt := range themeOptions {
		if opt == settings.ThemeName {
			themeSelected = i
			break
		}
	}

	cursorOptions := []string{"block", "underline"}
	cursorSelected := 0
	for i, opt := range cursorOptions {
		if opt == settings.CursorType {
			cursorSelected = i
			break
		}
	}

	gameModeOptions := []string{GameModeNormal, GameModeSimple}
	gameModeSelected := 0
	for i, opt := range gameModeOptions {
		if opt == settings.GameMode {
			gameModeSelected = i
			break
		}
	}

	textLengthOptions := []string{TextLengthShort, TextLengthMedium, TextLengthLong, TextLengthVeryLong}
	textLengthSelected := 0
	for i, opt := range textLengthOptions {
		if opt == settings.TextLength {
			textLengthSelected = i
			break
		}
	}

	refreshRateOptions := []string{"5", "10", "15", "20", "30"}
	refreshRateSelected := 0
	for i, opt := range refreshRateOptions {
		if opt == fmt.Sprintf("%d", settings.RefreshRate) {
			refreshRateSelected = i
			break
		}
	}

	return []list.Item{
		&SettingsItem{
			title:    "Theme",
			options:  themeOptions,
			details:  "Select your preferred theme",
			selected: themeSelected,
			key:      "theme",
		},
		&SettingsItem{
			title:    "Cursor Type",
			options:  cursorOptions,
			details:  "Choose cursor appearance",
			selected: cursorSelected,
			key:      "cursor",
		},
		&SettingsItem{
			title:    "Game Mode",
			options:  gameModeOptions,
			details:  "Select game difficulty mode",
			selected: gameModeSelected,
			key:      "game_mode",
		},
		&SettingsItem{
			title:    "Text Length",
			options:  textLengthOptions,
			details:  "Choose text length for typing",
			selected: textLengthSelected,
			key:      "text_length",
		},
		&SettingsItem{
			title:    "Refresh Rate",
			options:  refreshRateOptions,
			details:  "Set UI refresh rate (FPS)",
			selected: refreshRateSelected,
			key:      "refresh_rate",
		},
	}
}

func initialSettingsModel() SettingsModel {
	settings := CurrentSettings
	items := createSettingsItems(settings)

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)
	l.SetShowStatusBar(false)
	l.Title = "Settings"
	l.Styles.Title = SettingsTitleStyle

	return SettingsModel{
		list:     l,
		settings: settings,
	}
}

func (m SettingsModel) Init() tea.Cmd { return nil }

func (m SettingsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.list.SetSize(m.width/3, m.height/2)

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			if i, ok := m.list.SelectedItem().(*SettingsItem); ok {
				if !i.expanded {
					i.expanded = true
				} else {
					i.selected = (i.selected + 1) % len(i.options)
					switch i.key {
					case "theme":
						m.settings.ThemeName = i.options[i.selected]
					case "cursor":
						m.settings.CursorType = i.options[i.selected]
					case "game_mode":
						m.settings.GameMode = i.options[i.selected]
					case "text_length":
						m.settings.TextLength = i.options[i.selected]
					case "refresh_rate":
						fmt.Sscanf(i.options[i.selected], "%d", &m.settings.RefreshRate)
					}
					UpdateSettings(m.settings)
				}
			}
		case "esc":
			if i, ok := m.list.SelectedItem().(*SettingsItem); ok {
				i.expanded = false
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m SettingsModel) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	listView := SettingsListStyle.Render(m.list.View())

	details := "Select an item to view details"
	if i, ok := m.list.SelectedItem().(*SettingsItem); ok {
		if i.expanded {
			details = fmt.Sprintf("%s\n\nCurrent: %s\n\nOptions:\n",
				i.details,
				i.options[i.selected],
			)
			for idx, opt := range i.options {
				bullet := "•"
				if idx == i.selected {
					bullet = ">"
				}
				details += fmt.Sprintf("%s %s\n", bullet, opt)
			}
		} else {
			details = i.details
		}
	}

	detailsView := SettingsDetailsStyle.Render(details)

	content := lipgloss.Place(
		m.width,
		m.height-1,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinHorizontal(lipgloss.Top, listView, detailsView),
	)

	help := lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Center,
		SettingsHelpStyle.Render("↑/↓: Navigate • Enter: Select • Esc: Back • q: Quit"),
	)

	return lipgloss.JoinVertical(lipgloss.Bottom, content, help)
}

func ShowSettings() error {
	p := tea.NewProgram(initialSettingsModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running settings program: %w", err)
	}
	return nil
}
