package snapshotwiring

import snapshotapp "github.com/yyhuni/lunafox/server/internal/modules/snapshot/application"

var _ snapshotapp.SnapshotScanRefLookup = (*snapshotScanRefLookupAdapter)(nil)

var _ snapshotapp.WebsiteSnapshotQueryStore = (*snapshotWebsiteStoreAdapter)(nil)
var _ snapshotapp.WebsiteSnapshotCommandStore = (*snapshotWebsiteStoreAdapter)(nil)
var _ snapshotapp.EndpointSnapshotQueryStore = (*snapshotEndpointStoreAdapter)(nil)
var _ snapshotapp.EndpointSnapshotCommandStore = (*snapshotEndpointStoreAdapter)(nil)
var _ snapshotapp.DirectorySnapshotQueryStore = (*snapshotDirectoryStoreAdapter)(nil)
var _ snapshotapp.DirectorySnapshotCommandStore = (*snapshotDirectoryStoreAdapter)(nil)
var _ snapshotapp.SubdomainSnapshotQueryStore = (*snapshotSubdomainStoreAdapter)(nil)
var _ snapshotapp.SubdomainSnapshotCommandStore = (*snapshotSubdomainStoreAdapter)(nil)
var _ snapshotapp.HostPortSnapshotQueryStore = (*snapshotHostPortStoreAdapter)(nil)
var _ snapshotapp.HostPortSnapshotCommandStore = (*snapshotHostPortStoreAdapter)(nil)
var _ snapshotapp.ScreenshotSnapshotQueryStore = (*snapshotScreenshotStoreAdapter)(nil)
var _ snapshotapp.ScreenshotSnapshotCommandStore = (*snapshotScreenshotStoreAdapter)(nil)
var _ snapshotapp.VulnerabilitySnapshotQueryStore = (*snapshotVulnerabilityStoreAdapter)(nil)
var _ snapshotapp.VulnerabilitySnapshotCommandStore = (*snapshotVulnerabilityStoreAdapter)(nil)

var _ snapshotapp.WebsiteAssetSync = (*snapshotWebsiteAssetSyncAdapter)(nil)
var _ snapshotapp.EndpointAssetSync = (*snapshotEndpointAssetSyncAdapter)(nil)
var _ snapshotapp.DirectoryAssetSync = (*snapshotDirectoryAssetSyncAdapter)(nil)
var _ snapshotapp.SubdomainAssetSync = (*snapshotSubdomainAssetSyncAdapter)(nil)
var _ snapshotapp.HostPortAssetSync = (*snapshotHostPortAssetSyncAdapter)(nil)
var _ snapshotapp.ScreenshotAssetSync = (*snapshotScreenshotAssetSyncAdapter)(nil)
var _ snapshotapp.VulnerabilityAssetSync = (*snapshotVulnerabilityAssetSyncAdapter)(nil)
var _ snapshotapp.VulnerabilityRawOutputCodec = (*snapshotVulnerabilityRawOutputCodec)(nil)
