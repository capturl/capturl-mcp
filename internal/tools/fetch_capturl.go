package tools

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// This regex matches the path pattern: /o/{orgId}/c/{capturlId}
var capturlURLRegex = regexp.MustCompile(`^/o/([^/]+)/c/([^/]+)$`)

// FetchCapturlTool returns the tool definition for the fetch_capturl tool.
func FetchCapturlTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "fetch_capturl",
		Description: "Fetch the content of a capturl url (https://capturl.com/o/{orgId}/c/{capturlId}). Encrypted capturls are not supported.",
	}
}

// FetchCapturlToolInput is the input for the fetch_capturl tool.
type FetchCapturlToolInput struct {
	URL string `json:"url" jsonschema:"The URL of the capturl to fetch"`
}

// FetchCapturlToolHandler is the handler for the fetch_capturl tool.
func FetchCapturlToolHandler(tc ToolContext, opts ...FetchCapturlToolOption) mcp.ToolHandlerFor[FetchCapturlToolInput, any] {
	config := &FetchCapturlToolConfig{
		HTTPClient: http.DefaultClient,
		BaseURL:    "https://capturl.com",
	}
	for _, opt := range opts {
		opt(config)
	}
	return func(ctx context.Context, request *mcp.CallToolRequest, input FetchCapturlToolInput) (*mcp.CallToolResult, any, error) {
		httpClient := config.HTTPClient
		baseURL := config.BaseURL

		// Parse the URL to extract orgId and snipId
		parsedURL, err := url.Parse(input.URL)
		if err != nil {
			slog.Error("failed to parse input URL", "error", err)
			return nil, nil, errors.New("input URL was not formatted correctly. Please use the following format: https://capturl.com/o/{orgId}/c/{capturlId}")
		}
		matches := capturlURLRegex.FindStringSubmatch(parsedURL.Path)
		if matches == nil || len(matches) != 3 {
			slog.Error("input URL did not match expected pattern", "url", input.URL)
			return nil, nil, errors.New("input URL was not formatted correctly. Please use the following format: https://capturl.com/o/{orgId}/c/{capturlId}")
		}

		orgId := matches[1]
		snipId := matches[2]

		// Get a valid authentication token (auto-refreshes if needed)
		token, err := tc.TokenProvider.GetToken()
		if err != nil {
			slog.Error("failed to get authentication token", "error", err)
			return nil, nil, errors.New("unable to get an fresh id token. Did you correctly set the CAPTURL_AUTH_TOKEN environment variable?")
		}

		// Fetch the rendered capturl from the capturl web app.
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/o/%s/c/%s.jpeg", strings.TrimRight(baseURL, "/"), orgId, snipId), nil)
		if err != nil {
			slog.Error("failed to create http request", "error", err)
			return nil, nil, errors.New("server was unable to create the http request to fetch the capturl")
		}
		req.Header.Set("Authorization", "Bearer "+token.IDToken)
		httpResp, err := httpClient.Do(req)
		if err != nil {
			slog.Error("failed to fetch capturl", "error", err)
			return nil, nil, errors.New("server was unable to fetch the capturl")
		}
		defer httpResp.Body.Close()
		// Check the content type of the response to see if we were able to successfully fetch the capturl.
		if httpResp.Header.Get("Content-Type") == "application/json" {
			// This API returns a JSON response when an error occurs.
			body, err := io.ReadAll(httpResp.Body)
			if err != nil {
				slog.Error("failed to read capturl response body", "error", err)
				return nil, nil, errors.New("unable to read the capturl response body")
			}
			return nil, nil, fmt.Errorf("capturl API returned an error: %s", string(body))
		} else if strings.HasPrefix(httpResp.Header.Get("Content-Type"), "image/") {
			imageBytes, err := io.ReadAll(httpResp.Body)
			if err != nil {
				slog.Error("failed to read capturl response body", "error", err)
				return nil, nil, errors.New("unable to read the capturl image")
			}
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.ImageContent{
						Data:     imageBytes,
						MIMEType: httpResp.Header.Get("Content-Type"),
					},
				},
			}, nil, nil
		} else {
			slog.Error("unsupported content type", "content-type", httpResp.Header.Get("Content-Type"))
			return nil, nil, errors.New("capturl API returned a content type that is not supported")
		}
	}
}
