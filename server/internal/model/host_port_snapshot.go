package model

import (
	"time"
)

// HostPortSnapshot represents a host-port snapshot
type HostPortSnapshot struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	ScanID    int       `gorm:"column:scan_id;not null;index:idx_hpm_snap_scan;uniqueIndex:unique_scan_host_ip_port_snapshot,priority:1" json:"scanId"`
	Host      string    `gorm:"column:host;size:1000;not null;index:idx_hpm_snap_host;uniqueIndex:unique_scan_host_ip_port_snapshot,priority:2" json:"host"`
	IP        string    `gorm:"column:ip;type:inet;not null;index:idx_hpm_snap_ip;uniqueIndex:unique_scan_host_ip_port_snapshot,priority:3" json:"ip"`
	Port      int       `gorm:"column:port;not null;index:idx_hpm_snap_port;uniqueIndex:unique_scan_host_ip_port_snapshot,priority:4" json:"port"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime;index:idx_hpm_snap_created_at" json:"createdAt"`

	// Relationships
	Scan *Scan `gorm:"foreignKey:ScanID" json:"scan,omitempty"`
}

// TableName returns the table name for HostPortSnapshot
func (HostPortSnapshot) TableName() string {
	return "host_port_mapping_snapshot"
}
