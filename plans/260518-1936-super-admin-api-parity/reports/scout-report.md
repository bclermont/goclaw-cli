# Scout Report

## Summary

Compared GoClaw CLI command surface with GoClaw Gateway server routes after pulling both repos on 2026-05-18.

- CLI repo: `D:\www\nextlevelbuilder\goclaw-cli`, branch `feat/ai-first-cli-expansion`, head `697bba0`.
- Server repo: `D:\www\digitop\goclaw`, branch `dev`, head `99460a7f`.
- Existing CLI command count is high, but several operational routes are still missing or incorrectly wired.

## Relevant Files

| File | Finding |
|---|---|
| `cmd/api_keys.go` | `api-keys revoke` sends DELETE to `/v1/api-keys/{id}`. Server expects POST `/v1/api-keys/{id}/revoke`. |
| `cmd/packages.go` | Missing package update lifecycle commands. Existing commands stop at list/install/uninstall/runtimes/deny-groups/github-releases. |
| `cmd/admin_tts_media.go` | `media upload` is a placeholder, not a multipart upload. TTS lacks config/capabilities/synthesize HTTP routes. |
| `cmd/storage.go` | Missing storage upload and move commands. |
| `cmd/mcp_grants_requests.go` | Missing MCP per-user credential routes. |
| `cmd/admin_credentials_grants.go` | Missing `env:reveal` route exposure. |
| `D:\www\digitop\goclaw\internal\http\gateway_upgrade.go` | Server exposes gateway release upgrade status/start. CLI lacks command group. |
| `D:\www\digitop\goclaw\internal\http\workstations.go` | Server exposes workstation CRUD, permissions, activity. CLI lacks command group. |
| `D:\www\digitop\goclaw\internal\http\webhooks_admin.go` | Server exposes webhook admin CRUD. CLI lacks command group. |

## Route Gap Priority

P0:
- API key revoke route fix.
- Gateway upgrade status/start.
- Package updates list/refresh/apply/apply-all.
- Workstations CRUD/permissions/activity.

P1:
- Webhooks admin CRUD.
- MCP user credentials.
- CLI credential grant env reveal.
- Real media upload.
- TTS HTTP config/capabilities/synthesize.

P2:
- Storage upload/move.
- Contacts unmerge.
- Tenant-users list.
- Writer groups.
- Team attachment download.

Deferred:
- `/v1/chat/completions` and `/v1/responses`; existing `chat` command covers CLI interaction.
- Runtime webhook send helpers until auth/header contract is intentionally designed.

## Unresolved Questions

- None blocking.
