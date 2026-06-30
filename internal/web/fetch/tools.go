// SPDX-License-Identifier: MIT

package fetch

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rangertaha/web-mcp/internal/server"
	"github.com/rangertaha/web-mcp/internal/web"
)

// Register adds the fetch toolset to the server.
func Register(s *server.Server, c *web.Clients) {
	s.NoteToolset(Name)
	svc := &service{c: c}

	server.Register(s, server.ToolDef{
		Name:        "fetch_url",
		Title:       "Fetch URL",
		Description: "Fetch the contents of an http(s) URL and return the response body (truncated) with status and content type.",
	}, svc.fetch)
}

// FetchInput identifies the URL to fetch.
type FetchInput struct {
	URL string `json:"url" jsonschema:"absolute http:// or https:// URL to fetch"`
}

func (s *service) fetch(ctx context.Context, _ *mcp.CallToolRequest, in FetchInput) (*mcp.CallToolResult, *Result, error) {
	out, err := s.Fetch(ctx, in.URL)
	return nil, out, err
}
