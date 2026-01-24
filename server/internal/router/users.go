package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yyhuni/orbit/server/internal/handler"
)

func RegisterUserRoutes(protected *gin.RouterGroup, userHandler *handler.UserHandler) {
	protected.POST("/users", userHandler.Create)
	protected.GET("/users", userHandler.List)
	protected.PUT("/users/me/password", userHandler.UpdatePassword)
}
