---
phase: 4
title: "Credential and Integration Surfaces"
status: completed
effort: "5h"
---

# Phase 4: Credential and Integration Surfaces

## Overview

Priority: P1
Current status: completed

Expose secure credential and integration surfaces needed by super-admin agents: webhook admin, MCP per-user credentials, and guarded CLI credential env reveal.

## Context Links

- `D:\www\digitop\goclaw\internal\http\webhooks_admin.go`
- `D:\www\digitop\goclaw\internal\http\webhooks_message.go`
- `D:\www\digitop\goclaw\internal\http\webhooks_llm.go`
- `D:\www\digitop\goclaw\internal\http\mcp_user_credentials.go`
- `D:\www\digitop\goclaw\internal\http\secure_cli_agent_grants.go`
- `D:\www\nextlevelbuilder\goclaw-cli\cmd\admin_credentials_grants.go`
- `D:\www\nextlevelbuilder\goclaw-cli\cmd\mcp_grants_requests.go`

## Key Insights

- Webhook runtime endpoints are integration entrypoints, not normal CLI sends. Admin CRUD is the required CLI surface.
- MCP user credentials are separate from MCP grants.
- `env:reveal` is intentionally sensitive and rate-limited server-side. CLI must make intent explicit.

## Requirements

- Add `goclaw webhooks create|list|get|update|rotate|delete`.
- Add optional `webhooks send-message` and `send-llm` only if the intended auth headers can be provided safely.
- Add `goclaw mcp user-credentials get|set|delete`.
- Add `goclaw credentials agent-grants reveal-env <credentialID> <grantID>` with explicit reveal gating.
- Secret reveal must require both intent confirmation and explicit output permission.

## Architecture

```
webhooks -> REST /v1/webhooks/*
mcp user-credentials -> REST /v1/mcp/servers/{id}/user-credentials
credentials reveal -> REST /v1/cli-credentials/{id}/agent-grants/{grantId}/env:reveal
```

## Related Code Files

| File | Action |
|---|---|
| `cmd/webhooks.go` | create |
| `cmd/webhooks_test.go` | create |
| `cmd/mcp_user_credentials.go` | create or extend |
| `cmd/mcp_user_credentials_test.go` | create |
| `cmd/admin_credentials_grants.go` | modify |
| `cmd/admin_credentials_grants_test.go` | create |

## Implementation Steps

1. Add tests for webhook CRUD request shapes.
2. Implement webhooks admin commands with `--body` for update and create to avoid schema churn.
3. Implement rotate with table-mode secret warning.
4. Implement MCP user credential commands with `--body` JSON for set.
5. Add env reveal command:
   - require `--yes`
   - require `--show-secrets` or equivalent explicit flag
   - JSON mode can emit server payload only when explicit reveal flag is set
   - table mode must refuse to print secret values unless explicit reveal flag is set; prefer metadata-only output otherwise
6. Confirm all commands return errors centrally.

## Tests Before

- Webhook command tests fail because command does not exist.
- MCP user credential command tests fail because routes are unexposed.
- Env reveal test fails because subcommand missing.
- Env reveal tests must prove no network call without `--yes` and no secret output without explicit reveal flag.

## Refactor

- Keep sensitive command logic close to existing credential files.
- Avoid central secret masking framework unless needed by multiple commands.

## Tests After

- `go test ./cmd -run "TestWebhooks|TestMCPUserCredentials|TestAdminCredentialsGrantReveal"`
- Help smoke: `go run . webhooks --help`, `go run . mcp user-credentials --help`.

## Todo List

- [x] Webhook admin CRUD.
- [x] Webhook rotate/delete confirmations.
- [x] MCP user credential get/set/delete.
- [x] Secure CLI env reveal with explicit gate.
- [x] Add table-mode and JSON-mode tests for reveal output behavior.

## Success Criteria

- [x] Super-admin can manage webhook integrations from CLI.
- [x] MCP per-user secrets can be managed without dashboard.
- [x] Secret reveal has intentional UX and test coverage.

## Risk Assessment

- Risk: secrets leak in stdout/history. Mitigation: prefer env/file input, explicit reveal flag, JSON only for automation.
- Risk: runtime webhook send requires HMAC headers not modeled in CLI. Mitigation: defer runtime send unless contract is clear.

## Security Considerations

- Never print secret values in success banners.
- Add `--yes` for delete/rotate/reveal.
- Do not accept raw secrets as positional args when `@file` or JSON body is safer.
- Treat reveal payload as sensitive even in tests; use dummy names that cannot resemble real tokens.

## Regression Gate

```powershell
go test ./cmd -run "TestWebhooks|TestMCPUserCredentials|TestAdminCredentialsGrantReveal"
go vet ./...
```

## Next Steps

- Phase 5 covers lower-risk media/storage/channel fillers.
