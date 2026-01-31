# Migration Guide: tooladapter to toolfoundation/adapter

This guide helps you migrate from `github.com/jonwraymond/tooladapter` to `github.com/jonwraymond/toolfoundation/adapter`.

## Import Path Changes

| Old Import | New Import |
|------------|------------|
| `github.com/jonwraymond/tooladapter` | `github.com/jonwraymond/toolfoundation/adapter` |
| `github.com/jonwraymond/tooladapter/adapters` | `github.com/jonwraymond/toolfoundation/adapter/adapters` |

## Step-by-Step Migration

### 1. Update go.mod

Remove the old dependency and add the new one:

```bash
go get github.com/jonwraymond/toolfoundation@latest
go mod tidy
```

### 2. Update Imports

Find and replace all imports in your codebase:

```bash
# Using sed (macOS/Linux)
find . -name "*.go" -exec sed -i '' \
  's|github.com/jonwraymond/tooladapter/adapters|github.com/jonwraymond/toolfoundation/adapter/adapters|g' {} +

find . -name "*.go" -exec sed -i '' \
  's|github.com/jonwraymond/tooladapter|github.com/jonwraymond/toolfoundation/adapter|g' {} +
```

Or manually update your imports:

```go
// Before
import (
    "github.com/jonwraymond/tooladapter"
    "github.com/jonwraymond/tooladapter/adapters"
)

// After
import (
    "github.com/jonwraymond/toolfoundation/adapter"
    "github.com/jonwraymond/toolfoundation/adapter/adapters"
)
```

### 3. Verify Your Code

```bash
go build ./...
go test ./...
```

## API Compatibility

The API remains fully compatible. The following types and functions work identically:

- `adapter.NewRegistry()` (was `tooladapter.NewRegistry()`)
- `adapter.CanonicalTool` (was `tooladapter.CanonicalTool`)
- `adapter.ConversionResult` (was `tooladapter.ConversionResult`)
- `adapters.NewMCPAdapter()`
- `adapters.NewOpenAIAdapter()`
- `adapters.NewAnthropicAdapter()`

## Example Migration

### Before

```go
package main

import (
    "log"

    "github.com/jonwraymond/tooladapter"
    "github.com/jonwraymond/tooladapter/adapters"
)

func main() {
    registry := tooladapter.NewRegistry()
    registry.Register(adapters.NewMCPAdapter())
    registry.Register(adapters.NewOpenAIAdapter())

    result, err := registry.Convert(mcpTool, "mcp", "openai")
    if err != nil {
        log.Fatal(err)
    }

    for _, w := range result.Warnings {
        log.Printf("Warning: %s", w)
    }
}
```

### After

```go
package main

import (
    "log"

    "github.com/jonwraymond/toolfoundation/adapter"
    "github.com/jonwraymond/toolfoundation/adapter/adapters"
)

func main() {
    registry := adapter.NewRegistry()
    registry.Register(adapters.NewMCPAdapter())
    registry.Register(adapters.NewOpenAIAdapter())

    result, err := registry.Convert(mcpTool, "mcp", "openai")
    if err != nil {
        log.Fatal(err)
    }

    for _, w := range result.Warnings {
        log.Printf("Warning: %s", w)
    }
}
```

## Getting Help

If you encounter issues during migration, please open an issue in the [toolfoundation repository](https://github.com/jonwraymond/toolfoundation/issues).
