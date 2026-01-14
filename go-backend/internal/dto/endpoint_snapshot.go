package dto

import "time"

// EndpointSnapshotItem represents a single endpoint snapshot data for bulk upsert
type EndpointSnapshotItem struct {
	URL               string   `json:"url" binding:"required,url"`
	Host              string   `json:"host"`
	Title             string   `json:"title"`
	StatusCode        *int     `json:"statusCode"`
	ContentLength     *int     `json:"contentLength"`
	Location          string   `json:"location"`
	Webserver         string   `json:"webserver"`
	ContentType       string   `json:"contentType"`
	Tech              []string `json:"tech"`
	ResponseBody      string   `json:"responseBody"`
	Vhost             *bool    `json:"vhost"`
	MatchedGFPatterns []string `json:"matchedGfPatterns"`
	ResponseHeaders   string   `json:"responseHeaders"`
}

// BulkUpsertEndpointSnapshotsRequest represents bulk upsert endpoint snapshots request
type BulkUpsertEndpointSnapshotsRequest struct {
	TargetID  int                    `json:"targetId" binding:"required"`
	Endpoints []EndpointSnapshotItem `json:"endpoints" binding:"required,min=1,max=5000,dive"`
}

// BulkUpsertEndpointSnapshotsResponse represents bulk upsert endpoint snapshots response
type BulkUpsertEndpointSnapshotsResponse struct {
	SnapshotCount int `json:"snapshotCount"`
	AssetCount    int `json:"assetCount"`
}

// EndpointSnapshotListQuery represents endpoint snapshot list query parameters
type EndpointSnapshotListQuery struct {
	PaginationQuery
	Filter string `form:"filter"`
}

// EndpointSnapshotResponse represents endpoint snapshot response
type EndpointSnapshotResponse struct {
	ID                int       `json:"id"`
	ScanID            int       `json:"scanId"`
	URL               string    `json:"url"`
	Host              string    `json:"host"`
	Title             string    `json:"title"`
	StatusCode        *int      `json:"statusCode"`
	ContentLength     *int      `json:"contentLength"`
	Location          string    `json:"location"`
	Webserver         string    `json:"webserver"`
	ContentType       string    `json:"contentType"`
	Tech              []string  `json:"tech"`
	ResponseBody      string    `json:"responseBody"`
	Vhost             *bool     `json:"vhost"`
	MatchedGFPatterns []string  `json:"matchedGfPatterns"`
	ResponseHeaders   string    `json:"responseHeaders"`
	CreatedAt         time.Time `json:"createdAt"`
}
