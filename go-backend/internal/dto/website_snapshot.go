package dto

import "time"

// WebsiteSnapshotItem represents a single website snapshot data for bulk upsert
type WebsiteSnapshotItem struct {
	URL             string   `json:"url" binding:"required,url"`
	Host            string   `json:"host"`
	Title           string   `json:"title"`
	StatusCode      *int     `json:"statusCode"`
	ContentLength   *int64   `json:"contentLength"`
	Location        string   `json:"location"`
	Webserver       string   `json:"webserver"`
	ContentType     string   `json:"contentType"`
	Tech            []string `json:"tech"`
	ResponseBody    string   `json:"responseBody"`
	Vhost           *bool    `json:"vhost"`
	ResponseHeaders string   `json:"responseHeaders"`
}

// BulkUpsertWebsiteSnapshotsRequest represents bulk upsert website snapshots request
type BulkUpsertWebsiteSnapshotsRequest struct {
	TargetID int                   `json:"targetId" binding:"required"`
	Websites []WebsiteSnapshotItem `json:"websites" binding:"required,min=1,max=5000,dive"`
}

// BulkUpsertWebsiteSnapshotsResponse represents bulk upsert website snapshots response
type BulkUpsertWebsiteSnapshotsResponse struct {
	SnapshotCount int `json:"snapshotCount"`
	AssetCount    int `json:"assetCount"`
}

// WebsiteSnapshotListQuery represents website snapshot list query parameters
type WebsiteSnapshotListQuery struct {
	PaginationQuery
	Filter string `form:"filter"`
}

// WebsiteSnapshotResponse represents website snapshot response
type WebsiteSnapshotResponse struct {
	ID              int       `json:"id"`
	ScanID          int       `json:"scanId"`
	URL             string    `json:"url"`
	Host            string    `json:"host"`
	Title           string    `json:"title"`
	StatusCode      *int      `json:"statusCode"`
	ContentLength   *int64    `json:"contentLength"`
	Location        string    `json:"location"`
	Webserver       string    `json:"webserver"`
	ContentType     string    `json:"contentType"`
	Tech            []string  `json:"tech"`
	ResponseBody    string    `json:"responseBody"`
	Vhost           *bool     `json:"vhost"`
	ResponseHeaders string    `json:"responseHeaders"`
	CreatedAt       time.Time `json:"createdAt"`
}
