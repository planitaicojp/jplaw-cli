package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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
		resp := model.LawDataResponse{
			LawInfo: model.LawInfo{LawID: "test-id"},
		}
		json.NewEncoder(w).Encode(resp)
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

func TestGetLawDataWithAttachedFiles(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := model.LawDataResponse{
			LawInfo: model.LawInfo{LawID: "test-id"},
			AttachedFilesInfo: &model.AttachedFilesInfo{
				AttachedFiles: []model.AttachedFile{
					{LawRevisionID: "rev-1", Src: "image001.png", Updated: "2024-01-01"},
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
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

func TestGetLawFile(t *testing.T) {
	expected := []byte("<Law>テスト法令XML</Law>")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/law_file/xml/405AC0000000088" {
			t.Errorf("path = %q, want /law_file/xml/405AC0000000088", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(expected)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	data, err := client.GetLawFile("xml", "405AC0000000088", "")
	if err != nil {
		t.Fatalf("GetLawFile() error: %v", err)
	}
	if string(data) != string(expected) {
		t.Errorf("data = %q, want %q", string(data), string(expected))
	}
}

func TestGetLawFileWithAsof(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("asof") != "2024-01-01" {
			t.Errorf("asof = %q, want 2024-01-01", r.URL.Query().Get("asof"))
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write([]byte("data"))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	data, err := client.GetLawFile("json", "test-id", "2024-01-01")
	if err != nil {
		t.Fatalf("GetLawFile() error: %v", err)
	}
	if len(data) == 0 {
		t.Error("GetLawFile() returned empty data")
	}
}
