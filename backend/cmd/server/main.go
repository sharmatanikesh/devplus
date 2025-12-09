package main

import (
	"os"

	"devplus-backend/internal/config"
	"devplus-backend/internal/db"
	"devplus-backend/internal/router"
	"devplus-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func main() {
	// Initialize Logger
	logger.InitLogger()

	// Load Configuration
	cfg := config.LoadConfig()

	// Initialize Database
	db.InitDB(cfg)

	// Initialize Router
	r := gin.Default()

	// Health Check
	r.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Setup Routes
	router.SetupRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Info().Msgf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
