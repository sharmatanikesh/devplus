package github_service

import (
	"context"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"devplus-backend/internal/models"
)

type GithubService struct {
	db *gorm.DB
}

func NewGithubService(db *gorm.DB) *GithubService {
	return &GithubService{db: db}
}

func (s *GithubService) SyncRepositories(ctx context.Context, token string) ([]*models.Repository, error) {
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
			Name:         repo.Name,
			Owner:        repo.Owner.Login,
			URL:          repo.HTMLURL,
			UpdatedAt:    &repo.UpdatedAt.Time,
		}

		// Upsert into DB
		if err := s.db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "github_repo_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"name", "owner", "url", "updated_at"}),
		}).Create(r).Error; err != nil {
			// Log error but maybe continue? For now let's return error
			return nil, err
		}
		result = append(result, r)
	}

	return result, nil
}

func (s *GithubService) GetRepositories(ctx context.Context) ([]*models.Repository, error) {
	var repos []*models.Repository
	if err := s.db.Order("updated_at desc").Find(&repos).Error; err != nil {
		return nil, err
	}
	return repos, nil
}

func (s *GithubService) FetchRepositoryPullRequests(ctx context.Context, token string, owner, repo string) ([]*models.PullRequest, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	opt := &github.PullRequestListOptions{
		State:       "all",
		Sort:        "created",
		Direction:   "desc",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	prs, _, err := client.PullRequests.List(ctx, owner, repo, opt)
	if err != nil {
		return nil, err
	}

	var result []*models.PullRequest
	for _, pr := range prs {
		result = append(result, &models.PullRequest{
			GithubPRID: pr.ID,
			Number:     interfaceToInt64(pr.Number),
			Title:      pr.Title,
			State:      pr.State,
			// Assuming AuthorName is user login
			AuthorName: pr.User.Login,
			// RepoID would be linked to our internal Repo ID if we had it here,
			// but we are returning API data.
		})
	}
	return result, nil
}

func (s *GithubService) GetDashboardStats(ctx context.Context) (*models.DashboardStats, error) {
	var stats models.DashboardStats

	// Count PRs
	if err := s.db.Model(&models.PullRequest{}).Count(&stats.TotalPRs).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&models.PullRequest{}).Where("state = ?", "open").Count(&stats.OpenPRs).Error; err != nil {
		return nil, err
	}
	// "merged" state in GitHub API is usually "closed" with merged_at != null, or explicit "merged" in our schema depending on how we saved it.
	// For simplicity, assuming "closed" state and we'd check merged status, but for now let's check for 'closed' as a proxy or if we saved 'merged'.
	// The PR model has 'State', let's assume 'closed' for now.
	// Or better, let's just count 'closed'.
	// If we want 'merged', we need to check the data we save.
	// Let's assume we save 'closed' for merged ones too in 'state'. A more accurate check requires a 'MergedAt' field or 'Merged' boolean.
	// I'll query for 'closed' for now.
	if err := s.db.Model(&models.PullRequest{}).Where("state = ?", "closed").Count(&stats.MergedPRs).Error; err != nil {
		return nil, err
	}

	// Active Repos (Total Repos for now)
	if err := s.db.Model(&models.Repository{}).Count(&stats.ActiveRepos).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

func (s *GithubService) GetRecentPullRequests(ctx context.Context, limit int) ([]*models.PullRequest, error) {
	var prs []*models.PullRequest
	// Preload Repository to show repo name details
	if err := s.db.Preload("Repository").Order("created_at desc").Limit(limit).Find(&prs).Error; err != nil {
		return nil, err
	}
	return prs, nil
}

func interfaceToInt64(i interface{}) *int64 {
	x, ok := i.(int)
	if ok {
		v := int64(x)
		return &v
	}
	// Handle other types if necessary or return nil
	return nil
}
