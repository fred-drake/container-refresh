package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

// Container defines the structure for a container configuration
type Container struct {
	Name         string `json:"name"`
	Image        string `json:"image"`
	PullInterval string `json:"pull_interval,omitempty"`
}

// Service defines the structure for a service configuration
type Service struct {
	Name           string `json:"name"`
	RestartCommand string `json:"restart_command,omitempty"`
}

// Config defines the structure for the application's configuration.
type Config struct {
	Token           string      // Token for authentication
	Containers      []Container // Container configurations
	SystemdServices []Service   // Service configurations
	ServerPort      string      // Port to listen on
	Executable      string      // Container runtime executable (docker or podman)
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	cfg := &Config{}

	// Read token from file
	tokenFile := os.Getenv("TOKEN_FILE")
	if tokenFile != "" {
		tokenData, err := ioutil.ReadFile(tokenFile)
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

	// Parse containers JSON from environment
	containersJSON := os.Getenv("CONTAINERS")
	if containersJSON != "" {
		var containers []Container
		if err := json.Unmarshal([]byte(containersJSON), &containers); err != nil {
			return nil, err
		}
		cfg.Containers = containers
	}

	// Parse services JSON from environment
	servicesJSON := os.Getenv("SERVICES")
	if servicesJSON != "" {
		var services []Service
		if err := json.Unmarshal([]byte(servicesJSON), &services); err != nil {
			return nil, err
		}
		cfg.SystemdServices = services
	}

	return cfg, nil
}
