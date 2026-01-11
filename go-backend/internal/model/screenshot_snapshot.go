package model

import (
	"time"
)

// ScreenshotSnapshot represents a screenshot snapshot
type ScreenshotSnapshot struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`
	ScanID     int       `gorm:"column:scan_id;not null;index:idx_screenshot_snap_scan;uniqueIndex:unique_screenshot_per_scan_snapshot,priority:1" json:"scanId"`
	URL        string    `gorm:"column:url;type:text;uniqueIndex:unique_screenshot_per_scan_snapshot,priority:2" json:"url"`
	StatusCode *int16    `gorm:"column:status_code" json:"statusCode"`
	Image      []byte    `gorm:"column:image;type:bytea" json:"-"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime;index:idx_screenshot_snap_created_at" json:"createdAt"`

	// Relationships
	Scan *Scan `gorm:"foreignKey:ScanID" json:"scan,omitempty"`
}

// TableName returns the table name for ScreenshotSnapshot
func (ScreenshotSnapshot) TableName() string {
	return "screenshot_snapshot"
}
