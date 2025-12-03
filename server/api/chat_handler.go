package api

import (
	"context"
	"log"
	"net/http"
	"time"
	"veritas-server/db"
	"veritas-server/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openai/openai-go"
)

// Chat handles chat requests, creates conversations if needed, and interacts with LLM
func Chat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create conversation if not provided
	conversationID, err := ensureConversation(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create conversation"})
		return
	}
	req.ConversationID = conversationID

	// Save user message
	if err := saveUserMessage(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save message"})
		return
	}

	// Get LLM response
	responseContent := getLLMResponse(req)

	// Save assistant message
	if err := saveAssistantMessage(req.ConversationID, responseContent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save response"})
		return
	}

	c.JSON(http.StatusOK, ChatResponse{
		Response:       responseContent,
		ConversationID: req.ConversationID,
	})
}

// ensureConversation creates a conversation if one doesn't exist
func ensureConversation(req ChatRequest) (string, error) {
	if req.ConversationID != "" {
		return req.ConversationID, nil
	}

	title := generateConversationTitle(req.Message)
	conv := models.Conversation{
		ID:        uuid.New().String(),
		Title:     title,
		CreatedAt: time.Now(),
	}

	if err := db.DB.Create(&conv).Error; err != nil {
		return "", err
	}

	return conv.ID, nil
}

// generateConversationTitle creates a title from the first message
func generateConversationTitle(message string) string {
	title := "New Chat"
	if len(message) > 0 {
		if len(message) > 30 {
			title = message[:30] + "..."
		} else {
			title = message
		}
	}
	return title
}

// saveUserMessage saves the user's message to the database
func saveUserMessage(req ChatRequest) error {
	userMsg := models.Message{
		ConversationID: req.ConversationID,
		Role:           "user",
		Content:        req.Message,
		CreatedAt:      time.Now(),
	}
	return db.DB.Create(&userMsg).Error
}

// saveAssistantMessage saves the assistant's response to the database
func saveAssistantMessage(conversationID, content string) error {
	assistantMsg := models.Message{
		ConversationID: conversationID,
		Role:           "assistant",
		Content:        content,
		CreatedAt:      time.Now(),
	}
	return db.DB.Create(&assistantMsg).Error
}

// getLLMResponse calls the LLM API with full conversation history and returns the response
func getLLMResponse(req ChatRequest) string {
	client := createLLMClient()
	if client == nil {
		return "Error: OPENAI_API_KEY is not configured"
	}

	// Load full conversation history so the model has memory
	var history []models.Message
	if err := db.DB.
		Where("conversation_id = ?", req.ConversationID).
		Order("created_at asc").
		Find(&history).Error; err != nil {
		log.Printf("Failed to load conversation history: %v", err)
	}

	var chatMessages []openai.ChatCompletionMessageParamUnion

	// Convert stored messages into OpenAI chat messages
	for _, m := range history {
		switch m.Role {
		case "user":
			chatMessages = append(chatMessages, openai.UserMessage(m.Content))
		case "assistant":
			chatMessages = append(chatMessages, openai.AssistantMessage(m.Content))
		default:
			// Ignore unknown roles for now
		}
	}

	// Fallback: if for some reason we have no history, at least send the current message
	if len(chatMessages) == 0 {
		chatMessages = append(chatMessages, openai.UserMessage(req.Message))
	}

	resp, err := client.Chat.Completions.New(
		context.Background(),
		openai.ChatCompletionNewParams{
			Model:    req.ModelID,
			Messages: chatMessages,
		},
	)

	if err != nil {
		log.Printf("OpenAI API error: %v", err)
		return "Error: Failed to get response from LLM provider. " + err.Error()
	}

	return resp.Choices[0].Message.Content
}
