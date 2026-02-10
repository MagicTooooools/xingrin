package scanwiring

import scanapp "github.com/yyhuni/lunafox/server/internal/modules/scan/application"

var _ scanapp.ScanQueryStore = (*scanStoreAdapter)(nil)
var _ scanapp.ScanCommandStore = (*scanStoreAdapter)(nil)
var _ scanapp.ScanStore = (*scanStoreAdapter)(nil)
var _ scanapp.CreateTargetLookup = (*scanTargetLookupAdapter)(nil)
var _ scanapp.ScanTaskCanceller = (*scanTaskCancellerAdapter)(nil)
var _ scanapp.TaskStore = (*scanTaskStoreAdapter)(nil)
var _ scanapp.TaskRuntimeScanStore = (*scanTaskRuntimeStoreAdapter)(nil)
