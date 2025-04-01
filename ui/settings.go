package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// UserSettings holds user preferences that should be persisted
type UserSettings struct {
	ThemeName  string `json:"theme"`
	CursorType string `json:"cursor_type"`
	GameMode   string `json:"game_mode"`
	UseNumbers bool   `json:"use_numbers"`
}

// Game modes
const (
	GameModeNormal = "normal"
	GameModeSimple = "simple" // No punctuation, lowercase only
)

// DefaultSettings provides default values for user settings
var DefaultSettings = UserSettings{
	ThemeName:  "default",
	CursorType: "block",
	GameMode:   GameModeNormal,
	UseNumbers: true,
}

// CurrentSettings holds the active user settings
var CurrentSettings UserSettings

// GetConfigDir returns the path to the application's configuration directory
func GetConfigDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user config directory: %w", err)
	}

	appConfigDir := filepath.Join(configDir, "go-typer")

	// Ensure the directory exists
	if err := os.MkdirAll(appConfigDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	return appConfigDir, nil
}

// GetSettingsFilePath returns the full path to the settings file
func GetSettingsFilePath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "settings.json"), nil
}

// LoadSettings loads user settings from the config file
func LoadSettings() error {
	// Start with default settings
	CurrentSettings = DefaultSettings

	// Try to read the settings file
	settingsPath, err := GetSettingsFilePath()
	if err != nil {
		return fmt.Errorf("failed to get settings file path: %w", err)
	}

	data, err := os.ReadFile(settingsPath)
	if err != nil {
		if os.IsNotExist(err) {
			// If file doesn't exist, create it with default settings
			return SaveSettings()
		}
		return fmt.Errorf("error reading settings file: %w", err)
	}

	// Parse the JSON
	if err := json.Unmarshal(data, &CurrentSettings); err != nil {
		return fmt.Errorf("error parsing settings file: %w", err)
	}

	// Apply the loaded settings
	ApplySettings()

	return nil
}

// SaveSettings saves the current settings to the config file
func SaveSettings() error {
	settingsPath, err := GetSettingsFilePath()
	if err != nil {
		return fmt.Errorf("failed to get settings file path: %w", err)
	}

	// Marshal the settings to JSON
	data, err := json.MarshalIndent(CurrentSettings, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling settings: %w", err)
	}

	// Write to the file
	if err := os.WriteFile(settingsPath, data, 0644); err != nil {
		return fmt.Errorf("error writing settings file: %w", err)
	}

	return nil
}

// UpdateSettings changes a setting and saves the changes
func UpdateSettings(settings UserSettings) error {
	if settings.ThemeName != "" {
		CurrentSettings.ThemeName = settings.ThemeName
	}

	if settings.CursorType != "" {
		CurrentSettings.CursorType = settings.CursorType
	}

	// Update game mode if provided
	if settings.GameMode != "" {
		CurrentSettings.GameMode = settings.GameMode
	}

	// Update numbers setting if provided (use zero value to skip)
	if settings.UseNumbers != CurrentSettings.UseNumbers {
		CurrentSettings.UseNumbers = settings.UseNumbers
	}

	// Apply the changes
	ApplySettings()

	// Save to disk
	return SaveSettings()
}

// ApplySettings applies the current settings to the application
func ApplySettings() {
	// Apply theme
	if CurrentSettings.ThemeName != "" {
		LoadTheme(CurrentSettings.ThemeName)
		UpdateStyles()
	}

	// Apply cursor type
	if CurrentSettings.CursorType == "underline" {
		DefaultCursorType = UnderlineCursor
	} else {
		DefaultCursorType = BlockCursor
	}
}

// InitSettings initializes the settings system
func InitSettings() {
	if err := LoadSettings(); err != nil {
		fmt.Printf("Warning: Could not load settings: %v\n", err)
		fmt.Println("Using default settings")

		// Use defaults
		CurrentSettings = DefaultSettings
		ApplySettings()
	}
}
