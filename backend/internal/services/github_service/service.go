package github_service

import (
	"context"
	"fmt"
	"time"

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

	// 1. Get Repo to find owner/name
	repo, err := s.repo.GetRepository(ctx, "", repoID)
	if err != nil {
		log.Error().Err(err).Msg("[Service.SyncPullRequests] Failed to get repo")
		return nil, err
	}
	log.Info().Str("owner", repo.Owner).Str("repo", repo.Name).Msg("[Service.SyncPullRequests] Found Repo")

	// 2. Get existing PRs from database for this repo
	existingPRs, err := s.repo.GetPullRequestsByRepoID(ctx, repoID)
	if err != nil {
		log.Error().Err(err).Msg("[Service.SyncPullRequests] Failed to get existing PRs")
		return nil, err
	}
	
	// Create map of existing open PRs by number
	existingOpenPRs := make(map[int64]bool)
	for _, pr := range existingPRs {
		if pr.State != nil && *pr.State == "open" {
			existingOpenPRs[*pr.Number] = true
		}
	}
	log.Info().Int("existing_open_count", len(existingOpenPRs)).Msg("[Service.SyncPullRequests] Found existing open PRs in DB")

	// 3. Setup GitHub Client
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// 4. Fetch open PRs from GitHub
	prOpt := &github.PullRequestListOptions{
		State:       "open",
		ListOptions: github.ListOptions{PerPage: 100},
	}
	log.Info().Msg("[Service.SyncPullRequests] Fetching open PRs from GitHub...")
	openPRsFromGitHub, _, err := client.PullRequests.List(ctx, repo.Owner, repo.Name, prOpt)
	if err != nil {
		log.Error().Err(err).Msg("[Service.SyncPullRequests] GitHub API Error")
		return nil, err
	}
	log.Info().Int("count", len(openPRsFromGitHub)).Msg("[Service.SyncPullRequests] Fetched open PRs from GitHub")

	// Create map of open PRs from GitHub
	openPRNumbers := make(map[int64]bool)
	for _, pr := range openPRsFromGitHub {
		openPRNumbers[int64(pr.GetNumber())] = true
	}

	// 5. Check for PRs that were open in DB but not in GitHub's open list
	// These PRs have been closed or merged
	for prNumber := range existingOpenPRs {
		if !openPRNumbers[prNumber] {
			log.Info().Int64("pr_number", prNumber).Msg("[Service.SyncPullRequests] PR was open in DB but not in GitHub open list, checking status...")
			
			// Fetch individual PR to check if merged or just closed
			pr, _, err := client.PullRequests.Get(ctx, repo.Owner, repo.Name, int(prNumber))
			if err != nil {
				log.Error().Err(err).Int64("pr_number", prNumber).Msg("[Service.SyncPullRequests] Failed to fetch PR details")
				continue
			}
			
			state := pr.GetState()
			if pr.GetMerged() {
				state = "merged"
				log.Info().Int64("pr_number", prNumber).Msg("[Service.SyncPullRequests] PR was merged")
			} else {
				log.Info().Int64("pr_number", prNumber).Msg("[Service.SyncPullRequests] PR was closed without merging")
			}
			
			// Update the PR in database
			prModel := &models.PullRequest{
				GithubPRID: pr.ID,
				Number:     github.Int64(int64(pr.GetNumber())),
				Title:      github.String(pr.GetTitle()),
				State:      github.String(state),
				RepoID:     &repo.ID,
				AuthorID:   pr.GetUser().ID,
				AuthorName: github.String(pr.GetUser().GetLogin()),
				CreatedAt:  &pr.CreatedAt.Time,
				UpdatedAt:  &pr.UpdatedAt.Time,
			}
			if err := s.repo.UpsertPullRequest(ctx, prModel); err != nil {
				log.Error().Int("pr_number", pr.GetNumber()).Err(err).Msg("[Service.SyncPullRequests] Upsert Error")
			}
		}
	}

	// 6. Upsert all open PRs from GitHub
	var syncedPRs []*models.PullRequest
	for _, pr := range openPRsFromGitHub {
		log.Info().Int("pr_number", pr.GetNumber()).Str("title", pr.GetTitle()).Msg("[Service.SyncPullRequests] Processing open PR")
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

func (s *GithubService) GetRepositories(ctx context.Context, userID string) ([]*models.Repository, error) {
	return s.repo.GetRepositories(ctx, userID)
}

func (s *GithubService) GetRepository(ctx context.Context, userID string, id string) (*models.Repository, error) {
	return s.repo.GetRepository(ctx, userID, id)
}

func (s *GithubService) GetPullRequests(ctx context.Context, userID string, owner, repo string) ([]*models.PullRequest, error) {
	return s.repo.GetPullRequests(ctx, userID, owner, repo)
}

func (s *GithubService) GetPullRequest(ctx context.Context, userID string, repoID string, number int) (*models.PullRequest, error) {
	return s.repo.GetPullRequest(ctx, userID, repoID, number)
}

func (s *GithubService) GetPullRequestByID(ctx context.Context, prID string) (*models.PullRequest, error) {
	return s.repo.GetPullRequestByID(ctx, prID)
}

func (s *GithubService) GetMetrics(ctx context.Context, userID string, filter models.MetricsFilter) (*models.DashboardStats, error) {
	return s.repo.GetMetrics(ctx, userID, filter)
}

func (s *GithubService) GetDashboardStats(ctx context.Context, userID string) (*models.DashboardStats, error) {
	return s.repo.GetDashboardStats(ctx, userID)
}

func (s *GithubService) GetRecentPullRequests(ctx context.Context, userID string, limit int) ([]*models.PullRequest, error) {
	return s.repo.GetRecentPullRequests(ctx, userID, limit)
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
	// Using empty userID for internal operations - repo lookup by ID
	repo, err := s.repo.GetRepository(ctx, "", repoID)
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
	// Using empty userID for internal operations - PR lookup by repoID and number
	pr, err := s.repo.GetPullRequest(ctx, "", repoID, prNumber)
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

func (s *GithubService) GetPersonalMetrics(ctx context.Context, userID string, token string, username string, days int) (*models.PersonalMetrics, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	metrics := &models.PersonalMetrics{
		LanguageStats:      []models.LanguageStat{},
		ContributionByRepo: []models.RepositoryContribution{},
		RecentActivity:     []models.Activity{},
		TopRepos:           []models.TopRepository{},
		TimeRange:          fmt.Sprintf("Last %d days", days),
	}

	// Define time range based on days parameter
	since := time.Now().AddDate(0, 0, -days)
	sinceStr := since.Format("2006-01-02")

	log.Info().Str("since", sinceStr).Str("username", username).Msg("Fetching personal metrics")

	// Use GitHub Events API for recent activity (much more efficient)
	events, _, err := client.Activity.ListEventsPerformedByUser(ctx, username, false, &github.ListOptions{PerPage: 100})
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch user events")
		return nil, err
	}

	repoMap := make(map[string]bool)
	languageMap := make(map[string]int64)
	var totalBytes int64
	repoContributions := make(map[string]*models.RepositoryContribution)
	repoStatsMap := make(map[string]*models.TopRepository)

	// Process events to extract meaningful data
	for _, event := range events {
		if event.Repo == nil {
			continue
		}
		repoName := event.Repo.GetName()
		repoMap[repoName] = true

		// Extract activities from different event types
		switch *event.Type {
		case "PushEvent":
			if payload, ok := event.Payload().(*github.PushEvent); ok {
				for _, commit := range payload.Commits {
					metrics.RecentActivity = append(metrics.RecentActivity, models.Activity{
						Type:     "commit",
						RepoName: repoName,
						Message:  commit.GetMessage(),
						Date:     event.CreatedAt.Format("2006-01-02 15:04"),
						URL:      commit.GetURL(),
					})
				}
				
				// Track repo contribution
				if _, exists := repoContributions[repoName]; !exists {
					parts := splitRepoName(repoName)
					repoContributions[repoName] = &models.RepositoryContribution{
						RepoName:  parts[1],
						RepoOwner: parts[0],
					}
				}
				repoContributions[repoName].CommitCount += int64(len(payload.Commits))
			}

		case "PullRequestEvent":
			if payload, ok := event.Payload().(*github.PullRequestEvent); ok {
				if payload.PullRequest.GetUser().GetLogin() == username {
					metrics.RecentActivity = append(metrics.RecentActivity, models.Activity{
						Type:      "pr",
						RepoName:  repoName,
						Message:   payload.PullRequest.GetTitle(),
						Date:      event.CreatedAt.Format("2006-01-02 15:04"),
						URL:       payload.PullRequest.GetHTMLURL(),
						Additions: int64(payload.PullRequest.GetAdditions()),
						Deletions: int64(payload.PullRequest.GetDeletions()),
					})
					
					if _, exists := repoContributions[repoName]; !exists {
						parts := splitRepoName(repoName)
						repoContributions[repoName] = &models.RepositoryContribution{
							RepoName:  parts[1],
							RepoOwner: parts[0],
						}
					}
					repoContributions[repoName].PullRequests++
					repoContributions[repoName].AddedLines += int64(payload.PullRequest.GetAdditions())
					repoContributions[repoName].DeletedLines += int64(payload.PullRequest.GetDeletions())
				}
			}

		case "IssuesEvent":
			if payload, ok := event.Payload().(*github.IssuesEvent); ok {
				if payload.GetAction() == "opened" {
					metrics.IssuesOpened++
				}
			}
		}
	}

	// Use Search API for commits in the time range (more efficient than listing per repo)
	query := fmt.Sprintf("author:%s author-date:>%s", username, sinceStr)
	commitResults, _, err := client.Search.Commits(ctx, query, &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	})
	if err != nil {
		log.Warn().Err(err).Msg("Failed to search commits, falling back to events data")
	} else {
		metrics.TotalCommits = int64(commitResults.GetTotal())
		
		// Process commit results for language stats
		processedRepos := make(map[string]bool)
		for _, commit := range commitResults.Commits {
			if commit.Repository == nil {
				continue
			}
			repoFullName := commit.Repository.GetFullName()
			
			// Get language stats for repos we haven't processed yet (limit to 30 most active repos)
			if !processedRepos[repoFullName] && len(processedRepos) < 30 {
				processedRepos[repoFullName] = true
				parts := splitRepoName(repoFullName)
				
				languages, _, err := client.Repositories.ListLanguages(ctx, parts[0], parts[1])
				if err != nil {
					log.Debug().Err(err).Str("repo", repoFullName).Msg("Failed to fetch languages")
					continue
				}
				
				for lang, bytes := range languages {
					languageMap[lang] += int64(bytes)
					totalBytes += int64(bytes)
				}
				
				// Get repo details for top repos
				repo, _, err := client.Repositories.Get(ctx, parts[0], parts[1])
				if err == nil {
					repoStatsMap[repoFullName] = &models.TopRepository{
						Name:        repo.GetName(),
						Owner:       repo.GetOwner().GetLogin(),
						Stars:       int64(repo.GetStargazersCount()),
						Forks:       int64(repo.GetForksCount()),
						Language:    repo.GetLanguage(),
						Description: repo.GetDescription(),
					}
				}
			}
		}
	}

	// Use Search API for PRs in the time range
	prQuery := fmt.Sprintf("author:%s is:pr created:>%s", username, sinceStr)
	prResults, _, err := client.Search.Issues(ctx, prQuery, &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	})
	if err != nil {
		log.Warn().Err(err).Msg("Failed to search PRs")
	} else {
		metrics.PullRequestCount = int64(prResults.GetTotal())
	}

	// Convert language map to sorted stats
	for lang, bytes := range languageMap {
		if totalBytes > 0 {
			percent := float64(bytes) / float64(totalBytes) * 100
			metrics.LanguageStats = append(metrics.LanguageStats, models.LanguageStat{
				Language: lang,
				Bytes:    bytes,
				Percent:  percent,
			})
		}
	}

	// Convert maps to slices
	for _, contrib := range repoContributions {
		if contrib.CommitCount > 0 || contrib.PullRequests > 0 {
			metrics.ContributionByRepo = append(metrics.ContributionByRepo, *contrib)
		}
	}

	for _, repo := range repoStatsMap {
		metrics.TopRepos = append(metrics.TopRepos, *repo)
	}

	// Sort language stats by bytes (descending)
	for i := 0; i < len(metrics.LanguageStats); i++ {
		for j := i + 1; j < len(metrics.LanguageStats); j++ {
			if metrics.LanguageStats[j].Bytes > metrics.LanguageStats[i].Bytes {
				metrics.LanguageStats[i], metrics.LanguageStats[j] = metrics.LanguageStats[j], metrics.LanguageStats[i]
			}
		}
	}

	// Sort contribution by repo by commit count (descending)
	for i := 0; i < len(metrics.ContributionByRepo); i++ {
		for j := i + 1; j < len(metrics.ContributionByRepo); j++ {
			if metrics.ContributionByRepo[j].CommitCount > metrics.ContributionByRepo[i].CommitCount {
				metrics.ContributionByRepo[i], metrics.ContributionByRepo[j] = metrics.ContributionByRepo[j], metrics.ContributionByRepo[i]
			}
		}
	}

	// Sort top repos by stars
	for i := 0; i < len(metrics.TopRepos); i++ {
		for j := i + 1; j < len(metrics.TopRepos); j++ {
			if metrics.TopRepos[j].Stars > metrics.TopRepos[i].Stars {
				metrics.TopRepos[i], metrics.TopRepos[j] = metrics.TopRepos[j], metrics.TopRepos[i]
			}
		}
	}

	// Sort recent activity by date (most recent first)
	for i := 0; i < len(metrics.RecentActivity); i++ {
		for j := i + 1; j < len(metrics.RecentActivity); j++ {
			if metrics.RecentActivity[j].Date > metrics.RecentActivity[i].Date {
				metrics.RecentActivity[i], metrics.RecentActivity[j] = metrics.RecentActivity[j], metrics.RecentActivity[i]
			}
		}
	}

	// Limit arrays to top results
	if len(metrics.LanguageStats) > 10 {
		metrics.LanguageStats = metrics.LanguageStats[:10]
	}
	if len(metrics.ContributionByRepo) > 20 {
		metrics.ContributionByRepo = metrics.ContributionByRepo[:20]
	}
	if len(metrics.TopRepos) > 10 {
		metrics.TopRepos = metrics.TopRepos[:10]
	}
	if len(metrics.RecentActivity) > 30 {
		metrics.RecentActivity = metrics.RecentActivity[:30]
	}

	// Set aggregate metrics
	metrics.TotalContributions = metrics.TotalCommits + metrics.PullRequestCount

	log.Info().
		Int64("total_commits", metrics.TotalCommits).
		Int64("total_prs", metrics.PullRequestCount).
		Int("languages", len(metrics.LanguageStats)).
		Int("repos", len(metrics.ContributionByRepo)).
		Msg("Successfully computed personal metrics")

	return metrics, nil
}

// Helper function to split "owner/repo" format
func splitRepoName(fullName string) []string {
	parts := make([]string, 2)
	idx := 0
	for i, char := range fullName {
		if char == '/' {
			parts[0] = fullName[:i]
			parts[1] = fullName[i+1:]
			return parts
		}
	}
	if idx == 0 {
		parts[0] = fullName
		parts[1] = ""
	}
	return parts
}
