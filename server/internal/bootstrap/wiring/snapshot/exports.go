package snapshotwiring

import (
	assetapp "github.com/yyhuni/lunafox/server/internal/modules/asset/application"
	scanrepo "github.com/yyhuni/lunafox/server/internal/modules/scan/repository"
	securityapp "github.com/yyhuni/lunafox/server/internal/modules/security/application"
	snapshotapp "github.com/yyhuni/lunafox/server/internal/modules/snapshot/application"
	snapshotrepo "github.com/yyhuni/lunafox/server/internal/modules/snapshot/repository"
)

func NewSnapshotScanRefLookupAdapter(repo *scanrepo.ScanRepository) snapshotapp.SnapshotScanRefLookup {
	return newSnapshotScanRefLookupAdapter(repo)
}

func NewSnapshotWebsiteStoreAdapter(repo *snapshotrepo.WebsiteSnapshotRepository) snapshotapp.WebsiteSnapshotStore {
	return newSnapshotWebsiteStoreAdapter(repo)
}

func NewSnapshotSubdomainStoreAdapter(repo *snapshotrepo.SubdomainSnapshotRepository) snapshotapp.SubdomainSnapshotStore {
	return newSnapshotSubdomainStoreAdapter(repo)
}

func NewSnapshotEndpointStoreAdapter(repo *snapshotrepo.EndpointSnapshotRepository) snapshotapp.EndpointSnapshotStore {
	return newSnapshotEndpointStoreAdapter(repo)
}

func NewSnapshotDirectoryStoreAdapter(repo *snapshotrepo.DirectorySnapshotRepository) snapshotapp.DirectorySnapshotStore {
	return newSnapshotDirectoryStoreAdapter(repo)
}

func NewSnapshotHostPortStoreAdapter(repo *snapshotrepo.HostPortSnapshotRepository) snapshotapp.HostPortSnapshotStore {
	return newSnapshotHostPortStoreAdapter(repo)
}

func NewSnapshotScreenshotStoreAdapter(repo *snapshotrepo.ScreenshotSnapshotRepository) snapshotapp.ScreenshotSnapshotStore {
	return newSnapshotScreenshotStoreAdapter(repo)
}

func NewSnapshotVulnerabilityStoreAdapter(repo *snapshotrepo.VulnerabilitySnapshotRepository) snapshotapp.VulnerabilitySnapshotStore {
	return newSnapshotVulnerabilityStoreAdapter(repo)
}

func NewSnapshotWebsiteAssetSyncAdapter(service *assetapp.WebsiteFacade) snapshotapp.WebsiteAssetSync {
	return newSnapshotWebsiteAssetSyncAdapter(service)
}

func NewSnapshotSubdomainAssetSyncAdapter(service *assetapp.SubdomainFacade) snapshotapp.SubdomainAssetSync {
	return newSnapshotSubdomainAssetSyncAdapter(service)
}

func NewSnapshotEndpointAssetSyncAdapter(service *assetapp.EndpointFacade) snapshotapp.EndpointAssetSync {
	return newSnapshotEndpointAssetSyncAdapter(service)
}

func NewSnapshotDirectoryAssetSyncAdapter(service *assetapp.DirectoryFacade) snapshotapp.DirectoryAssetSync {
	return newSnapshotDirectoryAssetSyncAdapter(service)
}

func NewSnapshotHostPortAssetSyncAdapter(service *assetapp.HostPortFacade) snapshotapp.HostPortAssetSync {
	return newSnapshotHostPortAssetSyncAdapter(service)
}

func NewSnapshotScreenshotAssetSyncAdapter(service *assetapp.ScreenshotFacade) snapshotapp.ScreenshotAssetSync {
	return newSnapshotScreenshotAssetSyncAdapter(service)
}

func NewSnapshotVulnerabilityAssetSyncAdapter(service *securityapp.VulnerabilityFacade) snapshotapp.VulnerabilityAssetSync {
	return newSnapshotVulnerabilityAssetSyncAdapter(service)
}

func NewSnapshotVulnerabilityRawOutputCodec() snapshotapp.VulnerabilityRawOutputCodec {
	return newSnapshotVulnerabilityRawOutputCodec()
}

func NewSnapshotWebsiteApplicationService(
	store snapshotapp.WebsiteSnapshotStore,
	scanLookup snapshotapp.SnapshotScanRefLookup,
	assetSync snapshotapp.WebsiteAssetSync,
) *snapshotapp.WebsiteSnapshotFacade {
	queryService := snapshotapp.NewWebsiteSnapshotQueryService(store, scanLookup)
	commandService := snapshotapp.NewWebsiteSnapshotCommandService(store, scanLookup, assetSync)
	return snapshotapp.NewWebsiteSnapshotFacade(queryService, commandService)
}

func NewSnapshotSubdomainApplicationService(
	store snapshotapp.SubdomainSnapshotStore,
	scanLookup snapshotapp.SnapshotScanRefLookup,
	assetSync snapshotapp.SubdomainAssetSync,
) *snapshotapp.SubdomainSnapshotFacade {
	queryService := snapshotapp.NewSubdomainSnapshotQueryService(store, scanLookup)
	commandService := snapshotapp.NewSubdomainSnapshotCommandService(store, scanLookup, assetSync)
	return snapshotapp.NewSubdomainSnapshotFacade(queryService, commandService)
}

func NewSnapshotEndpointApplicationService(
	store snapshotapp.EndpointSnapshotStore,
	scanLookup snapshotapp.SnapshotScanRefLookup,
	assetSync snapshotapp.EndpointAssetSync,
) *snapshotapp.EndpointSnapshotFacade {
	queryService := snapshotapp.NewEndpointSnapshotQueryService(store, scanLookup)
	commandService := snapshotapp.NewEndpointSnapshotCommandService(store, scanLookup, assetSync)
	return snapshotapp.NewEndpointSnapshotFacade(queryService, commandService)
}

func NewSnapshotDirectoryApplicationService(
	store snapshotapp.DirectorySnapshotStore,
	scanLookup snapshotapp.SnapshotScanRefLookup,
	assetSync snapshotapp.DirectoryAssetSync,
) *snapshotapp.DirectorySnapshotFacade {
	queryService := snapshotapp.NewDirectorySnapshotQueryService(store, scanLookup)
	commandService := snapshotapp.NewDirectorySnapshotCommandService(store, scanLookup, assetSync)
	return snapshotapp.NewDirectorySnapshotFacade(queryService, commandService)
}

func NewSnapshotHostPortApplicationService(
	store snapshotapp.HostPortSnapshotStore,
	scanLookup snapshotapp.SnapshotScanRefLookup,
	assetSync snapshotapp.HostPortAssetSync,
) *snapshotapp.HostPortSnapshotFacade {
	queryService := snapshotapp.NewHostPortSnapshotQueryService(store, scanLookup)
	commandService := snapshotapp.NewHostPortSnapshotCommandService(store, scanLookup, assetSync)
	return snapshotapp.NewHostPortSnapshotFacade(queryService, commandService)
}

func NewSnapshotScreenshotApplicationService(
	store snapshotapp.ScreenshotSnapshotStore,
	scanLookup snapshotapp.SnapshotScanRefLookup,
	assetSync snapshotapp.ScreenshotAssetSync,
) *snapshotapp.ScreenshotSnapshotFacade {
	queryService := snapshotapp.NewScreenshotSnapshotQueryService(store, scanLookup)
	commandService := snapshotapp.NewScreenshotSnapshotCommandService(store, scanLookup, assetSync)
	return snapshotapp.NewScreenshotSnapshotFacade(queryService, commandService)
}

func NewSnapshotVulnerabilityApplicationService(
	store snapshotapp.VulnerabilitySnapshotStore,
	scanLookup snapshotapp.SnapshotScanRefLookup,
	assetSync snapshotapp.VulnerabilityAssetSync,
	rawOutputCodec snapshotapp.VulnerabilityRawOutputCodec,
) *snapshotapp.VulnerabilitySnapshotFacade {
	queryService := snapshotapp.NewVulnerabilitySnapshotQueryService(store, scanLookup)
	commandService := snapshotapp.NewVulnerabilitySnapshotCommandService(store, scanLookup, assetSync, rawOutputCodec)
	return snapshotapp.NewVulnerabilitySnapshotFacade(queryService, commandService)
}
