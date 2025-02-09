package cmd

import (
	"crypto/sha256"
	"fmt"
	"log"
	"strings"
	"time"

	"termpilot/db"
	"termpilot/models"
	"termpilot/ollamaclient"

	"github.com/spf13/cobra"
)

func init() {
	chatCmd.Flags().Bool("list", false, "list all conversations")
	chatCmd.Flags().String("continue", "", "continue a conversation")
	chatCmd.Flags().String("continue-last", "", "continue the last conversation")
}

func listConversations() {
	conversations, err := db.GetAllConversations()
	if err != nil {
		log.Fatalf("Failed to list conversations: %v", err)
	}
	fmt.Println(conversations)
}

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

		list, err := cmd.Flags().GetBool("list")
		if err != nil {
			log.Fatalf("Failed to get list: %v", err)
		}

		if list {
			listConversations()
			return
		}

		conversationId, err := cmd.Flags().GetString("continue")
		if err != nil {
			log.Fatalf("Failed to get continue: %v", err)
		}

		if conversationId != "" {
			conversation, err := db.GetConversation(conversationId)
			if err != nil {
				log.Fatalf("Failed to get conversation: %v", err)
			}

			if len(conversation.Messages) == 0 {
				log.Fatalf("Conversation has no messages")
			}

			var messages []ollamaclient.Message

			for _, message := range conversation.Messages {
				messages = append(messages, ollamaclient.Message{
					Role:    message.Role,
					Content: message.Content,
				})
			}

			prompt := strings.Join(args, " ")

			ollamaClient := ollamaclient.NewOllamaClient(baseUrl, model, port, version)
			response, err := ollamaClient.ChatCompletion(prompt, messages)
			if err != nil {
				log.Fatalf("Failed to get response: %v", err)
			}

			conversation.Messages = append(conversation.Messages, models.Message{Content: prompt, Role: "user"})
			conversation.Messages = append(conversation.Messages, models.Message{Content: response, Role: "assistant"})

			db.UpdateConversation(*conversation)

			fmt.Println(response)
			return
		}

		prompt := strings.Join(args, " ")

		ollamaClient := ollamaclient.NewOllamaClient(baseUrl, model, port, version)
		response, err := ollamaClient.ChatCompletion(prompt, []ollamaclient.Message{})
		if err != nil {
			log.Fatalf("Failed to get response: %v", err)
		}

		db.CreateConversation(models.Conversation{
			ID:        fmt.Sprintf("%x", sha256.Sum256([]byte(time.Now().String())))[:8],
			Title:     prompt[:min(len(prompt), 20)],
			Messages:  []models.Message{{Content: prompt, Role: "user"}, {Content: response, Role: "assistant"}},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})

		fmt.Println(response)
	},
}
