package application

import (
	"context"

	scandto "github.com/yyhuni/lunafox/server/internal/modules/scan/dto"
)

// AgentTaskService orchestrates task pull/status updates for agent runtime endpoints.
type AgentTaskService struct {
	runtime ScanTaskRuntimePort
}

func NewAgentTaskService(runtime ScanTaskRuntimePort) *AgentTaskService {
	return &AgentTaskService{runtime: runtime}
}

func (service *AgentTaskService) PullTask(ctx context.Context, agentID int) (*scandto.TaskAssignment, error) {
	return service.runtime.PullTask(ctx, agentID)
}

func (service *AgentTaskService) UpdateStatus(ctx context.Context, agentID, taskID int, status, errorMessage string) error {
	return service.runtime.UpdateStatus(ctx, agentID, taskID, status, errorMessage)
}
