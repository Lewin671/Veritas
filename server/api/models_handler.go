package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// GetModels returns the list of available LLM models
func GetModels(c *gin.Context) {
	var modelList []Model

	// Try to load models from environment variable
	modelsConfig := os.Getenv("AVAILABLE_MODELS")
	if modelsConfig != "" {
		log.Printf("Loading models from AVAILABLE_MODELS environment variable")
		// Parse JSON format: [{"id":"model-id","name":"Model Name","description":"Description"}]
		if err := json.Unmarshal([]byte(modelsConfig), &modelList); err != nil {
			log.Printf("Failed to parse AVAILABLE_MODELS, using defaults: %v", err)
			log.Printf("AVAILABLE_MODELS value: %s", modelsConfig)
			modelList = getDefaultModels()
		} else {
			log.Printf("Successfully loaded %d models from AVAILABLE_MODELS", len(modelList))
		}
	} else {
		log.Printf("AVAILABLE_MODELS not set, using default models")
		modelList = getDefaultModels()
	}

	c.JSON(http.StatusOK, modelList)
}

// getDefaultModels returns the default list of models
func getDefaultModels() []Model {
	return []Model{
		{ID: "gpt-4o", Name: "GPT-4o", Description: "Most capable model for complex tasks"},
		{ID: "gpt-4o-mini", Name: "GPT-4o Mini", Description: "Fast and efficient for simple tasks"},
		{ID: "claude-3-5-sonnet", Name: "Claude 3.5 Sonnet", Description: "High intelligence and speed"},
	}
}
