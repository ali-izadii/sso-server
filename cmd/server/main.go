package main

import (
	"log"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("SSO Server starting on port %s", port)

	// TODO: Initialize database
	// TODO: Initialize services
	// TODO: Initialize handlers
	// TODO: Start HTTP server
}
