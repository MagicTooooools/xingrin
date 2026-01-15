package repository

import (
	"github.com/xingrin/go-backend/internal/model"
	"github.com/xingrin/go-backend/internal/pkg/scope"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ScreenshotSnapshotRepository handles screenshot snapshot database operations
type ScreenshotSnapshotRepository struct {
	db *gorm.DB
}

// NewScreenshotSnapshotRepository creates a new screenshot snapshot repository
func NewScreenshotSnapshotRepository(db *gorm.DB) *ScreenshotSnapshotRepository {
	return &ScreenshotSnapshotRepository{db: db}
}

// ScreenshotSnapshotFilterMapping defines field mapping for screenshot snapshot filtering
var ScreenshotSnapshotFilterMapping = scope.FilterMapping{
	"url":    {Column: "url"},
	"status": {Column: "status_code", IsNumeric: true},
}

// BulkUpsert creates or updates multiple screenshot snapshots
// Uses ON CONFLICT (scan_id, url) DO UPDATE with COALESCE for non-null updates.
// Note: created_at is set only on insert (keeps the first capture time).
func (r *ScreenshotSnapshotRepository) BulkUpsert(snapshots []model.ScreenshotSnapshot) (int64, error) {
	if len(snapshots) == 0 {
		return 0, nil
	}

	var totalAffected int64

	// Process in batches to avoid parameter limits (4 fields per record)
	batchSize := 500
	for i := 0; i < len(snapshots); i += batchSize {
		end := i + batchSize
		if end > len(snapshots) {
			end = len(snapshots)
		}
		batch := snapshots[i:end]

		result := r.db.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "scan_id"}, {Name: "url"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"status_code": gorm.Expr("COALESCE(EXCLUDED.status_code, screenshot_snapshot.status_code)"),
				"image":       gorm.Expr("COALESCE(EXCLUDED.image, screenshot_snapshot.image)"),
			}),
		}).Create(&batch)
		if result.Error != nil {
			return totalAffected, result.Error
		}
		totalAffected += result.RowsAffected
	}

	return totalAffected, nil
}

// FindByScanID finds screenshot snapshots by scan ID with pagination and filter.
// This method intentionally excludes the image blob to avoid large payloads.
func (r *ScreenshotSnapshotRepository) FindByScanID(scanID int, page, pageSize int, filter string) ([]model.ScreenshotSnapshot, int64, error) {
	var snapshots []model.ScreenshotSnapshot
	var total int64

	baseQuery := r.db.Model(&model.ScreenshotSnapshot{}).Where("scan_id = ?", scanID)
	baseQuery = baseQuery.Scopes(scope.WithFilterDefault(filter, ScreenshotSnapshotFilterMapping, "url"))

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := baseQuery.
		Select("id, scan_id, url, status_code, created_at").
		Scopes(
			scope.WithPagination(page, pageSize),
			scope.OrderByCreatedAtDesc(),
		).
		Find(&snapshots).Error

	return snapshots, total, err
}

// FindByIDAndScanID finds a screenshot snapshot by ID under a scan (includes image data)
func (r *ScreenshotSnapshotRepository) FindByIDAndScanID(id int, scanID int) (*model.ScreenshotSnapshot, error) {
	var snapshot model.ScreenshotSnapshot
	err := r.db.Where("id = ? AND scan_id = ?", id, scanID).First(&snapshot).Error
	if err != nil {
		return nil, err
	}
	return &snapshot, nil
}

// CountByScanID returns the count of screenshot snapshots for a scan
func (r *ScreenshotSnapshotRepository) CountByScanID(scanID int) (int64, error) {
	var count int64
	err := r.db.Model(&model.ScreenshotSnapshot{}).Where("scan_id = ?", scanID).Count(&count).Error
	return count, err
}
