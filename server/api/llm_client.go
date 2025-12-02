package api

import (
	"log"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

// createLLMClient creates an OpenAI client with optional custom base URL
func createLLMClient() *openai.Client {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Printf("Warning: OPENAI_API_KEY is not configured")
		return nil
	}

	// Support OpenAI compatible APIs with custom base URL
	baseURL := os.Getenv("OPENAI_BASE_URL")
	if baseURL != "" {
		log.Printf("Using custom base URL: %s", baseURL)
		config := openai.DefaultConfig(apiKey)
		config.BaseURL = baseURL
		return openai.NewClientWithConfig(config)
	}

	log.Printf("Using default OpenAI base URL: https://api.openai.com/v1")
	return openai.NewClient(apiKey)
}
