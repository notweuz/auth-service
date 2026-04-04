package main

import (
	"auth-service/internal"
	"auth-service/internal/config"
	"auth-service/internal/logger"

	"github.com/rs/zerolog/log"
)

func main() {
	logger.SetupLogger()
	cfg := config.MustSetupConfig(".")
	db, err := internal.SetupDatabase(&cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to setup database")
	}
	app, err := internal.SetupApp(&cfg, db)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to setup app")
	}

	log.Info().Int("port", cfg.GRPCPort).Msg("Starting server at")
	if err := app.Server.Serve(app.Listener); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}

	defer app.Close()
}
