package middleware

import "testing"

func TestNormalizeTokenKey(t *testing.T) {
	t.Run("authorization_bearer", func(t *testing.T) {
		key, src := normalizeTokenKey("Bearer abcdef", "")
		if src != "authorization" {
			t.Fatalf("unexpected source: %s", src)
		}
		if key != "abcdef" {
			t.Fatalf("unexpected key: %q", key)
		}
	})

	t.Run("authorization_bearer_lower", func(t *testing.T) {
		key, _ := normalizeTokenKey("bearer abcdef", "")
		if key != "abcdef" {
			t.Fatalf("unexpected key: %q", key)
		}
	})

	t.Run("x_subscription_fallback", func(t *testing.T) {
		key, src := normalizeTokenKey("", "zzz")
		if src != "x-subscription-token" {
			t.Fatalf("unexpected source: %s", src)
		}
		if key != "zzz" {
			t.Fatalf("unexpected key: %q", key)
		}
	})

	t.Run("sk_prefix_compat", func(t *testing.T) {
		// 48-char key prefixed with "sk-" => total length 51
		raw := "sk-" + "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUV" // 48 chars after prefix
		if len(raw) != 51 {
			t.Fatalf("unexpected raw length: %d", len(raw))
		}
		key, _ := normalizeTokenKey("Bearer "+raw, "")
		if key != "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUV" {
			t.Fatalf("unexpected key: %q", key)
		}
	})

	t.Run("sk_prefix_not_stripped_when_short", func(t *testing.T) {
		key, _ := normalizeTokenKey("Bearer sk-abc", "")
		if key != "sk-abc" {
			t.Fatalf("unexpected key: %q", key)
		}
	})
}
