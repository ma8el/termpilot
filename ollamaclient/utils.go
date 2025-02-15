package ollamaclient

import (
	"fmt"
	"net/http"
	"os/exec"
	"time"
)

func IsOllamaRunning(baseURL string, port string) bool {
	url := fmt.Sprintf("%s:%s/", baseURL, port)
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func StartOllama() error {
	cmd := exec.Command("ollama", "serve")

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start Ollama: %v", err)
	}

	for i := 0; i < 10; i++ {
		if IsOllamaRunning("http://localhost", "11434") {
			return nil
		}
		time.Sleep(1 * time.Second)
	}

	return fmt.Errorf("ollama failed to start after 10 seconds")
}

func AskToStartOllama() bool {
	fmt.Print("Ollama is not running. Would you like to start it? (y/n): ")
	var response string
	fmt.Scanln(&response)
	return response == "y" || response == "Y"
}

func StartOllamaIfNotRunning() error {
	if !IsOllamaRunning("http://localhost", "11434") {
		if AskToStartOllama() {
			return StartOllama()
		}
		return fmt.Errorf("ollama is not running")
	}
	return nil
}
