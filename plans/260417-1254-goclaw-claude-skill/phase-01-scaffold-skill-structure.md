---
phase: 1
title: Scaffold skill structure
status: pending
priority: high
effort_hours: 2
---

# Phase 1 — Scaffold skill structure

## Context Links
- Parent: [plan.md](plan.md)
- Researcher: `plans/reports/researcher-260417-1254-claude-skill-authoring.md`
- Explorer: `plans/reports/explore-260417-1254-goclaw-command-inventory.md`

## Overview
Bootstrap thư mục `claude-skill/` trong repo `goclaw-cli` với cấu trúc chuẩn skill, SKILL.md skeleton, thư mục `references/` rỗng (stub cho 15 files), và `install.sh` chưa active. Không implement logic patching settings.json ở phase này.

## Key Insights
- Skill ở TRONG repo `goclaw-cli` (không tách repo) → version sync với CLI
- SKILL.md < 100 lines (thin index + navigation)
- References là progressive disclosure: Claude dùng Read tool khi cần
- Frontmatter description LÀ keyword matcher → tune kỹ sau

## Requirements

### Functional
- Tạo thư mục `claude-skill/` ở root repo
- `SKILL.md` với frontmatter đầy đủ (name, description, when_to_use, allowed-tools)
- Thư mục `references/` với 16 file .md rỗng (stub) — thêm `media.md` mà explorer miss
- `install.sh` skeleton (chưa có logic, chỉ echo placeholder)
- `README.md` skeleton
- `LICENSE` — **MIT** (user confirmed); viết trực tiếp từ template SPDX MIT
- Re-verify top-level command count: chạy `goclaw --help` sau khi build binary, confirm 36 groups (không phải 38 như explorer report). Update inventory nếu lệch.

### Non-functional
- File encoding UTF-8
- Unix line endings
- Kebab-case cho tên reference files

## Architecture

```
goclaw-cli/
├── claude-skill/
│   ├── SKILL.md
│   ├── README.md
│   ├── LICENSE                 # MIT
│   ├── install.sh
│   ├── check-drift.sh          # CI flag-validation script (stub, impl Phase 4)
│   └── references/             # 16 files
│       ├── exec-workflow.md         # hero
│       ├── auth-and-config.md
│       ├── agents-core.md
│       ├── agents-advanced.md
│       ├── chat-sessions.md
│       ├── monitoring-ops.md
│       ├── knowledge-memory.md
│       ├── teams-collaboration.md
│       ├── channels-messaging.md
│       ├── data-movement.md
│       ├── providers-skills-tools.md
│       ├── automation-scheduling.md
│       ├── mcp-integration.md
│       ├── admin-system.md
│       ├── media.md                 # NEW: was missing from explorer
│       └── docs-api.md
```

## Related Code Files

### To create
- `claude-skill/SKILL.md`
- `claude-skill/README.md`
- `claude-skill/LICENSE`
- `claude-skill/install.sh`
- `claude-skill/references/*.md` (15 files, stub)

### To read for context
- `goclaw-cli/README.md` — tone/voice match
- `goclaw-cli/LICENSE` — copy same license
- `cmd/root.go` — verify global flags (--output, --yes, --profile, --tenant-id)

## Implementation Steps

1. Check license file exists at repo root. If yes, copy to `claude-skill/LICENSE`. If no, note for Phase 4.
2. Create directory `claude-skill/` and `claude-skill/references/`.
3. Write `SKILL.md` with:
   - Frontmatter: `name: goclaw`, `description` (front-load keywords: "goclaw", "gateway server", "AI agent", "exec remote command", "GoClaw Gateway"), `when_to_use`, `allowed-tools: Bash(goclaw:*)`
   - Body: overview (3 paragraphs max), convention rules (always `--output json`, always read credential store from `~/.goclaw/`, NEVER hardcode token), navigation list linking to each reference file
   - Length target: 80-100 lines
4. Write `README.md` skeleton: title, install one-liner placeholder, 3-example prompts placeholder, requirements (goclaw in PATH, Claude Code installed), license link.
5. Write `install.sh` skeleton: shebang `#!/usr/bin/env bash`, `set -euo pipefail`, echo-only placeholder (no patching yet). Flag Phase 4 for real implementation.
6. Create 15 stub reference files. Each stub contains: 1-line title (H1), "TODO: populate in Phase 2/3" comment. Explicit phase assignment in TODO.
7. Verify: `tree claude-skill/` produces expected structure; `bash -n install.sh` passes syntax check.

## Todo List

- [ ] Build `goclaw` binary; run `goclaw --help` → re-verify 36 top-level groups (confirm `delegations`, `media` standalone; no phantom entries)
- [ ] Scaffold directory `claude-skill/` + `references/`
- [ ] Write `LICENSE` (MIT, copyright NextLevelBuilder 2026)
- [ ] Write `SKILL.md` with keyword-tuned frontmatter (< 100 lines)
- [ ] Write `README.md` skeleton (install via tarball+SHA256 placeholder)
- [ ] Write `install.sh` shebang + `set -euo pipefail` + arg parser stub
- [ ] Write `check-drift.sh` stub (Phase 4 fills logic)
- [ ] Create 16 stub reference files (15 + `media.md`) với phase assignment comment
- [ ] Run `bash -n install.sh check-drift.sh` syntax check
- [ ] Commit phase 1 artifact

## Success Criteria
- `tree claude-skill/` shows 21 files (SKILL.md, README.md, LICENSE, install.sh, check-drift.sh, 16 references/)
- `SKILL.md` < 100 lines + valid YAML frontmatter
- `bash -n claude-skill/install.sh check-drift.sh` exits 0
- Mỗi stub reference có 1 TODO line chỉ phase sẽ được fill
- `goclaw --help` output committed as `claude-skill/.verified-commands.txt` cho reference traceability

## Risk
- Nếu `goclaw --help` list không match 36 → explorer report sai → cluster mapping cần re-balance (adjust Phase 2/3 trước khi start)
- MIT LICENSE yêu cầu copyright holder rõ — dùng "NextLevelBuilder" (từ repo owner `nextlevelbuilder`)

## Security
- Không hardcode endpoint/tenant trong stub
- install.sh chưa touch `~/.claude/settings.json` ở phase này

## Next
- Phase 2 viết 5 reference priority (hero use case: exec)
