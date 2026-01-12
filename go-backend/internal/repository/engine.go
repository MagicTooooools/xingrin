package repository

import (
	"github.com/xingrin/go-backend/internal/model"
	"github.com/xingrin/go-backend/internal/pkg/scope"
	"gorm.io/gorm"
)

// EngineRepository handles scan engine database operations
type EngineRepository struct {
	db *gorm.DB
}

// NewEngineRepository creates a new engine repository
func NewEngineRepository(db *gorm.DB) *EngineRepository {
	return &EngineRepository{db: db}
}

// Create creates a new engine
func (r *EngineRepository) Create(engine *model.ScanEngine) error {
	return r.db.Create(engine).Error
}

// FindByID finds an engine by ID
func (r *EngineRepository) FindByID(id int) (*model.ScanEngine, error) {
	var engine model.ScanEngine
	err := r.db.First(&engine, id).Error
	if err != nil {
		return nil, err
	}
	return &engine, nil
}

// FindAll finds all engines with pagination
func (r *EngineRepository) FindAll(page, pageSize int) ([]model.ScanEngine, int64, error) {
	var engines []model.ScanEngine
	var total int64

	if err := r.db.Model(&model.ScanEngine{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Scopes(
		scope.WithPagination(page, pageSize),
		scope.OrderByCreatedAtDesc(),
	).Find(&engines).Error

	return engines, total, err
}

// Update updates an engine
func (r *EngineRepository) Update(engine *model.ScanEngine) error {
	return r.db.Save(engine).Error
}

// Delete deletes an engine
func (r *EngineRepository) Delete(id int) error {
	return r.db.Delete(&model.ScanEngine{}, id).Error
}

// ExistsByName checks if engine name exists
func (r *EngineRepository) ExistsByName(name string, excludeID ...int) (bool, error) {
	var count int64
	query := r.db.Model(&model.ScanEngine{}).Where("name = ?", name)
	if len(excludeID) > 0 {
		query = query.Where("id != ?", excludeID[0])
	}
	err := query.Count(&count).Error
	return count > 0, err
}
