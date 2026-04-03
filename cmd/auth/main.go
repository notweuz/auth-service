package main

import (
	"auth-service/internal/config"
	"auth-service/internal/logger"
)

func main() {
	logger.SetupLogger()
	config.MustSetupConfig(".")
}
