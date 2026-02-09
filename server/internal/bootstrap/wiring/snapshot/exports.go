package snapshotwiring

import (
	assetapp "github.com/yyhuni/lunafox/server/internal/modules/asset/application"
	scanrepo "github.com/yyhuni/lunafox/server/internal/modules/scan/repository"
	securityapp "github.com/yyhuni/lunafox/server/internal/modules/security/application"
	snapshotrepo "github.com/yyhuni/lunafox/server/internal/modules/snapshot/repository"
)

func NewScanLookupAdapter(repo *scanrepo.ScanRepository) *snapshotScanLookupAdapter {
	return newSnapshotScanLookupAdapter(repo)
}

func NewWebsiteStoreAdapter(repo *snapshotrepo.WebsiteSnapshotRepository) *snapshotWebsiteStoreAdapter {
	return newSnapshotWebsiteStoreAdapter(repo)
}

func NewSubdomainStoreAdapter(repo *snapshotrepo.SubdomainSnapshotRepository) *snapshotSubdomainStoreAdapter {
	return newSnapshotSubdomainStoreAdapter(repo)
}

func NewEndpointStoreAdapter(repo *snapshotrepo.EndpointSnapshotRepository) *snapshotEndpointStoreAdapter {
	return newSnapshotEndpointStoreAdapter(repo)
}

func NewDirectoryStoreAdapter(repo *snapshotrepo.DirectorySnapshotRepository) *snapshotDirectoryStoreAdapter {
	return newSnapshotDirectoryStoreAdapter(repo)
}

func NewHostPortStoreAdapter(repo *snapshotrepo.HostPortSnapshotRepository) *snapshotHostPortStoreAdapter {
	return newSnapshotHostPortStoreAdapter(repo)
}

func NewScreenshotStoreAdapter(repo *snapshotrepo.ScreenshotSnapshotRepository) *snapshotScreenshotStoreAdapter {
	return newSnapshotScreenshotStoreAdapter(repo)
}

func NewVulnerabilityStoreAdapter(repo *snapshotrepo.VulnerabilitySnapshotRepository) *snapshotVulnerabilityStoreAdapter {
	return newSnapshotVulnerabilityStoreAdapter(repo)
}

func NewWebsiteAssetSyncAdapter(service *assetapp.WebsiteService) *snapshotWebsiteAssetSyncAdapter {
	return newSnapshotWebsiteAssetSyncAdapter(service)
}

func NewSubdomainAssetSyncAdapter(service *assetapp.SubdomainService) *snapshotSubdomainAssetSyncAdapter {
	return newSnapshotSubdomainAssetSyncAdapter(service)
}

func NewEndpointAssetSyncAdapter(service *assetapp.EndpointService) *snapshotEndpointAssetSyncAdapter {
	return newSnapshotEndpointAssetSyncAdapter(service)
}

func NewDirectoryAssetSyncAdapter(service *assetapp.DirectoryService) *snapshotDirectoryAssetSyncAdapter {
	return newSnapshotDirectoryAssetSyncAdapter(service)
}

func NewHostPortAssetSyncAdapter(service *assetapp.HostPortService) *snapshotHostPortAssetSyncAdapter {
	return newSnapshotHostPortAssetSyncAdapter(service)
}

func NewScreenshotAssetSyncAdapter(service *assetapp.ScreenshotService) *snapshotScreenshotAssetSyncAdapter {
	return newSnapshotScreenshotAssetSyncAdapter(service)
}

func NewVulnerabilityAssetSyncAdapter(service *securityapp.VulnerabilityService) *snapshotVulnerabilityAssetSyncAdapter {
	return newSnapshotVulnerabilityAssetSyncAdapter(service)
}

func NewVulnerabilityRawOutputCodec() *snapshotVulnerabilityRawOutputCodec {
	return newSnapshotVulnerabilityRawOutputCodec()
}
