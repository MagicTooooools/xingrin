package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yyhuni/lunafox/server/internal/modules/asset/handler"
	websitehandler "github.com/yyhuni/lunafox/server/internal/modules/asset/handler/website"
)

func registerWebsiteRoutes(protected *gin.RouterGroup, websiteHandler *websitehandler.WebsiteHandler) {
	protected.GET("/targets/:id/websites", websiteHandler.List)
	protected.GET("/targets/:id/websites/export", websiteHandler.Export)
	protected.POST("/targets/:id/websites/bulk-create", websiteHandler.BulkCreate)
	protected.POST("/targets/:id/websites/bulk-upsert", websiteHandler.BulkUpsert)

	protected.DELETE("/websites/:id", websiteHandler.Delete)
	protected.POST("/websites/bulk-delete", websiteHandler.BulkDelete)
}

func registerSubdomainRoutes(protected *gin.RouterGroup, subdomainHandler *handler.SubdomainHandler) {
	protected.GET("/targets/:id/subdomains", subdomainHandler.List)
	protected.GET("/targets/:id/subdomains/export", subdomainHandler.Export)
	protected.POST("/targets/:id/subdomains/bulk-create", subdomainHandler.BulkCreate)

	protected.POST("/subdomains/bulk-delete", subdomainHandler.BulkDelete)
}

func registerDirectoryRoutes(protected *gin.RouterGroup, directoryHandler *handler.DirectoryHandler) {
	protected.GET("/targets/:id/directories", directoryHandler.List)
	protected.GET("/targets/:id/directories/export", directoryHandler.Export)
	protected.POST("/targets/:id/directories/bulk-create", directoryHandler.BulkCreate)
	protected.POST("/targets/:id/directories/bulk-upsert", directoryHandler.BulkUpsert)

	protected.POST("/directories/bulk-delete", directoryHandler.BulkDelete)
}
