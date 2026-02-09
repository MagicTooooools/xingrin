package assetwiring

import assetapp "github.com/yyhuni/lunafox/server/internal/modules/asset/application"

var _ assetapp.WebsiteStore = (*assetWebsiteStoreAdapter)(nil)
var _ assetapp.EndpointStore = (*assetEndpointStoreAdapter)(nil)
var _ assetapp.DirectoryStore = (*assetDirectoryStoreAdapter)(nil)
var _ assetapp.SubdomainStore = (*assetSubdomainStoreAdapter)(nil)
var _ assetapp.ScreenshotStore = (*assetScreenshotStoreAdapter)(nil)
var _ assetapp.HostPortStore = (*assetHostPortStoreAdapter)(nil)

var _ assetapp.WebsiteTargetLookup = (*assetTargetLookupAdapter)(nil)
var _ assetapp.EndpointTargetLookup = (*assetTargetLookupAdapter)(nil)
var _ assetapp.DirectoryTargetLookup = (*assetTargetLookupAdapter)(nil)
var _ assetapp.SubdomainTargetLookup = (*assetTargetLookupAdapter)(nil)
var _ assetapp.ScreenshotTargetLookup = (*assetTargetLookupAdapter)(nil)
var _ assetapp.HostPortTargetLookup = (*assetTargetLookupAdapter)(nil)
