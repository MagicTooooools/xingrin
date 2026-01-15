package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xingrin/server/internal/dto"
	"github.com/xingrin/server/internal/service"
)

// ScanLogHandler handles scan log HTTP requests
type ScanLogHandler struct {
	svc *service.ScanLogService
}

// NewScanLogHandler creates a new scan log handler
func NewScanLogHandler(svc *service.ScanLogService) *ScanLogHandler {
	return &ScanLogHandler{svc: svc}
}

// List returns logs for a scan with cursor pagination
// GET /api/scans/:id/logs?afterId=123&limit=200
func (h *ScanLogHandler) List(c *gin.Context) {
	scanID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid scan ID")
		return
	}

	var query dto.ScanLogListQuery
	if !dto.BindQuery(c, &query) {
		return
	}

	resp, err := h.svc.ListByScanID(scanID, &query)
	if err != nil {
		if errors.Is(err, service.ErrScanNotFound) {
			dto.NotFound(c, "Scan not found")
			return
		}
		dto.InternalError(c, "Failed to get scan logs")
		return
	}

	dto.Success(c, resp)
}

// BulkCreate creates multiple logs for a scan (for worker to write logs)
// POST /api/scans/:id/logs
func (h *ScanLogHandler) BulkCreate(c *gin.Context) {
	scanID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid scan ID")
		return
	}

	var req dto.BulkCreateScanLogsRequest
	if !dto.BindJSON(c, &req) {
		return
	}

	createdCount, err := h.svc.BulkCreate(scanID, req.Logs)
	if err != nil {
		if errors.Is(err, service.ErrScanNotFound) {
			dto.NotFound(c, "Scan not found")
			return
		}
		dto.InternalError(c, "Failed to create scan logs")
		return
	}

	dto.Created(c, dto.BulkCreateScanLogsResponse{
		CreatedCount: createdCount,
	})
}
