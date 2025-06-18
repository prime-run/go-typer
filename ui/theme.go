package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/prime-run/go-typer/utils"
	"gopkg.in/yaml.v3"
)

const (
	ThemeDefault    = "default"    // Default theme name
	ThemeDark       = "dark"       // Dark theme name
	ThemeMonochrome = "monochrome" // Monochrome theme name

	YMLSuffix = ".yml" // YAML file suffix
)

type ThemeColors struct {
	HelpText string `yaml:"help_text"` // Help text color
	Timer    string `yaml:"timer"`     // Timer color
	Border   string `yaml:"border"`    // Border color

	TextDim          string `yaml:"text_dim"`           // Dimmed text color
	TextPreview      string `yaml:"text_preview"`       // Preview text color
	TextCorrect      string `yaml:"text_correct"`       // Correct text color
	TextError        string `yaml:"text_error"`         // Error text color
	TextPartialError string `yaml:"text_partial_error"` // Partial error text color

	CursorFg        string `yaml:"cursor_fg"`        // Cursor foreground color
	CursorBg        string `yaml:"cursor_bg"`        // Cursor background color
	CursorUnderline string `yaml:"cursor_underline"` // Cursor underline color

	Padding string `yaml:"padding"` // Padding color
}

var (
	CurrentTheme ThemeColors

	themeColorMap = map[string]*string{
		"help_text":          &CurrentTheme.HelpText,
		"timer":              &CurrentTheme.Timer,
		"border":             &CurrentTheme.Border,
		"text_dim":           &CurrentTheme.TextDim,
		"text_preview":       &CurrentTheme.TextPreview,
		"text_correct":       &CurrentTheme.TextCorrect,
		"text_error":         &CurrentTheme.TextError,
		"text_partial_error": &CurrentTheme.TextPartialError,
		"cursor_fg":          &CurrentTheme.CursorFg,
		"cursor_bg":          &CurrentTheme.CursorBg,
		"cursor_underline":   &CurrentTheme.CursorUnderline,
		"padding":            &CurrentTheme.Padding,
	}

	// Default theme colors
	DefaultTheme = ThemeColors{
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
)

// InitTheme initializes the theme by loading the default theme file.
func InitTheme() {
	ensureDefaultThemesExist()

	if err := LoadTheme(ThemeDefault); err != nil {
		fmt.Printf("Warning: Could not load theme file: %v\n", err)
		fmt.Println("Using default theme")
	}
}

// LoadTheme loads the theme from the specified file or directory.
func LoadTheme(themeNameOrPath string) error {
	CurrentTheme = DefaultTheme

	if strings.TrimSpace(themeNameOrPath) == "" {
		return fmt.Errorf("empty theme name")
	}

	themePath := utils.GetThemePath(themeNameOrPath)

	data, err := os.ReadFile(themePath)
	if err != nil {
		if os.IsNotExist(err) {
			themeName := filepath.Base(themePath)
			themeName = strings.TrimSuffix(themeName, YMLSuffix)
			if !utils.IsValidThemeName(themeName) {
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

// GetColor returns the color associated with the given color name.
// If the color name is not found, it returns a default color (white).
func GetColor(colorName string) lipgloss.Color {
	var hexColor string

	if color, ok := themeColorMap[colorName]; ok {
		hexColor = *color
	} else {
		hexColor = "#FFFFFF"
	}

	return lipgloss.Color(hexColor)
}

func ListAvailableThemes() []string {
	themes := []string{ThemeDefault, ThemeDark, ThemeMonochrome}

	configDir, err := utils.GetAppConfigDir()
	if err == nil {
		colorschemesDir := filepath.Join(configDir, "colorschemes")
		files, err := os.ReadDir(colorschemesDir)
		if err == nil {
			for _, file := range files {
				if !file.IsDir() && strings.HasSuffix(file.Name(), YMLSuffix) {
					themeName := strings.TrimSuffix(file.Name(), YMLSuffix)
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
			if !file.IsDir() && strings.HasSuffix(file.Name(), YMLSuffix) {
				themeName := strings.TrimSuffix(file.Name(), YMLSuffix)
				if themeName != ThemeDefault && themeName != ThemeDark && themeName != ThemeMonochrome {
					themes = append(themes, themeName)
				}
			}
		}
	}

	return themes
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

	configDir, err := utils.GetAppConfigDir()
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
		themePath := filepath.Join(colorschemesDir, themeName+YMLSuffix)

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
