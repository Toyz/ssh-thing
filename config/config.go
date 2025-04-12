package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

type SSHServer struct {
	Name           string   `toml:"name"`
	Host           string   `toml:"host"`
	User           string   `toml:"user"`
	Port           int      `toml:"port"`
	PrivateKeyPath string   `toml:"private_key_path"`
	Password       string   `toml:"password"`
	Commands       []string `toml:"commands"`
}

type Config struct {
	Servers []SSHServer `toml:"servers"`
}

func LoadConfig(filePath string) (*Config, error) {
	if filePath == "" {
		filePath = "servers.toml"

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			if _, err := os.Stat("config.toml"); err == nil {
				filePath = "config.toml"
			}
		}
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", filePath, err)
	}

	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", filePath, err)
	}

	for i, server := range cfg.Servers {
		if cfg.Servers[i].Port == 0 {
			cfg.Servers[i].Port = 22
		}

		if strings.HasPrefix(server.PrivateKeyPath, "~/") {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return nil, fmt.Errorf("failed to get user home directory: %w", err)
			}
			cfg.Servers[i].PrivateKeyPath = filepath.Join(homeDir, server.PrivateKeyPath[2:])
		}
	}

	return &cfg, nil
}

func SaveConfig(cfg *Config, filePath string) error {
	if filePath == "" {
		filePath = "servers.toml"
	}

	data, err := toml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config to %s: %w", filePath, err)
	}

	return nil
}
