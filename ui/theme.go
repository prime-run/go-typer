package ui

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

const (
	ThemeDefault    = "default"
	ThemeDark       = "dark"
	ThemeMonochrome = "monochrome"
)

type ThemeColors struct {
	HelpText string `yaml:"help_text"`
	Timer    string `yaml:"timer"`
	Border   string `yaml:"border"`

	TextDim          string `yaml:"text_dim"`
	TextPreview      string `yaml:"text_preview"`
	TextCorrect      string `yaml:"text_correct"`
	TextError        string `yaml:"text_error"`
	TextPartialError string `yaml:"text_partial_error"`

	CursorFg        string `yaml:"cursor_fg"`
	CursorBg        string `yaml:"cursor_bg"`
	CursorUnderline string `yaml:"cursor_underline"`

	Padding string `yaml:"padding"`
}

var DefaultTheme = ThemeColors{

	HelpText: "#626262",
	Timer:    "#FFDB58",
	Border:   "#7F9ABE",

	TextDim:          "#555555",
	TextPreview:      "#7F9ABE",
	TextCorrect:      "#00FF00",
	TextError:        "#FF0000",
	TextPartialError: "#FF8C00",

	CursorFg:        "#FFFFFF",
	CursorBg:        "#00AAFF",
	CursorUnderline: "#00AAFF",

	Padding: "#888888",
}

var CurrentTheme ThemeColors

func GetThemePath(themeName string) string {
	themeName = strings.TrimPrefix(themeName, "-")

	if strings.HasSuffix(themeName, ".yml") {
		return themeName
	}

	configDir, err := GetConfigDir()
	if err != nil {
		return filepath.Join("colorschemes", themeName+".yml")
	}

	colorschemesDir := filepath.Join(configDir, "colorschemes")
	if err := os.MkdirAll(colorschemesDir, 0755); err != nil {
		return filepath.Join("colorschemes", themeName+".yml")
	}

	return filepath.Join(colorschemesDir, themeName+".yml")
}

func LoadTheme(themeNameOrPath string) error {
	CurrentTheme = DefaultTheme

	if strings.TrimSpace(themeNameOrPath) == "" {
		return fmt.Errorf("empty theme name")
	}

	themePath := GetThemePath(themeNameOrPath)

	data, err := os.ReadFile(themePath)
	if err != nil {
		if os.IsNotExist(err) {
			themeName := filepath.Base(themePath)
			themeName = strings.TrimSuffix(themeName, ".yml")
			if !isValidThemeName(themeName) {
				return fmt.Errorf("invalid theme name: %s", themeName)
			}

			themeDir := filepath.Dir(themePath)
			if err := os.MkdirAll(themeDir, 0755); err != nil {
				return fmt.Errorf("error creating theme directory: %w", err)
			}

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

	if err := yaml.Unmarshal(data, &CurrentTheme); err != nil {
		return fmt.Errorf("error parsing theme file: %w", err)
	}

	return nil
}

func isValidThemeName(name string) bool {
	if name == "" {
		return false
	}

	for _, c := range name {
		if !isValidThemeNameChar(c) {
			return false
		}
	}

	return true
}

func isValidThemeNameChar(c rune) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') ||
		c == '_' || c == '-'
}

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
		hexColor = "#FFFFFF"
	}

	return lipgloss.Color(hexColor)
}

func ListAvailableThemes() []string {
	themes := []string{ThemeDefault, ThemeDark, ThemeMonochrome}

	configDir, err := GetConfigDir()
	if err == nil {
		colorschemesDir := filepath.Join(configDir, "colorschemes")
		files, err := os.ReadDir(colorschemesDir)
		if err == nil {
			for _, file := range files {
				if !file.IsDir() && strings.HasSuffix(file.Name(), ".yml") {
					themeName := strings.TrimSuffix(file.Name(), ".yml")
					if themeName != ThemeDefault && themeName != ThemeDark && themeName != ThemeMonochrome {
						themes = append(themes, themeName)
					}
				}
			}
		}
	}

	files, err := os.ReadDir("colorschemes")
	if err == nil {
		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".yml") {
				themeName := strings.TrimSuffix(file.Name(), ".yml")
				if themeName != ThemeDefault && themeName != ThemeDark && themeName != ThemeMonochrome {
					themes = append(themes, themeName)
				}
			}
		}
	}

	return themes
}

func InitTheme() {
	ensureDefaultThemesExist()

	themeFile := GetThemePath(ThemeDefault)
	if err := LoadTheme(themeFile); err != nil {
		fmt.Printf("Warning: Could not load theme file: %v\n", err)
		fmt.Println("Using default theme")
	}
}

func ensureDefaultThemesExist() {
	defaultThemes := map[string]ThemeColors{
		ThemeDefault: DefaultTheme,
		ThemeDark: {
			HelpText: "#888888",
			Timer:    "#A177FF",
			Border:   "#5661B3",

			TextDim:          "#444444",
			TextPreview:      "#8892BF",
			TextCorrect:      "#36D399",
			TextError:        "#F87272",
			TextPartialError: "#FBBD23",

			CursorFg:        "#222222",
			CursorBg:        "#7B93DB",
			CursorUnderline: "#7B93DB",

			Padding: "#666666",
		},
		ThemeMonochrome: {
			HelpText: "#AAAAAA",
			Timer:    "#FFFFFF",
			Border:   "#DDDDDD",

			TextDim:          "#777777",
			TextPreview:      "#DDDDDD",
			TextCorrect:      "#FFFFFF",
			TextError:        "#444444",
			TextPartialError: "#BBBBBB",

			CursorFg:        "#000000",
			CursorBg:        "#FFFFFF",
			CursorUnderline: "#FFFFFF",

			Padding: "#999999",
		},
	}

	configDir, err := GetConfigDir()
	if err != nil {
		fmt.Printf("Warning: Could not get config directory: %v\n", err)
		return
	}

	colorschemesDir := filepath.Join(configDir, "colorschemes")
	if err := os.MkdirAll(colorschemesDir, 0755); err != nil {
		fmt.Printf("Warning: Could not create colorschemes directory: %v\n", err)
		return
	}

	for themeName, colors := range defaultThemes {
		themePath := filepath.Join(colorschemesDir, themeName+".yml")

		if _, err := os.Stat(themePath); err == nil {
			continue
		}

		yamlData, err := yaml.Marshal(colors)
		if err != nil {
			fmt.Printf("Warning: Could not marshal %s theme: %v\n", themeName, err)
			continue
		}

		if err := os.WriteFile(themePath, yamlData, 0644); err != nil {
			fmt.Printf("Warning: Could not create %s theme file: %v\n", themeName, err)
			continue
		}

		fmt.Printf("Created %s theme file at %s\n", themeName, themePath)
	}
}
