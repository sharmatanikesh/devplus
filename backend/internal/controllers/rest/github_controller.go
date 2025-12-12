package rest

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"devplus-backend/internal/interfaces"
	"devplus-backend/internal/middleware"
)

type GithubController struct {
	service interfaces.GithubService
}

func NewGithubController(service interfaces.GithubService) *GithubController {
	return &GithubController{service: service}
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

func (c *GithubController) SyncRepositories(w http.ResponseWriter, r *http.Request) {
	// 1. Get Provider Token from Context
	token, ok := r.Context().Value(middleware.GithubTokenContextKey).(string)
	if !ok || token == "" {
		http.Error(w, "GitHub token not found in context", http.StatusUnauthorized)
		return
	}

	// 2. Call Service (Fetch from GitHub & Upsert)
	repos, err := c.service.SyncRepositories(r.Context(), token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Return JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(repos)
}

func (c *GithubController) GetPullRequests(w http.ResponseWriter, r *http.Request) {
	// 1. Get Provider Token from Context
	token, ok := r.Context().Value(middleware.GithubTokenContextKey).(string)
	if !ok || token == "" {
		http.Error(w, "GitHub token not found in context", http.StatusUnauthorized)
		return
	}

	// 2. Parse Path Params
	vars := mux.Vars(r)
	owner := vars["owner"]
	repo := vars["repo"]

	if owner == "" || repo == "" {
		http.Error(w, "Owner and Repo are required", http.StatusBadRequest)
		return
	}

	// 3. Call Service
	prs, err := c.service.FetchRepositoryPullRequests(r.Context(), token, owner, repo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. Return JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prs)
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
