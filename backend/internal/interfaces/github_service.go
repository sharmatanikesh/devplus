package interfaces

import (
	"context"
	"devplus-backend/internal/models"
)

type GithubService interface {
	SyncRepositories(ctx context.Context, token string) ([]*models.Repository, error)
	GetRepositories(ctx context.Context) ([]*models.Repository, error)
	FetchRepositoryPullRequests(ctx context.Context, token string, owner, repo string) ([]*models.PullRequest, error)
	GetDashboardStats(ctx context.Context) (*models.DashboardStats, error)
	GetRecentPullRequests(ctx context.Context, limit int) ([]*models.PullRequest, error)
}
