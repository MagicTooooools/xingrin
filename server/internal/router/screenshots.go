package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yyhuni/orbit/server/internal/handler"
)

func RegisterScreenshotRoutes(
	protected *gin.RouterGroup,
	screenshotHandler *handler.ScreenshotHandler,
	screenshotSnapshotHandler *handler.ScreenshotSnapshotHandler,
) {
	protected.GET("/targets/:id/screenshots", screenshotHandler.ListByTargetID)
	protected.POST("/targets/:id/screenshots/bulk-upsert", screenshotSnapshotHandler.BulkUpsert)

	protected.POST("/screenshots/bulk-delete", screenshotHandler.BulkDelete)
}
