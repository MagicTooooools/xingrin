package repository

import (
	"github.com/yyhuni/lunafox/server/internal/modules/snapshot/repository/persistence"
	"gorm.io/gorm/clause"
)

// BulkCreate creates multiple directory snapshots, ignoring duplicates
func (r *DirectorySnapshotRepository) BulkCreate(snapshots []model.DirectorySnapshot) (int64, error) {
	if len(snapshots) == 0 {
		return 0, nil
	}

	var totalAffected int64

	batchSize := 100
	for i := 0; i < len(snapshots); i += batchSize {
		end := i + batchSize
		if end > len(snapshots) {
			end = len(snapshots)
		}
		batch := snapshots[i:end]

		result := r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&batch)
		if result.Error != nil {
			return totalAffected, result.Error
		}
		totalAffected += result.RowsAffected
	}

	return totalAffected, nil
}
