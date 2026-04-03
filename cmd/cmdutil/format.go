package cmdutil

import (
	"github.com/spf13/cobra"

	"github.com/planitaicojp/jplaw-cli/internal/config"
)

func GetFormat(cmd *cobra.Command, defaultFormat string) string {
	format, _ := cmd.Flags().GetString("format")
	if format != "" {
		return format
	}
	if f := config.EnvOr(config.EnvFormat, ""); f != "" {
		return f
	}
	cfg, err := config.Load()
	if err != nil {
		return defaultFormat
	}
	if cfg.Format != "" && cfg.Format != config.DefaultFormat {
		return cfg.Format
	}
	return defaultFormat
}
