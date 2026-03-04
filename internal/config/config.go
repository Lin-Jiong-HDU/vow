package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	// DefaultDailyWordCount is the default number of words to learn per day
	DefaultDailyWordCount = 20
	// DefaultConfigDir is the default configuration directory
	DefaultConfigDir = ".vow"
	// DefaultConfigFile is the default configuration file name
	DefaultConfigFile = "config.json"
)

// Config represents the application configuration
type Config struct {
	DailyWordCount int `json:"dailyWordCount"`
}

// New creates a new Config with default values
func New() *Config {
	return &Config{
		DailyWordCount: DefaultDailyWordCount,
	}
}

// DefaultConfigPath returns the default configuration file path
func DefaultConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Join(home, DefaultConfigDir, DefaultConfigFile), nil
}

// LoadConfig loads the configuration from the given path.
// If the path is empty, it uses the default path (~/.vow/config.json).
// If the configuration file or directory does not exist, it creates them with default values.
func LoadConfig(path string) (*Config, error) {
	// Use default path if not specified
	if path == "" {
		var err error
		path, err = DefaultConfigPath()
		if err != nil {
			return nil, err
		}
	}

	// Try to read existing config
	fileData, err := os.ReadFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// File doesn't exist, create default config
		return createDefaultConfig(path)
	}

	// Parse existing config
	var config Config
	if len(fileData) == 0 {
		return nil, fmt.Errorf("config file is empty: %s", path)
	}

	if err := json.Unmarshal(fileData, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// Save saves the configuration to the given path
func (c *Config) Save(path string) error {
	if path == "" {
		var err error
		path, err = DefaultConfigPath()
		if err != nil {
			return err
		}
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal config to JSON with indentation
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// createDefaultConfig creates a default configuration file at the given path
func createDefaultConfig(path string) (*Config, error) {
	config := New()

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	// Save default config
	if err := config.Save(path); err != nil {
		return nil, fmt.Errorf("failed to create default config: %w", err)
	}

	return config, nil
}
