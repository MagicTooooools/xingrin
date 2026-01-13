package service

import (
	"database/sql"
	"errors"

	"github.com/xingrin/go-backend/internal/dto"
	"github.com/xingrin/go-backend/internal/model"
	"github.com/xingrin/go-backend/internal/pkg/validator"
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
	target, err := s.targetRepo.FindByID(targetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrTargetNotFound
		}
		return 0, err
	}

	// Filter valid URLs that match target
	websites := make([]model.Website, 0, len(urls))
	for _, u := range urls {
		if validator.IsURLMatchTarget(u, target.Name, target.Type) {
			host := repository.ExtractHostFromURL(u)
			websites = append(websites, model.Website{
				TargetID: targetID,
				URL:      u,
				Host:     host,
			})
		}
	}

	if len(websites) == 0 {
		return 0, nil
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

// BulkUpsert creates or updates multiple websites for a target (for scanner import)
func (s *WebsiteService) BulkUpsert(targetID int, items []dto.WebsiteUpsertItem) (int64, error) {
	target, err := s.targetRepo.FindByID(targetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrTargetNotFound
		}
		return 0, err
	}

	// Filter valid URLs that match target and convert to models
	websites := make([]model.Website, 0, len(items))
	for _, item := range items {
		if validator.IsURLMatchTarget(item.URL, target.Name, target.Type) {
			host := item.Host
			if host == "" {
				host = repository.ExtractHostFromURL(item.URL)
			}

			websites = append(websites, model.Website{
				TargetID:        targetID,
				URL:             item.URL,
				Host:            host,
				Location:        item.Location,
				Title:           item.Title,
				Webserver:       item.Webserver,
				ResponseBody:    item.ResponseBody,
				ContentType:     item.ContentType,
				Tech:            item.Tech,
				StatusCode:      item.StatusCode,
				ContentLength:   item.ContentLength,
				Vhost:           item.Vhost,
				ResponseHeaders: item.ResponseHeaders,
			})
		}
	}

	if len(websites) == 0 {
		return 0, nil
	}

	return s.repo.BulkUpsert(websites)
}
