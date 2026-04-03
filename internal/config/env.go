package config

import "os"

const (
	EnvFormat    = "JPLAW_FORMAT"
	EnvBaseURL   = "JPLAW_BASE_URL"
	EnvConfigDir = "JPLAW_CONFIG_DIR"
	EnvNoColor   = "JPLAW_NO_COLOR"
	EnvVerbose   = "JPLAW_VERBOSE"
	EnvNoInput   = "JPLAW_NO_INPUT"
)

func EnvOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
