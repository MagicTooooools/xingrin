package model

import (
	"time"
)

// Screenshot represents a screenshot asset
type Screenshot struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`
	TargetID   int       `gorm:"column:target_id;not null;index:idx_screenshot_target;uniqueIndex:unique_screenshot_per_target,priority:1" json:"targetId"`
	URL        string    `gorm:"column:url;type:text;uniqueIndex:unique_screenshot_per_target,priority:2" json:"url"`
	StatusCode *int16    `gorm:"column:status_code" json:"statusCode"`
	Image      []byte    `gorm:"column:image;type:bytea" json:"-"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime;index:idx_screenshot_created_at" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`

	// Relationships
	Target *Target `gorm:"foreignKey:TargetID" json:"target,omitempty"`
}

// TableName returns the table name for Screenshot
func (Screenshot) TableName() string {
	return "screenshot"
}
