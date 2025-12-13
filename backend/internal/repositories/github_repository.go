package repositories

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"devplus-backend/internal/models"
)

type GithubRepository interface {
	UpsertRepository(ctx context.Context, repo *models.Repository) error
	GetRepositories(ctx context.Context, userID string) ([]*models.Repository, error)
	GetRepository(ctx context.Context, userID string, id string) (*models.Repository, error)
	GetPullRequests(ctx context.Context, userID string, owner, repo string) ([]*models.PullRequest, error)
	GetPullRequest(ctx context.Context, userID string, repoID string, number int) (*models.PullRequest, error)
	GetPullRequestByID(ctx context.Context, prID string) (*models.PullRequest, error)
	GetDashboardStats(ctx context.Context, userID string) (*models.DashboardStats, error)
	GetMetrics(ctx context.Context, userID string, filter models.MetricsFilter) (*models.DashboardStats, error)
	GetRecentPullRequests(ctx context.Context, userID string, limit int) ([]*models.PullRequest, error)
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

func (r *gormGithubRepository) GetRepositories(ctx context.Context, userID string) ([]*models.Repository, error) {
	var repos []*models.Repository
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("updated_at desc").Find(&repos).Error; err != nil {
		return nil, err
	}
	return repos, nil
}

func (r *gormGithubRepository) GetRepository(ctx context.Context, userID string, id string) (*models.Repository, error) {
	var repo models.Repository
	query := r.db.WithContext(ctx).Where("id = ?", id)
	// Only filter by userID if it's provided (for user-specific queries)
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if err := query.First(&repo).Error; err != nil {
		return nil, err
	}
	return &repo, nil
}

func (r *gormGithubRepository) GetPullRequests(ctx context.Context, userID string, owner, repo string) ([]*models.PullRequest, error) {
	var prs []*models.PullRequest
	// We need to join with Repositories to find by owner/repo name and filter by userID
	err := r.db.WithContext(ctx).Joins("JOIN repositories ON repositories.id = pull_requests.repo_id").
		Where("repositories.owner = ? AND repositories.name = ? AND repositories.user_id = ?", owner, repo, userID).
		Order("pull_requests.number desc").
		Find(&prs).Error
	if err != nil {
		return nil, err
	}
	return prs, nil
}

func (r *gormGithubRepository) GetDashboardStats(ctx context.Context, userID string) (*models.DashboardStats, error) {
	var stats models.DashboardStats

	// Count PRs for repos owned by this user
	prScope := r.db.WithContext(ctx).Model(&models.PullRequest{}).Joins("JOIN repositories ON repositories.id = pull_requests.repo_id").Where("repositories.user_id = ?", userID)
	
	if err := prScope.Count(&stats.TotalPRs).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Model(&models.PullRequest{}).Joins("JOIN repositories ON repositories.id = pull_requests.repo_id").Where("repositories.user_id = ? AND pull_requests.state = ?", userID, "open").Count(&stats.OpenPRs).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Model(&models.PullRequest{}).Joins("JOIN repositories ON repositories.id = pull_requests.repo_id").Where("repositories.user_id = ? AND pull_requests.state = ?", userID, "closed").Count(&stats.MergedPRs).Error; err != nil {
		return nil, err
	}

	// Active Repos (Total Repos for this user)
	if err := r.db.WithContext(ctx).Model(&models.Repository{}).Where("user_id = ?", userID).Count(&stats.ActiveRepos).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

func (r *gormGithubRepository) GetMetrics(ctx context.Context, userID string, filter models.MetricsFilter) (*models.DashboardStats, error) {
	var stats models.DashboardStats

	// Base scopes for PRs and Repos - filter by userID
	prScope := r.db.WithContext(ctx).Model(&models.PullRequest{}).Joins("JOIN repositories ON repositories.id = pull_requests.repo_id").Where("repositories.user_id = ?", userID)

	// Apply filters to PR Scope
	if filter.RepoID != "" {
		prScope = prScope.Where("pull_requests.repo_id = ?", filter.RepoID)
	}
	if filter.StartDate != nil {
		prScope = prScope.Where("pull_requests.created_at >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		prScope = prScope.Where("pull_requests.created_at <= ?", *filter.EndDate)
	}

	// Count PRs
	if err := prScope.Count(&stats.TotalPRs).Error; err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Model(&models.PullRequest{}).Joins("JOIN repositories ON repositories.id = pull_requests.repo_id").Where("repositories.user_id = ?", userID).Scopes(func(db *gorm.DB) *gorm.DB {
		if filter.RepoID != "" {
			db = db.Where("pull_requests.repo_id = ?", filter.RepoID)
		}
		if filter.StartDate != nil {
			db = db.Where("pull_requests.created_at >= ?", *filter.StartDate)
		}
		if filter.EndDate != nil {
			db = db.Where("pull_requests.created_at <= ?", *filter.EndDate)
		}
		return db
	}).Where("pull_requests.state = ?", "open").Count(&stats.OpenPRs).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Model(&models.PullRequest{}).Joins("JOIN repositories ON repositories.id = pull_requests.repo_id").Where("repositories.user_id = ?", userID).Scopes(func(db *gorm.DB) *gorm.DB {
		if filter.RepoID != "" {
			db = db.Where("pull_requests.repo_id = ?", filter.RepoID)
		}
		if filter.StartDate != nil {
			db = db.Where("pull_requests.created_at >= ?", *filter.StartDate)
		}
		if filter.EndDate != nil {
			db = db.Where("pull_requests.created_at <= ?", *filter.EndDate)
		}
		return db
	}).Where("pull_requests.state = ?", "closed").Count(&stats.MergedPRs).Error; err != nil {
		return nil, err
	}

	// Active Repos - filter by userID
	repoScope := r.db.WithContext(ctx).Model(&models.Repository{}).Where("user_id = ?", userID)
	if filter.RepoID != "" {
		repoScope = repoScope.Where("id = ?", filter.RepoID)
	}
	if err := repoScope.Count(&stats.ActiveRepos).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

func (r *gormGithubRepository) GetPullRequest(ctx context.Context, userID string, repoID string, number int) (*models.PullRequest, error) {
	var pr models.PullRequest
	query := r.db.WithContext(ctx).Preload("Repository").Where("pull_requests.repo_id = ? AND pull_requests.number = ?", repoID, number)
	// Only filter by userID if it's provided (for user-specific queries)
	if userID != "" {
		query = query.Joins("JOIN repositories ON repositories.id = pull_requests.repo_id").Where("repositories.user_id = ?", userID)
	}
	if err := query.First(&pr).Error; err != nil {
		return nil, err
	}
	return &pr, nil
}

func (r *gormGithubRepository) GetPullRequestByID(ctx context.Context, prID string) (*models.PullRequest, error) {
	var pr models.PullRequest
	if err := r.db.WithContext(ctx).Preload("Repository").Where("id = ?", prID).First(&pr).Error; err != nil {
		return nil, err
	}
	return &pr, nil
}

func (r *gormGithubRepository) GetRecentPullRequests(ctx context.Context, userID string, limit int) ([]*models.PullRequest, error) {
	var prs []*models.PullRequest
	if err := r.db.WithContext(ctx).Preload("Repository").Joins("JOIN repositories ON repositories.id = pull_requests.repo_id").Where("repositories.user_id = ?", userID).Order("pull_requests.created_at desc").Limit(limit).Find(&prs).Error; err != nil {
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
