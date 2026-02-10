package scanwiring

import (
	scanapp "github.com/yyhuni/lunafox/server/internal/modules/scan/application"
	scanrepo "github.com/yyhuni/lunafox/server/internal/modules/scan/repository"
)

type scanStoreAdapter struct {
	repo *scanrepo.ScanRepository
}

func newScanStoreAdapter(repo *scanrepo.ScanRepository) *scanStoreAdapter {
	return &scanStoreAdapter{repo: repo}
}

func (adapter *scanStoreAdapter) FindAll(page, pageSize int, targetID int, status, search string) ([]scanapp.QueryScan, int64, error) {
	scans, total, err := adapter.repo.FindAll(page, pageSize, targetID, status, search)
	if err != nil {
		return nil, 0, err
	}
	results := make([]scanapp.QueryScan, 0, len(scans))
	for index := range scans {
		results = append(results, *scanRecordToQueryScan(&scans[index]))
	}
	return results, total, nil
}

func (adapter *scanStoreAdapter) FindByIDWithTarget(id int) (*scanapp.QueryScan, error) {
	scan, err := adapter.repo.FindByIDWithTarget(id)
	if err != nil {
		return nil, err
	}
	return scanRecordToQueryScan(scan), nil
}

func (adapter *scanStoreAdapter) GetActiveByID(id int) (*scanapp.QueryScan, error) {
	scan, err := adapter.repo.GetActiveByID(id)
	if err != nil {
		return nil, err
	}
	return scanRecordToQueryScan(scan), nil
}

func (adapter *scanStoreAdapter) FindByIDs(ids []int) ([]scanapp.QueryScan, error) {
	scans, err := adapter.repo.FindByIDs(ids)
	if err != nil {
		return nil, err
	}
	results := make([]scanapp.QueryScan, 0, len(scans))
	for index := range scans {
		results = append(results, *scanRecordToQueryScan(&scans[index]))
	}
	return results, nil
}

func (adapter *scanStoreAdapter) CreateWithInputTargetsAndTasks(scan *scanapp.CreateScan, inputs []scanapp.CreateScanInputTarget, tasks []scanapp.CreateScanTask) error {
	repoScan := &scanrepo.ScanCreateRecord{
		TargetID:          scan.TargetID,
		EngineIDs:         scan.EngineIDs,
		EngineNames:       append([]byte(nil), scan.EngineNames...),
		YamlConfiguration: scan.YamlConfiguration,
		ScanMode:          scan.ScanMode,
		Status:            scan.Status,
	}

	repoInputs := make([]scanrepo.ScanInputTargetRecord, 0, len(inputs))
	for index := range inputs {
		repoInputs = append(repoInputs, scanrepo.ScanInputTargetRecord{Value: inputs[index].Value, InputType: inputs[index].InputType})
	}

	repoTasks := make([]scanrepo.ScanTaskCreateRecord, 0, len(tasks))
	for index := range tasks {
		repoTasks = append(repoTasks, scanrepo.ScanTaskCreateRecord{Stage: tasks[index].Stage, WorkflowName: tasks[index].WorkflowName, Status: tasks[index].Status})
	}

	if err := adapter.repo.CreateWithInputTargetsAndTasks(repoScan, repoInputs, repoTasks); err != nil {
		return err
	}

	scan.ID = repoScan.ID
	scan.CreatedAt = repoScan.CreatedAt
	return nil
}

func (adapter *scanStoreAdapter) BulkSoftDelete(ids []int) (int64, []string, error) {
	return adapter.repo.BulkSoftDelete(ids)
}

func (adapter *scanStoreAdapter) GetStatistics() (*scanapp.QueryStatistics, error) {
	stats, err := adapter.repo.GetStatistics()
	if err != nil {
		return nil, err
	}
	return &scanapp.QueryStatistics{
		Total:           stats.Total,
		Running:         stats.Running,
		Completed:       stats.Completed,
		Failed:          stats.Failed,
		TotalVulns:      stats.TotalVulns,
		TotalSubdomains: stats.TotalSubdomains,
		TotalEndpoints:  stats.TotalEndpoints,
		TotalWebsites:   stats.TotalWebsites,
		TotalAssets:     stats.TotalAssets,
	}, nil
}

func (adapter *scanStoreAdapter) UpdateStatus(id int, status string, errorMessage ...string) error {
	return adapter.repo.UpdateStatus(id, status, errorMessage...)
}
