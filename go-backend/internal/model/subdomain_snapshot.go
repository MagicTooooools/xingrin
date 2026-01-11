package model

import (
	"time"
)

// SubdomainSnapshot represents a subdomain discovered in a scan
type SubdomainSnapshot struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	ScanID    int       `gorm:"column:scan_id;not null" json:"scanId"`
	Name      string    `gorm:"column:name;size:1000" json:"name"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`

	// Relationships
	Scan *Scan `gorm:"foreignKey:ScanID" json:"scan,omitempty"`
}

// TableName returns the table name for SubdomainSnapshot
func (SubdomainSnapshot) TableName() string {
	return "subdomain_snapshot"
}
