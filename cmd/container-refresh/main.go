package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"container-refresh/internal/config"
	"container-refresh/internal/handler"
)

const defaultConfigPath = "/etc/container-refresh.toml"
const defaultPort = "8080"

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		// For development, allow reading from the local directory if /etc path fails or isn't set
		// In a real deployment, you'd likely want to remove this fallback or make it more explicit.
		if _, err := os.Stat(defaultConfigPath); os.IsNotExist(err) {
			altConfigPath := "./container-refresh.toml"
			log.Printf("WARN: %s not found, attempting to use %s for development.", defaultConfigPath, altConfigPath)
			configPath = altConfigPath
		} else {
			configPath = defaultConfigPath
		}
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration from %s: %v", configPath, err)
	}

	h := handler.NewHandler(cfg)

	http.HandleFunc("/update", h.UpdateHandler)

	port := cfg.ServerPort
	if port == "" {
		port = defaultPort
	}

	log.Printf("Starting container-refresh server on port %s", port)
	log.Printf("Listening for POST requests on /update")
	log.Printf("Configuration loaded: Token (hidden), %d containers, %d systemd services", len(cfg.Containers), len(cfg.SystemdServices))

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
