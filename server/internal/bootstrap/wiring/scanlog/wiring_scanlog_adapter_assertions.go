package scanlogwiring

import scanapp "github.com/yyhuni/lunafox/server/internal/modules/scan/application"

var _ scanapp.ScanLogStore = (*scanLogStoreAdapter)(nil)
var _ scanapp.ScanLogScanLookup = (*scanLogLookupAdapter)(nil)
