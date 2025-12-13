package models

type DashboardStats struct {
	TotalPRs    int64 `json:"total_prs"`
	OpenPRs     int64 `json:"open_prs"`
	MergedPRs   int64 `json:"merged_prs"`
	ActiveRepos int64 `json:"active_repos"`
}

type MetricsFilter struct {
	RepoID    string
	StartDate *string // RFC3339
	EndDate   *string // RFC3339
}
