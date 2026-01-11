package model

import (
	"time"
)

// ScreenshotSnapshot represents a screenshot captured in a scan
type ScreenshotSnapshot struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`
	ScanID     int       `gorm:"column:scan_id;not null" json:"scanId"`
	URL        string    `gorm:"column:url;type:text" json:"url"`
	StatusCode *int16    `gorm:"column:status_code" json:"statusCode"`
	Image      []byte    `gorm:"column:image;type:bytea" json:"-"` // Hidden from JSON
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`

	// Relationships
	Scan *Scan `gorm:"foreignKey:ScanID" json:"scan,omitempty"`
}

// TableName returns the table name for ScreenshotSnapshot
func (ScreenshotSnapshot) TableName() string {
	return "screenshot_snapshot"
}
