package cmd

import (
	"github.com/prime-run/go-typer/ui"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a new game",
	Long:  `Start a new game of Go Typer. This command will initialize a new game session.`,
	Run:   ui.ProgressBar,
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().BoolP("verbose", "v", false, "Enable verbose output")
}
