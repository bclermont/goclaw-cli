---
phase: 1
title: "Contract Audit and TDD Harness"
status: completed
effort: "4h"
---

# Phase 1: Contract Audit and TDD Harness

## Overview

Priority: P1
Current status: completed

Create the safety net before implementation. The existing CLI claims broad API coverage, but server routes drifted. This phase turns the route comparison into executable tests and a compact parity report.

## Context Links

- [Plan Overview](./plan.md)
- [Scout Report](./reports/scout-report.md)
- [Endpoint Inventory](./research/server-endpoint-inventory.md)
- `D:\www\nextlevelbuilder\goclaw-cli\README.md`
- `D:\www\nextlevelbuilder\goclaw-cli\docs\code-standards.md`
- `D:\www\digitop\goclaw\internal\http\api_keys.go`
- `D:\www\digitop\goclaw\internal\http\packages.go`
- `D:\www\digitop\goclaw\internal\http\workstations.go`

## Key Insights

- Route drift exists now: CLI revokes API keys with DELETE while server registers POST revoke.
- Some README claims are ahead of implementation, for example media upload exists as a command but currently does not upload.
- Server endpoints are split across REST and WS; command coverage cannot be checked by OpenAPI alone.
- `docs/development-rules.md` is missing. Use `CLAUDE.md`, README, code standards, and current command patterns.

## Requirements

- Capture current server route inventory for P0/P1 domains.
- Add tests before changing behavior.
- Keep tests deterministic with `httptest.NewServer`.
- Avoid real server dependency.

## Architecture

```
server route source -> parity matrix -> failing command tests -> implementation phases
```

Test harness should validate:
- HTTP method and path.
- JSON body shape.
- Multipart request shape where relevant.
- Confirmation gate behavior before network call.

## Related Code Files

| File | Action |
|---|---|
| `cmd/api_keys_test.go` | create |
| `cmd/packages_updates_test.go` | create |
| `cmd/workstations_test.go` | create |
| `cmd/gateway_upgrade_test.go` | create |
| `plans/260518-1936-super-admin-api-parity/research/server-endpoint-inventory.md` | create/update |
| `plans/260518-1936-super-admin-api-parity/reports/scout-report.md` | create/update |

## Implementation Steps

1. Add test utilities only if existing tests do not already cover setting global `cfg`, `printer`, and mock HTTP server.
2. Write failing test for API key revoke method/path.
3. Write contract tests for package update endpoints, gateway upgrade, and workstation CRUD path building.
4. Write non-network confirmation tests for delete/revoke/reveal flows.
5. Add parity report update command as a manual report step, not as runtime CLI code.

## Tests Before

- `go test ./cmd -run TestAPIKeysRevokeUsesPostRevoke`
- `go test ./cmd -run TestGatewayUpgrade`
- `go test ./cmd -run TestPackagesUpdates`
- `go test ./cmd -run TestWorkstations`

Expected at start: selected tests fail or do not compile until commands are added.

## Refactor

- Keep shared test helpers local to `cmd/*_test.go` unless duplicated three times.
- Do not introduce a generic command framework.
- Reuse `newHTTP` and `internal/client.HTTPClient`.

## Tests After

- All focused tests pass.
- `go test ./cmd -run "Test(APIKeys|GatewayUpgrade|PackagesUpdates|Workstations)"`

## Todo List

- [x] Add route contract tests for P0 bug and new command groups.
- [x] Capture route inventory report.
- [x] Confirm missing `docs/development-rules.md` in plan notes.
- [x] Keep command files below 200 lines or split by domain.

## Success Criteria

- [x] Tests fail for real missing behavior before implementation.
- [x] Test names map to server route contracts.
- [x] No test requires a live GoClaw server.

## Risk Assessment

- Risk: overbuilding a route scanner. Mitigation: keep report manual and tests focused.
- Risk: brittle tests tied to command text. Mitigation: assert method/path/body, not help formatting.

## Security Considerations

- Secret-bearing tests must use dummy values only.
- Do not commit `.env` or real tokens.

## Regression Gate

```powershell
go test ./cmd -run "Test(APIKeys|GatewayUpgrade|PackagesUpdates|Workstations)"
go build ./...
```

## Next Steps

- Phase 2 can start after API key and package/gateway tests exist.
