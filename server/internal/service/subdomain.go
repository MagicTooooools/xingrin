package service

import (
	"database/sql"
	"errors"

	"github.com/yyhuni/lunafox/server/internal/dto"
	"github.com/yyhuni/lunafox/server/internal/model"
	"github.com/yyhuni/lunafox/server/internal/pkg/validator"
	"github.com/yyhuni/lunafox/server/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrSubdomainNotFound = errors.New("subdomain not found")
	ErrInvalidTargetType = errors.New("target type must be domain for subdomains")
	ErrSubdomainNotMatch = errors.New("subdomain does not match target domain")
)

// SubdomainService handles subdomain business logic
type SubdomainService struct {
	repo       *repository.SubdomainRepository
	targetRepo *repository.TargetRepository
}

// NewSubdomainService creates a new subdomain service
func NewSubdomainService(repo *repository.SubdomainRepository, targetRepo *repository.TargetRepository) *SubdomainService {
	return &SubdomainService{repo: repo, targetRepo: targetRepo}
}

// ListByTarget returns paginated subdomains for a target
func (s *SubdomainService) ListByTarget(targetID int, query *dto.SubdomainListQuery) ([]model.Subdomain, int64, error) {
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

// BulkCreate creates multiple subdomains for a target
func (s *SubdomainService) BulkCreate(targetID int, names []string) (int, error) {
	// Check if target exists and get target info
	target, err := s.targetRepo.FindByID(targetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrTargetNotFound
		}
		return 0, err
	}

	// Validate target type must be domain
	if target.Type != "domain" {
		return 0, ErrInvalidTargetType
	}

	// Filter valid subdomains that match target domain
	subdomains := make([]model.Subdomain, 0, len(names))
	for _, name := range names {
		// Check if subdomain matches target domain (includes DNS name validation)
		if validator.IsSubdomainOfTarget(name, target.Name) {
			subdomains = append(subdomains, model.Subdomain{
				TargetID: targetID,
				Name:     name,
			})
		}
	}

	if len(subdomains) == 0 {
		return 0, nil
	}

	return s.repo.BulkCreate(subdomains)
}

// BulkDelete deletes multiple subdomains by IDs
func (s *SubdomainService) BulkDelete(ids []int) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	return s.repo.BulkDelete(ids)
}

// StreamByTarget returns a cursor for streaming export
func (s *SubdomainService) StreamByTarget(targetID int) (*sql.Rows, error) {
	_, err := s.targetRepo.FindByID(targetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTargetNotFound
		}
		return nil, err
	}
	return s.repo.StreamByTargetID(targetID)
}

// CountByTarget returns the count of subdomains for a target
func (s *SubdomainService) CountByTarget(targetID int) (int64, error) {
	_, err := s.targetRepo.FindByID(targetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrTargetNotFound
		}
		return 0, err
	}
	return s.repo.CountByTargetID(targetID)
}

// ScanRow scans a row into Subdomain model
func (s *SubdomainService) ScanRow(rows *sql.Rows) (*model.Subdomain, error) {
	return s.repo.ScanRow(rows)
}
