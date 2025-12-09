package tools

import (
	"github.com/capturl/capturl-mcp/internal/auth"
)

// ToolContext provides shared dependencies for all MCP tool handlers.
type ToolContext struct {
	// The token provider to use when making API calls to Capturl.
	TokenProvider auth.TokenProvider
}
