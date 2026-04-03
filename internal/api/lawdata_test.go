package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	cerrors "github.com/planitaicojp/jplaw-cli/internal/errors"
	"github.com/planitaicojp/jplaw-cli/internal/model"
)

func TestGetLawData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/law_data/405AC0000000088" {
			t.Errorf("path = %q, want /law_data/405AC0000000088", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("response_format") != "json" {
			t.Errorf("response_format = %q, want json", q.Get("response_format"))
		}
		if q.Get("law_full_text_format") != "json" {
			t.Errorf("law_full_text_format = %q, want json", q.Get("law_full_text_format"))
		}
		if q.Get("json_format") != "full" {
			t.Errorf("json_format = %q, want full", q.Get("json_format"))
		}
		w.WriteHeader(http.StatusOK)
		resp := model.LawDataResponse{
			LawInfo:     model.LawInfo{LawID: "405AC0000000088", LawNum: "平成15年法律第88号"},
			LawFullText: json.RawMessage(`{"tag":"Law","children":[]}`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	resp, err := client.GetLawData("405AC0000000088", nil)
	if err != nil {
		t.Fatalf("GetLawData() error: %v", err)
	}
	if resp.LawInfo.LawID != "405AC0000000088" {
		t.Errorf("LawID = %q, want %q", resp.LawInfo.LawID, "405AC0000000088")
	}
	if resp.LawInfo.LawNum != "平成15年法律第88号" {
		t.Errorf("LawNum = %q, want %q", resp.LawInfo.LawNum, "平成15年法律第88号")
	}
}

func TestGetLawDataWithParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("asof") != "2024-01-01" {
			t.Errorf("asof = %q, want 2024-01-01", q.Get("asof"))
		}
		if q.Get("elm") != "第一条" {
			t.Errorf("elm = %q, want 第一条", q.Get("elm"))
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(model.LawDataResponse{
			LawInfo: model.LawInfo{LawID: "test-id"},
		})
	}))
	defer server.Close()

	client := NewClient(server.URL)
	params := &LawDataParams{
		Asof: "2024-01-01",
		Elm:  "第一条",
	}
	resp, err := client.GetLawData("test-id", params)
	if err != nil {
		t.Fatalf("GetLawData() error: %v", err)
	}
	if resp.LawInfo.LawID != "test-id" {
		t.Errorf("LawID = %q, want %q", resp.LawInfo.LawID, "test-id")
	}
}

func TestGetLawDataFormatOverrides(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("json_format") != "simple" {
			t.Errorf("json_format = %q, want simple", q.Get("json_format"))
		}
		if q.Get("law_full_text_format") != "xml" {
			t.Errorf("law_full_text_format = %q, want xml", q.Get("law_full_text_format"))
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(model.LawDataResponse{
			LawInfo: model.LawInfo{LawID: "test-id"},
		})
	}))
	defer server.Close()

	client := NewClient(server.URL)
	params := &LawDataParams{
		JsonFormat:        "simple",
		LawFullTextFormat: "xml",
	}
	resp, err := client.GetLawData("test-id", params)
	if err != nil {
		t.Fatalf("GetLawData() error: %v", err)
	}
	if resp.LawInfo.LawID != "test-id" {
		t.Errorf("LawID = %q, want %q", resp.LawInfo.LawID, "test-id")
	}
}

func TestGetLawDataWithAttachedFiles(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(model.LawDataResponse{
			LawInfo: model.LawInfo{LawID: "test-id"},
			AttachedFilesInfo: &model.AttachedFilesInfo{
				AttachedFiles: []model.AttachedFile{
					{LawRevisionID: "rev-1", Src: "image001.png", Updated: "2024-01-01"},
				},
			},
		})
	}))
	defer server.Close()

	client := NewClient(server.URL)
	resp, err := client.GetLawData("test-id", nil)
	if err != nil {
		t.Fatalf("GetLawData() error: %v", err)
	}
	if resp.AttachedFilesInfo == nil {
		t.Fatal("AttachedFilesInfo is nil")
	}
	if len(resp.AttachedFilesInfo.AttachedFiles) != 1 {
		t.Fatalf("AttachedFiles count = %d, want 1", len(resp.AttachedFilesInfo.AttachedFiles))
	}
	af := resp.AttachedFilesInfo.AttachedFiles[0]
	if af.Src != "image001.png" {
		t.Errorf("Src = %q, want %q", af.Src, "image001.png")
	}
}

func TestGetLawDataNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"code":    "NOT_FOUND",
			"message": "該当する法令が見つかりません",
		})
	}))
	defer server.Close()

	client := NewClient(server.URL)
	_, err := client.GetLawData("nonexistent", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var nfe *cerrors.NotFoundError
	if !isNotFoundError(err, &nfe) {
		t.Errorf("expected NotFoundError, got %T: %v", err, err)
	}
}

func TestGetLawFile(t *testing.T) {
	tests := []struct {
		name     string
		fileType string
		idOrNum  string
		asof     string
		wantPath string
		wantAsof bool
	}{
		{
			name:     "xml without asof",
			fileType: "xml",
			idOrNum:  "405AC0000000088",
			asof:     "",
			wantPath: "/law_file/xml/405AC0000000088",
			wantAsof: false,
		},
		{
			name:     "json with asof",
			fileType: "json",
			idOrNum:  "test-id",
			asof:     "2024-01-01",
			wantPath: "/law_file/json/test-id",
			wantAsof: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := []byte("<Law>テスト法令</Law>")
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != tt.wantPath {
					t.Errorf("path = %q, want %q", r.URL.Path, tt.wantPath)
				}
				if tt.wantAsof {
					if r.URL.Query().Get("asof") != tt.asof {
						t.Errorf("asof = %q, want %q", r.URL.Query().Get("asof"), tt.asof)
					}
				} else {
					if r.URL.RawQuery != "" {
						t.Errorf("unexpected query: %q", r.URL.RawQuery)
					}
				}
				w.Header().Set("Content-Type", "application/octet-stream")
				w.WriteHeader(http.StatusOK)
				w.Write(expected)
			}))
			defer server.Close()

			client := NewClient(server.URL)
			data, err := client.GetLawFile(tt.fileType, tt.idOrNum, tt.asof)
			if err != nil {
				t.Fatalf("GetLawFile() error: %v", err)
			}
			if !bytes.Equal(data, expected) {
				t.Errorf("data = %q, want %q", data, expected)
			}
		})
	}
}

func TestGetLawFileError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"code":    "INTERNAL_ERROR",
			"message": "サーバーエラー",
		})
	}))
	defer server.Close()

	client := NewClient(server.URL)
	_, err := client.GetLawFile("xml", "test-id", "")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// isNotFoundError checks if err is a *cerrors.NotFoundError via errors.As-like check.
func isNotFoundError(err error, target **cerrors.NotFoundError) bool {
	if e, ok := err.(*cerrors.NotFoundError); ok {
		*target = e
		return true
	}
	return false
}
