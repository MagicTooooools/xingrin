package repository

import (
	"database/sql"

	"github.com/orbit/server/internal/model"
	"github.com/orbit/server/internal/pkg/scope"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// DirectorySnapshotRepository handles directory snapshot database operations
type DirectorySnapshotRepository struct {
	db *gorm.DB
}

// NewDirectorySnapshotRepository creates a new directory snapshot repository
func NewDirectorySnapshotRepository(db *gorm.DB) *DirectorySnapshotRepository {
	return &DirectorySnapshotRepository{db: db}
}

// DirectorySnapshotFilterMapping defines field mapping for directory snapshot filtering
var DirectorySnapshotFilterMapping = scope.FilterMapping{
	"url":         {Column: "url"},
	"status":      {Column: "status", IsNumeric: true},
	"contentType": {Column: "content_type"},
}

// BulkCreate creates multiple directory snapshots, ignoring duplicates
func (r *DirectorySnapshotRepository) BulkCreate(snapshots []model.DirectorySnapshot) (int64, error) {
	if len(snapshots) == 0 {
		return 0, nil
	}

	var totalAffected int64

	batchSize := 100
	for i := 0; i < len(snapshots); i += batchSize {
		end := i + batchSize
		if end > len(snapshots) {
			end = len(snapshots)
		}
		batch := snapshots[i:end]

		result := r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&batch)
		if result.Error != nil {
			return totalAffected, result.Error
		}
		totalAffected += result.RowsAffected
	}

	return totalAffected, nil
}

// FindByScanID finds directory snapshots by scan ID with pagination and filter
func (r *DirectorySnapshotRepository) FindByScanID(scanID int, page, pageSize int, filter string) ([]model.DirectorySnapshot, int64, error) {
	var snapshots []model.DirectorySnapshot
	var total int64

	baseQuery := r.db.Model(&model.DirectorySnapshot{}).Where("scan_id = ?", scanID)
	baseQuery = baseQuery.Scopes(scope.WithFilter(filter, DirectorySnapshotFilterMapping))

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := baseQuery.Scopes(
		scope.WithPagination(page, pageSize),
		scope.OrderByCreatedAtDesc(),
	).Find(&snapshots).Error

	return snapshots, total, err
}

// StreamByScanID returns a sql.Rows cursor for streaming export
func (r *DirectorySnapshotRepository) StreamByScanID(scanID int) (*sql.Rows, error) {
	return r.db.Model(&model.DirectorySnapshot{}).
		Where("scan_id = ?", scanID).
		Order("created_at DESC").
		Rows()
}

// CountByScanID returns the count of directory snapshots for a scan
func (r *DirectorySnapshotRepository) CountByScanID(scanID int) (int64, error) {
	var count int64
	err := r.db.Model(&model.DirectorySnapshot{}).Where("scan_id = ?", scanID).Count(&count).Error
	return count, err
}

// ScanRow scans a single row into DirectorySnapshot model
func (r *DirectorySnapshotRepository) ScanRow(rows *sql.Rows) (*model.DirectorySnapshot, error) {
	var snapshot model.DirectorySnapshot
	if err := r.db.ScanRows(rows, &snapshot); err != nil {
		return nil, err
	}
	return &snapshot, nil
}
