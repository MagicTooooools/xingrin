package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yyhuni/lunafox/server/internal/modules/asset/handler"
	endpointhandler "github.com/yyhuni/lunafox/server/internal/modules/asset/handler/endpoint"
	websitehandler "github.com/yyhuni/lunafox/server/internal/modules/asset/handler/website"
	snapshothandler "github.com/yyhuni/lunafox/server/internal/modules/snapshot/handler"
	snapshotrouter "github.com/yyhuni/lunafox/server/internal/modules/snapshot/router"
)

// RegisterAssetRoutes registers asset routes.
func RegisterAssetRoutes(
	api *gin.RouterGroup,
	protected *gin.RouterGroup,
	screenshotHandler *handler.ScreenshotHandler,
	screenshotSnapshotHandler *snapshothandler.ScreenshotSnapshotHandler,
	websiteHandler *websitehandler.WebsiteHandler,
	subdomainHandler *handler.SubdomainHandler,
	endpointHandler *endpointhandler.EndpointHandler,
	directoryHandler *handler.DirectoryHandler,
	hostPortHandler *handler.HostPortHandler,
	endpointSnapshotHandler *snapshothandler.EndpointSnapshotHandler,
	hostPortSnapshotHandler *snapshothandler.HostPortSnapshotHandler,
	websiteSnapshotHandler *snapshothandler.WebsiteSnapshotHandler,
	subdomainSnapshotHandler *snapshothandler.SubdomainSnapshotHandler,
	directorySnapshotHandler *snapshothandler.DirectorySnapshotHandler,
	vulnerabilitySnapshotHandler *snapshothandler.VulnerabilitySnapshotHandler,
) {
	RegisterPublicRoutes(api, screenshotHandler, screenshotSnapshotHandler)
	RegisterWebsiteRoutes(protected, websiteHandler)
	RegisterSubdomainRoutes(protected, subdomainHandler)
	RegisterEndpointRoutes(protected, endpointHandler, endpointSnapshotHandler)
	RegisterDirectoryRoutes(protected, directoryHandler)
	RegisterHostPortRoutes(protected, hostPortHandler, hostPortSnapshotHandler)
	RegisterScreenshotRoutes(protected, screenshotHandler, screenshotSnapshotHandler)
	snapshotrouter.RegisterScanSnapshotRoutes(
		protected,
		websiteSnapshotHandler,
		subdomainSnapshotHandler,
		endpointSnapshotHandler,
		directorySnapshotHandler,
		hostPortSnapshotHandler,
		screenshotSnapshotHandler,
		vulnerabilitySnapshotHandler,
	)
}
