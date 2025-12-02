package api

import (
	"net/http"
	"time"
	"veritas-server/db"
	"veritas-server/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateConversation creates a new conversation
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

// GetConversations returns all conversations ordered by creation time
func GetConversations(c *gin.Context) {
	var convs []models.Conversation
	if err := db.DB.Order("created_at desc").Find(&convs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch conversations"})
		return
	}
	c.JSON(http.StatusOK, convs)
}

// GetConversation returns a specific conversation with its messages
func GetConversation(c *gin.Context) {
	id := c.Param("id")
	var conv models.Conversation
	if err := db.DB.Preload("Messages").First(&conv, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Conversation not found"})
		return
	}
	c.JSON(http.StatusOK, conv)
}
