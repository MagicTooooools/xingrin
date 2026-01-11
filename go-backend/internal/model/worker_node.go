package model

import (
	"time"
)

// WorkerNode represents a worker node
type WorkerNode struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"column:name;size:100;uniqueIndex:unique_worker_name" json:"name"`
	IPAddress string    `gorm:"column:ip_address;type:inet" json:"ipAddress"`
	SSHPort   int       `gorm:"column:ssh_port;default:22" json:"sshPort"`
	Username  string    `gorm:"column:username;size:50;default:'root'" json:"username"`
	Password  string    `gorm:"column:password;size:200" json:"-"`
	IsLocal   bool      `gorm:"column:is_local;default:false" json:"isLocal"`
	Status    string    `gorm:"column:status;size:20;default:'pending'" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime;index:idx_worker_node_created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

// TableName returns the table name for WorkerNode
func (WorkerNode) TableName() string {
	return "worker_node"
}

// WorkerStatus constants
const (
	WorkerStatusPending   = "pending"
	WorkerStatusDeploying = "deploying"
	WorkerStatusOnline    = "online"
	WorkerStatusOffline   = "offline"
	WorkerStatusUpdating  = "updating"
	WorkerStatusOutdated  = "outdated"
)
