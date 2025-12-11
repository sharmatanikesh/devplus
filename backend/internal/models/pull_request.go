package models

import (
	"time"
)

type PullRequest struct {
	ID         int64       `gorm:"column:id;primaryKey;autoIncrement"`
	CreatedAt  *time.Time  `gorm:"column:created_at"`
	UpdatedAt  *time.Time  `gorm:"column:updated_at"`
	DeletedAt  *time.Time  `gorm:"column:deleted_at;index:idx_pull_requests_deleted_at"`
	GithubPRID *int64      `gorm:"column:github_pr_id;uniqueIndex:idx_pull_requests_github_pr_id"`
	Number     *int64      `gorm:"column:number"`
	Title      *string     `gorm:"column:title"`
	State      *string     `gorm:"column:state"`
	RepoID     *int64      `gorm:"column:repo_id"`
	Repository *Repository `gorm:"foreignKey:RepoID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	AuthorID   *int64      `gorm:"column:author_id"`
	AuthorName *string     `gorm:"column:author_name"`
	AISummary  *string     `gorm:"column:ai_summary"`
	AIDecision *string     `gorm:"column:a_idecision"`
}

func (PullRequest) TableName() string {
	return "public.pull_requests"
}
