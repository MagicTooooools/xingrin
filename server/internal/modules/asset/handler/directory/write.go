package directory

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	service "github.com/yyhuni/lunafox/server/internal/modules/asset/application"
	"github.com/yyhuni/lunafox/server/internal/modules/asset/dto"
)

// BulkCreate creates multiple directories for a target.
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

	dto.Created(c, dto.BulkCreateDirectoriesResponse{CreatedCount: createdCount})
}

// BulkDelete deletes multiple directories by IDs.
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

// BulkUpsert creates or updates multiple directories for a target.
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

	dto.Success(c, dto.BulkUpsertDirectoriesResponse{AffectedCount: affectedCount})
}
