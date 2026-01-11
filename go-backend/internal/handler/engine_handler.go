package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xingrin/go-backend/internal/dto"
	"github.com/xingrin/go-backend/internal/service"
)

// EngineHandler handles engine endpoints
type EngineHandler struct {
	svc *service.EngineService
}

// NewEngineHandler creates a new engine handler
func NewEngineHandler(svc *service.EngineService) *EngineHandler {
	return &EngineHandler{svc: svc}
}

// Create creates a new engine
// POST /api/engines
func (h *EngineHandler) Create(c *gin.Context) {
	var req dto.CreateEngineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.BadRequest(c, "Invalid request body")
		return
	}

	engine, err := h.svc.Create(&req)
	if err != nil {
		if errors.Is(err, service.ErrEngineExists) {
			dto.BadRequest(c, "Engine name already exists")
			return
		}
		dto.InternalError(c, "Failed to create engine")
		return
	}

	dto.Created(c, dto.EngineResponse{
		ID:            engine.ID,
		Name:          engine.Name,
		Configuration: engine.Configuration,
		CreatedAt:     engine.CreatedAt,
		UpdatedAt:     engine.UpdatedAt,
	})
}

// List returns paginated engines
// GET /api/engines
func (h *EngineHandler) List(c *gin.Context) {
	var query dto.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		dto.BadRequest(c, "Invalid query parameters")
		return
	}

	engines, total, err := h.svc.List(&query)
	if err != nil {
		dto.InternalError(c, "Failed to list engines")
		return
	}

	var resp []dto.EngineResponse
	for _, e := range engines {
		resp = append(resp, dto.EngineResponse{
			ID:            e.ID,
			Name:          e.Name,
			Configuration: e.Configuration,
			CreatedAt:     e.CreatedAt,
			UpdatedAt:     e.UpdatedAt,
		})
	}

	dto.Paginated(c, resp, total, query.GetPage(), query.GetPageSize())
}

// GetByID returns an engine by ID
// GET /api/engines/:id
func (h *EngineHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid engine ID")
		return
	}

	engine, err := h.svc.GetByID(id)
	if err != nil {
		if errors.Is(err, service.ErrEngineNotFound) {
			dto.NotFound(c, "Engine not found")
			return
		}
		dto.InternalError(c, "Failed to get engine")
		return
	}

	dto.Success(c, dto.EngineResponse{
		ID:            engine.ID,
		Name:          engine.Name,
		Configuration: engine.Configuration,
		CreatedAt:     engine.CreatedAt,
		UpdatedAt:     engine.UpdatedAt,
	})
}

// Update updates an engine
// PUT /api/engines/:id
func (h *EngineHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid engine ID")
		return
	}

	var req dto.UpdateEngineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.BadRequest(c, "Invalid request body")
		return
	}

	engine, err := h.svc.Update(id, &req)
	if err != nil {
		if errors.Is(err, service.ErrEngineNotFound) {
			dto.NotFound(c, "Engine not found")
			return
		}
		if errors.Is(err, service.ErrEngineExists) {
			dto.BadRequest(c, "Engine name already exists")
			return
		}
		dto.InternalError(c, "Failed to update engine")
		return
	}

	dto.Success(c, dto.EngineResponse{
		ID:            engine.ID,
		Name:          engine.Name,
		Configuration: engine.Configuration,
		CreatedAt:     engine.CreatedAt,
		UpdatedAt:     engine.UpdatedAt,
	})
}

// Delete deletes an engine
// DELETE /api/engines/:id
func (h *EngineHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		dto.BadRequest(c, "Invalid engine ID")
		return
	}

	err = h.svc.Delete(id)
	if err != nil {
		if errors.Is(err, service.ErrEngineNotFound) {
			dto.NotFound(c, "Engine not found")
			return
		}
		dto.InternalError(c, "Failed to delete engine")
		return
	}

	dto.NoContent(c)
}
