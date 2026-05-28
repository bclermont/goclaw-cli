---
phase: 2
title: Priority references (hero path)
status: pending
priority: high
effort_hours: 10
blockedBy: [1]
---

# Phase 2 — Priority references (hero path)

## Context Links
- Parent: [plan.md](plan.md)
- Phase 1: [phase-01-scaffold-skill-structure.md](phase-01-scaffold-skill-structure.md)
- Explorer: `plans/reports/explore-260417-1254-goclaw-command-inventory.md`

## Overview
Viết 5 reference files phục vụ use case chính (exec + nhóm thường dùng nhất). Sau phase này skill đã đủ dùng cho 70% tình huống thực tế.

## Priority order (ship nhanh use case "exec trên server" trước)

1. `exec-workflow.md` — HERO. `tools invoke exec` + approvals flow + examples.
2. `auth-and-config.md` — login, whoami, profile, tenant switch (mọi skill flow đều cần context auth)
3. `chat-sessions.md` — chat single-shot + sessions list/preview/delete
4. `agents-core.md` — agents list/get/create/delete + instances
5. `monitoring-ops.md` — status, health, logs (non-streaming), traces, usage

## Content template (áp dụng cho mọi reference file)

```markdown
# <Cluster Name>

## When to use
<1 sentence — khi nào Claude load reference này>

## Commands in scope
<bullet list: group → subcommands>

## Verified flags (từ code source)
<table: flag | type | default | purpose>

## JSON output
<Y/N per subcommand, cite printer.Print usage>

## Destructive ops
<bullet list — which subcommands require --yes, which require user confirm>

## Common patterns
<3-5 worked examples, copy-pasteable bash>

## Edge cases & gotchas
<bullet list>

## Cross-refs
<links tới references khác nếu có overlap>
```

## Specific content per file

### exec-workflow.md (target 200-250 lines) — VERIFIED SCHEMA

**Server impl:** `/Volumes/GOON/www/nlb/goclaw/internal/tools/shell.go:112-128` + builtin registry `gateway_builtin_tools.go:24`.

**Tool name:** `exec`
**Parameters schema:**
```json
{
  "command": "string (REQUIRED) — shell command to execute",
  "working_dir": "string (optional) — defaults to workspace root"
}
```
**Response:** stdout/stderr từ `*Result` struct (check struct shape in `internal/tools/executor.go` trước khi viết example).
**Approval gating:** package install commands + shell deny patterns → creates approval request via `ExecApprovalManager` (`exec_approval.go:91` — field `Command`). User approves via `goclaw approvals approve <id>`.

**Content outline:**
- Cover: `goclaw tools invoke exec --param command="..."` + `--param working_dir=/tmp`
- Cover: `goclaw tools invoke <tool-name>` generic pattern (link vào file `providers-skills-tools.md` cho tool list)
- Cover: `goclaw approvals list/approve/deny` — cần chạy qua ws.Call (xem `cmd/admin.go:26-80`)
- **Canonical home của approvals** = file này (M4 resolution); `chat-sessions.md` cross-link tới đây
- Example 1: one-shot command with auto-approve (innocuous command, không trigger approval)
- Example 2: command hits approval gate → user approves → retry
- Example 3: parse stdout/stderr/exit_code fields JSON
- Gotcha: `approvals watch` streaming — flag "not Bash-friendly, use polling `approvals list --output json`"
- Gotcha: NUL byte rejection, shell deny patterns, Unicode normalization (từ shell.go:138-144)

### auth-and-config.md (target 150-200 lines)
- Cover: `auth login/logout/whoami/use-context/list-contexts/pair`
- Cover: `config get/set/permissions`
- Cover: `credentials list/create/delete/rotate`
- Cover: `api-keys list/create/reveal/revoke/extend`
- Cover: global flags `--profile`, `--tenant-id`, `--server`, `--token`
- Example: switch profile mid-session
- Gotcha: `auth pair` = device pairing, long-running polling — document "not Bash-friendly; skill REFUSE to run, tell user chạy shell riêng"
- Gotcha: **token expiry handling (U2 fix)** — nếu `goclaw` exit code indicate 401, Claude phải suggest `goclaw auth login` (không try-retry loop). Document exit codes nếu có.

### chat-sessions.md (target 200 lines)
- Cover: `chat -m "message" --no-stream --output json` (single-shot ONLY)
- Cover: `chat abort <session>` (destructive, no --yes flag) — **include trong destructive section** (N3 fix)
- Cover: `chat inject`, `chat status`
- Cover: `sessions list/preview/delete/reset/label`
- Cover: NDJSON parsing từ `chat --output json`
- Cross-link: approvals → xem `exec-workflow.md` (không duplicate)
- Gotcha: `chat` interactive mode = TUI, không dùng được qua Bash; chỉ dùng single-shot
- Example: chat → get session ID → preview history

### agents-core.md (target 250 lines — cluster lớn)
- Cover: `agents list/get/create/update/delete`
- Cover: `agents files list/get/create/delete`
- Cover: `agents instances list/get/create/delete/trigger/reset`
- Cover: `agents wake`
- Example: full lifecycle — create agent → upload file → instance → wake
- Destructive: delete agent + instances (both need --yes)

### monitoring-ops.md (target 150 lines)
- Cover: `status`, `health`, `version`
- Cover: `traces list/get`
- Cover: `usage summary/breakdown/trends/export`
- Cover: `logs tail` — **FLAG: streaming, NOT Bash-friendly**
- Example: diagnose server health flow

## Related Code Files (to read)

- `cmd/tools.go` — verify exec invoke schema
- `cmd/admin.go:11-82` — approvals flow (ws.Call signatures)
- `cmd/auth.go` — auth flow + pairing
- `cmd/chat.go` — single-shot vs interactive modes
- `cmd/sessions.go` — session CRUD
- `cmd/agents.go` + `agents_*.go` — agent lifecycle
- `cmd/status.go`, `health.go`, `logs.go`, `traces.go`, `usage.go` — monitoring commands
- `cmd/root.go` — global flags

## Implementation Steps

1. Đọc source files listed above, extract flag definitions + subcommand schemas
2. Viết `exec-workflow.md` trước (hero use case)
3. Test: trong Claude Code, chạy thử "run `pwd` trên goclaw server" → verify Claude dùng `tools invoke exec --param command="pwd"` đúng
4. Viết 4 file còn lại theo template
5. Mỗi file: check mọi flag listed phải xuất hiện trong `cmd/*.go` tương ứng (không invent)
6. Cross-link giữa các file (vd exec-workflow link tới auth-and-config cho context "token phải set trước")

## Todo List

- [ ] Read source files để extract verified flags
- [ ] Write `exec-workflow.md` (hero)
- [ ] Manual test exec-workflow trong Claude Code
- [ ] Write `auth-and-config.md`
- [ ] Write `chat-sessions.md`
- [ ] Write `agents-core.md`
- [ ] Write `monitoring-ops.md`
- [ ] Cross-link review
- [ ] Commit phase 2

## Success Criteria
- 5 reference files mỗi file 150-250 lines
- Mọi flag listed verify được từ `cmd/*.go`
- Claude Code test: ≥ 3 intents về exec chạy đúng lệnh mà không hallucinate flag
- Destructive ops mỗi file có section "user confirm required"

## Risk
- ~~Schema response của `tools invoke exec` chưa known~~ **RESOLVED:** verified trong `shell.go:114-128`, `exec_approval.go:91`
- Response body (`*Result` struct) chưa đọc — verify `internal/tools/executor.go` Result shape trước khi viết Example 3
- Verify flag match tốn thời gian — chấp nhận vì chất lượng reference quyết định UX
- Approvals đã đưa vào exec-workflow nhưng `chat-sessions.md` dễ duplicate — phải strict cross-link, không copy content

## Security
- Mọi example dùng placeholder `<agent-id>`, `<tenant>`, `<user-id>` — không hardcode data thật

## Next
- Phase 3: viết 10 reference còn lại
