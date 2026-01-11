package model

import (
	"time"
)

// DirectorySnapshot represents a directory discovered in a scan
type DirectorySnapshot struct {
	ID            int       `gorm:"primaryKey;autoIncrement" json:"id"`
	ScanID        int       `gorm:"column:scan_id;not null" json:"scanId"`
	URL           string    `gorm:"column:url;size:2000" json:"url"`
	Status        *int      `gorm:"column:status" json:"status"`
	ContentLength *int64    `gorm:"column:content_length" json:"contentLength"`
	Words         *int      `gorm:"column:words" json:"words"`
	Lines         *int      `gorm:"column:lines" json:"lines"`
	ContentType   string    `gorm:"column:content_type;size:200" json:"contentType"`
	Duration      *int64    `gorm:"column:duration" json:"duration"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`

	// Relationships
	Scan *Scan `gorm:"foreignKey:ScanID" json:"scan,omitempty"`
}

// TableName returns the table name for DirectorySnapshot
func (DirectorySnapshot) TableName() string {
	return "directory_snapshot"
}
