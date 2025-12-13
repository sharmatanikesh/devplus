package repositories

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"devplus-backend/internal/models"
)

type GithubRepository interface {
	UpsertRepository(ctx context.Context, repo *models.Repository) error
	GetRepositories(ctx context.Context) ([]*models.Repository, error)
	GetRepository(ctx context.Context, id string) (*models.Repository, error)
	GetPullRequests(ctx context.Context, owner, repo string) ([]*models.PullRequest, error)
	GetPullRequest(ctx context.Context, repoID string, number int) (*models.PullRequest, error)
	GetDashboardStats(ctx context.Context) (*models.DashboardStats, error)
	GetMetrics(ctx context.Context, filter models.MetricsFilter) (*models.DashboardStats, error)
	GetRecentPullRequests(ctx context.Context, limit int) ([]*models.PullRequest, error)
	GetRepositoryByGithubID(ctx context.Context, githubID int64) (*models.Repository, error)
	UpsertPullRequest(ctx context.Context, pr *models.PullRequest) error
	UpdatePullRequestAnalysis(ctx context.Context, prID string, summary, decision string) error
	UpdateRepositoryAnalysis(ctx context.Context, repoID string, summary string) error
}

type gormGithubRepository struct {
	db *gorm.DB
}

func NewGithubRepository(db *gorm.DB) GithubRepository {
	return &gormGithubRepository{db: db}
}

func (r *gormGithubRepository) UpsertRepository(ctx context.Context, repo *models.Repository) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "github_repo_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "owner", "url", "updated_at", "user_id"}),
	}).Create(repo).Error
}

func (r *gormGithubRepository) GetRepositories(ctx context.Context) ([]*models.Repository, error) {
	var repos []*models.Repository
	if err := r.db.WithContext(ctx).Order("updated_at desc").Find(&repos).Error; err != nil {
		return nil, err
	}
	return repos, nil
}

func (r *gormGithubRepository) GetRepository(ctx context.Context, id string) (*models.Repository, error) {
	var repo models.Repository
	if err := r.db.WithContext(ctx).First(&repo, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &repo, nil
}

func (r *gormGithubRepository) GetPullRequests(ctx context.Context, owner, repo string) ([]*models.PullRequest, error) {
	var prs []*models.PullRequest
	// We need to join with Repositories to find by owner/repo name
	err := r.db.WithContext(ctx).Joins("JOIN repositories ON repositories.id = pull_requests.repo_id").
		Where("repositories.owner = ? AND repositories.name = ?", owner, repo).
		Order("pull_requests.number desc").
		Find(&prs).Error
	if err != nil {
		return nil, err
	}
	return prs, nil
}

func (r *gormGithubRepository) GetDashboardStats(ctx context.Context) (*models.DashboardStats, error) {
	var stats models.DashboardStats

	// Count PRs
	if err := r.db.WithContext(ctx).Model(&models.PullRequest{}).Count(&stats.TotalPRs).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Model(&models.PullRequest{}).Where("state = ?", "open").Count(&stats.OpenPRs).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Model(&models.PullRequest{}).Where("state = ?", "closed").Count(&stats.MergedPRs).Error; err != nil {
		return nil, err
	}

	// Active Repos (Total Repos for now)
	if err := r.db.WithContext(ctx).Model(&models.Repository{}).Count(&stats.ActiveRepos).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

func (r *gormGithubRepository) GetMetrics(ctx context.Context, filter models.MetricsFilter) (*models.DashboardStats, error) {
	var stats models.DashboardStats

	// Base scopes for PRs and Repos
	prScope := r.db.WithContext(ctx).Model(&models.PullRequest{})

	// Apply filters to PR Scope
	if filter.RepoID != "" {
		prScope = prScope.Where("repo_id = ?", filter.RepoID)
	}
	if filter.StartDate != nil {
		prScope = prScope.Where("created_at >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		prScope = prScope.Where("created_at <= ?", *filter.EndDate)
	}

	// Count PRs
	if err := prScope.Count(&stats.TotalPRs).Error; err != nil {
		return nil, err
	}
	// Re-apply scope for each count query, or clone.
	// GORM DB instance is safe to reuse if we build chain again or use Session.
	// But `prScope` variable holds the accumulated conditions.
	// Wait, `prScope` is `*gorm.DB`. When we call .Where it returns a NEW instance (cloned).
	// So `prScope` HAS the filters.
	// But when we call .Count, it executes.
	// We should be careful not to mutate `prScope` if we need it again?
	// Actually GORM methods like Where return a new instance.
	// `prScope = prScope.Where(...)` updates the variable to point to the new instance.
	// So `prScope` is the base with filters.

	if err := prScope.Where("state = ?", "open").Count(&stats.OpenPRs).Error; err != nil {
		return nil, err
	}
	if err := prScope.Where("state = ?", "closed").Count(&stats.MergedPRs).Error; err != nil {
		return nil, err
	}

	// Active Repos (Total Repos for now)
	// If filtering by RepoID, active repos is 1 (if exists) or 0?
	// Or should "Active Repos" ignore date filters?
	// Dashboard usually shows total repos under management.
	// If filtering by RepoID, it's 1.
	// Let's apply RepoID filter if present.
	repoScope := r.db.WithContext(ctx).Model(&models.Repository{})
	if filter.RepoID != "" {
		repoScope = repoScope.Where("id = ?", filter.RepoID)
	}
	// Date scope might apply to "updated_at" for repositories?
	// For now, let's just filter by RepoID for the "Active Repos" count.
	if err := repoScope.Count(&stats.ActiveRepos).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

func (r *gormGithubRepository) GetPullRequest(ctx context.Context, repoID string, number int) (*models.PullRequest, error) {
	var pr models.PullRequest
	if err := r.db.WithContext(ctx).Where("repo_id = ? AND number = ?", repoID, number).First(&pr).Error; err != nil {
		return nil, err
	}
	return &pr, nil
}

func (r *gormGithubRepository) GetRecentPullRequests(ctx context.Context, limit int) ([]*models.PullRequest, error) {
	var prs []*models.PullRequest
	if err := r.db.WithContext(ctx).Preload("Repository").Order("created_at desc").Limit(limit).Find(&prs).Error; err != nil {
		return nil, err
	}
	return prs, nil
}

func (r *gormGithubRepository) GetRepositoryByGithubID(ctx context.Context, githubID int64) (*models.Repository, error) {
	var repo models.Repository
	if err := r.db.WithContext(ctx).Where("github_repo_id = ?", githubID).First(&repo).Error; err != nil {
		return nil, err
	}
	return &repo, nil
}

func (r *gormGithubRepository) UpsertPullRequest(ctx context.Context, pr *models.PullRequest) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "github_pr_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"number", "title", "state", "updated_at", "author_id", "author_name"}),
	}).Create(pr).Error
}

func (r *gormGithubRepository) UpdatePullRequestAnalysis(ctx context.Context, prID string, summary, decision string) error {
	return r.db.WithContext(ctx).Model(&models.PullRequest{}).
		Where("id = ?", prID).
		Updates(map[string]interface{}{
			"ai_summary":  summary,
			"ai_decision": decision,
		}).Error
}

func (r *gormGithubRepository) UpdateRepositoryAnalysis(ctx context.Context, repoID string, summary string) error {
	return r.db.WithContext(ctx).Model(&models.Repository{}).
		Where("id = ?", repoID).
		Update("ai_summary", summary).Error
}
