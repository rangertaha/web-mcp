// SPDX-License-Identifier: MIT

// Package fetch exposes a single tool that retrieves the contents of a URL.
package fetch

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/rangertaha/web-mcp/internal/web"
)

// Name is the toolset name used for enable/disable filtering.
const Name = "fetch"

// service wraps the shared web client.
type service struct {
	c *web.Clients
}

// Result is the outcome of fetching a URL, trimmed to the fields useful to an
// LLM. Body is truncated to the client's MaxBytes.
type Result struct {
	URL         string `json:"url"`
	StatusCode  int    `json:"statusCode"`
	ContentType string `json:"contentType,omitempty"`
	Bytes       int    `json:"bytes"`
	Truncated   bool   `json:"truncated"`
	Body        string `json:"body"`
}

// Fetch retrieves the URL over HTTP(S) and returns its (possibly truncated)
// body along with response metadata.
func (s *service) Fetch(ctx context.Context, rawURL string) (*Result, error) {
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		return nil, fmt.Errorf("url must start with http:// or https://")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", s.c.UserAgent)
	req.Header.Set("Accept", "*/*")

	resp, err := s.c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching %s: %w", rawURL, err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Read one extra byte past the cap to detect truncation.
	limited := io.LimitReader(resp.Body, int64(s.c.MaxBytes)+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", rawURL, err)
	}
	truncated := len(data) > s.c.MaxBytes
	if truncated {
		data = data[:s.c.MaxBytes]
	}

	return &Result{
		URL:         rawURL,
		StatusCode:  resp.StatusCode,
		ContentType: resp.Header.Get("Content-Type"),
		Bytes:       len(data),
		Truncated:   truncated,
		Body:        string(data),
	}, nil
}
