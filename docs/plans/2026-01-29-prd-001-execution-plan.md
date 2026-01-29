# PRD-001 Execution Plan — tooladapter (Strict TDD)

**Status:** Ready
**Date:** 2026-01-29

This is a strict TDD execution guide for PRD-001. Each task must follow:
1) Red (write failing test)
2) Red verification (run test)
3) Green (minimal implementation)
4) Green verification (run test)
5) Commit (one commit per task)

---

## Task 0 — Repo scaffolding

**Goal:** Establish Go module and baseline docs without behavior.

**Steps:**
- Create `go.mod`:
  - `module github.com/jonwraymond/tooladapter`
  - `go 1.24.4` (match existing stack).
  - Require `github.com/jonwraymond/toolmodel` (align with current version in `ai-tools-stack/go.mod`).
- Add `doc.go` with package description and scope.
- Add minimal `README.md` (what it is / not).

**Commands:**
- `go mod init github.com/jonwraymond/tooladapter`
- `go mod edit -go=1.24.4`
- `go get github.com/jonwraymond/toolmodel@<current>`

**Commit:** `chore(tooladapter): initialize module skeleton`

---

## Task 1 — CanonicalTool + JSONSchema

**Goal:** Implement canonical tool representation and JSON schema superset utilities.

**Test cases:**
- `CanonicalTool.ID()` behavior with and without namespace.
- `CanonicalTool.Validate()` requires name and input schema.
- `JSONSchema.DeepCopy()` does not alias slices/maps.
- `JSONSchema.ToMap()` omits zero fields and renders nested properties.

**Implementation details:**
- `canonical.go` with:
  - `CanonicalTool`, `JSONSchema`, `ToolHandler` (if required by tests).
  - `ID()` and `Validate()`.
  - `DeepCopy()` and `ToMap()` with full recursion for properties/items/defs.
- Guard against nil receivers.

**Commit:** `feat(tooladapter): add CanonicalTool and JSONSchema`

---

## Task 2 — Adapter interface + feature modeling

**Goal:** Provide adapter abstraction and feature loss modeling.

**Test cases:**
- `SchemaFeature.String()` for all enum values.
- `AllFeatures()` returns stable ordering.
- `ConversionError.Error()` contains adapter, direction, cause.
- `FeatureLossWarning.String()` includes feature + adapter.

**Implementation details:**
- `adapter.go`:
  - `Adapter` interface.
  - `SchemaFeature` enum + `String()`.
  - `AllFeatures()`.
  - `ConversionError` with `Unwrap()`.
  - `FeatureLossWarning`.

**Commit:** `feat(tooladapter): add Adapter interface and SchemaFeature`

---

## Task 3 — AdapterRegistry

**Goal:** Thread-safe adapter registry with conversion helper.

**Test cases:**
- Register, duplicate register, get, list, unregister.
- Convert: missing adapter errors, conversion errors wrapped.
- Convert: warnings emitted for unsupported features.

**Implementation details:**
- `registry.go`:
  - `AdapterRegistry` struct with RWMutex.
  - `Register/Get/List/Unregister`.
  - `Convert`:
    - `fromAdapter.ToCanonical` -> canonical tool.
    - `toAdapter.SupportsFeature` -> build warning list.
    - `toAdapter.FromCanonical` -> final tool.
    - Wrap errors with `ConversionError`.

**Commit:** `feat(tooladapter): add AdapterRegistry`

---

## Task 4 — MCP adapter

**Goal:** Bidirectional MCP <-> Canonical conversion using MCP SDK types.

**Pre-step:** Inspect MCP SDK types to avoid schema mismatch:
- Locate `mcp.Tool` and input schema types in `github.com/modelcontextprotocol/go-sdk/mcp`.
- Confirm schema fields and types (likely map-based input schema).

**Test cases:**
- `Name()` returns `mcp`.
- ToCanonical maps MCP `Name`, `Description`, input schema.
- FromCanonical emits MCP tool with input schema.
- Round-trip preserves name/description and core schema shape.

**Implementation details:**
- `adapters/mcp.go`:
  - `MCPAdapter` with `ToCanonical`/`FromCanonical`.
  - Schema conversion helpers for map <-> JSONSchema.
  - `SupportsFeature` returns true for all features.

**Commit:** `feat(tooladapter): add MCP adapter`

---

## Task 5 — OpenAI adapter

**Goal:** OpenAI function conversion with strict-mode behavior.

**Test cases:**
- Name, ToCanonical, FromCanonical.
- Strict mode is captured in output.
- Feature support: no `$ref`, `pattern` only in strict.

**Implementation details:**
- `adapters/openai.go`:
  - Define `OpenAIFunction` struct in this module.
  - Convert `parameters` (OpenAI JSON schema subset) to JSONSchema.
  - Strict mode sets `additionalProperties=false` (and pattern only if present).

**Commit:** `feat(tooladapter): add OpenAI adapter`

---

## Task 6 — Anthropic adapter

**Goal:** Anthropic conversion with `input_schema` mapping.

**Test cases:**
- Name, ToCanonical, FromCanonical.
- Feature support (no `$ref`).
- Round-trip preserves name/description and schema shape.

**Implementation details:**
- `adapters/anthropic.go`:
  - Define `AnthropicTool` struct (with `input_schema`).
  - Convert map schema to JSONSchema recursively.

**Commit:** `feat(tooladapter): add Anthropic adapter`

---

## Task 7 — Optional schema helpers (only if adapters need them)

**Goal:** Centralize validation or draft conversion.

**Test cases:**
- Minimal schema validates.
- draft-07 conversion (if implemented) is deterministic.

**Implementation details:**
- `schema/validate.go` and/or `schema/convert.go`.
- Keep dependency footprint minimal.

**Commit:** `feat(tooladapter): add schema helpers`

---

## Task 8 — Documentation & quality gates

- Create repo docs:
  - `docs/index.md`
  - `docs/design-notes.md`
  - `docs/user-journey.md`
- Add at least one **Mermaid** diagram showing conversion flow.
- Expand `README.md` with examples for MCP/OpenAI/Anthropic conversion.
- Add GoDoc to all exported types.
- Run:
  - `go test ./...`
  - `go vet ./...`
  - `golangci-lint run` (if configured)
  - `gosec ./...` (if configured)

**Commit:** `docs(tooladapter): add usage docs`

---

## Task 9 — Stack docs integration (ai-tools-stack)

**Goal:** Ensure tooladapter docs appear in the unified docs site with diagrams.

**Steps:**
- Add `ai-tools-stack/docs/components/tooladapter.md` (overview + usage + diagram embeds).
- Add `tooladapter` to `ai-tools-stack/mkdocs.yml`:
  - Components nav entry.
  - Library Docs multirepo import entry.
- Add D2 diagram:
  - `ai-tools-stack/docs/diagrams/component-tooladapter.d2`
  - Run `ai-tools-stack/scripts/render-d2.sh` to generate SVG.
- If stack overview diagrams need an update, add a note or update `docs/architecture/overview.md`.
- Optional preview:
  - `ai-tools-stack/scripts/prepare-mkdocs-multirepo.sh`
  - `mkdocs serve`

**Commit:** `docs(ai-tools-stack): add tooladapter docs and diagrams`

---

## Task 10 — Version matrix propagation

**Goal:** Align stack versioning with `ai-tools-stack` and propagate to repos.

**Steps:**
- Tag `tooladapter` with `vX.Y.Z`.
- Update `ai-tools-stack/go.mod` to include `tooladapter vX.Y.Z`.
- Run `ai-tools-stack/scripts/update-version-matrix.sh --apply`.
- Verify `VERSIONS.md` synced into each repo.

**Commit:** `chore(ai-tools-stack): add tooladapter to version matrix`

---

## Verification Checklist

- [ ] `go test ./...` passes
- [ ] `go vet ./...` passes
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
10) docs(ai-tools-stack): add tooladapter docs and diagrams
11) chore(ai-tools-stack): add tooladapter to version matrix
