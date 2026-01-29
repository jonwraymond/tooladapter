# tooladapter

Protocol-agnostic tool format conversion library for Go.

## Overview

tooladapter enables bidirectional transformation between MCP, OpenAI, and Anthropic tool definitions through a canonical intermediate representation (`CanonicalTool`). It is a pure data-transform library with no I/O, network access, or runtime execution.

## Design Goals

1. **Pure transforms**: All conversions are stateless data transformations. No network calls, no file I/O, no tool execution.

2. **Determinism**: The same input always produces the same output, enabling caching, deduplication, and reliable testing.

3. **Loss visibility**: When converting between formats, feature loss is tracked and reported as warnings, not hidden or silently dropped.

4. **Minimal dependencies**: Only depends on the MCP Go SDK for the MCP adapter. OpenAI and Anthropic adapters use self-contained types with no external SDK dependencies.

5. **Go idioms**: Errors are wrapped with context and support `errors.Unwrap()`. The registry is thread-safe. All exported types have GoDoc comments.

## Position in the Stack

tooladapter sits alongside (not inside) other ApertureStack libraries:

```
toolmodel --> tooladapter --> toolset --> metatools-mcp
    |              |
    v              v
  MCP-aligned    Protocol-agnostic
  definitions    conversion layer
```

**Dependency order (DAG):**

- `toolmodel` provides MCP-aligned tool definitions
- `tooladapter` provides protocol conversion (depends on toolmodel for alignment)
- `toolset` composes tools and can export to multiple formats via tooladapter
- `metatools-mcp` can optionally wire toolset to expose non-MCP formats

tooladapter does **not** change the MCP surface. It provides a normalization layer that other components can use to support multiple LLM providers.

## Key Features

- **Canonical representation**: `CanonicalTool` stores the superset of schema information from all supported formats
- **JSONSchema superset**: `JSONSchema` type supports JSON Schema 2020-12 features including `$ref`, `$defs`, combinators (`anyOf`, `oneOf`, `allOf`, `not`), and draft-07 compatibility
- **Feature loss tracking**: `FeatureLossWarning` indicates when schema features cannot be preserved in the target format
- **Thread-safe registry**: `AdapterRegistry` safely manages adapters across goroutines
- **Round-trip preservation**: Format-specific metadata stored in `SourceMeta` enables better round-trip conversions

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

| Format | Adapter | Supports $ref | Supports Combinators | Notes |
|--------|---------|:-------------:|:--------------------:|-------|
| MCP | `MCPAdapter` | Yes | Yes | Full JSON Schema support |
| OpenAI | `OpenAIAdapter` | No | No | Strict mode enforces `additionalProperties=false` |
| Anthropic | `AnthropicAdapter` | No | Yes | Uses `input_schema` field name |

## Architecture

```
+-----------+     +---------------+     +-----------+
| MCP Tool  |---->| CanonicalTool |---->| OpenAI Fn |
+-----------+     +---------------+     +-----------+
      ^                  |                   |
      |                  v                   |
      |          +---------------+           |
      +----------|   Registry    |<----------+
                 +---------------+
```

The registry manages adapters and provides the `Convert()` method that:

1. Calls `ToCanonical()` on the source adapter
2. Checks `SupportsFeature()` on the target adapter to generate warnings
3. Calls `FromCanonical()` on the target adapter
4. Returns the converted tool with any feature loss warnings

## Versioning

tooladapter follows semantic versioning. The authoritative version matrix is maintained in `ai-tools-stack/go.mod` and propagated to `VERSIONS.md` in each repository.

Current version: See [VERSIONS.md](../VERSIONS.md)

## Next Steps

- [Design Notes](design-notes.md) - Schema decisions, conversion semantics, and limitations
- [User Journey](user-journey.md) - Complete conversion examples with Mermaid diagrams
