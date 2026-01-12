package repository

import (
	"time"

	"github.com/xingrin/go-backend/internal/model"
	"github.com/xingrin/go-backend/internal/pkg/scope"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// TargetRepository handles target database operations
type TargetRepository struct {
	db *gorm.DB
}

// NewTargetRepository creates a new target repository
func NewTargetRepository(db *gorm.DB) *TargetRepository {
	return &TargetRepository{db: db}
}

// TargetFilterMapping defines field mapping for target filtering
var TargetFilterMapping = scope.FilterMapping{
	"name": {Column: "name"},
	"type": {Column: "type"},
}

// Create creates a new target
func (r *TargetRepository) Create(target *model.Target) error {
	return r.db.Create(target).Error
}

// FindByID finds a target by ID (excluding soft deleted)
func (r *TargetRepository) FindByID(id int) (*model.Target, error) {
	var target model.Target
	err := r.db.Scopes(scope.WithNotDeleted()).
		Where("id = ?", id).
		First(&target).Error
	if err != nil {
		return nil, err
	}
	return &target, nil
}

// FindAll finds all targets with pagination and filters (excluding soft deleted)
// Preloads organizations for each target
func (r *TargetRepository) FindAll(page, pageSize int, targetType, filter string) ([]model.Target, int64, error) {
	var targets []model.Target
	var total int64

	// Build base query with scopes
	baseQuery := r.db.Model(&model.Target{}).Scopes(scope.WithNotDeleted())

	// Apply type filter
	if targetType != "" {
		baseQuery = baseQuery.Where("type = ?", targetType)
	}

	// Apply smart filter (supports plain text as name search)
	if filter != "" {
		baseQuery = baseQuery.Scopes(scope.WithFilterDefault(filter, TargetFilterMapping, "name"))
	}

	// Count total
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch with preload and pagination
	err := baseQuery.
		Preload("Organizations", "deleted_at IS NULL").
		Scopes(
			scope.WithPagination(page, pageSize),
			scope.OrderByCreatedAtDesc(),
		).
		Find(&targets).Error

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

// BulkSoftDelete soft deletes multiple targets by IDs
func (r *TargetRepository) BulkSoftDelete(ids []int) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	now := time.Now()
	result := r.db.Model(&model.Target{}).
		Scopes(scope.WithNotDeleted()).
		Where("id IN ?", ids).
		Update("deleted_at", now)
	return result.RowsAffected, result.Error
}

// ExistsByName checks if target name exists (excluding soft deleted)
func (r *TargetRepository) ExistsByName(name string, excludeID ...int) (bool, error) {
	var count int64
	query := r.db.Model(&model.Target{}).
		Scopes(scope.WithNotDeleted()).
		Where("name = ?", name)
	if len(excludeID) > 0 {
		query = query.Where("id != ?", excludeID[0])
	}
	err := query.Count(&count).Error
	return count > 0, err
}

// BulkCreateIgnoreConflicts creates multiple targets, ignoring duplicates
func (r *TargetRepository) BulkCreateIgnoreConflicts(targets []model.Target) (int, error) {
	if len(targets) == 0 {
		return 0, nil
	}

	result := r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&targets)
	if result.Error != nil {
		return 0, result.Error
	}

	return int(result.RowsAffected), nil
}

// FindByNames finds targets by names (excluding soft deleted)
func (r *TargetRepository) FindByNames(names []string) ([]model.Target, error) {
	if len(names) == 0 {
		return nil, nil
	}

	var targets []model.Target
	err := r.db.Scopes(scope.WithNotDeleted()).
		Where("name IN ?", names).
		Find(&targets).Error
	return targets, err
}
