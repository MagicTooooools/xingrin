package model

import (
	"time"
)

// ScanLog represents a scan log entry
type ScanLog struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ScanID    int       `gorm:"column:scan_id;not null;index:idx_scan_log_scan" json:"scanId"`
	Level     string    `gorm:"column:level;size:10;default:'info'" json:"level"`
	Content   string    `gorm:"column:content;type:text" json:"content"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime;index:idx_scan_log_created_at" json:"createdAt"`

	// Relationships
	Scan *Scan `gorm:"foreignKey:ScanID" json:"scan,omitempty"`
}

// TableName returns the table name for ScanLog
func (ScanLog) TableName() string {
	return "scan_log"
}

// ScanLogLevel constants
const (
	ScanLogLevelInfo    = "info"
	ScanLogLevelWarning = "warning"
	ScanLogLevelError   = "error"
)
