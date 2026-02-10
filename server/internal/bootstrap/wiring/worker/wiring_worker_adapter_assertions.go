package workerwiring

import catalogapp "github.com/yyhuni/lunafox/server/internal/modules/catalog/application"

var _ catalogapp.WorkerScanGuard = (*workerScanGuardAdapter)(nil)
var _ catalogapp.WorkerProviderSettingsStore = (*workerSettingsStoreAdapter)(nil)
