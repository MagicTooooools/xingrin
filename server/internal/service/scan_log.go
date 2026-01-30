package service

import (
	"errors"

	"github.com/yyhuni/lunafox/server/internal/dto"
	"github.com/yyhuni/lunafox/server/internal/model"
	"github.com/yyhuni/lunafox/server/internal/repository"
	"gorm.io/gorm"
)

// ScanLogService handles scan log business logic
type ScanLogService struct {
	repo     *repository.ScanLogRepository
	scanRepo *repository.ScanRepository
}

// NewScanLogService creates a new scan log service
func NewScanLogService(repo *repository.ScanLogRepository, scanRepo *repository.ScanRepository) *ScanLogService {
	return &ScanLogService{
		repo:     repo,
		scanRepo: scanRepo,
	}
}

// ListByScanID returns logs for a scan with cursor pagination
func (s *ScanLogService) ListByScanID(scanID int, query *dto.ScanLogListQuery) (*dto.ScanLogListResponse, error) {
	// Check if scan exists
	_, err := s.scanRepo.FindByID(scanID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrScanNotFound
		}
		return nil, err
	}

	// Get limit (default 200, max 1000)
	limit := query.Limit
	if limit <= 0 {
		limit = 200
	}
	if limit > 1000 {
		limit = 1000
	}

	// Query logs with cursor pagination
	logs, err := s.repo.FindByScanIDWithCursor(scanID, query.AfterID, limit+1)
	if err != nil {
		return nil, err
	}

	// Check if there are more logs
	hasMore := len(logs) > limit
	if hasMore {
		logs = logs[:limit]
	}

	// Convert to response DTOs
	results := make([]dto.ScanLogResponse, len(logs))
	for i, log := range logs {
		results[i] = dto.ScanLogResponse{
			ID:        log.ID,
			ScanID:    log.ScanID,
			Level:     log.Level,
			Content:   log.Content,
			CreatedAt: log.CreatedAt,
		}
	}

	return &dto.ScanLogListResponse{
		Results: results,
		HasMore: hasMore,
	}, nil
}

// BulkCreate creates multiple logs for a scan
func (s *ScanLogService) BulkCreate(scanID int, items []dto.ScanLogItem) (int, error) {
	// Check if scan exists
	_, err := s.scanRepo.FindByID(scanID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrScanNotFound
		}
		return 0, err
	}

	if len(items) == 0 {
		return 0, nil
	}

	// Convert to models
	logs := make([]model.ScanLog, len(items))
	for i, item := range items {
		logs[i] = model.ScanLog{
			ScanID:  scanID,
			Level:   item.Level,
			Content: item.Content,
		}
	}

	if err := s.repo.BulkCreate(logs); err != nil {
		return 0, err
	}

	return len(logs), nil
}
