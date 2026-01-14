package model

import (
	"time"
)

// DirectorySnapshot represents a directory snapshot
type DirectorySnapshot struct {
	ID            int       `gorm:"primaryKey;autoIncrement" json:"id"`
	ScanID        int       `gorm:"column:scan_id;not null;index:idx_directory_snap_scan;uniqueIndex:unique_directory_per_scan_snapshot,priority:1" json:"scanId"`
	URL           string    `gorm:"column:url;size:2000;index:idx_directory_snap_url;uniqueIndex:unique_directory_per_scan_snapshot,priority:2" json:"url"`
	Status        *int      `gorm:"column:status;index:idx_directory_snap_status" json:"status"`
	ContentLength *int      `gorm:"column:content_length" json:"contentLength"`
	ContentType   string    `gorm:"column:content_type;size:200;index:idx_directory_snap_content_type" json:"contentType"`
	Duration      *int      `gorm:"column:duration" json:"duration"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime;index:idx_directory_snap_created_at" json:"createdAt"`

	// Relationships
	Scan *Scan `gorm:"foreignKey:ScanID" json:"scan,omitempty"`
}

// TableName returns the table name for DirectorySnapshot
func (DirectorySnapshot) TableName() string {
	return "directory_snapshot"
}
