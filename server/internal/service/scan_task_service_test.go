package service

import (
	"context"
	"errors"
	"testing"

	"github.com/yyhuni/orbit/server/internal/model"
	"github.com/yyhuni/orbit/server/internal/repository"
	"gorm.io/gorm"
)

type mockScanTaskRepo struct {
	task         *model.ScanTask
	lastStatus   string
	lastErrorMsg string
	statusCount  map[string]int // for GetStatusCountsByScanID
	activeCount  int
}

func (m *mockScanTaskRepo) FindByID(ctx context.Context, id int) (*model.ScanTask, error) {
	if m.task == nil || m.task.ID != id {
		return nil, gorm.ErrRecordNotFound
	}
	return m.task, nil
}

func (m *mockScanTaskRepo) PullTask(ctx context.Context, agentID int) (*model.ScanTask, error) {
	return nil, errors.New("not implemented")
}

func (m *mockScanTaskRepo) UpdateStatus(ctx context.Context, id int, status string, errorMessage string) error {
	m.lastStatus = status
	m.lastErrorMsg = errorMessage
	if m.task != nil && m.task.ID == id {
		m.task.Status = status
	}
	return nil
}

func (m *mockScanTaskRepo) GetStatusCountsByScanID(ctx context.Context, scanID int) (pending, running, completed, failed, cancelled int, err error) {
	if m.statusCount != nil {
		return m.statusCount["pending"], m.statusCount["running"], m.statusCount["completed"], m.statusCount["failed"], m.statusCount["cancelled"], nil
	}
	// Default: all tasks completed (so scan status will be updated)
	return 0, 0, 1, 0, 0, nil
}

func (m *mockScanTaskRepo) CountActiveByScanAndStage(ctx context.Context, scanID, stage int) (int, error) {
	return m.activeCount, nil
}

func (m *mockScanTaskRepo) UnlockNextStage(ctx context.Context, scanID, stage int) (int64, error) {
	return 0, nil
}

func (m *mockScanTaskRepo) CancelTasksByScanID(ctx context.Context, scanID int) ([]repository.CancelledTaskInfo, error) {
	return nil, nil
}

func (m *mockScanTaskRepo) FailTasksForOfflineAgent(ctx context.Context, agentID int) error {
	return nil
}

type mockScanRepo struct {
	lastStatus string
	lastErr    string
	updated    []int
}

func (m *mockScanRepo) FindByIDWithTarget(id int) (*model.Scan, error) {
	return &model.Scan{ID: id}, nil
}

func (m *mockScanRepo) UpdateStatus(id int, status string, errorMessage ...string) error {
	m.lastStatus = status
	m.updated = append(m.updated, id)
	if len(errorMessage) > 0 {
		m.lastErr = errorMessage[0]
	}
	return nil
}

func TestScanTaskServiceUpdateStatusCompleted(t *testing.T) {
	task := &model.ScanTask{
		ID:     1,
		ScanID: 10,
		Status: "running",
	}
	agentID := 7
	task.AgentID = &agentID

	repo := &mockScanTaskRepo{task: task}
	scanRepo := &mockScanRepo{}
	svc := NewScanTaskService(repo, scanRepo)

	if err := svc.UpdateStatus(context.Background(), agentID, task.ID, "completed", ""); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if repo.lastStatus != "completed" {
		t.Fatalf("expected status updated, got %s", repo.lastStatus)
	}
	if scanRepo.lastStatus != model.ScanStatusCompleted {
		t.Fatalf("expected scan status completed, got %s", scanRepo.lastStatus)
	}
}

func TestScanTaskServiceUpdateStatusIdempotent(t *testing.T) {
	task := &model.ScanTask{
		ID:     2,
		ScanID: 11,
		Status: "completed",
	}
	agentID := 7
	task.AgentID = &agentID

	repo := &mockScanTaskRepo{task: task}
	scanRepo := &mockScanRepo{}
	svc := NewScanTaskService(repo, scanRepo)

	if err := svc.UpdateStatus(context.Background(), agentID, task.ID, "completed", ""); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if repo.lastStatus != "" {
		t.Fatalf("expected no update call")
	}
}

func TestScanTaskServiceUpdateStatusInvalidTransition(t *testing.T) {
	task := &model.ScanTask{
		ID:     3,
		ScanID: 12,
		Status: "pending",
	}
	agentID := 7
	task.AgentID = &agentID

	repo := &mockScanTaskRepo{task: task}
	svc := NewScanTaskService(repo, &mockScanRepo{})

	err := svc.UpdateStatus(context.Background(), agentID, task.ID, "completed", "")
	if !errors.Is(err, ErrScanTaskInvalidTransition) {
		t.Fatalf("expected invalid transition error, got %v", err)
	}
}

func TestScanTaskServiceUpdateStatusOwnership(t *testing.T) {
	task := &model.ScanTask{
		ID:     4,
		ScanID: 13,
		Status: "running",
	}
	ownerID := 9
	task.AgentID = &ownerID

	repo := &mockScanTaskRepo{task: task}
	svc := NewScanTaskService(repo, &mockScanRepo{})

	err := svc.UpdateStatus(context.Background(), 7, task.ID, "completed", "")
	if !errors.Is(err, ErrScanTaskNotOwned) {
		t.Fatalf("expected ownership error, got %v", err)
	}
}

func TestScanTaskServiceUpdateStatusFailedNeedsMessage(t *testing.T) {
	task := &model.ScanTask{
		ID:     5,
		ScanID: 14,
		Status: "running",
	}
	agentID := 7
	task.AgentID = &agentID

	repo := &mockScanTaskRepo{task: task}
	svc := NewScanTaskService(repo, &mockScanRepo{})

	err := svc.UpdateStatus(context.Background(), agentID, task.ID, "failed", "")
	if !errors.Is(err, ErrScanTaskInvalidUpdate) {
		t.Fatalf("expected invalid update error, got %v", err)
	}
}

func TestUpdateStatusScanNotUpdatedWhenTasksPending(t *testing.T) {
	task := &model.ScanTask{ID: 10, ScanID: 100, Status: "running"}
	agentID := 1
	task.AgentID = &agentID

	repo := &mockScanTaskRepo{
		task:        task,
		statusCount: map[string]int{"pending": 2, "completed": 1},
	}
	scanRepo := &mockScanRepo{}
	svc := NewScanTaskService(repo, scanRepo)

	if err := svc.UpdateStatus(context.Background(), agentID, task.ID, "completed", ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(scanRepo.updated) != 0 {
		t.Fatalf("expected no scan update when tasks still pending, got %d", len(scanRepo.updated))
	}
}

func TestUpdateStatusScanNotUpdatedWhenTasksRunning(t *testing.T) {
	task := &model.ScanTask{ID: 11, ScanID: 101, Status: "running"}
	agentID := 1
	task.AgentID = &agentID

	repo := &mockScanTaskRepo{
		task:        task,
		statusCount: map[string]int{"running": 1, "completed": 2},
	}
	scanRepo := &mockScanRepo{}
	svc := NewScanTaskService(repo, scanRepo)

	if err := svc.UpdateStatus(context.Background(), agentID, task.ID, "completed", ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(scanRepo.updated) != 0 {
		t.Fatalf("expected no scan update when tasks still running, got %d", len(scanRepo.updated))
	}
}

func TestUpdateStatusScanFailedWhenAnyTaskFailed(t *testing.T) {
	task := &model.ScanTask{ID: 12, ScanID: 102, Status: "running"}
	agentID := 1
	task.AgentID = &agentID

	repo := &mockScanTaskRepo{
		task:        task,
		statusCount: map[string]int{"completed": 3, "failed": 1},
	}
	scanRepo := &mockScanRepo{}
	svc := NewScanTaskService(repo, scanRepo)

	if err := svc.UpdateStatus(context.Background(), agentID, task.ID, "failed", "some error"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if scanRepo.lastStatus != model.ScanStatusFailed {
		t.Fatalf("expected scan status failed, got %s", scanRepo.lastStatus)
	}
}

func TestUpdateStatusScanCancelledWhenAllCancelled(t *testing.T) {
	task := &model.ScanTask{ID: 13, ScanID: 103, Status: "running"}
	agentID := 1
	task.AgentID = &agentID

	repo := &mockScanTaskRepo{
		task:        task,
		statusCount: map[string]int{"cancelled": 3},
	}
	scanRepo := &mockScanRepo{}
	svc := NewScanTaskService(repo, scanRepo)

	if err := svc.UpdateStatus(context.Background(), agentID, task.ID, "cancelled", ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if scanRepo.lastStatus != model.ScanStatusCancelled {
		t.Fatalf("expected scan status cancelled, got %s", scanRepo.lastStatus)
	}
}
