package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yyhuni/orbit/server/internal/dto"
	"github.com/yyhuni/orbit/server/internal/service"
)

// AgentHandler handles agent API endpoints
type AgentHandler struct {
	svc *service.AgentService
}

// NewAgentHandler creates a new agent handler
func NewAgentHandler(svc *service.AgentService) *AgentHandler {
	return &AgentHandler{svc: svc}
}

// UpdateStatus updates scan status (called by Agent based on Worker exit code)
// PATCH /api/agent/scans/:id/status
func (h *AgentHandler) UpdateStatus(c *gin.Context) {
	scanID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid scan ID")
		return
	}

	var req dto.AgentUpdateStatusRequest
	if !dto.BindJSON(c, &req) {
		return
	}

	if err := h.svc.UpdateStatus(scanID, req.Status, req.ErrorMessage); err != nil {
		if errors.Is(err, service.ErrAgentScanNotFound) {
			dto.NotFound(c, "Scan not found")
			return
		}
		if errors.Is(err, service.ErrAgentInvalidTransition) {
			dto.BadRequest(c, "Invalid status transition")
			return
		}
		dto.InternalError(c, "Failed to update status")
		return
	}

	dto.Success(c, dto.AgentUpdateStatusResponse{Success: true})
}
