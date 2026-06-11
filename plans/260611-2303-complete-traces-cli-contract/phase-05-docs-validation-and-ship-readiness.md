---
phase: 5
title: "Docs Validation and Ship Readiness"
status: complete
priority: P1
effort: "4h"
dependencies: [2, 3, 4]
---

# Phase 5: Docs Validation and Ship Readiness

## Overview

Update docs and run validation gates after all trace command gaps are implemented.

## Requirements

- Functional: README and changelog describe actual trace commands and filters.
- Functional: docs mention unsupported replay only as out of scope, not as command.
- Non-functional: compile/test/vet gates pass before handoff.

## Architecture

No architecture changes. This phase reconciles docs with final CLI behavior and performs a whole-plan consistency sweep.

## Related Code Files

- Modify: `/Users/duynguyen/.codex/worktrees/202b/goclaw-cli/README.md`
- Modify: `/Users/duynguyen/.codex/worktrees/202b/goclaw-cli/CHANGELOG.md`
- Modify: `/Users/duynguyen/.codex/worktrees/202b/goclaw-cli/docs/project-roadmap.md`
- Modify: `/Users/duynguyen/.codex/worktrees/202b/goclaw-cli/docs/codebase-summary.md`
- Read: `/Users/duynguyen/.codex/worktrees/202b/goclaw-cli/plans/260611-2303-complete-traces-cli-contract/*.md`

## Tests Before

- None; docs-only after implementation. Re-run all focused command tests first to protect behavior.

## Refactor

- Remove stale README examples:
  - `traces list --since`
  - `traces list --root-only`
  unless implementation intentionally preserves/deprecates them.
- Add examples:
  - `traces list --session-key ... --channel ... --offset ...`
  - `traces get <trace-id>`
  - `traces export <trace-id> --output trace.json.gz`
  - `traces follow --session-key ... --since ...`
  - `traces timeline <run-id> --session-key ...`

## Tests After

Run:
1. `go test -count=1 ./cmd -run 'TestTraces(List|Get|Follow|Timeline|Export)'`
2. `go test -count=1 ./cmd ./internal/client ./internal/output`
3. `go test -count=1 ./...`
4. `go vet ./...`
5. `go build ./...`

If local PATH lacks Go, record exact missing command and run with absolute Go path or stop before claiming done.

## Implementation Steps

1. Update README command table/details.
2. Add CHANGELOG bullet for trace contract alignment.
3. Update docs status so future agents do not trust old flat fixture plan.
4. Search all docs/plans for stale phrases:
   - `trace_id` as top-level field for `get`
   - `input_tokens`/`cost` trace table fields
   - `root-only`
   - `trace replay`
5. Run validation gates.
6. Capture final report in `plans/260611-2303-complete-traces-cli-contract/reports/final-validation.md`.

## Success Criteria

- [x] Docs match implemented flags and response behavior.
- [x] Unsupported replay is not presented as CLI-supported.
- [x] Full validation gates pass or exact environment blocker recorded.
- [x] Whole-plan consistency sweep finds no stale contract contradictions.

## Risk Assessment

Risk: broad doc scan touches unrelated historical plans.
Mitigation: update active docs and this plan; leave old reports immutable except adding supersession note if needed.

## Regression Gate

- `go test -count=1 ./...`
- `go vet ./...`
- `go build ./...`
