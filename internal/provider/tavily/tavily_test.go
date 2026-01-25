package tavily

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"uni-search-hub/internal/dto"
	"uni-search-hub/internal/model"
)

func TestClient_Search(t *testing.T) {
	var gotAuth string
	var gotBody map[string]any

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/search" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		gotAuth = r.Header.Get("Authorization")
		defer r.Body.Close()
		_ = json.NewDecoder(r.Body).Decode(&gotBody)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
  "query": "hello",
  "answer": "world",
  "results": [
    {"title":"t1","url":"u1","content":"c1","score":0.9,"raw_content":"rc1"}
  ],
  "response_time": 0.12,
  "request_id": "rid"
}`))
	}))
	defer srv.Close()

	c := NewClient()
	c.BaseURL = srv.URL
	c.HTTPClient = srv.Client()

	ch := &model.Channel{ApiKey: "k"}
	req := &dto.UnifiedSearchRequest{
		Query:          "hello",
		IncludeContent: "full",
	}
	resp, err := c.Search(context.Background(), ch, req)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if gotAuth != "Bearer k" {
		t.Fatalf("unexpected auth header: %s", gotAuth)
	}
	if gotBody["query"] != "hello" {
		t.Fatalf("unexpected request query: %v", gotBody["query"])
	}
	// full -> include_raw_content=true
	if gotBody["include_raw_content"] != true {
		t.Fatalf("expected include_raw_content=true, got: %v", gotBody["include_raw_content"])
	}

	if resp.Answer != "world" || resp.Query != "hello" {
		t.Fatalf("unexpected response: %+v", resp)
	}
	if len(resp.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(resp.Results))
	}
	if resp.Results[0].Content != "rc1" {
		t.Fatalf("expected raw_content mapped to content, got %q", resp.Results[0].Content)
	}
}

