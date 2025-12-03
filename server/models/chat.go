package models

import (
	"time"
)

type Conversation struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	Messages  []Message `gorm:"foreignKey:ConversationID" json:"messages"`
}

type Message struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ConversationID string    `json:"conversationId"`
	Role           string    `json:"role"`
	Content        string    `json:"content"`
	ModelConfigID  string    `json:"modelConfigId"` // Track which model was used
	CreatedAt      time.Time `json:"createdAt"`
}
