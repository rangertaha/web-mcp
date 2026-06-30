// SPDX-License-Identifier: MIT

// Package app assembles the fully-configured web-mcp server from configuration.
// It is shared by the command entry point (cmd/web) so the exact server the
// binary runs is the one under test.
package app

import (
	"log"
	"os"

	"github.com/rangertaha/web-mcp/internal/config"
	"github.com/rangertaha/web-mcp/internal/prompts"
	"github.com/rangertaha/web-mcp/internal/server"
	"github.com/rangertaha/web-mcp/internal/web"
	"github.com/rangertaha/web-mcp/internal/web/fetch"
)

// Assemble builds the fully-configured server (all enabled toolsets and
// prompts) and returns it with a cleanup function. version is reported to
// clients.
func Assemble(cfg *config.Config, version string) (*server.Server, func(), error) {
	clients, err := web.NewClients(cfg.UserAgent, cfg.MaxBytes)
	if err != nil {
		return nil, nil, err
	}

	srv := server.New("web-mcp", version, cfg.ReadOnly)

	for _, ts := range toolsets() {
		if cfg.ToolsetEnabled(ts.Name) {
			ts.Register(srv, clients)
		}
	}

	// Diagnostics go to stderr; stdout is reserved for the MCP protocol.
	log.SetOutput(os.Stderr)

	prompts.Register(srv)

	return srv, func() {}, nil
}

// toolsets returns every toolset registrar, in registration order. New service
// areas are added here.
func toolsets() []server.Toolset {
	return []server.Toolset{
		{Name: fetch.Name, Register: fetch.Register},
	}
}
