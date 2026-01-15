package dto

import "time"

// ScreenshotSnapshotItem represents a single screenshot snapshot data for bulk upsert
// Image is expected to be base64 encoded in JSON and will be decoded into []byte automatically.
type ScreenshotSnapshotItem struct {
	URL        string `json:"url" binding:"required,url"`
	StatusCode *int16 `json:"statusCode"`
	Image      []byte `json:"image"`
}

// BulkUpsertScreenshotSnapshotsRequest represents bulk upsert screenshot snapshots request
type BulkUpsertScreenshotSnapshotsRequest struct {
	TargetID    int                      `json:"targetId" binding:"required"`
	Screenshots []ScreenshotSnapshotItem `json:"screenshots" binding:"required,min=1,max=5000,dive"`
}

// BulkUpsertScreenshotSnapshotsResponse represents bulk upsert screenshot snapshots response
type BulkUpsertScreenshotSnapshotsResponse struct {
	SnapshotCount int `json:"snapshotCount"`
	AssetCount    int `json:"assetCount"`
}

// ScreenshotSnapshotListQuery represents screenshot snapshot list query parameters
type ScreenshotSnapshotListQuery struct {
	PaginationQuery
	Filter string `form:"filter"`
}

// ScreenshotSnapshotResponse represents screenshot snapshot response (without image data)
type ScreenshotSnapshotResponse struct {
	ID         int       `json:"id"`
	ScanID     int       `json:"scanId"`
	URL        string    `json:"url"`
	StatusCode *int16    `json:"statusCode"`
	CreatedAt  time.Time `json:"createdAt"`
}
