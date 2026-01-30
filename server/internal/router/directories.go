package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yyhuni/lunafox/server/internal/handler"
)

func RegisterDirectoryRoutes(protected *gin.RouterGroup, directoryHandler *handler.DirectoryHandler) {
	protected.GET("/targets/:id/directories", directoryHandler.List)
	protected.GET("/targets/:id/directories/export", directoryHandler.Export)
	protected.POST("/targets/:id/directories/bulk-create", directoryHandler.BulkCreate)
	protected.POST("/targets/:id/directories/bulk-upsert", directoryHandler.BulkUpsert)

	protected.POST("/directories/bulk-delete", directoryHandler.BulkDelete)
}
