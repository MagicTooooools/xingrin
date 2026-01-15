package repository

import (
	"database/sql"
	"net/url"

	"github.com/xingrin/server/internal/model"
	"github.com/xingrin/server/internal/pkg/scope"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// WebsiteRepository handles website database operations
type WebsiteRepository struct {
	db *gorm.DB
}

// NewWebsiteRepository creates a new website repository
func NewWebsiteRepository(db *gorm.DB) *WebsiteRepository {
	return &WebsiteRepository{db: db}
}

// WebsiteFilterMapping defines field mapping for website filtering
var WebsiteFilterMapping = scope.FilterMapping{
	"url":    {Column: "url"},
	"host":   {Column: "host"},
	"title":  {Column: "title"},
	"status": {Column: "status_code", IsNumeric: true},
	"tech":   {Column: "tech", IsArray: true},
}

// FindByTargetID finds websites by target ID with pagination and filter
func (r *WebsiteRepository) FindByTargetID(targetID int, page, pageSize int, filter string) ([]model.Website, int64, error) {
	var websites []model.Website
	var total int64

	// Base query
	baseQuery := r.db.Model(&model.Website{}).Where("target_id = ?", targetID)

	// Apply filter scope
	baseQuery = baseQuery.Scopes(scope.WithFilter(filter, WebsiteFilterMapping))

	// Count total
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch with pagination and ordering
	err := baseQuery.Scopes(
		scope.WithPagination(page, pageSize),
		scope.OrderByCreatedAtDesc(),
	).Find(&websites).Error

	return websites, total, err
}

// FindByID finds a website by ID
func (r *WebsiteRepository) FindByID(id int) (*model.Website, error) {
	var website model.Website
	err := r.db.First(&website, id).Error
	if err != nil {
		return nil, err
	}
	return &website, nil
}

// BulkCreate creates multiple websites, ignoring duplicates
func (r *WebsiteRepository) BulkCreate(websites []model.Website) (int, error) {
	if len(websites) == 0 {
		return 0, nil
	}

	// Use ON CONFLICT DO NOTHING to ignore duplicates (url + target_id unique)
	result := r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&websites)
	if result.Error != nil {
		return 0, result.Error
	}

	return int(result.RowsAffected), nil
}

// Delete deletes a website by ID
func (r *WebsiteRepository) Delete(id int) error {
	return r.db.Delete(&model.Website{}, id).Error
}

// BulkDelete deletes multiple websites by IDs
func (r *WebsiteRepository) BulkDelete(ids []int) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	result := r.db.Where("id IN ?", ids).Delete(&model.Website{})
	return result.RowsAffected, result.Error
}

// ExtractHostFromURL extracts host from URL
func ExtractHostFromURL(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	return parsed.Host
}

// StreamByTargetID returns a sql.Rows cursor for streaming export
func (r *WebsiteRepository) StreamByTargetID(targetID int) (*sql.Rows, error) {
	return r.db.Model(&model.Website{}).
		Where("target_id = ?", targetID).
		Order("created_at DESC").
		Rows()
}

// CountByTargetID returns the count of websites for a target
func (r *WebsiteRepository) CountByTargetID(targetID int) (int64, error) {
	var count int64
	err := r.db.Model(&model.Website{}).Where("target_id = ?", targetID).Count(&count).Error
	return count, err
}

// ScanRow scans a single row into Website model
func (r *WebsiteRepository) ScanRow(rows *sql.Rows) (*model.Website, error) {
	var website model.Website
	if err := r.db.ScanRows(rows, &website); err != nil {
		return nil, err
	}
	return &website, nil
}

// BulkUpsert creates or updates multiple websites
// Uses ON CONFLICT DO UPDATE with COALESCE for non-null updates
// Tech array is merged and deduplicated
func (r *WebsiteRepository) BulkUpsert(websites []model.Website) (int64, error) {
	if len(websites) == 0 {
		return 0, nil
	}

	var totalAffected int64

	// Process in batches to avoid parameter limits
	batchSize := 100
	for i := 0; i < len(websites); i += batchSize {
		end := i + batchSize
		if end > len(websites) {
			end = len(websites)
		}
		batch := websites[i:end]

		affected, err := r.upsertBatch(batch)
		if err != nil {
			return totalAffected, err
		}
		totalAffected += affected
	}

	return totalAffected, nil
}

// upsertBatch upserts a single batch of websites
func (r *WebsiteRepository) upsertBatch(websites []model.Website) (int64, error) {
	if len(websites) == 0 {
		return 0, nil
	}

	// Use GORM's OnConflict with custom UpdateAll
	// For tech array merge, we need raw SQL in the update clause
	result := r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "url"}, {Name: "target_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"host":             gorm.Expr("COALESCE(NULLIF(EXCLUDED.host, ''), website.host)"),
			"location":         gorm.Expr("COALESCE(NULLIF(EXCLUDED.location, ''), website.location)"),
			"title":            gorm.Expr("COALESCE(NULLIF(EXCLUDED.title, ''), website.title)"),
			"webserver":        gorm.Expr("COALESCE(NULLIF(EXCLUDED.webserver, ''), website.webserver)"),
			"response_body":    gorm.Expr("COALESCE(NULLIF(EXCLUDED.response_body, ''), website.response_body)"),
			"content_type":     gorm.Expr("COALESCE(NULLIF(EXCLUDED.content_type, ''), website.content_type)"),
			"status_code":      gorm.Expr("COALESCE(EXCLUDED.status_code, website.status_code)"),
			"content_length":   gorm.Expr("COALESCE(EXCLUDED.content_length, website.content_length)"),
			"vhost":            gorm.Expr("COALESCE(EXCLUDED.vhost, website.vhost)"),
			"response_headers": gorm.Expr("COALESCE(NULLIF(EXCLUDED.response_headers, ''), website.response_headers)"),
			// Merge tech arrays and deduplicate
			"tech": gorm.Expr(`(
				SELECT ARRAY(SELECT DISTINCT unnest FROM unnest(
					COALESCE(website.tech, ARRAY[]::varchar(100)[]) ||
					COALESCE(EXCLUDED.tech, ARRAY[]::varchar(100)[])
				) ORDER BY unnest)
			)`),
		}),
	}).Create(&websites)

	return result.RowsAffected, result.Error
}
