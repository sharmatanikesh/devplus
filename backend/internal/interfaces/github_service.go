package interfaces

import (
	"context"
	"devplus-backend/internal/models"
)

type GithubService interface {
	GetRepositories(ctx context.Context, userID string) ([]*models.Repository, error)
	SyncRepositories(ctx context.Context, userID string, token string) ([]*models.Repository, error)
	SyncPullRequests(ctx context.Context, repoID string, token string) ([]*models.PullRequest, error)
	GetPullRequests(ctx context.Context, userID string, owner, repo string) ([]*models.PullRequest, error)
	GetDashboardStats(ctx context.Context, userID string) (*models.DashboardStats, error)
	GetRecentPullRequests(ctx context.Context, userID string, limit int) ([]*models.PullRequest, error)
	GetRepository(ctx context.Context, userID string, id string) (*models.Repository, error)
	GetPullRequest(ctx context.Context, userID string, repoID string, number int) (*models.PullRequest, error)
	GetPullRequestByID(ctx context.Context, prID string) (*models.PullRequest, error)
	GetMetrics(ctx context.Context, userID string, filter models.MetricsFilter) (*models.DashboardStats, error)
	GetPersonalMetrics(ctx context.Context, userID string, token string, username string, days int) (*models.PersonalMetrics, error)
	GetRepositoryByGithubID(ctx context.Context, githubID int64) (*models.Repository, error)
	UpsertPullRequest(ctx context.Context, pr *models.PullRequest) error
	UpdatePullRequestAnalysis(ctx context.Context, prID string, summary, decision string) error
	AnalyzeRepository(ctx context.Context, repoID string) error
	UpdateRepositoryAnalysis(ctx context.Context, repoID string, summary string) error
	AnalyzePullRequest(ctx context.Context, repoID string, prNumber int) error
	TriggerReleaseRiskAnalysis(ctx context.Context, repoID string, owner string, name string, prData string) error
	UpdateReleaseRiskAnalysis(ctx context.Context, repoID string, riskScore int, changelog string, rawAnalysis string) error
	GetPullRequestsByRepoID(ctx context.Context, repoID string) ([]*models.PullRequest, error)
}
