# jplaw-cli Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a Go CLI tool that wraps the e-Gov Law API v2 for searching, listing, and reading Japanese laws.

**Architecture:** Cobra-based CLI with `cmd/` subcommands and `internal/` packages for API client, models, config, output formatting, error handling, and law text conversion. Follows conoha-cli patterns with auth layer removed.

**Tech Stack:** Go 1.26+, github.com/spf13/cobra, gopkg.in/yaml.v3

---

### Task 1: Project Scaffolding

**Files:**
- Create: `go.mod`
- Create: `main.go`
- Create: `Makefile`
- Create: `.goreleaser.yaml`
- Create: `.golangci.yml`

- [ ] **Step 1: Initialize Go module**

Run: `cd /root/dev/planitai/jplaw-cli && go mod init github.com/planitaicojp/jplaw-cli`
Expected: `go.mod` created

- [ ] **Step 2: Create main.go**

Create `main.go`:

```go
package main

import "github.com/planitaicojp/jplaw-cli/cmd"

func main() {
	cmd.Execute()
}
```

- [ ] **Step 3: Create Makefile**

Create `Makefile`:

```makefile
BINARY := jplaw
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-s -w -X github.com/planitaicojp/jplaw-cli/cmd.version=$(VERSION)"

.PHONY: build test lint clean install

build:
	go build $(LDFLAGS) -o $(BINARY) .

install:
	go install $(LDFLAGS) .

test:
	go test ./... -v

lint:
	golangci-lint run ./...

clean:
	rm -f $(BINARY)
	rm -rf dist/

coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
```

- [ ] **Step 4: Create .goreleaser.yaml**

Create `.goreleaser.yaml`:

```yaml
version: 2

builds:
  - binary: jplaw
    ldflags:
      - -s -w -X github.com/planitaicojp/jplaw-cli/cmd.version={{.Version}}
    goos:
      - linux
      - darwin
      - windows
    goarch:
      - amd64
      - arm64

archives:
  - format: tar.gz
    name_template: "{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}"
    format_overrides:
      - goos: windows
        format: zip

checksum:
  name_template: "checksums.txt"

changelog:
  sort: asc
  filters:
    exclude:
      - "^docs:"
      - "^test:"
      - "^chore:"
```

- [ ] **Step 5: Create .golangci.yml**

Create `.golangci.yml`:

```yaml
version: "2"
linters:
  default: none
  enable:
    - govet
    - ineffassign
    - staticcheck
    - unused
    - errcheck
  settings:
    errcheck:
      check-type-assertions: false
      check-blank: false
      disable-default-exclusions: false
      exclude-functions:
        - io.Copy
        - (io.Closer).Close
        - (*os.File).Close
        - (net/http.ResponseWriter).Write
        - (*encoding/json.Encoder).Encode
        - (*encoding/json.Decoder).Decode
formatters:
  enable:
    - gofmt
```

- [ ] **Step 6: Commit**

```bash
git add go.mod main.go Makefile .goreleaser.yaml .golangci.yml
git commit -m "feat: project scaffolding with go.mod, main.go, Makefile, goreleaser, golangci"
```

---

### Task 2: Error Types and Exit Codes

**Files:**
- Create: `internal/errors/exitcodes.go`
- Create: `internal/errors/errors.go`
- Create: `internal/errors/errors_test.go`

- [ ] **Step 1: Write the failing test**

Create `internal/errors/errors_test.go`:

```go
package errors

import (
	"testing"
)

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
		{
			"with code",
			&APIError{StatusCode: 400, Code: "INVALID", Message: "bad request"},
			"APIエラー (HTTP 400, INVALID): bad request",
		},
		{
			"without code",
			&APIError{StatusCode: 500, Message: "server error"},
			"APIエラー (HTTP 500): server error",
		},
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
		{
			"with field",
			&ValidationError{Field: "law-type", Message: "不正な値"},
			"バリデーションエラー (law-type): 不正な値",
		},
		{
			"without field",
			&ValidationError{Message: "不正な入力"},
			"バリデーションエラー: 不正な入力",
		},
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
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /root/dev/planitai/jplaw-cli && go test ./internal/errors/ -v`
Expected: FAIL — package does not exist

- [ ] **Step 3: Create exit codes**

Create `internal/errors/exitcodes.go`:

```go
package errors

const (
	ExitOK         = 0
	ExitGeneral    = 1
	ExitValidation = 2
	ExitNotFound   = 3
	ExitAPI        = 4
	ExitNetwork    = 5
)
```

- [ ] **Step 4: Create error types**

Create `internal/errors/errors.go`:

```go
package errors

import "fmt"

// ExitCoder is implemented by errors that carry a process exit code.
type ExitCoder interface {
	ExitCode() int
}

// APIError represents an error returned by the e-Gov API.
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

func (e *APIError) ExitCode() int {
	return ExitAPI
}

// ValidationError represents invalid user input.
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

func (e *ValidationError) ExitCode() int {
	return ExitValidation
}

// NotFoundError indicates that a requested law was not found.
type NotFoundError struct {
	Resource string
	ID       string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%sが見つかりません: %s", e.Resource, e.ID)
}

func (e *NotFoundError) ExitCode() int {
	return ExitNotFound
}

// NetworkError wraps an underlying network-level error.
type NetworkError struct {
	Err error
}

func (e *NetworkError) Error() string {
	return fmt.Sprintf("ネットワークエラー: %v", e.Err)
}

func (e *NetworkError) Unwrap() error {
	return e.Err
}

func (e *NetworkError) ExitCode() int {
	return ExitNetwork
}

// GetExitCode returns the exit code for the given error.
func GetExitCode(err error) int {
	if err == nil {
		return ExitOK
	}
	if ec, ok := err.(ExitCoder); ok {
		return ec.ExitCode()
	}
	return ExitGeneral
}
```

- [ ] **Step 5: Run tests**

Run: `cd /root/dev/planitai/jplaw-cli && go test ./internal/errors/ -v`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add internal/errors/
git commit -m "feat: add error types and exit codes"
```

---

### Task 3: Configuration Management

**Files:**
- Create: `internal/config/config.go`
- Create: `internal/config/env.go`
- Create: `internal/config/config_test.go`

- [ ] **Step 1: Write the failing test**

Create `internal/config/config_test.go`:

```go
package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefault(t *testing.T) {
	// Set config dir to temp
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
	if cfg.BaseURL != "https://example.com/api/2" {
		t.Errorf("BaseURL = %q, want %q", cfg.BaseURL, "https://example.com/api/2")
	}
}

func TestSaveAndLoad(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("JPLAW_CONFIG_DIR", tmp)

	cfg := &Config{
		Format:  "json",
		BaseURL: "https://example.com/api/2",
	}
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
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /root/dev/planitai/jplaw-cli && go test ./internal/config/ -v`
Expected: FAIL — package does not exist

- [ ] **Step 3: Create env.go**

Create `internal/config/env.go`:

```go
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

// EnvOr returns the environment variable value if set, otherwise the fallback.
func EnvOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
```

- [ ] **Step 4: Create config.go**

Create `internal/config/config.go`:

```go
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

// Config holds the user configuration.
type Config struct {
	Format  string `yaml:"format"`
	BaseURL string `yaml:"base_url"`
}

// ConfigDir returns the configuration directory path.
func ConfigDir() string {
	if d := os.Getenv(EnvConfigDir); d != "" {
		return d
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "jplaw")
}

// Load reads the config file. Returns defaults if the file does not exist.
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

// Save writes the config to file.
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
```

- [ ] **Step 5: Add yaml.v3 dependency**

Run: `cd /root/dev/planitai/jplaw-cli && go get gopkg.in/yaml.v3`

- [ ] **Step 6: Run tests**

Run: `cd /root/dev/planitai/jplaw-cli && go test ./internal/config/ -v`
Expected: PASS

- [ ] **Step 7: Commit**

```bash
git add internal/config/ go.mod go.sum
git commit -m "feat: add configuration management with env var support"
```

---

### Task 4: Output Formatters

**Files:**
- Create: `internal/output/formatter.go`
- Create: `internal/output/json.go`
- Create: `internal/output/table.go`
- Create: `internal/output/text.go`
- Create: `internal/output/formatter_test.go`

- [ ] **Step 1: Write the failing test**

Create `internal/output/formatter_test.go`:

```go
package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

type testRow struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func TestJSONFormatter(t *testing.T) {
	var buf bytes.Buffer
	f := New("json")

	rows := []testRow{
		{ID: "1", Name: "テスト法"},
	}

	if err := f.Format(&buf, rows); err != nil {
		t.Fatalf("Format() error: %v", err)
	}

	var result []testRow
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if len(result) != 1 || result[0].Name != "テスト法" {
		t.Errorf("unexpected result: %+v", result)
	}
}

func TestTableFormatter(t *testing.T) {
	var buf bytes.Buffer
	f := New("table")

	rows := []testRow{
		{ID: "1", Name: "テスト法"},
		{ID: "2", Name: "サンプル法"},
	}

	if err := f.Format(&buf, rows); err != nil {
		t.Fatalf("Format() error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "ID") {
		t.Error("table output missing ID header")
	}
	if !strings.Contains(out, "テスト法") {
		t.Error("table output missing row data")
	}
}

func TestTableFormatterEmpty(t *testing.T) {
	var buf bytes.Buffer
	f := New("table")

	rows := []testRow{}
	if err := f.Format(&buf, rows); err != nil {
		t.Fatalf("Format() error: %v", err)
	}

	if buf.Len() != 0 {
		t.Errorf("expected empty output, got: %q", buf.String())
	}
}

func TestTextFormatter(t *testing.T) {
	var buf bytes.Buffer
	f := New("text")

	text := "第一条　この法律は…"
	if err := f.Format(&buf, text); err != nil {
		t.Fatalf("Format() error: %v", err)
	}

	if !strings.Contains(buf.String(), "第一条") {
		t.Error("text output missing content")
	}
}

func TestNewFormatterDefault(t *testing.T) {
	f := New("unknown")
	if _, ok := f.(*TableFormatter); !ok {
		t.Error("New() with unknown format should return TableFormatter")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /root/dev/planitai/jplaw-cli && go test ./internal/output/ -v`
Expected: FAIL — package does not exist

- [ ] **Step 3: Create formatter.go**

Create `internal/output/formatter.go`:

```go
package output

import "io"

// Formatter formats and writes data to a writer.
type Formatter interface {
	Format(w io.Writer, data any) error
}

// New creates a formatter for the given format name.
func New(format string) Formatter {
	switch format {
	case "json":
		return &JSONFormatter{}
	case "text":
		return &TextFormatter{}
	default:
		return &TableFormatter{}
	}
}
```

- [ ] **Step 4: Create json.go**

Create `internal/output/json.go`:

```go
package output

import (
	"encoding/json"
	"io"
)

// JSONFormatter outputs data as indented JSON.
type JSONFormatter struct{}

func (f *JSONFormatter) Format(w io.Writer, data any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	return enc.Encode(data)
}
```

- [ ] **Step 5: Create table.go**

Create `internal/output/table.go`:

```go
package output

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"text/tabwriter"
)

// TableFormatter outputs data as an aligned table.
type TableFormatter struct{}

func (f *TableFormatter) Format(w io.Writer, data any) error {
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Slice {
		_, err := fmt.Fprintf(w, "%v\n", data)
		return err
	}

	if val.Len() == 0 {
		return nil
	}

	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)

	// Get headers from json tags
	elem := val.Index(0)
	if elem.Kind() == reflect.Ptr {
		elem = elem.Elem()
	}
	elemType := elem.Type()

	headers := make([]string, elemType.NumField())
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		name := field.Tag.Get("json")
		if idx := strings.IndexByte(name, ','); idx != -1 {
			name = name[:idx]
		}
		if name == "" || name == "-" {
			name = field.Name
		}
		headers[i] = strings.ToUpper(name)
	}
	if _, err := fmt.Fprintln(tw, strings.Join(headers, "\t")); err != nil {
		return err
	}

	// Write rows
	for i := 0; i < val.Len(); i++ {
		row := val.Index(i)
		if row.Kind() == reflect.Ptr {
			row = row.Elem()
		}
		fields := make([]string, row.NumField())
		for j := 0; j < row.NumField(); j++ {
			fields[j] = fmt.Sprintf("%v", row.Field(j).Interface())
		}
		if _, err := fmt.Fprintln(tw, strings.Join(fields, "\t")); err != nil {
			return err
		}
	}

	return tw.Flush()
}
```

- [ ] **Step 6: Create text.go**

Create `internal/output/text.go`:

```go
package output

import (
	"fmt"
	"io"
)

// TextFormatter outputs data as plain text.
type TextFormatter struct{}

func (f *TextFormatter) Format(w io.Writer, data any) error {
	_, err := fmt.Fprintln(w, data)
	return err
}
```

- [ ] **Step 7: Run tests**

Run: `cd /root/dev/planitai/jplaw-cli && go test ./internal/output/ -v`
Expected: PASS

- [ ] **Step 8: Commit**

```bash
git add internal/output/
git commit -m "feat: add output formatters (json, table, text)"
```

---

### Task 5: API Models

**Files:**
- Create: `internal/model/enums.go`
- Create: `internal/model/law.go`
- Create: `internal/model/revision.go`
- Create: `internal/model/keyword.go`
- Create: `internal/model/lawdata.go`
- Create: `internal/model/error.go`
- Create: `internal/model/enums_test.go`

- [ ] **Step 1: Write the failing test**

Create `internal/model/enums_test.go`:

```go
package model

import (
	"encoding/json"
	"testing"
)

func TestLawTypeLabel(t *testing.T) {
	tests := []struct {
		input    LawType
		expected string
	}{
		{LawTypeConstitution, "憲法"},
		{LawTypeAct, "法律"},
		{LawTypeCabinetOrder, "政令"},
		{LawTypeImperialOrder, "勅令"},
		{LawTypeMinisterialOrdinance, "府省令"},
		{LawTypeRule, "規則"},
		{LawTypeMisc, "その他"},
		{LawType("Unknown"), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(string(tt.input), func(t *testing.T) {
			if got := tt.input.Label(); got != tt.expected {
				t.Errorf("Label() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestEraLabel(t *testing.T) {
	tests := []struct {
		input    Era
		expected string
	}{
		{EraMeiji, "明治"},
		{EraTaisho, "大正"},
		{EraShowa, "昭和"},
		{EraHeisei, "平成"},
		{EraReiwa, "令和"},
	}

	for _, tt := range tests {
		t.Run(string(tt.input), func(t *testing.T) {
			if got := tt.input.Label(); got != tt.expected {
				t.Errorf("Label() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestLawTypeJSON(t *testing.T) {
	type wrapper struct {
		Type LawType `json:"type"`
	}
	data := []byte(`{"type":"Act"}`)
	var w wrapper
	if err := json.Unmarshal(data, &w); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if w.Type != LawTypeAct {
		t.Errorf("Type = %q, want %q", w.Type, LawTypeAct)
	}
}

func TestLawTypeLabelFromAPI(t *testing.T) {
	lt := LawTypeLabelFromAPI("法律")
	if lt != LawTypeAct {
		t.Errorf("LawTypeLabelFromAPI('法律') = %q, want %q", lt, LawTypeAct)
	}

	lt = LawTypeLabelFromAPI("unknown")
	if lt != "" {
		t.Errorf("LawTypeLabelFromAPI('unknown') = %q, want empty", lt)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /root/dev/planitai/jplaw-cli && go test ./internal/model/ -v`
Expected: FAIL — package does not exist

- [ ] **Step 3: Create enums.go**

Create `internal/model/enums.go`:

```go
package model

// LawType represents the type of a law.
type LawType string

const (
	LawTypeConstitution         LawType = "Constitution"
	LawTypeAct                  LawType = "Act"
	LawTypeCabinetOrder         LawType = "CabinetOrder"
	LawTypeImperialOrder        LawType = "ImperialOrder"
	LawTypeMinisterialOrdinance LawType = "MinisterialOrdinance"
	LawTypeRule                 LawType = "Rule"
	LawTypeMisc                 LawType = "Misc"
)

var lawTypeLabels = map[LawType]string{
	LawTypeConstitution:         "憲法",
	LawTypeAct:                  "法律",
	LawTypeCabinetOrder:         "政令",
	LawTypeImperialOrder:        "勅令",
	LawTypeMinisterialOrdinance: "府省令",
	LawTypeRule:                 "規則",
	LawTypeMisc:                 "その他",
}

var labelToLawType = map[string]LawType{
	"憲法":  LawTypeConstitution,
	"法律":  LawTypeAct,
	"政令":  LawTypeCabinetOrder,
	"勅令":  LawTypeImperialOrder,
	"府省令": LawTypeMinisterialOrdinance,
	"規則":  LawTypeRule,
	"その他": LawTypeMisc,
}

// Label returns the Japanese label for the law type.
func (lt LawType) Label() string {
	if label, ok := lawTypeLabels[lt]; ok {
		return label
	}
	return string(lt)
}

// LawTypeLabelFromAPI returns the LawType for a Japanese label.
func LawTypeLabelFromAPI(label string) LawType {
	if lt, ok := labelToLawType[label]; ok {
		return lt
	}
	return ""
}

// Era represents the Japanese era.
type Era string

const (
	EraMeiji  Era = "Meiji"
	EraTaisho Era = "Taisho"
	EraShowa  Era = "Showa"
	EraHeisei Era = "Heisei"
	EraReiwa  Era = "Reiwa"
)

var eraLabels = map[Era]string{
	EraMeiji:  "明治",
	EraTaisho: "大正",
	EraShowa:  "昭和",
	EraHeisei: "平成",
	EraReiwa:  "令和",
}

var labelToEra = map[string]Era{
	"明治": EraMeiji,
	"大正": EraTaisho,
	"昭和": EraShowa,
	"平成": EraHeisei,
	"令和": EraReiwa,
}

// Label returns the Japanese label for the era.
func (e Era) Label() string {
	if label, ok := eraLabels[e]; ok {
		return label
	}
	return string(e)
}

// EraFromLabel returns the Era for a Japanese label.
func EraFromLabel(label string) Era {
	if e, ok := labelToEra[label]; ok {
		return e
	}
	return ""
}

// RepealStatus represents the repeal status of a law.
type RepealStatus string

const (
	RepealStatusNone                RepealStatus = "None"
	RepealStatusRepeal              RepealStatus = "Repeal"
	RepealStatusExpire              RepealStatus = "Expire"
	RepealStatusSuspend             RepealStatus = "Suspend"
	RepealStatusLossOfEffectiveness RepealStatus = "LossOfEffectiveness"
)

// AmendmentType represents the amendment type.
type AmendmentType string

const (
	AmendmentTypeNew     AmendmentType = "1"
	AmendmentTypeAmended AmendmentType = "3"
	AmendmentTypeRepeal  AmendmentType = "8"
)

// Mission represents new/partial amendment.
type Mission string

const (
	MissionNew     Mission = "New"
	MissionPartial Mission = "Partial"
)

// CurrentRevisionStatus represents the current revision status.
type CurrentRevisionStatus string

const (
	CurrentRevisionStatusCurrentEnforced  CurrentRevisionStatus = "CurrentEnforced"
	CurrentRevisionStatusUnEnforced       CurrentRevisionStatus = "UnEnforced"
	CurrentRevisionStatusPreviousEnforced CurrentRevisionStatus = "PreviousEnforced"
	CurrentRevisionStatusRepeal           CurrentRevisionStatus = "Repeal"
)
```

- [ ] **Step 4: Create law.go**

Create `internal/model/law.go`:

```go
package model

// LawInfo contains non-revision-dependent law metadata.
type LawInfo struct {
	LawType     LawType `json:"law_type" yaml:"law_type"`
	LawID       string  `json:"law_id" yaml:"law_id"`
	LawNum      string  `json:"law_num" yaml:"law_num"`
	LawNumEra   Era     `json:"law_num_era" yaml:"law_num_era"`
	LawNumYear  int     `json:"law_num_year" yaml:"law_num_year"`
	LawNumType  LawType `json:"law_num_type" yaml:"law_num_type"`
	LawNumNum   string  `json:"law_num_num" yaml:"law_num_num"`
	Promulgation string `json:"promulgation_date" yaml:"promulgation_date"`
}

// LawItem represents a single law in the list response.
type LawItem struct {
	LawInfo             LawInfo      `json:"law_info" yaml:"law_info"`
	RevisionInfo        RevisionInfo `json:"revision_info" yaml:"revision_info"`
	CurrentRevisionInfo RevisionInfo `json:"current_revision_info" yaml:"current_revision_info"`
}

// LawsResponse is the response from GET /laws.
type LawsResponse struct {
	TotalCount int64     `json:"total_count" yaml:"total_count"`
	Count      int64     `json:"count" yaml:"count"`
	NextOffset *int64    `json:"next_offset" yaml:"next_offset"`
	Laws       []LawItem `json:"laws" yaml:"laws"`
}

// LawListRow is a flattened row for table output.
type LawListRow struct {
	LawID       string `json:"law_id"`
	LawNum      string `json:"law_num"`
	LawTitle    string `json:"law_title"`
	LawType     string `json:"law_type"`
	Promulgation string `json:"promulgation_date"`
}
```

- [ ] **Step 5: Create revision.go**

Create `internal/model/revision.go`:

```go
package model

// RevisionInfo contains revision-dependent law metadata.
type RevisionInfo struct {
	LawRevisionID               string                `json:"law_revision_id" yaml:"law_revision_id"`
	LawType                     LawType               `json:"law_type" yaml:"law_type"`
	LawTitle                    string                `json:"law_title" yaml:"law_title"`
	LawTitleKana                string                `json:"law_title_kana" yaml:"law_title_kana"`
	Abbrev                      string                `json:"abbrev" yaml:"abbrev"`
	Category                    string                `json:"category" yaml:"category"`
	Updated                     string                `json:"updated" yaml:"updated"`
	AmendmentPromulgateDate     string                `json:"amendment_promulgate_date" yaml:"amendment_promulgate_date"`
	AmendmentEnforcementDate    string                `json:"amendment_enforcement_date" yaml:"amendment_enforcement_date"`
	AmendmentEnforcementComment string                `json:"amendment_enforcement_comment" yaml:"amendment_enforcement_comment"`
	AmendmentLawID              string                `json:"amendment_law_id" yaml:"amendment_law_id"`
	AmendmentLawTitle           string                `json:"amendment_law_title" yaml:"amendment_law_title"`
	AmendmentLawTitleKana       string                `json:"amendment_law_title_kana" yaml:"amendment_law_title_kana"`
	AmendmentLawNum             string                `json:"amendment_law_num" yaml:"amendment_law_num"`
	AmendmentType               AmendmentType         `json:"amendment_type" yaml:"amendment_type"`
	RepealStatus                RepealStatus          `json:"repeal_status" yaml:"repeal_status"`
	RepealDate                  *string               `json:"repeal_date" yaml:"repeal_date"`
	RemainInForce               bool                  `json:"remain_in_force" yaml:"remain_in_force"`
	Mission                     Mission               `json:"mission" yaml:"mission"`
	CurrentRevisionStatus       CurrentRevisionStatus `json:"current_revision_status" yaml:"current_revision_status"`
}

// RevisionsResponse is the response from GET /law_revisions.
type RevisionsResponse struct {
	LawInfo   LawInfo        `json:"law_info" yaml:"law_info"`
	Revisions []RevisionInfo `json:"revisions" yaml:"revisions"`
}

// RevisionListRow is a flattened row for table output.
type RevisionListRow struct {
	AmendmentDate  string `json:"amendment_date"`
	AmendmentTitle string `json:"amendment_title"`
	AmendmentNum   string `json:"amendment_num"`
	AmendmentType  string `json:"amendment_type"`
}
```

- [ ] **Step 6: Create keyword.go**

Create `internal/model/keyword.go`:

```go
package model

// Sentence is a matched sentence in keyword search.
type Sentence struct {
	Position string `json:"position" yaml:"position"`
	Text     string `json:"text" yaml:"text"`
}

// KeywordItem is a single result in keyword search.
type KeywordItem struct {
	LawInfo      LawInfo      `json:"law_info" yaml:"law_info"`
	RevisionInfo RevisionInfo `json:"revision_info" yaml:"revision_info"`
	Sentences    []Sentence   `json:"sentences" yaml:"sentences"`
}

// KeywordResponse is the response from GET /keyword.
type KeywordResponse struct {
	TotalCount    int64         `json:"total_count" yaml:"total_count"`
	SentenceCount int64         `json:"sentence_count" yaml:"sentence_count"`
	NextOffset    *int64        `json:"next_offset" yaml:"next_offset"`
	Items         []KeywordItem `json:"items" yaml:"items"`
}

// KeywordListRow is a flattened row for table output.
type KeywordListRow struct {
	LawNum   string `json:"law_num"`
	LawTitle string `json:"law_title"`
	Match    string `json:"match"`
}
```

- [ ] **Step 7: Create lawdata.go**

Create `internal/model/lawdata.go`:

```go
package model

import "encoding/json"

// LawDataResponse is the response from GET /law_data.
type LawDataResponse struct {
	LawInfo           LawInfo         `json:"law_info" yaml:"law_info"`
	RevisionInfo      RevisionInfo    `json:"revision_info" yaml:"revision_info"`
	LawFullText       json.RawMessage `json:"law_full_text" yaml:"law_full_text"`
	AttachedFilesInfo *AttachedFilesInfo `json:"attached_files_info,omitempty" yaml:"attached_files_info,omitempty"`
}

// AttachedFilesInfo contains attachment metadata.
type AttachedFilesInfo struct {
	ImageData     string         `json:"image_data,omitempty" yaml:"image_data,omitempty"`
	AttachedFiles []AttachedFile `json:"attached_files,omitempty" yaml:"attached_files,omitempty"`
}

// AttachedFile represents a single attached file.
type AttachedFile struct {
	LawRevisionID string `json:"law_revision_id" yaml:"law_revision_id"`
	Src           string `json:"src" yaml:"src"`
	Updated       string `json:"updated" yaml:"updated"`
}
```

- [ ] **Step 8: Create error.go**

Create `internal/model/error.go`:

```go
package model

// ErrorInfo is the error response from the e-Gov API.
type ErrorInfo struct {
	Code    string `json:"code" yaml:"code"`
	Message string `json:"message" yaml:"message"`
}
```

- [ ] **Step 9: Run tests**

Run: `cd /root/dev/planitai/jplaw-cli && go test ./internal/model/ -v`
Expected: PASS

- [ ] **Step 10: Commit**

```bash
git add internal/model/
git commit -m "feat: add API models (law, revision, keyword, lawdata, enums)"
```

---

### Task 6: Base API Client

**Files:**
- Create: `internal/api/client.go`
- Create: `internal/api/debug.go`
- Create: `internal/api/client_test.go`

- [ ] **Step 1: Write the failing test**

Create `internal/api/client_test.go`:

```go
package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientGet(t *testing.T) {
	type response struct {
		Name string `json:"name"`
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if ua := r.Header.Get("User-Agent"); ua == "" {
			t.Error("missing User-Agent header")
		}
		if accept := r.Header.Get("Accept"); accept != "application/json" {
			t.Errorf("Accept = %q, want application/json", accept)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response{Name: "テスト"})
	}))
	defer server.Close()

	client := NewClient(server.URL)
	var result response
	if err := client.Get("/test", &result); err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if result.Name != "テスト" {
		t.Errorf("Name = %q, want %q", result.Name, "テスト")
	}
}

func TestClientGetAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"code":    "INVALID",
			"message": "bad request",
		})
	}))
	defer server.Close()

	client := NewClient(server.URL)
	err := client.Get("/test", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestClientGet404(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"code":    "NOT_FOUND",
			"message": "not found",
		})
	}))
	defer server.Close()

	client := NewClient(server.URL)
	err := client.Get("/test", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestClientGetRaw(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write([]byte("binary data"))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	data, err := client.GetRaw("/test")
	if err != nil {
		t.Fatalf("GetRaw() error: %v", err)
	}
	if string(data) != "binary data" {
		t.Errorf("got %q, want %q", string(data), "binary data")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /root/dev/planitai/jplaw-cli && go test ./internal/api/ -v`
Expected: FAIL — package does not exist

- [ ] **Step 3: Create debug.go**

Create `internal/api/debug.go`:

```go
package api

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	DebugNone    = 0
	DebugAPI     = 1
	DebugVerbose = 2
)

var debugLevel int

// SetDebugLevel sets the debug logging level.
func SetDebugLevel(level int) {
	debugLevel = level
}

func debugLogRequest(req *http.Request) {
	if debugLevel < DebugAPI {
		return
	}
	fmt.Fprintf(os.Stderr, ">> %s %s\n", req.Method, req.URL)
	if debugLevel >= DebugVerbose {
		for k, v := range req.Header {
			fmt.Fprintf(os.Stderr, ">> %s: %s\n", k, v)
		}
	}
}

func debugLogResponse(resp *http.Response, elapsed time.Duration) {
	if debugLevel < DebugAPI {
		return
	}
	fmt.Fprintf(os.Stderr, "<< %d %s (%s)\n", resp.StatusCode, resp.Status, elapsed)
}
```

- [ ] **Step 4: Create client.go**

Create `internal/api/client.go`:

```go
package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	cerrors "github.com/planitaicojp/jplaw-cli/internal/errors"
	"github.com/planitaicojp/jplaw-cli/internal/model"
)

// UserAgent is sent with every request.
var UserAgent = "jplaw-cli/dev"

const (
	defaultTimeout = 30 * time.Second
	maxRetries     = 3
)

// Client is the base HTTP client for the e-Gov Law API.
type Client struct {
	HTTP    *http.Client
	BaseURL string
}

// NewClient creates a new API client.
func NewClient(baseURL string) *Client {
	return &Client{
		HTTP:    &http.Client{Timeout: defaultTimeout},
		BaseURL: baseURL,
	}
}

// Get performs a GET request and decodes the JSON response into result.
func (c *Client) Get(path string, result any) error {
	resp, err := c.do(http.MethodGet, path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}
	return nil
}

// GetRaw performs a GET request and returns the raw response body.
func (c *Client) GetRaw(path string) ([]byte, error) {
	resp, err := c.do(http.MethodGet, path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (c *Client) do(method, path string) (*http.Response, error) {
	url := c.BaseURL + path

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("リクエスト作成エラー: %w", err)
	}

	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Accept", "application/json")

	debugLogRequest(req)
	start := time.Now()

	var resp *http.Response
	for attempt := 0; attempt <= maxRetries; attempt++ {
		resp, err = c.HTTP.Do(req)
		if err != nil {
			if attempt == maxRetries {
				return nil, &cerrors.NetworkError{Err: err}
			}
			time.Sleep(time.Duration(attempt+1) * time.Second)
			continue
		}

		if resp.StatusCode == 429 || resp.StatusCode >= 500 {
			if attempt < maxRetries {
				resp.Body.Close()
				time.Sleep(time.Duration(attempt+1) * time.Second)
				continue
			}
		}
		break
	}

	elapsed := time.Since(start)
	if resp == nil {
		return nil, &cerrors.NetworkError{Err: fmt.Errorf("リトライ後もレスポンスなし")}
	}
	debugLogResponse(resp, elapsed)

	if resp.StatusCode >= 400 {
		return resp, parseAPIError(resp)
	}

	return resp, nil
}

func parseAPIError(resp *http.Response) error {
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var errInfo model.ErrorInfo
	if json.Unmarshal(body, &errInfo) == nil && errInfo.Message != "" {
		if resp.StatusCode == 404 {
			return &cerrors.NotFoundError{Resource: "法令", ID: errInfo.Message}
		}
		return &cerrors.APIError{
			StatusCode: resp.StatusCode,
			Code:       errInfo.Code,
			Message:    errInfo.Message,
		}
	}

	if resp.StatusCode == 404 {
		return &cerrors.NotFoundError{Resource: "法令", ID: ""}
	}

	return &cerrors.APIError{
		StatusCode: resp.StatusCode,
		Message:    string(body),
	}
}
```

- [ ] **Step 5: Run tests**

Run: `cd /root/dev/planitai/jplaw-cli && go test ./internal/api/ -v`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add internal/api/
git commit -m "feat: add base API client with retry and debug logging"
```

---

### Task 7: API Service Methods

**Files:**
- Create: `internal/api/laws.go`
- Create: `internal/api/keyword.go`
- Create: `internal/api/lawdata.go`
- Create: `internal/api/revisions.go`
- Create: `internal/api/attachment.go`
- Create: `internal/api/laws_test.go`

- [ ] **Step 1: Write the failing test**

Create `internal/api/laws_test.go`:

```go
package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/planitaicojp/jplaw-cli/internal/model"
)

func TestListLaws(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/laws" {
			t.Errorf("path = %q, want /laws", r.URL.Path)
		}
		if r.URL.Query().Get("law_type") != "Act" {
			t.Errorf("law_type = %q, want Act", r.URL.Query().Get("law_type"))
		}
		if r.URL.Query().Get("limit") != "10" {
			t.Errorf("limit = %q, want 10", r.URL.Query().Get("limit"))
		}
		resp := model.LawsResponse{
			TotalCount: 1,
			Count:      1,
			Laws: []model.LawItem{
				{
					LawInfo: model.LawInfo{
						LawID:  "test-id",
						LawNum: "テスト法令番号",
					},
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	params := &LawsParams{
		LawType: []model.LawType{model.LawTypeAct},
		Limit:   10,
	}
	resp, err := client.ListLaws(params)
	if err != nil {
		t.Fatalf("ListLaws() error: %v", err)
	}
	if resp.TotalCount != 1 {
		t.Errorf("TotalCount = %d, want 1", resp.TotalCount)
	}
	if resp.Laws[0].LawInfo.LawID != "test-id" {
		t.Errorf("LawID = %q, want %q", resp.Laws[0].LawInfo.LawID, "test-id")
	}
}

func TestSearchKeyword(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/keyword" {
			t.Errorf("path = %q, want /keyword", r.URL.Path)
		}
		if r.URL.Query().Get("keyword") != "個人情報" {
			t.Errorf("keyword = %q, want 個人情報", r.URL.Query().Get("keyword"))
		}
		resp := model.KeywordResponse{
			TotalCount:    1,
			SentenceCount: 1,
			Items: []model.KeywordItem{
				{
					LawInfo: model.LawInfo{LawID: "test-id"},
					Sentences: []model.Sentence{
						{Text: "個人情報の保護"},
					},
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	params := &KeywordParams{
		Keyword: "個人情報",
	}
	resp, err := client.SearchKeyword(params)
	if err != nil {
		t.Fatalf("SearchKeyword() error: %v", err)
	}
	if resp.TotalCount != 1 {
		t.Errorf("TotalCount = %d, want 1", resp.TotalCount)
	}
}

func TestGetRevisions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/law_revisions/test-id" {
			t.Errorf("path = %q, want /law_revisions/test-id", r.URL.Path)
		}
		resp := model.RevisionsResponse{
			LawInfo: model.LawInfo{LawID: "test-id"},
			Revisions: []model.RevisionInfo{
				{LawRevisionID: "rev-1"},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	resp, err := client.GetRevisions("test-id", nil)
	if err != nil {
		t.Fatalf("GetRevisions() error: %v", err)
	}
	if len(resp.Revisions) != 1 {
		t.Errorf("Revisions count = %d, want 1", len(resp.Revisions))
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /root/dev/planitai/jplaw-cli && go test ./internal/api/ -v -run TestListLaws`
Expected: FAIL — `ListLaws` not defined

- [ ] **Step 3: Create laws.go**

Create `internal/api/laws.go`:

```go
package api

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/planitaicojp/jplaw-cli/internal/model"
)

// LawsParams are the query parameters for GET /laws.
type LawsParams struct {
	LawID             string
	LawNum            string
	LawNumEra         model.Era
	LawTitle          string
	LawType           []model.LawType
	CategoryCd        []string
	PromulgationFrom  string
	PromulgationTo    string
	RepealStatus      []model.RepealStatus
	Limit             int
	Offset            int
	Order             string
}

// ListLaws calls GET /laws with the given parameters.
func (c *Client) ListLaws(params *LawsParams) (*model.LawsResponse, error) {
	q := url.Values{}
	q.Set("response_format", "json")

	if params != nil {
		if params.LawID != "" {
			q.Set("law_id", params.LawID)
		}
		if params.LawNum != "" {
			q.Set("law_num", params.LawNum)
		}
		if params.LawNumEra != "" {
			q.Set("law_num_era", string(params.LawNumEra))
		}
		if params.LawTitle != "" {
			q.Set("law_title", params.LawTitle)
		}
		for _, lt := range params.LawType {
			q.Add("law_type", string(lt))
		}
		for _, cd := range params.CategoryCd {
			q.Add("category_cd", cd)
		}
		if params.PromulgationFrom != "" {
			q.Set("promulgation_date_from", params.PromulgationFrom)
		}
		if params.PromulgationTo != "" {
			q.Set("promulgation_date_to", params.PromulgationTo)
		}
		for _, rs := range params.RepealStatus {
			q.Add("repeal_status", string(rs))
		}
		if params.Limit > 0 {
			q.Set("limit", strconv.Itoa(params.Limit))
		}
		if params.Offset > 0 {
			q.Set("offset", strconv.Itoa(params.Offset))
		}
		if params.Order != "" {
			q.Set("order", params.Order)
		}
	}

	path := fmt.Sprintf("/laws?%s", q.Encode())
	var resp model.LawsResponse
	if err := c.Get(path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
```

- [ ] **Step 4: Create keyword.go**

Create `internal/api/keyword.go`:

```go
package api

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/planitaicojp/jplaw-cli/internal/model"
)

// KeywordParams are the query parameters for GET /keyword.
type KeywordParams struct {
	Keyword          string
	LawType          []model.LawType
	LawNumEra        model.Era
	Asof             string
	CategoryCd       []string
	PromulgationFrom string
	PromulgationTo   string
	Limit            int
	Offset           int
	Order            string
}

// SearchKeyword calls GET /keyword with the given parameters.
func (c *Client) SearchKeyword(params *KeywordParams) (*model.KeywordResponse, error) {
	q := url.Values{}
	q.Set("response_format", "json")

	if params != nil {
		q.Set("keyword", params.Keyword)
		for _, lt := range params.LawType {
			q.Add("law_type", string(lt))
		}
		if params.LawNumEra != "" {
			q.Set("law_num_era", string(params.LawNumEra))
		}
		if params.Asof != "" {
			q.Set("asof", params.Asof)
		}
		for _, cd := range params.CategoryCd {
			q.Add("category_cd", cd)
		}
		if params.PromulgationFrom != "" {
			q.Set("promulgation_date_from", params.PromulgationFrom)
		}
		if params.PromulgationTo != "" {
			q.Set("promulgation_date_to", params.PromulgationTo)
		}
		if params.Limit > 0 {
			q.Set("limit", strconv.Itoa(params.Limit))
		}
		if params.Offset > 0 {
			q.Set("offset", strconv.Itoa(params.Offset))
		}
		if params.Order != "" {
			q.Set("order", params.Order)
		}
	}

	path := fmt.Sprintf("/keyword?%s", q.Encode())
	var resp model.KeywordResponse
	if err := c.Get(path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
```

- [ ] **Step 5: Create lawdata.go**

Create `internal/api/lawdata.go`:

```go
package api

import (
	"fmt"
	"net/url"

	"github.com/planitaicojp/jplaw-cli/internal/model"
)

// LawDataParams are the query parameters for GET /law_data.
type LawDataParams struct {
	Asof          string
	Elm           string
	JsonFormat    string // "full" or "light"
	LawFullTextFormat string // "json" or "xml"
}

// GetLawData calls GET /law_data/{id} with the given parameters.
func (c *Client) GetLawData(idOrNum string, params *LawDataParams) (*model.LawDataResponse, error) {
	q := url.Values{}
	q.Set("response_format", "json")
	q.Set("law_full_text_format", "json")
	q.Set("json_format", "full")

	if params != nil {
		if params.Asof != "" {
			q.Set("asof", params.Asof)
		}
		if params.Elm != "" {
			q.Set("elm", params.Elm)
		}
		if params.JsonFormat != "" {
			q.Set("json_format", params.JsonFormat)
		}
		if params.LawFullTextFormat != "" {
			q.Set("law_full_text_format", params.LawFullTextFormat)
		}
	}

	path := fmt.Sprintf("/law_data/%s?%s", url.PathEscape(idOrNum), q.Encode())
	var resp model.LawDataResponse
	if err := c.Get(path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetLawFile calls GET /law_file/{type}/{id} and returns the raw file bytes.
func (c *Client) GetLawFile(fileType, idOrNum, asof string) ([]byte, error) {
	q := url.Values{}
	if asof != "" {
		q.Set("asof", asof)
	}

	path := fmt.Sprintf("/law_file/%s/%s", url.PathEscape(fileType), url.PathEscape(idOrNum))
	if len(q) > 0 {
		path += "?" + q.Encode()
	}
	return c.GetRaw(path)
}
```

- [ ] **Step 6: Create revisions.go**

Create `internal/api/revisions.go`:

```go
package api

import (
	"fmt"
	"net/url"

	"github.com/planitaicojp/jplaw-cli/internal/model"
)

// RevisionsParams are the query parameters for GET /law_revisions.
type RevisionsParams struct {
	AmendmentDateFrom string
	AmendmentDateTo   string
}

// GetRevisions calls GET /law_revisions/{id} with the given parameters.
func (c *Client) GetRevisions(idOrNum string, params *RevisionsParams) (*model.RevisionsResponse, error) {
	q := url.Values{}
	q.Set("response_format", "json")

	if params != nil {
		if params.AmendmentDateFrom != "" {
			q.Set("amendment_date_from", params.AmendmentDateFrom)
		}
		if params.AmendmentDateTo != "" {
			q.Set("amendment_date_to", params.AmendmentDateTo)
		}
	}

	path := fmt.Sprintf("/law_revisions/%s?%s", url.PathEscape(idOrNum), q.Encode())
	var resp model.RevisionsResponse
	if err := c.Get(path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
```

- [ ] **Step 7: Create attachment.go**

Create `internal/api/attachment.go`:

```go
package api

import (
	"fmt"
	"net/url"
)

// GetAttachment calls GET /attachment/{revision_id} and returns the raw bytes.
func (c *Client) GetAttachment(revisionID, src string) ([]byte, error) {
	q := url.Values{}
	if src != "" {
		q.Set("src", src)
	}

	path := fmt.Sprintf("/attachment/%s", url.PathEscape(revisionID))
	if len(q) > 0 {
		path += "?" + q.Encode()
	}
	return c.GetRaw(path)
}
```

- [ ] **Step 8: Run tests**

Run: `cd /root/dev/planitai/jplaw-cli && go test ./internal/api/ -v`
Expected: PASS

- [ ] **Step 9: Commit**

```bash
git add internal/api/
git commit -m "feat: add API service methods (laws, keyword, lawdata, revisions, attachment)"
```

---

### Task 8: Command Utilities (cmdutil)

**Files:**
- Create: `cmd/cmdutil/client.go`
- Create: `cmd/cmdutil/format.go`
- Create: `cmd/cmdutil/args.go`

- [ ] **Step 1: Create args.go**

Create `cmd/cmdutil/args.go`:

```go
package cmdutil

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ExactArgs returns a PositionalArgs that reports the command's Use line on mismatch.
func ExactArgs(n int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != n {
			return fmt.Errorf("%d個の引数が必要ですが、%d個指定されました\n\nUsage:\n  %s", n, len(args), cmd.UseLine())
		}
		return nil
	}
}
```

- [ ] **Step 2: Create client.go**

Create `cmd/cmdutil/client.go`:

```go
package cmdutil

import (
	"github.com/planitaicojp/jplaw-cli/internal/api"
	"github.com/planitaicojp/jplaw-cli/internal/config"
)

// NewClient creates an API client from configuration.
func NewClient() (*api.Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	baseURL := config.EnvOr(config.EnvBaseURL, cfg.BaseURL)
	return api.NewClient(baseURL), nil
}
```

- [ ] **Step 3: Create format.go**

Create `cmd/cmdutil/format.go`:

```go
package cmdutil

import (
	"github.com/spf13/cobra"

	"github.com/planitaicojp/jplaw-cli/internal/config"
)

// GetFormat returns the output format, respecting flag > env > config > default.
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
```

- [ ] **Step 4: Add cobra dependency**

Run: `cd /root/dev/planitai/jplaw-cli && go get github.com/spf13/cobra@v1.10.2`

- [ ] **Step 5: Commit**

```bash
git add cmd/cmdutil/ go.mod go.sum
git commit -m "feat: add command utilities (client, format, args)"
```

---

### Task 9: Root Command, Version, Completion

**Files:**
- Create: `cmd/root.go`
- Create: `cmd/version.go`
- Create: `cmd/completion.go`

- [ ] **Step 1: Create root.go**

Create `cmd/root.go`:

```go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/jplaw-cli/cmd/attachment"
	"github.com/planitaicojp/jplaw-cli/cmd/get"
	"github.com/planitaicojp/jplaw-cli/cmd/history"
	"github.com/planitaicojp/jplaw-cli/cmd/list"
	"github.com/planitaicojp/jplaw-cli/cmd/search"
	"github.com/planitaicojp/jplaw-cli/internal/api"
	cerrors "github.com/planitaicojp/jplaw-cli/internal/errors"
)

var (
	version = "dev"

	flagFormat  string
	flagVerbose bool
	flagNoColor bool
	flagNoInput bool
)

var rootCmd = &cobra.Command{
	Use:           "jplaw",
	Short:         "e-Gov法令API CLI",
	Long:          "e-Gov法令APIを操作するためのコマンドラインツール",
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		api.UserAgent = "jplaw-cli/" + version
		if flagVerbose {
			api.SetDebugLevel(api.DebugVerbose)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flagFormat, "format", "", "出力フォーマット: table, json, text")
	rootCmd.PersistentFlags().BoolVar(&flagVerbose, "verbose", false, "詳細出力（HTTPデバッグ）")
	rootCmd.PersistentFlags().BoolVar(&flagNoColor, "no-color", false, "色出力を無効化")
	rootCmd.PersistentFlags().BoolVar(&flagNoInput, "no-input", false, "対話プロンプトを無効化")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(search.Cmd)
	rootCmd.AddCommand(list.Cmd)
	rootCmd.AddCommand(get.Cmd)
	rootCmd.AddCommand(history.Cmd)
	rootCmd.AddCommand(attachment.Cmd)
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(cerrors.GetExitCode(err))
	}
}
```

- [ ] **Step 2: Create version.go**

Create `cmd/version.go`:

```go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "バージョン情報を表示",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("jplaw version %s\n", version)
		fmt.Println("e-Gov法令API (https://laws.e-gov.go.jp/api/2) CLI")
	},
}
```

- [ ] **Step 3: Create completion.go**

Create `cmd/completion.go`:

```go
package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/jplaw-cli/cmd/cmdutil"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "シェル補完スクリプトを生成",
	Long: `指定されたシェルの補完スクリプトを生成します。

使用例:
  # Bash
  jplaw completion bash > /etc/bash_completion.d/jplaw

  # Zsh
  jplaw completion zsh > "${fpath[1]}/_jplaw"

  # Fish
  jplaw completion fish > ~/.config/fish/completions/jplaw.fish`,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cmdutil.ExactArgs(1),
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			return rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			return rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			return rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		default:
			return nil
		}
	},
}
```

- [ ] **Step 4: Note — do not compile yet**

root.go imports subcommand packages (search, list, get, history, attachment) that don't exist yet. These are created in Tasks 10-14. The project will compile after all subcommand tasks are done.

- [ ] **Step 5: Commit**

```bash
git add cmd/root.go cmd/version.go cmd/completion.go
git commit -m "feat: add root command, version, and shell completion"
```

---

### Task 10: search Subcommand

**Files:**
- Create: `cmd/search/search.go`

- [ ] **Step 1: Create search.go**

Create `cmd/search/search.go`:

```go
package search

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/jplaw-cli/cmd/cmdutil"
	"github.com/planitaicojp/jplaw-cli/internal/api"
	"github.com/planitaicojp/jplaw-cli/internal/model"
	"github.com/planitaicojp/jplaw-cli/internal/output"
)

var (
	flagLawType  []string
	flagEra      string
	flagAsof     string
	flagCategory []string
	flagLimit    int
	flagOffset   int
)

// Cmd is the search subcommand.
var Cmd = &cobra.Command{
	Use:   "search <キーワード>",
	Short: "法令をキーワードで全文検索",
	Long:  "法令本文をキーワードで全文検索します。ワイルドカード、AND、OR、NOT検索に対応。",
	Args:  cmdutil.ExactArgs(1),
	RunE:  run,
}

func init() {
	Cmd.Flags().StringSliceVar(&flagLawType, "law-type", nil, "法令種別（法律,政令,府省令等、複数可）")
	Cmd.Flags().StringVar(&flagEra, "era", "", "時代（明治,大正,昭和,平成,令和）")
	Cmd.Flags().StringVar(&flagAsof, "asof", "", "時点指定（YYYY-MM-DD）")
	Cmd.Flags().StringSliceVar(&flagCategory, "category", nil, "分類コード")
	Cmd.Flags().IntVar(&flagLimit, "limit", 100, "結果数上限")
	Cmd.Flags().IntVar(&flagOffset, "offset", 0, "開始位置")
}

func run(cmd *cobra.Command, args []string) error {
	client, err := cmdutil.NewClient()
	if err != nil {
		return err
	}

	params := &api.KeywordParams{
		Keyword:    args[0],
		Limit:      flagLimit,
		Offset:     flagOffset,
		Asof:       flagAsof,
		CategoryCd: flagCategory,
	}

	for _, lt := range flagLawType {
		apiType := model.LawTypeLabelFromAPI(lt)
		if apiType == "" {
			params.LawType = append(params.LawType, model.LawType(lt))
		} else {
			params.LawType = append(params.LawType, apiType)
		}
	}

	if flagEra != "" {
		era := model.EraFromLabel(flagEra)
		if era == "" {
			params.LawNumEra = model.Era(flagEra)
		} else {
			params.LawNumEra = era
		}
	}

	resp, err := client.SearchKeyword(params)
	if err != nil {
		return err
	}

	format := cmdutil.GetFormat(cmd, "table")

	if format == "json" {
		return output.New("json").Format(os.Stdout, resp)
	}

	if len(resp.Items) == 0 {
		fmt.Fprintln(os.Stderr, "検索結果がありません")
		return nil
	}

	rows := make([]model.KeywordListRow, 0)
	for _, item := range resp.Items {
		for _, s := range item.Sentences {
			text := stripHTML(s.Text)
			rows = append(rows, model.KeywordListRow{
				LawNum:   item.LawInfo.LawNum,
				LawTitle: item.RevisionInfo.LawTitle,
				Match:    text,
			})
		}
	}

	fmt.Fprintf(os.Stderr, "検索結果: %d件\n", resp.TotalCount)
	return output.New("table").Format(os.Stdout, rows)
}

func stripHTML(s string) string {
	result := s
	for {
		start := strings.Index(result, "<")
		if start == -1 {
			break
		}
		end := strings.Index(result[start:], ">")
		if end == -1 {
			break
		}
		result = result[:start] + result[start+end+1:]
	}
	return result
}
```

- [ ] **Step 2: Commit**

```bash
git add cmd/search/
git commit -m "feat: add search subcommand for keyword search"
```

---

### Task 11: list Subcommand

**Files:**
- Create: `cmd/list/list.go`

- [ ] **Step 1: Create list.go**

Create `cmd/list/list.go`:

```go
package list

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/jplaw-cli/cmd/cmdutil"
	"github.com/planitaicojp/jplaw-cli/internal/api"
	"github.com/planitaicojp/jplaw-cli/internal/model"
	"github.com/planitaicojp/jplaw-cli/internal/output"
)

var (
	flagLawType         []string
	flagEra             string
	flagTitle           string
	flagPromulgationFrom string
	flagPromulgationTo   string
	flagCategory        []string
	flagRepealStatus    []string
	flagLimit           int
	flagOffset          int
)

// Cmd is the list subcommand.
var Cmd = &cobra.Command{
	Use:   "list",
	Short: "法令一覧を取得",
	Long:  "条件に一致する法令の一覧を取得します。",
	RunE:  run,
}

func init() {
	Cmd.Flags().StringSliceVar(&flagLawType, "law-type", nil, "法令種別（法律,政令,府省令等、複数可）")
	Cmd.Flags().StringVar(&flagEra, "era", "", "時代（明治,大正,昭和,平成,令和）")
	Cmd.Flags().StringVar(&flagTitle, "title", "", "法令名（部分一致）")
	Cmd.Flags().StringVar(&flagPromulgationFrom, "promulgation-from", "", "公布日FROM（YYYY-MM-DD）")
	Cmd.Flags().StringVar(&flagPromulgationTo, "promulgation-to", "", "公布日TO（YYYY-MM-DD）")
	Cmd.Flags().StringSliceVar(&flagCategory, "category", nil, "分類コード")
	Cmd.Flags().StringSliceVar(&flagRepealStatus, "repeal-status", nil, "廃止状態")
	Cmd.Flags().IntVar(&flagLimit, "limit", 100, "結果数上限")
	Cmd.Flags().IntVar(&flagOffset, "offset", 0, "開始位置")
}

func run(cmd *cobra.Command, args []string) error {
	client, err := cmdutil.NewClient()
	if err != nil {
		return err
	}

	params := &api.LawsParams{
		LawTitle:         flagTitle,
		PromulgationFrom: flagPromulgationFrom,
		PromulgationTo:   flagPromulgationTo,
		CategoryCd:       flagCategory,
		Limit:            flagLimit,
		Offset:           flagOffset,
	}

	for _, lt := range flagLawType {
		apiType := model.LawTypeLabelFromAPI(lt)
		if apiType == "" {
			params.LawType = append(params.LawType, model.LawType(lt))
		} else {
			params.LawType = append(params.LawType, apiType)
		}
	}

	if flagEra != "" {
		era := model.EraFromLabel(flagEra)
		if era == "" {
			params.LawNumEra = model.Era(flagEra)
		} else {
			params.LawNumEra = era
		}
	}

	for _, rs := range flagRepealStatus {
		params.RepealStatus = append(params.RepealStatus, model.RepealStatus(rs))
	}

	resp, err := client.ListLaws(params)
	if err != nil {
		return err
	}

	format := cmdutil.GetFormat(cmd, "table")

	if format == "json" {
		return output.New("json").Format(os.Stdout, resp)
	}

	if len(resp.Laws) == 0 {
		fmt.Fprintln(os.Stderr, "該当する法令がありません")
		return nil
	}

	rows := make([]model.LawListRow, len(resp.Laws))
	for i, law := range resp.Laws {
		rows[i] = model.LawListRow{
			LawID:        law.LawInfo.LawID,
			LawNum:       law.LawInfo.LawNum,
			LawTitle:     law.RevisionInfo.LawTitle,
			LawType:      law.LawInfo.LawType.Label(),
			Promulgation: law.LawInfo.Promulgation,
		}
	}

	fmt.Fprintf(os.Stderr, "法令一覧: %d/%d件\n", resp.Count, resp.TotalCount)
	return output.New("table").Format(os.Stdout, rows)
}
```

- [ ] **Step 2: Commit**

```bash
git add cmd/list/
git commit -m "feat: add list subcommand for law listing"
```

---

### Task 12: get Subcommand

**Files:**
- Create: `cmd/get/get.go`

- [ ] **Step 1: Create get.go**

Create `cmd/get/get.go`:

```go
package get

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/jplaw-cli/cmd/cmdutil"
	"github.com/planitaicojp/jplaw-cli/internal/api"
	"github.com/planitaicojp/jplaw-cli/internal/output"
)

var (
	flagAsof   string
	flagElm    string
	flagFile   string
	flagOutput string
)

// Cmd is the get subcommand.
var Cmd = &cobra.Command{
	Use:   "get <法令IDまたは法令番号>",
	Short: "法令本文を取得",
	Long:  "法令IDまたは法令番号を指定して法令本文を取得します。",
	Args:  cmdutil.ExactArgs(1),
	RunE:  run,
}

func init() {
	Cmd.Flags().StringVar(&flagAsof, "asof", "", "時点指定（YYYY-MM-DD）")
	Cmd.Flags().StringVar(&flagElm, "elm", "", "特定条文指定（例: 第一条, MainProvision）")
	Cmd.Flags().StringVar(&flagFile, "file", "", "ファイルダウンロード（xml,json,html,rtf,docx）")
	Cmd.Flags().StringVarP(&flagOutput, "output", "o", "", "ファイル保存先パス")
}

func run(cmd *cobra.Command, args []string) error {
	client, err := cmdutil.NewClient()
	if err != nil {
		return err
	}

	idOrNum := args[0]

	// File download mode
	if flagFile != "" {
		return downloadFile(client, idOrNum)
	}

	// Law data mode
	params := &api.LawDataParams{
		Asof: flagAsof,
		Elm:  flagElm,
	}

	resp, err := client.GetLawData(idOrNum, params)
	if err != nil {
		return err
	}

	format := cmdutil.GetFormat(cmd, "text")

	switch format {
	case "json":
		return output.New("json").Format(os.Stdout, resp)
	case "xml":
		// Re-fetch with XML format
		params.LawFullTextFormat = "xml"
		xmlResp, err := client.GetLawData(idOrNum, params)
		if err != nil {
			return err
		}
		return output.New("json").Format(os.Stdout, xmlResp)
	default:
		// Text output: print law title + full text as readable text
		fmt.Fprintf(os.Stdout, "%s\n", resp.RevisionInfo.LawTitle)
		fmt.Fprintf(os.Stdout, "（%s）\n\n", resp.LawInfo.LawNum)

		// For now, output raw JSON full text
		// TODO: lawtext converter in Task 15
		var prettyJSON json.RawMessage
		if err := json.Unmarshal(resp.LawFullText, &prettyJSON); err == nil {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			enc.SetEscapeHTML(false)
			return enc.Encode(prettyJSON)
		}
		_, err = os.Stdout.Write(resp.LawFullText)
		return err
	}
}

func downloadFile(client *api.Client, idOrNum string) error {
	data, err := client.GetLawFile(flagFile, idOrNum, flagAsof)
	if err != nil {
		return err
	}

	if flagOutput == "" {
		_, err = os.Stdout.Write(data)
		return err
	}

	if err := os.WriteFile(flagOutput, data, 0644); err != nil {
		return fmt.Errorf("ファイル書き込みエラー: %w", err)
	}
	fmt.Fprintf(os.Stderr, "保存しました: %s (%d bytes)\n", flagOutput, len(data))
	return nil
}
```

- [ ] **Step 2: Commit**

```bash
git add cmd/get/
git commit -m "feat: add get subcommand for law text retrieval"
```

---

### Task 13: history Subcommand

**Files:**
- Create: `cmd/history/history.go`

- [ ] **Step 1: Create history.go**

Create `cmd/history/history.go`:

```go
package history

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/jplaw-cli/cmd/cmdutil"
	"github.com/planitaicojp/jplaw-cli/internal/api"
	"github.com/planitaicojp/jplaw-cli/internal/model"
	"github.com/planitaicojp/jplaw-cli/internal/output"
)

var (
	flagAmendmentFrom string
	flagAmendmentTo   string
)

// Cmd is the history subcommand.
var Cmd = &cobra.Command{
	Use:   "history <法令IDまたは法令番号>",
	Short: "法令の改正履歴を取得",
	Long:  "法令の改正履歴一覧を取得します。",
	Args:  cmdutil.ExactArgs(1),
	RunE:  run,
}

func init() {
	Cmd.Flags().StringVar(&flagAmendmentFrom, "amendment-from", "", "改正施行日FROM（YYYY-MM-DD）")
	Cmd.Flags().StringVar(&flagAmendmentTo, "amendment-to", "", "改正施行日TO（YYYY-MM-DD）")
}

func run(cmd *cobra.Command, args []string) error {
	client, err := cmdutil.NewClient()
	if err != nil {
		return err
	}

	params := &api.RevisionsParams{
		AmendmentDateFrom: flagAmendmentFrom,
		AmendmentDateTo:   flagAmendmentTo,
	}

	resp, err := client.GetRevisions(args[0], params)
	if err != nil {
		return err
	}

	format := cmdutil.GetFormat(cmd, "table")

	if format == "json" {
		return output.New("json").Format(os.Stdout, resp)
	}

	if len(resp.Revisions) == 0 {
		fmt.Fprintln(os.Stderr, "改正履歴がありません")
		return nil
	}

	rows := make([]model.RevisionListRow, len(resp.Revisions))
	for i, rev := range resp.Revisions {
		rows[i] = model.RevisionListRow{
			AmendmentDate:  rev.AmendmentEnforcementDate,
			AmendmentTitle: rev.AmendmentLawTitle,
			AmendmentNum:   rev.AmendmentLawNum,
			AmendmentType:  string(rev.AmendmentType),
		}
	}

	fmt.Fprintf(os.Stderr, "%s の改正履歴: %d件\n", resp.LawInfo.LawNum, len(resp.Revisions))
	return output.New("table").Format(os.Stdout, rows)
}
```

- [ ] **Step 2: Commit**

```bash
git add cmd/history/
git commit -m "feat: add history subcommand for revision history"
```

---

### Task 14: attachment Subcommand

**Files:**
- Create: `cmd/attachment/attachment.go`

- [ ] **Step 1: Create attachment.go**

Create `cmd/attachment/attachment.go`:

```go
package attachment

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/jplaw-cli/cmd/cmdutil"
)

var (
	flagSrc    string
	flagOutput string
)

// Cmd is the attachment subcommand.
var Cmd = &cobra.Command{
	Use:   "attachment <リビジョンID>",
	Short: "添付ファイルをダウンロード",
	Long:  "法令リビジョンに添付されたファイルをダウンロードします。",
	Args:  cmdutil.ExactArgs(1),
	RunE:  run,
}

func init() {
	Cmd.Flags().StringVar(&flagSrc, "src", "", "特定添付ファイルのsrc属性値")
	Cmd.Flags().StringVarP(&flagOutput, "output", "o", "", "保存先パス")
}

func run(cmd *cobra.Command, args []string) error {
	client, err := cmdutil.NewClient()
	if err != nil {
		return err
	}

	data, err := client.GetAttachment(args[0], flagSrc)
	if err != nil {
		return err
	}

	if flagOutput == "" {
		_, err = os.Stdout.Write(data)
		return err
	}

	if err := os.WriteFile(flagOutput, data, 0644); err != nil {
		return fmt.Errorf("ファイル書き込みエラー: %w", err)
	}
	fmt.Fprintf(os.Stderr, "保存しました: %s (%d bytes)\n", flagOutput, len(data))
	return nil
}
```

- [ ] **Step 2: Commit**

```bash
git add cmd/attachment/
git commit -m "feat: add attachment subcommand for file download"
```

---

### Task 15: Build and Smoke Test

**Files:**
- Modify: `main.go` (already created)

- [ ] **Step 1: Tidy dependencies**

Run: `cd /root/dev/planitai/jplaw-cli && go mod tidy`

- [ ] **Step 2: Build**

Run: `cd /root/dev/planitai/jplaw-cli && make build`
Expected: `jplaw` binary created

- [ ] **Step 3: Verify version command**

Run: `cd /root/dev/planitai/jplaw-cli && ./jplaw version`
Expected: `jplaw version dev`

- [ ] **Step 4: Verify help output**

Run: `cd /root/dev/planitai/jplaw-cli && ./jplaw --help`
Expected: Shows all subcommands (search, list, get, history, attachment, version, completion)

- [ ] **Step 5: Verify search subcommand help**

Run: `cd /root/dev/planitai/jplaw-cli && ./jplaw search --help`
Expected: Shows search flags (--law-type, --era, --asof, --limit, etc.)

- [ ] **Step 6: Smoke test with real API**

Run: `cd /root/dev/planitai/jplaw-cli && ./jplaw list --law-type 憲法 --limit 5`
Expected: Table with Constitution entries

- [ ] **Step 7: Run all tests**

Run: `cd /root/dev/planitai/jplaw-cli && make test`
Expected: All tests PASS

- [ ] **Step 8: Commit**

```bash
git add go.mod go.sum
git commit -m "feat: complete initial build with all subcommands"
```

---

### Task 16: Lawtext Converter (Phase 1 — Basic Structure)

**Files:**
- Create: `internal/lawtext/structure.go`
- Create: `internal/lawtext/converter.go`
- Create: `internal/lawtext/converter_test.go`
- Modify: `cmd/get/get.go`

- [ ] **Step 1: Write the failing test**

Create `internal/lawtext/converter_test.go`:

```go
package lawtext

import (
	"encoding/json"
	"testing"
)

func TestConvertSimpleArticle(t *testing.T) {
	// Minimal JSON full format representing a simple article
	raw := json.RawMessage(`{
		"tag": "Law",
		"attr": {},
		"children": [
			{
				"tag": "LawBody",
				"attr": {},
				"children": [
					{
						"tag": "LawTitle",
						"attr": {},
						"children": ["テスト法"]
					},
					{
						"tag": "MainProvision",
						"attr": {},
						"children": [
							{
								"tag": "Article",
								"attr": {"Num": "1"},
								"children": [
									{
										"tag": "ArticleCaption",
										"attr": {},
										"children": ["（目的）"]
									},
									{
										"tag": "ArticleTitle",
										"attr": {},
										"children": ["第一条"]
									},
									{
										"tag": "Paragraph",
										"attr": {"Num": "1"},
										"children": [
											{
												"tag": "ParagraphNum",
												"attr": {},
												"children": []
											},
											{
												"tag": "ParagraphSentence",
												"attr": {},
												"children": [
													{
														"tag": "Sentence",
														"attr": {},
														"children": ["この法律は、テストを目的とする。"]
													}
												]
											}
										]
									}
								]
							}
						]
					}
				]
			}
		]
	}`)

	result, err := Convert(raw)
	if err != nil {
		t.Fatalf("Convert() error: %v", err)
	}

	if result == "" {
		t.Fatal("Convert() returned empty string")
	}

	// Should contain the article title
	if !contains(result, "第一条") {
		t.Errorf("result missing '第一条', got:\n%s", result)
	}

	// Should contain the sentence
	if !contains(result, "この法律は、テストを目的とする。") {
		t.Errorf("result missing sentence, got:\n%s", result)
	}
}

func TestConvertEmpty(t *testing.T) {
	raw := json.RawMessage(`{}`)
	result, err := Convert(raw)
	if err != nil {
		t.Fatalf("Convert() error: %v", err)
	}
	if result != "" {
		t.Errorf("expected empty string, got: %q", result)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsStr(s, substr))
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd /root/dev/planitai/jplaw-cli && go test ./internal/lawtext/ -v`
Expected: FAIL — package does not exist

- [ ] **Step 3: Create structure.go**

Create `internal/lawtext/structure.go`:

```go
package lawtext

import "encoding/json"

// Node represents a node in the law JSON full format.
type Node struct {
	Tag      string          `json:"tag"`
	Attr     map[string]string `json:"attr"`
	Children []json.RawMessage `json:"children"`
}

// parseNode parses a JSON raw message into a Node.
// Children can be either strings or nested Node objects.
func parseNode(raw json.RawMessage) (*Node, error) {
	var node Node
	if err := json.Unmarshal(raw, &node); err != nil {
		return nil, err
	}
	return &node, nil
}

// childText extracts all text content from a node's children recursively.
func childText(children []json.RawMessage) string {
	var result string
	for _, child := range children {
		// Try as string first
		var s string
		if json.Unmarshal(child, &s) == nil {
			result += s
			continue
		}
		// Try as node
		node, err := parseNode(child)
		if err != nil {
			continue
		}
		result += childText(node.Children)
	}
	return result
}
```

- [ ] **Step 4: Create converter.go**

Create `internal/lawtext/converter.go`:

```go
package lawtext

import (
	"encoding/json"
	"strings"
)

// Convert transforms law JSON full format into human-readable text.
func Convert(raw json.RawMessage) (string, error) {
	if len(raw) == 0 || string(raw) == "{}" || string(raw) == "null" {
		return "", nil
	}

	node, err := parseNode(raw)
	if err != nil {
		return "", err
	}

	var b strings.Builder
	renderNode(&b, node, 0)
	return b.String(), nil
}

func renderNode(b *strings.Builder, node *Node, depth int) {
	if node == nil {
		return
	}

	switch node.Tag {
	case "LawTitle":
		text := childText(node.Children)
		if text != "" {
			b.WriteString(text)
			b.WriteString("\n\n")
		}

	case "TOCLabel", "PartTitle", "ChapterTitle", "SectionTitle",
		"SubsectionTitle", "DivisionTitle":
		text := childText(node.Children)
		if text != "" {
			indent := strings.Repeat("  ", depth)
			b.WriteString(indent)
			b.WriteString(text)
			b.WriteString("\n\n")
		}

	case "ArticleTitle":
		text := childText(node.Children)
		if text != "" {
			b.WriteString(text)
		}

	case "ArticleCaption":
		text := childText(node.Children)
		if text != "" {
			b.WriteString(text)
			b.WriteString("\n")
		}

	case "ParagraphSentence", "ItemSentence", "Subitem1Sentence",
		"Subitem2Sentence", "Subitem3Sentence":
		text := childText(node.Children)
		if text != "" {
			indent := strings.Repeat("  ", depth)
			b.WriteString(indent)
			b.WriteString(text)
			b.WriteString("\n")
		}

	case "ParagraphNum":
		text := childText(node.Children)
		if text != "" {
			indent := strings.Repeat("  ", depth)
			b.WriteString(indent)
			b.WriteString(text)
			b.WriteString(" ")
		}

	case "ItemTitle":
		text := childText(node.Children)
		if text != "" {
			indent := strings.Repeat("    ", 1)
			b.WriteString(indent)
			b.WriteString(text)
			b.WriteString(" ")
		}

	case "Subitem1Title", "Subitem2Title", "Subitem3Title":
		text := childText(node.Children)
		if text != "" {
			indent := strings.Repeat("      ", 1)
			b.WriteString(indent)
			b.WriteString(text)
			b.WriteString(" ")
		}

	case "SupplProvisionLabel":
		text := childText(node.Children)
		if text != "" {
			b.WriteString("\n")
			b.WriteString(strings.Repeat("─", 40))
			b.WriteString("\n")
			b.WriteString(text)
			b.WriteString("\n\n")
		}

	case "Article":
		// Render children, then add newline after article
		for _, child := range node.Children {
			childNode, err := parseNode(child)
			if err != nil {
				continue
			}
			renderNode(b, childNode, depth)
		}
		b.WriteString("\n")
		return

	case "Paragraph":
		for _, child := range node.Children {
			childNode, err := parseNode(child)
			if err != nil {
				continue
			}
			renderNode(b, childNode, 1)
		}
		return

	case "Item":
		for _, child := range node.Children {
			childNode, err := parseNode(child)
			if err != nil {
				continue
			}
			renderNode(b, childNode, 2)
		}
		return

	case "Subitem1", "Subitem2", "Subitem3":
		for _, child := range node.Children {
			childNode, err := parseNode(child)
			if err != nil {
				continue
			}
			renderNode(b, childNode, 3)
		}
		return
	}

	// Default: recurse into children
	for _, child := range node.Children {
		// Skip string children (already handled by specific tags)
		var s string
		if json.Unmarshal(child, &s) == nil {
			continue
		}
		childNode, err := parseNode(child)
		if err != nil {
			continue
		}
		renderNode(b, childNode, depth)
	}
}
```

- [ ] **Step 5: Run tests**

Run: `cd /root/dev/planitai/jplaw-cli && go test ./internal/lawtext/ -v`
Expected: PASS

- [ ] **Step 6: Update cmd/get/get.go to use lawtext converter**

In `cmd/get/get.go`, replace the default case in the switch statement:

Replace the `default:` block (lines starting from `// Text output:` through the end of the default case) with:

```go
	default:
		// Text output: convert law JSON to readable text
		fmt.Fprintf(os.Stdout, "%s\n", resp.RevisionInfo.LawTitle)
		fmt.Fprintf(os.Stdout, "（%s）\n\n", resp.LawInfo.LawNum)

		text, err := lawtext.Convert(resp.LawFullText)
		if err != nil {
			// Fallback to raw JSON if conversion fails
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			enc.SetEscapeHTML(false)
			return enc.Encode(resp.LawFullText)
		}
		if text != "" {
			_, err = fmt.Fprint(os.Stdout, text)
		}
		return err
```

Add the import `"github.com/planitaicojp/jplaw-cli/internal/lawtext"` to the imports.

- [ ] **Step 7: Build and verify**

Run: `cd /root/dev/planitai/jplaw-cli && make build && make test`
Expected: Build succeeds, all tests pass

- [ ] **Step 8: Commit**

```bash
git add internal/lawtext/ cmd/get/get.go
git commit -m "feat: add lawtext converter for human-readable law text output"
```

---

### Task 17: End-to-End Verification

- [ ] **Step 1: Search test**

Run: `cd /root/dev/planitai/jplaw-cli && ./jplaw search "個人情報" --limit 3`
Expected: Table with matching laws and sentences

- [ ] **Step 2: List test**

Run: `cd /root/dev/planitai/jplaw-cli && ./jplaw list --law-type 法律 --limit 5`
Expected: Table with 5 laws

- [ ] **Step 3: Get test (text)**

Run: `cd /root/dev/planitai/jplaw-cli && ./jplaw get 405AC0000000088 --elm 第一条`
Expected: Readable text output of Article 1 of the Personal Information Protection Act

- [ ] **Step 4: Get test (json)**

Run: `cd /root/dev/planitai/jplaw-cli && ./jplaw get 405AC0000000088 --format json --elm 第一条`
Expected: JSON output

- [ ] **Step 5: History test**

Run: `cd /root/dev/planitai/jplaw-cli && ./jplaw history 405AC0000000088`
Expected: Table with revision history

- [ ] **Step 6: Verbose test**

Run: `cd /root/dev/planitai/jplaw-cli && ./jplaw list --law-type 憲法 --limit 1 --verbose`
Expected: HTTP debug output on stderr, result on stdout

- [ ] **Step 7: JSON pipe test**

Run: `cd /root/dev/planitai/jplaw-cli && ./jplaw list --law-type 憲法 --limit 1 --format json | python3 -m json.tool`
Expected: Pretty-printed JSON

- [ ] **Step 8: Run full test suite**

Run: `cd /root/dev/planitai/jplaw-cli && make test`
Expected: All tests PASS

- [ ] **Step 9: Final commit**

```bash
git add -A
git commit -m "chore: tidy up after end-to-end verification"
```
