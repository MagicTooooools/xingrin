package dto

import "time"

// IPAddressListQuery represents IP address list query parameters
type IPAddressListQuery struct {
	PaginationQuery
	Filter string `form:"filter"`
}

// IPAddressResponse represents aggregated IP address response (grouped by IP)
type IPAddressResponse struct {
	IP        string    `json:"ip"`
	Hosts     []string  `json:"hosts"`
	Ports     []int     `json:"ports"`
	CreatedAt time.Time `json:"createdAt"`
}

// IPAddressItem represents a single IP address mapping for bulk operations
type IPAddressItem struct {
	Host string `json:"host" binding:"required"`
	IP   string `json:"ip" binding:"required,ip"`
	Port int    `json:"port" binding:"required,min=1,max=65535"`
}

// BulkUpsertIPAddressesRequest represents bulk upsert request (for scanner import)
type BulkUpsertIPAddressesRequest struct {
	Mappings []IPAddressItem `json:"mappings" binding:"required,min=1,max=5000,dive"`
}

// BulkUpsertIPAddressesResponse represents bulk upsert response
type BulkUpsertIPAddressesResponse struct {
	UpsertedCount int `json:"upsertedCount"`
}

// BulkDeleteIPAddressesRequest represents bulk delete request (by IP list)
type BulkDeleteIPAddressesRequest struct {
	IPs []string `json:"ips" binding:"required,min=1"`
}

// BulkDeleteIPAddressesResponse represents bulk delete response
type BulkDeleteIPAddressesResponse struct {
	DeletedCount int64 `json:"deletedCount"`
}
