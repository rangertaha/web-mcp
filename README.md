# web-mcp

[![CI](https://github.com/rangertaha/web-mcp/actions/workflows/ci.yml/badge.svg)](https://github.com/rangertaha/web-mcp/actions/workflows/ci.yml)
[![Status: under construction](https://img.shields.io/badge/status-under%20construction-orange)](#-under-construction)

<div align="center">

## 🚧 &nbsp; UNDER CONSTRUCTION &nbsp; 🚧

**This server is an early scaffold — a work in progress.**

It runs over stdio with **one read-only toolset** wired end-to-end.<br>
More toolsets are on the way (see the **TODO** list below).<br>
APIs, configuration, and tool names may still change.

</div>

---

A [Model Context Protocol](https://modelcontextprotocol.io) (MCP) server, written
in Go, exposing **web fetch** (and, later, search) as tools an LLM client
(Claude Desktop/Code, Cursor, and others) can call.

## Features

- **Typed tools with schemas**: every tool has an auto-generated JSON Schema for
  its input and output, inferred from Go structs.
- **Read-only switch**: `WEB_READONLY=true` hides every mutating tool.
- **Toolset filtering**: enable only the areas you need with `WEB_TOOLSETS`.
- **No credentials required** for fetching.

## Install

```sh
go install github.com/rangertaha/web-mcp/cmd/web@latest
```

Or build from source:

```sh
git clone https://github.com/rangertaha/web-mcp
cd web-mcp
make build        # produces ./bin/web
```

## CLI

```sh
web mcp      # run the MCP server over stdio (default when no subcommand)
web test     # verify outbound connectivity
```

## Configuration

| Variable         | Required | Description                                                |
| ---------------- | :------: | ---------------------------------------------------------- |
| `WEB_USER_AGENT` |    no    | User-Agent for outbound requests (default `web-mcp`).      |
| `WEB_MAX_BYTES`  |    no    | Cap on returned body size in bytes (default `1048576`).    |
| `WEB_TOOLSETS`   |    no    | Comma-separated toolset names to enable, or `all`.         |
| `WEB_READONLY`   |    no    | `true` to expose only read-only tools.                     |

## Toolsets

| Toolset | Covers                                              |
| ------- | --------------------------------------------------- |
| `fetch` | fetch the contents of an http(s) URL (`fetch_url`)  |

### TODO toolsets

- `search` — web search via a provider API (Brave/Tavily/SerpAPI); needs a key.
- `extract` — HTML-to-text / readability extraction and metadata.

## License

MIT — see [LICENSE](LICENSE).
