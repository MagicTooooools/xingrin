package dto

import "time"

// HostPortListQuery represents host-port list query parameters
type HostPortListQuery struct {
	PaginationQuery
	Filter string `form:"filter"`
}

// HostPortResponse represents aggregated host-port response (grouped by IP)
type HostPortResponse struct {
	IP        string    `json:"ip"`
	Hosts     []string  `json:"hosts"`
	Ports     []int     `json:"ports"`
	CreatedAt time.Time `json:"createdAt"`
}

// HostPortItem represents a single host-port mapping for bulk operations
type HostPortItem struct {
	Host string `json:"host" binding:"required"`
	IP   string `json:"ip" binding:"required,ip"`
	Port int    `json:"port" binding:"required,min=1,max=65535"`
}

// BulkUpsertHostPortsRequest represents bulk upsert request (for scanner import)
type BulkUpsertHostPortsRequest struct {
	Mappings []HostPortItem `json:"mappings" binding:"required,min=1,max=5000,dive"`
}

// BulkUpsertHostPortsResponse represents bulk upsert response
type BulkUpsertHostPortsResponse struct {
	UpsertedCount int `json:"upsertedCount"`
}

// BulkDeleteHostPortsRequest represents bulk delete request (by IP list)
type BulkDeleteHostPortsRequest struct {
	IPs []string `json:"ips" binding:"required,min=1"`
}

// BulkDeleteHostPortsResponse represents bulk delete response
type BulkDeleteHostPortsResponse struct {
	DeletedCount int64 `json:"deletedCount"`
}
