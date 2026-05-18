---
phase: 3
title: "Workstations Command Surface"
status: completed
effort: "6h"
---

# Phase 3: Workstations Command Surface

## Overview

Priority: P1
Current status: completed

Expose Standard-edition workstation management. This is the main missing surface for coding agents that need remote shell/workspace execution, agent links, permissions, and activity logs.

## Context Links

- `D:\www\digitop\goclaw\internal\http\workstations.go`
- `D:\www\digitop\goclaw\internal\gateway\methods\workstations.go`
- `D:\www\nextlevelbuilder\goclaw-cli\cmd\agents_links.go`
- `D:\www\nextlevelbuilder\goclaw-cli\cmd\tools_custom.go`

## Key Insights

- Server exposes REST CRUD plus permissions/activity.
- Server also exposes WS methods for link/unlink agent.
- `workstations.testConnection` currently may return not implemented depending edition/phase; CLI should pass through server error.

## Requirements

- Add `goclaw workstations list|get|create|update|delete|test`.
- Add `goclaw workstations permissions list|add|remove|toggle`.
- Add `goclaw workstations activity`.
- Add `goclaw workstations link-agent|unlink-agent` if WS methods are live.
- Preserve tenant/admin auth errors.
- Link/unlink must call WS `workstations.linkAgent` / `workstations.unlinkAgent` with `agentId` and `workstationId`.

## Architecture

```
cmd/workstations.go
cmd/workstations_permissions.go
    -> REST /v1/workstations/*
    -> WS workstations.linkAgent / unlinkAgent for agent binding
```

## Related Code Files

| File | Action |
|---|---|
| `cmd/workstations.go` | create |
| `cmd/workstations_permissions.go` | create |
| `cmd/workstations_test.go` | create |
| `cmd/workstations_permissions_test.go` | create |
| `README.md` | update in Phase 6 |

## Implementation Steps

1. Write tests for CRUD request method/path/body.
2. Implement command tree with JSON body flags:
   - create: `--key`, `--name`, `--backend-type`, `--metadata`, `--default-cwd`, `--default-env`
   - update: `--body` JSON object for advanced safety and low churn
3. Add delete confirmation before network call.
4. Implement permissions commands:
   - list `<workstationID>`
   - add `<workstationID> --pattern`
   - remove `<workstationID> <permID> --yes`
   - toggle `<workstationID> <permID> --enabled=true|false`
5. Implement activity with `--limit` and optional `--cursor`.
6. Add link/unlink commands using WS only after verifying parameter names in server method:
   - method `workstations.linkAgent`
   - method `workstations.unlinkAgent`
   - params `agentId`, `workstationId`

## Tests Before

- Failing command tests for create/list/delete and permission add/remove.
- Failing help smoke for `go run . workstations --help`.
- Failing WS request-shape test for link/unlink if those commands are included.

## Refactor

- Keep core CRUD and permissions in separate files.
- Use `--body` for update to avoid broad flag explosion.
- Reuse table rendering patterns from `tenants`, `channels`, and `mcp`.

## Tests After

- `go test ./cmd -run TestWorkstations`
- `go run . workstations --help`
- `go run . workstations permissions --help`

## Todo List

- [x] Add workstation root command.
- [x] Add CRUD.
- [x] Add permissions.
- [x] Add activity.
- [x] Add link/unlink if server method params verified.
- [x] Add tests that assert exact WS method and camelCase params for link/unlink.

## Success Criteria

- [x] Super-admin agent can inspect and mutate workstations from CLI.
- [x] Permission allowlist can be managed without dashboard.
- [x] Activity logs are accessible for debugging remote coding-agent behavior.

## Risk Assessment

- Risk: metadata schema differs by backend. Mitigation: CLI accepts JSON and lets server validate.
- Risk: wrong WS link parameter names. Mitigation: read server method before implementation and add request tests.
- Risk: treating link/unlink as REST because most workstation commands are REST. Mitigation: isolate WS link tests.

## Security Considerations

- Workstation metadata/default env may contain secrets. Do not pretty-print secrets beyond server sanitized view.
- Delete/remove commands require `--yes`.

## Regression Gate

```powershell
go test ./cmd -run TestWorkstations
go build ./...
```

## Next Steps

- Phase 4 can proceed independently for credential/integration domains.
