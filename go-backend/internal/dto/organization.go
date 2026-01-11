package dto

import "time"

// CreateOrganizationRequest represents create organization request
type CreateOrganizationRequest struct {
	Name        string `json:"name" binding:"required,max=300"`
	Description string `json:"description" binding:"max=1000"`
}

// UpdateOrganizationRequest represents update organization request
type UpdateOrganizationRequest struct {
	Name        string `json:"name" binding:"required,max=300"`
	Description string `json:"description" binding:"max=1000"`
}

// OrganizationResponse represents organization response
type OrganizationResponse struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}
