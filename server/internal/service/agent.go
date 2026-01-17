package service

import (
	"errors"

	"github.com/orbit/server/internal/model"
	"github.com/orbit/server/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrAgentScanNotFound       = errors.New("scan not found")
	ErrAgentInvalidTransition  = errors.New("invalid status transition")
)

// validTransitions defines allowed status transitions
// key: current status, value: allowed next statuses
var validTransitions = map[string][]string{
	model.ScanStatusPending:   {model.ScanStatusScheduled, model.ScanStatusCancelled},
	model.ScanStatusScheduled: {model.ScanStatusRunning, model.ScanStatusCancelled},
	model.ScanStatusRunning:   {model.ScanStatusCompleted, model.ScanStatusFailed, model.ScanStatusCancelled},
	// Terminal states: no transitions allowed
	model.ScanStatusCompleted: {},
	model.ScanStatusFailed:    {},
	model.ScanStatusCancelled: {},
}

// AgentService handles agent-related operations
type AgentService struct {
	scanRepo *repository.ScanRepository
}

// NewAgentService creates a new agent service
func NewAgentService(scanRepo *repository.ScanRepository) *AgentService {
	return &AgentService{
		scanRepo: scanRepo,
	}
}

// UpdateStatus updates scan status (called by Agent based on Worker exit code)
func (s *AgentService) UpdateStatus(scanID int, status string, errorMessage string) error {
	scan, err := s.scanRepo.FindByID(scanID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrAgentScanNotFound
		}
		return err
	}

	// Validate status transition
	if !isValidTransition(scan.Status, status) {
		return ErrAgentInvalidTransition
	}

	if errorMessage != "" {
		return s.scanRepo.UpdateStatus(scanID, status, errorMessage)
	}
	return s.scanRepo.UpdateStatus(scanID, status)
}

// isValidTransition checks if the status transition is allowed
func isValidTransition(current, next string) bool {
	allowed, exists := validTransitions[current]
	if !exists {
		return false
	}
	for _, s := range allowed {
		if s == next {
			return true
		}
	}
	return false
}
