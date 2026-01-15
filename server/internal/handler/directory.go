package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xingrin/server/internal/dto"
	"github.com/xingrin/server/internal/model"
	"github.com/xingrin/server/internal/pkg/csv"
	"github.com/xingrin/server/internal/service"
)

// DirectoryHandler handles directory endpoints
type DirectoryHandler struct {
	svc *service.DirectoryService
}

// NewDirectoryHandler creates a new directory handler
func NewDirectoryHandler(svc *service.DirectoryService) *DirectoryHandler {
	return &DirectoryHandler{svc: svc}
}

// List returns paginated directories for a target
// GET /api/targets/:id/directories
func (h *DirectoryHandler) List(c *gin.Context) {
	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid target ID")
		return
	}

	var query dto.DirectoryListQuery
	if !dto.BindQuery(c, &query) {
		return
	}

	directories, total, err := h.svc.ListByTarget(targetID, &query)
	if err != nil {
		if errors.Is(err, service.ErrTargetNotFound) {
			dto.NotFound(c, "Target not found")
			return
		}
		dto.InternalError(c, "Failed to list directories")
		return
	}

	// Convert to response
	var resp []dto.DirectoryResponse
	for _, d := range directories {
		resp = append(resp, toDirectoryResponse(&d))
	}

	dto.Paginated(c, resp, total, query.GetPage(), query.GetPageSize())
}

// BulkCreate creates multiple directories for a target
// POST /api/targets/:id/directories/bulk-create
func (h *DirectoryHandler) BulkCreate(c *gin.Context) {
	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid target ID")
		return
	}

	var req dto.BulkCreateDirectoriesRequest
	if !dto.BindJSON(c, &req) {
		return
	}

	createdCount, err := h.svc.BulkCreate(targetID, req.URLs)
	if err != nil {
		if errors.Is(err, service.ErrTargetNotFound) {
			dto.NotFound(c, "Target not found")
			return
		}
		dto.InternalError(c, "Failed to create directories")
		return
	}

	dto.Created(c, dto.BulkCreateDirectoriesResponse{
		CreatedCount: createdCount,
	})
}

// BulkDelete deletes multiple directories by IDs
// POST /api/directories/bulk-delete
func (h *DirectoryHandler) BulkDelete(c *gin.Context) {
	var req dto.BulkDeleteRequest
	if !dto.BindJSON(c, &req) {
		return
	}

	deletedCount, err := h.svc.BulkDelete(req.IDs)
	if err != nil {
		dto.InternalError(c, "Failed to delete directories")
		return
	}

	dto.Success(c, dto.BulkDeleteResponse{DeletedCount: deletedCount})
}

// Export exports directories as CSV
// GET /api/targets/:id/directories/export
func (h *DirectoryHandler) Export(c *gin.Context) {
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
		dto.InternalError(c, "Failed to export directories")
		return
	}

	rows, err := h.svc.StreamByTarget(targetID)
	if err != nil {
		dto.InternalError(c, "Failed to export directories")
		return
	}

	headers := []string{
		"id", "target_id", "url", "status", "content_length",
		"content_type", "duration", "created_at",
	}

	filename := fmt.Sprintf("target-%d-directories.csv", targetID)

	mapper := func(rows *sql.Rows) ([]string, error) {
		directory, err := h.svc.ScanRow(rows)
		if err != nil {
			return nil, err
		}

		status := ""
		if directory.Status != nil {
			status = strconv.Itoa(*directory.Status)
		}

		contentLength := ""
		if directory.ContentLength != nil {
			contentLength = strconv.Itoa(*directory.ContentLength)
		}

		duration := ""
		if directory.Duration != nil {
			duration = strconv.Itoa(*directory.Duration)
		}

		return []string{
			strconv.Itoa(directory.ID),
			strconv.Itoa(directory.TargetID),
			directory.URL,
			status,
			contentLength,
			directory.ContentType,
			duration,
			directory.CreatedAt.Format("2006-01-02 15:04:05"),
		}, nil
	}

	if err := csv.StreamCSV(c, rows, headers, filename, mapper, count); err != nil {
		return
	}
}

// toDirectoryResponse converts model to response DTO
func toDirectoryResponse(d *model.Directory) dto.DirectoryResponse {
	return dto.DirectoryResponse{
		ID:            d.ID,
		TargetID:      d.TargetID,
		URL:           d.URL,
		Status:        d.Status,
		ContentLength: d.ContentLength,
		ContentType:   d.ContentType,
		Duration:      d.Duration,
		CreatedAt:     d.CreatedAt,
	}
}

// BulkUpsert creates or updates multiple directories for a target
// POST /api/targets/:id/directories/bulk-upsert
func (h *DirectoryHandler) BulkUpsert(c *gin.Context) {
	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid target ID")
		return
	}

	var req dto.BulkUpsertDirectoriesRequest
	if !dto.BindJSON(c, &req) {
		return
	}

	affectedCount, err := h.svc.BulkUpsert(targetID, req.Directories)
	if err != nil {
		if errors.Is(err, service.ErrTargetNotFound) {
			dto.NotFound(c, "Target not found")
			return
		}
		dto.InternalError(c, "Failed to upsert directories")
		return
	}

	dto.Success(c, dto.BulkUpsertDirectoriesResponse{
		AffectedCount: affectedCount,
	})
}
