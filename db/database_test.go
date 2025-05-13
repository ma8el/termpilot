package db

import (
	"os"
	"termpilot/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestDatabaseOperations(t *testing.T) {
	// Use in-memory SQLite for testing
	tempFile := "test.db"

	// Setup
	DB, err := initTestDB(tempFile)
	assert.NoError(t, err)
	assert.NotNil(t, DB)

	// Teardown
	defer os.Remove(tempFile)

	// Test creating a conversation
	conversation := models.Conversation{
		ID:        "test123",
		Title:     "Test Conversation",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Messages: []models.Message{
			{Content: "Hello", Role: "user"},
			{Content: "Hi there", Role: "assistant"},
		},
	}

	createdConv, err := CreateConversation(conversation)
	assert.NoError(t, err)
	assert.Equal(t, "test123", createdConv.ID)
	assert.Equal(t, "Test Conversation", createdConv.Title)

	// Test getting a conversation
	fetchedConv, err := GetConversation("test123")
	assert.NoError(t, err)
	assert.Equal(t, "test123", fetchedConv.ID)
	assert.Equal(t, 2, len(fetchedConv.Messages))
	assert.Equal(t, "Hello", fetchedConv.Messages[0].Content)

	// Test updating a conversation
	fetchedConv.Title = "Updated Title"
	fetchedConv.Messages = append(fetchedConv.Messages, models.Message{
		Content: "New message",
		Role:    "user",
	})

	updatedConv, err := UpdateConversation(*fetchedConv)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Title", updatedConv.Title)
	assert.Equal(t, 3, len(updatedConv.Messages))

	// Test getting all conversations
	allConvs, err := GetAllConversations()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(allConvs), 1)

	// Test getting last conversation
	lastConv, err := GetLastConversation()
	assert.NoError(t, err)
	assert.Equal(t, "test123", lastConv.ID)

	// Test deleting a conversation
	err = DeleteConversation("test123")
	assert.NoError(t, err)

	// Verify deletion
	_, err = GetConversation("test123")
	assert.Error(t, err) // Should get an error now
}

func initTestDB(path string) (*gorm.DB, error) {
	var err error
	DB, err = gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	DB.AutoMigrate(&models.Conversation{}, &models.Message{})
	return DB, nil
}
