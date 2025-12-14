package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
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
	// 1. Get User ID from context
	userVal, ok := r.Context().Value(middleware.UserContextKey).(models.User)
	if !ok {
		http.Error(w, "Unauthorized: User not found in context", http.StatusUnauthorized)
		return
	}

	// 2. Get ID from path
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Repository ID is required", http.StatusBadRequest)
		return
	}

	// 3. Call Service
	repo, err := c.service.GetRepository(r.Context(), userVal.ID, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. Return JSON
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
	// 1. Get User ID from context
	userVal, ok := r.Context().Value(middleware.UserContextKey).(models.User)
	if !ok {
		http.Error(w, "Unauthorized: User not found in context", http.StatusUnauthorized)
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
	prs, err := c.service.GetPullRequests(r.Context(), userVal.ID, owner, repo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. Return JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prs)
}

func (c *GithubController) GetMetrics(w http.ResponseWriter, r *http.Request) {
	// 1. Get User ID from context
	userVal, ok := r.Context().Value(middleware.UserContextKey).(models.User)
	if !ok {
		http.Error(w, "Unauthorized: User not found in context", http.StatusUnauthorized)
		return
	}

	// 2. Parse Query Params
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

	// 3. Call Service
	stats, err := c.service.GetMetrics(r.Context(), userVal.ID, filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. Return JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (c *GithubController) GetPersonalMetrics(w http.ResponseWriter, r *http.Request) {
	// 1. Get User from context
	userVal, ok := r.Context().Value(middleware.UserContextKey).(models.User)
	if !ok {
		http.Error(w, "Unauthorized: User not found in context", http.StatusUnauthorized)
		return
	}

	// 2. Get GitHub token from context
	token, ok := r.Context().Value(middleware.GithubTokenContextKey).(string)
	if !ok || token == "" {
		http.Error(w, "GitHub token not found in context", http.StatusUnauthorized)
		return
	}

	// 3. Get optional days parameter (default 90)
	days := 90
	if daysParam := r.URL.Query().Get("days"); daysParam != "" {
		if parsedDays, err := strconv.Atoi(daysParam); err == nil && parsedDays > 0 && parsedDays <= 365 {
			days = parsedDays
		}
	}

	// 4. Call Service to fetch personal metrics from GitHub
	metrics, err := c.service.GetPersonalMetrics(r.Context(), userVal.ID, token, userVal.Username, days)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch personal metrics")
		http.Error(w, "Failed to fetch personal metrics: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 5. Return JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

func (c *GithubController) GetDashboardStats(w http.ResponseWriter, r *http.Request) {
	// 1. Get User ID from context
	userVal, ok := r.Context().Value(middleware.UserContextKey).(models.User)
	if !ok {
		http.Error(w, "Unauthorized: User not found in context", http.StatusUnauthorized)
		return
	}

	// 2. Call Service
	stats, err := c.service.GetDashboardStats(r.Context(), userVal.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Return JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (c *GithubController) GetRecentActivity(w http.ResponseWriter, r *http.Request) {
	// 1. Get User ID from context
	userVal, ok := r.Context().Value(middleware.UserContextKey).(models.User)
	if !ok {
		http.Error(w, "Unauthorized: User not found in context", http.StatusUnauthorized)
		return
	}

	// 2. Call Service - Limit to 10 recent PRs
	prs, err := c.service.GetRecentPullRequests(r.Context(), userVal.ID, 10)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Return JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prs)
}

func (c *GithubController) GetPullRequestDetail(w http.ResponseWriter, r *http.Request) {
	// 1. Get User ID from context
	userVal, ok := r.Context().Value(middleware.UserContextKey).(models.User)
	if !ok {
		http.Error(w, "Unauthorized: User not found in context", http.StatusUnauthorized)
		return
	}

	// 2. Get ID and PR Number from path
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

	// 3. Call Service
	pr, err := c.service.GetPullRequest(r.Context(), userVal.ID, id, prNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. Return JSON
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

// StreamPullRequestAnalysis handles SSE for real-time PR analysis updates
func (c *GithubController) StreamPullRequestAnalysis(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repoID := vars["id"]
	prNumberStr := vars["pr_number"]

	if repoID == "" || prNumberStr == "" {
		http.Error(w, "Repo ID and PR number are required", http.StatusBadRequest)
		return
	}

	// Create unique key for this PR
	prKey := fmt.Sprintf("%s:%s", repoID, prNumberStr)

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
		RepoID:  prKey,
		Channel: make(chan string, 10),
	}

	// Register client
	GlobalSSEManager.AddClient(prKey, client)
	defer GlobalSSEManager.RemoveClient(prKey, client)

	log.Info().Str("pr_key", prKey).Msg("[StreamPullRequestAnalysis] Client connected")

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
			log.Info().Str("pr_key", prKey).Msg("[StreamPullRequestAnalysis] Client disconnected")
			return
		}
	}
}

func (c *GithubController) HandleAIWebhook(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("[HandleAIWebhook] Received webhook callback from Kestra")

	var payload struct {
		PRID        string `json:"pr_id"`
		RawAnalysis string `json:"raw_analysis"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Error().Err(err).Msg("[HandleAIWebhook] Failed to decode payload")
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	if payload.PRID == "" {
		log.Error().Msg("[HandleAIWebhook] pr_id is missing")
		http.Error(w, "pr_id is required", http.StatusBadRequest)
		return
	}

	if payload.RawAnalysis == "" {
		log.Error().Msg("[HandleAIWebhook] raw_analysis is missing")
		http.Error(w, "raw_analysis is required", http.StatusBadRequest)
		return
	}

	// Parse JSON response from AI
	var aiResponse struct {
		Summary  string `json:"summary"`
		Decision string `json:"decision"`
	}
	
	// Strip markdown code block markers if present
	cleanedAnalysis := strings.TrimSpace(payload.RawAnalysis)
	
	// Try regex extraction for fenced code blocks
	fencePattern := regexp.MustCompile(`(?si)^\s*` + "`" + "`" + "`" + `[A-Za-z0-9_-]*\s*(.*?)\s*` + "`" + "`" + "`" + `(?:\s*)?$`)
	if matches := fencePattern.FindStringSubmatch(cleanedAnalysis); len(matches) > 1 {
		cleanedAnalysis = strings.TrimSpace(matches[1])
	} else {
		// Fallback: remove any leading/trailing triple-backtick markers
		cleanedAnalysis = strings.TrimPrefix(cleanedAnalysis, "```")
		cleanedAnalysis = strings.TrimSuffix(cleanedAnalysis, "```")
		cleanedAnalysis = strings.TrimSpace(cleanedAnalysis)
	}
	
	// First attempt: Try to parse as-is (handles pretty-printed JSON)
	err := json.Unmarshal([]byte(cleanedAnalysis), &aiResponse)
	if err != nil {
		// Second attempt: Try to unmarshal and re-marshal to fix any formatting issues
		var rawJSON map[string]interface{}
		if err2 := json.Unmarshal([]byte(cleanedAnalysis), &rawJSON); err2 == nil {
			// Successfully parsed as generic JSON, now extract the fields
			if summary, ok := rawJSON["summary"].(string); ok {
				aiResponse.Summary = summary
			}
			if decision, ok := rawJSON["decision"].(string); ok {
				aiResponse.Decision = decision
			}
			log.Info().Str("summary_length", fmt.Sprintf("%d", len(aiResponse.Summary))).Str("decision", aiResponse.Decision).Msg("[HandleAIWebhook] Parsed AI response using fallback method")
		} else {
			log.Error().Err(err).Err(err2).Str("raw_analysis_preview", cleanedAnalysis[:min(500, len(cleanedAnalysis))]).Msg("[HandleAIWebhook] Failed to parse AI response JSON")
			http.Error(w, "Failed to parse AI response", http.StatusBadRequest)
			return
		}
	}

	// Update DB with parsed summary and decision
	ctx := r.Context()
	if err := c.service.UpdatePullRequestAnalysis(ctx, payload.PRID, aiResponse.Summary, aiResponse.Decision); err != nil {
		log.Error().Err(err).Msg("[HandleAIWebhook] Failed to update PR")
		http.Error(w, "Failed to update PR: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Notify all connected SSE clients
	pr, err := c.service.GetPullRequestByID(ctx, payload.PRID)
	if err == nil && pr != nil {
		prKey := fmt.Sprintf("%s:%d", *pr.RepoID, *pr.Number)
		notificationData := map[string]interface{}{
			"status":      "completed",
			"ai_summary":  aiResponse.Summary,
			"ai_decision": aiResponse.Decision,
			"pr_id":       payload.PRID,
		}
		notificationJSON, _ := json.Marshal(notificationData)
		GlobalSSEManager.NotifyClients(prKey, FormatSSEMessage(string(notificationJSON)))
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
		RawAnalysis string `json:"raw_analysis"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Error().Err(err).Msg("[HandleRepoAIWebhook] Failed to decode payload")
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	if payload.RepoID == "" {
		log.Error().Msg("[HandleRepoAIWebhook] repo_id is missing")
		http.Error(w, "repo_id is required", http.StatusBadRequest)
		return
	}

	if payload.RawAnalysis == "" {
		log.Error().Msg("[HandleRepoAIWebhook] raw_analysis is missing")
		http.Error(w, "raw_analysis is required", http.StatusBadRequest)
		return
	}

	// Strip markdown code block markers if present
	payload.RawAnalysis = strings.TrimSpace(payload.RawAnalysis)
	
	// Try regex extraction for fenced code blocks
	fencePattern := regexp.MustCompile(`(?si)^\s*` + "`" + "`" + "`" + `[A-Za-z0-9_-]*\s*(.*?)\s*` + "`" + "`" + "`" + `(?:\s*)?$`)
	if matches := fencePattern.FindStringSubmatch(payload.RawAnalysis); len(matches) > 1 {
		payload.RawAnalysis = strings.TrimSpace(matches[1])
	} else {
		// Fallback: remove any leading/trailing triple-backtick markers
		payload.RawAnalysis = strings.TrimPrefix(payload.RawAnalysis, "```")
		payload.RawAnalysis = strings.TrimSuffix(payload.RawAnalysis, "```")
		payload.RawAnalysis = strings.TrimSpace(payload.RawAnalysis)
	}

	// Update DB with raw markdown
	ctx := context.Background()
	if err := c.service.UpdateRepositoryAnalysis(ctx, payload.RepoID, payload.RawAnalysis); err != nil {
		log.Error().Err(err).Msg("[HandleRepoAIWebhook] Failed to update analysis")
		http.Error(w, "Failed to update analysis: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Notify all connected SSE clients
	notificationData := map[string]interface{}{
		"status":     "completed",
		"ai_summary": payload.RawAnalysis,
		"repo_id":    payload.RepoID,
	}
	notificationJSON, _ := json.Marshal(notificationData)
	GlobalSSEManager.NotifyClients(payload.RepoID, FormatSSEMessage(string(notificationJSON)))
	
	w.WriteHeader(http.StatusOK)
}

// GetPullRequestsByRepoID returns all pull requests for a repository
func (c *GithubController) GetPullRequestsByRepoID(w http.ResponseWriter, r *http.Request) {
	// Get User ID from context
	_, ok := r.Context().Value(middleware.UserContextKey).(models.User)
	if !ok {
		http.Error(w, "Unauthorized: User not found in context", http.StatusUnauthorized)
		return
	}

	// Get repository ID from path
	vars := mux.Vars(r)
	repoID := vars["id"]

	// Get all PRs for the repository
	prs, err := c.service.GetPullRequestsByRepoID(r.Context(), repoID)
	if err != nil {
		log.Error().Err(err).Str("repo_id", repoID).Msg("Failed to fetch pull requests")
		http.Error(w, "Failed to fetch pull requests", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prs)
}

// CalculateReleaseRisk triggers AI analysis for release risk calculation
func (c *GithubController) CalculateReleaseRisk(w http.ResponseWriter, r *http.Request) {
	// Get User ID from context
	userVal, ok := r.Context().Value(middleware.UserContextKey).(models.User)
	if !ok {
		http.Error(w, "Unauthorized: User not found in context", http.StatusUnauthorized)
		return
	}

	// Get repository ID from path
	vars := mux.Vars(r)
	repoID := vars["id"]

	// Parse request body to get selected PR IDs
	var requestBody struct {
		PRIDs []string `json:"pr_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		log.Error().Err(err).Msg("Failed to decode request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(requestBody.PRIDs) == 0 {
		http.Error(w, "No pull requests selected", http.StatusBadRequest)
		return
	}

	// Get all PRs for the repository
	allPRs, err := c.service.GetPullRequestsByRepoID(r.Context(), repoID)
	if err != nil {
		log.Error().Err(err).Str("repo_id", repoID).Msg("Failed to fetch pull requests")
		http.Error(w, "Failed to fetch pull requests", http.StatusInternalServerError)
		return
	}

	// Filter PRs based on selected IDs
	prMap := make(map[string]*models.PullRequest)
	for _, pr := range allPRs {
		prMap[pr.ID] = pr
	}

	var selectedPRs []*models.PullRequest
	for _, prID := range requestBody.PRIDs {
		if pr, exists := prMap[prID]; exists {
			selectedPRs = append(selectedPRs, pr)
		}
	}

	if len(selectedPRs) == 0 {
		http.Error(w, "No valid pull requests found from selection", http.StatusBadRequest)
		return
	}

	// Get repository details
	repo, err := c.service.GetRepository(r.Context(), userVal.ID, repoID)
	if err != nil {
		log.Error().Err(err).Str("repo_id", repoID).Msg("Failed to fetch repository")
		http.Error(w, "Failed to fetch repository", http.StatusNotFound)
		return
	}

	// Format PR data for AI (only selected PRs)
	var prData strings.Builder
	for i, pr := range selectedPRs {
		title := ""
		if pr.Title != nil {
			title = *pr.Title
		}
		prData.WriteString(fmt.Sprintf("\n### PR #%d: %s\n", i+1, title))
		
		if pr.Number != nil {
			prData.WriteString(fmt.Sprintf("- **PR Number**: %d\n", *pr.Number))
		}
		if pr.State != nil {
			prData.WriteString(fmt.Sprintf("- **State**: %s\n", *pr.State))
		}
		if pr.AuthorName != nil {
			prData.WriteString(fmt.Sprintf("- **Author**: %s\n", *pr.AuthorName))
		}
		if pr.AISummary != nil && *pr.AISummary != "" {
			prData.WriteString(fmt.Sprintf("- **AI Analysis**: %s\n", *pr.AISummary))
		}
		prData.WriteString("\n")
	}

	// Trigger Kestra workflow
	if err := c.service.TriggerReleaseRiskAnalysis(r.Context(), repoID, repo.Owner, repo.Name, prData.String()); err != nil {
		log.Error().Err(err).Str("repo_id", repoID).Msg("Failed to trigger release risk analysis")
		http.Error(w, "Failed to trigger release risk analysis", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "analyzing",
		"message": "Release risk analysis started for selected PRs",
	})
}

// HandleReleaseRiskCallback handles the callback from Kestra with release risk analysis results
func (c *GithubController) HandleReleaseRiskCallback(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		RepositoryID string `json:"repository_id"`
		RawAnalysis  string `json:"raw_analysis"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Error().Err(err).Msg("[HandleReleaseRiskCallback] Failed to decode JSON")
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if payload.RepositoryID == "" {
		log.Error().Msg("[HandleReleaseRiskCallback] repository_id is missing")
		http.Error(w, "repository_id is required", http.StatusBadRequest)
		return
	}

	if payload.RawAnalysis == "" {
		log.Error().Msg("[HandleReleaseRiskCallback] raw_analysis is missing")
		http.Error(w, "raw_analysis is required", http.StatusBadRequest)
		return
	}

	// Strip markdown code block markers if present
	payload.RawAnalysis = strings.TrimSpace(payload.RawAnalysis)
	
	fencePattern := regexp.MustCompile(`(?si)^\s*` + "`" + "`" + "`" + `[A-Za-z0-9_-]*\s*(.*?)\s*` + "`" + "`" + "`" + `(?:\s*)?$`)
	if matches := fencePattern.FindStringSubmatch(payload.RawAnalysis); len(matches) > 1 {
		payload.RawAnalysis = strings.TrimSpace(matches[1])
	} else {
		payload.RawAnalysis = strings.TrimPrefix(payload.RawAnalysis, "```")
		payload.RawAnalysis = strings.TrimSuffix(payload.RawAnalysis, "```")
		payload.RawAnalysis = strings.TrimSpace(payload.RawAnalysis)
	}

	// Parse the JSON response from AI
	var analysisResult struct {
		Changelog string `json:"changelog"`
		RiskScore int    `json:"risk_score"`
		Summary   string `json:"summary"`
	}

	if err := json.Unmarshal([]byte(payload.RawAnalysis), &analysisResult); err != nil {
		log.Error().Err(err).Str("raw_analysis", payload.RawAnalysis).Msg("[HandleReleaseRiskCallback] Failed to parse AI response")
		http.Error(w, "Failed to parse AI response", http.StatusBadRequest)
		return
	}

	// Update repository with release risk analysis
	ctx := context.Background()
	if err := c.service.UpdateReleaseRiskAnalysis(ctx, payload.RepositoryID, analysisResult.RiskScore, analysisResult.Changelog, payload.RawAnalysis); err != nil {
		log.Error().Err(err).Msg("[HandleReleaseRiskCallback] Failed to update release risk analysis")
		http.Error(w, "Failed to update analysis: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Notify all connected SSE clients
	notificationData := map[string]interface{}{
		"status":                "completed",
		"release_risk_score":    analysisResult.RiskScore,
		"release_changelog":     analysisResult.Changelog,
		"release_risk_analysis": payload.RawAnalysis,
		"repo_id":               payload.RepositoryID,
	}
	notificationJSON, _ := json.Marshal(notificationData)
	GlobalSSEManager.NotifyClients(payload.RepositoryID, FormatSSEMessage(string(notificationJSON)))

	log.Info().
		Str("repository_id", payload.RepositoryID).
		Int("risk_score", analysisResult.RiskScore).
		Msg("[HandleReleaseRiskCallback] Release risk analysis completed")

	w.WriteHeader(http.StatusOK)
}
