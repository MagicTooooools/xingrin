package model

import (
	"time"
)

// NucleiTemplateRepo represents a nuclei template git repository
type NucleiTemplateRepo struct {
	ID           int        `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string     `gorm:"column:name;size:200;uniqueIndex" json:"name"`
	RepoURL      string     `gorm:"column:repo_url;size:500" json:"repoUrl"`
	LocalPath    string     `gorm:"column:local_path;size:500" json:"localPath"`
	CommitHash   string     `gorm:"column:commit_hash;size:40" json:"commitHash"`
	LastSyncedAt *time.Time `gorm:"column:last_synced_at" json:"lastSyncedAt"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

// TableName returns the table name for NucleiTemplateRepo
func (NucleiTemplateRepo) TableName() string {
	return "nuclei_template_repo"
}
