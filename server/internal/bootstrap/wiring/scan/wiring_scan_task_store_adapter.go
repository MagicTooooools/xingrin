package scanwiring

import (
	"context"

	scanapp "github.com/yyhuni/lunafox/server/internal/modules/scan/application"
	scanrepo "github.com/yyhuni/lunafox/server/internal/modules/scan/repository"
)

type scanTaskStoreAdapter struct{ repo scanrepo.ScanTaskRepository }

func newScanTaskStoreAdapter(repo scanrepo.ScanTaskRepository) *scanTaskStoreAdapter {
	return &scanTaskStoreAdapter{repo: repo}
}

func (adapter *scanTaskStoreAdapter) FindByID(ctx context.Context, id int) (*scanapp.TaskRecord, error) {
	task, err := adapter.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &scanapp.TaskRecord{ID: task.ID, ScanID: task.ScanID, Stage: task.Stage, WorkflowName: task.WorkflowName, Status: task.Status, AgentID: task.AgentID, Config: task.Config}, nil
}

func (adapter *scanTaskStoreAdapter) PullTask(ctx context.Context, agentID int) (*scanapp.TaskRecord, error) {
	task, err := adapter.repo.PullTask(ctx, agentID)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, nil
	}
	return &scanapp.TaskRecord{ID: task.ID, ScanID: task.ScanID, Stage: task.Stage, WorkflowName: task.WorkflowName, Status: task.Status, AgentID: task.AgentID, Config: task.Config}, nil
}

func (adapter *scanTaskStoreAdapter) UpdateStatus(ctx context.Context, id int, status string, errorMessage string) error {
	return adapter.repo.UpdateStatus(ctx, id, status, errorMessage)
}

func (adapter *scanTaskStoreAdapter) GetStatusCountsByScanID(ctx context.Context, scanID int) (pending, running, completed, failed, cancelled int, err error) {
	return adapter.repo.GetStatusCountsByScanID(ctx, scanID)
}

func (adapter *scanTaskStoreAdapter) CountActiveByScanAndStage(ctx context.Context, scanID, stage int) (int, error) {
	return adapter.repo.CountActiveByScanAndStage(ctx, scanID, stage)
}

func (adapter *scanTaskStoreAdapter) UnlockNextStage(ctx context.Context, scanID, stage int) (int64, error) {
	return adapter.repo.UnlockNextStage(ctx, scanID, stage)
}
