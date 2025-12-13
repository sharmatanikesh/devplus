package interfaces

import (
	"context"
	"devplus-backend/internal/models"
)

type GithubService interface {
	GetRepositories(ctx context.Context) ([]*models.Repository, error)
	SyncRepositories(ctx context.Context, userID string, token string) ([]*models.Repository, error)
	SyncPullRequests(ctx context.Context, repoID string, token string) ([]*models.PullRequest, error)
	GetPullRequests(ctx context.Context, owner, repo string) ([]*models.PullRequest, error)
	GetDashboardStats(ctx context.Context) (*models.DashboardStats, error)
	GetRecentPullRequests(ctx context.Context, limit int) ([]*models.PullRequest, error)
	GetRepository(ctx context.Context, id string) (*models.Repository, error)
	GetPullRequest(ctx context.Context, repoID string, number int) (*models.PullRequest, error)
	GetMetrics(ctx context.Context, filter models.MetricsFilter) (*models.DashboardStats, error)
	GetRepositoryByGithubID(ctx context.Context, githubID int64) (*models.Repository, error)
	UpsertPullRequest(ctx context.Context, pr *models.PullRequest) error
	UpdatePullRequestAnalysis(ctx context.Context, prID string, summary, decision string) error
	AnalyzeRepository(ctx context.Context, repoID string) error
	UpdateRepositoryAnalysis(ctx context.Context, repoID string, summary string) error
	AnalyzePullRequest(ctx context.Context, repoID string, prNumber int) error
}
