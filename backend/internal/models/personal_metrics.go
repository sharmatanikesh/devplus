package models

// PersonalMetrics represents personalized GitHub statistics for a user
type PersonalMetrics struct {
	TotalCommits       int64                   `json:"total_commits"`
	TotalContributions int64                   `json:"total_contributions"`
	PullRequestCount   int64                   `json:"pull_request_count"`
	IssuesOpened       int64                   `json:"issues_opened"`
	TimeRange          string                  `json:"time_range"` // e.g., "Last 90 days"
	LanguageStats      []LanguageStat          `json:"language_stats"`
	ContributionByRepo []RepositoryContribution `json:"contribution_by_repo"`
	RecentActivity     []Activity              `json:"recent_activity"`
	TopRepos           []TopRepository         `json:"top_repos"`
}

// LanguageStat represents programming language usage statistics
type LanguageStat struct {
	Language string  `json:"language"`
	Count    int64   `json:"count"`
	Bytes    int64   `json:"bytes"`
	Percent  float64 `json:"percent"`
}

// RepositoryContribution represents contributions to a specific repository
type RepositoryContribution struct {
	RepoName       string `json:"repo_name"`
	RepoOwner      string `json:"repo_owner"`
	CommitCount    int64  `json:"commit_count"`
	AddedLines     int64  `json:"added_lines"`
	DeletedLines   int64  `json:"deleted_lines"`
	PullRequests   int64  `json:"pull_requests"`
	LastCommitDate string `json:"last_commit_date"`
}

// Activity represents a recent GitHub activity
type Activity struct {
	Type        string `json:"type"` // commit, pr, issue, etc.
	RepoName    string `json:"repo_name"`
	Message     string `json:"message"`
	Date        string `json:"date"`
	URL         string `json:"url"`
	Additions   int64  `json:"additions,omitempty"`
	Deletions   int64  `json:"deletions,omitempty"`
}

// TopRepository represents repository stats for top contributed repos
type TopRepository struct {
	Name         string `json:"name"`
	Owner        string `json:"owner"`
	Stars        int64  `json:"stars"`
	Forks        int64  `json:"forks"`
	Commits      int64  `json:"commits"`
	Language     string `json:"language"`
	Description  string `json:"description"`
}
