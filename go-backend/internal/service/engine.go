package service

import (
	"errors"

	"github.com/xingrin/go-backend/internal/dto"
	"github.com/xingrin/go-backend/internal/model"
	"github.com/xingrin/go-backend/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrEngineNotFound = errors.New("engine not found")
	ErrEngineExists   = errors.New("engine name already exists")
)

// EngineService handles engine business logic
type EngineService struct {
	repo *repository.EngineRepository
}

// NewEngineService creates a new engine service
func NewEngineService(repo *repository.EngineRepository) *EngineService {
	return &EngineService{repo: repo}
}

// Create creates a new engine
func (s *EngineService) Create(req *dto.CreateEngineRequest) (*model.ScanEngine, error) {
	exists, err := s.repo.ExistsByName(req.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEngineExists
	}

	engine := &model.ScanEngine{
		Name:          req.Name,
		Configuration: req.Configuration,
	}

	if err := s.repo.Create(engine); err != nil {
		return nil, err
	}

	return engine, nil
}

// List returns paginated engines
func (s *EngineService) List(query *dto.PaginationQuery) ([]model.ScanEngine, int64, error) {
	return s.repo.FindAll(query.GetPage(), query.GetPageSize())
}

// GetByID returns an engine by ID
func (s *EngineService) GetByID(id int) (*model.ScanEngine, error) {
	engine, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEngineNotFound
		}
		return nil, err
	}
	return engine, nil
}

// Update updates an engine
func (s *EngineService) Update(id int, req *dto.UpdateEngineRequest) (*model.ScanEngine, error) {
	engine, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEngineNotFound
		}
		return nil, err
	}

	if engine.Name != req.Name {
		exists, err := s.repo.ExistsByName(req.Name, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrEngineExists
		}
	}

	engine.Name = req.Name
	engine.Configuration = req.Configuration

	if err := s.repo.Update(engine); err != nil {
		return nil, err
	}

	return engine, nil
}

// Patch partially updates an engine
func (s *EngineService) Patch(id int, req *dto.PatchEngineRequest) (*model.ScanEngine, error) {
	engine, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEngineNotFound
		}
		return nil, err
	}

	// Only update fields that are provided
	if req.Name != nil && *req.Name != engine.Name {
		exists, err := s.repo.ExistsByName(*req.Name, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrEngineExists
		}
		engine.Name = *req.Name
	}

	if req.Configuration != nil {
		engine.Configuration = *req.Configuration
	}

	if err := s.repo.Update(engine); err != nil {
		return nil, err
	}

	return engine, nil
}

// Delete deletes an engine
func (s *EngineService) Delete(id int) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrEngineNotFound
		}
		return err
	}

	return s.repo.Delete(id)
}
