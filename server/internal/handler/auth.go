package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/orbit/server/internal/auth"
	"github.com/orbit/server/internal/dto"
	"github.com/orbit/server/internal/middleware"
	"github.com/orbit/server/internal/model"
	"gorm.io/gorm"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	db         *gorm.DB
	jwtManager *auth.JWTManager
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(db *gorm.DB, jwtManager *auth.JWTManager) *AuthHandler {
	return &AuthHandler{
		db:         db,
		jwtManager: jwtManager,
	}
}

// LoginRequest represents login request body
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents login response
type LoginResponse struct {
	AccessToken  string   `json:"accessToken"`
	RefreshToken string   `json:"refreshToken"`
	ExpiresIn    int64    `json:"expiresIn"`
	User         UserInfo `json:"user"`
}

// UserInfo represents basic user information
type UserInfo struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// RefreshRequest represents refresh token request
type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// RefreshResponse represents refresh token response
type RefreshResponse struct {
	AccessToken string `json:"accessToken"`
	ExpiresIn   int64  `json:"expiresIn"`
}

// Login handles user login
// POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if !dto.BindJSON(c, &req) {
		return
	}

	// Find user by username
	var user model.User
	if err := h.db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			dto.Unauthorized(c, "Invalid username or password")
			return
		}
		dto.InternalError(c, "Database error")
		return
	}

	// Check if user is active
	if !user.IsActive {
		dto.Unauthorized(c, "User account is disabled")
		return
	}

	// Verify password
	if !auth.VerifyPassword(req.Password, user.Password) {
		dto.Unauthorized(c, "Invalid username or password")
		return
	}

	// Generate token pair
	tokenPair, err := h.jwtManager.GenerateTokenPair(user.ID, user.Username)
	if err != nil {
		dto.InternalError(c, "Failed to generate token")
		return
	}

	dto.Success(c, LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
		User: UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
	})
}

// RefreshToken handles token refresh
// POST /api/auth/refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if !dto.BindJSON(c, &req) {
		return
	}

	// Validate refresh token
	claims, err := h.jwtManager.ValidateToken(req.RefreshToken)
	if err != nil {
		dto.Unauthorized(c, "Invalid or expired refresh token")
		return
	}

	// Generate new access token
	accessToken, expiresIn, err := h.jwtManager.GenerateAccessToken(claims.UserID, claims.Username)
	if err != nil {
		dto.InternalError(c, "Failed to generate token")
		return
	}

	dto.Success(c, RefreshResponse{
		AccessToken: accessToken,
		ExpiresIn:   expiresIn,
	})
}

// GetCurrentUser returns current user info
// GET /api/auth/me
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	claims, ok := middleware.GetUserClaims(c)
	if !ok {
		dto.Unauthorized(c, "Not authenticated")
		return
	}

	// Get user from database
	var user model.User
	if err := h.db.First(&user, claims.UserID).Error; err != nil {
		dto.NotFound(c, "User not found")
		return
	}

	dto.Success(c, UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	})
}
