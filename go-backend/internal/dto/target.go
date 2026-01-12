package dto

import "time"

// CreateTargetRequest represents create target request
type CreateTargetRequest struct {
	Name string `json:"name" binding:"required,max=300"`
}

// UpdateTargetRequest represents update target request
type UpdateTargetRequest struct {
	Name string `json:"name" binding:"required,max=300"`
}

// TargetListQuery represents target list query parameters
type TargetListQuery struct {
	PaginationQuery
	Type   string `form:"type" binding:"omitempty,oneof=domain ip cidr"`
	Search string `form:"search"`
}

// TargetResponse represents target response
type TargetResponse struct {
	ID            int                 `json:"id"`
	Name          string              `json:"name"`
	Type          string              `json:"type"`
	CreatedAt     time.Time           `json:"createdAt"`
	LastScannedAt *time.Time          `json:"lastScannedAt"`
	Organizations []OrganizationBrief `json:"organizations,omitempty"`
}

// OrganizationBrief represents brief organization info for target response
type OrganizationBrief struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// BatchCreateTargetRequest represents batch create targets request
type BatchCreateTargetRequest struct {
	Targets        []TargetItem `json:"targets" binding:"required,min=1,max=5000,dive"`
	OrganizationID *int         `json:"organizationId"`
}

// TargetItem represents a single target in batch create
type TargetItem struct {
	Name string `json:"name" binding:"required,max=300"`
}

// BatchCreateTargetResponse represents batch create targets response
type BatchCreateTargetResponse struct {
	CreatedCount  int            `json:"createdCount"`
	FailedCount   int            `json:"failedCount"`
	FailedTargets []FailedTarget `json:"failedTargets"`
	Message       string         `json:"message"`
}

// FailedTarget represents a failed target in batch create
type FailedTarget struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

// BulkDeleteRequest represents bulk delete request
type BulkDeleteRequest struct {
	IDs []int `json:"ids" binding:"required,min=1"`
}

// BulkDeleteResponse represents bulk delete response
type BulkDeleteResponse struct {
	DeletedCount int64 `json:"deletedCount"`
}
