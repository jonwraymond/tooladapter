# tooladapter

Protocol-agnostic tool format conversion library for Go.

## Overview

tooladapter enables bidirectional transformation between MCP, OpenAI, and Anthropic tool definitions through a canonical intermediate representation (`CanonicalTool`).

Key features:
- Pure data transforms (no I/O, network, or runtime execution)
- Feature loss tracking with warnings
- Thread-safe adapter registry
- Deterministic conversions

## Installation

```bash
go get github.com/jonwraymond/tooladapter
```

## Quick Start

```go
import (
    "github.com/jonwraymond/tooladapter"
    "github.com/jonwraymond/tooladapter/adapters"
)

// Create registry and register adapters
registry := tooladapter.NewRegistry()
registry.Register(adapters.NewMCPAdapter())
registry.Register(adapters.NewOpenAIAdapter())
registry.Register(adapters.NewAnthropicAdapter())

// Convert MCP tool to OpenAI format
result, err := registry.Convert(mcpTool, "mcp", "openai")
if err != nil {
    log.Fatal(err)
}

// Check for feature loss warnings
for _, w := range result.Warnings {
    log.Printf("Warning: %s", w)
}

openaiFunc := result.Tool.(adapters.OpenAIFunction)
```

## Supported Formats

| Format | Adapter | Notes |
|--------|---------|-------|
| MCP | `MCPAdapter` | Full feature support |
| OpenAI | `OpenAIAdapter` | Strict mode enforces `additionalProperties=false` |
| Anthropic | `AnthropicAdapter` | No `$ref` support |

## License

See LICENSE file.
