package application

import "context"

type ScanQueryStore interface {
	FindAll(page, pageSize int, targetID int, status, search string) ([]QueryScan, int64, error)
	GetQueryByID(id int) (*QueryScan, error)
	GetStatistics() (*QueryStatistics, error)
}

type ScanCommandStore interface {
	GetByIDNotDeleted(id int) (*QueryScan, error)
	FindByIDs(ids []int) ([]QueryScan, error)
	CreateWithInputTargetsAndTasks(scan *CreateScan, inputs []CreateScanInputTarget, tasks []CreateScanTask) error
	BulkSoftDelete(ids []int) (int64, []string, error)
	UpdateStatus(id int, status string, errorMessage ...string) error
}

type ScanCreateCommandStore interface {
	CreateWithInputTargetsAndTasks(scan *CreateScan, inputs []CreateScanInputTarget, tasks []CreateScanTask) error
}

type CreateTargetLookup interface {
	GetCreateTargetRefByID(id int) (*TargetRef, error)
}

type ScanTaskCanceller interface {
	CancelTasksByScanID(ctx context.Context, scanID int) ([]CancelledTaskInfo, error)
}

type TaskCancelNotifier interface {
	SendTaskCancel(agentID, taskID int)
}

type TaskStore interface {
	GetByID(ctx context.Context, id int) (*TaskRecord, error)
	PullTask(ctx context.Context, agentID int) (*TaskRecord, error)
	UpdateStatus(ctx context.Context, id int, status string, errorMessage string) error
	GetStatusCountsByScanID(ctx context.Context, scanID int) (pending, running, completed, failed, cancelled int, err error)
	CountActiveByScanAndStage(ctx context.Context, scanID, stage int) (int, error)
	UnlockNextStage(ctx context.Context, scanID, stage int) (int64, error)
}

type TaskRuntimeScanStore interface {
	GetTaskRuntimeByID(id int) (*TaskScanRecord, error)
	UpdateStatus(id int, status string, errorMessage string) error
}
