package db

import (
	"fmt"
	"log"
	"sync"

	"devplus-backend/internal/config"
	"devplus-backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB   *gorm.DB
	once sync.Once
)

func InitDB(cfg *config.Config) {
	once.Do(func() {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			cfg.DBHost,
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBName,
			cfg.DBPort,
		)

		var err error
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal("Failed to connect to database:", err)
		}

		log.Println("Database connection established")

		// Auto-migrate models
		err = DB.AutoMigrate(
			&models.User{},
			&models.Repository{},
			&models.PullRequest{},
			&models.Commit{},
			&models.Metric{},
		)
		if err != nil {
			log.Fatal("Failed to migrate database:", err)
		}
		log.Println("Database migration completed")
	})
}
