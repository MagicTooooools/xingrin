package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	service "github.com/yyhuni/lunafox/server/internal/modules/identity/application"
	"github.com/yyhuni/lunafox/server/internal/modules/identity/dto"
)

// ListOrganizationTargets returns paginated targets for an organization.
// GET /api/organizations/:id/targets
func (handler *OrganizationHandler) ListOrganizationTargets(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid organization ID")
		return
	}

	var query dto.TargetListQuery
	if !dto.BindQuery(c, &query) {
		return
	}

	targets, total, err := handler.svc.ListOrganizationTargets(id, &query)
	if err != nil {
		if errors.Is(err, service.ErrOrganizationNotFound) {
			dto.NotFound(c, "Organization not found")
			return
		}
		dto.InternalError(c, "Failed to list targets")
		return
	}

	resp := make([]dto.TargetResponse, 0, len(targets))
	for _, target := range targets {
		resp = append(resp, toTargetOutput(target))
	}

	dto.Paginated(c, resp, total, query.GetPage(), query.GetPageSize())
}

// LinkOrganizationTargets adds targets to an organization.
// POST /api/organizations/:id/link_targets
func (handler *OrganizationHandler) LinkOrganizationTargets(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid organization ID")
		return
	}

	var req dto.LinkTargetsRequest
	if !dto.BindJSON(c, &req) {
		return
	}

	err = handler.svc.LinkOrganizationTargets(id, req.TargetIDs)
	if err != nil {
		if errors.Is(err, service.ErrOrganizationNotFound) {
			dto.NotFound(c, "Organization not found")
			return
		}
		if errors.Is(err, service.ErrTargetNotFound) {
			dto.BadRequest(c, "One or more target IDs do not exist")
			return
		}
		dto.InternalError(c, "Failed to link targets")
		return
	}

	dto.NoContent(c)
}

// UnlinkOrganizationTargets removes targets from an organization.
// POST /api/organizations/:id/unlink_targets
func (handler *OrganizationHandler) UnlinkOrganizationTargets(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid organization ID")
		return
	}

	var req dto.LinkTargetsRequest
	if !dto.BindJSON(c, &req) {
		return
	}

	unlinkedCount, err := handler.svc.UnlinkOrganizationTargets(id, req.TargetIDs)
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
