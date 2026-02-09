package assetwiring

import (
	assetrepo "github.com/yyhuni/lunafox/server/internal/modules/asset/repository"
	catalogrepo "github.com/yyhuni/lunafox/server/internal/modules/catalog/repository"
)

func NewTargetLookupAdapter(repo *catalogrepo.TargetRepository) *assetTargetLookupAdapter {
	return newAssetTargetLookupAdapter(repo)
}

func NewWebsiteStoreAdapter(repo *assetrepo.WebsiteRepository) *assetWebsiteStoreAdapter {
	return newAssetWebsiteStoreAdapter(repo)
}

func NewSubdomainStoreAdapter(repo *assetrepo.SubdomainRepository) *assetSubdomainStoreAdapter {
	return newAssetSubdomainStoreAdapter(repo)
}

func NewEndpointStoreAdapter(repo *assetrepo.EndpointRepository) *assetEndpointStoreAdapter {
	return newAssetEndpointStoreAdapter(repo)
}

func NewDirectoryStoreAdapter(repo *assetrepo.DirectoryRepository) *assetDirectoryStoreAdapter {
	return newAssetDirectoryStoreAdapter(repo)
}

func NewHostPortStoreAdapter(repo *assetrepo.HostPortRepository) *assetHostPortStoreAdapter {
	return newAssetHostPortStoreAdapter(repo)
}

func NewScreenshotStoreAdapter(repo *assetrepo.ScreenshotRepository) *assetScreenshotStoreAdapter {
	return newAssetScreenshotStoreAdapter(repo)
}
