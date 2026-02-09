package snapshotwiring

import (
	assetapp "github.com/yyhuni/lunafox/server/internal/modules/asset/application"
	securityapp "github.com/yyhuni/lunafox/server/internal/modules/security/application"
	snapshotapp "github.com/yyhuni/lunafox/server/internal/modules/snapshot/application"
)

type snapshotWebsiteAssetSyncAdapter struct {
	service *assetapp.WebsiteService
}

func newSnapshotWebsiteAssetSyncAdapter(service *assetapp.WebsiteService) *snapshotWebsiteAssetSyncAdapter {
	return &snapshotWebsiteAssetSyncAdapter{service: service}
}

func (adapter *snapshotWebsiteAssetSyncAdapter) BulkUpsert(targetID int, items []snapshotapp.WebsiteAssetUpsertItem) (int64, error) {
	return adapter.service.BulkUpsert(targetID, snapshotWebsiteAssetUpsertItemsToDTO(items))
}

type snapshotEndpointAssetSyncAdapter struct {
	service *assetapp.EndpointService
}

func newSnapshotEndpointAssetSyncAdapter(service *assetapp.EndpointService) *snapshotEndpointAssetSyncAdapter {
	return &snapshotEndpointAssetSyncAdapter{service: service}
}

func (adapter *snapshotEndpointAssetSyncAdapter) BulkUpsert(targetID int, items []snapshotapp.EndpointAssetUpsertItem) (int64, error) {
	return adapter.service.BulkUpsert(targetID, snapshotEndpointAssetUpsertItemsToDTO(items))
}

type snapshotDirectoryAssetSyncAdapter struct {
	service *assetapp.DirectoryService
}

func newSnapshotDirectoryAssetSyncAdapter(service *assetapp.DirectoryService) *snapshotDirectoryAssetSyncAdapter {
	return &snapshotDirectoryAssetSyncAdapter{service: service}
}

func (adapter *snapshotDirectoryAssetSyncAdapter) BulkUpsert(targetID int, items []snapshotapp.DirectoryAssetUpsertItem) (int64, error) {
	return adapter.service.BulkUpsert(targetID, snapshotDirectoryAssetUpsertItemsToDTO(items))
}

type snapshotSubdomainAssetSyncAdapter struct {
	service *assetapp.SubdomainService
}

func newSnapshotSubdomainAssetSyncAdapter(service *assetapp.SubdomainService) *snapshotSubdomainAssetSyncAdapter {
	return &snapshotSubdomainAssetSyncAdapter{service: service}
}

func (adapter *snapshotSubdomainAssetSyncAdapter) BulkCreate(targetID int, names []string) (int, error) {
	return adapter.service.BulkCreate(targetID, names)
}

type snapshotHostPortAssetSyncAdapter struct {
	service *assetapp.HostPortService
}

func newSnapshotHostPortAssetSyncAdapter(service *assetapp.HostPortService) *snapshotHostPortAssetSyncAdapter {
	return &snapshotHostPortAssetSyncAdapter{service: service}
}

func (adapter *snapshotHostPortAssetSyncAdapter) BulkUpsert(targetID int, items []snapshotapp.HostPortAssetItem) (int64, error) {
	return adapter.service.BulkUpsert(targetID, snapshotHostPortAssetItemsToDTO(items))
}

type snapshotScreenshotAssetSyncAdapter struct {
	service *assetapp.ScreenshotService
}

func newSnapshotScreenshotAssetSyncAdapter(service *assetapp.ScreenshotService) *snapshotScreenshotAssetSyncAdapter {
	return &snapshotScreenshotAssetSyncAdapter{service: service}
}

func (adapter *snapshotScreenshotAssetSyncAdapter) BulkUpsert(targetID int, req *snapshotapp.ScreenshotAssetUpsertRequest) (int64, error) {
	return adapter.service.BulkUpsert(targetID, snapshotScreenshotAssetRequestToDTO(req))
}

type snapshotVulnerabilityAssetSyncAdapter struct {
	service *securityapp.VulnerabilityService
}

func newSnapshotVulnerabilityAssetSyncAdapter(service *securityapp.VulnerabilityService) *snapshotVulnerabilityAssetSyncAdapter {
	return &snapshotVulnerabilityAssetSyncAdapter{service: service}
}

func (adapter *snapshotVulnerabilityAssetSyncAdapter) BulkCreate(targetID int, items []snapshotapp.VulnerabilityAssetCreateItem) (int64, error) {
	return adapter.service.BulkCreate(targetID, snapshotVulnerabilityAssetCreateItemsToDTO(items))
}
