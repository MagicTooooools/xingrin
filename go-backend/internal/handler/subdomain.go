package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xingrin/go-backend/internal/dto"
	"github.com/xingrin/go-backend/internal/model"
	"github.com/xingrin/go-backend/internal/pkg/csv"
	"github.com/xingrin/go-backend/internal/service"
)

// SubdomainHandler handles subdomain endpoints
type SubdomainHandler struct {
	svc *service.SubdomainService
}

// NewSubdomainHandler creates a new subdomain handler
func NewSubdomainHandler(svc *service.SubdomainService) *SubdomainHandler {
	return &SubdomainHandler{svc: svc}
}

// List returns paginated subdomains for a target
// GET /api/targets/:id/subdomains
func (h *SubdomainHandler) List(c *gin.Context) {
	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid target ID")
		return
	}

	var query dto.SubdomainListQuery
	if !dto.BindQuery(c, &query) {
		return
	}

	subdomains, total, err := h.svc.ListByTarget(targetID, &query)
	if err != nil {
		if errors.Is(err, service.ErrTargetNotFound) {
			dto.NotFound(c, "Target not found")
			return
		}
		dto.InternalError(c, "Failed to list subdomains")
		return
	}

	// Convert to response
	var resp []dto.SubdomainResponse
	for _, s := range subdomains {
		resp = append(resp, toSubdomainResponse(&s))
	}

	dto.Paginated(c, resp, total, query.GetPage(), query.GetPageSize())
}

// BulkCreate creates multiple subdomains for a target
// POST /api/targets/:id/subdomains/bulk-create
func (h *SubdomainHandler) BulkCreate(c *gin.Context) {
	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid target ID")
		return
	}

	var req dto.BulkCreateSubdomainsRequest
	if !dto.BindJSON(c, &req) {
		return
	}

	createdCount, err := h.svc.BulkCreate(targetID, req.Names)
	if err != nil {
		if errors.Is(err, service.ErrTargetNotFound) {
			dto.NotFound(c, "Target not found")
			return
		}
		if errors.Is(err, service.ErrInvalidTargetType) {
			dto.BadRequest(c, "Target type must be domain for subdomains")
			return
		}
		dto.InternalError(c, "Failed to create subdomains")
		return
	}

	dto.Created(c, dto.BulkCreateSubdomainsResponse{
		CreatedCount: createdCount,
	})
}

// BulkDelete deletes multiple subdomains by IDs
// POST /api/subdomains/bulk-delete
func (h *SubdomainHandler) BulkDelete(c *gin.Context) {
	var req dto.BulkDeleteRequest
	if !dto.BindJSON(c, &req) {
		return
	}

	deletedCount, err := h.svc.BulkDelete(req.IDs)
	if err != nil {
		dto.InternalError(c, "Failed to delete subdomains")
		return
	}

	dto.Success(c, dto.BulkDeleteResponse{DeletedCount: deletedCount})
}

// Export exports subdomains as CSV
// GET /api/targets/:id/subdomains/export
func (h *SubdomainHandler) Export(c *gin.Context) {
	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid target ID")
		return
	}

	// Get count for progress estimation
	count, err := h.svc.CountByTarget(targetID)
	if err != nil {
		if errors.Is(err, service.ErrTargetNotFound) {
			dto.NotFound(c, "Target not found")
			return
		}
		dto.InternalError(c, "Failed to export subdomains")
		return
	}

	rows, err := h.svc.StreamByTarget(targetID)
	if err != nil {
		dto.InternalError(c, "Failed to export subdomains")
		return
	}

	headers := []string{"id", "target_id", "name", "created_at"}
	filename := fmt.Sprintf("target-%d-subdomains.csv", targetID)

	mapper := func(rows *sql.Rows) ([]string, error) {
		subdomain, err := h.svc.ScanRow(rows)
		if err != nil {
			return nil, err
		}

		return []string{
			strconv.Itoa(subdomain.ID),
			strconv.Itoa(subdomain.TargetID),
			subdomain.Name,
			subdomain.CreatedAt.Format("2006-01-02 15:04:05"),
		}, nil
	}

	if err := csv.StreamCSV(c, rows, headers, filename, mapper, count); err != nil {
		return
	}
}

// toSubdomainResponse converts model to response DTO
func toSubdomainResponse(s *model.Subdomain) dto.SubdomainResponse {
	return dto.SubdomainResponse{
		ID:        s.ID,
		TargetID:  s.TargetID,
		Name:      s.Name,
		CreatedAt: s.CreatedAt,
	}
}
