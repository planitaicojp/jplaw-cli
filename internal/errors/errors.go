package errors

import "fmt"

type ExitCoder interface {
	ExitCode() int
}

type APIError struct {
	StatusCode int
	Code       string
	Message    string
}

func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("APIエラー (HTTP %d, %s): %s", e.StatusCode, e.Code, e.Message)
	}
	return fmt.Sprintf("APIエラー (HTTP %d): %s", e.StatusCode, e.Message)
}

func (e *APIError) ExitCode() int { return ExitAPI }

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("バリデーションエラー (%s): %s", e.Field, e.Message)
	}
	return fmt.Sprintf("バリデーションエラー: %s", e.Message)
}

func (e *ValidationError) ExitCode() int { return ExitValidation }

type NotFoundError struct {
	Resource string
	ID       string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%sが見つかりません: %s", e.Resource, e.ID)
}

func (e *NotFoundError) ExitCode() int { return ExitNotFound }

type NetworkError struct {
	Err error
}

func (e *NetworkError) Error() string {
	return fmt.Sprintf("ネットワークエラー: %v", e.Err)
}

func (e *NetworkError) Unwrap() error { return e.Err }
func (e *NetworkError) ExitCode() int { return ExitNetwork }

func GetExitCode(err error) int {
	if err == nil {
		return ExitOK
	}
	if ec, ok := err.(ExitCoder); ok {
		return ec.ExitCode()
	}
	return ExitGeneral
}
