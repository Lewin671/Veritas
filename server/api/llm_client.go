package api

import (
	"fmt"
	"log"
	"os"
	"veritas-server/models"
	"veritas-server/services"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// createLLMClient creates an OpenAI client with optional custom base URL (legacy)
// nolint:unused
func createLLMClient() *openai.Client {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Println("Warning: OPENAI_API_KEY is not configured")
		return nil
	}

	// Support OpenAI compatible APIs with custom base URL
	baseURL := os.Getenv("OPENAI_BASE_URL")
	if baseURL != "" {
		log.Printf("Using custom base URL: %s", baseURL)
		client := openai.NewClient(
			option.WithAPIKey(apiKey),
			option.WithBaseURL(baseURL),
		)
		return &client
	}

	log.Println("Using default OpenAI base URL: https://api.openai.com/v1")
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)
	return &client
}

// createLLMClientFromConfig creates an OpenAI client from a ModelConfig
func createLLMClientFromConfig(config *models.ModelConfig, decrypt bool) (*openai.Client, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	apiKey := config.APIKey

	// Decrypt API key if needed and not empty
	if decrypt && apiKey != "" {
		decryptedKey, err := services.DecryptAPIKey(config.APIKey)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt API key: %w", err)
		}
		apiKey = decryptedKey
	}

	// For local models like Ollama, API key can be empty or a placeholder
	// Use "ollama" as default if empty for local models
	if apiKey == "" {
		apiKey = "ollama"
		log.Printf("Using placeholder API key for local model: %s", config.Name)
	}

	// Create client with custom base URL if provided
	if config.BaseURL != "" {
		log.Printf("Using custom base URL: %s for model: %s", config.BaseURL, config.Name)
		client := openai.NewClient(
			option.WithAPIKey(apiKey),
			option.WithBaseURL(config.BaseURL),
		)
		return &client, nil
	}

	log.Printf("Using default OpenAI base URL for model: %s", config.Name)
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)
	return &client, nil
}
