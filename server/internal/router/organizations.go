package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yyhuni/orbit/server/internal/handler"
)

func RegisterOrganizationRoutes(protected *gin.RouterGroup, orgHandler *handler.OrganizationHandler) {
	protected.POST("/organizations", orgHandler.Create)
	protected.POST("/organizations/bulk-delete", orgHandler.BulkDelete)
	protected.GET("/organizations", orgHandler.List)
	protected.GET("/organizations/:id", orgHandler.GetByID)
	protected.GET("/organizations/:id/targets", orgHandler.ListTargets)
	protected.POST("/organizations/:id/link_targets", orgHandler.LinkTargets)
	protected.POST("/organizations/:id/unlink_targets", orgHandler.UnlinkTargets)
	protected.PUT("/organizations/:id", orgHandler.Update)
	protected.DELETE("/organizations/:id", orgHandler.Delete)
}
