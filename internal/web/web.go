// SPDX-License-Identifier: MIT

// Package web holds the shared HTTP client used by the web tool packages
// (fetch, …). Unlike the other servers in this family, web-mcp talks to
// arbitrary URLs rather than a single API host, so it uses net/http directly.
package web

import (
	"net/http"
	"time"
)

// defaultMaxBytes bounds how much of a response body is read and returned.
const defaultMaxBytes = 1 << 20 // 1 MiB

// defaultTimeout bounds a single fetch.
const defaultTimeout = 30 * time.Second

// Clients bundles the shared HTTP client and fetch policy.
type Clients struct {
	// HTTP performs the outbound requests.
	HTTP *http.Client
	// UserAgent is sent on every request.
	UserAgent string
	// MaxBytes caps how many bytes of a response body are returned.
	MaxBytes int
}

// NewClients builds the shared web client. An empty userAgent or non-positive
// maxBytes falls back to sensible defaults.
func NewClients(userAgent string, maxBytes int) (*Clients, error) {
	if userAgent == "" {
		userAgent = "web-mcp"
	}
	if maxBytes <= 0 {
		maxBytes = defaultMaxBytes
	}
	return &Clients{
		HTTP:      &http.Client{Timeout: defaultTimeout},
		UserAgent: userAgent,
		MaxBytes:  maxBytes,
	}, nil
}
