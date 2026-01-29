package service

import (
	"context"
	"errors"
	"strings"

	"github.com/yyhuni/orbit/server/internal/dto"
	"github.com/yyhuni/orbit/server/internal/model"
	"github.com/yyhuni/orbit/server/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrScanTaskNotFound          = errors.New("scan task not found")
	ErrScanTaskNotOwned          = errors.New("scan task not owned by agent")
	ErrScanTaskInvalidTransition = errors.New("invalid scan task transition")
	ErrScanTaskInvalidUpdate     = errors.New("invalid scan task update")
)

// ScanTaskService handles scan task assignment and status updates.
type ScanTaskService struct {
	scanTaskRepo repository.ScanTaskRepository
	scanRepo     ScanRepository
}

// ScanRepository defines scan data access used by ScanTaskService.
type ScanRepository interface {
	FindByIDWithTarget(id int) (*model.Scan, error)
	UpdateStatus(id int, status string, errorMessage ...string) error
}

// NewScanTaskService creates a new scan task service.
func NewScanTaskService(scanTaskRepo repository.ScanTaskRepository, scanRepo ScanRepository) *ScanTaskService {
	return &ScanTaskService{
		scanTaskRepo: scanTaskRepo,
		scanRepo:     scanRepo,
	}
}

// PullTask assigns a pending task to the agent and returns task details.
func (s *ScanTaskService) PullTask(ctx context.Context, agentID int) (*dto.TaskAssignment, error) {
	task, err := s.scanTaskRepo.PullTask(ctx, agentID)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, nil
	}

	scan, err := s.scanRepo.FindByIDWithTarget(task.ScanID)
	if err != nil {
		return nil, err
	}

	if scan.Status == model.ScanStatusPending {
		if err := s.scanRepo.UpdateStatus(scan.ID, model.ScanStatusRunning); err != nil {
			return nil, err
		}
	}

	config := strings.TrimSpace(task.Config)
	if config == "" {
		config = strings.TrimSpace(scan.YamlConfiguration)
	}

	assignment := &dto.TaskAssignment{
		TaskID:       task.ID,
		ScanID:       task.ScanID,
		Stage:        task.Stage,
		WorkflowName: task.WorkflowName,
		TargetID:     scan.TargetID,
		WorkspaceDir: task.WorkspaceDir(),
		Config:       config,
	}

	if scan.Target != nil {
		assignment.TargetName = scan.Target.Name
		assignment.TargetType = scan.Target.Type
	}

	return assignment, nil
}

// UpdateStatus validates and updates task status for an agent.
func (s *ScanTaskService) UpdateStatus(ctx context.Context, agentID, taskID int, status, errorMessage string) error {
	status = strings.TrimSpace(status)
	errorMessage = strings.TrimSpace(errorMessage)

	if status == "" {
		return ErrScanTaskInvalidUpdate
	}
	if status == "failed" && errorMessage == "" {
		return ErrScanTaskInvalidUpdate
	}

	task, err := s.scanTaskRepo.FindByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrScanTaskNotFound
		}
		return err
	}

	if task.AgentID == nil || *task.AgentID != agentID {
		return ErrScanTaskNotOwned
	}

	if task.Status == status {
		return nil
	}

	if task.Status != "running" {
		return ErrScanTaskInvalidTransition
	}

	if status != "completed" && status != "failed" && status != "cancelled" {
		return ErrScanTaskInvalidTransition
	}

	if err := s.scanTaskRepo.UpdateStatus(ctx, taskID, status, errorMessage); err != nil {
		return err
	}

	switch status {
	case "completed":
		return s.scanRepo.UpdateStatus(task.ScanID, model.ScanStatusCompleted)
	case "failed":
		return s.scanRepo.UpdateStatus(task.ScanID, model.ScanStatusFailed, errorMessage)
	case "cancelled":
		return s.scanRepo.UpdateStatus(task.ScanID, model.ScanStatusCancelled)
	default:
		return nil
	}
}

// CancelRunningTasksForAgent cancels running tasks for an agent and updates scans.
func (s *ScanTaskService) CancelRunningTasksForAgent(ctx context.Context, agentID int) error {
	scanIDs, err := s.scanTaskRepo.CancelRunningTasksForAgent(ctx, agentID)
	if err != nil {
		return err
	}
	if len(scanIDs) == 0 {
		return nil
	}

	seen := map[int]struct{}{}
	for _, scanID := range scanIDs {
		if _, ok := seen[scanID]; ok {
			continue
		}
		seen[scanID] = struct{}{}
		_ = s.scanRepo.UpdateStatus(scanID, model.ScanStatusCancelled)
	}
	return nil
}
