package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"devplus-backend/internal/interfaces"
	"devplus-backend/internal/middleware"
	"devplus-backend/internal/models"
	"devplus-backend/internal/services/ai"
)

type GithubController struct {
	service    interfaces.GithubService
	aiFactory  *ai.AIFactory
	backendURL string
}

func NewGithubController(service interfaces.GithubService, aiFactory *ai.AIFactory, backendURL string) *GithubController {
	return &GithubController{
		service:    service,
		aiFactory:  aiFactory,
		backendURL: backendURL,
	}
}

func (c *GithubController) GetRepositories(w http.ResponseWriter, r *http.Request) {
	// 2. Call Service (Fetch from DB)
	repos, err := c.service.GetRepositories(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Return JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(repos)
}

func (c *GithubController) GetRepository(w http.ResponseWriter, r *http.Request) {
	// 1. Get ID from path
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Repository ID is required", http.StatusBadRequest)
		return
	}

	// 2. Call Service
	repo, err := c.service.GetRepository(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Return JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(repo)
}

func (c *GithubController) SyncRepositories(w http.ResponseWriter, r *http.Request) {
	// 1. Get Provider Token from Context
	token, ok := r.Context().Value(middleware.GithubTokenContextKey).(string)
	if !ok || token == "" {
		http.Error(w, "GitHub token not found in context", http.StatusUnauthorized)
		return
	}

	// Get User from Context
	user, ok := r.Context().Value(middleware.UserContextKey).(*models.User)
	if !ok || user == nil {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	// 2. Call Service (Fetch from GitHub & Upsert)
	repos, err := c.service.SyncRepositories(r.Context(), user.ID, token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Return JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(repos)
}

func (c *GithubController) GetPullRequests(w http.ResponseWriter, r *http.Request) {

	// 1. Parse Path Params
	vars := mux.Vars(r)
	owner := vars["owner"]
	repo := vars["repo"]

	if owner == "" || repo == "" {
		http.Error(w, "Owner and Repo are required", http.StatusBadRequest)
		return
	}

	// 2. Call Service
	prs, err := c.service.GetPullRequests(r.Context(), owner, repo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. Return JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prs)
}

func (c *GithubController) GetMetrics(w http.ResponseWriter, r *http.Request) {
	// Parse Query Params
	query := r.URL.Query()
	repoID := query.Get("repo_id")
	startDate := query.Get("start_date")
	endDate := query.Get("end_date")

	filter := models.MetricsFilter{
		RepoID: repoID,
	}

	if startDate != "" {
		filter.StartDate = &startDate
	}
	if endDate != "" {
		filter.EndDate = &endDate
	}

	stats, err := c.service.GetMetrics(r.Context(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (c *GithubController) GetDashboardStats(w http.ResponseWriter, r *http.Request) {
	// 1. Call Service
	stats, err := c.service.GetDashboardStats(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 2. Return JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (c *GithubController) GetRecentActivity(w http.ResponseWriter, r *http.Request) {
	// 1. Call Service
	// Limit to 10 recent PRs
	prs, err := c.service.GetRecentPullRequests(r.Context(), 10)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 2. Return JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prs)
}

func (c *GithubController) GetPullRequestDetail(w http.ResponseWriter, r *http.Request) {
	// 1. Get ID and PR Number from path
	vars := mux.Vars(r)
	id := vars["id"]
	prNumberStr := vars["pr_number"]

	if id == "" || prNumberStr == "" {
		http.Error(w, "Repository ID and PR Number are required", http.StatusBadRequest)
		return
	}

	prNumber, err := strconv.Atoi(prNumberStr)
	if err != nil {
		http.Error(w, "Invalid PR Number", http.StatusBadRequest)
		return
	}

	// 2. Call Service
	pr, err := c.service.GetPullRequest(r.Context(), id, prNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Return JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pr)
}

func (c *GithubController) AnalyzePullRequest(w http.ResponseWriter, r *http.Request) {
	// 1. Parse Path Params
	vars := mux.Vars(r)
	id := vars["id"]
	prNumberStr := vars["pr_number"]

	prNumber, err := strconv.Atoi(prNumberStr)
	if err != nil {
		http.Error(w, "Invalid PR Number", http.StatusBadRequest)
		return
	}

	// 2. Fetch PR context
	pr, err := c.service.GetPullRequest(r.Context(), id, prNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Trigger AI Service
	// Always use "kestra" as it's the confirmed workflow engine
	aiService, err := c.aiFactory.GetAIService("kestra")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Construct Callback URL
	callbackURL := c.backendURL + "/api/v1/webhook/ai"

	if err := aiService.AnalyzePR(r.Context(), pr, callbackURL); err != nil {
		http.Error(w, "Failed to trigger analysis: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"status": "queued", "message": "Analysis triggered"})
}

func (c *GithubController) HandleAIWebhook(w http.ResponseWriter, r *http.Request) {
	// Parse Webhook Payload
	var payload struct {
		PRID     string `json:"pr_id"`
		Summary  string `json:"summary"`
		Decision string `json:"decision"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	if payload.PRID == "" {
		http.Error(w, "pr_id is required", http.StatusBadRequest)
		return
	}

	// Update DB
	ctx := r.Context()
	if err := c.service.UpdatePullRequestAnalysis(ctx, payload.PRID, payload.Summary, payload.Decision); err != nil {
		http.Error(w, "Failed to update PR: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *GithubController) HandleGithubWebhook(w http.ResponseWriter, r *http.Request) {
	// Simple payload structure for relevant fields
	type WebhookPayload struct {
		Action      string `json:"action"`
		PullRequest struct {
			ID     int64  `json:"id"`
			Number int64  `json:"number"`
			Title  string `json:"title"`
			State  string `json:"state"`
			User   struct {
				ID    int64  `json:"id"`
				Login string `json:"login"`
			} `json:"user"`
		} `json:"pull_request"`
		Repository struct {
			ID int64 `json:"id"`
		} `json:"repository"`
	}

	var payload WebhookPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// Only process "opened" or "synchronize" (re-analysis on update)
	if payload.Action != "opened" && payload.Action != "synchronize" {
		w.WriteHeader(http.StatusOK) // Ignore other actions
		return
	}

	ctx := r.Context()

	// 1. Find Repository by GithubRepoID
	repo, err := c.service.GetRepositoryByGithubID(ctx, payload.Repository.ID)
	if err != nil {
		// If repo not found, strictly we ignore, but maybe log warning.
		http.Error(w, "Repository not found: "+err.Error(), http.StatusNotFound)
		return
	}

	// 2. Upsert Pull Request
	pr := &models.PullRequest{
		GithubPRID: &payload.PullRequest.ID,
		Number:     &payload.PullRequest.Number,
		Title:      &payload.PullRequest.Title,
		State:      &payload.PullRequest.State,
		RepoID:     &repo.ID,
		AuthorID:   &payload.PullRequest.User.ID,
		AuthorName: &payload.PullRequest.User.Login,
	}

	if err := c.service.UpsertPullRequest(ctx, pr); err != nil {
		http.Error(w, "Failed to upsert PR: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Refetch PR to get the generated ID (UUID)
	savedPR, err := c.service.GetPullRequest(ctx, repo.ID, int(*pr.Number))
	if err != nil {
		http.Error(w, "Failed to fetch saved PR: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Trigger AI Analysis
	aiService, err := c.aiFactory.GetAIService("kestra")
	if err != nil {
		// Log error but assume success for webhook response?
		// Better to return 500 so GitHub retries if transient.
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	callbackURL := c.backendURL + "/api/v1/webhook/ai"
	if err := aiService.AnalyzePR(ctx, savedPR, callbackURL); err != nil {
		http.Error(w, "Failed to trigger analysis: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
