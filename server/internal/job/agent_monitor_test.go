package job

import (
	"context"
	"testing"
	"time"

	"github.com/yyhuni/orbit/server/internal/model"
	"github.com/yyhuni/orbit/server/internal/repository"
)

type fakeAgentRepo struct {
	agents  []*model.Agent
	updated []int
}

func (f *fakeAgentRepo) Create(ctx context.Context, agent *model.Agent) error {
	return nil
}

func (f *fakeAgentRepo) FindByID(ctx context.Context, id int) (*model.Agent, error) {
	return nil, nil
}

func (f *fakeAgentRepo) FindByAPIKey(ctx context.Context, apiKey string) (*model.Agent, error) {
	return nil, nil
}

func (f *fakeAgentRepo) List(ctx context.Context, page, pageSize int, status string) ([]*model.Agent, int64, error) {
	return nil, 0, nil
}

func (f *fakeAgentRepo) FindStaleOnline(ctx context.Context, before time.Time) ([]*model.Agent, error) {
	return f.agents, nil
}

func (f *fakeAgentRepo) Update(ctx context.Context, agent *model.Agent) error {
	return nil
}

func (f *fakeAgentRepo) UpdateStatus(ctx context.Context, id int, status string) error {
	f.updated = append(f.updated, id)
	return nil
}

func (f *fakeAgentRepo) UpdateHeartbeat(ctx context.Context, id int, update repository.AgentHeartbeatUpdate) error {
	return nil
}

func (f *fakeAgentRepo) Delete(ctx context.Context, id int) error {
	return nil
}

type fakeScanTaskRepo struct {
	recovered []int
}

func (f *fakeScanTaskRepo) FindByID(ctx context.Context, id int) (*model.ScanTask, error) {
	return nil, nil
}

func (f *fakeScanTaskRepo) PullTask(ctx context.Context, agentID int) (*model.ScanTask, error) {
	return nil, nil
}

func (f *fakeScanTaskRepo) UpdateStatus(ctx context.Context, id int, status string, errorMessage string) error {
	return nil
}

func (f *fakeScanTaskRepo) GetStatusCountsByScanID(ctx context.Context, scanID int) (pending, running, completed, failed, cancelled int, err error) {
	return 0, 0, 1, 0, 0, nil
}

func (f *fakeScanTaskRepo) CancelTasksByScanID(ctx context.Context, scanID int) ([]repository.CancelledTaskInfo, error) {
	return nil, nil
}

func (f *fakeScanTaskRepo) FailTasksForOfflineAgent(ctx context.Context, agentID int) error {
	f.recovered = append(f.recovered, agentID)
	return nil
}

func TestAgentMonitorMarksOfflineAndRecovers(t *testing.T) {
	agentRepo := &fakeAgentRepo{
		agents: []*model.Agent{{ID: 1}, {ID: 2}},
	}
	taskRepo := &fakeScanTaskRepo{}

	monitor := NewAgentMonitor(agentRepo, taskRepo, time.Minute, 2*time.Minute)
	monitor.check(context.Background())

	if len(agentRepo.updated) != 2 {
		t.Fatalf("expected 2 agents updated, got %d", len(agentRepo.updated))
	}
	if len(taskRepo.recovered) != 2 {
		t.Fatalf("expected 2 agents recovered, got %d", len(taskRepo.recovered))
	}
}
