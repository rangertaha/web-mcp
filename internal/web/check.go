// SPDX-License-Identifier: MIT

package web

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// Check verifies outbound connectivity by fetching a stable test URL. It
// returns the HTTP status code observed.
func Check(ctx context.Context, c *Clients) (int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://example.com", nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("User-Agent", c.UserAgent)
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return 0, fmt.Errorf("fetching https://example.com: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1024))
	return resp.StatusCode, nil
}
