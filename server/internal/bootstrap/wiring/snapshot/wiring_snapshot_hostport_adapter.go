package snapshotwiring

import (
	"database/sql"

	snapshotdomain "github.com/yyhuni/lunafox/server/internal/modules/snapshot/domain"
	snapshotrepo "github.com/yyhuni/lunafox/server/internal/modules/snapshot/repository"
	snapshotmodel "github.com/yyhuni/lunafox/server/internal/modules/snapshot/repository/persistence"
)

type snapshotHostPortStoreAdapter struct {
	repo *snapshotrepo.HostPortSnapshotRepository
}

func newSnapshotHostPortStoreAdapter(repo *snapshotrepo.HostPortSnapshotRepository) *snapshotHostPortStoreAdapter {
	return &snapshotHostPortStoreAdapter{repo: repo}
}

func (adapter *snapshotHostPortStoreAdapter) BulkCreate(snapshots []snapshotdomain.HostPortSnapshot) (int64, error) {
	items := make([]snapshotmodel.HostPortSnapshot, 0, len(snapshots))
	for index := range snapshots {
		items = append(items, *snapshotHostPortDomainToModel(&snapshots[index]))
	}
	return adapter.repo.BulkCreate(items)
}

func (adapter *snapshotHostPortStoreAdapter) FindByScanID(scanID int, page, pageSize int, filter string) ([]snapshotdomain.HostPortSnapshot, int64, error) {
	items, total, err := adapter.repo.FindByScanID(scanID, page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}
	results := make([]snapshotdomain.HostPortSnapshot, 0, len(items))
	for index := range items {
		results = append(results, *snapshotHostPortModelToDomain(&items[index]))
	}
	return results, total, nil
}

func (adapter *snapshotHostPortStoreAdapter) StreamByScanID(scanID int) (*sql.Rows, error) {
	return adapter.repo.StreamByScanID(scanID)
}

func (adapter *snapshotHostPortStoreAdapter) CountByScanID(scanID int) (int64, error) {
	return adapter.repo.CountByScanID(scanID)
}

func (adapter *snapshotHostPortStoreAdapter) ScanRow(rows *sql.Rows) (*snapshotdomain.HostPortSnapshot, error) {
	item, err := adapter.repo.ScanRow(rows)
	if err != nil {
		return nil, err
	}
	return snapshotHostPortModelToDomain(item), nil
}
