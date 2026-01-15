package model

import (
	"time"

	"gorm.io/datatypes"
)

// NotificationSettings represents notification settings (singleton)
type NotificationSettings struct {
	ID                int            `gorm:"primaryKey" json:"id"`
	DiscordEnabled    bool           `gorm:"column:discord_enabled;default:false" json:"discordEnabled"`
	DiscordWebhookURL string         `gorm:"column:discord_webhook_url;size:500" json:"discordWebhookUrl"`
	WecomEnabled      bool           `gorm:"column:wecom_enabled;default:false" json:"wecomEnabled"`
	WecomWebhookURL   string         `gorm:"column:wecom_webhook_url;size:500" json:"wecomWebhookUrl"`
	Categories        datatypes.JSON `gorm:"column:categories;type:jsonb" json:"categories"`
	CreatedAt         time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt         time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

// TableName returns the table name for NotificationSettings
func (NotificationSettings) TableName() string {
	return "notification_settings"
}

// DefaultCategories returns the default category configuration
func DefaultCategories() map[string]bool {
	return map[string]bool{
		"scan":          true,
		"vulnerability": true,
		"asset":         true,
		"system":        false,
	}
}
