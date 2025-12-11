package models

import (
	"time"
)

type Repository struct {
	ID             int64      `gorm:"column:id;primaryKey;autoIncrement"`
	CreatedAt      *time.Time `gorm:"column:created_at"`
	UpdatedAt      *time.Time `gorm:"column:updated_at"`
	DeletedAt      *time.Time `gorm:"column:deleted_at;index:idx_repositories_deleted_at"`
	GithubRepoID   *int64     `gorm:"column:github_repo_id;uniqueIndex:idx_repositories_github_repo_id"`
	Name           *string    `gorm:"column:name"`
	Owner          *string    `gorm:"column:owner"`
	URL            *string    `gorm:"column:url"`
	InstallationID *int64     `gorm:"column:installation_id"`
}

func (Repository) TableName() string {
	return "public.repositories"
}
