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
	return adapter.repo.FindAll(page, pageSize, targetID, status, search)
}

func (adapter *scanStoreAdapter) FindByIDWithTarget(id int) (*scanapp.QueryScan, error) {
	return adapter.repo.FindByIDWithTarget(id)
}

func (adapter *scanStoreAdapter) GetActiveByID(id int) (*scanapp.QueryScan, error) {
	return adapter.repo.GetActiveByID(id)
}

func (adapter *scanStoreAdapter) FindByIDs(ids []int) ([]scanapp.QueryScan, error) {
	return adapter.repo.FindByIDs(ids)
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
	return adapter.repo.GetStatistics()
}

func (adapter *scanStoreAdapter) UpdateStatus(id int, status string, errorMessage ...string) error {
	return adapter.repo.UpdateStatus(id, status, errorMessage...)
}
