package repository

import (
	"database/sql"

	"github.com/orbit/server/internal/model"
	"github.com/orbit/server/internal/pkg/scope"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SubdomainRepository handles subdomain database operations
type SubdomainRepository struct {
	db *gorm.DB
}

// NewSubdomainRepository creates a new subdomain repository
func NewSubdomainRepository(db *gorm.DB) *SubdomainRepository {
	return &SubdomainRepository{db: db}
}

// SubdomainFilterMapping defines field mapping for subdomain filtering
var SubdomainFilterMapping = scope.FilterMapping{
	"name": {Column: "name"},
}

// FindByTargetID finds subdomains by target ID with pagination and filter
func (r *SubdomainRepository) FindByTargetID(targetID int, page, pageSize int, filter string) ([]model.Subdomain, int64, error) {
	var subdomains []model.Subdomain
	var total int64

	// Base query
	baseQuery := r.db.Model(&model.Subdomain{}).Where("target_id = ?", targetID)

	// Apply filter scope with default field "name"
	baseQuery = baseQuery.Scopes(scope.WithFilterDefault(filter, SubdomainFilterMapping, "name"))

	// Count total
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch with pagination and ordering
	err := baseQuery.Scopes(
		scope.WithPagination(page, pageSize),
		scope.OrderByCreatedAtDesc(),
	).Find(&subdomains).Error

	return subdomains, total, err
}

// BulkCreate creates multiple subdomains, ignoring duplicates
func (r *SubdomainRepository) BulkCreate(subdomains []model.Subdomain) (int, error) {
	if len(subdomains) == 0 {
		return 0, nil
	}

	// Use ON CONFLICT DO NOTHING to ignore duplicates (name + target_id unique)
	result := r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&subdomains)
	if result.Error != nil {
		return 0, result.Error
	}

	return int(result.RowsAffected), nil
}

// BulkDelete deletes multiple subdomains by IDs
func (r *SubdomainRepository) BulkDelete(ids []int) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	result := r.db.Where("id IN ?", ids).Delete(&model.Subdomain{})
	return result.RowsAffected, result.Error
}

// StreamByTargetID returns a sql.Rows cursor for streaming export
func (r *SubdomainRepository) StreamByTargetID(targetID int) (*sql.Rows, error) {
	return r.db.Model(&model.Subdomain{}).
		Where("target_id = ?", targetID).
		Order("created_at DESC").
		Rows()
}

// CountByTargetID returns the count of subdomains for a target
func (r *SubdomainRepository) CountByTargetID(targetID int) (int64, error) {
	var count int64
	err := r.db.Model(&model.Subdomain{}).Where("target_id = ?", targetID).Count(&count).Error
	return count, err
}

// ScanRow scans a single row into Subdomain model
func (r *SubdomainRepository) ScanRow(rows *sql.Rows) (*model.Subdomain, error) {
	var subdomain model.Subdomain
	if err := r.db.ScanRows(rows, &subdomain); err != nil {
		return nil, err
	}
	return &subdomain, nil
}
