package dto

// AgentUpdateStatusRequest represents update scan status request from agent
type AgentUpdateStatusRequest struct {
	Status       string `json:"status" binding:"required,oneof=scheduled running completed failed cancelled"`
	ErrorMessage string `json:"errorMessage" binding:"omitempty"`
}

// AgentUpdateStatusResponse represents update scan status response
type AgentUpdateStatusResponse struct {
	Success bool `json:"success"`
}
