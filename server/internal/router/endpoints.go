package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yyhuni/lunafox/server/internal/handler"
)

func RegisterEndpointRoutes(
	protected *gin.RouterGroup,
	endpointHandler *handler.EndpointHandler,
	endpointSnapshotHandler *handler.EndpointSnapshotHandler,
) {
	protected.GET("/targets/:id/endpoints", endpointHandler.List)
	protected.GET("/targets/:id/endpoints/export", endpointHandler.Export)
	protected.POST("/targets/:id/endpoints/bulk-create", endpointHandler.BulkCreate)
	protected.POST("/targets/:id/endpoints/bulk-upsert", endpointSnapshotHandler.BulkUpsert)

	protected.GET("/endpoints/:id", endpointHandler.GetByID)
	protected.DELETE("/endpoints/:id", endpointHandler.Delete)
	protected.POST("/endpoints/bulk-delete", endpointHandler.BulkDelete)
}
