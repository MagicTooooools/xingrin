package router

import (
	"github.com/gin-gonic/gin"
	assethandler "github.com/yyhuni/lunafox/server/internal/modules/asset/handler"
	snapshothandler "github.com/yyhuni/lunafox/server/internal/modules/snapshot/handler"
)

func registerHostPortRoutes(
	protected *gin.RouterGroup,
	hostPortHandler *assethandler.HostPortHandler,
	hostPortSnapshotHandler *snapshothandler.HostPortSnapshotHandler,
) {
	protected.GET("/targets/:id/host-ports", hostPortHandler.List)
	protected.GET("/targets/:id/host-ports/export", hostPortHandler.Export)
	protected.POST("/targets/:id/host-ports/bulk-upsert", hostPortSnapshotHandler.BulkUpsert)

	protected.POST("/host-ports/bulk-delete", hostPortHandler.BulkDelete)
}
