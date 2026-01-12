package service

import (
	"database/sql"
	"errors"

	"github.com/xingrin/go-backend/internal/dto"
	"github.com/xingrin/go-backend/internal/model"
	"github.com/xingrin/go-backend/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrWebsiteNotFound = errors.New("website not found")
)

// WebsiteService handles website business logic
type WebsiteService struct {
	repo       *repository.WebsiteRepository
	targetRepo *repository.TargetRepository
}

// NewWebsiteService creates a new website service
func NewWebsiteService(repo *repository.WebsiteRepository, targetRepo *repository.TargetRepository) *WebsiteService {
	return &WebsiteService{repo: repo, targetRepo: targetRepo}
}

// ListByTarget returns paginated websites for a target
func (s *WebsiteService) ListByTarget(targetID int, query *dto.WebsiteListQuery) ([]model.Website, int64, error) {
	return s.repo.FindByTargetID(targetID, query.GetPage(), query.GetPageSize(), query.Filter)
}

// GetByID returns a website by ID
func (s *WebsiteService) GetByID(id int) (*model.Website, error) {
	website, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrWebsiteNotFound
		}
		return nil, err
	}
	return website, nil
}

// BulkCreate creates multiple websites for a target
func (s *WebsiteService) BulkCreate(targetID int, urls []string) (int, error) {
	_, err := s.targetRepo.FindByID(targetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrTargetNotFound
		}
		return 0, err
	}

	websites := make([]model.Website, 0, len(urls))
	for _, u := range urls {
		host := repository.ExtractHostFromURL(u)
		websites = append(websites, model.Website{
			TargetID: targetID,
			URL:      u,
			Host:     host,
		})
	}

	return s.repo.BulkCreate(websites)
}

// Delete deletes a website by ID
func (s *WebsiteService) Delete(id int) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrWebsiteNotFound
		}
		return err
	}
	return s.repo.Delete(id)
}

// BulkDelete deletes multiple websites by IDs
func (s *WebsiteService) BulkDelete(ids []int) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	return s.repo.BulkDelete(ids)
}

// StreamByTarget returns a cursor for streaming export
func (s *WebsiteService) StreamByTarget(targetID int) (*sql.Rows, error) {
	_, err := s.targetRepo.FindByID(targetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTargetNotFound
		}
		return nil, err
	}
	return s.repo.StreamByTargetID(targetID)
}

// CountByTarget returns the count of websites for a target
func (s *WebsiteService) CountByTarget(targetID int) (int64, error) {
	// Check if target exists first
	_, err := s.targetRepo.FindByID(targetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrTargetNotFound
		}
		return 0, err
	}
	return s.repo.CountByTargetID(targetID)
}

// ScanRow scans a row into Website model
func (s *WebsiteService) ScanRow(rows *sql.Rows) (*model.Website, error) {
	return s.repo.ScanRow(rows)
}
