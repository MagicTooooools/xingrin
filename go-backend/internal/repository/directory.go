package repository

import (
	"database/sql"

	"github.com/xingrin/go-backend/internal/model"
	"github.com/xingrin/go-backend/internal/pkg/scope"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// DirectoryRepository handles directory database operations
type DirectoryRepository struct {
	db *gorm.DB
}

// NewDirectoryRepository creates a new directory repository
func NewDirectoryRepository(db *gorm.DB) *DirectoryRepository {
	return &DirectoryRepository{db: db}
}

// DirectoryFilterMapping defines field mapping for directory filtering
var DirectoryFilterMapping = scope.FilterMapping{
	"url":    {Column: "url"},
	"status": {Column: "status", IsNumeric: true},
}

// FindByTargetID finds directories by target ID with pagination and filter
func (r *DirectoryRepository) FindByTargetID(targetID int, page, pageSize int, filter string) ([]model.Directory, int64, error) {
	var directories []model.Directory
	var total int64

	// Base query
	baseQuery := r.db.Model(&model.Directory{}).Where("target_id = ?", targetID)

	// Apply filter scope with default field "url"
	baseQuery = baseQuery.Scopes(scope.WithFilterDefault(filter, DirectoryFilterMapping, "url"))

	// Count total
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch with pagination and ordering
	err := baseQuery.Scopes(
		scope.WithPagination(page, pageSize),
		scope.OrderByCreatedAtDesc(),
	).Find(&directories).Error

	return directories, total, err
}

// BulkCreate creates multiple directories, ignoring duplicates
func (r *DirectoryRepository) BulkCreate(directories []model.Directory) (int, error) {
	if len(directories) == 0 {
		return 0, nil
	}

	// Use ON CONFLICT DO NOTHING to ignore duplicates (url + target_id unique)
	result := r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&directories)
	if result.Error != nil {
		return 0, result.Error
	}

	return int(result.RowsAffected), nil
}

// BulkDelete deletes multiple directories by IDs
func (r *DirectoryRepository) BulkDelete(ids []int) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	result := r.db.Where("id IN ?", ids).Delete(&model.Directory{})
	return result.RowsAffected, result.Error
}

// StreamByTargetID returns a sql.Rows cursor for streaming export
func (r *DirectoryRepository) StreamByTargetID(targetID int) (*sql.Rows, error) {
	return r.db.Model(&model.Directory{}).
		Where("target_id = ?", targetID).
		Order("created_at DESC").
		Rows()
}

// CountByTargetID returns the count of directories for a target
func (r *DirectoryRepository) CountByTargetID(targetID int) (int64, error) {
	var count int64
	err := r.db.Model(&model.Directory{}).Where("target_id = ?", targetID).Count(&count).Error
	return count, err
}

// ScanRow scans a single row into Directory model
func (r *DirectoryRepository) ScanRow(rows *sql.Rows) (*model.Directory, error) {
	var directory model.Directory
	if err := r.db.ScanRows(rows, &directory); err != nil {
		return nil, err
	}
	return &directory, nil
}

// BulkUpsert creates or updates multiple directories
// Uses ON CONFLICT DO UPDATE with COALESCE for non-null updates
func (r *DirectoryRepository) BulkUpsert(directories []model.Directory) (int64, error) {
	if len(directories) == 0 {
		return 0, nil
	}

	var totalAffected int64

	// Process in batches to avoid parameter limits
	batchSize := 100
	for i := 0; i < len(directories); i += batchSize {
		end := i + batchSize
		if end > len(directories) {
			end = len(directories)
		}
		batch := directories[i:end]

		affected, err := r.upsertBatch(batch)
		if err != nil {
			return totalAffected, err
		}
		totalAffected += affected
	}

	return totalAffected, nil
}

// upsertBatch upserts a single batch of directories
func (r *DirectoryRepository) upsertBatch(directories []model.Directory) (int64, error) {
	if len(directories) == 0 {
		return 0, nil
	}

	result := r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "url"}, {Name: "target_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"status":         gorm.Expr("COALESCE(EXCLUDED.status, directory.status)"),
			"content_length": gorm.Expr("COALESCE(EXCLUDED.content_length, directory.content_length)"),
			"content_type":   gorm.Expr("COALESCE(NULLIF(EXCLUDED.content_type, ''), directory.content_type)"),
			"duration":       gorm.Expr("COALESCE(EXCLUDED.duration, directory.duration)"),
		}),
	}).Create(&directories)

	return result.RowsAffected, result.Error
}
