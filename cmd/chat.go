package cmd

import (
	"fmt"
	"log"
	"strings"

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

		baseUrl, err := cmd.Flags().GetString("base-url")
		if err != nil {
			log.Fatalf("Failed to get base-url: %v", err)
		}

		port, err := cmd.Flags().GetString("port")
		if err != nil {
			log.Fatalf("Failed to get port: %v", err)
		}

		version, err := cmd.Flags().GetString("version")
		if err != nil {
			log.Fatalf("Failed to get version: %v", err)
		}

		prompt := strings.Join(args, " ")

		ollamaClient := ollamaclient.NewOllamaClient(baseUrl, model, port, version)
		response, err := ollamaClient.ChatCompletion(prompt)
		if err != nil {
			log.Fatalf("Failed to get response: %v", err)
		}

		fmt.Println(response)
	},
}
