package api

import (
	"context"
	"net/http"
	"time"
	"veritas-server/db"
	"veritas-server/models"
	"veritas-server/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openai/openai-go"
)

// ModelConfigRequest represents the request body for creating/updating model configs
type ModelConfigRequest struct {
	Name      string `json:"name" binding:"required"`
	Provider  string `json:"provider" binding:"required"`
	BaseURL   string `json:"baseUrl"`
	ModelID   string `json:"modelId" binding:"required"`
	APIKey    string `json:"apiKey" binding:"required"`
	IsDefault bool   `json:"isDefault"`
}

// CreateModelConfig creates a new model configuration
func CreateModelConfig(c *gin.Context) {
	var req ModelConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Encrypt API key
	encryptedKey, err := services.EncryptAPIKey(req.APIKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt API key"})
		return
	}

	// If this is set as default, unset other defaults
	if req.IsDefault {
		if err := db.DB.Model(&models.ModelConfig{}).Where("is_default = ?", true).Update("is_default", false).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update default status"})
			return
		}
	}

	config := models.ModelConfig{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Provider:  req.Provider,
		BaseURL:   req.BaseURL,
		ModelID:   req.ModelID,
		APIKey:    encryptedKey,
		IsDefault: req.IsDefault,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.DB.Create(&config).Error; err != nil {
		// Check for duplicate name
		if err.Error() == "UNIQUE constraint failed: model_configs.name" ||
			err.Error() == "duplicate key value violates unique constraint \"model_configs_name_key\"" {
			c.JSON(http.StatusConflict, gin.H{"error": "A configuration with this name already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create configuration"})
		return
	}

	c.JSON(http.StatusCreated, config.ToResponse())
}

// GetModelConfigs returns all model configurations (with masked API keys)
func GetModelConfigs(c *gin.Context) {
	var configs []models.ModelConfig
	if err := db.DB.Find(&configs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve configurations"})
		return
	}

	responses := make([]models.ModelConfigResponse, len(configs))
	for i, config := range configs {
		responses[i] = config.ToResponse()
	}

	c.JSON(http.StatusOK, responses)
}

// GetModelConfig returns a specific model configuration by ID
func GetModelConfig(c *gin.Context) {
	id := c.Param("id")

	var config models.ModelConfig
	if err := db.DB.First(&config, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Configuration not found"})
		return
	}

	c.JSON(http.StatusOK, config.ToResponse())
}

// UpdateModelConfig updates an existing model configuration
func UpdateModelConfig(c *gin.Context) {
	id := c.Param("id")

	var req ModelConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	var config models.ModelConfig
	if err := db.DB.First(&config, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Configuration not found"})
		return
	}

	// Encrypt new API key
	encryptedKey, err := services.EncryptAPIKey(req.APIKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt API key"})
		return
	}

	// If this is set as default, unset other defaults
	if req.IsDefault && !config.IsDefault {
		if err := db.DB.Model(&models.ModelConfig{}).Where("is_default = ? AND id != ?", true, id).Update("is_default", false).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update default status"})
			return
		}
	}

	// Update fields
	config.Name = req.Name
	config.Provider = req.Provider
	config.BaseURL = req.BaseURL
	config.ModelID = req.ModelID
	config.APIKey = encryptedKey
	config.IsDefault = req.IsDefault
	config.UpdatedAt = time.Now()

	if err := db.DB.Save(&config).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update configuration"})
		return
	}

	c.JSON(http.StatusOK, config.ToResponse())
}

// DeleteModelConfig deletes a model configuration
func DeleteModelConfig(c *gin.Context) {
	id := c.Param("id")

	var config models.ModelConfig
	if err := db.DB.First(&config, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Configuration not found"})
		return
	}

	// Check if config is referenced by any messages
	var messageCount int64
	if err := db.DB.Model(&models.Message{}).Where("model_config_id = ?", id).Count(&messageCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check configuration usage"})
		return
	}

	if messageCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete configuration that is referenced by messages"})
		return
	}

	if err := db.DB.Delete(&config).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete configuration"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Configuration deleted successfully"})
}

// TestModelConfigRequest represents the request body for testing a model config
type TestModelConfigRequest struct {
	BaseURL string `json:"baseUrl"`
	ModelID string `json:"modelId" binding:"required"`
	APIKey  string `json:"apiKey" binding:"required"`
}

// TestModelConfigResponse represents the response for testing a model config
type TestModelConfigResponse struct {
	Success      bool              `json:"success"`
	Message      string            `json:"message"`
	Details      map[string]string `json:"details,omitempty"`
	ErrorDetails string            `json:"errorDetails,omitempty"`
}

// TestModelConfig tests a model configuration without saving
func TestModelConfig(c *gin.Context) {
	var req TestModelConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Create a temporary model config for testing
	tempConfig := &models.ModelConfig{
		BaseURL: req.BaseURL,
		ModelID: req.ModelID,
		APIKey:  req.APIKey, // Use plaintext for testing (not encrypted)
	}

	// Test the connection with timeout
	startTime := time.Now()
	client, err := createLLMClientFromConfig(tempConfig, false) // false = don't decrypt
	if err != nil {
		c.JSON(http.StatusOK, TestModelConfigResponse{
			Success:      false,
			Message:      "Failed to create LLM client",
			ErrorDetails: err.Error(),
		})
		return
	}

	// Try a simple completion to test the connection
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	_, err = client.Chat.Completions.New(
		ctx,
		openai.ChatCompletionNewParams{
			Model: openai.ChatModel(req.ModelID), //nolint:unconvert
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.UserMessage("Hello"),
			},
			MaxTokens: openai.Int(10),
		},
	)

	responseTime := time.Since(startTime)

	if err != nil {
		// Determine error type
		errorMsg := err.Error()
		statusCode := http.StatusOK // Still return 200, but with success: false

		if ctx.Err() != nil {
			c.JSON(statusCode, TestModelConfigResponse{
				Success:      false,
				Message:      "Connection timeout",
				ErrorDetails: "Request timed out after 30 seconds",
			})
			return
		}

		// Check for common error patterns
		if contains(errorMsg, "401") || contains(errorMsg, "unauthorized") || contains(errorMsg, "invalid_api_key") {
			c.JSON(statusCode, TestModelConfigResponse{
				Success:      false,
				Message:      "Authentication failed",
				ErrorDetails: "Invalid API key",
			})
			return
		}

		if contains(errorMsg, "404") || contains(errorMsg, "model_not_found") {
			c.JSON(statusCode, TestModelConfigResponse{
				Success:      false,
				Message:      "Invalid model ID",
				ErrorDetails: "The specified model does not exist",
			})
			return
		}

		if contains(errorMsg, "429") || contains(errorMsg, "rate_limit") {
			c.JSON(statusCode, TestModelConfigResponse{
				Success:      false,
				Message:      "Rate limit exceeded",
				ErrorDetails: "Too many requests, please try again later",
			})
			return
		}

		c.JSON(statusCode, TestModelConfigResponse{
			Success:      false,
			Message:      "Connection failed",
			ErrorDetails: errorMsg,
		})
		return
	}

	c.JSON(http.StatusOK, TestModelConfigResponse{
		Success: true,
		Message: "Connection successful",
		Details: map[string]string{
			"responseTime":   responseTime.String(),
			"modelAvailable": "true",
		},
	})
}

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
