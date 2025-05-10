package main

import (
	"fmt"
	"log"
	"net/http"

	"container-refresh/internal/config"
	"container-refresh/internal/handler"
)

const defaultPort = "8080"

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration from environment: %v", err)
	}

	h := handler.NewHandler(cfg)

	http.HandleFunc("/update", h.UpdateHandler)

	port := cfg.ServerPort
	if port == "" {
		port = defaultPort
	}

	log.Printf("Starting container-refresh server on port %s", port)
	log.Printf("Listening for POST requests on /update")
	log.Printf("Configuration loaded: Token (hidden), %d containers, %d systemd services", 
		len(cfg.Containers), len(cfg.SystemdServices))
	log.Printf("Using container executable: %s", cfg.Executable)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
