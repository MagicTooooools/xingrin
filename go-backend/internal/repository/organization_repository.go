package repository

import (
	"time"

	"github.com/xingrin/go-backend/internal/model"
	"gorm.io/gorm"
)

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
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&org).Error
	if err != nil {
		return nil, err
	}
	return &org, nil
}

// FindAll finds all organizations with pagination (excluding soft deleted)
func (r *OrganizationRepository) FindAll(offset, limit int) ([]model.Organization, int64, error) {
	var orgs []model.Organization
	var total int64

	query := r.db.Model(&model.Organization{}).Where("deleted_at IS NULL")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&orgs).Error
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

// ExistsByName checks if organization name exists (excluding soft deleted)
func (r *OrganizationRepository) ExistsByName(name string, excludeID ...int) (bool, error) {
	var count int64
	query := r.db.Model(&model.Organization{}).Where("name = ? AND deleted_at IS NULL", name)
	if len(excludeID) > 0 {
		query = query.Where("id != ?", excludeID[0])
	}
	err := query.Count(&count).Error
	return count > 0, err
}
