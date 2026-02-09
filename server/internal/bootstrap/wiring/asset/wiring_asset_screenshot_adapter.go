package assetwiring

import (
	assetdomain "github.com/yyhuni/lunafox/server/internal/modules/asset/domain"
	assetrepo "github.com/yyhuni/lunafox/server/internal/modules/asset/repository"
)

type assetScreenshotStoreAdapter struct {
	repo *assetrepo.ScreenshotRepository
}

func newAssetScreenshotStoreAdapter(repo *assetrepo.ScreenshotRepository) *assetScreenshotStoreAdapter {
	return &assetScreenshotStoreAdapter{repo: repo}
}

func (adapter *assetScreenshotStoreAdapter) FindByTargetID(targetID int, page, pageSize int, filter string) ([]assetdomain.Screenshot, int64, error) {
	items, total, err := adapter.repo.FindByTargetID(targetID, page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}
	results := make([]assetdomain.Screenshot, 0, len(items))
	for index := range items {
		results = append(results, *assetModelScreenshotToDomain(&items[index]))
	}
	return results, total, nil
}

func (adapter *assetScreenshotStoreAdapter) FindByID(id int) (*assetdomain.Screenshot, error) {
	item, err := adapter.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return assetModelScreenshotToDomain(item), nil
}

func (adapter *assetScreenshotStoreAdapter) BulkDelete(ids []int) (int64, error) {
	return adapter.repo.BulkDelete(ids)
}

func (adapter *assetScreenshotStoreAdapter) BulkUpsert(items []assetdomain.Screenshot) (int64, error) {
	return adapter.repo.BulkUpsert(assetDomainScreenshotListToModel(items))
}
