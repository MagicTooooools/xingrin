package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yyhuni/lunafox/server/internal/middleware"
	agenthandler "github.com/yyhuni/lunafox/server/internal/modules/agent/handler"
	agentrepo "github.com/yyhuni/lunafox/server/internal/modules/agent/repository"
	scanhandler "github.com/yyhuni/lunafox/server/internal/modules/scan/handler"
)

// RegisterAgentRoutes registers agent-facing routes.
func RegisterAgentRoutes(
	api *gin.RouterGroup,
	protected *gin.RouterGroup,
	agentHandler *agenthandler.AgentHandler,
	agentWSHandler *agenthandler.AgentWebSocketHandler,
	agentTaskHandler *scanhandler.AgentTaskHandler,
	agentRepo agentrepo.AgentRepository,
) {
	api.POST("/agents/registrations", agentHandler.Register)
	api.GET("/agents/install-script", agentHandler.InstallScript)
	api.GET("/agents/ws", middleware.AgentAuthMiddleware(agentRepo), agentWSHandler.Handle)

	agentAPI := api.Group("/agents")
	agentAPI.Use(middleware.AgentAuthMiddleware(agentRepo))
	agentAPI.Use(middleware.AgentValidationMiddleware())
	{
		agentAPI.POST("/tasks/pull", agentTaskHandler.PullTask)
		agentAPI.PATCH("/tasks/:taskId/status", agentTaskHandler.UpdateTaskStatus)
	}

	protected.POST("/agents/registration-tokens", agentHandler.CreateRegistrationToken)
	protected.GET("/agents", agentHandler.ListAgents)
	protected.GET("/agents/:id", agentHandler.GetAgent)
	protected.DELETE("/agents/:id", agentHandler.DeleteAgent)
	protected.PUT("/agents/:id/config", agentHandler.UpdateAgentConfig)
}
