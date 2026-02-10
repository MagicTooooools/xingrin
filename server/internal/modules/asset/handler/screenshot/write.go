package screenshot

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	service "github.com/yyhuni/lunafox/server/internal/modules/asset/application"
	"github.com/yyhuni/lunafox/server/internal/modules/asset/dto"
)

// BulkDelete deletes multiple screenshots.
// POST /api/screenshots/bulk-delete
func (h *ScreenshotHandler) BulkDelete(c *gin.Context) {
	var req dto.BulkDeleteRequest
	if !dto.BindJSON(c, &req) {
		return
	}

	deletedCount, err := h.svc.BulkDelete(req.IDs)
	if err != nil {
		dto.InternalError(c, "Failed to delete screenshots")
		return
	}

	dto.Success(c, dto.BulkDeleteResponse{DeletedCount: deletedCount})
}

// BulkUpsert creates or updates multiple screenshots for a target.
// POST /api/targets/:id/screenshots/bulk-upsert
func (h *ScreenshotHandler) BulkUpsert(c *gin.Context) {
	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid target ID")
		return
	}

	var req dto.BulkUpsertScreenshotRequest
	if !dto.BindJSON(c, &req) {
		return
	}

	upsertedCount, err := h.svc.BulkUpsert(targetID, &req)
	if err != nil {
		if errors.Is(err, service.ErrTargetNotFound) {
			dto.NotFound(c, "Target not found")
			return
		}
		dto.InternalError(c, "Failed to upsert screenshots")
		return
	}

	dto.Success(c, dto.BulkUpsertScreenshotResponse{UpsertedCount: upsertedCount})
}
