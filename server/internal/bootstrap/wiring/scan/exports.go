package scanwiring

import (
	catalogrepo "github.com/yyhuni/lunafox/server/internal/modules/catalog/repository"
	scandomain "github.com/yyhuni/lunafox/server/internal/modules/scan/domain"
	scanrepo "github.com/yyhuni/lunafox/server/internal/modules/scan/repository"
)

func NewStoreAdapter(repo *scanrepo.ScanRepository) *scanStoreAdapter {
	return newScanStoreAdapter(repo)
}

func NewCommandStore(repo *scanrepo.ScanRepository) scandomain.ScanRepository {
	return newScanCommandStore(repo)
}

func NewTaskCancellerAdapter(repo scanrepo.ScanTaskRepository) *scanTaskCancellerAdapter {
	return newScanTaskCancellerAdapter(repo)
}

func NewCreateTargetLookupAdapter(repo *catalogrepo.TargetRepository) *scanCreateTargetLookupAdapter {
	return newScanCreateTargetLookupAdapter(repo)
}

func NewTaskStoreAdapter(repo scanrepo.ScanTaskRepository) *scanTaskStoreAdapter {
	return newScanTaskStoreAdapter(repo)
}

func NewTaskRuntimeScanStoreAdapter(repo *scanrepo.ScanRepository) *taskRuntimeScanStoreAdapter {
	return newTaskRuntimeScanStoreAdapter(repo)
}
