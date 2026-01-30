package repository

import (
	"github.com/yyhuni/lunafox/server/internal/model"
	"gorm.io/gorm"
)

type SubfinderProviderSettingsRepository struct {
	db *gorm.DB
}

func NewSubfinderProviderSettingsRepository(db *gorm.DB) *SubfinderProviderSettingsRepository {
	return &SubfinderProviderSettingsRepository{db: db}
}

// GetInstance returns the singleton settings (id=1)
func (r *SubfinderProviderSettingsRepository) GetInstance() (*model.SubfinderProviderSettings, error) {
	var settings model.SubfinderProviderSettings
	if err := r.db.First(&settings, 1).Error; err != nil {
		return nil, err
	}
	return &settings, nil
}

// Update updates the settings
func (r *SubfinderProviderSettingsRepository) Update(settings *model.SubfinderProviderSettings) error {
	settings.ID = 1 // Force singleton
	return r.db.Save(settings).Error
}
