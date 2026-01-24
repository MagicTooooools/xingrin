package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yyhuni/orbit/server/internal/handler"
)

func RegisterHostPortRoutes(
	protected *gin.RouterGroup,
	hostPortHandler *handler.HostPortHandler,
	hostPortSnapshotHandler *handler.HostPortSnapshotHandler,
) {
	protected.GET("/targets/:id/host-ports", hostPortHandler.List)
	protected.GET("/targets/:id/host-ports/export", hostPortHandler.Export)
	protected.POST("/targets/:id/host-ports/bulk-upsert", hostPortSnapshotHandler.BulkUpsert)

	protected.POST("/host-ports/bulk-delete", hostPortHandler.BulkDelete)
}
