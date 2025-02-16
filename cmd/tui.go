package cmd

import (
	"crypto/sha256"
	"fmt"
	"log"
	"strings"
	"termpilot/db"
	"termpilot/models"
	"termpilot/ollamaclient"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/viper"
)

type item struct {
	id    string
	title string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.id }
func (i item) FilterValue() string { return i.title }

type model struct {
	conversations list.Model
	messages      viewport.Model
	input         textinput.Model
	selectedConv  *models.Conversation
	width         int
	height        int
	state         uiState
	ollamaClient  *ollamaclient.OllamaClient
}

type uiState int

const (
	stateBrowsing uiState = iota
	stateChatting
	stateNewChat
)

func getOllamaClient() *ollamaclient.OllamaClient {
	return ollamaclient.NewOllamaClient(
		viper.GetString("base-url"),
		viper.GetString("model"),
		viper.GetString("port"),
		viper.GetString("version"),
	)
}

func initialModel() model {
	convs, _ := db.GetAllConversations()
	items := make([]list.Item, len(convs))
	for i, conv := range convs {
		items[i] = item{id: conv.ID, title: conv.Title}
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Conversations"

	ti := textinput.New()
	ti.Placeholder = "Type your message..."
	ti.Focus()

	m := model{
		conversations: l,
		input:         ti,
		state:         stateBrowsing,
		ollamaClient:  getOllamaClient(),
	}
	m.messages = viewport.New(80, 20)
	m.messages.HighPerformanceRendering = false
	m.messages.SetContent("Loading conversations...")
	m.messages.Style = m.messages.Style.
		Margin(1, 2).
		Padding(1, 1)
	return m
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		tea.SetWindowTitle("Termpilot"),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case stateBrowsing:
		return updateBrowsing(m, msg)
	case stateChatting:
		return updateChatting(m, msg)
	case stateNewChat:
		return updateNewChat(m, msg)
	}
	return m, nil
}

func (m model) View() string {
	switch m.state {
	case stateBrowsing:
		return browsingView(m)
	case stateChatting:
		return chatView(m)
	case stateNewChat:
		return newChatView(m)
	}
	return ""
}

func updateBrowsing(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			selected := m.conversations.SelectedItem().(item)
			conv, _ := db.GetConversation(selected.id)
			m.selectedConv = conv
			m.state = stateChatting
			m.messages.SetContent(formatMessages(conv.Messages))
			m.messages.GotoBottom()
			return m, nil
		case "n":
			m.state = stateNewChat
			m.input.Focus()
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.conversations.SetSize(msg.Width, msg.Height-4)
		m.messages.Width = msg.Width
		m.messages.Height = msg.Height - 4
	}

	var cmd tea.Cmd
	m.conversations, cmd = m.conversations.Update(msg)
	return m, cmd
}

func browsingView(m model) string {
	return m.conversations.View()
}

func chatView(m model) string {
	return fmt.Sprintf(
		"Chat: %s\n%s\n\n%s",
		m.selectedConv.Title,
		m.messages.View(),
		m.input.View(),
	)
}

func newChatView(m model) string {
	return fmt.Sprintf(
		"New Chat\n\n%s\n\n%s",
		"Type your message below (Press Esc to cancel)",
		m.input.View(),
	)
}

func updateNewChat(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			m.state = stateBrowsing
			m.input.Reset()
			return m, nil

		case tea.KeyEnter:
			prompt := m.input.Value()
			m.input.Reset()

			response, err := m.ollamaClient.ChatCompletion(prompt, []ollamaclient.Message{})
			if err != nil {
				log.Printf("Chat error: %v", err)
				return m, nil
			}

			conv := models.Conversation{
				ID:    fmt.Sprintf("%x", sha256.Sum256([]byte(time.Now().String())))[:8],
				Title: prompt[:min(len(prompt), 20)],
				Messages: []models.Message{
					{Content: prompt, Role: "user"},
					{Content: response, Role: "assistant"},
				},
			}

			if _, err := db.CreateConversation(conv); err != nil {
				log.Printf("Save error: %v", err)
			}

			convs, _ := db.GetAllConversations()
			items := make([]list.Item, len(convs))
			for i, conv := range convs {
				items[i] = item{id: conv.ID, title: conv.Title}
			}
			m.conversations.SetItems(items)

			m.messages.SetContent(formatMessages(conv.Messages))
			m.messages.GotoBottom()

			m.state = stateBrowsing
			return m, nil

		}

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func updateChatting(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.state = stateBrowsing
			m.selectedConv = nil
			m.messages.GotoBottom()
			return m, nil

		case "enter":
			prompt := m.input.Value()
			m.input.Reset()

			var messages []ollamaclient.Message
			for _, msg := range m.selectedConv.Messages {
				messages = append(messages, ollamaclient.Message{
					Role:    msg.Role,
					Content: msg.Content,
				})
			}

			messages = append(messages, ollamaclient.Message{
				Role:    "user",
				Content: prompt,
			})

			response, err := m.ollamaClient.ChatCompletion(prompt, messages)
			if err != nil {
				log.Printf("Chat error: %v", err)
				return m, nil
			}

			m.selectedConv.Messages = append(m.selectedConv.Messages,
				models.Message{Content: prompt, Role: "user"},
				models.Message{Content: response, Role: "assistant"},
			)

			if _, err := db.UpdateConversation(*m.selectedConv); err != nil {
				log.Printf("Save error: %v", err)
			}

			m.messages.SetContent(formatMessages(m.selectedConv.Messages))
			m.messages.GotoBottom()
			return m, tea.Batch(
				cmd,
				func() tea.Msg { return tea.WindowSizeMsg{Width: m.width, Height: m.height} },
			)

		case "ctrl+c", "q":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.messages.Width = msg.Width
		m.messages.Height = msg.Height - 4 // Account for input field
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func formatMessages(messages []models.Message) string {
	var content strings.Builder
	for _, msg := range messages {
		content.WriteString(
			fmt.Sprintf("**%s**\n%s\n\n",
				strings.ToUpper(msg.Role),
				strings.TrimSpace(msg.Content),
			))
	}
	return content.String()
}
