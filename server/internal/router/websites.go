package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yyhuni/lunafox/server/internal/handler"
)

func RegisterWebsiteRoutes(protected *gin.RouterGroup, websiteHandler *handler.WebsiteHandler) {
	protected.GET("/targets/:id/websites", websiteHandler.List)
	protected.GET("/targets/:id/websites/export", websiteHandler.Export)
	protected.POST("/targets/:id/websites/bulk-create", websiteHandler.BulkCreate)
	protected.POST("/targets/:id/websites/bulk-upsert", websiteHandler.BulkUpsert)

	protected.DELETE("/websites/:id", websiteHandler.Delete)
	protected.POST("/websites/bulk-delete", websiteHandler.BulkDelete)
}
