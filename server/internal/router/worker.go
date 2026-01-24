package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yyhuni/orbit/server/internal/handler"
	"github.com/yyhuni/orbit/server/internal/middleware"
)

func RegisterWorkerRoutes(
	api *gin.RouterGroup,
	workerToken string,
	workerHandler *handler.WorkerHandler,
	wordlistHandler *handler.WordlistHandler,
	subdomainSnapshotHandler *handler.SubdomainSnapshotHandler,
	websiteSnapshotHandler *handler.WebsiteSnapshotHandler,
	endpointSnapshotHandler *handler.EndpointSnapshotHandler,
) {
	workerAPI := api.Group("/worker")
	workerAPI.Use(middleware.WorkerAuthMiddleware(workerToken))
	{
		workerAPI.GET("/scans/:id/target-name", workerHandler.GetTargetName)
		workerAPI.GET("/scans/:id/provider-config", workerHandler.GetProviderConfig)
		workerAPI.GET("/wordlists/:name", wordlistHandler.GetByName)
		workerAPI.GET("/wordlists/:name/download", wordlistHandler.DownloadByName)
		workerAPI.POST("/scans/:id/subdomains/bulk-upsert", subdomainSnapshotHandler.BulkUpsert)
		workerAPI.POST("/scans/:id/websites/bulk-upsert", websiteSnapshotHandler.BulkUpsert)
		workerAPI.POST("/scans/:id/endpoints/bulk-upsert", endpointSnapshotHandler.BulkUpsert)
	}
}
