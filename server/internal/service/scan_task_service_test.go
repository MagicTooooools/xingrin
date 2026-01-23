package service

import (
	"context"
	"errors"
	"testing"

	"github.com/yyhuni/orbit/server/internal/model"
	"gorm.io/gorm"
)

type mockScanTaskRepo struct {
	task         *model.ScanTask
	lastStatus   string
	lastErrorMsg string
	cancelIDs    []int
}

func (m *mockScanTaskRepo) Create(ctx context.Context, task *model.ScanTask) error {
	return nil
}

func (m *mockScanTaskRepo) FindByID(ctx context.Context, id int) (*model.ScanTask, error) {
	if m.task == nil || m.task.ID != id {
		return nil, gorm.ErrRecordNotFound
	}
	return m.task, nil
}

func (m *mockScanTaskRepo) FindByScanID(ctx context.Context, scanID int) ([]*model.ScanTask, error) {
	return nil, errors.New("not implemented")
}

func (m *mockScanTaskRepo) FindByAgentID(ctx context.Context, agentID int) ([]*model.ScanTask, error) {
	return nil, errors.New("not implemented")
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

func (m *mockScanTaskRepo) CancelRunningTasksForAgent(ctx context.Context, agentID int) ([]int, error) {
	return m.cancelIDs, nil
}

func (m *mockScanTaskRepo) RecoverTasksForOfflineAgent(ctx context.Context, agentID int) error {
	return nil
}

func (m *mockScanTaskRepo) Delete(ctx context.Context, id int) error {
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

func TestCancelRunningTasksForAgentDedup(t *testing.T) {
	repo := &mockScanTaskRepo{
		cancelIDs: []int{1, 1, 2},
	}
	scanRepo := &mockScanRepo{}
	svc := NewScanTaskService(repo, scanRepo)

	if err := svc.CancelRunningTasksForAgent(context.Background(), 9); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(scanRepo.updated) != 2 {
		t.Fatalf("expected 2 scan updates, got %d", len(scanRepo.updated))
	}
}
