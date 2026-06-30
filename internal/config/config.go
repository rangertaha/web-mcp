// SPDX-License-Identifier: MIT

// Package config loads and validates runtime configuration for the web-mcp
// server from environment variables.
//
// All configuration is supplied via the environment so the server can run as a
// stdio subprocess launched by an MCP client (Claude Desktop/Code, Cursor, …),
// where command-line flags are awkward to pass.
package config

import (
	"os"
	"strconv"
	"strings"
)

// Environment variable names recognised by the server.
const (
	EnvUserAgent = "WEB_USER_AGENT" // User-Agent sent on outbound requests
	EnvMaxBytes  = "WEB_MAX_BYTES"  // cap on returned response body size
	EnvToolsets  = "WEB_TOOLSETS"   // comma-separated toolset names, or "all"
	EnvReadOnly  = "WEB_READONLY"   // "true" disables all write tools
)

// Config holds validated server configuration.
type Config struct {
	// UserAgent is sent on outbound HTTP requests. Empty uses a default.
	UserAgent string
	// MaxBytes caps how many bytes of a fetched body are returned. <=0 uses a
	// default.
	MaxBytes int
	// Toolsets is the set of enabled toolset names. A nil/empty set means "all".
	Toolsets []string
	// ReadOnly, when true, suppresses mutating tools at registration time.
	ReadOnly bool
}

// AllToolsets reports whether every toolset should be enabled.
func (c *Config) AllToolsets() bool {
	if len(c.Toolsets) == 0 {
		return true
	}
	for _, t := range c.Toolsets {
		if t == "all" {
			return true
		}
	}
	return false
}

// ToolsetEnabled reports whether the named toolset should be registered.
func (c *Config) ToolsetEnabled(name string) bool {
	if c.AllToolsets() {
		return true
	}
	for _, t := range c.Toolsets {
		if strings.EqualFold(t, name) {
			return true
		}
	}
	return false
}

// Load reads configuration from the process environment. web-mcp has no
// required configuration, so Load never fails today; it returns an error for
// signature parity with the other servers in this family.
func Load() (*Config, error) {
	maxBytes := 0
	if v := strings.TrimSpace(os.Getenv(EnvMaxBytes)); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			maxBytes = n
		}
	}
	return &Config{
		UserAgent: strings.TrimSpace(os.Getenv(EnvUserAgent)),
		MaxBytes:  maxBytes,
		Toolsets:  splitList(os.Getenv(EnvToolsets)),
		ReadOnly:  isTruthy(os.Getenv(EnvReadOnly)),
	}, nil
}

// splitList parses a comma-separated environment value into a trimmed,
// lower-cased slice, dropping empty entries.
func splitList(v string) []string {
	if strings.TrimSpace(v) == "" {
		return nil
	}
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.ToLower(strings.TrimSpace(p)); p != "" {
			out = append(out, p)
		}
	}
	return out
}

// isTruthy reports whether an environment value represents boolean true.
func isTruthy(v string) bool {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}
