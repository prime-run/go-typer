package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"gopkg.in/yaml.v3"
)

// Available themes
const (
	ThemeDefault    = "default"
	ThemeDark       = "dark"
	ThemeMonochrome = "monochrome"
)

// ThemeColors holds all the color values loaded from YAML
type ThemeColors struct {
	// UI Elements
	HelpText string `yaml:"help_text"`
	Timer    string `yaml:"timer"`
	Border   string `yaml:"border"`

	// Text Display
	TextDim          string `yaml:"text_dim"`
	TextPreview      string `yaml:"text_preview"`
	TextCorrect      string `yaml:"text_correct"`
	TextError        string `yaml:"text_error"`
	TextPartialError string `yaml:"text_partial_error"`

	// Cursor
	CursorFg        string `yaml:"cursor_fg"`
	CursorBg        string `yaml:"cursor_bg"`
	CursorUnderline string `yaml:"cursor_underline"`

	// Miscellaneous
	Padding string `yaml:"padding"`
}

// DefaultTheme provides fallback colors if theme file can't be loaded
var DefaultTheme = ThemeColors{
	// UI Elements
	HelpText: "#626262",
	Timer:    "#FFDB58",
	Border:   "#7F9ABE",

	// Text Display
	TextDim:          "#555555",
	TextPreview:      "#7F9ABE",
	TextCorrect:      "#00FF00",
	TextError:        "#FF0000",
	TextPartialError: "#FF8C00",

	// Cursor
	CursorFg:        "#FFFFFF",
	CursorBg:        "#00AAFF",
	CursorUnderline: "#00AAFF",

	// Miscellaneous
	Padding: "#888888",
}

// CurrentTheme holds the currently loaded theme colors
var CurrentTheme ThemeColors

// GetThemePath resolves a theme name to its file path
func GetThemePath(themeName string) string {
	// Sanitize theme name to prevent issues with flags or special characters
	if strings.HasPrefix(themeName, "-") {
		// Remove the leading dash to prevent command-line argument confusion
		themeName = strings.TrimPrefix(themeName, "-")
	}

	// If the themeName is already a file path, return it
	if strings.HasSuffix(themeName, ".yml") {
		return themeName
	}

	// Otherwise, construct the path for a built-in theme
	return filepath.Join("colorschemes", themeName+".yml")
}

// LoadTheme loads color theme from a YAML file
func LoadTheme(themeNameOrPath string) error {
	// Start with default theme
	CurrentTheme = DefaultTheme

	// Validate theme name
	if strings.TrimSpace(themeNameOrPath) == "" {
		return fmt.Errorf("empty theme name")
	}

	// Get the theme file path
	themePath := GetThemePath(themeNameOrPath)

	// Try to read the theme file
	data, err := os.ReadFile(themePath)
	if err != nil {
		// If file not found, create default theme file if it's a built-in theme
		if os.IsNotExist(err) {
			// If themePath is not in the colorschemes directory, return error
			if !strings.HasPrefix(themePath, filepath.Join("colorschemes", "")) {
				return fmt.Errorf("theme file not found: %s", themePath)
			}

			// Check if theme name is valid
			themeName := filepath.Base(themePath)
			themeName = strings.TrimSuffix(themeName, ".yml")
			if !isValidThemeName(themeName) {
				return fmt.Errorf("invalid theme name: %s", themeName)
			}

			// Ensure colorschemes directory exists
			if err := os.MkdirAll("colorschemes", 0755); err != nil {
				return fmt.Errorf("error creating colorschemes directory: %w", err)
			}

			// Create default theme file
			yamlData, err := yaml.Marshal(DefaultTheme)
			if err != nil {
				return fmt.Errorf("error marshaling default theme: %w", err)
			}

			if err := os.WriteFile(themePath, yamlData, 0644); err != nil {
				return fmt.Errorf("error writing default theme file: %w", err)
			}

			fmt.Printf("Created default theme file at %s\n", themePath)
			return nil
		}
		return fmt.Errorf("error reading theme file: %w", err)
	}

	// Parse the YAML
	if err := yaml.Unmarshal(data, &CurrentTheme); err != nil {
		return fmt.Errorf("error parsing theme file: %w", err)
	}

	return nil
}

// isValidThemeName checks if a theme name is valid
func isValidThemeName(name string) bool {
	if name == "" {
		return false
	}

	// Check each character
	for _, c := range name {
		if !isValidThemeNameChar(c) {
			return false
		}
	}

	return true
}

// isValidThemeNameChar checks if a character is valid in a theme name
func isValidThemeNameChar(c rune) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') ||
		c == '_' || c == '-'
}

// GetColor returns a lipgloss.Color from the current theme
func GetColor(colorName string) lipgloss.Color {
	var hexColor string

	switch colorName {
	case "help_text":
		hexColor = CurrentTheme.HelpText
	case "timer":
		hexColor = CurrentTheme.Timer
	case "border":
		hexColor = CurrentTheme.Border
	case "text_dim":
		hexColor = CurrentTheme.TextDim
	case "text_preview":
		hexColor = CurrentTheme.TextPreview
	case "text_correct":
		hexColor = CurrentTheme.TextCorrect
	case "text_error":
		hexColor = CurrentTheme.TextError
	case "text_partial_error":
		hexColor = CurrentTheme.TextPartialError
	case "cursor_fg":
		hexColor = CurrentTheme.CursorFg
	case "cursor_bg":
		hexColor = CurrentTheme.CursorBg
	case "cursor_underline":
		hexColor = CurrentTheme.CursorUnderline
	case "padding":
		hexColor = CurrentTheme.Padding
	default:
		hexColor = "#FFFFFF" // Default white
	}

	return lipgloss.Color(hexColor)
}

// ListAvailableThemes returns a list of available theme names
func ListAvailableThemes() []string {
	themes := []string{ThemeDefault, ThemeDark, ThemeMonochrome}

	// Try to list additional themes from the colorschemes directory
	files, err := os.ReadDir("colorschemes")
	if err == nil {
		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".yml") {
				themeName := strings.TrimSuffix(file.Name(), ".yml")
				// Add only if it's not already in the list
				if themeName != ThemeDefault && themeName != ThemeDark && themeName != ThemeMonochrome {
					themes = append(themes, themeName)
				}
			}
		}
	}

	return themes
}

// InitTheme initializes the theme system by loading the default theme
func InitTheme() {
	themeFile := GetThemePath(ThemeDefault)
	if err := LoadTheme(themeFile); err != nil {
		fmt.Printf("Warning: Could not load theme file: %v\n", err)
		fmt.Println("Using default theme")
	}
}
