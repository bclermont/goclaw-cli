---
phase: 6
title: "Docs Validation and Release Readiness"
status: completed
effort: "2h"
---

# Phase 6: Docs Validation and Release Readiness

## Overview

Priority: P1
Current status: completed

Validate the full implementation, update project docs, and prepare a focused commit/PR handoff. This phase does not ship unless tests and docs match actual implemented command behavior.

## Context Links

- [Plan Overview](./plan.md)
- [Red Team Review](./reports/red-team-review.md)
- `D:\www\nextlevelbuilder\goclaw-cli\README.md`
- `D:\www\nextlevelbuilder\goclaw-cli\docs\project-roadmap.md`
- `D:\www\nextlevelbuilder\goclaw-cli\docs\codebase-summary.md`
- `D:\www\nextlevelbuilder\goclaw-cli\docs\system-architecture.md`

## Key Insights

- Existing docs overclaim full coverage in places. Docs must state exact command support.
- Route parity should be validated with command help and source grep, not only README.
- Broad tests are feasible here because CLI repo is small compared with server repo.

## Requirements

- Update README command inventory and examples.
- Update roadmap/changelog-equivalent docs if present.
- Run full compile, vet, and tests.
- Produce final route coverage note in plan reports.
- Revise README "Full API coverage" claim if any server endpoint remains deferred after this slice.

## Architecture

```
implementation -> focused tests -> full tests -> docs sync -> final report
```

## Related Code Files

| File | Action |
|---|---|
| `README.md` | modify |
| `docs/project-roadmap.md` | modify |
| `docs/codebase-summary.md` | modify |
| `docs/system-architecture.md` | modify if command architecture changes |
| `plans/260518-1936-super-admin-api-parity/reports/final-validation.md` | create |

## Implementation Steps

1. Run focused tests from Phases 2-5.
2. Run `go test ./...`, `go vet ./...`, `go build ./...`.
3. Run help smoke commands for new command groups.
4. Update docs with only implemented commands.
5. Write final validation report with passed/failed commands and remaining deferred endpoints.
6. Prepare scoped git diff; do not stage unrelated `AGENTS.md` unless explicitly requested.
7. If any endpoint remains deferred, document it explicitly instead of keeping blanket "Full API coverage" wording.

## Tests Before

- No docs updates until commands compile and focused tests pass.

## Refactor

- If any command file exceeds 200 lines, split before final validation.
- Remove duplicated multipart/header code if repeated across three files.

## Tests After

- `go test ./...`
- `go vet ./...`
- `go build ./...`
- `go run . gateway upgrade --help`
- `go run . workstations --help`
- `go run . packages updates --help`
- `go run . webhooks --help`

## Todo List

- [x] Run focused tests.
- [x] Run full validation.
- [x] Update README and docs.
- [x] Write final validation report.
- [x] Review git diff for secrets and unrelated changes.
- [x] Replace or qualify README coverage claims based on actual final route report.

## Success Criteria

- [x] Full validation passes or failures are documented with exact commands.
- [x] Docs no longer overclaim unimplemented behavior.
- [x] Final diff is scoped to super-admin API parity work.
- [x] Final validation report lists deferred endpoints with reason.

## Risk Assessment

- Risk: server route changes again before implementation. Mitigation: rerun route scan at start of cook.
- Risk: docs churn unrelated to code. Mitigation: update only command inventory and current status.

## Security Considerations

- Run secret scan on staged diff before commit.
- Avoid storing fixture tokens that look real.

## Regression Gate

```powershell
go test ./...
go vet ./...
go build ./...
git diff --check
```

## Next Steps

- After implementation, use `ck:code-review`, then scoped git commit/push flow if requested.
