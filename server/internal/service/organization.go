package service

import (
	"errors"

	"github.com/yyhuni/orbit/server/internal/dto"
	"github.com/yyhuni/orbit/server/internal/model"
	"github.com/yyhuni/orbit/server/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrOrganizationNotFound = errors.New("organization not found")
	ErrOrganizationExists   = errors.New("organization name already exists")
)

// OrganizationService handles organization business logic
type OrganizationService struct {
	repo *repository.OrganizationRepository
}

// NewOrganizationService creates a new organization service
func NewOrganizationService(repo *repository.OrganizationRepository) *OrganizationService {
	return &OrganizationService{repo: repo}
}

// Create creates a new organization
func (s *OrganizationService) Create(req *dto.CreateOrganizationRequest) (*model.Organization, error) {
	exists, err := s.repo.ExistsByName(req.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrOrganizationExists
	}

	org := &model.Organization{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.repo.Create(org); err != nil {
		return nil, err
	}

	return org, nil
}

// List returns paginated organizations with target count
func (s *OrganizationService) List(query *dto.OrganizationListQuery) ([]repository.OrganizationWithCount, int64, error) {
	return s.repo.FindAll(query.GetPage(), query.GetPageSize(), query.Filter)
}

// GetByID returns an organization by ID with target count
func (s *OrganizationService) GetByID(id int) (*repository.OrganizationWithCount, error) {
	org, err := s.repo.FindByIDWithCount(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrganizationNotFound
		}
		return nil, err
	}
	return org, nil
}

// Update updates an organization
func (s *OrganizationService) Update(id int, req *dto.UpdateOrganizationRequest) (*model.Organization, error) {
	org, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrganizationNotFound
		}
		return nil, err
	}

	// Check name uniqueness if changed
	if org.Name != req.Name {
		exists, err := s.repo.ExistsByName(req.Name, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrOrganizationExists
		}
	}

	org.Name = req.Name
	org.Description = req.Description

	if err := s.repo.Update(org); err != nil {
		return nil, err
	}

	return org, nil
}

// Delete soft deletes an organization
func (s *OrganizationService) Delete(id int) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrOrganizationNotFound
		}
		return err
	}

	return s.repo.SoftDelete(id)
}

// BulkDelete soft deletes multiple organizations
func (s *OrganizationService) BulkDelete(ids []int) (int64, error) {
	return s.repo.BulkSoftDelete(ids)
}

// ListTargets returns paginated targets for an organization
func (s *OrganizationService) ListTargets(organizationID int, query *dto.TargetListQuery) ([]model.Target, int64, error) {
	_, err := s.repo.FindByID(organizationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, ErrOrganizationNotFound
		}
		return nil, 0, err
	}

	return s.repo.FindTargets(organizationID, query.GetPage(), query.GetPageSize(), query.Type, query.Filter)
}

// LinkTargets adds targets to an organization
func (s *OrganizationService) LinkTargets(organizationID int, targetIDs []int) error {
	_, err := s.repo.FindByID(organizationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrOrganizationNotFound
		}
		return err
	}

	return s.repo.BulkAddTargets(organizationID, targetIDs)
}

// UnlinkTargets removes targets from an organization
func (s *OrganizationService) UnlinkTargets(organizationID int, targetIDs []int) (int64, error) {
	_, err := s.repo.FindByID(organizationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrOrganizationNotFound
		}
		return 0, err
	}

	return s.repo.UnlinkTargets(organizationID, targetIDs)
}
