package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yyhuni/orbit/server/internal/handler"
)

func RegisterScanRoutes(protected *gin.RouterGroup, scanHandler *handler.ScanHandler) {
	protected.GET("/scans", scanHandler.List)
	protected.POST("/scans", scanHandler.Create)
	protected.GET("/scans/statistics", scanHandler.Statistics)
	protected.GET("/scans/:id", scanHandler.GetByID)
	protected.DELETE("/scans/:id", scanHandler.Delete)
	protected.POST("/scans/:id/stop", scanHandler.Stop)
	protected.POST("/scans/bulk-delete", scanHandler.BulkDelete)
}
