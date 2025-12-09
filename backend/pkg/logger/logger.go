package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger() {
	// Pretty print for development
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	// Default level info
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	log.Info().Msg("Logger initialized")
}
