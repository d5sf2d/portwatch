package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Config holds the portwatch runtime configuration.
type Config struct {
	Ports    []int         `json:"ports"`
	Interval time.Duration `json:"interval"`
	StateDir string        `json:"state_dir"`
	LogFile  string        `json:"log_file"`
}

// Default returns a Config with sensible defaults.
func Default() *Config {
	return &Config{
		Ports:    []int{},
		Interval: 60 * time.Second,
		StateDir: "/tmp/portwatch",
		LogFile:  "",
	}
}

// Load reads a JSON config file from path and merges it with defaults.
func Load(path string) (*Config, error) {
	cfg := Default()

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: open %s: %w", path, err)
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(cfg); err != nil {
		return nil, fmt.Errorf("config: decode: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks that the config values are acceptable.
func (c *Config) Validate() error {
	if c.Interval < time.Second {
		return fmt.Errorf("config: interval must be at least 1s, got %s", c.Interval)
	}
	for _, p := range c.Ports {
		if p < 1 || p > 65535 {
			return fmt.Errorf("config: invalid port %d", p)
		}
	}
	return nil
}
