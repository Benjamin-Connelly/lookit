package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds all application configuration.
type Config struct {
	// Root directory to browse
	Root string `mapstructure:"root"`

	// Theme: light, dark, or auto
	Theme string `mapstructure:"theme"`

	// Keymap preset: default, vim, emacs
	Keymap string `mapstructure:"keymap"`

	// Web server settings
	Server ServerConfig `mapstructure:"server"`

	// Git integration settings
	Git GitConfig `mapstructure:"git"`

	// File patterns to ignore (in addition to .gitignore)
	Ignore []string `mapstructure:"ignore"`
}

// ServerConfig holds web server settings.
type ServerConfig struct {
	Port    int    `mapstructure:"port"`
	Host    string `mapstructure:"host"`
	NoHTTPS bool   `mapstructure:"no_https"`
	Open    bool   `mapstructure:"open"`
}

// GitConfig holds git integration settings.
type GitConfig struct {
	Enabled    bool   `mapstructure:"enabled"`
	ShowStatus bool   `mapstructure:"show_status"`
	Remote     string `mapstructure:"remote"`
}

// DefaultConfig returns configuration with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Root:   ".",
		Theme:  "auto",
		Keymap: "default",
		Server: ServerConfig{
			Port: 7777,
			Host: "localhost",
		},
		Git: GitConfig{
			Enabled:    true,
			ShowStatus: true,
			Remote:     "origin",
		},
	}
}

// Load reads configuration from file and environment, merging with defaults.
func Load(cfgFile string) (*Config, error) {
	cfg := DefaultConfig()

	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	if cfgFile != "" {
		v.SetConfigFile(cfgFile)
	} else {
		configDir, err := ConfigDir()
		if err != nil {
			return cfg, nil // Use defaults if no config dir
		}
		v.AddConfigPath(configDir)
	}

	v.SetEnvPrefix("LOOKIT")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return cfg, nil // No config file is fine
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}

	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	return cfg, nil
}

// ConfigDir returns the lookit configuration directory.
func ConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "lookit"), nil
}
