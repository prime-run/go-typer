/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/prime-run/go-typer/ui"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-typer",
	Short: "A terminal-based typing game",
	Long: `Go Typer is a terminal-based typing game with features like:
- Real-time WPM calculation
- Customizable themes
- Multiple cursor styles
- And more!`,
	// When no subcommand is provided, show the start screen
	Run: func(cmd *cobra.Command, args []string) {
		ui.RunStartScreen()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Initialize settings
	ui.InitSettings()

	// Cobra supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
