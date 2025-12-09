package models

import "gorm.io/gorm"

type PullRequestState string

const (
	PROpen   PullRequestState = "open"
	PRClosed PullRequestState = "closed"
	PRMerged PullRequestState = "merged"
)

type AISummaryDecision string

const (
	DecisionApprove AISummaryDecision = "APPROVE"
	DecisionBlock   AISummaryDecision = "BLOCK"
	DecisionComment AISummaryDecision = "COMMENT"
)

type PullRequest struct {
	gorm.Model
	GithubPRID int64             `json:"github_pr_id" gorm:"uniqueIndex"`
	Number     int               `json:"number"`
	Title      string            `json:"title"`
	State      PullRequestState  `json:"state"`
	RepoID     uint              `json:"repo_id"`
	Repository Repository        `json:"repository" gorm:"foreignKey:RepoID"`
	AuthorID   uint              `json:"author_id"` // Simplified: User ID if in system
	AuthorName string            `json:"author_name"`
	AISummary  string            `json:"ai_summary" gorm:"type:text"`
	AIDecision AISummaryDecision `json:"ai_decision"`
}
