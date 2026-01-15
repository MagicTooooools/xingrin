package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xingrin/server/internal/dto"
	"github.com/xingrin/server/internal/model"
	"github.com/xingrin/server/internal/pkg/csv"
	"github.com/xingrin/server/internal/service"
)

// WebsiteHandler handles website endpoints
type WebsiteHandler struct {
	svc *service.WebsiteService
}

// NewWebsiteHandler creates a new website handler
func NewWebsiteHandler(svc *service.WebsiteService) *WebsiteHandler {
	return &WebsiteHandler{svc: svc}
}

// List returns paginated websites for a target
// GET /api/targets/:id/websites
func (h *WebsiteHandler) List(c *gin.Context) {
	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid target ID")
		return
	}

	var query dto.WebsiteListQuery
	if !dto.BindQuery(c, &query) {
		return
	}

	websites, total, err := h.svc.ListByTarget(targetID, &query)
	if err != nil {
		if errors.Is(err, service.ErrTargetNotFound) {
			dto.NotFound(c, "Target not found")
			return
		}
		dto.InternalError(c, "Failed to list websites")
		return
	}

	// Convert to response
	var resp []dto.WebsiteResponse
	for _, w := range websites {
		resp = append(resp, toWebsiteResponse(&w))
	}

	dto.Paginated(c, resp, total, query.GetPage(), query.GetPageSize())
}

// BulkCreate creates multiple websites for a target
// POST /api/targets/:id/websites/bulk-create
func (h *WebsiteHandler) BulkCreate(c *gin.Context) {
	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid target ID")
		return
	}

	var req dto.BulkCreateWebsitesRequest
	if !dto.BindJSON(c, &req) {
		return
	}

	createdCount, err := h.svc.BulkCreate(targetID, req.URLs)
	if err != nil {
		if errors.Is(err, service.ErrTargetNotFound) {
			dto.NotFound(c, "Target not found")
			return
		}
		dto.InternalError(c, "Failed to create websites")
		return
	}

	dto.Created(c, dto.BulkCreateWebsitesResponse{
		CreatedCount: createdCount,
	})
}

// Delete deletes a website by ID
// DELETE /api/websites/:id
func (h *WebsiteHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid website ID")
		return
	}

	err = h.svc.Delete(id)
	if err != nil {
		if errors.Is(err, service.ErrWebsiteNotFound) {
			dto.NotFound(c, "Website not found")
			return
		}
		dto.InternalError(c, "Failed to delete website")
		return
	}

	dto.NoContent(c)
}

// BulkDelete deletes multiple websites by IDs
// POST /api/websites/bulk-delete
func (h *WebsiteHandler) BulkDelete(c *gin.Context) {
	var req dto.BulkDeleteRequest
	if !dto.BindJSON(c, &req) {
		return
	}

	deletedCount, err := h.svc.BulkDelete(req.IDs)
	if err != nil {
		dto.InternalError(c, "Failed to delete websites")
		return
	}

	dto.Success(c, dto.BulkDeleteResponse{DeletedCount: deletedCount})
}

// BulkUpsert creates or updates multiple websites for a target (for scanner import)
// POST /api/targets/:id/websites/bulk-upsert
func (h *WebsiteHandler) BulkUpsert(c *gin.Context) {
	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid target ID")
		return
	}

	var req dto.BulkUpsertWebsitesRequest
	if !dto.BindJSON(c, &req) {
		return
	}

	upsertedCount, err := h.svc.BulkUpsert(targetID, req.Websites)
	if err != nil {
		if errors.Is(err, service.ErrTargetNotFound) {
			dto.NotFound(c, "Target not found")
			return
		}
		dto.InternalError(c, "Failed to upsert websites")
		return
	}

	dto.Success(c, dto.BulkUpsertWebsitesResponse{
		UpsertedCount: int(upsertedCount),
	})
}

// toWebsiteResponse converts model to response DTO
func toWebsiteResponse(w *model.Website) dto.WebsiteResponse {
	tech := w.Tech
	if tech == nil {
		tech = []string{}
	}
	return dto.WebsiteResponse{
		ID:              w.ID,
		URL:             w.URL,
		Host:            w.Host,
		Location:        w.Location,
		Title:           w.Title,
		Webserver:       w.Webserver,
		ContentType:     w.ContentType,
		StatusCode:      w.StatusCode,
		ContentLength:   w.ContentLength,
		ResponseBody:    w.ResponseBody,
		Tech:            tech,
		Vhost:           w.Vhost,
		ResponseHeaders: w.ResponseHeaders,
		CreatedAt:       w.CreatedAt,
	}
}

// Export exports websites as CSV
// GET /api/targets/:id/websites/export
func (h *WebsiteHandler) Export(c *gin.Context) {
	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid target ID")
		return
	}

	// Get count for Content-Length estimation (enables browser progress bar)
	count, err := h.svc.CountByTarget(targetID)
	if err != nil {
		if errors.Is(err, service.ErrTargetNotFound) {
			dto.NotFound(c, "Target not found")
			return
		}
		dto.InternalError(c, "Failed to export websites")
		return
	}

	rows, err := h.svc.StreamByTarget(targetID)
	if err != nil {
		dto.InternalError(c, "Failed to export websites")
		return
	}

	headers := []string{
		"id", "target_id", "url", "host", "location", "title", "status_code",
		"content_length", "content_type", "webserver", "tech",
		"response_body", "response_headers", "vhost", "created_at",
	}

	filename := fmt.Sprintf("target-%d-websites.csv", targetID)

	mapper := func(rows *sql.Rows) ([]string, error) {
		website, err := h.svc.ScanRow(rows)
		if err != nil {
			return nil, err
		}

		statusCode := ""
		if website.StatusCode != nil {
			statusCode = strconv.Itoa(*website.StatusCode)
		}

		contentLength := ""
		if website.ContentLength != nil {
			contentLength = strconv.Itoa(*website.ContentLength)
		}

		vhost := ""
		if website.Vhost != nil {
			vhost = strconv.FormatBool(*website.Vhost)
		}

		tech := ""
		if len(website.Tech) > 0 {
			tech = strings.Join(website.Tech, "|")
		}

		return []string{
			strconv.Itoa(website.ID),
			strconv.Itoa(website.TargetID),
			website.URL,
			website.Host,
			website.Location,
			website.Title,
			statusCode,
			contentLength,
			website.ContentType,
			website.Webserver,
			tech,
			website.ResponseBody,
			website.ResponseHeaders,
			vhost,
			website.CreatedAt.Format("2006-01-02 15:04:05"),
		}, nil
	}

	if err := csv.StreamCSV(c, rows, headers, filename, mapper, count); err != nil {
		// Response already started, can't send error
		return
	}
}
