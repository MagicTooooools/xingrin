package model

import "time"

// Agent represents a persistent agent service running on remote VPS
type Agent struct {
	ID     int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name   string `gorm:"type:varchar(100);not null" json:"name"`
	APIKey string `gorm:"type:varchar(8);not null;uniqueIndex" json:"api_key"`
	Status string `gorm:"type:varchar(20);default:'online'" json:"status"` // online/offline

	// Connection info
	Hostname  string `gorm:"type:varchar(255)" json:"hostname"`
	IPAddress string `gorm:"type:varchar(45)" json:"ip_address"`
	Version   string `gorm:"type:varchar(20)" json:"version"`

	// Scheduling configuration (dynamically modifiable via API)
	MaxTasks      int `gorm:"default:5" json:"max_tasks"`
	CPUThreshold  int `gorm:"default:85" json:"cpu_threshold"`
	MemThreshold  int `gorm:"default:85" json:"mem_threshold"`
	DiskThreshold int `gorm:"default:90" json:"disk_threshold"`

	// Self-registration related
	RegistrationToken string `gorm:"type:varchar(8)" json:"registration_token,omitempty"`

	// Timestamps
	ConnectedAt   *time.Time `gorm:"type:timestamp" json:"connected_at,omitempty"`
	LastHeartbeat *time.Time `gorm:"type:timestamp" json:"last_heartbeat,omitempty"`
	CreatedAt     time.Time  `gorm:"type:timestamp;default:now()" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"type:timestamp;default:now()" json:"updated_at"`
}

// TableName specifies the table name for Agent model
func (Agent) TableName() string {
	return "agent"
}

// RegistrationToken represents a token for agent self-registration
type RegistrationToken struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Token     string    `gorm:"type:varchar(8);not null;uniqueIndex" json:"token"`
	ExpiresAt time.Time `gorm:"type:timestamp;not null;default:now() + interval '1 hour'" json:"expires_at"`
	CreatedAt time.Time `gorm:"type:timestamp;default:now()" json:"created_at"`
}

// TableName specifies the table name for RegistrationToken model
func (RegistrationToken) TableName() string {
	return "registration_token"
}
