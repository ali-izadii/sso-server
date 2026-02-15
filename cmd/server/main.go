package main

import (
	"log"

	"github.com/ali/sso-server/internal/config"
	"github.com/ali/sso-server/internal/server"
	"github.com/ali/sso-server/pkg/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize global logger
	err = logger.Init(logger.Config{
		Level:  cfg.Log.Level,
		Format: cfg.Log.Format,
		File:   cfg.Log.File,
	})
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	logger.Info("config loaded", "env", config.GetEnv())

	srv, err := server.New(cfg)
	if err != nil {
		logger.Fatal("failed to create server", "error", err)
	}

	if err := srv.Start(); err != nil {
		logger.Fatal("server error", "error", err)
	}
}
