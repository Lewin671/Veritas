package services

import (
	"log"
	"os"
	"time"
	"veritas-server/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MigrateDefaultModelConfig creates a default model configuration from environment variables
// if no configurations exist in the database
func MigrateDefaultModelConfig(db *gorm.DB) error {
	// Check if any model configs exist
	var count int64
	if err := db.Model(&models.ModelConfig{}).Count(&count).Error; err != nil {
		return err
	}

	// If configs already exist, skip migration
	if count > 0 {
		log.Println("Model configurations already exist, skipping default config migration")
		return nil
	}

	// Check if environment variables are set
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Println("OPENAI_API_KEY not set, skipping default config migration")
		return nil
	}

	log.Println("Creating default model configuration from environment variables")

	// Encrypt API key
	encryptedKey, err := EncryptAPIKey(apiKey)
	if err != nil {
		return err
	}

	// Get base URL if set
	baseURL := os.Getenv("OPENAI_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}

	// Create default configuration
	config := models.ModelConfig{
		ID:        uuid.New().String(),
		Name:      "Default OpenAI",
		Provider:  "openai",
		BaseURL:   baseURL,
		ModelID:   "gpt-4o-mini", // Use a reasonable default
		APIKey:    encryptedKey,
		IsDefault: true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.Create(&config).Error; err != nil {
		return err
	}

	log.Printf("Default model configuration created: %s (ID: %s)", config.Name, config.ID)
	return nil
}
