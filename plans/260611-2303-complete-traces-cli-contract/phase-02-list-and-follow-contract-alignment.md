---
phase: 2
title: "List and Follow Contract Alignment"
status: complete
priority: P1
effort: "4h"
dependencies: [1]
---

# Phase 2: List and Follow Contract Alignment

## Overview

Fix trace list and follow output against server `dev` contracts. Keep `follow` one-shot; do not introduce streaming or loops.

## Requirements

- Functional: `traces list` parses server object envelope and supports current server filters.
- Functional: `traces follow` preserves existing one-shot behavior and uses current field names in table output.
- Non-functional: JSON/YAML modes preserve raw payload shape for automation.

## Architecture

`traces list` should decode an object with `traces`, `total`, `limit`, and `offset`. Table mode renders `traces`; JSON/YAML prints the full object. `traces follow` already decodes an object; only field mapping and tests need hardening.

## Related Code Files

- Modify: `/Users/duynguyen/.codex/worktrees/202b/goclaw-cli/cmd/traces.go`
- Modify: `/Users/duynguyen/.codex/worktrees/202b/goclaw-cli/cmd/traces_follow.go`
- Modify: `/Users/duynguyen/.codex/worktrees/202b/goclaw-cli/cmd/p3_commands_test.go`
- Modify: `/Users/duynguyen/.codex/worktrees/202b/goclaw-cli/cmd/traces_follow_test.go`

## Tests Before

1. Update `TestTracesListAddsP3Filters` so response is raw server object, not top-level array.
2. Add tests for new list filters:
   - `--session-key` -> `session_key`
   - `--user` -> `user_id`
   - `--channel` -> `channel`
   - `--offset` -> `offset`
3. Add/keep follow one-request test with atomic call count.
4. Add table tests expecting `TOTAL_INPUT_TOKENS`, `TOTAL_OUTPUT_TOKENS`, `TOTAL_COST` or agreed concise aliases sourced from server fields.

## Refactor

- Add a small local helper only if needed:
  `traceRowsFromEnvelope(data json.RawMessage) (map[string]any, []any, error)`.
- Remove or deprecate unsupported `--root-only`; either no-op with hidden deprecation or delete if tests/docs allow.
- Do not pass `since` to `traces list` unless server supports it; keep `since` only on `traces follow`.

## Tests After

- JSON mode for `traces list` includes `total`, `limit`, and `offset`.
- Table mode for `traces list` and `traces follow` renders rows from `traces`.

## Implementation Steps

1. Decode list response into map and read `traces`.
2. Adjust table columns to server field names:
   `ID`, `AGENT`, `STATUS`, `DURATION_MS`, `TOTAL_INPUT_TOKENS`, `TOTAL_OUTPUT_TOKENS`, `TOTAL_COST`.
3. Add list flags for server-supported filters.
4. Revisit README examples that mention `--root-only` or `--since` on list.
5. Ensure follow still validates exactly one target and RFC3339 `since`.

## Success Criteria

- [x] `traces list` works with raw server object response.
- [x] `traces list -o json` prints the object, not just the array.
- [x] Current server filters are available.
- [x] `traces follow` still makes one HTTP request and no watch loop.
- [x] Unsupported list flags removed from docs/help or explicitly hidden/deprecated.

## Risk Assessment

Risk: removing `--root-only` breaks scripts.
Mitigation: prefer hidden/deprecated flag that returns validation error explaining server no longer supports it, unless maintainer chooses deletion.

## Regression Gate

- `go test -count=1 ./cmd -run 'TestTracesList|TestTracesFollow'`
