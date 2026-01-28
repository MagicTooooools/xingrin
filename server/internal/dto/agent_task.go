package dto

// TaskAssignment represents a task assigned to an agent.
type TaskAssignment struct {
	TaskID       int    `json:"taskId"`
	ScanID       int    `json:"scanId"`
	Stage        int    `json:"stage"`
	WorkflowName string `json:"workflowName"`
	TargetID     int    `json:"targetId"`
	TargetName   string `json:"targetName"`
	TargetType   string `json:"targetType"`
	WorkspaceDir string `json:"workspaceDir"`
	Config       string `json:"config"`
}

// TaskStatusUpdateRequest represents a task status update from an agent.
type TaskStatusUpdateRequest struct {
	Status       string `json:"status" binding:"required"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}
