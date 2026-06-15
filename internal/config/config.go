package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Tunnels []Tunnel `yaml:"tunnels"`
}

type Tunnel struct {
	Name    string `yaml:"name"`
	SSHHost string `yaml:"ssh_host"`
	Dynamic int    `yaml:"dynamic"`
}

func FilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".tunnels"), nil
}

func Load() (*Config, error) {
	path, err := FilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}

	return &cfg, nil
}
