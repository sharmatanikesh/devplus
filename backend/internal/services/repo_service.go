package services

import (
	"devplus-backend/internal/db"
	"devplus-backend/internal/models"

	"encoding/json"
	"fmt"
	"net/http"
)

type RepoService struct{}

// ListRepos returns all repositories in the DB
func (s *RepoService) ListRepos() ([]models.Repository, error) {
	var repos []models.Repository
	if err := db.DB.Find(&repos).Error; err != nil {
		return nil, err
	}
	return repos, nil
}

// SyncRepo fetches latest PRs and Commits from GitHub for a given repo
// This is a naive implementation using public API or provided token
// Context: "manual repo sync / fetch (protected)"
func (s *RepoService) SyncRepo(repoID uint, token string) error {
	var repo models.Repository
	if err := db.DB.First(&repo, repoID).Error; err != nil {
		return err
	}

	// Fetch PRs
	// https://api.github.com/repos/OWNER/REPO/pulls
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls?state=all", repo.Owner, repo.Name)
	req, _ := http.NewRequest("GET", url, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("github api error: %s", resp.Status)
	}

	var pulls []struct {
		ID     int64  `json:"id"`
		Number int    `json:"number"`
		Title  string `json:"title"`
		State  string `json:"state"`
		User   struct {
			Login string `json:"login"`
			ID    int64  `json:"id"`
		} `json:"user"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&pulls); err != nil {
		return err
	}

	for _, p := range pulls {
		var pr models.PullRequest
		db.DB.Where("github_pr_id = ?", p.ID).First(&pr)
		pr.GithubPRID = p.ID
		pr.Number = p.Number
		pr.Title = p.Title
		pr.State = models.PullRequestState(p.State)
		pr.RepoID = repo.ID
		pr.AuthorName = p.User.Login
		db.DB.Save(&pr)
	}

	return nil
}

func (s *RepoService) ListPRs(repoID uint) ([]models.PullRequest, error) {
	var prs []models.PullRequest
	if err := db.DB.Where("repo_id = ?", repoID).Find(&prs).Error; err != nil {
		return nil, err
	}
	return prs, nil
}

func (s *RepoService) GetPR(repoID uint, prNumber int) (*models.PullRequest, error) {
	var pr models.PullRequest
	if err := db.DB.Preload("Repository").Where("repo_id = ? AND number = ?", repoID, prNumber).First(&pr).Error; err != nil {
		return nil, err
	}
	return &pr, nil
}
