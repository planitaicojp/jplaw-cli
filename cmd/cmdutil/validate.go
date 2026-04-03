package cmdutil

import (
	"time"

	cerrors "github.com/planitaicojp/jplaw-cli/internal/errors"
)

// ValidateDate checks that value is in YYYY-MM-DD format.
func ValidateDate(value, flagName string) error {
	if value == "" {
		return nil
	}
	if _, err := time.Parse("2006-01-02", value); err != nil {
		return &cerrors.ValidationError{Field: flagName, Message: "YYYY-MM-DD形式で指定してください"}
	}
	return nil
}
