package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestModels(t *testing.T) {
	// Setup in-memory DB with foreign key constraints enabled
	db, err := gorm.Open(sqlite.Open(":memory:?_foreign_keys=on"), &gorm.Config{})
	assert.NoError(t, err)

	// Execute PRAGMA to ensure foreign keys are enabled
	db.Exec("PRAGMA foreign_keys = ON")

	// Auto migrate schemas
	err = db.AutoMigrate(&Conversation{}, &Message{})
	assert.NoError(t, err)

	// Test conversation and messages creation
	t.Run("CreateModels", func(t *testing.T) {
		// Create a conversation with messages
		conversation := Conversation{
			ID:        "test-id",
			Title:     "Test Conversation",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		result := db.Create(&conversation)
		assert.NoError(t, result.Error)

		// Add messages to the conversation
		messages := []Message{
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

		result = db.Create(&messages)
		assert.NoError(t, result.Error)

		// Verify conversation was created
		var fetchedConv Conversation
		result = db.First(&fetchedConv, "id = ?", "test-id")
		assert.NoError(t, result.Error)
		assert.Equal(t, "Test Conversation", fetchedConv.Title)

		// Verify messages were created and associated with conversation
		var fetchedMsgs []Message
		result = db.Where("conversation_id = ?", "test-id").Find(&fetchedMsgs)
		assert.NoError(t, result.Error)
		assert.Equal(t, 2, len(fetchedMsgs))
		assert.Equal(t, "Hello", fetchedMsgs[0].Content)
		assert.Equal(t, "user", fetchedMsgs[0].Role)
	})

	// Test relationship between Conversation and Message
	t.Run("TestRelationships", func(t *testing.T) {
		// Create new conversation with preloaded messages
		var convWithMessages Conversation
		result := db.Preload("Messages").First(&convWithMessages, "id = ?", "test-id")
		assert.NoError(t, result.Error)

		// Verify preloading worked
		assert.Equal(t, 2, len(convWithMessages.Messages))

		// Test cascade delete (messages should be deleted when conversation is deleted)
		result = db.Delete(&convWithMessages)
		assert.NoError(t, result.Error)

		// Verify conversation was deleted
		var deletedConv Conversation
		result = db.First(&deletedConv, "id = ?", "test-id")
		assert.Error(t, result.Error) // Should be "record not found"

		// Verify messages were deleted (or orphaned based on your FK constraints)
		var remainingMsgs []Message
		result = db.Where("conversation_id = ?", "test-id").Find(&remainingMsgs)
		assert.NoError(t, result.Error)        // Query itself should work
		assert.Equal(t, 0, len(remainingMsgs)) // But no messages should be found
	})
}
