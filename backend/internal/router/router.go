package router

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"devplus-backend/internal/controllers/rest"
	"devplus-backend/internal/middleware"
	"devplus-backend/internal/services/auth_service"
	"devplus-backend/internal/services/github_service"
)

// SetupRouter configures all HTTP routes for the application.
func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	// Initialize Services & Controllers
	authService := auth_service.NewAuthService()
	authController := rest.NewAuthController(authService)

	githubService := github_service.NewGithubService()
	githubController := rest.NewGithubController(githubService)

	// Global OPTIONS handler
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		http.NotFound(w, r)
	}).Methods("OPTIONS")

	// Public Routes
	router.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"message": "Server is running",
		})
	}).Methods("GET")

	// Auth Routes
	router.HandleFunc("/api/auth/github/login", authController.Login).Methods("GET")
	router.HandleFunc("/api/auth/github/callback", authController.Callback).Methods("GET")
	router.HandleFunc("/api/auth/logout", authController.Logout).Methods("POST")

	// Protected Routes (Session Based)
	protected := router.PathPrefix("/api/v1").Subrouter()
	protected.Use(middleware.SessionMiddleware)

	// Register existing controllers to protected routes
	protected.HandleFunc("/repos", githubController.GetRepositories).Methods("GET")
	protected.HandleFunc("/repos/{owner}/{repo}/pulls", githubController.GetPullRequests).Methods("GET")
	// Add other protected routes here...

	return router
}

// SetupProtectedRouter is deprecated or needs to be adapted if still used by tests/main
// functionality moved to SetupRouter for now
