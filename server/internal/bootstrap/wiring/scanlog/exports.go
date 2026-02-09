package scanlogwiring

import (
	scanapp "github.com/yyhuni/lunafox/server/internal/modules/scan/application"
	scanrepo "github.com/yyhuni/lunafox/server/internal/modules/scan/repository"
)

func NewApplicationService(scanLogRepo *scanrepo.ScanLogRepository, scanRepo *scanrepo.ScanRepository) *scanapp.ScanLogService {
	return newScanLogApplicationService(scanLogRepo, scanRepo)
}
