package repository

import (
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/yyhuni/orbit/server/internal/model"
	"github.com/yyhuni/orbit/server/internal/pkg/scope"
	"gorm.io/gorm"
)

// ErrTargetNotFound indicates one or more target IDs do not exist
var ErrTargetNotFound = errors.New("one or more target IDs do not exist")

// PostgreSQL error codes
const (
	pgForeignKeyViolation = "23503"
)

// OrganizationFilterMapping defines filter fields for organization
var OrganizationFilterMapping = scope.FilterMapping{
	"name": {Column: "organization.name", IsArray: false},
}

// OrganizationWithCount represents organization with target count
type OrganizationWithCount struct {
	model.Organization
	TargetCount int64 `gorm:"column:target_count"`
}

// OrganizationRepository handles organization database operations
type OrganizationRepository struct {
	db *gorm.DB
}

// NewOrganizationRepository creates a new organization repository
func NewOrganizationRepository(db *gorm.DB) *OrganizationRepository {
	return &OrganizationRepository{db: db}
}

// Create creates a new organization
func (r *OrganizationRepository) Create(org *model.Organization) error {
	return r.db.Create(org).Error
}

// FindByID finds an organization by ID (excluding soft deleted)
func (r *OrganizationRepository) FindByID(id int) (*model.Organization, error) {
	var org model.Organization
	err := r.db.Scopes(scope.WithNotDeleted()).
		Where("id = ?", id).
		First(&org).Error
	if err != nil {
		return nil, err
	}
	return &org, nil
}

// FindByIDWithCount finds an organization by ID with target count (excluding soft deleted)
func (r *OrganizationRepository) FindByIDWithCount(id int) (*OrganizationWithCount, error) {
	var org OrganizationWithCount
	err := r.db.Table("organization").
		Select(`organization.*,
			(SELECT COUNT(*) FROM organization_target
			 INNER JOIN target ON target.id = organization_target.target_id
			 WHERE organization_target.organization_id = organization.id
			 AND target.deleted_at IS NULL) as target_count`).
		Where("organization.id = ? AND organization.deleted_at IS NULL", id).
		First(&org).Error
	if err != nil {
		return nil, err
	}
	return &org, nil
}

// FindAll finds all organizations with pagination and target count (excluding soft deleted)
func (r *OrganizationRepository) FindAll(page, pageSize int, filter string) ([]OrganizationWithCount, int64, error) {
	var orgs []OrganizationWithCount
	var total int64

	// Base query for counting
	countQuery := r.db.Model(&model.Organization{}).
		Scopes(scope.WithNotDeleted()).
		Scopes(scope.WithFilterDefault(filter, OrganizationFilterMapping, "name"))
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Query with target count using subquery (excluding soft deleted targets)
	query := r.db.Table("organization").
		Select(`organization.*,
			(SELECT COUNT(*) FROM organization_target
			 INNER JOIN target ON target.id = organization_target.target_id
			 WHERE organization_target.organization_id = organization.id
			 AND target.deleted_at IS NULL) as target_count`).
		Where("organization.deleted_at IS NULL").
		Scopes(scope.WithFilterDefault(filter, OrganizationFilterMapping, "name"))

	err := query.Scopes(
		scope.WithPagination(page, pageSize),
		scope.OrderBy("organization.created_at", true),
	).Find(&orgs).Error

	return orgs, total, err
}

// Update updates an organization
func (r *OrganizationRepository) Update(org *model.Organization) error {
	return r.db.Save(org).Error
}

// SoftDelete soft deletes an organization
func (r *OrganizationRepository) SoftDelete(id int) error {
	now := time.Now()
	return r.db.Model(&model.Organization{}).Where("id = ?", id).Update("deleted_at", now).Error
}

// BulkSoftDelete soft deletes multiple organizations
func (r *OrganizationRepository) BulkSoftDelete(ids []int) (int64, error) {
	now := time.Now()
	result := r.db.Model(&model.Organization{}).
		Scopes(scope.WithNotDeleted()).
		Where("id IN ?", ids).
		Update("deleted_at", now)
	return result.RowsAffected, result.Error
}

// ExistsByName checks if organization name exists (excluding soft deleted)
func (r *OrganizationRepository) ExistsByName(name string, excludeID ...int) (bool, error) {
	var count int64
	query := r.db.Model(&model.Organization{}).
		Scopes(scope.WithNotDeleted()).
		Where("name = ?", name)
	if len(excludeID) > 0 {
		query = query.Where("id != ?", excludeID[0])
	}
	err := query.Count(&count).Error
	return count > 0, err
}

// Exists checks if organization exists by ID (excluding soft deleted)
func (r *OrganizationRepository) Exists(id int) (bool, error) {
	var count int64
	err := r.db.Model(&model.Organization{}).
		Scopes(scope.WithNotDeleted()).
		Where("id = ?", id).
		Count(&count).Error
	return count > 0, err
}

// FindTargets finds targets belonging to an organization with pagination
func (r *OrganizationRepository) FindTargets(organizationID int, page, pageSize int, targetType, filter string) ([]model.Target, int64, error) {
	var targets []model.Target
	var total int64

	// Define filter mapping for target
	targetFilterMapping := scope.FilterMapping{
		"name": {Column: "target.name", IsArray: false},
	}

	// Base query: join organization_target to filter by organization
	query := r.db.Model(&model.Target{}).
		Joins("INNER JOIN organization_target ON organization_target.target_id = target.id").
		Where("organization_target.organization_id = ? AND target.deleted_at IS NULL", organizationID)

	if targetType != "" {
		query = query.Where("target.type = ?", targetType)
	}

	query = query.Scopes(scope.WithFilterDefault(filter, targetFilterMapping, "name"))

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Scopes(
		scope.WithPagination(page, pageSize),
		scope.OrderBy("target.created_at", true),
	).Find(&targets).Error

	return targets, total, err
}

// BulkAddTargets adds multiple targets to an organization (ignore duplicates)
// Returns ErrTargetNotFound if any target ID does not exist
func (r *OrganizationRepository) BulkAddTargets(organizationID int, targetIDs []int) error {
	if len(targetIDs) == 0 {
		return nil
	}

	values := make([]any, 0, len(targetIDs)*2)
	placeholders := make([]string, 0, len(targetIDs))

	for _, targetID := range targetIDs {
		placeholders = append(placeholders, "(?, ?)")
		values = append(values, organizationID, targetID)
	}

	sql := "INSERT INTO organization_target (organization_id, target_id) VALUES " +
		strings.Join(placeholders, ", ") +
		" ON CONFLICT DO NOTHING"

	err := r.db.Exec(sql, values...).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgForeignKeyViolation {
			return ErrTargetNotFound
		}
		return err
	}
	return nil
}

// UnlinkTargets removes targets from an organization
func (r *OrganizationRepository) UnlinkTargets(organizationID int, targetIDs []int) (int64, error) {
	if len(targetIDs) == 0 {
		return 0, nil
	}

	result := r.db.Exec(
		"DELETE FROM organization_target WHERE organization_id = ? AND target_id IN ?",
		organizationID, targetIDs,
	)
	return result.RowsAffected, result.Error
}
