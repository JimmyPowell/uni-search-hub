package metaso

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_ProxySearch_Passthrough(t *testing.T) {
	rawReq := []byte(`{"q":"hello","scope":"webpage","page":3}`)

	var gotAuth string
	var gotCT string
	var gotBody string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotCT = r.Header.Get("Content-Type")
		defer r.Body.Close()
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"error":"rate_limited"}`))
	}))
	defer srv.Close()

	c := NewClient()
	c.BaseURL = srv.URL
	c.HTTPClient = srv.Client()

	status, respCT, respBody, err := c.ProxySearch(context.Background(), "k", rawReq, "application/json")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if gotAuth != "Bearer k" {
		t.Fatalf("unexpected auth: %s", gotAuth)
	}
	if gotCT != "application/json" {
		t.Fatalf("unexpected content-type: %s", gotCT)
	}
	if gotBody != string(rawReq) {
		t.Fatalf("unexpected body: %s", gotBody)
	}
	if status != http.StatusTooManyRequests {
		t.Fatalf("unexpected status: %d", status)
	}
	if respCT != "application/json" {
		t.Fatalf("unexpected resp content-type: %s", respCT)
	}
	if string(respBody) != `{"error":"rate_limited"}` {
		t.Fatalf("unexpected resp body: %s", string(respBody))
	}
}

