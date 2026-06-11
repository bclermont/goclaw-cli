---
phase: 3
title: "Trace Detail and Export Hardening"
status: complete
priority: P1
effort: "5h"
dependencies: [1, 2]
---

# Phase 3: Trace Detail and Export Hardening

## Overview

Fix `traces get` for the real `{trace,spans}` response and harden `traces export` path handling.

## Requirements

- Functional: `traces get <traceID>` prints the full server payload in JSON/YAML and a useful header + span tree in table mode.
- Functional: `traces export <traceID>` validates IDs before HTTP and path-escapes route params.
- Non-functional: keep error routing through central error handler.

## Architecture

Table rendering should unwrap `trace` and `spans` only for display. JSON/YAML should preserve `{trace,spans}` exactly. Span tree should use server fields: `id`, `parent_span_id`, `span_type`, `name`, `duration_ms`, `status`.

## Related Code Files

- Modify: `/Users/duynguyen/.codex/worktrees/202b/goclaw-cli/cmd/traces.go`
- Modify: `/Users/duynguyen/.codex/worktrees/202b/goclaw-cli/cmd/traces_get_test.go`
- Modify: `/Users/duynguyen/.codex/worktrees/202b/goclaw-cli/cmd/testdata/trace_detail_get.json`
- Possibly create: `/Users/duynguyen/.codex/worktrees/202b/goclaw-cli/cmd/traces_render_helpers.go` only if `cmd/traces.go` becomes too large.

## Tests Before

1. Replace flat trace detail fixture with raw `{trace,spans}` fixture.
2. Add test that table mode includes trace `id`, `run_id`, `session_key`, `status`, `total_*` values.
3. Add test that span tree uses `id`/`parent_span_id`, not `span_id`.
4. Add export validation tests:
   - bad IDs do not call HTTP
   - valid ID path is escaped
5. Keep existing exit-code tests for 404/403/400/5xx.

## Refactor

- Change `renderTraceTable` to accept full response map, unwrap `trace` and `spans`.
- Update `buildSpanTree` to use `id` and `span_type`.
- Reuse `validateTraceID` for `export`.
- Use `url.PathEscape(id)` in export path.

## Tests After

- `traces get -o json` preserves top-level `trace` and `spans`.
- `traces get -o table` no longer prints empty headers for server-shaped payload.
- `traces export ../bad` returns validation error before HTTP.

## Implementation Steps

1. Decode `traces get` response with checked `json.Unmarshal`.
2. If response contains `trace` object, render that; if legacy flat response appears, optionally keep backward-compatible fallback.
3. Update span tree mapping:
   - key: `id`
   - parent: `parent_span_id`
   - type: `span_type`
4. Harden export ID before file path default is computed.
5. Ensure output filename for export remains safe: use validated ID, not raw arg.

## Success Criteria

- [x] `traces get` works against server-shaped `{trace,spans}`.
- [x] JSON/YAML mode round-trips server structure.
- [x] Table mode shows useful trace header and span tree.
- [x] `traces export` cannot send malformed/path-like IDs.
- [x] Existing exit-code contract remains unchanged.

## Risk Assessment

Risk: table render may expose untrusted LLM/tool text.
Mitigation: render compact fields only (`name`, `tool_name`, previews if explicitly chosen); avoid raw content dumps.

## Regression Gate

- `go test -count=1 ./cmd -run 'TestTracesGet|TestTracesExport'`
