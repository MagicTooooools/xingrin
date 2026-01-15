package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xingrin/server/internal/dto"
	"github.com/xingrin/server/internal/service"
)

// ScreenshotHandler handles screenshot endpoints
type ScreenshotHandler struct {
	svc *service.ScreenshotService
}

// NewScreenshotHandler creates a new screenshot handler
func NewScreenshotHandler(svc *service.ScreenshotService) *ScreenshotHandler {
	return &ScreenshotHandler{svc: svc}
}

// ListByTargetID returns screenshots for a target
// GET /api/targets/:id/screenshots
func (h *ScreenshotHandler) ListByTargetID(c *gin.Context) {
	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid target ID")
		return
	}

	var query dto.ScreenshotListQuery
	if !dto.BindQuery(c, &query) {
		return
	}

	screenshots, total, err := h.svc.ListByTargetID(targetID, &query)
	if err != nil {
		if errors.Is(err, service.ErrTargetNotFound) {
			dto.NotFound(c, "Target not found")
			return
		}
		dto.InternalError(c, "Failed to list screenshots")
		return
	}

	// Convert to response (exclude image data)
	var resp []dto.ScreenshotResponse
	for _, s := range screenshots {
		resp = append(resp, dto.ScreenshotResponse{
			ID:         s.ID,
			URL:        s.URL,
			StatusCode: s.StatusCode,
			CreatedAt:  s.CreatedAt,
			UpdatedAt:  s.UpdatedAt,
		})
	}

	dto.Paginated(c, resp, total, query.GetPage(), query.GetPageSize())
}

// GetImage returns screenshot image binary data
// GET /api/screenshots/:id/image
func (h *ScreenshotHandler) GetImage(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid screenshot ID")
		return
	}

	screenshot, err := h.svc.GetByID(id)
	if err != nil {
		if errors.Is(err, service.ErrScreenshotNotFound) {
			dto.NotFound(c, "Screenshot not found")
			return
		}
		dto.InternalError(c, "Failed to get screenshot")
		return
	}

	if len(screenshot.Image) == 0 {
		dto.NotFound(c, "Screenshot image not found")
		return
	}

	// Return WebP image
	c.Header("Content-Type", "image/webp")
	c.Header("Content-Disposition", "inline; filename=\"screenshot_"+strconv.Itoa(id)+".webp\"")
	c.Data(200, "image/webp", screenshot.Image)
}

// BulkDelete deletes multiple screenshots
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

// BulkUpsert creates or updates multiple screenshots for a target
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
