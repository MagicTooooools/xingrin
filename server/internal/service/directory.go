package service

import (
	"database/sql"
	"errors"

	"github.com/xingrin/server/internal/dto"
	"github.com/xingrin/server/internal/model"
	"github.com/xingrin/server/internal/pkg/validator"
	"github.com/xingrin/server/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrDirectoryNotFound = errors.New("directory not found")
)

// DirectoryService handles directory business logic
type DirectoryService struct {
	repo       *repository.DirectoryRepository
	targetRepo *repository.TargetRepository
}

// NewDirectoryService creates a new directory service
func NewDirectoryService(repo *repository.DirectoryRepository, targetRepo *repository.TargetRepository) *DirectoryService {
	return &DirectoryService{repo: repo, targetRepo: targetRepo}
}

// ListByTarget returns paginated directories for a target
func (s *DirectoryService) ListByTarget(targetID int, query *dto.DirectoryListQuery) ([]model.Directory, int64, error) {
	// Check if target exists
	_, err := s.targetRepo.FindByID(targetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, ErrTargetNotFound
		}
		return nil, 0, err
	}
	return s.repo.FindByTargetID(targetID, query.GetPage(), query.GetPageSize(), query.Filter)
}

// BulkCreate creates multiple directories for a target
func (s *DirectoryService) BulkCreate(targetID int, urls []string) (int, error) {
	// Check if target exists and get target info
	target, err := s.targetRepo.FindByID(targetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrTargetNotFound
		}
		return 0, err
	}

	// Filter valid URLs that match target
	directories := make([]model.Directory, 0, len(urls))
	for _, u := range urls {
		if validator.IsURLMatchTarget(u, target.Name, target.Type) {
			directories = append(directories, model.Directory{
				TargetID: targetID,
				URL:      u,
			})
		}
	}

	if len(directories) == 0 {
		return 0, nil
	}

	return s.repo.BulkCreate(directories)
}

// BulkDelete deletes multiple directories by IDs
func (s *DirectoryService) BulkDelete(ids []int) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	return s.repo.BulkDelete(ids)
}

// StreamByTarget returns a cursor for streaming export
func (s *DirectoryService) StreamByTarget(targetID int) (*sql.Rows, error) {
	_, err := s.targetRepo.FindByID(targetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTargetNotFound
		}
		return nil, err
	}
	return s.repo.StreamByTargetID(targetID)
}

// CountByTarget returns the count of directories for a target
func (s *DirectoryService) CountByTarget(targetID int) (int64, error) {
	_, err := s.targetRepo.FindByID(targetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrTargetNotFound
		}
		return 0, err
	}
	return s.repo.CountByTargetID(targetID)
}

// ScanRow scans a row into Directory model
func (s *DirectoryService) ScanRow(rows *sql.Rows) (*model.Directory, error) {
	return s.repo.ScanRow(rows)
}

// BulkUpsert creates or updates multiple directories for a target
func (s *DirectoryService) BulkUpsert(targetID int, items []dto.DirectoryUpsertItem) (int64, error) {
	// Check if target exists and get target info
	target, err := s.targetRepo.FindByID(targetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrTargetNotFound
		}
		return 0, err
	}

	// Filter valid items that match target
	directories := make([]model.Directory, 0, len(items))
	for _, item := range items {
		if validator.IsURLMatchTarget(item.URL, target.Name, target.Type) {
			directories = append(directories, model.Directory{
				TargetID:      targetID,
				URL:           item.URL,
				Status:        item.Status,
				ContentLength: item.ContentLength,
				ContentType:   item.ContentType,
				Duration:      item.Duration,
			})
		}
	}

	if len(directories) == 0 {
		return 0, nil
	}

	return s.repo.BulkUpsert(directories)
}
