package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"devplus-backend/internal/db"
	"devplus-backend/internal/models"
)

type WebhookService struct{}

// signature validation
func (s *WebhookService) ValidateSignature(req *http.Request, body []byte) error {
	signature := req.Header.Get("X-Hub-Signature-256")
	if signature == "" {
		return errors.New("missing signature")
	}

	secret := os.Getenv("GITHUB_WEBHOOK_SECRET")
	if secret == "" {
		// In dev/hackathon, maybe allow loose validation but logging warning
		return nil
	}

	parts := strings.SplitN(signature, "=", 2)
	if len(parts) != 2 || parts[0] != "sha256" {
		return errors.New("invalid signature format")
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expectedMAC := mac.Sum(nil)
	expectedSignature := hex.EncodeToString(expectedMAC)

	if !hmac.Equal([]byte(parts[1]), []byte(expectedSignature)) {
		return errors.New("signature mismatch")
	}

	return nil
}

func (s *WebhookService) ProcessWebhook(eventType string, payload []byte) error {
	switch eventType {
	case "pull_request":
		return s.handlePullRequest(payload)
	case "push":
		return s.handlePush(payload)
	default:
		// Ignore other events
		return nil
	}
}

func (s *WebhookService) handlePullRequest(payload []byte) error {
	var event struct {
		Action      string `json:"action"`
		Number      int    `json:"number"`
		PullRequest struct {
			ID    int64  `json:"id"`
			Title string `json:"title"`
			State string `json:"state"`
			User  struct {
				Login string `json:"login"`
				ID    int64  `json:"id"`
			} `json:"user"`
			Head struct {
				Repo struct {
					ID int64 `json:"id"`
				} `json:"repo"`
			} `json:"head"`
		} `json:"pull_request"`
		Repository struct {
			ID       int64  `json:"id"`
			Name     string `json:"name"`
			CloneURL string `json:"clone_url"`
			Owner    struct {
				Login string `json:"login"`
			} `json:"owner"`
		} `json:"repository"`
	}

	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	// Upsert Repository
	if err := s.upsertRepo(event.Repository.ID, event.Repository.Name, event.Repository.Owner.Login, event.Repository.CloneURL); err != nil {
		fmt.Printf("Error upserting repo: %v\n", err)
	}

	// Find Repo ID
	var repo models.Repository
	if err := db.DB.Where("github_repo_id = ?", event.Repository.ID).First(&repo).Error; err != nil {
		return err
	}

	// Upsert PR
	var pr models.PullRequest
	db.DB.Where("github_pr_id = ?", event.PullRequest.ID).First(&pr)

	pr.GithubPRID = event.PullRequest.ID
	pr.Number = event.Number
	pr.Title = event.PullRequest.Title
	pr.State = models.PullRequestState(event.PullRequest.State)
	pr.RepoID = repo.ID
	pr.AuthorName = event.PullRequest.User.Login

	// Ideally link AuthorID if User exists in DB
	var user models.User
	if err := db.DB.Where("github_id = ?", event.PullRequest.User.ID).First(&user).Error; err == nil {
		pr.AuthorID = user.ID
	}

	return db.DB.Save(&pr).Error
}

func (s *WebhookService) handlePush(payload []byte) error {
	var event struct {
		Ref        string `json:"ref"`
		Repository struct {
			ID       int64  `json:"id"`
			Name     string `json:"name"`
			CloneURL string `json:"clone_url"`
			Owner    struct {
				Login string `json:"login"`
			} `json:"owner"`
		} `json:"repository"`
		Commits []struct {
			ID      string `json:"id"`
			Message string `json:"message"`
			Author  struct {
				Name string `json:"name"`
			} `json:"author"`
		} `json:"commits"`
	}

	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	// Only care about main/master or development branches usually, but ingest all for now

	// Ensure repo exists
	if err := s.upsertRepo(event.Repository.ID, event.Repository.Name, event.Repository.Owner.Login, event.Repository.CloneURL); err != nil {
		fmt.Printf("Error upserting repo: %v\n", err)
	}

	var repo models.Repository
	if err := db.DB.Where("github_repo_id = ?", event.Repository.ID).First(&repo).Error; err != nil {
		return err
	}

	for _, c := range event.Commits {
		var commit models.Commit
		db.DB.Where("sha = ?", c.ID).First(&commit)
		commit.SHA = c.ID
		commit.Message = c.Message
		commit.AuthorName = c.Author.Name
		commit.RepoID = repo.ID

		db.DB.Save(&commit)
	}

	return nil
}

func (s *WebhookService) upsertRepo(ghId int64, name, owner, url string) error {
	var repo models.Repository
	res := db.DB.Where("github_repo_id = ?", ghId).First(&repo)
	repo.GithubRepoID = ghId
	repo.Name = name
	repo.Owner = owner
	repo.URL = url

	if res.Error != nil {
		return db.DB.Create(&repo).Error
	}
	return db.DB.Save(&repo).Error
}
