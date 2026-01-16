package repository

import (
	"database/sql"

	"github.com/orbit/server/internal/model"
	"github.com/orbit/server/internal/pkg/scope"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SubdomainSnapshotRepository handles subdomain snapshot database operations
type SubdomainSnapshotRepository struct {
	db *gorm.DB
}

// NewSubdomainSnapshotRepository creates a new subdomain snapshot repository
func NewSubdomainSnapshotRepository(db *gorm.DB) *SubdomainSnapshotRepository {
	return &SubdomainSnapshotRepository{db: db}
}

// SubdomainSnapshotFilterMapping defines field mapping for subdomain snapshot filtering
var SubdomainSnapshotFilterMapping = scope.FilterMapping{
	"name": {Column: "name"},
}

// BulkCreate creates multiple subdomain snapshots, ignoring duplicates
func (r *SubdomainSnapshotRepository) BulkCreate(snapshots []model.SubdomainSnapshot) (int64, error) {
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

// FindByScanID finds subdomain snapshots by scan ID with pagination and filter
func (r *SubdomainSnapshotRepository) FindByScanID(scanID int, page, pageSize int, filter string) ([]model.SubdomainSnapshot, int64, error) {
	var snapshots []model.SubdomainSnapshot
	var total int64

	baseQuery := r.db.Model(&model.SubdomainSnapshot{}).Where("scan_id = ?", scanID)
	baseQuery = baseQuery.Scopes(scope.WithFilter(filter, SubdomainSnapshotFilterMapping))

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
func (r *SubdomainSnapshotRepository) StreamByScanID(scanID int) (*sql.Rows, error) {
	return r.db.Model(&model.SubdomainSnapshot{}).
		Where("scan_id = ?", scanID).
		Order("created_at DESC").
		Rows()
}

// CountByScanID returns the count of subdomain snapshots for a scan
func (r *SubdomainSnapshotRepository) CountByScanID(scanID int) (int64, error) {
	var count int64
	err := r.db.Model(&model.SubdomainSnapshot{}).Where("scan_id = ?", scanID).Count(&count).Error
	return count, err
}

// ScanRow scans a single row into SubdomainSnapshot model
func (r *SubdomainSnapshotRepository) ScanRow(rows *sql.Rows) (*model.SubdomainSnapshot, error) {
	var snapshot model.SubdomainSnapshot
	if err := r.db.ScanRows(rows, &snapshot); err != nil {
		return nil, err
	}
	return &snapshot, nil
}
