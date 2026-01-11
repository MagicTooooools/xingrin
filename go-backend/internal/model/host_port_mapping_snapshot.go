package model

import (
	"time"
)

// HostPortMappingSnapshot represents a host-IP-port mapping discovered in a scan
type HostPortMappingSnapshot struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	ScanID    int       `gorm:"column:scan_id;not null" json:"scanId"`
	Host      string    `gorm:"column:host;size:1000;not null" json:"host"`
	IP        string    `gorm:"column:ip;type:inet;not null" json:"ip"`
	Port      int       `gorm:"column:port;not null" json:"port"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`

	// Relationships
	Scan *Scan `gorm:"foreignKey:ScanID" json:"scan,omitempty"`
}

// TableName returns the table name for HostPortMappingSnapshot
func (HostPortMappingSnapshot) TableName() string {
	return "host_port_mapping_snapshot"
}
