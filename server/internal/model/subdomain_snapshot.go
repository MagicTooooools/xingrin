package model

import (
	"time"
)

// SubdomainSnapshot represents a subdomain snapshot
type SubdomainSnapshot struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	ScanID    int       `gorm:"column:scan_id;not null;index:idx_subdomain_snap_scan;uniqueIndex:unique_subdomain_per_scan_snapshot,priority:1" json:"scanId"`
	Name      string    `gorm:"column:name;size:1000;index:idx_subdomain_snap_name;uniqueIndex:unique_subdomain_per_scan_snapshot,priority:2" json:"name"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime;index:idx_subdomain_snap_created_at" json:"createdAt"`

	// Relationships
	Scan *Scan `gorm:"foreignKey:ScanID" json:"scan,omitempty"`
}

// TableName returns the table name for SubdomainSnapshot
func (SubdomainSnapshot) TableName() string {
	return "subdomain_snapshot"
}
