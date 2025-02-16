package cmd

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(uiCmd)
}

var uiCmd = &cobra.Command{
	Use:   "ui",
	Short: "Start the interactive TUI",
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
			log.Fatalf("Error running TUI: %v", err)
		}
	},
}
