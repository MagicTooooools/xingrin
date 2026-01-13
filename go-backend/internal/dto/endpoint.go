package dto

import "time"

// EndpointListQuery represents endpoint list query parameters
type EndpointListQuery struct {
	PaginationQuery
	Filter string `form:"filter"`
}

// EndpointResponse represents endpoint response
type EndpointResponse struct {
	ID                int       `json:"id"`
	TargetID          int       `json:"targetId"`
	URL               string    `json:"url"`
	Host              string    `json:"host"`
	Location          string    `json:"location"`
	Title             string    `json:"title"`
	Webserver         string    `json:"webserver"`
	ContentType       string    `json:"contentType"`
	StatusCode        *int      `json:"statusCode"`
	ContentLength     *int      `json:"contentLength"`
	ResponseBody      string    `json:"responseBody"`
	Tech              []string  `json:"tech"`
	Vhost             *bool     `json:"vhost"`
	MatchedGFPatterns []string  `json:"matchedGfPatterns"`
	ResponseHeaders   string    `json:"responseHeaders"`
	CreatedAt         time.Time `json:"createdAt"`
}

// BulkCreateEndpointsRequest represents bulk create endpoints request
type BulkCreateEndpointsRequest struct {
	URLs []string `json:"urls" binding:"required,min=1,max=5000"`
}

// BulkCreateEndpointsResponse represents bulk create endpoints response
type BulkCreateEndpointsResponse struct {
	CreatedCount int `json:"createdCount"`
}

// EndpointUpsertItem represents a single endpoint for bulk upsert
type EndpointUpsertItem struct {
	URL               string   `json:"url" binding:"required,url"`
	Host              string   `json:"host"`
	Location          string   `json:"location"`
	Title             string   `json:"title"`
	Webserver         string   `json:"webserver"`
	ContentType       string   `json:"contentType"`
	StatusCode        *int     `json:"statusCode"`
	ContentLength     *int     `json:"contentLength"`
	ResponseBody      string   `json:"responseBody"`
	Tech              []string `json:"tech"`
	Vhost             *bool    `json:"vhost"`
	MatchedGFPatterns []string `json:"matchedGfPatterns"`
	ResponseHeaders   string   `json:"responseHeaders"`
}

// BulkUpsertEndpointsRequest represents bulk upsert endpoints request
type BulkUpsertEndpointsRequest struct {
	Endpoints []EndpointUpsertItem `json:"endpoints" binding:"required,min=1,max=5000,dive"`
}

// BulkUpsertEndpointsResponse represents bulk upsert endpoints response
type BulkUpsertEndpointsResponse struct {
	AffectedCount int64 `json:"affectedCount"`
}
