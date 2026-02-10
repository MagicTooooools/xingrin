package repository

import (
	"github.com/yyhuni/lunafox/server/internal/modules/snapshot/repository/persistence"
	"gorm.io/gorm/clause"
)

// BulkCreate creates multiple website snapshots, ignoring duplicates
// Uses ON CONFLICT DO NOTHING based on unique constraint (scan_id + url)
func (r *WebsiteSnapshotRepository) BulkCreate(snapshots []model.WebsiteSnapshot) (int64, error) {
	if len(snapshots) == 0 {
		return 0, nil
	}

	var totalAffected int64

	// Process in batches to avoid SQL statement size limits
	batchSize := 500
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
