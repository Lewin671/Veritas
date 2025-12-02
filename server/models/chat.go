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
	CreatedAt      time.Time `json:"createdAt"`
}
