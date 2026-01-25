package zhipu_web_search

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"uni-search-hub/internal/dto"
	"uni-search-hub/internal/model"
)

func TestClient_Search_Mapping(t *testing.T) {
	var gotAuth string
	var gotBody map[string]any

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/api/paas/v4/web_search" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		gotAuth = r.Header.Get("Authorization")
		defer r.Body.Close()
		_ = json.NewDecoder(r.Body).Decode(&gotBody)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
  "id": "tid",
  "created": 1,
  "request_id": "rid",
  "search_intent": [{"query":"hello","intent":"SEARCH_ALWAYS","keywords":"hello"}],
  "search_result": [{"title":"t1","content":"c1","link":"u1","media":"m","icon":"i","refer":"ref_1","publish_date":"2025-01-01"}]
}`))
	}))
	defer srv.Close()

	c := NewClient()
	c.BaseURL = srv.URL + "/api"
	c.HTTPClient = srv.Client()

	ch := &model.Channel{ApiKey: "k"}
	req := &dto.UnifiedSearchRequest{
		Query:          "hello",
		SearchDepth:    "basic",
		MaxResults:     20,
		TimeRange:      "week",
		IncludeContent: "full",
		Site:           "example.com",
		AutoParameters: true,
	}

	resp, err := c.Search(context.Background(), ch, req)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if gotAuth != "Bearer k" {
		t.Fatalf("unexpected auth header: %s", gotAuth)
	}

	if gotBody["search_query"] != "hello" {
		t.Fatalf("unexpected search_query: %v", gotBody["search_query"])
	}
	if gotBody["search_engine"] != "search_std" {
		t.Fatalf("unexpected search_engine: %v", gotBody["search_engine"])
	}
	if gotBody["search_intent"] != true {
		t.Fatalf("expected search_intent=true, got: %v", gotBody["search_intent"])
	}
	if gotBody["count"] != float64(20) {
		t.Fatalf("unexpected count: %v", gotBody["count"])
	}
	if gotBody["search_domain_filter"] != "example.com" {
		t.Fatalf("unexpected search_domain_filter: %v", gotBody["search_domain_filter"])
	}
	if gotBody["search_recency_filter"] != "oneWeek" {
		t.Fatalf("unexpected search_recency_filter: %v", gotBody["search_recency_filter"])
	}
	if gotBody["content_size"] != "high" {
		t.Fatalf("unexpected content_size: %v", gotBody["content_size"])
	}

	if resp.Query != "hello" {
		t.Fatalf("unexpected resp query: %q", resp.Query)
	}
	if resp.RequestID != "rid" {
		t.Fatalf("unexpected resp request_id: %q", resp.RequestID)
	}
	if len(resp.Results) != 1 || resp.Results[0].URL != "u1" {
		t.Fatalf("unexpected results: %+v", resp.Results)
	}
	if resp.Results[0].Favicon != "i" {
		t.Fatalf("unexpected favicon: %q", resp.Results[0].Favicon)
	}
	if resp.Results[0].PublishedAt != "2025-01-01" {
		t.Fatalf("unexpected published_at: %q", resp.Results[0].PublishedAt)
	}
}
