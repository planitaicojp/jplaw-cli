package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAttachment(t *testing.T) {
	expected := []byte{0x89, 0x50, 0x4E, 0x47} // PNG magic bytes
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/attachment/rev-001" {
			t.Errorf("path = %q, want /attachment/rev-001", r.URL.Path)
		}
		if r.URL.Query().Get("src") != "" {
			t.Errorf("unexpected src param: %q", r.URL.Query().Get("src"))
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(expected)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	data, err := client.GetAttachment("rev-001", "")
	if err != nil {
		t.Fatalf("GetAttachment() error: %v", err)
	}
	if len(data) != len(expected) {
		t.Fatalf("data length = %d, want %d", len(data), len(expected))
	}
	for i, b := range expected {
		if data[i] != b {
			t.Errorf("data[%d] = %x, want %x", i, data[i], b)
		}
	}
}

func TestGetAttachmentWithSrc(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/attachment/rev-002" {
			t.Errorf("path = %q, want /attachment/rev-002", r.URL.Path)
		}
		if r.URL.Query().Get("src") != "image001.pdf" {
			t.Errorf("src = %q, want image001.pdf", r.URL.Query().Get("src"))
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write([]byte("pdf-content"))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	data, err := client.GetAttachment("rev-002", "image001.pdf")
	if err != nil {
		t.Fatalf("GetAttachment() error: %v", err)
	}
	if string(data) != "pdf-content" {
		t.Errorf("data = %q, want %q", string(data), "pdf-content")
	}
}
