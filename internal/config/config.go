package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DailyWordCount int `json:"dailyWordCount"`
}

func LoadConfig(path string) (*Config, error) {
	var config Config
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get user home directory: %w", err)
		}
		path = home + "/.vow/config.json"
	}

	fileData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if len(fileData) == 0 {
		return nil, fmt.Errorf("config file is empty")
	}

	err = json.Unmarshal(fileData, &config)

	return &config, err
}
