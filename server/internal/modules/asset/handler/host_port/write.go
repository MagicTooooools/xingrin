package hostport

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	service "github.com/yyhuni/lunafox/server/internal/modules/asset/application"
	"github.com/yyhuni/lunafox/server/internal/modules/asset/dto"
)

// BulkUpsert creates multiple host-port mappings and ignores duplicates.
// POST /api/targets/:id/host-ports/bulk-upsert
func (h *HostPortHandler) BulkUpsert(c *gin.Context) {
	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid target ID")
		return
	}

	var req dto.BulkUpsertHostPortsRequest
	if !dto.BindJSON(c, &req) {
		return
	}

	upsertedCount, err := h.svc.BulkUpsert(targetID, req.Mappings)
	if err != nil {
		if errors.Is(err, service.ErrTargetNotFound) {
			dto.NotFound(c, "Target not found")
			return
		}
		dto.InternalError(c, "Failed to upsert host-ports")
		return
	}

	dto.Success(c, dto.BulkUpsertHostPortsResponse{UpsertedCount: int(upsertedCount)})
}

// BulkDelete deletes host-port mappings by IP list.
// POST /api/host-ports/bulk-delete
func (h *HostPortHandler) BulkDelete(c *gin.Context) {
	var req dto.BulkDeleteHostPortsRequest
	if !dto.BindJSON(c, &req) {
		return
	}

	deletedCount, err := h.svc.BulkDeleteByIPs(req.IPs)
	if err != nil {
		dto.InternalError(c, "Failed to delete host-ports")
		return
	}

	dto.Success(c, dto.BulkDeleteHostPortsResponse{DeletedCount: deletedCount})
}
