package scanwiring

import (
	scanapp "github.com/yyhuni/lunafox/server/internal/modules/scan/application"
	scanmodel "github.com/yyhuni/lunafox/server/internal/modules/scan/repository/persistence"
)

func scanModelToQueryScan(scan *scanmodel.Scan) *scanapp.QueryScan {
	if scan == nil {
		return nil
	}
	result := &scanapp.QueryScan{ID: scan.ID, TargetID: scan.TargetID, EngineIDs: scan.EngineIDs, EngineNames: scan.EngineNames, YamlConfiguration: scan.YamlConfiguration, ScanMode: scan.ScanMode, Status: scan.Status, ResultsDir: scan.ResultsDir, WorkerID: scan.WorkerID, ErrorMessage: scan.ErrorMessage, Progress: scan.Progress, CurrentStage: scan.CurrentStage, StageProgress: scan.StageProgress, CreatedAt: scan.CreatedAt, StoppedAt: scan.StoppedAt, CachedSubdomainsCount: scan.CachedSubdomainsCount, CachedWebsitesCount: scan.CachedWebsitesCount, CachedEndpointsCount: scan.CachedEndpointsCount, CachedIPsCount: scan.CachedIPsCount, CachedDirectoriesCount: scan.CachedDirectoriesCount, CachedScreenshotsCount: scan.CachedScreenshotsCount, CachedVulnsTotal: scan.CachedVulnsTotal, CachedVulnsCritical: scan.CachedVulnsCritical, CachedVulnsHigh: scan.CachedVulnsHigh, CachedVulnsMedium: scan.CachedVulnsMedium, CachedVulnsLow: scan.CachedVulnsLow}
	if scan.Target != nil {
		result.Target = &scanapp.QueryTargetRef{ID: scan.Target.ID, Name: scan.Target.Name, Type: scan.Target.Type}
	}
	return result
}
