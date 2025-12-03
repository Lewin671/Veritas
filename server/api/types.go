package api

// Model represents an available LLM model
type Model struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ChatRequest represents a chat message request
type ChatRequest struct {
	ModelConfigID  string `json:"modelConfigId"` // ID of the model configuration to use
	Message        string `json:"message"`
	ConversationID string `json:"conversationId"`
}

// ChatResponse represents a chat message response
type ChatResponse struct {
	Response       string `json:"response"`
	ConversationID string `json:"conversationId"`
}
