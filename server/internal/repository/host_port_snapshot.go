package repository

import (
	"database/sql"

	"github.com/yyhuni/orbit/server/internal/model"
	"github.com/yyhuni/orbit/server/internal/pkg/scope"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// HostPortSnapshotRepository handles host-port mapping snapshot database operations
type HostPortSnapshotRepository struct {
	db *gorm.DB
}

// NewHostPortSnapshotRepository creates a new host-port mapping snapshot repository
func NewHostPortSnapshotRepository(db *gorm.DB) *HostPortSnapshotRepository {
	return &HostPortSnapshotRepository{db: db}
}

// HostPortSnapshotFilterMapping defines field mapping for host-port snapshot filtering
var HostPortSnapshotFilterMapping = scope.FilterMapping{
	"host": {Column: "host"},
	"ip":   {Column: "ip"},
	"port": {Column: "port"},
}

// BulkCreate creates multiple host-port snapshots, ignoring duplicates
func (r *HostPortSnapshotRepository) BulkCreate(snapshots []model.HostPortSnapshot) (int64, error) {
	if len(snapshots) == 0 {
		return 0, nil
	}

	var totalAffected int64

	// Batch to avoid PostgreSQL parameter limit (65535)
	// 5 fields per record: scan_id, host, ip, port, created_at
	batchSize := 500
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

// FindByScanID finds host-port snapshots by scan ID with pagination and filter
func (r *HostPortSnapshotRepository) FindByScanID(scanID int, page, pageSize int, filter string) ([]model.HostPortSnapshot, int64, error) {
	var snapshots []model.HostPortSnapshot
	var total int64

	baseQuery := r.db.Model(&model.HostPortSnapshot{}).Where("scan_id = ?", scanID)
	baseQuery = baseQuery.Scopes(scope.WithFilter(filter, HostPortSnapshotFilterMapping))

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
func (r *HostPortSnapshotRepository) StreamByScanID(scanID int) (*sql.Rows, error) {
	return r.db.Model(&model.HostPortSnapshot{}).
		Where("scan_id = ?", scanID).
		Order("created_at DESC").
		Rows()
}

// CountByScanID returns the count of host-port snapshots for a scan
func (r *HostPortSnapshotRepository) CountByScanID(scanID int) (int64, error) {
	var count int64
	err := r.db.Model(&model.HostPortSnapshot{}).Where("scan_id = ?", scanID).Count(&count).Error
	return count, err
}

// ScanRow scans a single row into HostPortSnapshot model
func (r *HostPortSnapshotRepository) ScanRow(rows *sql.Rows) (*model.HostPortSnapshot, error) {
	var snapshot model.HostPortSnapshot
	if err := r.db.ScanRows(rows, &snapshot); err != nil {
		return nil, err
	}
	return &snapshot, nil
}
