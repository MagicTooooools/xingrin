package repository

import (
	"github.com/orbit/server/internal/model"
	"github.com/orbit/server/internal/pkg/scope"
	"gorm.io/gorm"
)

// ScanLogRepository handles scan log database operations
type ScanLogRepository struct {
	db *gorm.DB
}

// NewScanLogRepository creates a new scan log repository
func NewScanLogRepository(db *gorm.DB) *ScanLogRepository {
	return &ScanLogRepository{db: db}
}

// Create creates a new scan log
func (r *ScanLogRepository) Create(log *model.ScanLog) error {
	return r.db.Create(log).Error
}

// BulkCreate creates multiple scan logs
func (r *ScanLogRepository) BulkCreate(logs []model.ScanLog) error {
	if len(logs) == 0 {
		return nil
	}
	return r.db.CreateInBatches(logs, 100).Error
}

// FindByScanID finds logs by scan ID with pagination
func (r *ScanLogRepository) FindByScanID(scanID int, page, pageSize int, level string) ([]model.ScanLog, int64, error) {
	var logs []model.ScanLog
	var total int64

	// Build base query
	baseQuery := r.db.Model(&model.ScanLog{}).Where("scan_id = ?", scanID)

	// Apply level filter
	if level != "" {
		baseQuery = baseQuery.Where("level = ?", level)
	}

	// Count total
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch with pagination (ordered by created_at ASC for logs)
	err := baseQuery.
		Scopes(scope.WithPagination(page, pageSize)).
		Order("created_at ASC").
		Find(&logs).Error

	return logs, total, err
}

// FindByScanIDWithCursor finds logs by scan ID with cursor pagination
func (r *ScanLogRepository) FindByScanIDWithCursor(scanID int, afterID int64, limit int) ([]model.ScanLog, error) {
	var logs []model.ScanLog

	query := r.db.Where("scan_id = ?", scanID)

	// Apply cursor filter
	if afterID > 0 {
		query = query.Where("id > ?", afterID)
	}

	// Order by ID (auto-increment, guarantees consistent order)
	err := query.Order("id ASC").Limit(limit).Find(&logs).Error
	return logs, err
}

// DeleteByScanID deletes all logs for a scan
func (r *ScanLogRepository) DeleteByScanID(scanID int) error {
	return r.db.Where("scan_id = ?", scanID).Delete(&model.ScanLog{}).Error
}

// DeleteByScanIDs deletes all logs for multiple scans
func (r *ScanLogRepository) DeleteByScanIDs(scanIDs []int) error {
	if len(scanIDs) == 0 {
		return nil
	}
	return r.db.Where("scan_id IN ?", scanIDs).Delete(&model.ScanLog{}).Error
}
