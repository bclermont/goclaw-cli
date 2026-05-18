---
phase: 2
title: "Critical Ops Fixes"
status: completed
effort: "5h"
---

# Phase 2: Critical Ops Fixes

## Overview

Priority: P1
Current status: completed

Fix broken API key revoke contract, then expose gateway release upgrade and package update lifecycle. These are super-admin operations agents need during GoClaw development and recovery.

## Context Links

- [Phase 1](./phase-01-contract-audit-and-tdd-harness.md)
- `D:\www\nextlevelbuilder\goclaw-cli\cmd\api_keys.go`
- `D:\www\nextlevelbuilder\goclaw-cli\cmd\packages.go`
- `D:\www\digitop\goclaw\internal\http\api_keys.go`
- `D:\www\digitop\goclaw\internal\http\gateway_upgrade.go`
- `D:\www\digitop\goclaw\internal\http\packages.go`

## Key Insights

- API key revoke is a live bug: server route is `POST /v1/api-keys/{id}/revoke`, CLI sends DELETE.
- Gateway upgrade requires `X-GoClaw-Upgrade-Token`; bearer auth alone is insufficient.
- Package update endpoints return operational payloads agents can inspect for partial failures.

## Requirements

- Fix `goclaw api-keys revoke`.
- Add `goclaw gateway upgrade status`.
- Add `goclaw gateway upgrade start --tag <tag> --upgrade-token <token>`.
- Add `goclaw packages updates list`, `refresh`, `apply`, `apply-all`.
- Use confirmations for write operations.
- `packages updates apply-all` must exit non-zero when `failed[]` is non-empty unless `--allow-partial` is passed.

## Architecture

```
goclaw api-keys revoke -> POST /v1/api-keys/{id}/revoke
goclaw gateway upgrade -> /v1/system/gateway/upgrade/*
goclaw packages updates -> /v1/packages/updates + /v1/packages/update
```

Gateway upgrade needs a custom-header request path because current HTTP helpers do not support arbitrary headers. Keep the helper local to `cmd/gateway_upgrade.go` unless another phase needs it.

## Related Code Files

| File | Action |
|---|---|
| `cmd/api_keys.go` | modify |
| `cmd/api_keys_test.go` | create |
| `cmd/gateway_upgrade.go` | create |
| `cmd/gateway_upgrade_test.go` | create |
| `cmd/packages_updates.go` | create |
| `cmd/packages_updates_test.go` | create |
| `README.md` | update later in Phase 6 |

## Implementation Steps

1. Change API key revoke call from DELETE to POST revoke endpoint.
2. Add `gatewayCmd` root if no existing gateway group exists.
3. Implement upgrade status with custom token header.
4. Implement upgrade start with tag validation mirroring server accepted values: `latest` or `vMAJOR.MINOR.PATCH[-beta.N|-rc.N]`.
   - Token precedence: `--upgrade-token` > `GOCLAW_UPGRADE_TRIGGER_TOKEN`.
   - Tests must assert `X-GoClaw-Upgrade-Token` is sent and never printed.
5. Implement package update commands:
   - `updates list` -> GET `/v1/packages/updates`
   - `updates refresh` -> POST `/v1/packages/updates/refresh`
   - `updates apply <source:name>` -> POST `/v1/packages/update`
   - `updates apply-all [packages...]` -> POST `/v1/packages/updates/apply-all`
6. Keep output as raw server payload in JSON/YAML mode and concise table in table mode.
7. For `apply-all`, inspect server payload. If `failed[]` has entries, return a non-zero error after printing payload unless caller passed `--allow-partial`.

## Tests Before

- API revoke path test fails on current DELETE behavior.
- Gateway/package command tests fail because commands do not exist.
- Package apply-all test must cover HTTP 200 with `failed[]` and expect non-zero command result.

## Refactor

- Split package update commands into `cmd/packages_updates.go` to keep `packages.go` small.
- Avoid adding a general REST wrapper just for custom headers unless media/storage later needs it.

## Tests After

- `go test ./cmd -run "TestAPIKeysRevoke|TestGatewayUpgrade|TestPackagesUpdates"`
- Verify `go run . packages updates --help`.
- Verify `go run . gateway upgrade --help`.
- Verify apply-all partial failure behavior with a mock server response.

## Todo List

- [x] Fix API key revoke route.
- [x] Add gateway upgrade command tree.
- [x] Add package update command tree.
- [x] Add tests for write confirmation and request shape.
- [x] Add tests for custom upgrade header and package partial-failure exit behavior.

## Success Criteria

- [x] API key revoke matches server contract.
- [x] Gateway upgrade can be driven by agents without shelling into server.
- [x] Package update partial failures are visible in machine-readable output.

## Risk Assessment

- Risk: upgrade token passed on CLI leaks in shell history. Mitigation: support env fallback and document env preference.
- Risk: package update all is destructive. Mitigation: require `--yes`.

## Security Considerations

- Do not print upgrade token.
- Do not log API keys.
- Require `--yes` for package update apply/apply-all and upgrade start.
- Prefer `GOCLAW_UPGRADE_TRIGGER_TOKEN` over CLI token arguments in docs/examples.

## Regression Gate

```powershell
go test ./cmd -run "TestAPIKeysRevoke|TestGatewayUpgrade|TestPackagesUpdates"
go vet ./...
go build ./...
```

## Next Steps

- Phase 3 can proceed in parallel after Phase 1 because it owns separate files.
