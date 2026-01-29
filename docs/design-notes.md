# Design Notes

This document explains the schema design decisions, conversion semantics, and limitations of the tooladapter library.

## CanonicalTool Design

### Purpose

`CanonicalTool` is the protocol-agnostic intermediate representation. It stores a superset of information from all supported formats, enabling lossless conversion when the target format supports the features used.

### Fields

| Field | Type | Purpose |
|-------|------|---------|
| `Namespace` | `string` | Optional grouping (e.g., "github", "slack") |
| `Name` | `string` | **Required**. Tool identifier |
| `Version` | `string` | Semantic version |
| `Description` | `string` | Human-readable explanation |
| `Category` | `string` | Classification |
| `Tags` | `[]string` | Discovery keywords |
| `InputSchema` | `*JSONSchema` | **Required**. Parameter schema |
| `OutputSchema` | `*JSONSchema` | Optional. Return schema |
| `Timeout` | `time.Duration` | Execution timeout hint |
| `SourceFormat` | `string` | Original format (e.g., "mcp") |
| `SourceMeta` | `map[string]any` | Format-specific metadata for round-trip |
| `RequiredScopes` | `[]string` | Authorization scopes |

### ID Generation

The `ID()` method returns a stable identifier:

- With namespace: `"namespace:name"` (e.g., `"github:get_repo"`)
- Without namespace: `"name"` (e.g., `"get_weather"`)

### Validation

The `Validate()` method checks:

1. `Name` is not empty
2. `InputSchema` is not nil

This matches MCP requirements where tools must have a name and input schema.

---

## JSONSchema Design

### JSON Schema Version Support

The `JSONSchema` struct targets **JSON Schema 2020-12** while maintaining **draft-07 compatibility**. This aligns with:

- MCP SDK which uses `github.com/google/jsonschema-go` (2020-12)
- OpenAI which uses a subset of JSON Schema
- Anthropic which uses standard JSON Schema features

### Supported Keywords

| Category | Keywords |
|----------|----------|
| **Type** | `type` |
| **Validation** | `minimum`, `maximum`, `minLength`, `maxLength`, `pattern`, `format`, `enum`, `const` |
| **Object** | `properties`, `required`, `additionalProperties` |
| **Array** | `items` |
| **Composition** | `anyOf`, `oneOf`, `allOf`, `not` |
| **References** | `$ref`, `$defs` |
| **Metadata** | `description`, `default` |

### Pointer Types for Optional Fields

Numeric constraints use pointers to distinguish "not set" from "set to zero":

```go
type JSONSchema struct {
    Minimum   *float64  // nil = not specified, &0.0 = minimum is 0
    Maximum   *float64
    MinLength *int
    MaxLength *int
}
```

Similarly, `AdditionalProperties` uses `*bool`:

- `nil`: not specified (default JSON Schema behavior)
- `true`: explicitly allow additional properties
- `false`: explicitly forbid additional properties

### DeepCopy Semantics

`DeepCopy()` creates a complete independent copy with:

- No aliasing of slices or maps
- Recursive copying of nested `*JSONSchema` in Properties, Items, Defs, combinators
- Copied pointers for numeric constraints

This is essential for safe concurrent use and modification without side effects.

### ToMap Conversion

`ToMap()` returns a `map[string]any` suitable for JSON serialization with:

- **Zero-value omission**: Empty strings, nil pointers, empty slices are omitted
- **Recursive conversion**: Nested Properties and Defs are converted to nested maps
- **Correct key names**: Uses JSON Schema keywords (`$ref`, `$defs`, `additionalProperties`)

---

## Feature Support Matrix

| Feature | MCP | OpenAI | Anthropic | Notes |
|---------|:---:|:------:|:---------:|-------|
| `$ref` | Yes | **No** | **No** | Schema references |
| `$defs` | Yes | **No** | **No** | Schema definitions |
| `anyOf` | Yes | **No** | Yes | Any of listed schemas |
| `oneOf` | Yes | **No** | Yes | Exactly one of listed schemas |
| `allOf` | Yes | **No** | Yes | All of listed schemas |
| `not` | Yes | **No** | Yes | Schema negation |
| `pattern` | Yes | Yes* | Yes | Regex pattern |
| `format` | Yes | Yes | Yes | Semantic format |
| `additionalProperties` | Yes | Yes | Yes | Extra properties control |
| `minimum`/`maximum` | Yes | Yes | Yes | Numeric bounds |
| `minLength`/`maxLength` | Yes | Yes | Yes | String length bounds |
| `enum` | Yes | Yes | Yes | Value enumeration |
| `const` | Yes | Yes | Yes | Single value |
| `default` | Yes | Yes | Yes | Default value |

*OpenAI supports `pattern` in strict mode only.

---

## Conversion Semantics

### Feature Loss Warnings

When the target adapter doesn't support a feature used in the source schema, the library generates a `FeatureLossWarning`:

```go
type FeatureLossWarning struct {
    Feature     SchemaFeature  // e.g., FeatureRef
    FromAdapter string         // e.g., "mcp"
    ToAdapter   string         // e.g., "openai"
}
```

**Important**: Feature loss is a **warning**, not an error. The conversion proceeds, but the target schema will be missing those features.

Example: Converting MCP tool with `$ref` to OpenAI:

```go
mcpTool := mcp.Tool{
    Name: "example",
    InputSchema: map[string]any{
        "$ref": "#/$defs/Person",
        "$defs": map[string]any{
            "Person": map[string]any{
                "type": "object",
                "properties": map[string]any{
                    "name": map[string]any{"type": "string"},
                },
            },
        },
    },
}

result, _ := registry.Convert(mcpTool, "mcp", "openai")

for _, w := range result.Warnings {
    // "feature $ref lost converting from mcp to openai"
    // "feature $defs lost converting from mcp to openai"
    fmt.Println(w)
}
```

### Recursive Feature Detection

Feature loss detection is **recursive**. If a schema has nested properties, items, or definitions that use unsupported features, warnings are generated for each occurrence.

### Round-Trip Preservation

Format-specific metadata is stored in `SourceMeta` to improve round-trip conversions:

| Format | Preserved in SourceMeta |
|--------|------------------------|
| MCP | `title` field |
| OpenAI | `strict` mode flag |
| Anthropic | (none currently) |

Example: MCP Title preservation:

```go
// MCP -> Canonical: Title stored in SourceMeta
canonical, _ := mcpAdapter.ToCanonical(mcpTool)
// canonical.SourceMeta["title"] == "My Tool Title"

// Canonical -> MCP: Title restored from SourceMeta
restored, _ := mcpAdapter.FromCanonical(canonical)
// restored.(mcp.Tool).Title == "My Tool Title"
```

### Determinism Guarantees

All conversions are deterministic:

1. **Order-insensitive**: Map iteration order doesn't affect output
2. **No random elements**: No UUIDs, timestamps, or random values injected
3. **Stable output**: Same input always produces identical output

This enables:

- Reliable testing with exact output comparison
- Caching of conversion results
- Version control of generated tool definitions

---

## Adapter-Specific Behavior

### MCP Adapter

- **Full support**: All JSONSchema features are preserved
- **SDK types**: Uses `github.com/modelcontextprotocol/go-sdk/mcp.Tool`
- **Schema format**: `InputSchema` is `map[string]any` (standard JSON marshaling)
- **Title handling**: Stored in SourceMeta for round-trip

### OpenAI Adapter

- **Self-contained types**: `OpenAIFunction` struct defined in this module
- **Strict mode**: When `strict: true`:
  - Sets `additionalProperties: false` on output schema
  - Pattern validation is enabled
- **Limited features**: No `$ref`, `$defs`, or combinators
- **Field mapping**: `Parameters` (not `InputSchema`)

### Anthropic Adapter

- **Self-contained types**: `AnthropicTool` struct defined in this module
- **Field mapping**: Uses `input_schema` (not `InputSchema` or `Parameters`)
- **Combinator support**: Supports `anyOf`, `oneOf`, `allOf`, `not`
- **No references**: Does not support `$ref` or `$defs`

---

## Error Handling

### ConversionError

All conversion errors are wrapped in `ConversionError`:

```go
type ConversionError struct {
    Adapter   string  // "mcp", "openai", "anthropic"
    Direction string  // "to_canonical" or "from_canonical"
    Cause     error   // Underlying error
}

func (e *ConversionError) Error() string
func (e *ConversionError) Unwrap() error  // For errors.Is/As
```

This enables:

```go
result, err := registry.Convert(tool, "mcp", "openai")
if err != nil {
    var convErr *tooladapter.ConversionError
    if errors.As(err, &convErr) {
        log.Printf("Conversion failed in %s adapter (%s): %v",
            convErr.Adapter, convErr.Direction, convErr.Cause)
    }
}
```

### Type Validation

Adapters reject incorrect input types with descriptive errors:

```go
_, err := mcpAdapter.ToCanonical("not a tool")
// Error: "expected mcp.Tool or *mcp.Tool"

_, err := openaiAdapter.ToCanonical(mcpTool)
// Error: "expected OpenAIFunction or *OpenAIFunction"
```

---

## Limitations

1. **No schema validation**: This library converts schemas but does not validate data against schemas. Use a JSON Schema validator for that.

2. **No I/O**: Pure data transforms only. No network calls, file operations, or tool execution.

3. **No $ref resolution**: References are preserved as-is. The consuming system must resolve them.

4. **No schema inference**: Schemas must be explicitly provided. The library does not infer schemas from sample data.

5. **No streaming**: Conversion is synchronous. For large batches, callers should implement their own parallelization.

6. **Conservative conversion**: The library does not fabricate schema fields. If a feature isn't in the source, it won't appear in the output.
