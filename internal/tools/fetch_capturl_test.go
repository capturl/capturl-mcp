package tools

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/capturl/capturl-mcp/internal/auth"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type stubTokenProvider struct {
	token *auth.Token
	err   error
	calls int
}

func (s *stubTokenProvider) GetToken() (*auth.Token, error) {
	s.calls++
	if s.err != nil {
		return nil, s.err
	}
	return s.token, nil
}

func TestFetchCapturlInputValidation(t *testing.T) {
	tc := ToolContext{TokenProvider: &stubTokenProvider{}}
	handler := FetchCapturlToolHandler(tc)

	_, _, err := handler(context.Background(), nil, FetchCapturlToolInput{URL: "not-a-url"})
	if err == nil || !strings.Contains(err.Error(), "input URL was not formatted correctly") {
		t.Fatalf("expected URL format error, got %v", err)
	}
}

func TestFetchCapturlJSONErrorResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"not found"}`))
	}))
	defer server.Close()

	tokenProvider := &stubTokenProvider{token: &auth.Token{IDToken: "id-token"}}
	tc := ToolContext{TokenProvider: tokenProvider}
	handler := FetchCapturlToolHandler(tc, WithHTTPClient(server.Client()), WithBaseURL(server.URL))

	_, _, err := handler(context.Background(), nil, FetchCapturlToolInput{URL: server.URL + "/o/org/c/id"})
	if err == nil || !strings.Contains(err.Error(), "capturl API returned an error") {
		t.Fatalf("expected API error, got %v", err)
	}
	if tokenProvider.calls != 1 {
		t.Fatalf("expected token provider to be called once, got %d", tokenProvider.calls)
	}
}

func TestFetchCapturlReturnsImageContent(t *testing.T) {
	expectedImage := []byte{0x01, 0x02, 0x03}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Fatalf("unexpected Authorization header: %q", got)
		}
		w.Header().Set("Content-Type", "image/png")
		_, _ = w.Write(expectedImage)
	}))
	defer server.Close()

	tokenProvider := &stubTokenProvider{token: &auth.Token{IDToken: "test-token"}}
	tc := ToolContext{TokenProvider: tokenProvider}
	handler := FetchCapturlToolHandler(tc, WithHTTPClient(server.Client()), WithBaseURL(server.URL))

	result, _, err := handler(context.Background(), nil, FetchCapturlToolInput{URL: server.URL + "/o/org/c/id"})
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if result == nil || len(result.Content) != 1 {
		t.Fatalf("expected a single content item, got %#v", result)
	}
	imageContent, ok := result.Content[0].(*mcp.ImageContent)
	if !ok {
		t.Fatalf("expected ImageContent, got %T", result.Content[0])
	}
	if !bytes.Equal(imageContent.Data, expectedImage) || imageContent.MIMEType != "image/png" {
		t.Fatalf("unexpected image content: %+v", imageContent)
	}
}

func TestFetchCapturlUnsupportedContentType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("unsupported"))
	}))
	defer server.Close()

	tokenProvider := &stubTokenProvider{token: &auth.Token{IDToken: "id-token"}}
	tc := ToolContext{TokenProvider: tokenProvider}
	handler := FetchCapturlToolHandler(tc, WithHTTPClient(server.Client()), WithBaseURL(server.URL))

	_, _, err := handler(context.Background(), nil, FetchCapturlToolInput{URL: server.URL + "/o/org/c/id"})
	if err == nil || !strings.Contains(err.Error(), "not supported") {
		t.Fatalf("expected unsupported content type error, got %v", err)
	}
}
