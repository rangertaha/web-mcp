// SPDX-License-Identifier: MIT

// Package prompts registers MCP prompts: user-invoked, parameterized templates
// that clients surface as slash commands. Each prompt encodes a multi-step
// workflow by guiding the model to call the right tools in order.
package prompts

import (
	"fmt"

	"github.com/rangertaha/web-mcp/internal/server"
)

// Register adds the built-in workflow prompts to the server.
func Register(s *server.Server) {
	s.AddPrompt(
		"summarize_url",
		"Fetch a web page and summarize its main points.",
		[]server.PromptArg{
			{Name: "url", Description: "absolute http(s) URL", Required: true},
		},
		func(a map[string]string) string {
			return fmt.Sprintf(`Summarize the page at %s.

Steps:
1. Call fetch_url (url="%s") to retrieve the content.
2. If the body is HTML, focus on the main article text and ignore navigation/boilerplate.
3. Give a 3-5 sentence summary of the main points, and note if the content was truncated.`,
				a["url"], a["url"])
		},
	)
}
