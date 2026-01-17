package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/orbit/server/internal/dto"
	"github.com/orbit/server/internal/middleware"
	"github.com/orbit/server/internal/service"
)

// UserHandler handles user endpoints
type UserHandler struct {
	svc *service.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// Create creates a new user
// POST /api/users
func (h *UserHandler) Create(c *gin.Context) {
	var req dto.CreateUserRequest
	if !dto.BindJSON(c, &req) {
		return
	}

	user, err := h.svc.Create(&req)
	if err != nil {
		if errors.Is(err, service.ErrUsernameExists) {
			dto.BadRequest(c, "Username already exists")
			return
		}
		dto.InternalError(c, "Failed to create user")
		return
	}

	dto.Created(c, dto.UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		IsActive:    user.IsActive,
		IsSuperuser: user.IsSuperuser,
		DateJoined:  user.DateJoined,
		LastLogin:   user.LastLogin,
	})
}

// List returns paginated users
// GET /api/users
func (h *UserHandler) List(c *gin.Context) {
	var query dto.PaginationQuery
	if !dto.BindQuery(c, &query) {
		return
	}

	users, total, err := h.svc.List(&query)
	if err != nil {
		dto.InternalError(c, "Failed to list users")
		return
	}

	var resp []dto.UserResponse
	for _, u := range users {
		resp = append(resp, dto.UserResponse{
			ID:          u.ID,
			Username:    u.Username,
			Email:       u.Email,
			IsActive:    u.IsActive,
			IsSuperuser: u.IsSuperuser,
			DateJoined:  u.DateJoined,
			LastLogin:   u.LastLogin,
		})
	}

	dto.Paginated(c, resp, total, query.GetPage(), query.GetPageSize())
}

// UpdatePassword updates current user's password
// PUT /api/users/me/password
func (h *UserHandler) UpdatePassword(c *gin.Context) {
	// Get current user from context
	claims, ok := middleware.GetUserClaims(c)
	if !ok {
		dto.Unauthorized(c, "Not authenticated")
		return
	}

	var req dto.UpdatePasswordRequest
	if !dto.BindJSON(c, &req) {
		return
	}

	err := h.svc.UpdatePassword(claims.UserID, &req)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			dto.NotFound(c, "User not found")
			return
		}
		if errors.Is(err, service.ErrInvalidPassword) {
			dto.BadRequest(c, "Invalid old password")
			return
		}
		dto.InternalError(c, "Failed to update password")
		return
	}

	dto.Success(c, gin.H{"message": "Password updated"})
}
