# Configuration

All configuration is read from the environment. No credentials are required.

| Variable         | Required | Description                                            |
| ---------------- | :------: | ------------------------------------------------------ |
| `WEB_USER_AGENT` |    no    | User-Agent for outbound requests (default `web-mcp`).  |
| `WEB_MAX_BYTES`  |    no    | Cap on returned body size in bytes (default 1048576).  |
| `WEB_TOOLSETS`   |    no    | Comma-separated toolset names to enable, or `all`.     |
| `WEB_READONLY`   |    no    | `true` to expose only read-only tools.                 |
