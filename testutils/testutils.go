package testutils

import (
	"net/http"
	"net/http/httptest"
	"os"
	"termpilot/models"
	"termpilot/ollamaclient"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SetupTestDB creates a temporary test database
func SetupTestDB(path string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Migrate models
	db.AutoMigrate(&models.Conversation{}, &models.Message{})

	return db, nil
}

// CleanupTestDB removes the temporary test database
func CleanupTestDB(path string) {
	os.Remove(path)
}

// CreateTestConversation creates a test conversation
func CreateTestConversation(db *gorm.DB) (*models.Conversation, error) {
	conversation := models.Conversation{
		ID:        "test123",
		Title:     "Test Conversation",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.Create(&conversation).Error; err != nil {
		return nil, err
	}

	messages := []models.Message{
		{
			Content:        "Hello",
			Role:           "user",
			ConversationID: conversation.ID,
		},
		{
			Content:        "Hi there!",
			Role:           "assistant",
			ConversationID: conversation.ID,
		},
	}

	if err := db.Create(&messages).Error; err != nil {
		return nil, err
	}

	return &conversation, nil
}

// MockOllamaServer creates a mock server for Ollama API testing
func MockOllamaServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle different API endpoints
		switch r.URL.Path {
		case "/v1/chat/completions":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"id": "test-id",
				"model": "test-model",
				"created": 1630000000,
				"choices": [
					{
						"index": 0,
						"message": {
							"role": "assistant",
							"content": "I'm a test response"
						}
					}
				]
			}`))
		case "/v1/models":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"object": "list",
				"data": [
					{
						"id": "test-model",
						"object": "model",
						"owned_by": "user",
						"created": 1630000000
					}
				]
			}`))
		default:
			w.WriteHeader(http.StatusOK)
		}
	}))
}

// NewTestOllamaClient creates an Ollama client for testing
func NewTestOllamaClient(server *httptest.Server) *ollamaclient.OllamaClient {
	return ollamaclient.NewOllamaClient(
		server.URL,
		"test-model",
		"",
		"v1",
	)
}
