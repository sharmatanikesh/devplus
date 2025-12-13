package models

import (
	"time"
)

type PullRequest struct {
	ID         string      `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CreatedAt  *time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  *time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt  *time.Time  `gorm:"column:deleted_at;index:idx_pull_requests_deleted_at" json:"deleted_at"`
	GithubPRID *int64      `gorm:"column:github_pr_id;uniqueIndex:idx_pull_requests_github_pr_id" json:"github_pr_id"`
	Number     *int64      `gorm:"column:number" json:"number"`
	Title      *string     `gorm:"column:title" json:"title"`
	State      *string     `gorm:"column:state" json:"state"`
	RepoID     *string     `gorm:"column:repo_id" json:"repo_id"`
	Repository *Repository `gorm:"foreignKey:RepoID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"repository"`
	AuthorID   *int64      `gorm:"column:author_id" json:"author_id"`
	AuthorName *string     `gorm:"column:author_name" json:"author_name"`
	AISummary  *string     `gorm:"column:ai_summary" json:"ai_summary"`
	AIDecision *string     `gorm:"column:ai_decision" json:"ai_decision"`
}

func (PullRequest) TableName() string {
	return "public.pull_requests"
}
