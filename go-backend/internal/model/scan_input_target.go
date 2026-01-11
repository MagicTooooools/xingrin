package model

import (
	"time"
)

// ScanInputTarget represents a scan input target entry
type ScanInputTarget struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	ScanID    int       `gorm:"column:scan_id;not null;index" json:"scanId"`
	Value     string    `gorm:"column:value;size:2000" json:"value"`
	InputType string    `gorm:"column:input_type;size:10" json:"inputType"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`

	// Relationships
	Scan *Scan `gorm:"foreignKey:ScanID" json:"scan,omitempty"`
}

// TableName returns the table name for ScanInputTarget
func (ScanInputTarget) TableName() string {
	return "scan_input_target"
}

// InputType constants
const (
	InputTypeDomain = "domain"
	InputTypeIP     = "ip"
	InputTypeCIDR   = "cidr"
	InputTypeURL    = "url"
)
