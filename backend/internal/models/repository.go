package models

import (
	"time"
)

type Repository struct {
	ID             int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CreatedAt      *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt      *time.Time `gorm:"column:deleted_at;index:idx_repositories_deleted_at" json:"deleted_at"`
	GithubRepoID   *int64     `gorm:"column:github_repo_id;uniqueIndex:idx_repositories_github_repo_id" json:"github_repo_id"`
	Name           *string    `gorm:"column:name" json:"name"`
	Owner          *string    `gorm:"column:owner" json:"owner"`
	URL            *string    `gorm:"column:url" json:"url"`
	InstallationID *int64     `gorm:"column:installation_id" json:"installation_id"`
}

func (Repository) TableName() string {
	return "public.repositories"
}
