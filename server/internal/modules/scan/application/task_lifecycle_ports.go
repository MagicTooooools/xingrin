package application

import "context"

type ScanTaskCanceller interface {
	CancelTasksByScanID(ctx context.Context, scanID int) ([]CancelledTaskInfo, error)
}

type TaskCancelNotifier interface {
	SendTaskCancel(agentID, taskID int)
}
