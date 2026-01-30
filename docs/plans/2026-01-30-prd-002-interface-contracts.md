# PRD-002 Interface Contracts — tooladapter

**Status:** Done
**Date:** 2026-01-30


## Overview
Define explicit interface contracts (GoDoc + documented semantics) for all interfaces in this repo. Contracts must state concurrency guarantees, error semantics, ownership of inputs/outputs, and context handling.


## Goals
- Every interface has explicit GoDoc describing behavioral contract.
- Contract behavior is codified in tests (contract tests).
- Docs/README updated where behavior is user-facing.


## Non-Goals
- No API shape changes unless required to satisfy the contract tests.
- No new features beyond contract clarity and tests.


## Interface Inventory
| Interface | File | Methods |
| --- | --- | --- |
| `Adapter` | `tooladapter/adapter.go:97` | Name() string<br/>ToCanonical(raw any) (*CanonicalTool, error)<br/>FromCanonical(tool *CanonicalTool) (any, error)<br/>SupportsFeature(feature SchemaFeature) bool |

## Contract Template (apply per interface)
- **Thread-safety:** explicitly state if safe for concurrent use.
- **Context:** cancellation/deadline handling (if context is a parameter).
- **Errors:** classification, retryability, and wrapping expectations.
- **Ownership:** who owns/allocates inputs/outputs; mutation expectations.
- **Determinism/order:** ordering guarantees for returned slices/maps/streams.
- **Nil/zero handling:** behavior for nil inputs or empty values.


## Acceptance Criteria
- All interfaces have GoDoc with explicit behavioral contract.
- Contract tests exist and pass.
- No interface contract contradictions across repos.
