package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yyhuni/orbit/server/internal/dto"
	"github.com/yyhuni/orbit/server/internal/service"
)

// WorkerHandler handles worker API endpoints
type WorkerHandler struct {
	svc *service.WorkerService
}

// NewWorkerHandler creates a new worker handler
func NewWorkerHandler(svc *service.WorkerService) *WorkerHandler {
	return &WorkerHandler{svc: svc}
}

// GetTargetName returns target name for a scan
// GET /api/worker/scans/:id/target-name
func (h *WorkerHandler) GetTargetName(c *gin.Context) {
	scanID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid scan ID")
		return
	}

	target, err := h.svc.GetTargetName(scanID)
	if err != nil {
		if errors.Is(err, service.ErrWorkerScanNotFound) {
			dto.NotFound(c, "Scan not found")
			return
		}
		dto.InternalError(c, "Failed to get target name")
		return
	}

	dto.Success(c, dto.WorkerTargetNameResponse{
		Name: target.Name,
		Type: target.Type,
	})
}

// GetProviderConfig returns provider config for a tool
// GET /api/worker/scans/:id/provider-config?tool=subfinder
func (h *WorkerHandler) GetProviderConfig(c *gin.Context) {
	scanID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid scan ID")
		return
	}

	toolName := c.Query("tool")
	config, err := h.svc.GetProviderConfig(scanID, toolName)
	if err != nil {
		if errors.Is(err, service.ErrWorkerScanNotFound) {
			dto.NotFound(c, "Scan not found")
			return
		}
		if errors.Is(err, service.ErrWorkerToolRequired) {
			dto.BadRequest(c, "Tool parameter required")
			return
		}
		dto.InternalError(c, "Failed to get provider config")
		return
	}

	dto.Success(c, config)
}
