package repository

import (
	"database/sql"

	"github.com/orbit/server/internal/model"
	"github.com/orbit/server/internal/pkg/scope"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// EndpointSnapshotRepository handles endpoint snapshot database operations
type EndpointSnapshotRepository struct {
	db *gorm.DB
}

// NewEndpointSnapshotRepository creates a new endpoint snapshot repository
func NewEndpointSnapshotRepository(db *gorm.DB) *EndpointSnapshotRepository {
	return &EndpointSnapshotRepository{db: db}
}

// EndpointSnapshotFilterMapping defines field mapping for endpoint snapshot filtering
var EndpointSnapshotFilterMapping = scope.FilterMapping{
	"url":       {Column: "url"},
	"host":      {Column: "host"},
	"title":     {Column: "title"},
	"status":    {Column: "status_code", IsNumeric: true},
	"webserver": {Column: "webserver"},
	"tech":      {Column: "tech", IsArray: true},
}

// BulkCreate creates multiple endpoint snapshots, ignoring duplicates
func (r *EndpointSnapshotRepository) BulkCreate(snapshots []model.EndpointSnapshot) (int64, error) {
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

// FindByScanID finds endpoint snapshots by scan ID with pagination and filter
func (r *EndpointSnapshotRepository) FindByScanID(scanID int, page, pageSize int, filter string) ([]model.EndpointSnapshot, int64, error) {
	var snapshots []model.EndpointSnapshot
	var total int64

	baseQuery := r.db.Model(&model.EndpointSnapshot{}).Where("scan_id = ?", scanID)
	baseQuery = baseQuery.Scopes(scope.WithFilter(filter, EndpointSnapshotFilterMapping))

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
func (r *EndpointSnapshotRepository) StreamByScanID(scanID int) (*sql.Rows, error) {
	return r.db.Model(&model.EndpointSnapshot{}).
		Where("scan_id = ?", scanID).
		Order("created_at DESC").
		Rows()
}

// CountByScanID returns the count of endpoint snapshots for a scan
func (r *EndpointSnapshotRepository) CountByScanID(scanID int) (int64, error) {
	var count int64
	err := r.db.Model(&model.EndpointSnapshot{}).Where("scan_id = ?", scanID).Count(&count).Error
	return count, err
}

// ScanRow scans a single row into EndpointSnapshot model
func (r *EndpointSnapshotRepository) ScanRow(rows *sql.Rows) (*model.EndpointSnapshot, error) {
	var snapshot model.EndpointSnapshot
	if err := r.db.ScanRows(rows, &snapshot); err != nil {
		return nil, err
	}
	return &snapshot, nil
}
