package scanlogwiring

import (
	scanapp "github.com/yyhuni/lunafox/server/internal/modules/scan/application"
	scanrepo "github.com/yyhuni/lunafox/server/internal/modules/scan/repository"
)

func NewScanLogStoreAdapter(repo *scanrepo.ScanLogRepository) scanapp.ScanLogStore {
	return newScanLogStoreAdapter(repo)
}

func NewScanLogScanLookupAdapter(repo *scanrepo.ScanRepository) scanapp.ScanLogScanLookup {
	return newScanLogLookupAdapter(repo)
}

func NewScanLogApplicationService(logStore scanapp.ScanLogStore, scanLookup scanapp.ScanLogScanLookup) scanapp.ScanLogApplicationService {
	return scanapp.NewScanLogService(logStore, scanLookup)
}
