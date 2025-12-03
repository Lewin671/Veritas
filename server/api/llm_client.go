package api

import (
	"log"
	"os"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
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
		client := openai.NewClient(
			option.WithAPIKey(apiKey),
			option.WithBaseURL(baseURL),
		)
		return &client
	}

	log.Printf("Using default OpenAI base URL: https://api.openai.com/v1")
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)
	return &client
}
