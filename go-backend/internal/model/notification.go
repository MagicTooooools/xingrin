package model

import (
	"time"
)

// Notification represents a notification entry
type Notification struct {
	ID        int        `gorm:"primaryKey;autoIncrement" json:"id"`
	Category  string     `gorm:"column:category;size:20;index:idx_notification_category" json:"category"`
	Level     string     `gorm:"column:level;size:20;index:idx_notification_level" json:"level"`
	Title     string     `gorm:"column:title;size:200" json:"title"`
	Message   string     `gorm:"column:message;size:2000" json:"message"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime;index:idx_notification_created_at" json:"createdAt"`
	IsRead    bool       `gorm:"column:is_read;default:false;index:idx_notification_is_read" json:"isRead"`
	ReadAt    *time.Time `gorm:"column:read_at" json:"readAt"`
}

// TableName returns the table name for Notification
func (Notification) TableName() string {
	return "notification"
}

// NotificationCategory constants
const (
	NotificationCategoryScan          = "scan"
	NotificationCategoryVulnerability = "vulnerability"
	NotificationCategoryAsset         = "asset"
	NotificationCategorySystem        = "system"
)

// NotificationLevel constants
const (
	NotificationLevelLow      = "low"
	NotificationLevelMedium   = "medium"
	NotificationLevelHigh     = "high"
	NotificationLevelCritical = "critical"
)
