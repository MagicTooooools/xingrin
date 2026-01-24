package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yyhuni/orbit/server/internal/handler"
)

func RegisterPublicRoutes(
	api *gin.RouterGroup,
	screenshotHandler *handler.ScreenshotHandler,
	screenshotSnapshotHandler *handler.ScreenshotSnapshotHandler,
) {
	api.GET("/screenshots/:id/image", screenshotHandler.GetImage)
	api.GET("/scans/:id/screenshots/:snapshotId/image", screenshotSnapshotHandler.GetImage)
}
