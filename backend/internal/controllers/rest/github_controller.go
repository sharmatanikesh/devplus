package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"devplus-backend/internal/interfaces"
	"devplus-backend/internal/middleware"
	"devplus-backend/internal/models"
)

type GithubController struct {
	service interfaces.GithubService
}

func NewGithubController(service interfaces.GithubService) *GithubController {
	return &GithubController{
		service: service,
	}
}

func (c *GithubController) GetRepositories(w http.ResponseWriter, r *http.Request) {
	// 1. Get User ID from context
	userVal, ok := r.Context().Value(middleware.UserContextKey).(models.User)
	if !ok {
		http.Error(w, "Unauthorized: User not found in context", http.StatusUnauthorized)
		return
	}

	// 2. Call Service (Fetch from DB)
	repos, err := c.service.GetRepositories(r.Context(), userVal.ID)
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
	// Get User from Context
	userVal, ok := r.Context().Value(middleware.UserContextKey).(models.User)
	if !ok {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}
	user := &userVal

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

func (c *GithubController) SyncRepository(w http.ResponseWriter, r *http.Request) {
	// 1. Get Repo ID from Path
	vars := mux.Vars(r)
	id := vars["id"]

	log.Info().Str("repo_id", id).Msg("[SyncRepository] Request received")

	if id == "" {
		log.Error().Msg("[SyncRepository] Error: Repository ID is empty")
		http.Error(w, "Repository ID is required", http.StatusBadRequest)
		return
	}

	// 2. Get Token
	token, ok := r.Context().Value(middleware.GithubTokenContextKey).(string)
	if !ok || token == "" {
		log.Error().Msg("[SyncRepository] Error: Token not found in context")
		http.Error(w, "GitHub token not found in context", http.StatusUnauthorized)
		return
	}
	log.Info().Msg("[SyncRepository] Token found, triggering service...")

	// 3. Call Service to Sync PRs
	prs, err := c.service.SyncPullRequests(r.Context(), id, token)
	if err != nil {
		log.Error().Err(err).Msg("[SyncRepository] Service error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info().Int("count", len(prs)).Msg("[SyncRepository] Success. Returning PRs")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prs)
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

	// 2. Trigger Analysis via Service
	if err := c.service.AnalyzePullRequest(r.Context(), id, prNumber); err != nil {
		http.Error(w, "Failed to trigger analysis: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"status": "queued", "message": "Analysis triggered"})
}

func (c *GithubController) HandleAIWebhook(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("[HandleAIWebhook] Received webhook callback from Kestra")

	// Parse Webhook Payload
	var payload struct {
		PRID        string `json:"pr_id"`
		Summary     string `json:"summary"`
		Decision    string `json:"decision"`
		RawAnalysis string `json:"raw_analysis"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Error().Err(err).Msg("[HandleAIWebhook] Failed to decode payload")
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	log.Info().Str("pr_id", payload.PRID).Str("raw_analysis_length", strconv.Itoa(len(payload.RawAnalysis))).Msg("[HandleAIWebhook] Payload decoded")

	if payload.PRID == "" {
		log.Error().Msg("[HandleAIWebhook] pr_id is missing")
		http.Error(w, "pr_id is required", http.StatusBadRequest)
		return
	}

	// If raw_analysis is provided, parse it to extract summary and decision
	if payload.RawAnalysis != "" {
		log.Info().Msg("[HandleAIWebhook] Parsing raw_analysis")
		// Extract JSON from markdown code blocks
		// Format: ```json\n{...}\n```
		start := strings.Index(payload.RawAnalysis, "```json")
		end := strings.LastIndex(payload.RawAnalysis, "```")

		if start != -1 && end != -1 && end > start {
			jsonStr := payload.RawAnalysis[start+7 : end] // +7 to skip "```json\n"
			jsonStr = strings.TrimSpace(jsonStr)

			var analysis struct {
				Summary  string `json:"summary"`
				Decision string `json:"decision"`
			}

			if err := json.Unmarshal([]byte(jsonStr), &analysis); err == nil {
				payload.Summary = analysis.Summary
				payload.Decision = analysis.Decision
				log.Info().Str("summary", payload.Summary).Str("decision", payload.Decision).Msg("[HandleAIWebhook] Parsed analysis")
			} else {
				log.Error().Err(err).Msg("[HandleAIWebhook] Failed to parse JSON from raw_analysis")
			}
		} else {
			log.Warn().Msg("[HandleAIWebhook] Could not find JSON code blocks in raw_analysis")
		}
	}

	// Update DB
	ctx := r.Context()
	log.Info().Str("pr_id", payload.PRID).Msg("[HandleAIWebhook] Updating PR analysis in database")
	if err := c.service.UpdatePullRequestAnalysis(ctx, payload.PRID, payload.Summary, payload.Decision); err != nil {
		log.Error().Err(err).Msg("[HandleAIWebhook] Failed to update PR")
		http.Error(w, "Failed to update PR: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info().Msg("[HandleAIWebhook] Successfully updated PR analysis")
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

	// 3. Trigger AI Analysis
	if err := c.service.AnalyzePullRequest(ctx, repo.ID, int(*pr.Number)); err != nil {
		http.Error(w, "Failed to trigger analysis: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *GithubController) AnalyzeRepository(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repoID := vars["id"]

	if repoID == "" {
		http.Error(w, "Repo ID is required", http.StatusBadRequest)
		return
	}

	// 1. Trigger Analysis via Service
	if err := c.service.AnalyzeRepository(r.Context(), repoID); err != nil {
		http.Error(w, "Failed to trigger analysis: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 2. Respond
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "analysis_triggered"})
}

// StreamRepositoryAnalysis handles SSE for real-time analysis updates
func (c *GithubController) StreamRepositoryAnalysis(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repoID := vars["id"]

	if repoID == "" {
		http.Error(w, "Repo ID is required", http.StatusBadRequest)
		return
	}

	// Set headers for SSE with explicit CORS
	origin := r.Header.Get("Origin")
	if origin == "" {
		origin = "http://localhost:3000"
	}
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	// Create SSE client
	client := &SSEClient{
		RepoID:  repoID,
		Channel: make(chan string, 10),
	}

	// Register client
	GlobalSSEManager.AddClient(repoID, client)
	defer GlobalSSEManager.RemoveClient(repoID, client)

	log.Info().Str("repo_id", repoID).Msg("[StreamRepositoryAnalysis] Client connected")

	// Send initial connection message
	fmt.Fprintf(w, "data: %s\n\n", `{"status":"connected"}`)
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	// Listen for messages or client disconnect
	for {
		select {
		case msg := <-client.Channel:
			fmt.Fprintf(w, "%s", msg)
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		case <-r.Context().Done():
			log.Info().Str("repo_id", repoID).Msg("[StreamRepositoryAnalysis] Client disconnected")
			return
		}
	}
}

func (c *GithubController) HandleRepoAIWebhook(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("[HandleRepoAIWebhook] Received webhook callback from Kestra")

	var payload struct {
		RepoID      string `json:"repo_id"`
		Summary     string `json:"summary"`
		RawAnalysis string `json:"raw_analysis"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Error().Err(err).Msg("[HandleRepoAIWebhook] Failed to decode payload")
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	log.Info().Str("repo_id", payload.RepoID).Int("raw_analysis_length", len(payload.RawAnalysis)).Msg("[HandleRepoAIWebhook] Payload decoded")

	if payload.RepoID == "" {
		log.Error().Msg("[HandleRepoAIWebhook] repo_id is missing")
		http.Error(w, "repo_id is required", http.StatusBadRequest)
		return
	}

	// If raw_analysis is provided, parse it to extract summary
	if payload.RawAnalysis != "" {
		log.Info().Msg("[HandleRepoAIWebhook] Parsing raw_analysis")
		// Try to extract JSON from markdown code blocks first
		// Format: ```json\n{...}\n```
		start := strings.Index(payload.RawAnalysis, "```json")
		end := strings.LastIndex(payload.RawAnalysis, "```")

		if start != -1 && end != -1 && end > start {
			jsonStr := payload.RawAnalysis[start+7 : end] // +7 to skip "```json\n"
			jsonStr = strings.TrimSpace(jsonStr)

			var analysis struct {
				Summary     string   `json:"summary"`
				Suggestions []string `json:"suggestions"`
			}

			if err := json.Unmarshal([]byte(jsonStr), &analysis); err == nil {
				payload.Summary = analysis.Summary
				log.Info().Str("summary", payload.Summary).Msg("[HandleRepoAIWebhook] Parsed analysis from JSON")
			} else {
				log.Warn().Err(err).Msg("[HandleRepoAIWebhook] Failed to parse JSON from raw_analysis, using raw content")
				// If JSON parsing fails, use the entire raw_analysis as summary
				payload.Summary = payload.RawAnalysis
			}
		} else {
			log.Warn().Msg("[HandleRepoAIWebhook] Could not find JSON code blocks in raw_analysis, using raw content")
			// If no code blocks found, use the entire raw_analysis as summary
			payload.Summary = payload.RawAnalysis
		}
	}

	if payload.Summary == "" {
		log.Error().Msg("[HandleRepoAIWebhook] No summary could be extracted")
		http.Error(w, "summary or raw_analysis is required", http.StatusBadRequest)
		return
	}

	// Update DB
	log.Info().Str("repo_id", payload.RepoID).Msg("[HandleRepoAIWebhook] Updating repository analysis in database")
	ctx := context.Background()
	if err := c.service.UpdateRepositoryAnalysis(ctx, payload.RepoID, payload.Summary); err != nil {
		log.Error().Err(err).Msg("[HandleRepoAIWebhook] Failed to update analysis")
		http.Error(w, "Failed to update analysis: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info().Msg("[HandleRepoAIWebhook] Successfully updated repository analysis")
	
	// Notify all connected SSE clients
	notificationData := map[string]interface{}{
		"status":     "completed",
		"ai_summary": payload.Summary,
		"repo_id":    payload.RepoID,
	}
	notificationJSON, _ := json.Marshal(notificationData)
	GlobalSSEManager.NotifyClients(payload.RepoID, FormatSSEMessage(string(notificationJSON)))
	
	log.Info().Str("repo_id", payload.RepoID).Int("clients", GlobalSSEManager.GetClientCount(payload.RepoID)).Msg("[HandleRepoAIWebhook] Notified SSE clients")
	
	w.WriteHeader(http.StatusOK)
}
