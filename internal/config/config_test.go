package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefault(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("JPLAW_CONFIG_DIR", tmp)
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if cfg.Format != DefaultFormat {
		t.Errorf("Format = %q, want %q", cfg.Format, DefaultFormat)
	}
	if cfg.BaseURL != DefaultBaseURL {
		t.Errorf("BaseURL = %q, want %q", cfg.BaseURL, DefaultBaseURL)
	}
}

func TestLoadFromFile(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("JPLAW_CONFIG_DIR", tmp)
	data := []byte("format: json\nbase_url: https://example.com/api/2\n")
	if err := os.WriteFile(filepath.Join(tmp, "config.yaml"), data, 0600); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if cfg.Format != "json" {
		t.Errorf("Format = %q, want %q", cfg.Format, "json")
	}
}

func TestSaveAndLoad(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("JPLAW_CONFIG_DIR", tmp)
	cfg := &Config{Format: "json", BaseURL: "https://example.com/api/2"}
	if err := cfg.Save(); err != nil {
		t.Fatalf("Save() error: %v", err)
	}
	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if loaded.Format != "json" {
		t.Errorf("Format = %q, want %q", loaded.Format, "json")
	}
}

func TestEnvOr(t *testing.T) {
	t.Setenv("JPLAW_FORMAT", "json")
	if got := EnvOr(EnvFormat, "table"); got != "json" {
		t.Errorf("EnvOr() = %q, want %q", got, "json")
	}
	t.Setenv("JPLAW_FORMAT", "")
	if got := EnvOr(EnvFormat, "table"); got != "table" {
		t.Errorf("EnvOr() = %q, want %q", got, "table")
	}
}

func TestConfigDir(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("JPLAW_CONFIG_DIR", tmp)
	if got := ConfigDir(); got != tmp {
		t.Errorf("ConfigDir() = %q, want %q", got, tmp)
	}
}
