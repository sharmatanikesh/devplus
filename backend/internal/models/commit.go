package models

import (
	"time"
)

type Commit struct {
	ID         int64       `gorm:"column:id;primaryKey;autoIncrement"`
	CreatedAt  *time.Time  `gorm:"column:created_at"`
	UpdatedAt  *time.Time  `gorm:"column:updated_at"`
	DeletedAt  *time.Time  `gorm:"column:deleted_at;index:idx_commits_deleted_at"`
	SHA        *string     `gorm:"column:sha;uniqueIndex:idx_commits_sha"`
	Message    *string     `gorm:"column:message"`
	AuthorName *string     `gorm:"column:author_name"`
	RepoID     *int64      `gorm:"column:repo_id"`
	Repository *Repository `gorm:"foreignKey:RepoID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (Commit) TableName() string {
	return "public.commits"
}
