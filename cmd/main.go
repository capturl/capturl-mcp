package main

import (
	"context"
	"log"
	"os"

	"github.com/capturl/capturl-mcp/internal/auth"
	"github.com/capturl/capturl-mcp/internal/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	ctx := context.Background()

	authToken := os.Getenv("CAPTURL_AUTH_TOKEN")
	if authToken == "" {
		log.Fatal("CAPTURL_AUTH_TOKEN is not set. Please go to https://capturl.com/profile/settings to get a new token and set it in the MCP server's environment variable.")
	}

	// Create common tool context.
	toolContext := tools.ToolContext{
		TokenProvider: auth.NewTokenProvider(authToken),
	}
	server := mcp.NewServer(&mcp.Implementation{Name: "capturl-mcp", Title: "Capturl MCP Server"}, nil)
	mcp.AddTool(server, tools.FetchCapturlTool(), tools.FetchCapturlToolHandler(toolContext))
	if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
