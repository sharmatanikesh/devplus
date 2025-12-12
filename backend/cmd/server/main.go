package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"devplus-backend/internal/db"
	"devplus-backend/internal/router"
	"devplus-backend/pkg/logger"

	"github.com/rs/zerolog/log"
)

func main() {
	// Initialize Logger
	logger.InitLogger()

	// Initialize Database
	db.GetInstance()

	// Initialize Router
	r := router.SetupRouter()

	// Initialize Services & Controllers
	// Note: Services and Controllers are now initialized inside SetupRouter or can be injected if we refactor SetupRouter to accept them.
	// For this refactor, I moved initialization inside SetupRouter to keep main clean, or I should update SetupRouter to take dependencies.
	// Looking at my router.go change, I initialized them INSIDE SetupRouter.
	// So I don't need to re-initialize them here.

	// However, I need to check if I removed the old route registrations from main.go?
	// The previous main.go had:
	// ghService := github_service.NewGithubService()
	// ghController := rest.NewGithubController(ghService)
	// protectedRouter.HandleFunc("/github/repos", ghController.GetRepositories).Methods("GET")

	// My new router.go has:
	// protected.HandleFunc("/repos", githubController.GetRepositories).Methods("GET")

	// So I can remove the explicit additional route registration here unless I want to keep specific custom routes.
	// The user asked for "make a implementation plan... using gorm...". I have replaced the old structure with the new one.

	// Let's remove the redundancy.

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	log.Info().Str("address", addr).Msg("Server starting")

	// Channel to listen for interrupt signals
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	log.Info().Msg("Server started successfully")

	// Wait for interrupt signal
	<-done

	log.Info().Msg("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to shutdown server gracefully")
		return
	}

	// Close database connection
	if err := db.Close(); err != nil {
		log.Error().Err(err).Msg("Failed to close database connection")
	} else {
		log.Info().Msg("Database connection closed successfully")
	}

	log.Info().Msg("Server shutdown successfully")
}
