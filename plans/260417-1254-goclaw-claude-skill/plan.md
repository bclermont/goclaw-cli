---
title: GoClaw Claude Code Skill
status: in_progress
created: 2026-04-17
priority: high
blockedBy: []
blocks: []
source: skill
license: MIT
distribution: github-releases-tarball
implementation_progress: "Phases 1-3 done. Phase 4 95% done (install.sh + check-drift shellcheck-clean; TTY probe fixed; tarball/SHA256 rebuilt; root README updated; test matrix passed: modes 1/2/3, idempotent, piped-default-3, piped-mode-1-refused). Remaining: user runs 15-prompt smoke test in Claude Code + GitHub release publish."
---

# GoClaw Claude Code Skill — Implementation Plan

Build a Claude Code skill that lets Claude (agent) invoke goclaw CLI to interact with GoClaw Gateway server. Hero use case verified: `goclaw tools invoke exec --param command="<cmd>"` (server has `exec` builtin tool at `goclaw/internal/tools/shell.go:112-128`, params: `command` required, `working_dir` optional). Scope: 36 top-level command groups.

## Context Links

- Brainstorm synthesis: conversation turns `4/16 14:48` + `4/17 12:41`
- Researcher (skill authoring): `plans/reports/researcher-260417-1254-claude-skill-authoring.md`
- Explorer (command inventory): `plans/reports/explore-260417-1254-goclaw-command-inventory.md`
- Red team review: `plans/reports/code-reviewer-260417-1254-goclaw-skill-red-team.md`
- CLI source: `cmd/*.go` (~60 files, 36 groups verified)
- Server source: `/Volumes/GOON/www/nlb/goclaw/cmd/gateway_builtin_tools.go` (builtin tool registry), `/Volumes/GOON/www/nlb/goclaw/internal/tools/shell.go` (exec impl)
- Previous CLI plans (completed): `plans/260315-0007-goclaw-cli-implementation/`, `plans/260326-1350-cli-feature-parity-update/`

## Goals

1. Claude Code agent có thể gọi `goclaw <cmd>` autonomously qua Bash, parse JSON output, thao tác GoClaw server
2. Skill trigger tự động khi user nói về GoClaw (keyword matching qua description)
3. Cover toàn bộ 38 command groups qua progressive disclosure (14-15 reference files)
4. OSS-publishable: install script + README + LICENSE trong repo `goclaw-cli` (thư mục `claude-skill/`)
5. Safety: default prompt destructive, opt-in full-auto qua install flag

## Non-goals

- MCP server (bỏ per brainstorm YAGNI)
- Streaming commands (chat interactive, logs tail) — document limitation, không wrap
- Windows PowerShell script — macOS/Linux only cho v1
- Auto-generate references từ `--help` output — viết tay để control UX

## Tech Stack

- Markdown (SKILL.md + references/*)
- Bash (install.sh)
- Python3 (idempotent settings.json patcher, dùng venv `~/.claude/skills/.venv/bin/python3`)
- Hosting: trong repo `goclaw-cli` dưới `claude-skill/` để sync version với CLI

## Architecture Decisions

### D1 — Single repo (goclaw-cli/claude-skill/) vs separate repo
**Chọn:** đặt trong `goclaw-cli/claude-skill/`. Lý do: version sync tự động với CLI releases; user star 1 repo có cả 2; tránh divergence khi CLI đổi flags.

### D2 — Skill name
**Chọn:** `goclaw` (ngắn, match binary name). Không dùng `goclaw-cli` (redundant).

### D3 — Permission default (post-red-team revision)
**Problem cũ:** Red team chỉ ra `Bash(goclaw status)` không match `goclaw status --output json` (D6 force append `--output json`), và `Bash(goclaw list:*)` không match `goclaw agents list` (list không phải top-level).

**Chọn mới:** Install-time user interactive prompt, không hardcode pattern:
1. Install.sh mặc định hỏi "Choose permission mode: [1] Full auto (Bash(goclaw:*)) [2] Readonly verbs only [3] Manual (no patching)". Default = 3.
2. Mode 1 (`Bash(goclaw:*)`) — empirically verify wildcard matches `goclaw anything --output json` (Claude Code doc: `Bash(cmd:*)` = prefix match with any args). Nếu không match, fallback Mode 2.
3. Mode 2: enumerate readonly verbs per resource với trailing wildcard: `Bash(goclaw agents list:*)`, `Bash(goclaw agents get:*)`, `Bash(goclaw sessions list:*)`, ... (~20 rules). Generated từ inventory, không hardcode tay.
4. Mode 3: in JSON snippet, user tự copy.
5. `--mode 1|2|3` flag cho scripting (bypass prompt).

Loại bỏ `--full-auto` flag riêng; dùng `--mode 1` thay thế.

### D4 — Settings.json patching
**Chọn:** Python3 idempotent merge (per researcher report). Lý do: cross-platform, không cần jq, dùng venv sẵn có. Fallback: in hướng dẫn manual nếu venv không tồn tại.

### D5 — References structure
**Chọn:** 15 reference files (14 clusters từ explorer + 1 thêm `exec-workflow.md` cho hero use case). File nào > 300 lines tách tiếp.

### D6 — Convention trong SKILL.md
**Chọn:** Claude LUÔN append `--output json` khi chạy goclaw (trừ streaming). Quy tắc này viết ở SKILL.md và lặp ở mỗi reference.

## Phases Overview

| Phase | Title | Status | Effort (hrs) | Deliverable |
|-------|-------|--------|--------------|-------------|
| 1 | [Scaffold skill structure](phase-01-scaffold-skill-structure.md) | ✅ completed | 2 | `claude-skill/` dir + SKILL.md + MIT LICENSE + 16 stub references + install.sh + check-drift.sh skeletons + `.verified-commands.txt` (38 groups from binary) |
| 2 | [Priority references (hero path)](phase-02-priority-references.md) | ✅ completed | 10 | 5 refs: exec-workflow (verified schema from `shell.go:114-128`), auth-and-config, agents-core, chat-sessions, monitoring-ops |
| 3 | [Remaining references](phase-03-remaining-references.md) | ✅ completed | 18 | 11 refs including `media` cluster; total 16 refs, 1709 lines markdown |
| 4 | [Install script + README + release + test](phase-04-install-readme-test.md) | 🟡 95% done | 8 | install.sh + check-drift.sh shellcheck-clean; has_tty probe fixes pipe edge case; test matrix passed (3 modes + idempotent + piped default + Mode1 piped refusal); tarball + SHA256 rebuilt; root README linked. **Remaining: user runs 15-prompt smoke test in Claude Code + GitHub release publish.** |

**Total:** ~38 hrs ≈ 5 working days.

## Dependencies

- `goclaw` binary compiled + in PATH (user prereq)
- `~/.claude/skills/` directory exists (Claude Code installed)
- Python3 trong venv `~/.claude/skills/.venv/bin/python3` (nếu muốn auto-patch settings)

## Success Criteria

1. User nói "list agents trên goclaw" → Claude auto-invoke skill → chạy `goclaw agents list --output json` → summarize
2. User nói "run `ls -la` trên goclaw server" → Claude chạy `goclaw tools invoke exec --param command="ls -la"` → trả output + approval status
3. Destructive command (e.g. `goclaw agents delete xyz`) → Claude prompt user confirm trước khi thêm `--yes`
4. `./install.sh` chạy idempotent: chạy lần 2 không duplicate permissions
5. README chứa install one-liner + 3 example prompts
6. Manual test: skill trigger đúng cho ≥ 5 intent khác nhau, ≥ 90% commands Claude phát ra đúng flag

## Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| `allowed-tools` trong frontmatter không enforce (issue #14956) | High | Medium | Dùng settings.json permissions làm source of truth |
| Claude hallucinate flag name không có trong reference | Medium | Low | `check-drift.sh` grep every flag trong references vs `cmd/*.go` |
| settings.json malformed → install.sh phá config user | Low | High | `if [[ -f ]]; then cp ... \|\| exit 1; fi` (not `&&` chain); abort nếu invalid JSON |
| goclaw CLI đổi flag → reference outdated | Medium | Medium | `check-drift.sh` trong CI; version compat note ở SKILL.md frontmatter |
| User chọn Mode 1 trên prod → xóa data | Low | Critical | Mode 3 là default; Mode 1 có interactive confirm qua `/dev/tty` |
| Description keywords quá hẹp → skill không trigger | Medium | Low | 15-prompt smoke test (5 positive + 5 negative + 5 destructive) |
| Tarball download corruption / MITM | Low | High | SHA256 in install.sh; verify before extract |
| exec tool param schema đổi ở server release | Low | Medium | Pin supported goclaw versions trong SKILL.md; `check-drift.sh` optional server check |
| `curl|bash` install có stdin piped → interactive prompts fail | Medium | Medium | `read ... < /dev/tty`; refuse Mode 1 nếu `! -t 0` |

## Security Considerations

- **Token handling:** Skill KHÔNG lưu token — dùng credential store của goclaw CLI (đã có)
- **settings.json backup:** Trước patch, copy sang `.bak` với timestamp
- **Deny list default:** Include `Bash(goclaw tenants delete:*)`, `Bash(goclaw * --yes *)` trong deny cho non-full-auto mode
- **Public OSS hygiene:** KHÔNG hardcode server URL, tenant ID, token trong skill files

## Next Steps

1. Phase 1: scaffold (độc lập, có thể bắt đầu ngay)
2. Phase 2 & 3 có thể song song sau phase 1
3. Phase 4 chốt sau khi 2+3 xong
