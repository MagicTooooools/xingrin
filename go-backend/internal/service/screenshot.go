package service

import (
	"errors"

	"github.com/xingrin/go-backend/internal/dto"
	"github.com/xingrin/go-backend/internal/model"
	"github.com/xingrin/go-backend/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrScreenshotNotFound = errors.New("screenshot not found")
)

// ScreenshotService handles screenshot business logic
type ScreenshotService struct {
	repo       *repository.ScreenshotRepository
	targetRepo *repository.TargetRepository
}

// NewScreenshotService creates a new screenshot service
func NewScreenshotService(repo *repository.ScreenshotRepository, targetRepo *repository.TargetRepository) *ScreenshotService {
	return &ScreenshotService{
		repo:       repo,
		targetRepo: targetRepo,
	}
}

// ListByTargetID returns paginated screenshots for a target
func (s *ScreenshotService) ListByTargetID(targetID int, query *dto.ScreenshotListQuery) ([]model.Screenshot, int64, error) {
	// Verify target exists
	_, err := s.targetRepo.FindByID(targetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, ErrTargetNotFound
		}
		return nil, 0, err
	}

	return s.repo.FindByTargetID(targetID, query.GetPage(), query.GetPageSize(), query.Filter)
}

// GetByID returns a screenshot by ID (including image data)
func (s *ScreenshotService) GetByID(id int) (*model.Screenshot, error) {
	screenshot, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrScreenshotNotFound
		}
		return nil, err
	}
	return screenshot, nil
}

// BulkDelete deletes multiple screenshots by IDs
func (s *ScreenshotService) BulkDelete(ids []int) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	return s.repo.BulkDelete(ids)
}

// BulkUpsert creates or updates multiple screenshots for a target
func (s *ScreenshotService) BulkUpsert(targetID int, req *dto.BulkUpsertScreenshotRequest) (int64, error) {
	// Verify target exists
	_, err := s.targetRepo.FindByID(targetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrTargetNotFound
		}
		return 0, err
	}

	// Convert DTO to models
	screenshots := make([]model.Screenshot, 0, len(req.Screenshots))
	for _, item := range req.Screenshots {
		screenshots = append(screenshots, model.Screenshot{
			TargetID:   targetID,
			URL:        item.URL,
			StatusCode: item.StatusCode,
			Image:      item.Image,
		})
	}

	return s.repo.BulkUpsert(screenshots)
}
