package repository

import (
	"time"

	"github.com/yyhuni/lunafox/server/internal/model"
	"github.com/yyhuni/lunafox/server/internal/pkg/scope"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// TargetRepository handles target database operations
type TargetRepository struct {
	db *gorm.DB
}

// NewTargetRepository creates a new target repository
func NewTargetRepository(db *gorm.DB) *TargetRepository {
	return &TargetRepository{db: db}
}

// TargetFilterMapping defines field mapping for target filtering
var TargetFilterMapping = scope.FilterMapping{
	"name": {Column: "name"},
	"type": {Column: "type"},
}

// Create creates a new target
func (r *TargetRepository) Create(target *model.Target) error {
	return r.db.Create(target).Error
}

// FindByID finds a target by ID (excluding soft deleted)
func (r *TargetRepository) FindByID(id int) (*model.Target, error) {
	var target model.Target
	err := r.db.Scopes(scope.WithNotDeleted()).
		Where("id = ?", id).
		First(&target).Error
	if err != nil {
		return nil, err
	}
	return &target, nil
}

// FindAll finds all targets with pagination and filters (excluding soft deleted)
// Preloads organizations for each target
func (r *TargetRepository) FindAll(page, pageSize int, targetType, filter string) ([]model.Target, int64, error) {
	var targets []model.Target
	var total int64

	// Build base query with scopes
	baseQuery := r.db.Model(&model.Target{}).Scopes(scope.WithNotDeleted())

	// Apply type filter
	if targetType != "" {
		baseQuery = baseQuery.Where("type = ?", targetType)
	}

	// Apply smart filter (supports plain text as name search)
	if filter != "" {
		baseQuery = baseQuery.Scopes(scope.WithFilterDefault(filter, TargetFilterMapping, "name"))
	}

	// Count total
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch with preload and pagination
	err := baseQuery.
		Preload("Organizations", "deleted_at IS NULL").
		Scopes(
			scope.WithPagination(page, pageSize),
			scope.OrderByCreatedAtDesc(),
		).
		Find(&targets).Error

	return targets, total, err
}

// Update updates a target
func (r *TargetRepository) Update(target *model.Target) error {
	return r.db.Save(target).Error
}

// SoftDelete soft deletes a target
func (r *TargetRepository) SoftDelete(id int) error {
	now := time.Now()
	return r.db.Model(&model.Target{}).Where("id = ?", id).Update("deleted_at", now).Error
}

// BulkSoftDelete soft deletes multiple targets by IDs
func (r *TargetRepository) BulkSoftDelete(ids []int) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	now := time.Now()
	result := r.db.Model(&model.Target{}).
		Scopes(scope.WithNotDeleted()).
		Where("id IN ?", ids).
		Update("deleted_at", now)
	return result.RowsAffected, result.Error
}

// ExistsByName checks if target name exists (excluding soft deleted)
func (r *TargetRepository) ExistsByName(name string, excludeID ...int) (bool, error) {
	var count int64
	query := r.db.Model(&model.Target{}).
		Scopes(scope.WithNotDeleted()).
		Where("name = ?", name)
	if len(excludeID) > 0 {
		query = query.Where("id != ?", excludeID[0])
	}
	err := query.Count(&count).Error
	return count > 0, err
}

// BulkCreateIgnoreConflicts creates multiple targets, ignoring duplicates
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

// FindByNames finds targets by names (excluding soft deleted)
func (r *TargetRepository) FindByNames(names []string) ([]model.Target, error) {
	if len(names) == 0 {
		return nil, nil
	}

	var targets []model.Target
	err := r.db.Scopes(scope.WithNotDeleted()).
		Where("name IN ?", names).
		Find(&targets).Error
	return targets, err
}

// TargetAssetCounts holds asset count statistics for a target
type TargetAssetCounts struct {
	Subdomains  int64
	Websites    int64
	Endpoints   int64
	IPs         int64
	Directories int64
	Screenshots int64
}

// VulnerabilityCounts holds vulnerability count statistics by severity
type VulnerabilityCounts struct {
	Total    int64
	Critical int64
	High     int64
	Medium   int64
	Low      int64
}

// GetAssetCounts returns asset counts for a target
func (r *TargetRepository) GetAssetCounts(targetID int) (*TargetAssetCounts, error) {
	counts := &TargetAssetCounts{}

	// Count subdomains
	if err := r.db.Table("subdomain").
		Where("target_id = ?", targetID).
		Count(&counts.Subdomains).Error; err != nil {
		return nil, err
	}

	// Count websites
	if err := r.db.Table("website").
		Where("target_id = ?", targetID).
		Count(&counts.Websites).Error; err != nil {
		return nil, err
	}

	// Count endpoints
	if err := r.db.Table("endpoint").
		Where("target_id = ?", targetID).
		Count(&counts.Endpoints).Error; err != nil {
		return nil, err
	}

	// Count distinct IPs from host_port_mapping
	if err := r.db.Table("host_port_mapping").
		Where("target_id = ?", targetID).
		Select("COUNT(DISTINCT ip)").
		Scan(&counts.IPs).Error; err != nil {
		return nil, err
	}

	// Count directories
	if err := r.db.Table("directory").
		Where("target_id = ?", targetID).
		Count(&counts.Directories).Error; err != nil {
		return nil, err
	}

	// Count screenshots
	if err := r.db.Table("screenshot").
		Where("target_id = ?", targetID).
		Count(&counts.Screenshots).Error; err != nil {
		return nil, err
	}

	return counts, nil
}

// GetVulnerabilityCounts returns vulnerability counts by severity for a target
func (r *TargetRepository) GetVulnerabilityCounts(targetID int) (*VulnerabilityCounts, error) {
	counts := &VulnerabilityCounts{}

	// Count total vulnerabilities
	if err := r.db.Table("vulnerability").
		Where("target_id = ?", targetID).
		Count(&counts.Total).Error; err != nil {
		return nil, err
	}

	// Count by severity
	type severityCount struct {
		Severity string
		Count    int64
	}
	var severityCounts []severityCount

	if err := r.db.Table("vulnerability").
		Select("severity, COUNT(*) as count").
		Where("target_id = ?", targetID).
		Group("severity").
		Scan(&severityCounts).Error; err != nil {
		return nil, err
	}

	for _, sc := range severityCounts {
		switch sc.Severity {
		case "critical":
			counts.Critical = sc.Count
		case "high":
			counts.High = sc.Count
		case "medium":
			counts.Medium = sc.Count
		case "low":
			counts.Low = sc.Count
		}
	}

	return counts, nil
}
