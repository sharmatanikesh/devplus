package github_service

import (
	"context"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"

	"devplus-backend/internal/models"
)

type GithubService struct {
}

func NewGithubService() *GithubService {
	return &GithubService{}
}

func (s *GithubService) FetchUserRepositories(ctx context.Context, token string) ([]*models.Repository, error) {
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
		result = append(result, &models.Repository{
			GithubRepoID: repo.ID,
			Name:         repo.Name,
			Owner:        repo.Owner.Login,
			URL:          repo.HTMLURL,
			// InstallationID might need a different call or context if using GitHub Apps,
			// but for OAuth app flow it's less direct per repo unless installed.
			// We'll skip InstallationID for now as it wasn't requested explicitly for this step.
		})
	}

	return result, nil
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

func interfaceToInt64(i interface{}) *int64 {
	x, ok := i.(int)
	if ok {
		v := int64(x)
		return &v
	}
	// Handle other types if necessary or return nil
	return nil
}
