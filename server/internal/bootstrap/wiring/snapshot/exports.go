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

func NewWebsiteAssetSyncAdapter(service *assetapp.WebsiteFacade) *snapshotWebsiteAssetSyncAdapter {
	return newSnapshotWebsiteAssetSyncAdapter(service)
}

func NewSubdomainAssetSyncAdapter(service *assetapp.SubdomainFacade) *snapshotSubdomainAssetSyncAdapter {
	return newSnapshotSubdomainAssetSyncAdapter(service)
}

func NewEndpointAssetSyncAdapter(service *assetapp.EndpointFacade) *snapshotEndpointAssetSyncAdapter {
	return newSnapshotEndpointAssetSyncAdapter(service)
}

func NewDirectoryAssetSyncAdapter(service *assetapp.DirectoryFacade) *snapshotDirectoryAssetSyncAdapter {
	return newSnapshotDirectoryAssetSyncAdapter(service)
}

func NewHostPortAssetSyncAdapter(service *assetapp.HostPortFacade) *snapshotHostPortAssetSyncAdapter {
	return newSnapshotHostPortAssetSyncAdapter(service)
}

func NewScreenshotAssetSyncAdapter(service *assetapp.ScreenshotFacade) *snapshotScreenshotAssetSyncAdapter {
	return newSnapshotScreenshotAssetSyncAdapter(service)
}

func NewVulnerabilityAssetSyncAdapter(service *securityapp.VulnerabilityFacade) *snapshotVulnerabilityAssetSyncAdapter {
	return newSnapshotVulnerabilityAssetSyncAdapter(service)
}

func NewVulnerabilityRawOutputCodec() *snapshotVulnerabilityRawOutputCodec {
	return newSnapshotVulnerabilityRawOutputCodec()
}
