package cmdutil

import (
	"github.com/planitaicojp/jplaw-cli/internal/api"
	"github.com/planitaicojp/jplaw-cli/internal/config"
)

func NewClient() (*api.Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	baseURL := config.EnvOr(config.EnvBaseURL, cfg.BaseURL)
	return api.NewClient(baseURL), nil
}
