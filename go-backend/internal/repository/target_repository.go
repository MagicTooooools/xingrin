package repository

import (
	"time"

	"github.com/xingrin/go-backend/internal/model"
	"gorm.io/gorm"
)

// TargetRepository handles target database operations
type TargetRepository struct {
	db *gorm.DB
}

// NewTargetRepository creates a new target repository
func NewTargetRepository(db *gorm.DB) *TargetRepository {
	return &TargetRepository{db: db}
}

// Create creates a new target
func (r *TargetRepository) Create(target *model.Target) error {
	return r.db.Create(target).Error
}

// FindByID finds a target by ID (excluding soft deleted)
func (r *TargetRepository) FindByID(id int) (*model.Target, error) {
	var target model.Target
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&target).Error
	if err != nil {
		return nil, err
	}
	return &target, nil
}

// FindAll finds all targets with pagination and filters (excluding soft deleted)
func (r *TargetRepository) FindAll(offset, limit int, targetType, search string) ([]model.Target, int64, error) {
	var targets []model.Target
	var total int64

	query := r.db.Model(&model.Target{}).Where("deleted_at IS NULL")

	if targetType != "" {
		query = query.Where("type = ?", targetType)
	}
	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&targets).Error
	return targets, total, err
}

// Update updates a target
func (r *TargetRepository) Update(target *model.Target) error {
	return r.db.Save(target).Error
}

// SoftDelete soft deletes a target
func (r *TargetRepository) SoftDelete(id int) error {
	now := time.Now()
	return r.db.Model(&model.Target{}).Where("id = ?", id).Update("deleted_at", now).Error
}

// ExistsByName checks if target name exists (excluding soft deleted)
func (r *TargetRepository) ExistsByName(name string, excludeID ...int) (bool, error) {
	var count int64
	query := r.db.Model(&model.Target{}).Where("name = ? AND deleted_at IS NULL", name)
	if len(excludeID) > 0 {
		query = query.Where("id != ?", excludeID[0])
	}
	err := query.Count(&count).Error
	return count > 0, err
}
