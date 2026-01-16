package repository

import (
	"time"

	"github.com/orbit/server/internal/model"
	"github.com/orbit/server/internal/pkg/scope"
	"gorm.io/gorm"
)

// ScanRepository handles scan database operations
type ScanRepository struct {
	db *gorm.DB
}

// NewScanRepository creates a new scan repository
func NewScanRepository(db *gorm.DB) *ScanRepository {
	return &ScanRepository{db: db}
}

// ScanFilterMapping defines field mapping for scan filtering
var ScanFilterMapping = scope.FilterMapping{
	"status":   {Column: "status"},
	"target":   {Column: "target_id"},
	"targetId": {Column: "target_id"},
}

// Create creates a new scan
func (r *ScanRepository) Create(scan *model.Scan) error {
	return r.db.Create(scan).Error
}

// FindByID finds a scan by ID (excluding soft deleted)
func (r *ScanRepository) FindByID(id int) (*model.Scan, error) {
	var scan model.Scan
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).
		First(&scan).Error
	if err != nil {
		return nil, err
	}
	return &scan, nil
}

// FindByIDWithTarget finds a scan by ID with target preloaded
func (r *ScanRepository) FindByIDWithTarget(id int) (*model.Scan, error) {
	var scan model.Scan
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).
		Preload("Target").
		First(&scan).Error
	if err != nil {
		return nil, err
	}
	return &scan, nil
}

// FindAll finds all scans with pagination and filters (excluding soft deleted)
func (r *ScanRepository) FindAll(page, pageSize int, targetID int, status, search string) ([]model.Scan, int64, error) {
	var scans []model.Scan
	var total int64

	// Build base query
	baseQuery := r.db.Model(&model.Scan{}).Where("scan.deleted_at IS NULL")

	// Apply target filter
	if targetID > 0 {
		baseQuery = baseQuery.Where("scan.target_id = ?", targetID)
	}

	// Apply status filter
	if status != "" {
		baseQuery = baseQuery.Where("scan.status = ?", status)
	}

	// Apply search filter (search by target name via join)
	if search != "" {
		baseQuery = baseQuery.Joins("LEFT JOIN target ON target.id = scan.target_id").
			Where("target.name ILIKE ?", "%"+search+"%")
	}

	// Count total
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch with preload and pagination
	err := baseQuery.
		Preload("Target").
		Scopes(
			scope.WithPagination(page, pageSize),
			scope.OrderByCreatedAtDesc(),
		).
		Find(&scans).Error

	return scans, total, err
}

// Update updates a scan
func (r *ScanRepository) Update(scan *model.Scan) error {
	return r.db.Save(scan).Error
}

// SoftDelete soft deletes a scan
func (r *ScanRepository) SoftDelete(id int) error {
	now := time.Now()
	return r.db.Model(&model.Scan{}).Where("id = ?", id).Update("deleted_at", now).Error
}

// BulkSoftDelete soft deletes multiple scans by IDs
func (r *ScanRepository) BulkSoftDelete(ids []int) (int64, []string, error) {
	if len(ids) == 0 {
		return 0, nil, nil
	}

	// Get scan names before deleting
	var scans []model.Scan
	if err := r.db.Select("id, target_id").
		Where("id IN ? AND deleted_at IS NULL", ids).
		Preload("Target", "deleted_at IS NULL").
		Find(&scans).Error; err != nil {
		return 0, nil, err
	}

	names := make([]string, 0, len(scans))
	for _, s := range scans {
		if s.Target != nil {
			names = append(names, s.Target.Name)
		}
	}

	// Soft delete
	now := time.Now()
	result := r.db.Model(&model.Scan{}).
		Where("id IN ? AND deleted_at IS NULL", ids).
		Update("deleted_at", now)

	return result.RowsAffected, names, result.Error
}

// UpdateStatus updates scan status
func (r *ScanRepository) UpdateStatus(id int, status string, errorMessage ...string) error {
	updates := map[string]interface{}{"status": status}
	if len(errorMessage) > 0 {
		updates["error_message"] = errorMessage[0]
	}
	if status == model.ScanStatusCompleted || status == model.ScanStatusFailed || status == model.ScanStatusStopped {
		now := time.Now()
		updates["stopped_at"] = &now
	}
	return r.db.Model(&model.Scan{}).Where("id = ?", id).Updates(updates).Error
}

// UpdateProgress updates scan progress
func (r *ScanRepository) UpdateProgress(id int, progress int, currentStage string) error {
	return r.db.Model(&model.Scan{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"progress":      progress,
			"current_stage": currentStage,
		}).Error
}

// GetStatistics returns scan statistics
func (r *ScanRepository) GetStatistics() (*ScanStatistics, error) {
	stats := &ScanStatistics{}

	// Count total (excluding soft deleted)
	if err := r.db.Model(&model.Scan{}).Where("deleted_at IS NULL").
		Count(&stats.Total).Error; err != nil {
		return nil, err
	}

	// Count by status
	if err := r.db.Model(&model.Scan{}).Where("deleted_at IS NULL AND status = ?", model.ScanStatusRunning).
		Count(&stats.Running).Error; err != nil {
		return nil, err
	}
	if err := r.db.Model(&model.Scan{}).Where("deleted_at IS NULL AND status = ?", model.ScanStatusCompleted).
		Count(&stats.Completed).Error; err != nil {
		return nil, err
	}
	if err := r.db.Model(&model.Scan{}).Where("deleted_at IS NULL AND status = ?", model.ScanStatusFailed).
		Count(&stats.Failed).Error; err != nil {
		return nil, err
	}

	// Sum cached counts from all scans
	type sumResult struct {
		TotalVulns      int64
		TotalSubdomains int64
		TotalEndpoints  int64
		TotalWebsites   int64
	}
	var sums sumResult
	if err := r.db.Model(&model.Scan{}).Where("deleted_at IS NULL").
		Select(`
			COALESCE(SUM(cached_vulns_total), 0) as total_vulns,
			COALESCE(SUM(cached_subdomains_count), 0) as total_subdomains,
			COALESCE(SUM(cached_endpoints_count), 0) as total_endpoints,
			COALESCE(SUM(cached_websites_count), 0) as total_websites
		`).
		Scan(&sums).Error; err != nil {
		return nil, err
	}

	stats.TotalVulns = sums.TotalVulns
	stats.TotalSubdomains = sums.TotalSubdomains
	stats.TotalEndpoints = sums.TotalEndpoints
	stats.TotalWebsites = sums.TotalWebsites
	stats.TotalAssets = sums.TotalSubdomains + sums.TotalEndpoints + sums.TotalWebsites

	return stats, nil
}

// ScanStatistics holds scan statistics
type ScanStatistics struct {
	Total           int64
	Running         int64
	Completed       int64
	Failed          int64
	TotalVulns      int64
	TotalSubdomains int64
	TotalEndpoints  int64
	TotalWebsites   int64
	TotalAssets     int64
}

// FindByTargetIDs finds scans by target IDs
func (r *ScanRepository) FindByTargetIDs(targetIDs []int) ([]model.Scan, error) {
	if len(targetIDs) == 0 {
		return nil, nil
	}

	var scans []model.Scan
	err := r.db.Where("target_id IN ? AND deleted_at IS NULL", targetIDs).
		Preload("Target").
		Order("created_at DESC").
		Find(&scans).Error
	return scans, err
}

// HasActiveScan checks if target has an active scan
func (r *ScanRepository) HasActiveScan(targetID int) (bool, error) {
	var count int64
	err := r.db.Model(&model.Scan{}).
		Where("target_id = ? AND deleted_at IS NULL AND status IN ?", targetID,
			[]string{model.ScanStatusInitiated, model.ScanStatusRunning, model.ScanStatusPending}).
		Count(&count).Error
	return count > 0, err
}

// GetTargetByScanID returns the target associated with a scan
func (r *ScanRepository) GetTargetByScanID(scanID int) (*model.Target, error) {
	var scan model.Scan
	err := r.db.Where("id = ? AND deleted_at IS NULL", scanID).
		Preload("Target").
		First(&scan).Error
	if err != nil {
		return nil, err
	}
	if scan.Target == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return scan.Target, nil
}
