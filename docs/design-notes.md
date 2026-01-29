# Design Notes

## Schema Representation

### JSONSchema Type

The `JSONSchema` struct is a superset of JSON Schema features used across MCP, OpenAI, and Anthropic formats. It includes:

- Basic types and constraints (type, minimum, maximum, etc.)
- Object properties and required fields
- Array items
- Combinators (anyOf, oneOf, allOf, not)
- References ($ref, $defs)

### Design Decisions

1. **Pointer types for optional fields**: Numeric constraints use pointers (`*float64`, `*int`) to distinguish between "not set" and "set to zero".

2. **AdditionalProperties as *bool**: Allows distinguishing between:
   - `nil`: not specified (default behavior)
   - `true`: explicitly allow additional properties
   - `false`: explicitly forbid additional properties

3. **ToMap() omits zero values**: The map representation only includes non-zero fields to produce clean, minimal JSON output.

## Feature Support Matrix

| Feature | MCP | OpenAI | Anthropic |
|---------|-----|--------|-----------|
| `$ref` | ✅ | ❌ | ❌ |
| `$defs` | ✅ | ❌ | ❌ |
| `anyOf` | ✅ | ❌ | ✅ |
| `oneOf` | ✅ | ❌ | ✅ |
| `allOf` | ✅ | ❌ | ✅ |
| `not` | ✅ | ❌ | ✅ |
| `pattern` | ✅ | ✅* | ✅ |
| `format` | ✅ | ✅ | ✅ |
| `additionalProperties` | ✅ | ✅ | ✅ |

*OpenAI supports pattern in strict mode

## Conversion Behavior

### Feature Loss Warnings

When converting between formats, features not supported by the target adapter generate `FeatureLossWarning`. This is a warning, not an error - the conversion proceeds but with reduced fidelity.

Example: Converting an MCP tool with `$ref` to OpenAI format:
```go
result, _ := registry.Convert(mcpTool, "mcp", "openai")
// result.Warnings contains FeatureLossWarning for FeatureRef
```

### Round-Trip Preservation

The library preserves format-specific metadata in `SourceMeta` to enable better round-trip conversions:

- **MCP**: Title stored in SourceMeta
- **OpenAI**: Strict mode stored in SourceMeta

### Determinism

All conversions are deterministic - the same input always produces the same output. This is important for:

- Caching and deduplication
- Testing and verification
- Version control of generated definitions

## Limitations

1. **No schema validation**: This library converts schemas but doesn't validate that input data matches schemas. Use a JSON Schema validator for that.

2. **No I/O**: This is a pure data-transform library. No network calls, no file operations.

3. **No runtime execution**: Tool definitions are converted, but tools are not executed.

4. **$ref resolution**: References are preserved as-is, not resolved. The target system must handle reference resolution.
