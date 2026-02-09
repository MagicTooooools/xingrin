package endpoint

import (
	service "github.com/yyhuni/lunafox/server/internal/modules/asset/application"
)

// EndpointHandler handles endpoint endpoints.
type EndpointHandler struct {
	svc *service.EndpointService
}

// NewEndpointHandler creates a new endpoint handler.
func NewEndpointHandler(svc *service.EndpointService) *EndpointHandler {
	return &EndpointHandler{svc: svc}
}
