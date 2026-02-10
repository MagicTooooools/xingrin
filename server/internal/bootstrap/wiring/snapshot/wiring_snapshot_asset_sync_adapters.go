package snapshotwiring

import (
	assetapp "github.com/yyhuni/lunafox/server/internal/modules/asset/application"
	securityapp "github.com/yyhuni/lunafox/server/internal/modules/security/application"
	snapshotapp "github.com/yyhuni/lunafox/server/internal/modules/snapshot/application"
)

type snapshotWebsiteAssetSyncAdapter struct {
	service *assetapp.WebsiteFacade
}

func newSnapshotWebsiteAssetSyncAdapter(service *assetapp.WebsiteFacade) *snapshotWebsiteAssetSyncAdapter {
	return &snapshotWebsiteAssetSyncAdapter{service: service}
}

func (adapter *snapshotWebsiteAssetSyncAdapter) BulkUpsert(targetID int, items []snapshotapp.WebsiteAssetUpsertItem) (int64, error) {
	return adapter.service.BulkUpsert(targetID, snapshotWebsiteAssetUpsertItemsToDTO(items))
}

type snapshotEndpointAssetSyncAdapter struct {
	service *assetapp.EndpointFacade
}

func newSnapshotEndpointAssetSyncAdapter(service *assetapp.EndpointFacade) *snapshotEndpointAssetSyncAdapter {
	return &snapshotEndpointAssetSyncAdapter{service: service}
}

func (adapter *snapshotEndpointAssetSyncAdapter) BulkUpsert(targetID int, items []snapshotapp.EndpointAssetUpsertItem) (int64, error) {
	return adapter.service.BulkUpsert(targetID, snapshotEndpointAssetUpsertItemsToDTO(items))
}

type snapshotDirectoryAssetSyncAdapter struct {
	service *assetapp.DirectoryFacade
}

func newSnapshotDirectoryAssetSyncAdapter(service *assetapp.DirectoryFacade) *snapshotDirectoryAssetSyncAdapter {
	return &snapshotDirectoryAssetSyncAdapter{service: service}
}

func (adapter *snapshotDirectoryAssetSyncAdapter) BulkUpsert(targetID int, items []snapshotapp.DirectoryAssetUpsertItem) (int64, error) {
	return adapter.service.BulkUpsert(targetID, snapshotDirectoryAssetUpsertItemsToDTO(items))
}

type snapshotSubdomainAssetSyncAdapter struct {
	service *assetapp.SubdomainFacade
}

func newSnapshotSubdomainAssetSyncAdapter(service *assetapp.SubdomainFacade) *snapshotSubdomainAssetSyncAdapter {
	return &snapshotSubdomainAssetSyncAdapter{service: service}
}

func (adapter *snapshotSubdomainAssetSyncAdapter) BulkCreate(targetID int, names []string) (int, error) {
	return adapter.service.BulkCreate(targetID, names)
}

type snapshotHostPortAssetSyncAdapter struct {
	service *assetapp.HostPortFacade
}

func newSnapshotHostPortAssetSyncAdapter(service *assetapp.HostPortFacade) *snapshotHostPortAssetSyncAdapter {
	return &snapshotHostPortAssetSyncAdapter{service: service}
}

func (adapter *snapshotHostPortAssetSyncAdapter) BulkUpsert(targetID int, items []snapshotapp.HostPortAssetItem) (int64, error) {
	return adapter.service.BulkUpsert(targetID, snapshotHostPortAssetItemsToDTO(items))
}

type snapshotScreenshotAssetSyncAdapter struct {
	service *assetapp.ScreenshotFacade
}

func newSnapshotScreenshotAssetSyncAdapter(service *assetapp.ScreenshotFacade) *snapshotScreenshotAssetSyncAdapter {
	return &snapshotScreenshotAssetSyncAdapter{service: service}
}

func (adapter *snapshotScreenshotAssetSyncAdapter) BulkUpsert(targetID int, req *snapshotapp.ScreenshotAssetUpsertRequest) (int64, error) {
	return adapter.service.BulkUpsert(targetID, snapshotScreenshotAssetRequestToDTO(req))
}

type snapshotVulnerabilityAssetSyncAdapter struct {
	service *securityapp.VulnerabilityFacade
}

func newSnapshotVulnerabilityAssetSyncAdapter(service *securityapp.VulnerabilityFacade) *snapshotVulnerabilityAssetSyncAdapter {
	return &snapshotVulnerabilityAssetSyncAdapter{service: service}
}

func (adapter *snapshotVulnerabilityAssetSyncAdapter) BulkCreate(targetID int, items []snapshotapp.VulnerabilityAssetCreateItem) (int64, error) {
	return adapter.service.BulkCreate(targetID, snapshotVulnerabilityAssetCreateItemsToDTO(items))
}
