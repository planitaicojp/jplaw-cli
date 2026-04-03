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
				{LawInfo: model.LawInfo{LawID: "test-id", LawNum: "テスト法令番号"}},
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
					LawInfo:   model.LawInfo{LawID: "test-id"},
					Sentences: []model.Sentence{{Text: "個人情報の保護"}},
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	params := &KeywordParams{Keyword: "個人情報"}
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
			LawInfo:   model.LawInfo{LawID: "test-id"},
			Revisions: []model.RevisionInfo{{LawRevisionID: "rev-1"}},
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
