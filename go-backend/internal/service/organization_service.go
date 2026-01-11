package service

import (
	"errors"

	"github.com/xingrin/go-backend/internal/dto"
	"github.com/xingrin/go-backend/internal/model"
	"github.com/xingrin/go-backend/internal/repository"
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

// List returns paginated organizations
func (s *OrganizationService) List(query *dto.PaginationQuery) ([]model.Organization, int64, error) {
	return s.repo.FindAll(query.GetOffset(), query.GetPageSize())
}

// GetByID returns an organization by ID
func (s *OrganizationService) GetByID(id int) (*model.Organization, error) {
	org, err := s.repo.FindByID(id)
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
	// Check if exists
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrOrganizationNotFound
		}
		return err
	}

	return s.repo.SoftDelete(id)
}
