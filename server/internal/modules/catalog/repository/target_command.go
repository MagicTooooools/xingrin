package repository

import (
	"time"

	"github.com/yyhuni/lunafox/server/internal/modules/catalog/repository/persistence"
	"github.com/yyhuni/lunafox/server/internal/pkg/scope"
	"gorm.io/gorm/clause"
)

// Create creates a new target.
func (r *TargetRepository) Create(target *model.Target) error {
	return r.db.Create(target).Error
}

// Update updates a target.
func (r *TargetRepository) Update(target *model.Target) error {
	return r.db.Save(target).Error
}

// SoftDelete soft deletes a target.
func (r *TargetRepository) SoftDelete(id int) error {
	now := time.Now().UTC()
	return r.db.Model(&model.Target{}).Where("id = ?", id).Update("deleted_at", now).Error
}

// BulkSoftDelete soft deletes multiple targets by IDs.
func (r *TargetRepository) BulkSoftDelete(ids []int) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	now := time.Now().UTC()
	result := r.db.Model(&model.Target{}).
		Scopes(scope.WithNotDeleted()).
		Where("id IN ?", ids).
		Update("deleted_at", now)

	return result.RowsAffected, result.Error
}

// BulkCreateIgnoreConflicts creates multiple targets, ignoring duplicates.
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
