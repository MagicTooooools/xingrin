package model

import (
	"time"
)

// ScanEngine represents a scan engine
type ScanEngine struct {
	ID            int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string    `gorm:"column:name;size:200;uniqueIndex:unique_scan_engine_name" json:"name"`
	Configuration string    `gorm:"column:configuration;size:10000" json:"configuration"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime;index:idx_scan_engine_created_at" json:"createdAt"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

// TableName returns the table name for ScanEngine
func (ScanEngine) TableName() string {
	return "scan_engine"
}
