package dto

import "time"

// SubdomainListQuery represents subdomain list query parameters
type SubdomainListQuery struct {
	PaginationQuery
	Filter string `form:"filter"`
}

// SubdomainResponse represents subdomain response
type SubdomainResponse struct {
	ID        int       `json:"id"`
	TargetID  int       `json:"targetId"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

// BulkCreateSubdomainsRequest represents bulk create subdomains request
type BulkCreateSubdomainsRequest struct {
	Names []string `json:"names" binding:"required,min=1,max=5000"`
}

// BulkCreateSubdomainsResponse represents bulk create subdomains response
type BulkCreateSubdomainsResponse struct {
	CreatedCount int `json:"createdCount"`
}
