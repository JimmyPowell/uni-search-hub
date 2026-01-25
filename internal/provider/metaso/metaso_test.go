package metaso

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"uni-search-hub/internal/dto"
	"uni-search-hub/internal/model"
)

func TestClient_Search(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/search" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer k" {
			t.Fatalf("unexpected auth: %s", r.Header.Get("Authorization"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
  "credits": 3,
  "searchParameters": {"q":"hello","scope":"webpage","page":2},
  "webpages":[{"title":"t","link":"u","score":"high","snippet":"s","position":1,"date":"2025-01-01"}],
  "total": 24
}`))
	}))
	defer srv.Close()

	c := NewClient()
	c.BaseURL = srv.URL
	c.HTTPClient = srv.Client()

	ch := &model.Channel{ApiKey: "k"}
	req := &dto.UnifiedSearchRequest{
		Query:          "hello",
		Offset:         1, // -> page=2
		IncludeContent: "summary",
	}
	resp, err := c.Search(context.Background(), ch, req)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if resp.Query != "hello" {
		t.Fatalf("unexpected query: %s", resp.Query)
	}
	if resp.Usage == nil || resp.Usage.Credits != 3 {
		t.Fatalf("unexpected usage: %+v", resp.Usage)
	}
	if len(resp.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(resp.Results))
	}
	if resp.Results[0].Score != 0.9 {
		t.Fatalf("expected score 0.9, got %v", resp.Results[0].Score)
	}
}

