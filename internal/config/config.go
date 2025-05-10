package config

import (
	"encoding/json"
	"os"
	"strings"
)

// Config defines the structure for the application's configuration.
type Config struct {
	Token           string      // Token for authentication
	Images          []string    // Container image tags to pull
	ContainerNames  []string    // Container names to stop after pulling images
	ServerPort      string      // Port to listen on
	Executable      string      // Container runtime executable (docker or podman)
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	cfg := &Config{}

	// Read token from file
	tokenFile := os.Getenv("TOKEN_FILE")
	if tokenFile != "" {
		tokenData, err := os.ReadFile(tokenFile)
		if err != nil {
			return nil, err
		}
		cfg.Token = strings.TrimSpace(string(tokenData))
	}

	// Get port from environment
	cfg.ServerPort = os.Getenv("PORT")

	// Get container executable from environment
	cfg.Executable = os.Getenv("CONTAINER_EXECUTABLE")
	if cfg.Executable == "" {
		cfg.Executable = "docker" // Default to docker if not specified
	}

	// Parse images JSON from environment
	imagesJSON := os.Getenv("CONTAINERS")
	if imagesJSON != "" {
		var images []string
		if err := json.Unmarshal([]byte(imagesJSON), &images); err != nil {
			return nil, err
		}
		cfg.Images = images
	}

	// Parse container names JSON from environment
	containerNamesJSON := os.Getenv("CONTAINER_NAMES")
	if containerNamesJSON != "" {
		var containerNames []string
		if err := json.Unmarshal([]byte(containerNamesJSON), &containerNames); err != nil {
			return nil, err
		}
		cfg.ContainerNames = containerNames
	}

	return cfg, nil
}
