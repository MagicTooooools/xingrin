package repository

import (
	"github.com/yyhuni/lunafox/server/internal/modules/asset/repository/persistence"
	"github.com/yyhuni/lunafox/server/internal/pkg/scope"
)

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
