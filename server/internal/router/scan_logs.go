package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yyhuni/lunafox/server/internal/handler"
)

func RegisterScanLogRoutes(protected *gin.RouterGroup, scanLogHandler *handler.ScanLogHandler) {
	protected.GET("/scans/:id/logs", scanLogHandler.List)
	protected.POST("/scans/:id/logs", scanLogHandler.BulkCreate)
}
