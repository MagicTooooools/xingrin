package dto

import "time"

// WebsiteListQuery represents website list query parameters
type WebsiteListQuery struct {
	PaginationQuery
	Filter string `form:"filter"`
}

// WebsiteResponse represents website response
type WebsiteResponse struct {
	ID              int       `json:"id"`
	URL             string    `json:"url"`
	Host            string    `json:"host"`
	Location        string    `json:"location"`
	Title           string    `json:"title"`
	Webserver       string    `json:"webserver"`
	ContentType     string    `json:"contentType"`
	StatusCode      *int      `json:"statusCode"`
	ContentLength   *int      `json:"contentLength"`
	ResponseBody    string    `json:"responseBody"`
	Tech            []string  `json:"tech"`
	Vhost           *bool     `json:"vhost"`
	ResponseHeaders string    `json:"responseHeaders"`
	CreatedAt       time.Time `json:"createdAt"`
}

// BulkCreateWebsitesRequest represents bulk create websites request
type BulkCreateWebsitesRequest struct {
	URLs []string `json:"urls" binding:"required,min=1,max=5000,dive,url"`
}

// BulkCreateWebsitesResponse represents bulk create websites response
type BulkCreateWebsitesResponse struct {
	CreatedCount int `json:"createdCount"`
}

// WebsiteUpsertItem represents a single website for upsert operation
type WebsiteUpsertItem struct {
	URL             string   `json:"url" binding:"required,url"`
	Host            string   `json:"host"`
	Location        string   `json:"location"`
	Title           string   `json:"title"`
	Webserver       string   `json:"webserver"`
	ContentType     string   `json:"contentType"`
	StatusCode      *int     `json:"statusCode"`
	ContentLength   *int     `json:"contentLength"`
	ResponseBody    string   `json:"responseBody"`
	Tech            []string `json:"tech"`
	Vhost           *bool    `json:"vhost"`
	ResponseHeaders string   `json:"responseHeaders"`
}

// BulkUpsertWebsitesRequest represents bulk upsert websites request (for scanner import)
type BulkUpsertWebsitesRequest struct {
	Websites []WebsiteUpsertItem `json:"websites" binding:"required,min=1,max=5000,dive"`
}

// BulkUpsertWebsitesResponse represents bulk upsert websites response
type BulkUpsertWebsitesResponse struct {
	UpsertedCount int `json:"upsertedCount"`
}
