package model

import (
	"time"
)

// Organization represents an organization in the system
type Organization struct {
	ID          int        `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string     `gorm:"column:name;size:300" json:"name"`
	Description string     `gorm:"column:description;size:1000" json:"description"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	DeletedAt   *time.Time `gorm:"column:deleted_at;index" json:"-"`
}

// TableName returns the table name for Organization
func (Organization) TableName() string {
	return "organization"
}
