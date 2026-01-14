package repository

import (
	"database/sql"

	"github.com/xingrin/go-backend/internal/model"
	"github.com/xingrin/go-backend/internal/pkg/scope"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// WebsiteSnapshotRepository handles website snapshot database operations
type WebsiteSnapshotRepository struct {
	db *gorm.DB
}

// NewWebsiteSnapshotRepository creates a new website snapshot repository
func NewWebsiteSnapshotRepository(db *gorm.DB) *WebsiteSnapshotRepository {
	return &WebsiteSnapshotRepository{db: db}
}

// WebsiteSnapshotFilterMapping defines field mapping for website snapshot filtering
var WebsiteSnapshotFilterMapping = scope.FilterMapping{
	"url":       {Column: "url"},
	"host":      {Column: "host"},
	"title":     {Column: "title"},
	"status":    {Column: "status_code", IsNumeric: true},
	"webserver": {Column: "webserver"},
	"tech":      {Column: "tech", IsArray: true},
}

// BulkCreate creates multiple website snapshots, ignoring duplicates
// Uses ON CONFLICT DO NOTHING based on unique constraint (scan_id + url)
func (r *WebsiteSnapshotRepository) BulkCreate(snapshots []model.WebsiteSnapshot) (int64, error) {
	if len(snapshots) == 0 {
		return 0, nil
	}

	var totalAffected int64

	// Process in batches to avoid parameter limits
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

// FindByScanID finds website snapshots by scan ID with pagination and filter
func (r *WebsiteSnapshotRepository) FindByScanID(scanID int, page, pageSize int, filter string) ([]model.WebsiteSnapshot, int64, error) {
	var snapshots []model.WebsiteSnapshot
	var total int64

	// Base query
	baseQuery := r.db.Model(&model.WebsiteSnapshot{}).Where("scan_id = ?", scanID)

	// Apply filter scope
	baseQuery = baseQuery.Scopes(scope.WithFilter(filter, WebsiteSnapshotFilterMapping))

	// Count total
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch with pagination and default ordering
	err := baseQuery.Scopes(
		scope.WithPagination(page, pageSize),
		scope.OrderByCreatedAtDesc(),
	).Find(&snapshots).Error

	return snapshots, total, err
}

// StreamByScanID returns a sql.Rows cursor for streaming export
func (r *WebsiteSnapshotRepository) StreamByScanID(scanID int) (*sql.Rows, error) {
	return r.db.Model(&model.WebsiteSnapshot{}).
		Where("scan_id = ?", scanID).
		Order("created_at DESC").
		Rows()
}

// CountByScanID returns the count of website snapshots for a scan
func (r *WebsiteSnapshotRepository) CountByScanID(scanID int) (int64, error) {
	var count int64
	err := r.db.Model(&model.WebsiteSnapshot{}).Where("scan_id = ?", scanID).Count(&count).Error
	return count, err
}

// ScanRow scans a single row into WebsiteSnapshot model
func (r *WebsiteSnapshotRepository) ScanRow(rows *sql.Rows) (*model.WebsiteSnapshot, error) {
	var snapshot model.WebsiteSnapshot
	if err := r.db.ScanRows(rows, &snapshot); err != nil {
		return nil, err
	}
	return &snapshot, nil
}
