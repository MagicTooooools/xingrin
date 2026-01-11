package model

import (
	"time"
)

// Screenshot represents a screenshot asset
type Screenshot struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`
	TargetID   int       `gorm:"column:target_id;not null" json:"targetId"`
	URL        string    `gorm:"column:url;type:text" json:"url"`
	StatusCode *int16    `gorm:"column:status_code" json:"statusCode"`
	Image      []byte    `gorm:"column:image;type:bytea" json:"-"` // Hidden from JSON
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`

	// Relationships
	Target *Target `gorm:"foreignKey:TargetID" json:"target,omitempty"`
}

// TableName returns the table name for Screenshot
func (Screenshot) TableName() string {
	return "screenshot"
}
