package interfaces

import (
	"context"
	"devplus-backend/internal/models"
)

type GithubService interface {
	FetchUserRepositories(ctx context.Context, token string) ([]*models.Repository, error)
	FetchRepositoryPullRequests(ctx context.Context, token string, owner, repo string) ([]*models.PullRequest, error)
}
