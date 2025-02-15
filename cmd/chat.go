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

	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"
)

func init() {
	chatCmd.Flags().Bool("list", false, "list all conversations")
	chatCmd.Flags().String("continue", "", "continue a conversation")
	chatCmd.Flags().Bool("continue-last", false, "continue the last conversation")
	chatCmd.Flags().Bool("list-models", false, "list all models")
}

func fancyPrint(text string) string {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
	)
	if err != nil {
		return text
	}
	out, err := renderer.Render(text)
	if err != nil {
		return text
	}
	return out
}

func listConversations() {
	conversations, err := db.GetAllConversations()
	if err != nil {
		log.Fatalf("Failed to list conversations: %v", err)
	}
	numberOfConversations := len(conversations)
	fmt.Println("Conversations (", numberOfConversations, "):")
	for _, conversation := range conversations {
		fmt.Println(conversation.ID, conversation.Title)
	}
}

func continueConversation(conversationId string, args []string, ollamaClient *ollamaclient.OllamaClient) {
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

	response, err := ollamaClient.ChatCompletion(prompt, messages)
	if err != nil {
		log.Fatalf("Failed to get response: %v", err)
	}

	conversation.Messages = append(conversation.Messages, models.Message{Content: prompt, Role: "user"})
	conversation.Messages = append(conversation.Messages, models.Message{Content: response, Role: "assistant"})

	db.UpdateConversation(*conversation)

	fmt.Print(fancyPrint(response))
}

func startConversation(args []string, ollamaClient *ollamaclient.OllamaClient) {
	prompt := strings.Join(args, " ")
	response, err := ollamaClient.ChatCompletion(prompt, []ollamaclient.Message{})
	if err != nil {
		log.Fatalf("Failed to get response: %v", err)
	}

	db.CreateConversation(models.Conversation{
		ID:       fmt.Sprintf("%x", sha256.Sum256([]byte(time.Now().String())))[:8],
		Title:    prompt[:min(len(prompt), 20)],
		Messages: []models.Message{{Content: prompt, Role: "user"}, {Content: response, Role: "assistant"}},
	})

	fmt.Print(fancyPrint(response))
}

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Chat with Termpilot",
	Run: func(cmd *cobra.Command, args []string) {
		baseUrl, err := cmd.Flags().GetString("base-url")
		if err != nil {
			log.Fatalf("Failed to get base-url: %v", err)
		}

		port, err := cmd.Flags().GetString("port")
		if err != nil {
			log.Fatalf("Failed to get port: %v", err)
		}

		model, err := cmd.Flags().GetString("model")
		if err != nil {
			log.Fatalf("Failed to get model: %v", err)
		}

		version, err := cmd.Flags().GetString("version")
		if err != nil {
			log.Fatalf("Failed to get version: %v", err)
		}

		if err := ollamaclient.StartOllamaIfNotRunning(); err != nil {
			log.Fatalf("Failed to start ollama: %v", err)
		}

		list, err := cmd.Flags().GetBool("list")
		if err != nil {
			log.Fatalf("Failed to get list: %v", err)
		}

		if list {
			listConversations()
			return
		}

		ollamaClient := ollamaclient.NewOllamaClient(baseUrl, model, port, version)

		listModels, err := cmd.Flags().GetBool("list-models")
		if err != nil {
			log.Fatalf("Failed to get list-models: %v", err)
		}

		if listModels {
			models, err := ollamaClient.ListModels()
			if err != nil {
				log.Fatalf("Failed to list models: %v", err)
			}
			fmt.Println(models)
			return
		}

		conversationId, err := cmd.Flags().GetString("continue")
		if err != nil {
			log.Fatalf("Failed to get continue: %v", err)
		}

		if conversationId != "" {
			continueConversation(conversationId, args, ollamaClient)
			return
		}

		continueLast, err := cmd.Flags().GetBool("continue-last")
		if err != nil {
			log.Fatalf("Failed to get continue-last: %v", err)
		}

		if continueLast {
			conversation, err := db.GetLastConversation()

			if err != nil {
				log.Fatalf("Failed to get last conversation: %v", err)
			}

			continueConversation(conversation.ID, args, ollamaClient)
			return
		}

		startConversation(args, ollamaClient)
	},
}
