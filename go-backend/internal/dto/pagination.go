package dto

import "github.com/gin-gonic/gin"

// PaginationQuery represents pagination query parameters
type PaginationQuery struct {
	Page     int `form:"page" binding:"omitempty,min=1"`
	PageSize int `form:"pageSize" binding:"omitempty,min=1,max=1000"`
}

// GetPage returns page number with default
func (p *PaginationQuery) GetPage() int {
	if p.Page <= 0 {
		return 1
	}
	return p.Page
}

// GetPageSize returns page size with default
func (p *PaginationQuery) GetPageSize() int {
	if p.PageSize <= 0 {
		return 20
	}
	if p.PageSize > 1000 {
		return 1000
	}
	return p.PageSize
}

// PaginatedResponse represents a paginated response (matches Python format)
type PaginatedResponse[T any] struct {
	Results    []T   `json:"results"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	TotalPages int   `json:"totalPages"`
}

// NewPaginatedResponse creates a new paginated response
func NewPaginatedResponse[T any](data []T, total int64, page, pageSize int) *PaginatedResponse[T] {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &PaginatedResponse[T]{
		Results:    data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}

// Paginated sends a paginated response
func Paginated[T any](c *gin.Context, data []T, total int64, page, pageSize int) {
	Success(c, NewPaginatedResponse(data, total, page, pageSize))
}
