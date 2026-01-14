package dto

import "time"

// HostPortSnapshotItem is an alias for HostPortItem used in snapshot operations
type HostPortSnapshotItem = HostPortItem

type BulkUpsertHostPortSnapshotsRequest struct {
	TargetID  int                    `json:"targetId" binding:"required"`
	HostPorts []HostPortSnapshotItem `json:"hostPorts" binding:"required,min=1,max=5000,dive"`
}

type BulkUpsertHostPortSnapshotsResponse struct {
	SnapshotCount int `json:"snapshotCount"`
	AssetCount    int `json:"assetCount"`
}

type HostPortSnapshotListQuery struct {
	PaginationQuery
	Filter string `form:"filter"`
}

type HostPortSnapshotResponse struct {
	ID        int       `json:"id"`
	ScanID    int       `json:"scanId"`
	Host      string    `json:"host"`
	IP        string    `json:"ip"`
	Port      int       `json:"port"`
	CreatedAt time.Time `json:"createdAt"`
}
