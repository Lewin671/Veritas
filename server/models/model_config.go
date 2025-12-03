package models

import (
	"time"
)

// ModelConfig represents a configured LLM model with connection details
type ModelConfig struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null;uniqueIndex" json:"name"`
	Provider  string    `gorm:"not null" json:"provider"`
	BaseURL   string    `json:"baseUrl"`
	ModelID   string    `gorm:"not null" json:"modelId"`
	APIKey    string    `gorm:"not null" json:"-"` // Encrypted, never sent to client
	IsDefault bool      `gorm:"default:false" json:"isDefault"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ModelConfigResponse is the sanitized version sent to clients
type ModelConfigResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Provider  string    `json:"provider"`
	BaseURL   string    `json:"baseUrl"`
	ModelID   string    `json:"modelId"`
	IsDefault bool      `json:"isDefault"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ToResponse converts ModelConfig to ModelConfigResponse (masks API key)
func (m *ModelConfig) ToResponse() ModelConfigResponse {
	return ModelConfigResponse{
		ID:        m.ID,
		Name:      m.Name,
		Provider:  m.Provider,
		BaseURL:   m.BaseURL,
		ModelID:   m.ModelID,
		IsDefault: m.IsDefault,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
