# Codex CLI Prompt: Implement goclaw-cli P6 Backend-Unblocked Commands

Workdir: `/Volumes/GOON/www/nlb/goclaw-cli`

You are working on `goclaw-cli`, the Go/Cobra CLI for GoClaw Gateway.

## Mission

Implement only the P6 CLI commands now unblocked by backend PR #37 in `digitopvn/goclaw`.

Backend evidence:
- Backend repo: `digitopvn/goclaw`
- PR: `https://github.com/digitopvn/goclaw/pull/37`
- Merge commit: `56e227c4030e85163cd882b29ab472f8ce3e1a27`
- Beta tag containing these APIs: `v3.12.0-beta.16`
- Backend files proving contracts:
  - `internal/http/traces.go`
  - `internal/http/providers.go`
  - `internal/http/openapi_spec.json`
  - `docs/18-http-api.md`

Before implementation, verify current backend release status. If `v3.12.0-beta.16` release/assets are still publishing, note that, but CLI code may proceed because the tag already points to PR #37.

## Current goclaw-cli repo context

Read first:
- `README.md`
- `CLAUDE.md`
- `AGENTS.md`
- Existing command patterns in:
  - `cmd/traces.go`
  - `cmd/providers.go`
  - `cmd/providers_crud.go`
  - `cmd/helpers.go`
  - `internal/client/http.go`
  - command tests under `cmd/*_test.go`

Important local state warning:
- The checkout may already be on a feature branch and may have unrelated untracked files such as `.claude/` or `AGENTS.md`.
- Do not overwrite or delete unrelated user/untracked files.
- If the working tree is dirty in unrelated files, either work around them safely or create a clean worktree/branch for this task.

## Scope: implement exactly 2 command surfaces

### 1. Trace polling-friendly follow

Add a CLI command under `traces`, probably:

```bash
goclaw traces follow --session-key <key> [--since <RFC3339>] [--limit N] [--include-spans] [--status <status>] [--channel <channel>] [-o json|yaml|table]
goclaw traces follow --agent <uuid> [same flags]
```

Backend endpoint:

```http
GET /v1/traces/follow
```

Query contract:
- Require one of:
  - `session_key`
  - `agent_id`
- Optional:
  - `status`
  - `channel`
  - `since` as RFC3339
  - `limit`, default server 50, max server 200
  - `include_spans`, boolean, default false
- Non-admin callers only receive their own traces. Admin may pass `user_id`, but do not add `--user-id` unless existing CLI conventions already expose admin trace filters.

Response payload after `internal/client.HTTPClient` envelope unwrap:

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
- For JSON/YAML, print the full response map.
- For table, print trace rows similar to `traces list`, and include enough fields to be useful:
  - `TRACE_ID`, `AGENT`, `STATUS`, `DURATION_MS`, `INPUT_TOKENS`, `OUTPUT_TOKENS`, `COST`
- Validate that exactly one target style is provided if that matches local CLI style; at minimum return a clear error if neither `--session-key` nor `--agent` is provided.
- Do not implement long-lived watch loops unless existing CLI has a clear polling/watch convention. This backend endpoint is polling-friendly; one request is acceptable for this slice.

Tests to add:
- Command builds correct path for `session_key`, `since`, `limit`, `include_spans`.
- Command builds correct path for `agent_id`.
- Missing both `--session-key` and `--agent` returns error before HTTP call.
- JSON output preserves `next_since` and `spans_by_trace_id`.

### 2. Provider reconnect

Add a CLI command under `providers`, probably:

```bash
goclaw providers reconnect <provider-id> [-o json|yaml|table]
```

Backend endpoint:

```http
POST /v1/providers/{id}/reconnect
```

Auth/permission:
- Backend requires admin role.

Request contract:
- No body by default.
- Do not send `{ "verify": true }`.
- Do not add a `--verify` flag. Backend explicitly rejects verify-on-reconnect; users should call `goclaw providers verify <id>` separately.

Response payload after envelope unwrap:

```json
{
  "status": "reconnected",
  "provider": {},
  "registry_updated": true,
  "cache_invalidated": true
}
```

`status` enum:
- `reconnected`
- `disabled`
- `not_registered`

CLI behavior:
- For JSON/YAML, print the full response map.
- For table, print a small single-row table or concise success message containing `status`, `registry_updated`, and `cache_invalidated`.
- Path-escape the provider ID exactly like existing provider commands.

Tests to add:
- Command uses `POST /v1/providers/{escaped-id}/reconnect`.
- Command sends no request body by default.
- JSON/table handling does not drop `registry_updated` or `cache_invalidated`.

## Explicitly out of scope

Do not add stubs, placeholders, hidden flags, or commands for these still-blocked P6 backend items:
- `POST /v1/traces/{id}/replay`
- `GET /v1/logs/aggregate`
- `POST /v1/channels/instances/{id}/writers/test`
- `POST /v1/chat/sessions/{key}/branch`
- WebSocket `chat.history.delta`
- SSE/HTTP chat history follow

If you discover an existing CLI command already maps to either new endpoint, preserve it and document the mapping instead of duplicating command names.

## Workflow

Use TDD:
1. Scout current command/test helpers.
2. Add focused failing tests for the two command surfaces.
3. Implement the smallest command changes.
4. Run:
   - `go test ./...`
   - `go vet ./...`
   - `go build ./...`
5. Red-team the diff for:
   - wrong endpoint path
   - body accidentally sent to reconnect
   - command name collision
   - output mode regression
   - accidental implementation of blocked endpoints
6. Commit with a clean conventional message.
7. Open PR against the correct integration branch used by this repo.

## Expected result

One small PR in `goclaw-cli` that unblocks:
- `goclaw traces follow`
- `goclaw providers reconnect`

No backend stubs. No fake commands for absent APIs.
