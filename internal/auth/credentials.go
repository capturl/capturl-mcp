package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const (
	// This is a public API key for the Capturl MCP Server.
	capturlApiKey = "AIzaSyBAFxFnmuatYs1Ikx97HNbuIqiJp138I1o"
)

// Token represents an ID token.
type Token struct {
	IDToken   string    `json:"id_token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// TokenProvider is the interface for a token provider.
type TokenProvider interface {
	// GetToken returns an authentication token which can be used to authenticate requests to the Capturl API.
	GetToken() (*Token, error)
}

type tokenProvider struct {
	config *tokenProviderConfig

	mu    sync.Mutex
	token *Token
}

type tokenProviderConfig struct {
	Client       *http.Client
	RefreshToken string
	RefreshURL   string
}

type TokenProviderOption func(*tokenProviderConfig)

func WithClient(client *http.Client) TokenProviderOption {
	return func(c *tokenProviderConfig) {
		c.Client = client
	}
}

func WithRefreshURL(refreshURL string) TokenProviderOption {
	return func(c *tokenProviderConfig) {
		c.RefreshURL = refreshURL
	}
}

// NewTokenProvider creates a new token provider.
func NewTokenProvider(refreshToken string, opts ...TokenProviderOption) TokenProvider {
	config := &tokenProviderConfig{
		Client:       http.DefaultClient,
		RefreshToken: refreshToken,
		RefreshURL:   "https://securetoken.googleapis.com/v1/token",
	}
	for _, opt := range opts {
		opt(config)
	}
	return &tokenProvider{
		config: config,
	}
}

func (p *tokenProvider) GetToken() (*Token, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// If token is still valid (with 1 minute buffer), return it
	if p.token != nil && time.Until(p.token.ExpiresAt) > time.Minute {
		return p.token, nil
	}

	// Refresh the token
	token, err := p.mintNewToken()
	if err != nil {
		return nil, err
	}
	p.token = token
	return token, nil
}

func (p *tokenProvider) mintNewToken() (*Token, error) {
	resp, err := p.config.Client.PostForm(
		p.config.RefreshURL+"?key="+capturlApiKey,
		url.Values{
			"grant_type":    {"refresh_token"},
			"refresh_token": {p.config.RefreshToken},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("refresh request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("refresh failed (%d): %s", resp.StatusCode, body)
	}
	var body struct {
		IDToken   string `json:"id_token"`
		ExpiresIn string `json:"expires_in"` // seconds as string
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("failed to decode refresh response: %w", err)
	}

	var expiresIn int
	if _, err := fmt.Sscanf(body.ExpiresIn, "%d", &expiresIn); err != nil {
		return nil, fmt.Errorf("failed to parse expires_in: %w", err)
	}
	return &Token{
		IDToken:   body.IDToken,
		ExpiresAt: time.Now().Add(time.Duration(expiresIn) * time.Second),
	}, nil
}
