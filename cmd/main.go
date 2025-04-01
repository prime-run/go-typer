package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/prime-run/go-typer/ui"
	"github.com/spf13/cobra"
)

var cursorType string
var themeName string
var listThemes bool
var debugMode bool

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a new game",
	Long:  `Start a new game of Go Typer. This command will initialize a new game session.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Enable debug logging if flag is set
		if debugMode {
			ui.DebugEnabled = true
			ui.InitDebugLog()
			defer ui.CloseDebugLog()

			cmd.Printf("Debug mode enabled, logging to %s\n", filepath.Join(getConfigDirPath(), "debug.log"))
		}

		// If list-themes flag is used, display available themes and exit
		if listThemes {
			themes := ui.ListAvailableThemes()
			fmt.Println("Available themes:")
			for _, theme := range themes {
				fmt.Printf("  - %s\n", theme)
			}
			return
		}

		// If a theme was specified, validate and load it
		if themeName != "" {
			// Check for invalid theme names that might be misinterpreted flags
			if strings.HasPrefix(themeName, "-") {
				cmd.Printf("Warning: Invalid theme name '%s'. Theme names cannot start with '-'.\n", themeName)
				cmd.Println("Using saved settings")
			} else if isValidThemeName(themeName) {
				if err := ui.LoadTheme(themeName); err != nil {
					cmd.Printf("Warning: Could not load theme '%s': %v\n", themeName, err)
					cmd.Println("Using saved settings")
				} else {
					// Update styles after loading the theme
					ui.UpdateStyles()
					cmd.Printf("Using theme: %s\n", getDisplayThemeName(themeName))

					// Override the saved settings for this session
					ui.CurrentSettings.ThemeName = themeName
				}
			} else {
				cmd.Printf("Warning: Invalid theme name '%s'. Using saved settings.\n", themeName)
			}
		}

		// If cursor type is specified, override the saved setting
		if cursorType != "" {
			ui.CurrentSettings.CursorType = cursorType
		}

		// Apply the settings
		ui.ApplySettings()

		// The start command always goes directly to the game
		ui.StartLoadingWithOptions(ui.CurrentSettings.CursorType)
	},
}

// Helper function to get config directory path
func getConfigDirPath() string {
	configDir, err := ui.GetConfigDir()
	if err != nil {
		return os.TempDir()
	}
	return configDir
}

// isValidThemeName checks if a theme name is valid
func isValidThemeName(name string) bool {
	// If it has an extension, it should be .yml
	if strings.Contains(name, ".") && !strings.HasSuffix(name, ".yml") {
		return false
	}

	// If it's a file path, check if it exists
	if strings.Contains(name, "/") || strings.Contains(name, "\\") {
		_, err := os.Stat(name)
		return err == nil
	}

	// Otherwise, it should be a valid theme name (alphanumeric + underscore)
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

// getDisplayThemeName returns a user-friendly name for the theme
func getDisplayThemeName(themeName string) string {
	// If it's a file path, extract just the filename
	if strings.Contains(themeName, "/") || strings.Contains(themeName, "\\") {
		themeName = filepath.Base(themeName)
	}

	// Remove .yml extension if present
	themeName = strings.TrimSuffix(themeName, ".yml")

	// Capitalize the first letter of each word
	words := strings.Split(themeName, "_")
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[0:1]) + word[1:]
		}
	}

	return strings.Join(words, " ")
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().BoolP("verbose", "v", false, "Enable verbose output")
	startCmd.Flags().StringVarP(&cursorType, "cursor", "c", "block", "Cursor type (block or underline)")
	startCmd.Flags().StringVarP(&themeName, "theme", "t", "", "Theme name or path to custom theme file (default: default)")
	startCmd.Flags().BoolVar(&listThemes, "list-themes", false, "List available themes and exit")
	startCmd.Flags().BoolVar(&debugMode, "debug", false, "Enable debug mode for performance analysis")
}
