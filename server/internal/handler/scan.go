package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/orbit/server/internal/dto"
	"github.com/orbit/server/internal/service"
)

// ScanHandler handles scan HTTP requests
type ScanHandler struct {
	service *service.ScanService
}

// NewScanHandler creates a new scan handler
func NewScanHandler(service *service.ScanService) *ScanHandler {
	return &ScanHandler{service: service}
}

// List returns paginated scans
// GET /api/scans
func (h *ScanHandler) List(c *gin.Context) {
	var query dto.ScanListQuery
	if !dto.BindQuery(c, &query) {
		return
	}

	scans, total, err := h.service.List(&query)
	if err != nil {
		dto.InternalError(c, "Failed to list scans")
		return
	}

	// Convert to response DTOs
	items := make([]dto.ScanResponse, len(scans))
	for i, scan := range scans {
		items[i] = *h.service.ToScanResponse(&scan)
	}

	dto.Paginated(c, items, total, query.GetPage(), query.GetPageSize())
}

// GetByID returns a scan by ID
// GET /api/scans/:id
func (h *ScanHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid scan ID")
		return
	}

	scan, err := h.service.GetByID(id)
	if err != nil {
		if errors.Is(err, service.ErrScanNotFound) {
			dto.NotFound(c, "Scan not found")
			return
		}
		dto.InternalError(c, "Failed to get scan")
		return
	}

	dto.Success(c, h.service.ToScanDetailResponse(scan))
}

// Delete soft deletes a scan
// DELETE /api/scans/:id
func (h *ScanHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid scan ID")
		return
	}

	deletedCount, deletedNames, err := h.service.Delete(id)
	if err != nil {
		if errors.Is(err, service.ErrScanNotFound) {
			dto.NotFound(c, "Scan not found")
			return
		}
		dto.InternalError(c, "Failed to delete scan")
		return
	}

	dto.Success(c, gin.H{
		"scanId":       id,
		"deletedCount": deletedCount,
		"deletedScans": deletedNames,
	})
}

// BulkDelete soft deletes multiple scans
// POST /api/scans/bulk-delete
func (h *ScanHandler) BulkDelete(c *gin.Context) {
	var req dto.BulkDeleteRequest
	if !dto.BindJSON(c, &req) {
		return
	}

	deletedCount, deletedNames, err := h.service.BulkDelete(req.IDs)
	if err != nil {
		dto.InternalError(c, "Failed to bulk delete scans")
		return
	}

	dto.Success(c, gin.H{
		"deletedCount": deletedCount,
		"deletedScans": deletedNames,
	})
}

// Statistics returns scan statistics
// GET /api/scans/statistics
func (h *ScanHandler) Statistics(c *gin.Context) {
	stats, err := h.service.GetStatistics()
	if err != nil {
		dto.InternalError(c, "Failed to get scan statistics")
		return
	}

	dto.Success(c, stats)
}

// Stop stops a running scan
// POST /api/scans/:id/stop
func (h *ScanHandler) Stop(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid scan ID")
		return
	}

	revokedCount, err := h.service.Stop(id)
	if err != nil {
		if errors.Is(err, service.ErrScanNotFound) {
			dto.NotFound(c, "Scan not found")
			return
		}
		if errors.Is(err, service.ErrScanCannotStop) {
			dto.BadRequest(c, "Cannot stop scan: scan is not running")
			return
		}
		dto.InternalError(c, "Failed to stop scan")
		return
	}

	dto.Success(c, dto.StopScanResponse{
		RevokedTaskCount: revokedCount,
	})
}

// Create starts a new scan
// POST /api/scans
//
// Request body:
//
//	{
//	  "mode": "normal" | "quick",     // scan mode (default: "normal")
//	  "targetId": 123,                // required for mode=normal
//	  "targets": ["example.com"],     // required for mode=quick (raw targets)
//	  "engineIds": [1, 2],            // engine IDs to run
//	  "config": {}                    // optional scan configuration
//	}
func (h *ScanHandler) Create(c *gin.Context) {
	var req dto.CreateScanRequest
	if !dto.BindJSON(c, &req) {
		return
	}

	// Default mode is "normal"
	if req.Mode == "" {
		req.Mode = "normal"
	}

	switch req.Mode {
	case "normal":
		// Normal scan: requires targetId
		if req.TargetID == 0 {
			dto.BadRequest(c, "targetId is required for normal mode")
			return
		}
		// TODO: Implement when worker integration is ready
		dto.Error(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Normal scan is not yet implemented")

	case "quick":
		// Quick scan: requires targets list
		if len(req.Targets) == 0 {
			dto.BadRequest(c, "targets is required for quick mode")
			return
		}
		// TODO: Implement when worker integration is ready
		dto.Error(c, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Quick scan is not yet implemented")

	default:
		dto.BadRequest(c, "Invalid mode, must be 'normal' or 'quick'")
	}
}
