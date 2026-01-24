package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yyhuni/orbit/server/internal/handler"
)

func RegisterSubdomainRoutes(protected *gin.RouterGroup, subdomainHandler *handler.SubdomainHandler) {
	protected.GET("/targets/:id/subdomains", subdomainHandler.List)
	protected.GET("/targets/:id/subdomains/export", subdomainHandler.Export)
	protected.POST("/targets/:id/subdomains/bulk-create", subdomainHandler.BulkCreate)

	protected.POST("/subdomains/bulk-delete", subdomainHandler.BulkDelete)
}
