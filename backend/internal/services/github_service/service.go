package github_service

import (
	"context"
	"fmt"

	"github.com/google/go-github/v50/github"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"

	"devplus-backend/internal/models"
	"devplus-backend/internal/repositories"
	"devplus-backend/internal/services/ai"
)

type GithubService struct {
	repo       repositories.GithubRepository
	aiFactory  *ai.AIFactory
	backendURL string
}

func NewGithubService(repo repositories.GithubRepository, aiFactory *ai.AIFactory, backendURL string) *GithubService {
	return &GithubService{
		repo:       repo,
		aiFactory:  aiFactory,
		backendURL: backendURL,
	}
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

	// ... existing repo sync logic ...

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

		// Ensure we have the latest ID (UUID)
		if r.GithubRepoID == nil {
			// Should not happen as we just set it from repo.ID
			continue
		}
		savedRepo, err := s.repo.GetRepositoryByGithubID(ctx, *r.GithubRepoID)
		if err != nil {
			return nil, err
		}
		result = append(result, savedRepo)

		// Call new SyncPullRequests
		if _, err := s.SyncPullRequests(ctx, savedRepo.ID, token); err != nil {
			// Log error but continue
			// fmt.Printf("Failed to sync pull requests for repo %s: %v\n", r.Name, err)
			// Return error for now as requested by user to integrate it
			return nil, err
		}
	}

	return result, nil
}

func (s *GithubService) SyncPullRequests(ctx context.Context, repoID string, token string) ([]*models.PullRequest, error) {
	log.Info().Str("repo_id", repoID).Msg("[Service.SyncPullRequests] Syncing PRs")

	// 1. Get Repo ID to find owner/name
	repo, err := s.repo.GetRepository(ctx, repoID)
	if err != nil {
		log.Error().Err(err).Msg("[Service.SyncPullRequests] Failed to get repo")
		return nil, err
	}
	log.Info().Str("owner", repo.Owner).Str("repo", repo.Name).Msg("[Service.SyncPullRequests] Found Repo")

	// 2. Setup Client
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// 3. Fetch Open PRs
	prOpt := &github.PullRequestListOptions{
		State:       "open",
		ListOptions: github.ListOptions{PerPage: 50},
	}
	log.Info().Msg("[Service.SyncPullRequests] Fetching PRs from GitHub...")
	prs, _, err := client.PullRequests.List(ctx, repo.Owner, repo.Name, prOpt)
	if err != nil {
		log.Error().Err(err).Msg("[Service.SyncPullRequests] GitHub API Error")
		return nil, err
	}
	log.Info().Int("count", len(prs)).Msg("[Service.SyncPullRequests] Fetched PRs from GitHub")

	var syncedPRs []*models.PullRequest
	for _, pr := range prs {
		log.Info().Int("pr_number", pr.GetNumber()).Str("title", pr.GetTitle()).Msg("[Service.SyncPullRequests] Processing PR")
		prModel := &models.PullRequest{
			GithubPRID: pr.ID,
			Number:     github.Int64(int64(pr.GetNumber())),
			Title:      github.String(pr.GetTitle()),
			State:      github.String(pr.GetState()),
			RepoID:     &repo.ID,
			AuthorID:   pr.GetUser().ID,
			AuthorName: github.String(pr.GetUser().GetLogin()),
			CreatedAt:  &pr.CreatedAt.Time,
			UpdatedAt:  &pr.UpdatedAt.Time,
		}
		if err := s.repo.UpsertPullRequest(ctx, prModel); err != nil {
			log.Error().Int("pr_number", pr.GetNumber()).Err(err).Msg("[Service.SyncPullRequests] Upsert Error")
			return nil, err
		}
		syncedPRs = append(syncedPRs, prModel)
	}

	log.Info().Int("count", len(syncedPRs)).Msg("[Service.SyncPullRequests] Successfully synced PRs")
	return syncedPRs, nil
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

func (s *GithubService) AnalyzeRepository(ctx context.Context, repoID string) error {
	repo, err := s.repo.GetRepository(ctx, repoID)
	if err != nil {
		return err
	}

	// Trigger AI Service
	// Construct Callback URL
	callbackURL := fmt.Sprintf("%s/api/v1/webhook/ai/repo", s.backendURL)

	// Get AI Service (Kestra)
	aiService, err := s.aiFactory.GetAIService("kestra")
	if err != nil {
		return err
	}

	return aiService.AnalyzeRepo(ctx, repo, callbackURL)
}

func (s *GithubService) UpdateRepositoryAnalysis(ctx context.Context, repoID string, summary string) error {
	return s.repo.UpdateRepositoryAnalysis(ctx, repoID, summary)
}

func (s *GithubService) AnalyzePullRequest(ctx context.Context, repoID string, prNumber int) error {
	pr, err := s.repo.GetPullRequest(ctx, repoID, prNumber)
	if err != nil {
		return err
	}

	// Trigger AI Service
	// Always use "kestra" as it's the confirmed workflow engine
	aiService, err := s.aiFactory.GetAIService("kestra")
	if err != nil {
		return err
	}

	// Construct Callback URL
	callbackURL := s.backendURL + "/api/v1/webhook/ai"

	return aiService.AnalyzePR(ctx, pr, callbackURL)
}
