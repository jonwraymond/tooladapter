# tooladapter

Protocol-agnostic tool format conversion library for Go.

## Overview

tooladapter enables bidirectional transformation between MCP, OpenAI, and Anthropic tool definitions through a canonical intermediate representation (`CanonicalTool`).

### Key Features

- **Pure data transforms** - No I/O, network, or runtime execution
- **Feature loss tracking** - Warnings when schema features aren't supported by target format
- **Thread-safe registry** - Concurrent adapter management
- **Deterministic conversions** - Same input always produces same output

## Installation

```bash
go get github.com/jonwraymond/tooladapter
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"

    "github.com/jonwraymond/tooladapter"
    "github.com/jonwraymond/tooladapter/adapters"
    "github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
    // Create registry and register adapters
    registry := tooladapter.NewRegistry()
    registry.Register(adapters.NewMCPAdapter())
    registry.Register(adapters.NewOpenAIAdapter())
    registry.Register(adapters.NewAnthropicAdapter())

    // Create an MCP tool
    mcpTool := mcp.Tool{
        Name:        "get_weather",
        Description: "Get weather for a location",
        InputSchema: map[string]any{
            "type": "object",
            "properties": map[string]any{
                "location": map[string]any{
                    "type":        "string",
                    "description": "City name",
                },
            },
            "required": []any{"location"},
        },
    }

    // Convert MCP tool to OpenAI format
    result, err := registry.Convert(mcpTool, "mcp", "openai")
    if err != nil {
        log.Fatal(err)
    }

    // Check for feature loss warnings
    for _, w := range result.Warnings {
        fmt.Printf("Warning: %s\n", w)
    }

    openaiFunc := result.Tool.(adapters.OpenAIFunction)
    fmt.Printf("OpenAI function: %s\n", openaiFunc.Name)
}
```

## Supported Formats

| Format | Adapter | Supports $ref | Supports Combinators |
|--------|---------|---------------|---------------------|
| MCP | `MCPAdapter` | Yes | Yes |
| OpenAI | `OpenAIAdapter` | No | No |
| Anthropic | `AnthropicAdapter` | No | Yes |

## Architecture

```
┌─────────────┐     ┌───────────────┐     ┌─────────────┐
│  MCP Tool   │────▶│ CanonicalTool │────▶│ OpenAI Fn   │
└─────────────┘     └───────────────┘     └─────────────┘
       ▲                    │                    │
       │                    ▼                    │
       │            ┌───────────────┐            │
       └────────────│   Registry    │◀───────────┘
                    └───────────────┘
```

## Next Steps

- [Design Notes](design-notes.md) - Schema decisions and limitations
- [User Journey](user-journey.md) - Complete conversion examples
