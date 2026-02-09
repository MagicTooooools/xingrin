package snapshotwiring

import (
	"database/sql"

	snapshotdomain "github.com/yyhuni/lunafox/server/internal/modules/snapshot/domain"
	snapshotrepo "github.com/yyhuni/lunafox/server/internal/modules/snapshot/repository"
	snapshotmodel "github.com/yyhuni/lunafox/server/internal/modules/snapshot/repository/persistence"
)

type snapshotDirectoryStoreAdapter struct {
	repo *snapshotrepo.DirectorySnapshotRepository
}

func newSnapshotDirectoryStoreAdapter(repo *snapshotrepo.DirectorySnapshotRepository) *snapshotDirectoryStoreAdapter {
	return &snapshotDirectoryStoreAdapter{repo: repo}
}

func (adapter *snapshotDirectoryStoreAdapter) BulkCreate(snapshots []snapshotdomain.DirectorySnapshot) (int64, error) {
	items := make([]snapshotmodel.DirectorySnapshot, 0, len(snapshots))
	for index := range snapshots {
		items = append(items, *snapshotDirectoryDomainToModel(&snapshots[index]))
	}
	return adapter.repo.BulkCreate(items)
}

func (adapter *snapshotDirectoryStoreAdapter) FindByScanID(scanID int, page, pageSize int, filter string) ([]snapshotdomain.DirectorySnapshot, int64, error) {
	items, total, err := adapter.repo.FindByScanID(scanID, page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}
	results := make([]snapshotdomain.DirectorySnapshot, 0, len(items))
	for index := range items {
		results = append(results, *snapshotDirectoryModelToDomain(&items[index]))
	}
	return results, total, nil
}

func (adapter *snapshotDirectoryStoreAdapter) StreamByScanID(scanID int) (*sql.Rows, error) {
	return adapter.repo.StreamByScanID(scanID)
}

func (adapter *snapshotDirectoryStoreAdapter) CountByScanID(scanID int) (int64, error) {
	return adapter.repo.CountByScanID(scanID)
}

func (adapter *snapshotDirectoryStoreAdapter) ScanRow(rows *sql.Rows) (*snapshotdomain.DirectorySnapshot, error) {
	item, err := adapter.repo.ScanRow(rows)
	if err != nil {
		return nil, err
	}
	return snapshotDirectoryModelToDomain(item), nil
}
