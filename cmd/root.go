package cmd

import (
	"os"

	"github.com/prime-run/go-typer/ui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-typer",
	Short: "A terminal-based typing game",
	Long: `Go Typer is a terminal-based typing game with features like:
- Real-time WPM calculation
- Customizable themes
- Multiple cursor styles
- And more!`,
	Run: func(cmd *cobra.Command, args []string) {
		ui.RunStartScreen()
	},
}


func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	ui.InitSettings()


	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
