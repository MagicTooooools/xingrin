package snapshotwiring

import (
	"database/sql"

	snapshotdomain "github.com/yyhuni/lunafox/server/internal/modules/snapshot/domain"
	snapshotrepo "github.com/yyhuni/lunafox/server/internal/modules/snapshot/repository"
)

type snapshotSubdomainStoreAdapter struct {
	repo *snapshotrepo.SubdomainSnapshotRepository
}

func newSnapshotSubdomainStoreAdapter(repo *snapshotrepo.SubdomainSnapshotRepository) *snapshotSubdomainStoreAdapter {
	return &snapshotSubdomainStoreAdapter{repo: repo}
}

func (adapter *snapshotSubdomainStoreAdapter) BulkCreate(snapshots []snapshotdomain.SubdomainSnapshot) (int64, error) {
	return adapter.repo.BulkCreate(snapshots)
}

func (adapter *snapshotSubdomainStoreAdapter) FindByScanID(scanID int, page, pageSize int, filter string) ([]snapshotdomain.SubdomainSnapshot, int64, error) {
	return adapter.repo.FindByScanID(scanID, page, pageSize, filter)
}

func (adapter *snapshotSubdomainStoreAdapter) StreamByScanID(scanID int) (*sql.Rows, error) {
	return adapter.repo.StreamByScanID(scanID)
}

func (adapter *snapshotSubdomainStoreAdapter) CountByScanID(scanID int) (int64, error) {
	return adapter.repo.CountByScanID(scanID)
}

func (adapter *snapshotSubdomainStoreAdapter) ScanRow(rows *sql.Rows) (*snapshotdomain.SubdomainSnapshot, error) {
	return adapter.repo.ScanRow(rows)
}
