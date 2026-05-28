---
phase: 3
title: Remaining references
status: pending
priority: medium
effort_hours: 18
blockedBy: [1]
---

# Phase 3 — Remaining references

## Context Links
- Parent: [plan.md](plan.md)
- Phase 2 (template established): [phase-02-priority-references.md](phase-02-priority-references.md)
- Explorer: `plans/reports/explore-260417-1254-goclaw-command-inventory.md`

## Overview
Viết 10 reference files còn lại, theo cùng template đã chốt ở Phase 2. Độ ưu tiên thấp hơn vì use case ít gặp, nhưng cần đủ để đạt 100% coverage command surface.

## Scope — 11 files (M5 fix: `media` promoted to own cluster)

| File | Clusters | Target lines |
|------|----------|--------------|
| `agents-advanced.md` | agents.links, agents.ops, delegations (standalone, xác nhận từ `cmd/admin.go:141`) | 200 |
| `knowledge-memory.md` | kg, kg.dedup, memory | 200 |
| `teams-collaboration.md` | teams, teams.members, teams.events, teams.tasks, teams.workspace | 250 |
| `channels-messaging.md` | channels, channels.contacts, channels.pending, channels.writers, contacts, pending-messages | 200 |
| `data-movement.md` | export, import, storage | 150 |
| `providers-skills-tools.md` | providers, skills, tools (builtin list/get/update/tenant-config), packages | 250 |
| `automation-scheduling.md` | cron, heartbeat, heartbeat.checklist, heartbeat.targets, devices | 200 |
| `mcp-integration.md` | mcp, mcp.servers, mcp.grants, mcp.requests, mcp.reconnect | 200 |
| `admin-system.md` | tenants, system-config, activity, tts | 180 |
| `media.md` | media (admin_media.go) — **NEW** per red team M5 | 100 |
| `docs-api.md` | api-docs | 80 |

## Template (giữ như phase 2)

Mỗi file có sections:
- When to use
- Commands in scope
- Verified flags
- JSON output support per subcommand
- Destructive ops
- Common patterns (3-5 examples)
- Edge cases & gotchas
- Cross-refs

## Key streaming/TUI flags per cluster

- `teams-collaboration.md` — ⚠️ **toàn bộ teams.* là WebSocket**. Vẫn chạy được qua Bash (ws.Call sync) nhưng `teams.events` + task subscriptions là streaming — flag rõ
- `automation-scheduling.md` — ⚠️ `devices approve/reject` có interactive TUI fallback; dùng `--yes` để auto
- `mcp-integration.md` — ⚠️ `mcp reconnect` async; `mcp servers test` sync
- `channels-messaging.md` — ⚠️ một số subcommand chỉ trả success text (không JSON)
- `data-movement.md` — ⚠️ import/export lớn → timeout bash tool; doc cách split hoặc tăng timeout

## Related Code Files (to read per file)

- `agents-advanced.md` — `cmd/agents_links.go`, `agents_ops.go`
- `knowledge-memory.md` — `cmd/knowledge_graph.go`, `knowledge_graph_dedup.go`, `memory.go`, `memory_index.go`
- `teams-collaboration.md` — `cmd/teams*.go` (7 files)
- `channels-messaging.md` — `cmd/channels*.go`, `contacts.go`, `pending_messages.go`
- `data-movement.md` — `cmd/export_import.go`, `storage.go`
- `providers-skills-tools.md` — `cmd/providers*.go`, `skills*.go`, `tools.go`, `packages.go`
- `automation-scheduling.md` — `cmd/cron.go`, `heartbeat*.go`, `devices.go`
- `mcp-integration.md` — `cmd/mcp.go`, `mcp_reconnect.go`
- `admin-system.md` — `cmd/tenants.go`, `system_config.go`, `admin*.go`
- `docs-api.md` — `cmd/api_docs.go`, `api_keys.go`

## Implementation Steps

1. Có thể song song 10 files (độc lập nhau) — nếu làm tay tuần tự: mỗi file 15-30 phút
2. Cho file nào cluster lớn (teams-collaboration, providers-skills-tools) — spawn 1 fullstack-developer subagent riêng để đỡ bloat main context
3. Mỗi file done → cross-check verified flags match source
4. Cross-link: `automation-scheduling` link tới `auth-and-config` (device pairing nhắc tới auth); `data-movement` link tới `providers-skills-tools` (export skills)
5. Update SKILL.md navigation links khi mỗi file xong

## Todo List

- [ ] Write `agents-advanced.md`
- [ ] Write `knowledge-memory.md`
- [ ] Write `teams-collaboration.md` (lớn, có thể spawn subagent)
- [ ] Write `channels-messaging.md`
- [ ] Write `data-movement.md`
- [ ] Write `providers-skills-tools.md` (lớn, có thể spawn subagent)
- [ ] Write `automation-scheduling.md`
- [ ] Write `mcp-integration.md`
- [ ] Write `admin-system.md`
- [ ] Write `media.md` (NEW, cluster bị miss)
- [ ] Write `docs-api.md`
- [ ] Cross-link pass trên tất cả 16 references (phase 2 + 3)
- [ ] SKILL.md navigation list review
- [ ] Commit phase 3

## Success Criteria
- 10 reference files, mỗi file 80-250 lines, kebab-case naming
- Tổng 15 references tương đương ~3000 lines markdown
- Mọi flag verify từ source code (grep flag name trong `cmd/*.go`)
- SKILL.md navigation list cover hết 15 files

## Risk
- Scope lớn, dễ fatigue → ship từng batch 2-3 files/commit
- Teams cluster phức tạp nhất (7 source files); để cuối hoặc delegate subagent

## Security
- Không leak server URL, tenant, credential trong examples
- Placeholder `<xxx>` cho mọi giá trị user-specific

## Next
- Phase 4: install.sh + README + manual test
