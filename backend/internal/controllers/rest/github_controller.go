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
	// 1. Get Provider Token from Context
	token, ok := r.Context().Value(middleware.GithubTokenContextKey).(string)
	if !ok || token == "" {
		http.Error(w, "GitHub token not found in context", http.StatusUnauthorized)
		return
	}

	// 2. Call Service
	repos, err := c.service.FetchUserRepositories(r.Context(), token)
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
