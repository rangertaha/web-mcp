// SPDX-License-Identifier: MIT

// Command web runs the web Model Context Protocol server (`web mcp`) and checks
// outbound connectivity (`web test`).
//
// Configuration is read from the environment (see package config). The `mcp`
// command communicates over stdio, the transport expected by MCP clients such
// as Claude Desktop/Code and Cursor.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/urfave/cli/v3"

	"github.com/rangertaha/web-mcp/internal"
	"github.com/rangertaha/web-mcp/internal/app"
	"github.com/rangertaha/web-mcp/internal/config"
	"github.com/rangertaha/web-mcp/internal/web"
)

func main() {
	cmd := &cli.Command{
		Name:    "web",
		Usage:   "Web fetch and search as an MCP server",
		Version: internal.Version(),
		// A bare `web` (no subcommand) runs the MCP server.
		Action: runMCP,
		Commands: []*cli.Command{
			mcpCommand(),
			testCommand(),
		},
		// Print errors ourselves so the MCP stdio stream is never touched.
		ExitErrHandler: func(context.Context, *cli.Command, error) {},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "web: %v\n", err)
		os.Exit(1)
	}
}

// mcpCommand runs the MCP server over stdio.
func mcpCommand() *cli.Command {
	return &cli.Command{
		Name:   "mcp",
		Usage:  "Run the MCP server over stdio (for Claude Desktop/Code, Cursor, ...)",
		Action: runMCP,
	}
}

// runMCP assembles and serves the MCP server over stdio.
func runMCP(ctx context.Context, _ *cli.Command) error {
	if err := config.LoadEnvFile(config.EnvFile); err != nil {
		log.Printf("web: reading %s: %v", config.EnvFile, err)
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("configuration error:\n%w", err)
	}

	ver := internal.Version()
	srv, cleanup, err := app.Assemble(cfg, ver)
	if err != nil {
		return err
	}
	defer cleanup()

	log.Printf("web-mcp %s starting: %d tools, %d prompts across toolsets %v (read-only=%v)",
		ver, srv.ToolCount(), srv.PromptCount(), srv.Toolsets(), cfg.ReadOnly)

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	return srv.Run(ctx, &mcp.StdioTransport{})
}

// testCommand verifies outbound connectivity.
func testCommand() *cli.Command {
	return &cli.Command{
		Name:  "test",
		Usage: "Test outbound HTTP connectivity",
		Action: func(ctx context.Context, _ *cli.Command) error {
			if err := config.LoadEnvFile(config.EnvFile); err != nil {
				log.Printf("web: reading %s: %v", config.EnvFile, err)
			}

			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("configuration error:\n%w", err)
			}

			clients, err := web.NewClients(cfg.UserAgent, cfg.MaxBytes)
			if err != nil {
				return err
			}

			ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
			defer cancel()

			status, err := web.Check(ctx, clients)
			if err != nil {
				return fmt.Errorf("outbound connectivity check failed: %w", err)
			}

			fmt.Printf("OK  fetched https://example.com (HTTP %d)\n", status)
			fmt.Printf("    read-only=%v\n", cfg.ReadOnly)
			return nil
		},
	}
}
