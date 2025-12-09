package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetTokenCachesValidToken(t *testing.T) {
	var callCount int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if err := r.ParseForm(); err != nil {
			t.Fatalf("failed to parse form: %v", err)
		}
		if got := r.FormValue("refresh_token"); got != "refresh-token" {
			t.Fatalf("unexpected refresh_token: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id_token":"token-123","expires_in":"120"}`))
	}))
	defer server.Close()

	provider := NewTokenProvider("refresh-token", WithClient(server.Client()), WithRefreshURL(server.URL+"/token"))

	tok1, err := provider.GetToken()
	if err != nil {
		t.Fatalf("first token fetch failed: %v", err)
	}
	tok2, err := provider.GetToken()
	if err != nil {
		t.Fatalf("second token fetch failed: %v", err)
	}
	if callCount != 1 {
		t.Fatalf("expected token to be cached, got %d requests", callCount)
	}
	if tok1 != tok2 {
		t.Fatalf("expected cached token pointer to be reused")
	}
}

func TestGetTokenRefreshesWhenExpired(t *testing.T) {
	var callCount int
	tokens := []string{"first-token", "second-token"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() { callCount++ }()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// expires_in less than the 1 minute freshness buffer to force refresh
		_, _ = w.Write([]byte(`{"id_token":"` + tokens[callCount] + `","expires_in":"30"}`))
	}))
	defer server.Close()

	provider := NewTokenProvider("refresh-token", WithClient(server.Client()), WithRefreshURL(server.URL+"/token"))

	tok1, err := provider.GetToken()
	if err != nil {
		t.Fatalf("first token fetch failed: %v", err)
	}
	tok2, err := provider.GetToken()
	if err != nil {
		t.Fatalf("second token fetch failed: %v", err)
	}
	if callCount != 2 {
		t.Fatalf("expected token to refresh, got %d requests", callCount)
	}
	if tok1.IDToken == tok2.IDToken {
		t.Fatalf("expected refreshed token to differ")
	}
}

func TestMintNewTokenHandlesErrorStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad refresh", http.StatusBadRequest)
	}))
	defer server.Close()

	provider := NewTokenProvider("refresh-token", WithClient(server.Client()), WithRefreshURL(server.URL+"/token"))
	_, err := provider.GetToken()
	if err == nil || !strings.Contains(err.Error(), "refresh failed (400)") {
		t.Fatalf("expected refresh failure error, got %v", err)
	}
}

func TestMintNewTokenHandlesBadExpiresIn(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id_token":"abc","expires_in":"not-a-number"}`))
	}))
	defer server.Close()

	provider := NewTokenProvider("refresh-token", WithClient(server.Client()), WithRefreshURL(server.URL+"/token"))
	_, err := provider.GetToken()
	if err == nil || !strings.Contains(err.Error(), "failed to parse expires_in") {
		t.Fatalf("expected parse error, got %v", err)
	}
}
