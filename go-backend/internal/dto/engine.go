package dto

import "time"

// CreateEngineRequest represents create engine request
type CreateEngineRequest struct {
	Name          string `json:"name" binding:"required,max=200"`
	Configuration string `json:"configuration"`
}

// UpdateEngineRequest represents update engine request
type UpdateEngineRequest struct {
	Name          string `json:"name" binding:"required,max=200"`
	Configuration string `json:"configuration"`
}

// EngineResponse represents engine response
type EngineResponse struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	Configuration string    `json:"configuration"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
