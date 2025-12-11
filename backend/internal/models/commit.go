package models

import (
	"time"
)

type Commit struct {
	ID         int64       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CreatedAt  *time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  *time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt  *time.Time  `gorm:"column:deleted_at;index:idx_commits_deleted_at" json:"deleted_at"`
	SHA        *string     `gorm:"column:sha;uniqueIndex:idx_commits_sha" json:"sha"`
	Message    *string     `gorm:"column:message" json:"message"`
	AuthorName *string     `gorm:"column:author_name" json:"author_name"`
	RepoID     *int64      `gorm:"column:repo_id" json:"repo_id"`
	Repository *Repository `gorm:"foreignKey:RepoID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"repository"`
}

func (Commit) TableName() string {
	return "public.commits"
}
