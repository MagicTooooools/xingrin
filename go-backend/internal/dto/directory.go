package dto

import "time"

// DirectoryListQuery represents directory list query parameters
type DirectoryListQuery struct {
	PaginationQuery
	Filter string `form:"filter"`
}

// DirectoryResponse represents directory response
type DirectoryResponse struct {
	ID            int       `json:"id"`
	TargetID      int       `json:"targetId"`
	URL           string    `json:"url"`
	Status        *int      `json:"status"`
	ContentLength *int64    `json:"contentLength"`
	ContentType   string    `json:"contentType"`
	Duration      *int64    `json:"duration"`
	CreatedAt     time.Time `json:"createdAt"`
}

// BulkCreateDirectoriesRequest represents bulk create directories request
type BulkCreateDirectoriesRequest struct {
	URLs []string `json:"urls" binding:"required,min=1,max=5000"`
}

// BulkCreateDirectoriesResponse represents bulk create directories response
type BulkCreateDirectoriesResponse struct {
	CreatedCount int `json:"createdCount"`
}

// DirectoryUpsertItem represents a single directory for bulk upsert
type DirectoryUpsertItem struct {
	URL           string `json:"url" binding:"required,url"`
	Status        *int   `json:"status"`
	ContentLength *int64 `json:"contentLength"`
	ContentType   string `json:"contentType"`
	Duration      *int64 `json:"duration"`
}

// BulkUpsertDirectoriesRequest represents bulk upsert directories request
type BulkUpsertDirectoriesRequest struct {
	Directories []DirectoryUpsertItem `json:"directories" binding:"required,min=1,max=5000,dive"`
}

// BulkUpsertDirectoriesResponse represents bulk upsert directories response
type BulkUpsertDirectoriesResponse struct {
	AffectedCount int64 `json:"affectedCount"`
}
