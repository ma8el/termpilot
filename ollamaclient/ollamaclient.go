package ollamaclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OllamaResponse struct {
	ID      string `json:"id"`
	Model   string `json:"model"`
	Created int64  `json:"created"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type OllamaModelResponse struct {
	Object string `json:"object"`
	Data   []struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		OwnedBy string `json:"owned_by"`
		Created int    `json:"created"`
	} `json:"data"`
}

type OllamaClient struct {
	BaseURL string
	Port    string
	Version string
	Model   string
}

func NewOllamaClient(baseURL string, model string, port string, version string) *OllamaClient {
	return &OllamaClient{
		BaseURL: baseURL,
		Port:    port,
		Version: version,
		Model:   model,
	}
}

func (c *OllamaClient) ChatCompletion(prompt string, messages []Message) (string, error) {
	url := fmt.Sprintf("%s:%s/%s/chat/completions", c.BaseURL, c.Port, c.Version)

	requestBody := map[string]interface{}{
		"model": c.Model,
		"messages": append(messages, Message{
			Role:    "user",
			Content: prompt,
		}),
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var ollamaResponse OllamaResponse
	err = json.Unmarshal(body, &ollamaResponse)
	if err != nil {
		return "", err
	}

	if len(ollamaResponse.Choices) == 0 {
		return "", errors.New("no choices returned")
	}

	return ollamaResponse.Choices[0].Message.Content, nil
}

func (c *OllamaClient) ListModels() ([]string, error) {
	url := fmt.Sprintf("%s:%s/%s/models", c.BaseURL, c.Port, c.Version)

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var models OllamaModelResponse
	err = json.Unmarshal(body, &models)
	if err != nil {
		return nil, err
	}

	var modelNames []string
	for _, model := range models.Data {
		modelNames = append(modelNames, model.ID)
	}

	return modelNames, nil
}
