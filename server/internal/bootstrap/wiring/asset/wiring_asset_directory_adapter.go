package assetwiring

import (
	"database/sql"

	assetdomain "github.com/yyhuni/lunafox/server/internal/modules/asset/domain"
	assetrepo "github.com/yyhuni/lunafox/server/internal/modules/asset/repository"
)

type assetDirectoryStoreAdapter struct {
	repo *assetrepo.DirectoryRepository
}

func newAssetDirectoryStoreAdapter(repo *assetrepo.DirectoryRepository) *assetDirectoryStoreAdapter {
	return &assetDirectoryStoreAdapter{repo: repo}
}

func (adapter *assetDirectoryStoreAdapter) FindByTargetID(targetID int, page, pageSize int, filter string) ([]assetdomain.Directory, int64, error) {
	items, total, err := adapter.repo.FindByTargetID(targetID, page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}
	results := make([]assetdomain.Directory, 0, len(items))
	for index := range items {
		results = append(results, *assetModelDirectoryToDomain(&items[index]))
	}
	return results, total, nil
}

func (adapter *assetDirectoryStoreAdapter) BulkCreate(items []assetdomain.Directory) (int, error) {
	return adapter.repo.BulkCreate(assetDomainDirectoryListToModel(items))
}

func (adapter *assetDirectoryStoreAdapter) BulkDelete(ids []int) (int64, error) {
	return adapter.repo.BulkDelete(ids)
}

func (adapter *assetDirectoryStoreAdapter) BulkUpsert(items []assetdomain.Directory) (int64, error) {
	return adapter.repo.BulkUpsert(assetDomainDirectoryListToModel(items))
}

func (adapter *assetDirectoryStoreAdapter) StreamByTargetID(targetID int) (*sql.Rows, error) {
	return adapter.repo.StreamByTargetID(targetID)
}

func (adapter *assetDirectoryStoreAdapter) CountByTargetID(targetID int) (int64, error) {
	return adapter.repo.CountByTargetID(targetID)
}

func (adapter *assetDirectoryStoreAdapter) ScanRow(rows *sql.Rows) (*assetdomain.Directory, error) {
	item, err := adapter.repo.ScanRow(rows)
	if err != nil {
		return nil, err
	}
	return assetModelDirectoryToDomain(item), nil
}
