package model

import (
	"time"
)

// Directory represents a directory asset
type Directory struct {
	ID            int       `gorm:"primaryKey;autoIncrement" json:"id"`
	TargetID      int       `gorm:"column:target_id;not null;index:idx_directory_target;uniqueIndex:unique_directory_url_target,priority:1" json:"targetId"`
	URL           string    `gorm:"column:url;size:2000;not null;index:idx_directory_url;uniqueIndex:unique_directory_url_target,priority:2" json:"url"`
	Status        *int      `gorm:"column:status;index:idx_directory_status" json:"status"`
	ContentLength *int64    `gorm:"column:content_length" json:"contentLength"`
	Words         *int      `gorm:"column:words" json:"words"`
	Lines         *int      `gorm:"column:lines" json:"lines"`
	ContentType   string    `gorm:"column:content_type;size:200" json:"contentType"`
	Duration      *int64    `gorm:"column:duration" json:"duration"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime;index:idx_directory_created_at" json:"createdAt"`

	// Relationships
	Target *Target `gorm:"foreignKey:TargetID" json:"target,omitempty"`
}

// TableName returns the table name for Directory
func (Directory) TableName() string {
	return "directory"
}
