package dto

import "time"

// CreateEngineRequest represents create engine request
type CreateEngineRequest struct {
	Name          string `json:"name" binding:"required,max=200"`
	Configuration string `json:"configuration"`
}

// UpdateEngineRequest represents update engine request (PUT - full update)
type UpdateEngineRequest struct {
	Name          string `json:"name" binding:"required,max=200"`
	Configuration string `json:"configuration"`
}

// PatchEngineRequest represents patch engine request (PATCH - partial update)
type PatchEngineRequest struct {
	Name          *string `json:"name,omitempty" binding:"omitempty,max=200"`
	Configuration *string `json:"configuration,omitempty"`
}

// EngineResponse represents engine response
type EngineResponse struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	Configuration string    `json:"configuration"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
