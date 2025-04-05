package ui

import (
	"encoding/json"
	"fmt"
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
}

const (
	GameModeNormal = "normal"
	GameModeSimple = "simple"

	TextLengthShort    = "short"     // 1 quote
	TextLengthMedium   = "medium"    // 2 quotes
	TextLengthLong     = "long"      // 3 quotes
	TextLengthVeryLong = "very long" // 5 quotes
)

var DefaultSettings = UserSettings{
	ThemeName:      "default",
	CursorType:     "block",
	GameMode:       GameModeNormal,
	UseNumbers:     true,
	TextLength:     TextLengthShort,
	HasSeenWelcome: false,
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
