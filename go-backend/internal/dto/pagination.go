package dto

// PaginationQuery represents pagination query parameters
type PaginationQuery struct {
	Page     int `form:"page" binding:"omitempty,min=1"`
	PageSize int `form:"pageSize" binding:"omitempty,min=1,max=100"`
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
	if p.PageSize > 100 {
		return 100
	}
	return p.PageSize
}

// GetOffset returns offset for database query
func (p *PaginationQuery) GetOffset() int {
	return (p.GetPage() - 1) * p.GetPageSize()
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Data     interface{} `json:"data"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

// NewPaginatedResponse creates a new paginated response
func NewPaginatedResponse(data interface{}, total int64, page, pageSize int) *PaginatedResponse {
	return &PaginatedResponse{
		Data:     data,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
}
