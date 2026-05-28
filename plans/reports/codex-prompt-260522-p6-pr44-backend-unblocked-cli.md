# Codex CLI Prompt: Implement goclaw-cli P6 backend-unblocked commands after backend PR 44

Workdir: `/Volumes/GOON/www/nlb/goclaw-cli`

You are working on `goclaw-cli`, the Go/Cobra CLI for GoClaw Gateway.

## Mission

Implement the P6 CLI commands now unblocked by backend PRs #37 and #44 in `digitopvn/goclaw`.

This is a CLI-consuming task only. Do not create backend stubs. Do not invent APIs.

## Backend evidence

Backend repo: `digitopvn/goclaw`

Already released in beta:
- PR #37: `https://github.com/digitopvn/goclaw/pull/37`
- Merge commit: `56e227c4030e85163cd882b29ab472f8ce3e1a27`
- Beta tag known to contain these APIs: `v3.12.0-beta.16`
- APIs:
  - `GET /v1/traces/follow`
  - `POST /v1/providers/{id}/reconnect`

Merged to `dev`, beta release may still be publishing:
- PR #44: `https://github.com/digitopvn/goclaw/pull/44`
- Merge commit: `43049d3b3fbb5f457477118252d1f21fdc0480de`
- As of prompt creation, latest listed backend beta was `v3.12.0-beta.18`, created before PR #44 landed.
- As of prompt creation, `Dev CI and Beta Release` for commit `43049d3b` was still pending:
  - `https://github.com/digitopvn/goclaw/actions/runs/26292214016`
- APIs:
  - `POST /v1/chat/sessions/{key}/branch`
  - `GET /v1/chat/sessions/{key}/history/follow`
  - `POST /v1/channels/instances/{id}/writers/test`
  - `GET /v1/activity/aggregate`
  - `GET /v1/logs/runtime/aggregate`

Before implementation, verify current backend release status:

```bash
gh run list --repo digitopvn/goclaw --branch dev --limit 5 \
  --json databaseId,name,status,conclusion,createdAt,headSha,url
gh release list --repo digitopvn/goclaw --limit 10
```

If PR #44 is not in a beta tag yet, CLI implementation may still proceed from merged `dev` contracts, but do not claim live beta support until a beta tag containing `43049d3b` exists.

Backend files proving contracts:
- `internal/http/traces.go`
- `internal/http/providers.go`
- `internal/http/sessions.go`
- `internal/http/channel_instances.go`
- `internal/http/activity.go`
- `internal/http/logs.go`
- `internal/http/openapi_spec.json`
- `docs/18-http-api.md`

## Current goclaw-cli repo context

Read first:
- `README.md`
- `CLAUDE.md` if present
- `AGENTS.md` if present
- Existing command patterns:
  - `cmd/traces.go`
  - `cmd/providers.go`
  - `cmd/providers_crud.go`
  - `cmd/sessions.go`
  - `cmd/channels_writers.go`
  - `cmd/admin_activity.go`
  - `cmd/logs.go`
  - `cmd/helpers.go`
  - `internal/client/http.go`
  - `internal/output/output.go`
  - command tests under `cmd/*_test.go`

Important local state warning:
- The current checkout may be on branch `feat/claude-skill-v0.1`.
- The current checkout may contain unrelated untracked files such as `.claude/` and `AGENTS.md`.
- Do not overwrite, delete, stage, or commit unrelated user files.
- Prefer creating a clean worktree or branch for this task before editing.

Suggested branch name:

```bash
codex/feat-p6-backend-unblocked-cli
```

## Scope: implement backend-unblocked P6 CLI surfaces

Implement exactly these command surfaces unless current CLI already contains one. If an equivalent command exists, preserve it and document the mapping instead of creating a duplicate command name.

### 1. Trace polling-friendly follow

Command:

```bash
goclaw traces follow --session-key <key> [--since <RFC3339>] [--limit N] [--include-spans] [--status <status>] [--channel <channel>] [-o json|yaml|table]
goclaw traces follow --agent <uuid-or-key> [same flags]
```

Endpoint:

```http
GET /v1/traces/follow
```

Query contract:
- require exactly one of `session_key` or `agent_id`
- optional: `status`, `channel`, `since`, `limit`, `include_spans`
- `since` must be RFC3339 if provided
- server default `limit=50`, server max `200`

Response after HTTP envelope unwrap:

```json
{
  "traces": [],
  "spans_by_trace_id": {},
  "server_time": "2026-05-21T00:00:00Z",
  "next_since": "2026-05-21T00:00:00Z",
  "limit": 50
}
```

CLI behavior:
- JSON/YAML: print full response.
- Table: print trace rows similar to `traces list`.
- Include useful columns: `TRACE_ID`, `AGENT`, `STATUS`, `DURATION_MS`, `INPUT_TOKENS`, `OUTPUT_TOKENS`, `COST`.
- One request only. Do not implement a watch loop in this slice.

Tests:
- session-key query builds correct path/query.
- agent query builds correct path/query.
- missing both target flags errors before HTTP call.
- setting both target flags errors before HTTP call.
- JSON output preserves `next_since` and `spans_by_trace_id`.

### 2. Provider reconnect

Command:

```bash
goclaw providers reconnect <provider-id> [-o json|yaml|table]
```

Endpoint:

```http
POST /v1/providers/{id}/reconnect
```

Contract:
- admin-only backend permission.
- no request body by default.
- do not send `{ "verify": true }`.
- do not add `--verify`; users should call `goclaw providers verify <id>` separately.

Response after HTTP envelope unwrap:

```json
{
  "status": "reconnected",
  "provider": {},
  "registry_updated": true,
  "cache_invalidated": true
}
```

Status enum:
- `reconnected`
- `disabled`
- `not_registered`

CLI behavior:
- JSON/YAML: print full response.
- Table: print status, registry_updated, cache_invalidated.
- Path-escape provider ID like existing provider commands.

Tests:
- POST path is `/v1/providers/{escaped-id}/reconnect`.
- no request body is sent by default.
- JSON/table output preserves `registry_updated` and `cache_invalidated`.

### 3. Session branch at message index

Command:

```bash
goclaw sessions branch <session-key> --up-to-index <n> [--new-session-key <key>] [--label <label>] [--metadata k=v]... [-o json|yaml|table]
```

Endpoint:

```http
POST /v1/chat/sessions/{key}/branch
```

Request:

```json
{
  "new_session_key": "optional",
  "up_to_index": 12,
  "label": "optional",
  "metadata": {
    "source": "cli"
  }
}
```

Response:

```json
{
  "ok": true,
  "source_key": "agent:default:ws:direct:abc",
  "session_key": "agent:default:branch:direct:uuid",
  "copied_messages": 12,
  "total_messages": 24,
  "label": "optional"
}
```

CLI behavior:
- `--up-to-index` is required and must be >= 0.
- `--metadata` parses repeated `key=value`; reject malformed entries before HTTP call.
- path-escape the source session key.
- JSON/YAML: full response.
- Table: source, new session key, copied/total, label.

Tests:
- required `--up-to-index`.
- negative index rejected before HTTP call.
- request body shape matches backend contract.
- path escaping works for session keys containing `:` and `/`.
- conflict response maps to the existing CLI error handling pattern.

### 4. Session history follow by cursor

Command:

```bash
goclaw sessions follow <session-key> [--cursor <n>] [--limit <n>] [-o json|yaml|table]
```

Endpoint:

```http
GET /v1/chat/sessions/{key}/history/follow
```

Query:
- `cursor`, default `0`, must be >= 0.
- `limit`, default `50`, max server `200`, must be > 0.

Response:

```json
{
  "session_key": "agent:default:ws:direct:abc",
  "cursor": 12,
  "next_cursor": 18,
  "total": 18,
  "messages": [],
  "reset": false,
  "updated": "2026-05-22T13:00:00Z"
}
```

CLI behavior:
- One polling request only. Do not implement SSE/WS watch in this slice.
- JSON/YAML: full response.
- Table: print cursor, next_cursor, total, reset, and compact message rows.
- Use existing output helpers; do not create a separate renderer unless needed.

Tests:
- query path includes cursor and limit.
- negative cursor rejected before HTTP call.
- non-positive limit rejected before HTTP call.
- JSON output preserves `reset`, `next_cursor`, and messages.

### 5. Channel writer permission test

Command:

```bash
goclaw channels writers test <instance-id> --group-id <group-scope> --user-id <user-id> [-o json|yaml|table]
```

Endpoint:

```http
POST /v1/channels/instances/{id}/writers/test
```

Request:

```json
{
  "group_id": "group:telegram:-100123",
  "user_id": "386246614"
}
```

Expected response shape:

```json
{
  "allowed": true,
  "reason": "writer",
  "instance_id": "uuid",
  "agent_id": "uuid",
  "group_id": "group:telegram:-100123",
  "user_id": "386246614",
  "writer_count": 3
}
```

Known `reason` values:
- `writer`
- `not_writer`
- `no_writers_configured`
- `invalid_group`

CLI behavior:
- require `--group-id` and `--user-id`.
- POST body only has `group_id` and `user_id`.
- JSON/YAML: full response.
- Table: allowed, reason, writer_count, group_id, user_id.

Tests:
- missing required flags fail before HTTP call.
- request path/body exact.
- table/json output keeps `allowed`, `reason`, and `writer_count`.

### 6. Activity log aggregate

Command:

```bash
goclaw activity aggregate --group-by <action|actor_type|entity_type|actor_id> [--from <RFC3339>] [--to <RFC3339>] [--limit <n>] [--actor-type <v>] [--actor-id <v>] [--action <v>] [--entity-type <v>] [--entity-id <v>] [-o json|yaml|table]
```

Endpoint:

```http
GET /v1/activity/aggregate
```

Contract:
- `group_by` required.
- valid values: `action`, `actor_type`, `entity_type`, `actor_id`.
- backend restricts `group_by=actor_id` to admin.
- optional filters: `from`, `to`, `limit`, `actor_type`, `actor_id`, `action`, `entity_type`, `entity_id`.
- non-admin backend callers are scoped to resolved user context.

Response:

```json
{
  "source": "activity",
  "group_by": "action",
  "total": 10,
  "limit": 50,
  "from": "2026-05-22T00:00:00Z",
  "to": "2026-05-23T00:00:00Z",
  "buckets": [
    {"key": "session.branch", "count": 7, "last_seen": "2026-05-22T11:00:00Z"}
  ]
}
```

CLI behavior:
- validate `--group-by` against allowed values before HTTP call.
- validate `--from` and `--to` as RFC3339 if provided.
- JSON/YAML: full response.
- Table: key, count, last_seen.

Tests:
- missing/invalid group-by rejected before HTTP call.
- filters build correct query string.
- JSON output preserves source, group_by, total, buckets.

### 7. Runtime log aggregate

Command:

```bash
goclaw logs aggregate [--group-by <level|source>] [--level <debug|info|warn|error>] [--source <source>] [--from <RFC3339>] [-o json|yaml|table]
```

Endpoint:

```http
GET /v1/logs/runtime/aggregate
```

Contract:
- admin-only backend permission.
- source is runtime ring buffer, not durable audit logs.
- `group_by` default `level`; valid values: `level`, `source`.
- optional filters: `level`, `source`, `from`.

Response:

```json
{
  "source": "runtime",
  "retention": "ring_buffer",
  "capacity": 100,
  "sample_size": 25,
  "group_by": "level",
  "buckets": [
    {"key": "warn", "count": 3, "last_seen": 1760000000000}
  ]
}
```

CLI behavior:
- JSON/YAML: full response.
- Table: key, count, last_seen, plus source/retention/capacity/sample_size summary if local output helpers support it.
- Do not confuse this with `goclaw logs tail`, which is WebSocket streaming.

Tests:
- default group_by omitted or set to `level` based on existing CLI style.
- invalid group-by rejected before HTTP call.
- filters build correct query.
- JSON output preserves retention, capacity, sample_size.

## Explicitly out of scope

Do not add stubs, placeholders, hidden flags, docs, or command names for APIs that still do not exist.

Still out of scope:
- `POST /v1/traces/{id}/replay`
- `GET /v1/logs/aggregate`
- WebSocket `chat.history.delta`
- SSE chat history follow
- any long-running watch loop for history follow or traces follow
- live backend smoke if no beta tag containing PR #44 exists yet

## Implementation workflow

Use TDD.

1. Scout:
   - read README/CLAUDE/AGENTS.
   - inspect command patterns and test helpers.
   - verify whether PR #37 commands are already implemented.
   - verify backend beta status for PR #44.
2. Branch/worktree:
   - avoid the dirty `feat/claude-skill-v0.1` checkout.
   - create/switch to a clean feature branch or worktree.
3. Tests first:
   - add focused command tests for path/query/body/validation/output.
   - tests should fail before implementation.
4. Implement smallest changes:
   - use existing `internal/client.HTTPClient` helpers.
   - use existing output helpers.
   - keep command modules small; follow existing file layout.
5. Validation:
   - `go test ./...`
   - `go vet ./...`
   - `go build ./...`
6. Red-team diff:
   - wrong endpoint path.
   - body accidentally sent to provider reconnect.
   - command name collision.
   - output mode regression.
   - missing required client-side validation.
   - accidental implementation of still-blocked replay/generic logs APIs.
   - dirty worktree accidentally staged.
7. Ship:
   - commit with conventional message.
   - push branch.
   - open PR against the correct integration branch for this repo.

## Acceptance criteria

Expected output is one focused `goclaw-cli` PR that adds CLI commands for every backend-unblocked P6 surface above.

Done means:
- all seven command surfaces exist, unless already present and explicitly mapped.
- each command has focused tests.
- all tests/build/vet pass.
- README/help text reflects the new commands if this repo normally updates README for command coverage.
- no CLI command exists for trace replay or generic `/v1/logs/aggregate`.
- PR body lists the backend PR/tag evidence and notes whether PR #44 beta tag was available at implementation time.

