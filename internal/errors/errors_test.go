package errors

import "testing"

func TestGetExitCode(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected int
	}{
		{"nil error", nil, ExitOK},
		{"api error", &APIError{StatusCode: 500, Message: "server error"}, ExitAPI},
		{"validation error", &ValidationError{Message: "bad input"}, ExitValidation},
		{"not found error", &NotFoundError{Resource: "法令", ID: "123"}, ExitNotFound},
		{"network error", &NetworkError{Err: nil}, ExitNetwork},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := GetExitCode(tt.err)
			if code != tt.expected {
				t.Errorf("GetExitCode() = %d, want %d", code, tt.expected)
			}
		})
	}
}

func TestAPIErrorMessage(t *testing.T) {
	tests := []struct {
		name     string
		err      *APIError
		expected string
	}{
		{"with code", &APIError{StatusCode: 400, Code: "INVALID", Message: "bad request"}, "APIエラー (HTTP 400, INVALID): bad request"},
		{"without code", &APIError{StatusCode: 500, Message: "server error"}, "APIエラー (HTTP 500): server error"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("Error() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestValidationErrorMessage(t *testing.T) {
	tests := []struct {
		name     string
		err      *ValidationError
		expected string
	}{
		{"with field", &ValidationError{Field: "law-type", Message: "不正な値"}, "バリデーションエラー (law-type): 不正な値"},
		{"without field", &ValidationError{Message: "不正な入力"}, "バリデーションエラー: 不正な入力"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("Error() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestNotFoundErrorMessage(t *testing.T) {
	err := &NotFoundError{Resource: "法令", ID: "405AC0000000088"}
	expected := "法令が見つかりません: 405AC0000000088"
	if got := err.Error(); got != expected {
		t.Errorf("Error() = %q, want %q", got, expected)
	}
}

func TestNetworkErrorUnwrap(t *testing.T) {
	inner := &ValidationError{Message: "inner"}
	err := &NetworkError{Err: inner}
	if err.Unwrap() != inner {
		t.Error("Unwrap() did not return inner error")
	}
}
