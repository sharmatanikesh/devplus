package github_service

import (
	"context"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"

	"devplus-backend/internal/models"
	"devplus-backend/internal/repositories"
)

type GithubService struct {
	repo repositories.GithubRepository
}

func NewGithubService(repo repositories.GithubRepository) *GithubService {
	return &GithubService{repo: repo}
}

func (s *GithubService) SyncRepositories(ctx context.Context, userID string, token string) ([]*models.Repository, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	opt := &github.RepositoryListOptions{
		Sort:        "updated",
		Direction:   "desc",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	repos, _, err := client.Repositories.List(ctx, "", opt)
	if err != nil {
		return nil, err
	}

	var result []*models.Repository
	for _, repo := range repos {
		r := &models.Repository{
			GithubRepoID: repo.ID,
			Name:         repo.GetName(),
			Owner:        repo.GetOwner().GetLogin(),
			URL:          repo.GetHTMLURL(),
			UpdatedAt:    &repo.UpdatedAt.Time,
			UserID:       userID,
		}

		// Use Repository for Upsert
		if err := s.repo.UpsertRepository(ctx, r); err != nil {
			return nil, err
		}
		result = append(result, r)
	}

	return result, nil
}

func (s *GithubService) GetRepositories(ctx context.Context) ([]*models.Repository, error) {
	return s.repo.GetRepositories(ctx)
}

func (s *GithubService) GetRepository(ctx context.Context, id string) (*models.Repository, error) {
	return s.repo.GetRepository(ctx, id)
}

func (s *GithubService) GetPullRequests(ctx context.Context, owner, repo string) ([]*models.PullRequest, error) {
	return s.repo.GetPullRequests(ctx, owner, repo)
}

func (s *GithubService) GetPullRequest(ctx context.Context, repoID string, number int) (*models.PullRequest, error) {
	return s.repo.GetPullRequest(ctx, repoID, number)
}

func (s *GithubService) GetMetrics(ctx context.Context, filter models.MetricsFilter) (*models.DashboardStats, error) {
	return s.repo.GetMetrics(ctx, filter)
}

func (s *GithubService) GetDashboardStats(ctx context.Context) (*models.DashboardStats, error) {
	return s.repo.GetDashboardStats(ctx)
}

func (s *GithubService) GetRecentPullRequests(ctx context.Context, limit int) ([]*models.PullRequest, error) {
	return s.repo.GetRecentPullRequests(ctx, limit)
}

func (s *GithubService) GetRepositoryByGithubID(ctx context.Context, githubID int64) (*models.Repository, error) {
	return s.repo.GetRepositoryByGithubID(ctx, githubID)
}

func (s *GithubService) UpsertPullRequest(ctx context.Context, pr *models.PullRequest) error {
	return s.repo.UpsertPullRequest(ctx, pr)
}

func (s *GithubService) UpdatePullRequestAnalysis(ctx context.Context, prID string, summary, decision string) error {
	return s.repo.UpdatePullRequestAnalysis(ctx, prID, summary, decision)
}
