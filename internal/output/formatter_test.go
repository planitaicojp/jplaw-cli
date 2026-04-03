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
	rows := []testRow{{ID: "1", Name: "テスト法"}}
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
	rows := []testRow{{ID: "1", Name: "テスト法"}, {ID: "2", Name: "サンプル法"}}
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
