package dto

import "time"

// DirectorySnapshotItem is an alias for DirectoryUpsertItem
// Snapshot items and asset items have identical fields
type DirectorySnapshotItem = DirectoryUpsertItem

// BulkUpsertDirectorySnapshotsRequest represents bulk upsert directory snapshots request
type BulkUpsertDirectorySnapshotsRequest struct {
	TargetID    int                     `json:"targetId" binding:"required"`
	Directories []DirectorySnapshotItem `json:"directories" binding:"required,min=1,max=5000,dive"`
}

// BulkUpsertDirectorySnapshotsResponse represents bulk upsert directory snapshots response
type BulkUpsertDirectorySnapshotsResponse struct {
	SnapshotCount int `json:"snapshotCount"`
	AssetCount    int `json:"assetCount"`
}

// DirectorySnapshotListQuery represents directory snapshot list query parameters
type DirectorySnapshotListQuery struct {
	PaginationQuery
	Filter string `form:"filter"`
}

// DirectorySnapshotResponse represents directory snapshot response
type DirectorySnapshotResponse struct {
	ID            int       `json:"id"`
	ScanID        int       `json:"scanId"`
	URL           string    `json:"url"`
	Status        *int      `json:"status"`
	ContentLength *int      `json:"contentLength"`
	ContentType   string    `json:"contentType"`
	Duration      *int      `json:"duration"`
	CreatedAt     time.Time `json:"createdAt"`
}
