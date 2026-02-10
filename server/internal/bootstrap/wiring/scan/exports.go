package scanwiring

import (
	catalogrepo "github.com/yyhuni/lunafox/server/internal/modules/catalog/repository"
	scanapp "github.com/yyhuni/lunafox/server/internal/modules/scan/application"
	scandomain "github.com/yyhuni/lunafox/server/internal/modules/scan/domain"
	scanrepo "github.com/yyhuni/lunafox/server/internal/modules/scan/repository"
)

func NewQueryStoreAdapter(repo *scanrepo.ScanRepository) scanapp.ScanQueryStore {
	return newScanStoreAdapter(repo)
}

func NewScanCommandStoreAdapter(repo *scanrepo.ScanRepository) scanapp.ScanCommandStore {
	return newScanStoreAdapter(repo)
}

func NewCommandStore(repo *scanrepo.ScanRepository) scandomain.ScanRepository {
	return newScanCommandStore(repo)
}

func NewTaskCancellerAdapter(repo scanrepo.ScanTaskRepository) scanapp.ScanTaskCanceller {
	return newScanTaskCancellerAdapter(repo)
}

func NewCreateTargetLookupAdapter(repo *catalogrepo.TargetRepository) scanapp.CreateTargetLookup {
	return newScanCreateTargetLookupAdapter(repo)
}

func NewTaskStoreAdapter(repo scanrepo.ScanTaskRepository) scanapp.TaskStore {
	return newScanTaskStoreAdapter(repo)
}

func NewTaskRuntimeScanStoreAdapter(repo *scanrepo.ScanRepository) scanapp.TaskRuntimeScanStore {
	return newTaskRuntimeScanStoreAdapter(repo)
}
