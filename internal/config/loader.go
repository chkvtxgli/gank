package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var configDir = "banks"

func SetConfigDir(dir string) {
	if dir != "" {
		configDir = dir
	}
}

func GetConfigDir() string {
	return configDir
}

func LoadBank(bank, account string) (*BankConfig, error) {
	var configPath string
	if account != "" {
		configPath = filepath.Join(configDir, bank, account+".yaml")
	} else {
		configPath = filepath.Join(configDir, bank+".yaml")
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("no config found at '%s': %w", configPath, err)
	}

	var cfg BankConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config at '%s': %w", configPath, err)
	}

	cfg.SetDefaults()

	return &cfg, nil
}
