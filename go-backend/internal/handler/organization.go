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
	if !dto.BindJSON(c, &req) {
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
	var query dto.OrganizationListQuery
	if !dto.BindQuery(c, &query) {
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
			TargetCount: o.TargetCount,
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
		TargetCount: org.TargetCount,
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
	if !dto.BindJSON(c, &req) {
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

// BulkDelete soft deletes multiple organizations
// POST /api/organizations/bulk-delete
func (h *OrganizationHandler) BulkDelete(c *gin.Context) {
	var req dto.BulkDeleteRequest
	if !dto.BindJSON(c, &req) {
		return
	}

	deletedCount, err := h.svc.BulkDelete(req.IDs)
	if err != nil {
		dto.InternalError(c, "Failed to delete organizations")
		return
	}

	dto.Success(c, dto.BulkDeleteResponse{DeletedCount: deletedCount})
}

// ListTargets returns paginated targets for an organization
// GET /api/organizations/:id/targets
func (h *OrganizationHandler) ListTargets(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid organization ID")
		return
	}

	var query dto.TargetListQuery
	if !dto.BindQuery(c, &query) {
		return
	}

	targets, total, err := h.svc.ListTargets(id, &query)
	if err != nil {
		if errors.Is(err, service.ErrOrganizationNotFound) {
			dto.NotFound(c, "Organization not found")
			return
		}
		dto.InternalError(c, "Failed to list targets")
		return
	}

	var resp []dto.TargetResponse
	for _, t := range targets {
		resp = append(resp, dto.TargetResponse{
			ID:        t.ID,
			Name:      t.Name,
			Type:      t.Type,
			CreatedAt: t.CreatedAt,
		})
	}

	dto.Paginated(c, resp, total, query.GetPage(), query.GetPageSize())
}

// LinkTargets adds targets to an organization
// POST /api/organizations/:id/link_targets
func (h *OrganizationHandler) LinkTargets(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid organization ID")
		return
	}

	var req dto.LinkTargetsRequest
	if !dto.BindJSON(c, &req) {
		return
	}

	err = h.svc.LinkTargets(id, req.TargetIDs)
	if err != nil {
		if errors.Is(err, service.ErrOrganizationNotFound) {
			dto.NotFound(c, "Organization not found")
			return
		}
		dto.InternalError(c, "Failed to link targets")
		return
	}

	dto.NoContent(c)
}

// UnlinkTargets removes targets from an organization
// POST /api/organizations/:id/unlink_targets
func (h *OrganizationHandler) UnlinkTargets(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid organization ID")
		return
	}

	var req dto.LinkTargetsRequest
	if !dto.BindJSON(c, &req) {
		return
	}

	unlinkedCount, err := h.svc.UnlinkTargets(id, req.TargetIDs)
	if err != nil {
		if errors.Is(err, service.ErrOrganizationNotFound) {
			dto.NotFound(c, "Organization not found")
			return
		}
		dto.InternalError(c, "Failed to unlink targets")
		return
	}

	dto.Success(c, gin.H{"unlinkedCount": unlinkedCount})
}
