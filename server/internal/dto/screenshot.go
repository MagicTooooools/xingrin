package dto

import "time"

// ScreenshotListQuery represents screenshot list query parameters
type ScreenshotListQuery struct {
	PaginationQuery
	Filter string `form:"filter"` // URL fuzzy search
}

// ScreenshotResponse represents screenshot response (without image data)
type ScreenshotResponse struct {
	ID         int       `json:"id"`
	URL        string    `json:"url"`
	StatusCode *int16    `json:"statusCode"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// ScreenshotItem represents a single screenshot for bulk upsert
type ScreenshotItem struct {
	URL        string `json:"url" binding:"required"`
	StatusCode *int16 `json:"statusCode"`
	Image      []byte `json:"image"` // Base64 encoded image data
}

// BulkUpsertScreenshotRequest represents bulk upsert request
type BulkUpsertScreenshotRequest struct {
	Screenshots []ScreenshotItem `json:"screenshots" binding:"required,min=1,max=5000,dive"`
}

// BulkUpsertScreenshotResponse represents bulk upsert response
type BulkUpsertScreenshotResponse struct {
	UpsertedCount int64 `json:"upsertedCount"`
}
