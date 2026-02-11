package workerwiring

import (
	catalogapp "github.com/yyhuni/lunafox/server/internal/modules/catalog/application"
	catalogrepo "github.com/yyhuni/lunafox/server/internal/modules/catalog/repository"
	scanrepo "github.com/yyhuni/lunafox/server/internal/modules/scan/repository"
)

func NewWorkerScanGuardAdapter(scanRepo *scanrepo.ScanRepository) catalogapp.WorkerScanGuard {
	return newWorkerScanGuardAdapter(scanRepo)
}

func NewWorkerProviderSettingsStoreAdapter(settingsRepo *catalogrepo.SubfinderProviderSettingsRepository) catalogapp.WorkerProviderSettingsStore {
	return newWorkerSettingsStoreAdapter(settingsRepo)
}

func NewWorkerApplicationService(scanGuard catalogapp.WorkerScanGuard, settingsStore catalogapp.WorkerProviderSettingsStore) *catalogapp.WorkerService {
	return catalogapp.NewWorkerService(scanGuard, settingsStore)
}
