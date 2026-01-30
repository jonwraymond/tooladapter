# PRD-002 Execution Plan — tooladapter (TDD)

**Status:** Done
**Date:** 2026-01-30
**PRD:** `2026-01-30-prd-002-interface-contracts.md`


## TDD Workflow (required)
1. Red — write failing contract tests
2. Red verification — run tests
3. Green — minimal code/doc changes
4. Green verification — run tests
5. Commit — one commit per task


## Tasks
### Task 0 — Inventory + contract outline
- Confirm interface list and method signatures.
- Draft explicit contract bullets for each interface.
- Update docs/plans/README.md with this PRD + plan.
### Task 1 — Contract tests (Red/Green)
- Add `*_contract_test.go` with tests for each interface listed below.
- Use stub implementations where needed.
### Task 2 — GoDoc contracts
- Add/expand GoDoc on each interface with explicit contract clauses (thread-safety, errors, context, ownership).
- Update README/design-notes if user-facing.
### Task 3 — Verification
- Run `go test ./...`
- Run linters if configured (golangci-lint / gosec).


## Test Skeletons (contract_test.go)
### Adapter
```go
func TestAdapter_Contract(t *testing.T) {
    // Methods:
    // - Name() string
    // - ToCanonical(raw any) (*CanonicalTool, error)
    // - FromCanonical(tool *CanonicalTool) (any, error)
    // - SupportsFeature(feature SchemaFeature) bool
    // Contract assertions:
    // - Concurrency guarantees documented and enforced
    // - Error semantics (types/wrapping) validated
    // - Context cancellation respected (if applicable)
    // - Deterministic ordering where required
    // - Nil/zero input handling specified
}
```
