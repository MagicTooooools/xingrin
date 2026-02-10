package scanlogwiring

import (
	scanapp "github.com/yyhuni/lunafox/server/internal/modules/scan/application"
	scanrepo "github.com/yyhuni/lunafox/server/internal/modules/scan/repository"
)

type scanLogLookupAdapter struct {
	repo *scanrepo.ScanRepository
}

func (adapter *scanLogLookupAdapter) GetActiveByID(id int) (*scanapp.ScanLogScanRef, error) {
	scan, err := adapter.repo.GetActiveByID(id)
	if err != nil {
		return nil, err
	}
	return &scanapp.ScanLogScanRef{ID: scan.ID}, nil
}

func newScanLogApplicationService(scanLogRepo *scanrepo.ScanLogRepository, scanRepo *scanrepo.ScanRepository) *scanapp.ScanLogService {
	return scanapp.NewScanLogService(
		&scanLogStoreAdapter{repo: scanLogRepo},
		&scanLogLookupAdapter{repo: scanRepo},
	)
}
