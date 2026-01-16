package repository

import (
	"github.com/orbit/server/internal/model"
	"github.com/orbit/server/internal/pkg/scope"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ScreenshotRepository handles screenshot database operations
type ScreenshotRepository struct {
	db *gorm.DB
}

// NewScreenshotRepository creates a new screenshot repository
func NewScreenshotRepository(db *gorm.DB) *ScreenshotRepository {
	return &ScreenshotRepository{db: db}
}

// FindByTargetID finds screenshots by target ID with pagination and filter
func (r *ScreenshotRepository) FindByTargetID(targetID int, page, pageSize int, filter string) ([]model.Screenshot, int64, error) {
	var screenshots []model.Screenshot
	var total int64

	// Build base query
	baseQuery := r.db.Model(&model.Screenshot{}).Where("target_id = ?", targetID)

	// Apply URL filter (fuzzy search)
	if filter != "" {
		baseQuery = baseQuery.Where("url ILIKE ?", "%"+filter+"%")
	}

	// Count total
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch with pagination
	err := baseQuery.
		Scopes(
			scope.WithPagination(page, pageSize),
			scope.OrderByCreatedAtDesc(),
		).
		Find(&screenshots).Error

	return screenshots, total, err
}

// FindByID finds a screenshot by ID
func (r *ScreenshotRepository) FindByID(id int) (*model.Screenshot, error) {
	var screenshot model.Screenshot
	err := r.db.Where("id = ?", id).First(&screenshot).Error
	if err != nil {
		return nil, err
	}
	return &screenshot, nil
}

// BulkDelete deletes multiple screenshots by IDs
func (r *ScreenshotRepository) BulkDelete(ids []int) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	result := r.db.Where("id IN ?", ids).Delete(&model.Screenshot{})
	return result.RowsAffected, result.Error
}

// BulkUpsert creates or updates multiple screenshots
// Uses ON CONFLICT DO UPDATE with COALESCE for non-null updates
func (r *ScreenshotRepository) BulkUpsert(screenshots []model.Screenshot) (int64, error) {
	if len(screenshots) == 0 {
		return 0, nil
	}

	var totalAffected int64

	// Process in batches to avoid parameter limits (5 fields per record)
	batchSize := 500
	for i := 0; i < len(screenshots); i += batchSize {
		end := i + batchSize
		if end > len(screenshots) {
			end = len(screenshots)
		}
		batch := screenshots[i:end]

		result := r.db.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "target_id"}, {Name: "url"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"status_code": gorm.Expr("COALESCE(EXCLUDED.status_code, screenshot.status_code)"),
				"image":       gorm.Expr("COALESCE(EXCLUDED.image, screenshot.image)"),
				"updated_at":  gorm.Expr("CURRENT_TIMESTAMP"),
			}),
		}).Create(&batch)
		if result.Error != nil {
			return totalAffected, result.Error
		}
		totalAffected += result.RowsAffected
	}

	return totalAffected, nil
}
