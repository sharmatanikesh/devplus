package db

import (
	"fmt"
	"sync"

	"devplus-backend/internal/config"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	instance *gorm.DB
	once     sync.Once
)

func GetInstance() *gorm.DB {
	once.Do(func() {
		log.Info().Msg("Establishing database connection...")

		cfg := config.LoadConfig()
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			cfg.DBHost,
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBName,
			cfg.DBPort,
		)
		log.Info().
			Str("host", cfg.DBHost).
			Str("dbname", cfg.DBName).
			Str("port", cfg.DBPort).
			Msg("Connecting to database")

		db, err := gorm.Open(postgres.New(postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: true,
		}), &gorm.Config{})
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to connect to database")
		}

		instance = db
		log.Info().Msg("Database connection established successfully")
	})

	return instance
}

// Close closes the database connection
func Close() error {
	if instance != nil {
		db, err := instance.DB()
		if err != nil {
			return err
		}
		log.Info().Msg("Closing database connection...")
		return db.Close()
	}
	return nil
}

// IsConnected checks if the database connection is alive
func IsConnected() bool {
	if instance == nil {
		return false
	}

	sqlDB, err := instance.DB()
	if err != nil {
		return false
	}

	return sqlDB.Ping() == nil
}
