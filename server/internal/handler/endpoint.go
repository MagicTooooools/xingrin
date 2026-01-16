package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/orbit/server/internal/dto"
	"github.com/orbit/server/internal/model"
	"github.com/orbit/server/internal/pkg/csv"
	"github.com/orbit/server/internal/service"
)

// EndpointHandler handles endpoint endpoints
type EndpointHandler struct {
	svc *service.EndpointService
}

// NewEndpointHandler creates a new endpoint handler
func NewEndpointHandler(svc *service.EndpointService) *EndpointHandler {
	return &EndpointHandler{svc: svc}
}

// List returns paginated endpoints for a target
// GET /api/targets/:id/endpoints
func (h *EndpointHandler) List(c *gin.Context) {
	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid target ID")
		return
	}

	var query dto.EndpointListQuery
	if !dto.BindQuery(c, &query) {
		return
	}

	endpoints, total, err := h.svc.ListByTarget(targetID, &query)
	if err != nil {
		if errors.Is(err, service.ErrTargetNotFound) {
			dto.NotFound(c, "Target not found")
			return
		}
		dto.InternalError(c, "Failed to list endpoints")
		return
	}

	// Convert to response
	var resp []dto.EndpointResponse
	for _, e := range endpoints {
		resp = append(resp, toEndpointResponse(&e))
	}

	dto.Paginated(c, resp, total, query.GetPage(), query.GetPageSize())
}

// GetByID returns an endpoint by ID
// GET /api/endpoints/:id
func (h *EndpointHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid endpoint ID")
		return
	}

	endpoint, err := h.svc.GetByID(id)
	if err != nil {
		if errors.Is(err, service.ErrEndpointNotFound) {
			dto.NotFound(c, "Endpoint not found")
			return
		}
		dto.InternalError(c, "Failed to get endpoint")
		return
	}

	dto.Success(c, toEndpointResponse(endpoint))
}

// BulkCreate creates multiple endpoints for a target
// POST /api/targets/:id/endpoints/bulk-create
func (h *EndpointHandler) BulkCreate(c *gin.Context) {
	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid target ID")
		return
	}

	var req dto.BulkCreateEndpointsRequest
	if !dto.BindJSON(c, &req) {
		return
	}

	createdCount, err := h.svc.BulkCreate(targetID, req.URLs)
	if err != nil {
		if errors.Is(err, service.ErrTargetNotFound) {
			dto.NotFound(c, "Target not found")
			return
		}
		dto.InternalError(c, "Failed to create endpoints")
		return
	}

	dto.Created(c, dto.BulkCreateEndpointsResponse{
		CreatedCount: createdCount,
	})
}

// Delete deletes an endpoint by ID
// DELETE /api/endpoints/:id
func (h *EndpointHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid endpoint ID")
		return
	}

	err = h.svc.Delete(id)
	if err != nil {
		if errors.Is(err, service.ErrEndpointNotFound) {
			dto.NotFound(c, "Endpoint not found")
			return
		}
		dto.InternalError(c, "Failed to delete endpoint")
		return
	}

	dto.NoContent(c)
}

// BulkDelete deletes multiple endpoints by IDs
// POST /api/endpoints/bulk-delete
func (h *EndpointHandler) BulkDelete(c *gin.Context) {
	var req dto.BulkDeleteRequest
	if !dto.BindJSON(c, &req) {
		return
	}

	deletedCount, err := h.svc.BulkDelete(req.IDs)
	if err != nil {
		dto.InternalError(c, "Failed to delete endpoints")
		return
	}

	dto.Success(c, dto.BulkDeleteResponse{DeletedCount: deletedCount})
}

// Export exports endpoints as CSV
// GET /api/targets/:id/endpoints/export
func (h *EndpointHandler) Export(c *gin.Context) {
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
		dto.InternalError(c, "Failed to export endpoints")
		return
	}

	rows, err := h.svc.StreamByTarget(targetID)
	if err != nil {
		dto.InternalError(c, "Failed to export endpoints")
		return
	}

	headers := []string{
		"id", "target_id", "url", "host", "location", "title", "status_code",
		"content_length", "content_type", "webserver", "tech", "matched_gf_patterns",
		"response_body", "response_headers", "vhost", "created_at",
	}

	filename := fmt.Sprintf("target-%d-endpoints.csv", targetID)

	mapper := func(rows *sql.Rows) ([]string, error) {
		endpoint, err := h.svc.ScanRow(rows)
		if err != nil {
			return nil, err
		}

		statusCode := ""
		if endpoint.StatusCode != nil {
			statusCode = strconv.Itoa(*endpoint.StatusCode)
		}

		contentLength := ""
		if endpoint.ContentLength != nil {
			contentLength = strconv.Itoa(*endpoint.ContentLength)
		}

		vhost := ""
		if endpoint.Vhost != nil {
			vhost = strconv.FormatBool(*endpoint.Vhost)
		}

		tech := ""
		if len(endpoint.Tech) > 0 {
			tech = strings.Join(endpoint.Tech, "|")
		}

		matchedGFPatterns := ""
		if len(endpoint.MatchedGFPatterns) > 0 {
			matchedGFPatterns = strings.Join(endpoint.MatchedGFPatterns, "|")
		}

		return []string{
			strconv.Itoa(endpoint.ID),
			strconv.Itoa(endpoint.TargetID),
			endpoint.URL,
			endpoint.Host,
			endpoint.Location,
			endpoint.Title,
			statusCode,
			contentLength,
			endpoint.ContentType,
			endpoint.Webserver,
			tech,
			matchedGFPatterns,
			endpoint.ResponseBody,
			endpoint.ResponseHeaders,
			vhost,
			endpoint.CreatedAt.Format("2006-01-02 15:04:05"),
		}, nil
	}

	if err := csv.StreamCSV(c, rows, headers, filename, mapper, count); err != nil {
		return
	}
}

// toEndpointResponse converts model to response DTO
func toEndpointResponse(e *model.Endpoint) dto.EndpointResponse {
	tech := []string(e.Tech)
	if tech == nil {
		tech = []string{}
	}
	matchedGFPatterns := []string(e.MatchedGFPatterns)
	if matchedGFPatterns == nil {
		matchedGFPatterns = []string{}
	}
	return dto.EndpointResponse{
		ID:                e.ID,
		TargetID:          e.TargetID,
		URL:               e.URL,
		Host:              e.Host,
		Location:          e.Location,
		Title:             e.Title,
		Webserver:         e.Webserver,
		ContentType:       e.ContentType,
		StatusCode:        e.StatusCode,
		ContentLength:     e.ContentLength,
		ResponseBody:      e.ResponseBody,
		Tech:              tech,
		Vhost:             e.Vhost,
		MatchedGFPatterns: matchedGFPatterns,
		ResponseHeaders:   e.ResponseHeaders,
		CreatedAt:         e.CreatedAt,
	}
}

// BulkUpsert creates or updates multiple endpoints for a target
// POST /api/targets/:id/endpoints/bulk-upsert
func (h *EndpointHandler) BulkUpsert(c *gin.Context) {
	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid target ID")
		return
	}

	var req dto.BulkUpsertEndpointsRequest
	if !dto.BindJSON(c, &req) {
		return
	}

	affectedCount, err := h.svc.BulkUpsert(targetID, req.Endpoints)
	if err != nil {
		if errors.Is(err, service.ErrTargetNotFound) {
			dto.NotFound(c, "Target not found")
			return
		}
		dto.InternalError(c, "Failed to upsert endpoints")
		return
	}

	dto.Success(c, dto.BulkUpsertEndpointsResponse{
		AffectedCount: affectedCount,
	})
}
