package endpoint

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	service "github.com/yyhuni/lunafox/server/internal/modules/asset/application"
	"github.com/yyhuni/lunafox/server/internal/modules/asset/dto"
)

// List returns paginated endpoints for a target.
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

	resp := make([]dto.EndpointResponse, 0, len(endpoints))
	for _, endpoint := range endpoints {
		resp = append(resp, toEndpointResponse(&endpoint))
	}

	dto.Paginated(c, resp, total, query.GetPage(), query.GetPageSize())
}

// GetByID returns an endpoint by ID.
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

func toEndpointResponse(endpoint *service.Endpoint) dto.EndpointResponse {
	tech := []string(endpoint.Tech)
	if tech == nil {
		tech = []string{}
	}
	matchedGFPatterns := []string(endpoint.MatchedGFPatterns)
	if matchedGFPatterns == nil {
		matchedGFPatterns = []string{}
	}

	return dto.EndpointResponse{
		ID:                endpoint.ID,
		TargetID:          endpoint.TargetID,
		URL:               endpoint.URL,
		Host:              endpoint.Host,
		Location:          endpoint.Location,
		Title:             endpoint.Title,
		Webserver:         endpoint.Webserver,
		ContentType:       endpoint.ContentType,
		StatusCode:        endpoint.StatusCode,
		ContentLength:     endpoint.ContentLength,
		ResponseBody:      endpoint.ResponseBody,
		Tech:              tech,
		Vhost:             endpoint.Vhost,
		MatchedGFPatterns: matchedGFPatterns,
		ResponseHeaders:   endpoint.ResponseHeaders,
		CreatedAt:         endpoint.CreatedAt,
	}
}
