package workerwiring

import (
	catalogapp "github.com/yyhuni/lunafox/server/internal/modules/catalog/application"
	catalogrepo "github.com/yyhuni/lunafox/server/internal/modules/catalog/repository"
	scanrepo "github.com/yyhuni/lunafox/server/internal/modules/scan/repository"
)

func NewApplicationService(scanRepo *scanrepo.ScanRepository, settingsRepo *catalogrepo.SubfinderProviderSettingsRepository) *catalogapp.WorkerService {
	return newWorkerApplicationService(scanRepo, settingsRepo)
}
