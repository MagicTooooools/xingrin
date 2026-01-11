package model

import (
	"time"
)

// HostPortMapping represents a host-IP-port mapping
type HostPortMapping struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	TargetID  int       `gorm:"column:target_id;not null" json:"targetId"`
	Host      string    `gorm:"column:host;size:1000;not null" json:"host"`
	IP        string    `gorm:"column:ip;type:inet;not null" json:"ip"`
	Port      int       `gorm:"column:port;not null" json:"port"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`

	// Relationships
	Target *Target `gorm:"foreignKey:TargetID" json:"target,omitempty"`
}

// TableName returns the table name for HostPortMapping
func (HostPortMapping) TableName() string {
	return "host_port_mapping"
}
