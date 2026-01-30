package model

import (
	"fmt"
	"time"
)

// ScanTask represents a task in the queue supporting priority scheduling
type ScanTask struct {
	ID           int    `gorm:"primaryKey;autoIncrement" json:"id"`
	ScanID       int    `gorm:"not null;index:idx_scan_task_scan_id" json:"scan_id"`
	Stage        int    `gorm:"not null;default:0;index:idx_scan_task_pending_order,priority:2" json:"stage"`
	WorkflowName string `gorm:"type:varchar(100);not null" json:"workflow_name"`
	Status       string `gorm:"type:varchar(20);default:'pending';index:idx_scan_task_pending_order,priority:1" json:"status"`

	// Assignment information
	AgentID      *int   `gorm:"index:idx_scan_task_agent_id" json:"agent_id,omitempty"`
	Config       string `gorm:"type:text" json:"config"`
	ErrorMessage string `gorm:"type:varchar(4096)" json:"error_message,omitempty"`

	// Timestamps
	CreatedAt   time.Time  `gorm:"type:timestamp;default:now();index:idx_scan_task_pending_order,priority:3" json:"created_at"`
	StartedAt   *time.Time `gorm:"type:timestamp" json:"started_at,omitempty"`
	CompletedAt *time.Time `gorm:"type:timestamp" json:"completed_at,omitempty"`
}

// TableName specifies the table name for ScanTask model
func (ScanTask) TableName() string {
	return "scan_task"
}

// WorkspaceDir returns the workspace path for this scan task.
func (t *ScanTask) WorkspaceDir() string {
	return fmt.Sprintf("/opt/lunafox/results/scan_%d/task_%d", t.ScanID, t.ID)
}
