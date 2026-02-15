package main

import (
	"log"

	"github.com/ali/sso-server/internal/config"
	"github.com/ali/sso-server/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("SSO Server starting on %s:%s [env=%s]", cfg.Server.Host, cfg.Server.Port, config.GetEnv())

	srv := server.New(cfg)
	if err := srv.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
