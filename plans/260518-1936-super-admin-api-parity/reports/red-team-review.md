# Red Team Review

## Summary

Adversarial review for `plans/260518-1936-super-admin-api-parity/`.

Result: 8 findings, 7 accepted, 1 rejected. Main theme: the plan is directionally right, but a few command contracts need sharper gates so implementation does not accidentally create secret leaks, false success exits, or half-tested custom HTTP behavior.

## Findings

| ID | Severity | Lens | Finding | Evidence | Disposition |
|---|---|---|---|---|---|
| RT-01 | Critical | Contract | API key revoke is broken and must stay first. The plan has it, but implementation must not batch it behind larger features. | `cmd/api_keys.go:116`, `D:\www\digitop\goclaw\internal\http\api_keys.go:33` | Accept |
| RT-02 | High | Security | Gateway upgrade needs a custom header path. Existing `HTTPClient` helpers only set Authorization/Content-Type, so a naive implementation will silently omit `X-GoClaw-Upgrade-Token`. | `internal/client/http.go:70`, `internal/client/http.go:77`, `internal/client/http.go:79`, `D:\www\digitop\goclaw\internal\http\gateway_upgrade.go:24`, `D:\www\digitop\goclaw\internal\http\gateway_upgrade.go:167` | Accept |
| RT-03 | High | Automation | `packages updates apply-all` must not exit 0 when `failed[]` is non-empty unless caller explicitly allows partial failure. Server returns HTTP 200 by design, so CLI must create the automation failure signal. | `D:\www\digitop\goclaw\internal\http\packages_updates.go:286`, `D:\www\digitop\goclaw\internal\http\packages_updates.go:369`, `D:\www\digitop\goclaw\internal\http\packages_updates.go:370`, `D:\www\digitop\goclaw\internal\http\packages_updates.go:381` | Accept |
| RT-04 | High | Security | `env:reveal` must not print secrets in table mode at all unless the command has both `--yes` and `--show-secrets`; JSON automation should still require explicit reveal. | `D:\www\digitop\goclaw\internal\http\secure_cli_agent_grants.go:72`, `D:\www\digitop\goclaw\internal\http\secure_cli_agent_grants.go:445`, `plans\260518-1936-super-admin-api-parity\phase-04-credential-and-integration-surfaces.md:68` | Accept |
| RT-05 | High | Failure mode | Workstation link/unlink is WS-only and requires exact camelCase params. If treated like REST CRUD, it will compile but fail at runtime. | `D:\www\digitop\goclaw\internal\gateway\methods\workstations.go:52`, `D:\www\digitop\goclaw\internal\gateway\methods\workstations.go:53`, `D:\www\digitop\goclaw\internal\gateway\methods\workstations.go:276`, `D:\www\digitop\goclaw\internal\gateway\methods\workstations.go:277` | Accept |
| RT-06 | Medium | Failure mode | TTS synthesize returns raw audio, not JSON. CLI must require `--file` and must never send binary audio through `printer.Print`. | `D:\www\digitop\goclaw\internal\http\tts.go:22`, `D:\www\digitop\goclaw\internal\http\tts.go:23`, `D:\www\digitop\goclaw\internal\http\tts.go:58`, `plans\260518-1936-super-admin-api-parity\phase-05-media-tts-storage-channel-fillers.md:69` | Accept |
| RT-07 | Medium | Scope | README already claims "Full API coverage"; Phase 6 must explicitly rewrite this if any endpoint remains deferred, otherwise docs keep lying after the plan. | `README.md:7`, `plans\260518-1936-super-admin-api-parity\phase-06-docs-validation-and-release-readiness.md:26` | Accept |
| RT-08 | Medium | Complexity | The plan could overbuild a route scanner. A scanner is not necessary for this implementation; focused source inventory and command tests are enough. | `plans\260518-1936-super-admin-api-parity\phase-01-contract-audit-and-tdd-harness.md:55` | Reject: already constrained as manual report, not runtime CLI. |

## Accepted Plan Changes

- Phase 2 now explicitly requires custom-header tests and env fallback for gateway upgrade.
- Phase 2 now requires `packages updates apply-all` to fail non-zero on non-empty `failed[]` unless `--allow-partial` is passed.
- Phase 3 now pins WS link/unlink method names and parameter names.
- Phase 4 now blocks accidental `env:reveal` output in table mode without explicit reveal flags.
- Phase 5 now requires TTS synthesize to write raw audio only to a file.
- Phase 6 now requires README to revise the "Full API coverage" claim if any endpoint remains deferred.

## Rejected Findings

- RT-08 route scanner overbuild: rejected because Phase 1 already says "manual report step, not runtime CLI code."

## Unresolved Questions

- None.
