package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	devlog "github.com/prime-run/go-typer/log"
	"github.com/prime-run/go-typer/ui"
	"github.com/prime-run/go-typer/utils"
	"github.com/spf13/cobra"
)

var (
	cursorType string // Cursor type (block or underline)
	themeName  string // Theme name or path to custom theme file
	listThemes bool   // List available themes and exit
	debugMode  bool   // Enable debug mode for performance analysis
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a new game",
	Long:  "Start a new game of Go Typer. This command will initialize a new game session.",
	Run: func(cmd *cobra.Command, args []string) {
		// check if the debug mode is enabled and set the log level accordingly
		if debugMode {
			devlog.DebugEnabled = true
			devlog.InitLog()
			defer devlog.CloseLog()

			cmd.Printf("Debug mode enabled, logging to %s\n", filepath.Join(utils.GetConfigDirPath(), "debug.log"))
		}

		// check if listThemes is enabled and list available themes
		if listThemes {
			themes := ui.ListAvailableThemes()
			fmt.Println("Available themes:")
			for _, theme := range themes {
				fmt.Printf("  - %s\n", theme)
			}
			return
		}

		// check for theme name and load the theme
		if themeName != "" {
			fmt.Printf("Theme name provided: %s", themeName)
			if strings.HasPrefix(themeName, "-") {
				cmd.Printf("Warning: Invalid theme name '%s'. Theme names cannot start with '-'.\n", themeName)
				cmd.Println("Using saved settings")
			} else if utils.IsValidThemeName(themeName) {
				if err := ui.LoadTheme(themeName); err != nil {
					cmd.Printf("Warning: Could not load theme '%s': %v\n", themeName, err)
					cmd.Println("Using saved settings")
				} else {
					ui.UpdateStyles()
					cmd.Printf("Using theme: %s\n", utils.GetDisplayThemeName(themeName))
					ui.CurrentSettings.ThemeName = themeName
				}
			} else {
				cmd.Printf("Warning: Invalid theme name '%s'. Using saved settings.\n", themeName)
			}
		}

		// check for cursor type and set it
		if cursorType != "" {
			ui.CurrentSettings.CursorType = cursorType
		}

		// apply settings and start loading
		ui.ApplySettings()
		ui.StartLoadingWithOptions(ui.CurrentSettings.CursorType)
	},
}

func init() {
	// Cursor and theme configuration
	startCmd.Flags().StringVarP(&cursorType, "cursor", "c", "block", "Cursor type (block or underline)")
	startCmd.Flags().StringVarP(&themeName, "theme", "t", "", "Theme name or path to custom theme file (default: default)")
	startCmd.Flags().BoolVar(&listThemes, "list-themes", false, "List available themes and exit")

	// Debug and logging options
	startCmd.Flags().BoolP("verbose", "v", false, "Enable verbose output")
	startCmd.Flags().BoolVar(&debugMode, "debug", false, "Enable debug mode for performance analysis")

	// Add command to root
	rootCmd.AddCommand(startCmd)
}
