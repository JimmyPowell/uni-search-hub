package handler

import (
	"strings"
	"testing"
)

func TestStripJSONField(t *testing.T) {
	in := []byte(`{"query":"hello","api_key":"secret","max_results":2}`)
	out := stripJSONField(in, "api_key")
	if string(out) == string(in) {
		t.Fatalf("expected api_key stripped, got unchanged")
	}
	s := string(out)
	if s == "" {
		t.Fatalf("unexpected empty output")
	}
	if s == `{"query":"hello","api_key":"secret","max_results":2}` {
		t.Fatalf("api_key still present: %s", s)
	}
	if !(strings.Contains(s, `"query":"hello"`) && strings.Contains(s, `"max_results":2`) && !strings.Contains(s, `"api_key"`)) {
		t.Fatalf("unexpected output: %s", s)
	}
}
