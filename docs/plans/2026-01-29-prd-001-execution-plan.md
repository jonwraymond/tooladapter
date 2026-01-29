# PRD-001 Execution Plan — tooladapter (Strict TDD)

**Status:** Ready
**Date:** 2026-01-29

This plan is a strict TDD execution guide for PRD-001. Each task must follow:
1) Red (write failing test)
2) Red verification (run test)
3) Green (minimal implementation)
4) Green verification (run test)
5) Commit (one commit per task)

---

## Task 0 — Repo scaffolding

**Goal:** Establish Go module and baseline docs without behavior.

- Create `go.mod` (Go 1.24+) with dependency on `toolmodel`.
- Add `doc.go` package documentation.
- Add minimal `README.md` with scope + usage note.

**Commit:** `chore(tooladapter): initialize module skeleton`

---

## Task 1 — CanonicalTool + JSONSchema

**Goal:** Implement canonical tool representation and JSON schema superset utilities.

**Tests:**
- `CanonicalTool.ID()` behavior with/without namespace.
- `CanonicalTool.Validate()` requires name + input schema.
- `JSONSchema.DeepCopy()` preserves original immutability.
- `JSONSchema.ToMap()` shape conversion.

**Implementation:**
- `canonical.go` with CanonicalTool, JSONSchema, DeepCopy, ToMap.
- Ensure all exported types have GoDoc comments.

**Commit:** `feat(tooladapter): add CanonicalTool and JSONSchema`

---

## Task 2 — Adapter interface + feature modeling

**Goal:** Provide adapter abstraction and feature loss modeling.

**Tests:**
- `SchemaFeature.String()` for all enum values.
- `ConversionError.Error()` formatting.
- `FeatureLossWarning.String()` formatting.

**Implementation:**
- `adapter.go` with Adapter interface, SchemaFeature enum, AllFeatures, ConversionError, FeatureLossWarning.

**Commit:** `feat(tooladapter): add Adapter interface and SchemaFeature`

---

## Task 3 — AdapterRegistry

**Goal:** Thread-safe adapter registry with conversion helper.

**Tests:**
- Register, duplicate register, get, list, unregister.
- Convert calls both adapters and returns warnings for unsupported features.

**Implementation:**
- `registry.go` with Register/Get/List/Unregister/Convert.
- Use RWMutex for concurrency safety.

**Commit:** `feat(tooladapter): add AdapterRegistry`

---

## Task 4 — MCP adapter

**Goal:** Bidirectional MCP <-> Canonical conversion.

**Tests:**
- `Name()` returns `mcp`.
- ToCanonical parses MCP tool with JSON schema.
- FromCanonical emits MCP tool with schema.
- Round-trip preserves Name/Description/schema basics.

**Implementation:**
- `adapters/mcp.go` with conversion helpers.
- Support all schema features.

**Commit:** `feat(tooladapter): add MCP adapter`

---

## Task 5 — OpenAI adapter

**Goal:** OpenAI function conversion with strict-mode behavior.

**Tests:**
- Name, ToCanonical, FromCanonical, strict mode flag.
- Feature support for $ref and pattern validation.

**Implementation:**
- `adapters/openai.go` with OpenAIFunction shape and conversion helpers.

**Commit:** `feat(tooladapter): add OpenAI adapter`

---

## Task 6 — Anthropic adapter

**Goal:** Anthropic conversion with `input_schema` mapping.

**Tests:**
- Name, ToCanonical, FromCanonical.
- Feature support (no $ref).
- Round-trip preserves key fields.

**Implementation:**
- `adapters/anthropic.go` with conversion helpers.

**Commit:** `feat(tooladapter): add Anthropic adapter`

---

## Task 7 — Optional schema helpers (if needed by adapters)

**Goal:** Centralize JSON schema validation or draft conversion.

**Tests:**
- Validate minimal schema.
- Convert draft-07 to 2020-12 where required (only if used).

**Implementation:**
- `schema/convert.go`, `schema/validate.go` as thin wrappers.

**Commit:** `feat(tooladapter): add schema helpers`

---

## Task 8 — Documentation & quality gates

- `README.md` usage examples.
- Run `go test ./...`.
- Run `golangci-lint run` if configured.

**Commit:** `docs(tooladapter): add usage docs`

---

## Verification Checklist

- [ ] `go test ./...` passes
- [ ] coverage > 80%
- [ ] no lints (if configured)
- [ ] README + doc.go updated

---

## Commit Order

1) chore(tooladapter): initialize module skeleton
2) feat(tooladapter): add CanonicalTool and JSONSchema
3) feat(tooladapter): add Adapter interface and SchemaFeature
4) feat(tooladapter): add AdapterRegistry
5) feat(tooladapter): add MCP adapter
6) feat(tooladapter): add OpenAI adapter
7) feat(tooladapter): add Anthropic adapter
8) feat(tooladapter): add schema helpers (optional)
9) docs(tooladapter): add usage docs
