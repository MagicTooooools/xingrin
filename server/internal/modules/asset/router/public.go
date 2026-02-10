package router

import (
	"github.com/gin-gonic/gin"
	assethandler "github.com/yyhuni/lunafox/server/internal/modules/asset/handler"
	snapshothandler "github.com/yyhuni/lunafox/server/internal/modules/snapshot/handler"
)

func RegisterPublicRoutes(
	api *gin.RouterGroup,
	screenshotHandler *assethandler.ScreenshotHandler,
	screenshotSnapshotHandler *snapshothandler.ScreenshotSnapshotHandler,
) {
	api.GET("/screenshots/:id/image", screenshotHandler.GetImage)
	api.GET("/scans/:id/screenshots/:snapshotId/image", screenshotSnapshotHandler.GetImage)
}
