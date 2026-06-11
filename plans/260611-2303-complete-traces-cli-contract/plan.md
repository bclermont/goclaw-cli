---
title: "Complete Traces CLI Contract"
description: "Bring goclaw-cli traces commands back in sync with digitopvn/goclaw dev: real list/get/follow/export contracts plus run timeline support."
status: completed
priority: P1
effort: 2d
branch: "detached-head"
tags: [bugfix, traces, cli, api, tdd]
blockedBy: []
blocks: [260528-1357-fix-trace-details-by-id]
created: "2026-06-11T16:03:22.437Z"
createdBy: "ck:plan"
source: skill
---

# Complete Traces CLI Contract

## Overview

`digitopvn/goclaw` `dev` now exposes a trace API surface wider and slightly different from what `goclaw-cli` currently assumes. The CLI has `traces list/get/export/follow`, but static contract scout found:

- `GET /v1/traces` returns `{traces,total,limit,offset}`, while CLI still parses a top-level array.
- `GET /v1/traces/{traceID}` returns `{trace,spans}`, while CLI renders fields as if they are top-level.
- `GET /v1/runs/{runID}/timeline` exists and is documented, but CLI has no command.
- `traces export` does not reuse trace ID validation/path escaping.
- `traces list` exposes stale flags (`--since`, `--root-only`) and misses server-supported filters (`--session-key`, `--user`, `--channel`, `--offset`).

This plan is TDD-first: each implementation phase begins by replacing stale mocks with server-shaped `httptest` fixtures, then lands the smallest CLI changes needed.

## Scope

In scope:
- Align `traces list`, `traces get`, `traces export`, and `traces follow` with server `dev` trace contracts.
- Add run timeline read command for `GET /v1/runs/{runID}/timeline`.
- Update tests, README/help text, changelog, and docs status.

Out of scope:
- `POST /v1/traces/{id}/replay` because server `dev` does not expose it.
- New watch loops, SSE, or WebSocket trace streaming.
- Server-side changes in `/Volumes/GOON/www/digitop/goclaw`.

## Contract Sources

- Server handlers: `/Volumes/GOON/www/digitop/goclaw/internal/http/traces.go`
- Server docs: `/Volumes/GOON/www/digitop/goclaw/docs/18-http-api.md`
- Server response structs: `/Volumes/GOON/www/digitop/goclaw/internal/store/tracing_store.go`, `/Volumes/GOON/www/digitop/goclaw/internal/store/run_timeline_store.go`
- Existing CLI: `cmd/traces.go`, `cmd/traces_follow.go`, `cmd/traces_get_test.go`, `cmd/traces_follow_test.go`
- Existing output/client patterns: `cmd/helpers.go`, `internal/client/http.go`, `internal/output/tree.go`

## Decisions

- Keep trace work in existing `cmd/traces.go` unless file growth forces one small `cmd/traces_timeline.go` split for the new timeline command.
- Do not add a new generic response framework. Use tiny local helpers only where repeated across trace list/get/follow.
- JSON/YAML mode should preserve server payload shape; table mode may flatten for human scan.
- Reuse `validateTraceID` for `get` and `export`; add a separate `validateRunID` only if needed for path safety.
- Treat existing plan `260528-1357-fix-trace-details-by-id` as superseded for implementation details; its security findings remain useful.

## Phases

| Phase | Name | Status |
|-------|------|--------|
| 1 | [Contract Lock and Fixture Refresh](./phase-01-contract-lock-and-fixture-refresh.md) | Complete |
| 2 | [List and Follow Contract Alignment](./phase-02-list-and-follow-contract-alignment.md) | Complete |
| 3 | [Trace Detail and Export Hardening](./phase-03-trace-detail-and-export-hardening.md) | Complete |
| 4 | [Run Timeline Command](./phase-04-run-timeline-command.md) | Complete |
| 5 | [Docs Validation and Ship Readiness](./phase-05-docs-validation-and-ship-readiness.md) | Complete |

## Dependencies

| Relationship | Plan | Reason |
|--------------|------|--------|
| Blocks / supersedes | `260528-1357-fix-trace-details-by-id` | Older issue #17 plan used a flat fixture; this plan replaces it with current server `dev` envelope `{trace,spans}`. |
| Related | `260527-1412-domain-coverage-p6-backend-unblocked` | P6 added `traces follow`; this plan keeps that command but revalidates current contract drift. |

## Validation Gates

- Focused red/green: `go test -count=1 ./cmd -run 'TestTraces(List|Get|Follow|Timeline|Export)'`
- Package compile: `go test -count=1 ./cmd ./internal/client ./internal/output`
- Full repo: `go test -count=1 ./...`
- Static check: `go vet ./...`
- Compile: `go build ./...`
- Optional live smoke when credentials exist: list one trace, get it, export it, and query timeline when `run_id` is present.

Local note: current shell did not have `go` in PATH during scout. Implementer must either fix PATH or use the installed absolute Go binary if available.

## Success Criteria

- [x] `traces list` reads `{traces,total,limit,offset}` and renders non-empty table rows for server-shaped payloads.
- [x] `traces get` reads `{trace,spans}` and renders header + span tree from `span_type`, `id`, `parent_span_id`.
- [x] `traces export` validates and path-escapes trace ID before HTTP.
- [x] `traces follow` remains one-shot, preserves envelope in JSON/YAML, and table mode uses real server field names.
- [x] New `traces timeline <runID>` command reads `/v1/runs/{runID}/timeline`.
- [x] Tests use server-shaped fixtures, not the legacy flat trace fixture.
- [x] README/CHANGELOG/docs reflect actual commands and omit unsupported replay.

## Risk Assessment

| Risk | Mitigation |
|------|------------|
| CLI test helpers wrap responses in `ok/payload`, while server `writeJSON` often returns raw JSON | Add tests that exercise raw server JSON and envelope JSON where existing `HTTPClient` supports both. |
| Old docs mention fields like `input_tokens`/`cost`, but server trace fields are `total_input_tokens`/`total_cost` | Lock field mapping in table tests before implementation. |
| `cmd/traces.go` grows past maintainability target | Prefer one small split file for timeline only; avoid broad refactor. |
| Local Go unavailable in PATH | Capture this as environment setup; do not mark implementation complete until compile/test commands run somewhere valid. |
| Live fixture may contain sensitive prompt/user data | Use synthetic server-shaped fixtures for committed tests; live smoke output stays local/scrubbed. |
