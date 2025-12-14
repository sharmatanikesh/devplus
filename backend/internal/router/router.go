package router

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"devplus-backend/internal/controllers/rest"
	"devplus-backend/internal/middleware"
)

// SetupRouter configures all HTTP routes for the application.
func SetupRouter(authController *rest.AuthController, githubController *rest.GithubController) *mux.Router {
	router := mux.NewRouter()

	// Apply Middleware
	router.Use(middleware.CORSMiddleware)

	// Global OPTIONS handler
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		http.NotFound(w, r)
	}).Methods("OPTIONS")

	// API V1 Config
	v1 := router.PathPrefix("/api/v1").Subrouter()

	// Public Routes
	v1.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"message": "Server is running",
		})
	}).Methods("GET")

	// Webhooks (Public)
	v1.HandleFunc("/webhook/ai", githubController.HandleAIWebhook).Methods("POST")
	v1.HandleFunc("/webhook/ai/repo", githubController.HandleRepoAIWebhook).Methods("POST")
	v1.HandleFunc("/webhook/github", githubController.HandleGithubWebhook).Methods("POST")

	// Auth Routes (Nested under /auth)
	auth := v1.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/github/login", authController.Login).Methods("GET")
	auth.HandleFunc("/github/callback", authController.Callback).Methods("GET")
	auth.HandleFunc("/logout", authController.Logout).Methods("POST")

	// Protected Routes (Session Based)
	// We create a new subrouter off v1 so wAdd new feature: AI-powered code reviewe can apply middleware ONLY to these routes
	protected := v1.PathPrefix("/").Subrouter()
	protected.Use(middleware.SessionMiddleware)

	// Register existing controllers to protected routes
	// User Routes
	protected.HandleFunc("/auth/me", authController.GetCurrentUser).Methods("GET")
	protected.HandleFunc("/repos", githubController.GetRepositories).Methods("GET")
	protected.HandleFunc("/repos/{id}", githubController.GetRepository).Methods("GET")
	protected.HandleFunc("/repos/sync", githubController.SyncRepositories).Methods("POST")
	protected.HandleFunc("/repos/{id}/sync", githubController.SyncRepository).Methods("POST")
	protected.HandleFunc("/repos/{owner}/{repo}/pulls", githubController.GetPullRequests).Methods("GET")
	protected.HandleFunc("/repos/{id}/prs/{pr_number}", githubController.GetPullRequestDetail).Methods("GET")

	// Dashboard Routes
	protected.HandleFunc("/metrics", githubController.GetMetrics).Methods("GET")
	protected.HandleFunc("/metrics/personal", githubController.GetPersonalMetrics).Methods("GET")
	protected.HandleFunc("/dashboard/stats", githubController.GetDashboardStats).Methods("GET")
	protected.HandleFunc("/dashboard/recent-prs", githubController.GetRecentActivity).Methods("GET")

	// AI Analysis Routes
	// AI Analysis Routes
	protected.HandleFunc("/repos/{id}/prs/{pr_number}/analyze", githubController.AnalyzePullRequest).Methods("POST")
	protected.HandleFunc("/repos/{id}/prs/{pr_number}/analyze/stream", githubController.StreamPullRequestAnalysis).Methods("GET")
	protected.HandleFunc("/repos/{id}/analyze", githubController.AnalyzeRepository).Methods("POST")
	protected.HandleFunc("/repos/{id}/analyze/stream", githubController.StreamRepositoryAnalysis).Methods("GET")

	// Webhooks (Should ideally be public or verified by signature, but putting under protected for now or separate if needed)
	// If it's a callback from Kestra/Gemini, it might not have the user session.
	// We need a public router for webhooks.
	// For now, let's assume Kestra can pass the auth token (unlikely) or we use a separate handler.
	// Since `router.go` likely has a public section, let's check.
	// I'll register it here for now but ideally it belongs to public routes.

	return router
}
