package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	DatabaseDSN string `mapstructure:"DATABASE_DSN"`
	LogLevel    string `mapstructure:"LOG_LEVEL"`
	JwtSecret   string `mapstructure:"JWT_SECRET"`
	GRPCPort    int    `mapstructure:"GRPC_PORT"`
}

func MustSetupConfig(path string) Config {
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Panic().Err(err).Msg("failed to read config")
	}

	config := Config{}
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Panic().Err(err).Msg("failed to unmarshal config")
	}

	return config
}
