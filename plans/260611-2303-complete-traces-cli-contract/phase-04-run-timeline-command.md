---
phase: 4
title: "Run Timeline Command"
status: complete
priority: P1
effort: "4h"
dependencies: [1]
---

# Phase 4: Run Timeline Command

## Overview

Add CLI coverage for `GET /v1/runs/{runID}/timeline`, the trace-adjacent run archive endpoint present in server `dev`.

## Requirements

- Functional: command calls `/v1/runs/{runID}/timeline`.
- Functional: supports `--session-key`, `--limit`, and `--offset`.
- Functional: JSON/YAML mode preserves full response; table mode renders timeline items.
- Non-functional: no server mutation; read-only command.

## Architecture

Preferred command shape: `goclaw traces timeline <runID>`. It keeps trace-adjacent observability under `traces` without creating a new top-level `runs` group. If implementation finds an existing `runs`/`sessions` convention that fits better, document before changing.

## Related Code Files

- Modify: `/Users/duynguyen/.codex/worktrees/202b/goclaw-cli/cmd/traces.go`
- Create optional split: `/Users/duynguyen/.codex/worktrees/202b/goclaw-cli/cmd/traces_timeline.go`
- Create: `/Users/duynguyen/.codex/worktrees/202b/goclaw-cli/cmd/traces_timeline_test.go`
- Modify: `/Users/duynguyen/.codex/worktrees/202b/goclaw-cli/cmd/cmd_test.go`

## Tests Before

1. `TestTracesTimeline_BuildsPathAndQuery`:
   - run ID path-escaped
   - `session_key`, `limit`, `offset` query params set.
2. `TestTracesTimeline_JSONPreservesEnvelope`.
3. `TestTracesTimeline_TableRendersItems` with columns:
   `SEQ`, `TYPE`, `STATUS`, `TITLE`, `TOOL`, `TRACE_ID`, `SPAN_ID`, `CREATED_AT`.
4. Validation test for empty/path-like run IDs if `validateRunID` is added.

## Refactor

- Add `tracesTimelineCmd`.
- Add small table renderer for `items`.
- Reuse existing URL/query patterns from `sessions_follow.go` and `traces_follow.go`.

## Tests After

- New command appears under `traces --help`.
- Focused timeline tests pass.

## Implementation Steps

1. Add Cobra command with `Use: "timeline <runID>"`.
2. Add flags:
   - `--session-key`
   - `--limit`
   - `--offset`
3. Build path `/v1/runs/{url.PathEscape(runID)}/timeline`.
4. Decode response map.
5. Table mode renders `items`; non-table prints full map.
6. Register command in `tracesCmd.AddCommand(...)`.

## Success Criteria

- [x] CLI covers server run timeline endpoint.
- [x] Query params match server docs.
- [x] Table and JSON modes covered by tests.
- [x] Command docs make clear it is read-only archive view, not replay.

## Risk Assessment

Risk: command belongs under future `runs` group, not `traces`.
Mitigation: use `traces timeline` now because endpoint is registered by `TracesHandler` and returns trace/span IDs; do not create a broad `runs` group for one read command.

## Regression Gate

- `go test -count=1 ./cmd -run 'TestTracesTimeline|TestAllCommandsRegistered|TestCommandUseFields'`
