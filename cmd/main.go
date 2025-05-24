package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	devlog "github.com/prime-run/go-typer/log"
	"github.com/prime-run/go-typer/ui"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
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
		if debugMode {
			devlog.DebugEnabled = true
			devlog.InitLog()
			defer devlog.CloseLog()

			cmd.Printf("Debug mode enabled, logging to %s\n", filepath.Join(getConfigDirPath(), "debug.log"))
		}

		if listThemes {
			themes := ui.ListAvailableThemes()
			fmt.Println("Available themes:")
			for _, theme := range themes {
				fmt.Printf("  - %s\n", theme)
			}
			return
		}

		if themeName != "" {
			if strings.HasPrefix(themeName, "-") {
				cmd.Printf("Warning: Invalid theme name '%s'. Theme names cannot start with '-'.\n", themeName)
				cmd.Println("Using saved settings")
			} else if isValidThemeName(themeName) {
				if err := ui.LoadTheme(themeName); err != nil {
					cmd.Printf("Warning: Could not load theme '%s': %v\n", themeName, err)
					cmd.Println("Using saved settings")
				} else {
					ui.UpdateStyles()
					cmd.Printf("Using theme: %s\n", getDisplayThemeName(themeName))

					ui.CurrentSettings.ThemeName = themeName
				}
			} else {
				cmd.Printf("Warning: Invalid theme name '%s'. Using saved settings.\n", themeName)
			}
		}

		if cursorType != "" {
			ui.CurrentSettings.CursorType = cursorType
		}

		ui.ApplySettings()

		ui.StartLoadingWithOptions(ui.CurrentSettings.CursorType)
	},
}

func getConfigDirPath() string {
	configDir, err := ui.GetConfigDir()
	if err != nil {
		return os.TempDir()
	}
	return configDir
}

func isValidThemeName(name string) bool {
	if strings.Contains(name, ".") && !strings.HasSuffix(name, ".yml") {
		return false
	}

	if strings.Contains(name, "/") || strings.Contains(name, "\\") {
		_, err := os.Stat(name)
		return err == nil
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

func getDisplayThemeName(themeName string) string {
	if strings.Contains(themeName, "/") || strings.Contains(themeName, "\\") {
		themeName = filepath.Base(themeName)
	}

	themeName = strings.TrimSuffix(themeName, ".yml")

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
