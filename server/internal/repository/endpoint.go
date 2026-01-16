package repository

import (
	"database/sql"

	"github.com/orbit/server/internal/model"
	"github.com/orbit/server/internal/pkg/scope"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// EndpointRepository handles endpoint database operations
type EndpointRepository struct {
	db *gorm.DB
}

// NewEndpointRepository creates a new endpoint repository
func NewEndpointRepository(db *gorm.DB) *EndpointRepository {
	return &EndpointRepository{db: db}
}

// EndpointFilterMapping defines field mapping for endpoint filtering
var EndpointFilterMapping = scope.FilterMapping{
	"url":    {Column: "url"},
	"host":   {Column: "host"},
	"title":  {Column: "title"},
	"status": {Column: "status_code", IsNumeric: true},
	"tech":   {Column: "tech", IsArray: true},
}

// FindByTargetID finds endpoints by target ID with pagination and filter
func (r *EndpointRepository) FindByTargetID(targetID int, page, pageSize int, filter string) ([]model.Endpoint, int64, error) {
	var endpoints []model.Endpoint
	var total int64

	// Base query
	baseQuery := r.db.Model(&model.Endpoint{}).Where("target_id = ?", targetID)

	// Apply filter scope with default field "url"
	baseQuery = baseQuery.Scopes(scope.WithFilterDefault(filter, EndpointFilterMapping, "url"))

	// Count total
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch with pagination and ordering
	err := baseQuery.Scopes(
		scope.WithPagination(page, pageSize),
		scope.OrderByCreatedAtDesc(),
	).Find(&endpoints).Error

	return endpoints, total, err
}

// FindByID finds an endpoint by ID
func (r *EndpointRepository) FindByID(id int) (*model.Endpoint, error) {
	var endpoint model.Endpoint
	err := r.db.First(&endpoint, id).Error
	if err != nil {
		return nil, err
	}
	return &endpoint, nil
}

// BulkCreate creates multiple endpoints, ignoring duplicates
func (r *EndpointRepository) BulkCreate(endpoints []model.Endpoint) (int, error) {
	if len(endpoints) == 0 {
		return 0, nil
	}

	// Use ON CONFLICT DO NOTHING to ignore duplicates (url + target_id unique)
	result := r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&endpoints)
	if result.Error != nil {
		return 0, result.Error
	}

	return int(result.RowsAffected), nil
}

// Delete deletes an endpoint by ID
func (r *EndpointRepository) Delete(id int) error {
	return r.db.Delete(&model.Endpoint{}, id).Error
}

// BulkDelete deletes multiple endpoints by IDs
func (r *EndpointRepository) BulkDelete(ids []int) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	result := r.db.Where("id IN ?", ids).Delete(&model.Endpoint{})
	return result.RowsAffected, result.Error
}

// StreamByTargetID returns a sql.Rows cursor for streaming export
func (r *EndpointRepository) StreamByTargetID(targetID int) (*sql.Rows, error) {
	return r.db.Model(&model.Endpoint{}).
		Where("target_id = ?", targetID).
		Order("created_at DESC").
		Rows()
}

// CountByTargetID returns the count of endpoints for a target
func (r *EndpointRepository) CountByTargetID(targetID int) (int64, error) {
	var count int64
	err := r.db.Model(&model.Endpoint{}).Where("target_id = ?", targetID).Count(&count).Error
	return count, err
}

// ScanRow scans a single row into Endpoint model
func (r *EndpointRepository) ScanRow(rows *sql.Rows) (*model.Endpoint, error) {
	var endpoint model.Endpoint
	if err := r.db.ScanRows(rows, &endpoint); err != nil {
		return nil, err
	}
	return &endpoint, nil
}

// BulkUpsert creates or updates multiple endpoints
// Uses ON CONFLICT DO UPDATE with COALESCE for non-null updates
// Tech and MatchedGFPatterns arrays are merged and deduplicated
func (r *EndpointRepository) BulkUpsert(endpoints []model.Endpoint) (int64, error) {
	if len(endpoints) == 0 {
		return 0, nil
	}

	var totalAffected int64

	// Process in batches to avoid parameter limits
	batchSize := 100
	for i := 0; i < len(endpoints); i += batchSize {
		end := i + batchSize
		if end > len(endpoints) {
			end = len(endpoints)
		}
		batch := endpoints[i:end]

		affected, err := r.upsertBatch(batch)
		if err != nil {
			return totalAffected, err
		}
		totalAffected += affected
	}

	return totalAffected, nil
}

// upsertBatch upserts a single batch of endpoints
func (r *EndpointRepository) upsertBatch(endpoints []model.Endpoint) (int64, error) {
	if len(endpoints) == 0 {
		return 0, nil
	}

	result := r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "url"}, {Name: "target_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"host":             gorm.Expr("COALESCE(NULLIF(EXCLUDED.host, ''), endpoint.host)"),
			"location":         gorm.Expr("COALESCE(NULLIF(EXCLUDED.location, ''), endpoint.location)"),
			"title":            gorm.Expr("COALESCE(NULLIF(EXCLUDED.title, ''), endpoint.title)"),
			"webserver":        gorm.Expr("COALESCE(NULLIF(EXCLUDED.webserver, ''), endpoint.webserver)"),
			"response_body":    gorm.Expr("COALESCE(NULLIF(EXCLUDED.response_body, ''), endpoint.response_body)"),
			"content_type":     gorm.Expr("COALESCE(NULLIF(EXCLUDED.content_type, ''), endpoint.content_type)"),
			"status_code":      gorm.Expr("COALESCE(EXCLUDED.status_code, endpoint.status_code)"),
			"content_length":   gorm.Expr("COALESCE(EXCLUDED.content_length, endpoint.content_length)"),
			"vhost":            gorm.Expr("COALESCE(EXCLUDED.vhost, endpoint.vhost)"),
			"response_headers": gorm.Expr("COALESCE(NULLIF(EXCLUDED.response_headers, ''), endpoint.response_headers)"),
			// Merge tech arrays and deduplicate
			"tech": gorm.Expr(`(
				SELECT ARRAY(SELECT DISTINCT unnest FROM unnest(
					COALESCE(endpoint.tech, ARRAY[]::varchar(100)[]) ||
					COALESCE(EXCLUDED.tech, ARRAY[]::varchar(100)[])
				) ORDER BY unnest)
			)`),
			// Merge matched_gf_patterns arrays and deduplicate
			"matched_gf_patterns": gorm.Expr(`(
				SELECT ARRAY(SELECT DISTINCT unnest FROM unnest(
					COALESCE(endpoint.matched_gf_patterns, ARRAY[]::varchar(100)[]) ||
					COALESCE(EXCLUDED.matched_gf_patterns, ARRAY[]::varchar(100)[])
				) ORDER BY unnest)
			)`),
		}),
	}).Create(&endpoints)

	return result.RowsAffected, result.Error
}
