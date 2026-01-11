package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xingrin/go-backend/internal/dto"
	"github.com/xingrin/go-backend/internal/service"
)

// OrganizationHandler handles organization endpoints
type OrganizationHandler struct {
	svc *service.OrganizationService
}

// NewOrganizationHandler creates a new organization handler
func NewOrganizationHandler(svc *service.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{svc: svc}
}

// Create creates a new organization
// POST /api/organizations
func (h *OrganizationHandler) Create(c *gin.Context) {
	var req dto.CreateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.BadRequest(c, "Invalid request body")
		return
	}

	org, err := h.svc.Create(&req)
	if err != nil {
		if errors.Is(err, service.ErrOrganizationExists) {
			dto.BadRequest(c, "Organization name already exists")
			return
		}
		dto.InternalError(c, "Failed to create organization")
		return
	}

	dto.Created(c, dto.OrganizationResponse{
		ID:          org.ID,
		Name:        org.Name,
		Description: org.Description,
		CreatedAt:   org.CreatedAt,
	})
}

// List returns paginated organizations
// GET /api/organizations
func (h *OrganizationHandler) List(c *gin.Context) {
	var query dto.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		dto.BadRequest(c, "Invalid query parameters")
		return
	}

	orgs, total, err := h.svc.List(&query)
	if err != nil {
		dto.InternalError(c, "Failed to list organizations")
		return
	}

	var resp []dto.OrganizationResponse
	for _, o := range orgs {
		resp = append(resp, dto.OrganizationResponse{
			ID:          o.ID,
			Name:        o.Name,
			Description: o.Description,
			CreatedAt:   o.CreatedAt,
		})
	}

	dto.Paginated(c, resp, total, query.GetPage(), query.GetPageSize())
}

// GetByID returns an organization by ID
// GET /api/organizations/:id
func (h *OrganizationHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid organization ID")
		return
	}

	org, err := h.svc.GetByID(id)
	if err != nil {
		if errors.Is(err, service.ErrOrganizationNotFound) {
			dto.NotFound(c, "Organization not found")
			return
		}
		dto.InternalError(c, "Failed to get organization")
		return
	}

	dto.Success(c, dto.OrganizationResponse{
		ID:          org.ID,
		Name:        org.Name,
		Description: org.Description,
		CreatedAt:   org.CreatedAt,
	})
}

// Update updates an organization
// PUT /api/organizations/:id
func (h *OrganizationHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid organization ID")
		return
	}

	var req dto.UpdateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.BadRequest(c, "Invalid request body")
		return
	}

	org, err := h.svc.Update(id, &req)
	if err != nil {
		if errors.Is(err, service.ErrOrganizationNotFound) {
			dto.NotFound(c, "Organization not found")
			return
		}
		if errors.Is(err, service.ErrOrganizationExists) {
			dto.BadRequest(c, "Organization name already exists")
			return
		}
		dto.InternalError(c, "Failed to update organization")
		return
	}

	dto.Success(c, dto.OrganizationResponse{
		ID:          org.ID,
		Name:        org.Name,
		Description: org.Description,
		CreatedAt:   org.CreatedAt,
	})
}

// Delete soft deletes an organization
// DELETE /api/organizations/:id
func (h *OrganizationHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid organization ID")
		return
	}

	err = h.svc.Delete(id)
	if err != nil {
		if errors.Is(err, service.ErrOrganizationNotFound) {
			dto.NotFound(c, "Organization not found")
			return
		}
		dto.InternalError(c, "Failed to delete organization")
		return
	}

	dto.NoContent(c)
}
