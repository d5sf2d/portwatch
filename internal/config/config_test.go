package config_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/user/portwatch/internal/config"
)

func writeTempConfig(t *testing.T, v any) string {
	t.Helper()
	f, err := os.CreateTemp("", "portwatch-cfg-*.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.NewEncoder(f).Encode(v); err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestLoad_ValidConfig(t *testing.T) {
	path := writeTempConfig(t, map[string]any{
		"ports":     []int{80, 443, 8080},
		"interval":  "30s",
		"state_dir": "/tmp/pw-test",
	})

	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Ports) != 3 {
		t.Errorf("expected 3 ports, got %d", len(cfg.Ports))
	}
	if cfg.Interval != 30*time.Second {
		t.Errorf("expected 30s interval, got %s", cfg.Interval)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/path/config.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestValidate_ShortInterval(t *testing.T) {
	cfg := config.Default()
	cfg.Interval = 500 * time.Millisecond
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error for short interval")
	}
}

func TestValidate_InvalidPort(t *testing.T) {
	cfg := config.Default()
	cfg.Ports = []int{80, 99999}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error for out-of-range port")
	}
}

func TestDefault_Sensible(t *testing.T) {
	cfg := config.Default()
	if cfg.Interval == 0 {
		t.Error("default interval should not be zero")
	}
	if cfg.StateDir == "" {
		t.Error("default state_dir should not be empty")
	}
}
