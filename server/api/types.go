package api

// Model represents an available LLM model
type Model struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ChatRequest represents a chat message request
type ChatRequest struct {
	ModelID        string `json:"modelId"`
	Message        string `json:"message"`
	ConversationID string `json:"conversationId"`
}

// ChatResponse represents a chat message response
type ChatResponse struct {
	Response string `json:"response"`
}
