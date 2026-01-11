package model

import (
	"time"

	"gorm.io/datatypes"
)

// SubfinderProviderSettings represents subfinder provider settings (singleton)
type SubfinderProviderSettings struct {
	ID        int            `gorm:"primaryKey" json:"id"`
	Providers datatypes.JSON `gorm:"column:providers;type:jsonb" json:"providers"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

// TableName returns the table name for SubfinderProviderSettings
func (SubfinderProviderSettings) TableName() string {
	return "subfinder_provider_settings"
}

// DefaultProviders returns the default provider configuration
func DefaultProviders() map[string]interface{} {
	return map[string]interface{}{
		"fofa":           map[string]interface{}{"enabled": false, "email": "", "api_key": ""},
		"hunter":         map[string]interface{}{"enabled": false, "api_key": ""},
		"shodan":         map[string]interface{}{"enabled": false, "api_key": ""},
		"censys":         map[string]interface{}{"enabled": false, "api_id": "", "api_secret": ""},
		"zoomeye":        map[string]interface{}{"enabled": false, "api_key": ""},
		"securitytrails": map[string]interface{}{"enabled": false, "api_key": ""},
		"threatbook":     map[string]interface{}{"enabled": false, "api_key": ""},
		"quake":          map[string]interface{}{"enabled": false, "api_key": ""},
	}
}
