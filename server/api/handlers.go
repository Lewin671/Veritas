package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
	"veritas-server/db"
	"veritas-server/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	openai "github.com/sashabaranov/go-openai"
)

type Model struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ChatRequest struct {
	ModelID        string `json:"modelId"`
	Message        string `json:"message"`
	ConversationID string `json:"conversationId"`
}

type ChatResponse struct {
	Response string `json:"response"`
}

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

func getDefaultModels() []Model {
	return []Model{
		{ID: "gpt-4o", Name: "GPT-4o", Description: "Most capable model for complex tasks"},
		{ID: "gpt-4o-mini", Name: "GPT-4o Mini", Description: "Fast and efficient for simple tasks"},
		{ID: "claude-3-5-sonnet", Name: "Claude 3.5 Sonnet", Description: "High intelligence and speed"},
	}
}

func CreateConversation(c *gin.Context) {
	conv := models.Conversation{
		ID:        uuid.New().String(),
		Title:     "New Chat",
		CreatedAt: time.Now(),
	}
	if err := db.DB.Create(&conv).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create conversation"})
		return
	}
	c.JSON(http.StatusOK, conv)
}

func GetConversations(c *gin.Context) {
	var convs []models.Conversation
	if err := db.DB.Order("created_at desc").Find(&convs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch conversations"})
		return
	}
	c.JSON(http.StatusOK, convs)
}

func GetConversation(c *gin.Context) {
	id := c.Param("id")
	var conv models.Conversation
	if err := db.DB.Preload("Messages").First(&conv, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Conversation not found"})
		return
	}
	c.JSON(http.StatusOK, conv)
}

func Chat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create conversation if not provided
	if req.ConversationID == "" {
		title := "New Chat"
		if len(req.Message) > 0 {
			if len(req.Message) > 30 {
				title = req.Message[:30] + "..."
			} else {
				title = req.Message
			}
		}

		conv := models.Conversation{
			ID:        uuid.New().String(),
			Title:     title,
			CreatedAt: time.Now(),
		}
		if err := db.DB.Create(&conv).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create conversation"})
			return
		}
		req.ConversationID = conv.ID
	}

	// Save user message
	userMsg := models.Message{
		ConversationID: req.ConversationID,
		Role:           "user",
		Content:        req.Message,
		CreatedAt:      time.Now(),
	}
	if err := db.DB.Create(&userMsg).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save message"})
		return
	}

	// Call OpenAI API
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "OPENAI_API_KEY is not configured"})
		return
	}

	// Support OpenAI compatible APIs with custom base URL
	var client *openai.Client
	baseURL := os.Getenv("OPENAI_BASE_URL")
	if baseURL != "" {
		log.Printf("Using custom base URL: %s", baseURL)
		config := openai.DefaultConfig(apiKey)
		config.BaseURL = baseURL
		client = openai.NewClientWithConfig(config)
	} else {
		log.Printf("Using default OpenAI base URL: https://api.openai.com/v1")
		client = openai.NewClient(apiKey)
	}

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: req.ModelID,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: req.Message,
				},
			},
		},
	)

	var responseContent string
	if err != nil {
		log.Printf("OpenAI API error: %v", err)
		responseContent = "Error: Failed to get response from LLM provider. " + err.Error()
	} else {
		responseContent = resp.Choices[0].Message.Content
	}

	// Save assistant message
	assistantMsg := models.Message{
		ConversationID: req.ConversationID,
		Role:           "assistant",
		Content:        responseContent,
		CreatedAt:      time.Now(),
	}
	if err := db.DB.Create(&assistantMsg).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save response"})
		return
	}

	c.JSON(http.StatusOK, ChatResponse{
		Response: responseContent,
	})
}
