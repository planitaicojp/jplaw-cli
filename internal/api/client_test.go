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
