package internal

import (
	"auth-service/internal/config"
	"auth-service/internal/model"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.DatabaseDSN
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}
	log.Info().Msg("Connected to database")

	err = db.AutoMigrate(
		&model.User{},
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to migrate database")
	}

	log.Info().Msg("Migrated database successfully")
	return db, err
}
