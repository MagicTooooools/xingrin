package scanwiring

import (
	scanapp "github.com/yyhuni/lunafox/server/internal/modules/scan/application"
	scanrepo "github.com/yyhuni/lunafox/server/internal/modules/scan/repository"
)

type taskRuntimeScanStoreAdapter struct{ repo *scanrepo.ScanRepository }

func newTaskRuntimeScanStoreAdapter(repo *scanrepo.ScanRepository) *taskRuntimeScanStoreAdapter {
	return &taskRuntimeScanStoreAdapter{repo: repo}
}

func (adapter *taskRuntimeScanStoreAdapter) FindByIDWithTarget(id int) (*scanapp.TaskScanRecord, error) {
	scan, err := adapter.repo.FindByIDWithTarget(id)
	if err != nil {
		return nil, err
	}
	result := &scanapp.TaskScanRecord{ID: scan.ID, TargetID: scan.TargetID, Status: scan.Status, YamlConfiguration: scan.YamlConfiguration}
	if scan.Target != nil {
		result.Target = &scanapp.TaskTargetRef{Name: scan.Target.Name, Type: scan.Target.Type}
	}
	return result, nil
}

func (adapter *taskRuntimeScanStoreAdapter) UpdateStatus(id int, status string, args ...string) error {
	return adapter.repo.UpdateStatus(id, status, args...)
}
