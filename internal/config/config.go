package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	DefaultFormat  = "table"
	DefaultBaseURL = "https://laws.e-gov.go.jp/api/2"
	configFile     = "config.yaml"
)

type Config struct {
	Format  string `yaml:"format"`
	BaseURL string `yaml:"base_url"`
}

func ConfigDir() string {
	if d := os.Getenv(EnvConfigDir); d != "" {
		return d
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "jplaw")
}

func Load() (*Config, error) {
	path := filepath.Join(ConfigDir(), configFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return defaultConfig(), nil
		}
		return nil, fmt.Errorf("設定ファイル読み込みエラー: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("設定ファイルパースエラー: %w", err)
	}
	if cfg.Format == "" {
		cfg.Format = DefaultFormat
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = DefaultBaseURL
	}
	return &cfg, nil
}

func (c *Config) Save() error {
	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("設定ディレクトリ作成エラー: %w", err)
	}
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("設定ファイルマーシャルエラー: %w", err)
	}
	return os.WriteFile(filepath.Join(dir, configFile), data, 0600)
}

func defaultConfig() *Config {
	return &Config{
		Format:  DefaultFormat,
		BaseURL: DefaultBaseURL,
	}
}
