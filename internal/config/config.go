package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

// Config defines the structure for the application's configuration.
type Config struct {
	Token           string   `toml:"Token"`
	Containers      []string `toml:"Containers"`
	SystemdServices []string `toml:"SystemdServices"`
	ServerPort      string   `toml:"ServerPort,omitempty"` // omitempty allows it to be optional
	Executable      string   `toml:"Executable,omitempty"` // e.g., "docker" or "podman", defaults to "docker"
}

// LoadConfig reads the configuration file from the given path.
func LoadConfig(path string) (*Config, error) {
	var cfg Config
	configData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if _, err := toml.Decode(string(configData), &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
