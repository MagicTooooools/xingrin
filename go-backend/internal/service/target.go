package service

import (
	"errors"
	"strings"

	"github.com/xingrin/go-backend/internal/dto"
	"github.com/xingrin/go-backend/internal/model"
	"github.com/xingrin/go-backend/internal/pkg/validator"
	"github.com/xingrin/go-backend/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrTargetNotFound = errors.New("target not found")
	ErrTargetExists   = errors.New("target name already exists")
	ErrInvalidTarget  = errors.New("invalid target format")
)

// TargetService handles target business logic
type TargetService struct {
	repo    *repository.TargetRepository
	orgRepo *repository.OrganizationRepository
}

// NewTargetService creates a new target service
func NewTargetService(repo *repository.TargetRepository, orgRepo *repository.OrganizationRepository) *TargetService {
	return &TargetService{repo: repo, orgRepo: orgRepo}
}

// Create creates a new target
func (s *TargetService) Create(req *dto.CreateTargetRequest) (*model.Target, error) {
	exists, err := s.repo.ExistsByName(req.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrTargetExists
	}

	// Auto-detect type from name
	targetType := validator.DetectTargetType(req.Name)
	if targetType == "" {
		return nil, ErrInvalidTarget
	}

	target := &model.Target{
		Name: req.Name,
		Type: targetType,
	}

	if err := s.repo.Create(target); err != nil {
		return nil, err
	}

	return target, nil
}

// List returns paginated targets
func (s *TargetService) List(query *dto.TargetListQuery) ([]model.Target, int64, error) {
	return s.repo.FindAll(query.GetOffset(), query.GetPageSize(), query.Type, query.Search)
}

// GetByID returns a target by ID
func (s *TargetService) GetByID(id int) (*model.Target, error) {
	target, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTargetNotFound
		}
		return nil, err
	}
	return target, nil
}

// Update updates a target
func (s *TargetService) Update(id int, req *dto.UpdateTargetRequest) (*model.Target, error) {
	target, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTargetNotFound
		}
		return nil, err
	}

	// Check name uniqueness if changed
	if target.Name != req.Name {
		exists, err := s.repo.ExistsByName(req.Name, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrTargetExists
		}
	}

	// Auto-detect type from name
	targetType := validator.DetectTargetType(req.Name)
	if targetType == "" {
		return nil, ErrInvalidTarget
	}

	target.Name = req.Name
	target.Type = targetType

	if err := s.repo.Update(target); err != nil {
		return nil, err
	}

	return target, nil
}

// Delete soft deletes a target
func (s *TargetService) Delete(id int) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTargetNotFound
		}
		return err
	}

	return s.repo.SoftDelete(id)
}

// BulkDelete soft deletes multiple targets by IDs
func (s *TargetService) BulkDelete(ids []int) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	return s.repo.BulkSoftDelete(ids)
}

// BatchCreate creates multiple targets at once
func (s *TargetService) BatchCreate(req *dto.BatchCreateTargetRequest) *dto.BatchCreateTargetResponse {
	var failedTargets []dto.FailedTarget
	var validTargets []model.Target
	var validNames []string
	seen := make(map[string]bool)

	// Step 1: Validate and deduplicate
	for _, item := range req.Targets {
		name := strings.TrimSpace(item.Name)
		if name == "" {
			continue
		}

		// Normalize: lowercase for domains
		normalized := strings.ToLower(name)

		// Skip duplicates in this batch
		if seen[normalized] {
			continue
		}
		seen[normalized] = true

		// Detect type
		targetType := validator.DetectTargetType(normalized)
		if targetType == "" {
			failedTargets = append(failedTargets, dto.FailedTarget{
				Name:   name,
				Reason: "unrecognized target format",
			})
			continue
		}

		validTargets = append(validTargets, model.Target{
			Name: normalized,
			Type: targetType,
		})
		validNames = append(validNames, normalized)
	}

	if len(validTargets) == 0 {
		return &dto.BatchCreateTargetResponse{
			CreatedCount:  0,
			FailedCount:   len(failedTargets),
			FailedTargets: failedTargets,
			Message:       "no valid targets",
		}
	}

	// Step 2: Validate organization exists (if provided)
	if req.OrganizationID != nil {
		exists, err := s.orgRepo.Exists(*req.OrganizationID)
		if err != nil {
			return &dto.BatchCreateTargetResponse{
				CreatedCount:  0,
				FailedCount:   len(req.Targets),
				FailedTargets: failedTargets,
				Message:       "failed to validate organization: " + err.Error(),
			}
		}
		if !exists {
			return &dto.BatchCreateTargetResponse{
				CreatedCount:  0,
				FailedCount:   len(req.Targets),
				FailedTargets: failedTargets,
				Message:       "organization not found",
			}
		}
	}

	// Step 3: Bulk create targets with ignore conflicts
	createdCount, err := s.repo.BulkCreateIgnoreConflicts(validTargets)
	if err != nil {
		return &dto.BatchCreateTargetResponse{
			CreatedCount:  0,
			FailedCount:   len(req.Targets),
			FailedTargets: failedTargets,
			Message:       "batch create failed: " + err.Error(),
		}
	}

	// Step 4: Associate targets with organization (if provided)
	if req.OrganizationID != nil {
		// Query all targets by names to get their IDs
		targets, err := s.repo.FindByNames(validNames)
		if err != nil {
			// Targets created but association failed - still return success for creation
			return &dto.BatchCreateTargetResponse{
				CreatedCount:  createdCount,
				FailedCount:   len(failedTargets),
				FailedTargets: failedTargets,
				Message:       "targets created, but failed to associate with organization: " + err.Error(),
			}
		}

		// Extract target IDs
		targetIDs := make([]int, len(targets))
		for i, t := range targets {
			targetIDs[i] = t.ID
		}

		// Bulk add targets to organization
		if err := s.orgRepo.BulkAddTargets(*req.OrganizationID, targetIDs); err != nil {
			return &dto.BatchCreateTargetResponse{
				CreatedCount:  createdCount,
				FailedCount:   len(failedTargets),
				FailedTargets: failedTargets,
				Message:       "targets created, but failed to associate with organization: " + err.Error(),
			}
		}
	}

	return &dto.BatchCreateTargetResponse{
		CreatedCount:  createdCount,
		FailedCount:   len(failedTargets),
		FailedTargets: failedTargets,
		Message:       "successfully created targets",
	}
}
