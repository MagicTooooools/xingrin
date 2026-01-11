package dto

import "time"

// CreateTargetRequest represents create target request
type CreateTargetRequest struct {
	Name string `json:"name" binding:"required,max=300"`
	Type string `json:"type" binding:"omitempty,oneof=domain ip cidr"`
}

// UpdateTargetRequest represents update target request
type UpdateTargetRequest struct {
	Name string `json:"name" binding:"required,max=300"`
	Type string `json:"type" binding:"omitempty,oneof=domain ip cidr"`
}

// TargetListQuery represents target list query parameters
type TargetListQuery struct {
	PaginationQuery
	Type   string `form:"type" binding:"omitempty,oneof=domain ip cidr"`
	Search string `form:"search"`
}

// TargetResponse represents target response
type TargetResponse struct {
	ID            int        `json:"id"`
	Name          string     `json:"name"`
	Type          string     `json:"type"`
	CreatedAt     time.Time  `json:"createdAt"`
	LastScannedAt *time.Time `json:"lastScannedAt"`
}
