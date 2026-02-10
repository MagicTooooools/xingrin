package snapshotwiring

import (
	scanrepo "github.com/yyhuni/lunafox/server/internal/modules/scan/repository"
	snapshotdomain "github.com/yyhuni/lunafox/server/internal/modules/snapshot/domain"
)

type snapshotScanLookupAdapter struct {
	repo *scanrepo.ScanRepository
}

func newSnapshotScanLookupAdapter(repo *scanrepo.ScanRepository) *snapshotScanLookupAdapter {
	return &snapshotScanLookupAdapter{repo: repo}
}

func (adapter *snapshotScanLookupAdapter) GetActiveByID(id int) (*snapshotdomain.ScanRef, error) {
	item, err := adapter.repo.GetActiveByID(id)
	if err != nil {
		return nil, err
	}
	return snapshotScanModelToDomain(item), nil
}

func (adapter *snapshotScanLookupAdapter) GetTargetByScanID(scanID int) (*snapshotdomain.ScanTargetRef, error) {
	item, err := adapter.repo.GetTargetByScanID(scanID)
	if err != nil {
		return nil, err
	}
	return snapshotScanTargetModelToDomain(item), nil
}
