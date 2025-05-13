package ollamaclient

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOllamaClient(t *testing.T) {
	// Setup mock server
	mockServer := setupMockServer()
	defer mockServer.Close()

	// Create a client that points to our mock server
	client := NewOllamaClient(mockServer.URL, "test-model", "", "v1")

	// Test ChatCompletion
	t.Run("ChatCompletion", func(t *testing.T) {
		response, err := client.ChatCompletion("Hello, how are you?", []Message{})
		assert.NoError(t, err)
		assert.Contains(t, response, "I'm doing well")
	})

	// Test ListModels
	t.Run("ListModels", func(t *testing.T) {
		models, err := client.ListModels()
		assert.NoError(t, err)
		assert.Equal(t, 2, len(models))
		assert.Contains(t, models, "llama3")
		assert.Contains(t, models, "mistral")
	})
}

func TestIsOllamaRunning(t *testing.T) {
	// Setup mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	// Extract host and port from mock server URL
	baseURL := "http://" + mockServer.Listener.Addr().String()
	port := "" // Empty since the full URL already includes the port

	// Test IsOllamaRunning with mock server
	running := IsOllamaRunning(baseURL, port)
	assert.True(t, running)

	// Test with non-existent server
	notRunning := IsOllamaRunning("http://localhost", "12345")
	assert.False(t, notRunning)
}

// Mock server setup for testing Ollama API
func setupMockServer() *httptest.Server {
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
							"content": "I'm doing well, thank you for asking!"
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
						"id": "llama3",
						"object": "model",
						"owned_by": "user",
						"created": 1630000000
					},
					{
						"id": "mistral",
						"object": "model",
						"owned_by": "user",
						"created": 1630000000
					}
				]
			}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}
