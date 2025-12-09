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

	// Initialize Logger
	logger.InitLogger()

	// Initialize Database
	db.GetInstance()

	// Initialize Router
	r := router.SetupRouter()

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
