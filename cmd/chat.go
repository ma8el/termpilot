package cmd

import (
	"fmt"
	"log"

	"termpilot/ollamaclient"

	"github.com/spf13/cobra"
)

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Chat with Termpilot",
	Run: func(cmd *cobra.Command, args []string) {
		model, err := cmd.Flags().GetString("model")
		if err != nil {
			log.Fatalf("Failed to get model: %v", err)
		}

		ollamaClient := ollamaclient.NewOllamaClient("http://localhost", model, "11434", "v1")
		response, err := ollamaClient.ChatCompletion("Hello, how are you?")
		if err != nil {
			log.Fatalf("Failed to get response: %v", err)
		}

		fmt.Println(response)
	},
}
