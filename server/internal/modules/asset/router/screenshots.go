package router

import (
	"github.com/gin-gonic/gin"
	assethandler "github.com/yyhuni/lunafox/server/internal/modules/asset/handler"
	snapshothandler "github.com/yyhuni/lunafox/server/internal/modules/snapshot/handler"
)

func RegisterScreenshotRoutes(
	protected *gin.RouterGroup,
	screenshotHandler *assethandler.ScreenshotHandler,
	screenshotSnapshotHandler *snapshothandler.ScreenshotSnapshotHandler,
) {
	protected.GET("/targets/:id/screenshots", screenshotHandler.ListByTargetID)
	protected.POST("/targets/:id/screenshots/bulk-upsert", screenshotSnapshotHandler.BulkUpsert)

	protected.POST("/screenshots/bulk-delete", screenshotHandler.BulkDelete)
}
