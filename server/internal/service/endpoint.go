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
	ErrEndpointNotFound = errors.New("endpoint not found")
)

// EndpointService handles endpoint business logic
type EndpointService struct {
	repo       *repository.EndpointRepository
	targetRepo *repository.TargetRepository
}

// NewEndpointService creates a new endpoint service
func NewEndpointService(repo *repository.EndpointRepository, targetRepo *repository.TargetRepository) *EndpointService {
	return &EndpointService{repo: repo, targetRepo: targetRepo}
}

// ListByTarget returns paginated endpoints for a target
func (s *EndpointService) ListByTarget(targetID int, query *dto.EndpointListQuery) ([]model.Endpoint, int64, error) {
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

// GetByID returns an endpoint by ID
func (s *EndpointService) GetByID(id int) (*model.Endpoint, error) {
	endpoint, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEndpointNotFound
		}
		return nil, err
	}
	return endpoint, nil
}

// BulkCreate creates multiple endpoints for a target
func (s *EndpointService) BulkCreate(targetID int, urls []string) (int, error) {
	// Check if target exists and get target info
	target, err := s.targetRepo.FindByID(targetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrTargetNotFound
		}
		return 0, err
	}

	// Filter valid URLs that match target
	endpoints := make([]model.Endpoint, 0, len(urls))
	for _, u := range urls {
		if validator.IsURLMatchTarget(u, target.Name, target.Type) {
			host := repository.ExtractHostFromURL(u)
			endpoints = append(endpoints, model.Endpoint{
				TargetID: targetID,
				URL:      u,
				Host:     host,
			})
		}
	}

	if len(endpoints) == 0 {
		return 0, nil
	}

	return s.repo.BulkCreate(endpoints)
}

// Delete deletes an endpoint by ID
func (s *EndpointService) Delete(id int) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrEndpointNotFound
		}
		return err
	}
	return s.repo.Delete(id)
}

// BulkDelete deletes multiple endpoints by IDs
func (s *EndpointService) BulkDelete(ids []int) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	return s.repo.BulkDelete(ids)
}

// StreamByTarget returns a cursor for streaming export
func (s *EndpointService) StreamByTarget(targetID int) (*sql.Rows, error) {
	_, err := s.targetRepo.FindByID(targetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTargetNotFound
		}
		return nil, err
	}
	return s.repo.StreamByTargetID(targetID)
}

// CountByTarget returns the count of endpoints for a target
func (s *EndpointService) CountByTarget(targetID int) (int64, error) {
	_, err := s.targetRepo.FindByID(targetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrTargetNotFound
		}
		return 0, err
	}
	return s.repo.CountByTargetID(targetID)
}

// ScanRow scans a row into Endpoint model
func (s *EndpointService) ScanRow(rows *sql.Rows) (*model.Endpoint, error) {
	return s.repo.ScanRow(rows)
}

// BulkUpsert creates or updates multiple endpoints for a target
func (s *EndpointService) BulkUpsert(targetID int, items []dto.EndpointUpsertItem) (int64, error) {
	// Check if target exists and get target info
	target, err := s.targetRepo.FindByID(targetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrTargetNotFound
		}
		return 0, err
	}

	// Filter valid items that match target
	endpoints := make([]model.Endpoint, 0, len(items))
	for _, item := range items {
		if validator.IsURLMatchTarget(item.URL, target.Name, target.Type) {
			host := repository.ExtractHostFromURL(item.URL)
			if item.Host != "" {
				host = item.Host
			}
			endpoints = append(endpoints, model.Endpoint{
				TargetID:          targetID,
				URL:               item.URL,
				Host:              host,
				Location:          item.Location,
				Title:             item.Title,
				Webserver:         item.Webserver,
				ContentType:       item.ContentType,
				StatusCode:        item.StatusCode,
				ContentLength:     item.ContentLength,
				ResponseBody:      item.ResponseBody,
				Tech:              item.Tech,
				Vhost:             item.Vhost,
				MatchedGFPatterns: item.MatchedGFPatterns,
				ResponseHeaders:   item.ResponseHeaders,
			})
		}
	}

	if len(endpoints) == 0 {
		return 0, nil
	}

	return s.repo.BulkUpsert(endpoints)
}
