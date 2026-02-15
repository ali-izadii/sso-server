package main

import (
	"log"

	"github.com/ali/sso-server/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("SSO Server starting on %s:%s [env=%s]", cfg.Server.Host, cfg.Server.Port, config.GetEnv())

	// TODO: Initialize database
	// TODO: Initialize services
	// TODO: Initialize handlers
	// TODO: Start HTTP server
}
