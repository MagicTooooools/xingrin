package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yyhuni/orbit/server/internal/handler"
)

func RegisterAuthRoutes(api *gin.RouterGroup, authHandler *handler.AuthHandler) {
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/refresh", authHandler.RefreshToken)
	}
}

func RegisterAuthProtectedRoutes(protected *gin.RouterGroup, authHandler *handler.AuthHandler) {
	protected.GET("/auth/me", authHandler.GetCurrentUser)
}
