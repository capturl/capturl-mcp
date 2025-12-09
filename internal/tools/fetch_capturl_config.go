package tools

import "net/http"

// FetchCapturlToolConfig is the configuration for the fetch_capturl tool.
type FetchCapturlToolConfig struct {
	// The base URL to use when fetching capturls.
	BaseURL string
	// The HTTP client to use when fetching capturls.
	HTTPClient *http.Client
}

// FetchCapturlToolOption is an option for the fetch_capturl tool handler.
type FetchCapturlToolOption func(*FetchCapturlToolConfig)

// WithHTTPClient is an option for the fetch_capturl tool to set the HTTP client.
func WithHTTPClient(httpClient *http.Client) FetchCapturlToolOption {
	return func(config *FetchCapturlToolConfig) {
		config.HTTPClient = httpClient
	}
}

// WithBaseURL is an option for the fetch_capturl tool to set the base URL.
func WithBaseURL(baseURL string) FetchCapturlToolOption {
	return func(config *FetchCapturlToolConfig) {
		config.BaseURL = baseURL
	}
}
