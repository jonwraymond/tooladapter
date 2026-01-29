# PRD-001: tooladapter Library Implementation

> **For agents:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Create a protocol-agnostic tool abstraction library that enables bidirectional conversion between MCP, OpenAI, Anthropic, LangChain (optional), and OpenAPI (import only) tool formats.

**Architecture:** Introduce a canonical tool representation that stores the superset of schema information, enabling lossless conversion between formats. Protocol adapters implement a common interface for bidirectional conversion and feature capability reporting.

**Tech Stack:** Go 1.24+, JSON Schema, MCP Go SDK (for adapter), toolmodel dependency (for version alignment only)

**Priority:** P1

---

## Overview

`tooladapter` provides protocol-agnostic tool handling. It allows tools from any source to be normalized into a canonical representation, then re-exported to other protocol shapes. This is the foundation for `toolset` and later protocol-agnostic tooling.

**Reference:** `metatools-mcp/docs/proposals/protocol-agnostic-tools.md`

---

## Scope

### In scope
- Canonical tool representation (`CanonicalTool`) and full JSON Schema representation (`JSONSchema`).
- Adapter interface and schema feature capability enumeration.
- Adapter registry with conversion helper.
- MCP adapter.
- OpenAI adapter.
- Anthropic adapter.
- Schema deep copy + map conversion utilities.
- Unit tests for all exported behavior.

### Out of scope (future)
- LangChain adapter (optional stub only).
- OpenAPI adapter (import only).
- Runtime execution or tool invocation.
- Storage or persistence.

---

## Directory Structure

```
tooladapter/
├── canonical.go
├── canonical_test.go
├── adapter.go
├── adapter_test.go
├── registry.go
├── registry_test.go
├── adapters/
│   ├── mcp.go
│   ├── mcp_test.go
│   ├── openai.go
│   ├── openai_test.go
│   ├── anthropic.go
│   ├── anthropic_test.go
│   ├── langchain.go        # optional stub
│   └── openapi.go          # optional import-only stub
├── schema/
│   ├── convert.go
│   ├── convert_test.go
│   ├── validate.go
│   └── validate_test.go
├── doc.go
├── go.mod
└── go.sum
```

---

## Requirements

### R1 — Canonical representation
- Canonical tool ID is `namespace:name` when namespace is set, else `name`.
- Input schema is required; output schema is optional.
- Schema supports JSON Schema 2020-12 + draft-07 compatibility.

### R2 — Adapter interface
- `ToCanonical(raw any)` and `FromCanonical(tool *CanonicalTool)`.
- `SupportsFeature(feature SchemaFeature)`.
- Typed errors for conversion failure.

### R3 — Adapter registry
- Thread-safe registry with Register/Get/List/Unregister.
- Conversion helper with feature loss warnings.

### R4 — Protocol adapters
- MCP adapter supports all JSON Schema features.
- OpenAI adapter supports strict mode and reports limited features.
- Anthropic adapter uses `input_schema` and reports unsupported `$ref`.

### R5 — Tests
- TDD: failing tests first.
- Coverage target >80% for the module.

---

## Acceptance Criteria

- CanonicalTool + JSONSchema implemented with DeepCopy and ToMap.
- Adapter interface + SchemaFeature + error/warning types implemented.
- AdapterRegistry implemented with conversion helper.
- MCP/OpenAI/Anthropic adapters implemented and tested.
- All tests pass with >80% coverage.
- Documentation (`doc.go`, README) present and accurate.

---

## Dependencies

- `github.com/modelcontextprotocol/go-sdk/mcp` (for MCP adapter).
- `github.com/jonwraymond/toolmodel` (version alignment and potential reuse later).

---

## Notes

- Keep JSON Schema conversion conservative; do not invent unsupported fields.
- Avoid UTCP terminology. MCP terminology only.
- All exported APIs must have GoDoc comments.
