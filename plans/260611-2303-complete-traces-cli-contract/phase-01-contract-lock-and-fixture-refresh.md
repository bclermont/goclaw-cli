---
phase: 1
title: "Contract Lock and Fixture Refresh"
status: complete
priority: P1
effort: "3h"
dependencies: []
---

# Phase 1: Contract Lock and Fixture Refresh

## Overview

Lock current trace API contracts before touching code. Replace stale assumptions from prior trace plans with server-shaped fixtures and failing tests.

## Requirements

- Functional: capture server `dev` shapes for `list`, `get`, `follow`, `export`, and run timeline.
- Functional: document unsupported trace replay as out of scope.
- Non-functional: no committed secrets or live trace content.

## Architecture

The CLI uses `internal/client.HTTPClient`, which supports both `{ok,payload}` envelopes and raw JSON. Server `writeJSON` returns raw JSON for trace handlers, so test fixtures must cover raw JSON directly.

## Related Code Files

- Modify: `/Users/duynguyen/.codex/worktrees/202b/goclaw-cli/cmd/traces_get_test.go`
- Modify: `/Users/duynguyen/.codex/worktrees/202b/goclaw-cli/cmd/traces_follow_test.go`
- Modify: `/Users/duynguyen/.codex/worktrees/202b/goclaw-cli/cmd/p3_commands_test.go`
- Create: `/Users/duynguyen/.codex/worktrees/202b/goclaw-cli/cmd/traces_timeline_test.go`
- Modify/Create: `/Users/duynguyen/.codex/worktrees/202b/goclaw-cli/cmd/testdata/trace_*.json`
- Read-only source: `/Volumes/GOON/www/digitop/goclaw/internal/http/traces.go`
- Read-only source: `/Volumes/GOON/www/digitop/goclaw/docs/18-http-api.md`

## Tests Before

1. Add red test: `TestTracesList_ServerEnvelope_TableRows` with raw response:
   `{"traces":[{"id":"...","agent_id":"...","status":"completed","total_input_tokens":10,"total_output_tokens":5,"total_cost":0.01}],"total":1,"limit":20,"offset":0}`.
2. Add red test: `TestTracesGet_ServerEnvelope_TableRendersTraceAndSpans` with raw response:
   `{"trace":{...},"spans":[{"id":"span-root","parent_span_id":null,"span_type":"agent","name":"agent.run"}]}`.
3. Add red test: `TestTracesTimeline_CommandMissing_RED` expecting the final command shape to exist and call `/v1/runs/{runID}/timeline`.
4. Add regression test that no command named `traces replay` exists.

## Refactor

- None in this phase except test fixture setup.

## Tests After

- Existing tests may fail; that is expected until later phases.
- Run focused tests and record red failures in a short report.

## Implementation Steps

1. Re-read current server route registration in `internal/http/traces.go`.
2. Update trace fixtures to use real JSON names: `id`, `total_input_tokens`, `total_output_tokens`, `total_cost`, `span_type`, `parent_span_id`.
3. Preserve old flat fixture only if needed for backward-compat test; mark it legacy.
4. Write failing tests before implementation.
5. Create `plans/260611-2303-complete-traces-cli-contract/reports/contract-lock.md` with exact route/field inventory.

## Success Criteria

- [x] Test fixtures match server `dev` response structures.
- [x] At least 3 red tests prove current CLI gaps.
- [x] Report lists unsupported `POST /v1/traces/{id}/replay` as absent on server.
- [x] No fixture contains bearer tokens, prompts, API keys, or real user IDs.

## Risk Assessment

Risk: accidental fake envelope via `okJSON` hides real raw server behavior.
Mitigation: write at least one raw JSON `httptest` response per trace command.

## Regression Gate

- `go test -count=1 ./cmd -run 'TestTraces(List|Get|Timeline|Replay)'`
