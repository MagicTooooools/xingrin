package service

import (
	"errors"
	"net"
	"regexp"
	"strings"

	"github.com/xingrin/go-backend/internal/dto"
	"github.com/xingrin/go-backend/internal/model"
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
	repo *repository.TargetRepository
}

// NewTargetService creates a new target service
func NewTargetService(repo *repository.TargetRepository) *TargetService {
	return &TargetService{repo: repo}
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

	// Auto-detect type if not provided
	targetType := req.Type
	if targetType == "" {
		targetType = detectTargetType(req.Name)
		if targetType == "" {
			return nil, ErrInvalidTarget
		}
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

	// Auto-detect type if not provided
	targetType := req.Type
	if targetType == "" {
		targetType = detectTargetType(req.Name)
		if targetType == "" {
			return nil, ErrInvalidTarget
		}
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

// detectTargetType auto-detects target type from name
func detectTargetType(name string) string {
	name = strings.TrimSpace(name)

	// Check CIDR
	if strings.Contains(name, "/") {
		_, _, err := net.ParseCIDR(name)
		if err == nil {
			return "cidr"
		}
	}

	// Check IP
	if ip := net.ParseIP(name); ip != nil {
		return "ip"
	}

	// Check domain
	if isValidDomain(name) {
		return "domain"
	}

	return ""
}

// isValidDomain validates domain format
func isValidDomain(domain string) bool {
	if len(domain) == 0 || len(domain) > 253 {
		return false
	}

	// Simple domain regex
	domainRegex := regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`)
	return domainRegex.MatchString(domain)
}
