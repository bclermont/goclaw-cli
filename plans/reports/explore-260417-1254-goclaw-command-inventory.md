# GoClaw CLI Command Surface Inventory

**Generated:** 2026-04-17 | **Thoroughness:** Medium | **Scope:** 38 top-level command groups + 60 files

---

## 1. TOP-LEVEL COMMAND MAP

Complete mapping of all 38 command groups registered in `rootCmd.AddCommand` (root.go + resource files).

| Group | Subcommands | Purpose | Key Flags | Destructive | JSON | Streaming | Notes |
|-------|-------------|---------|-----------|-------------|------|-----------|-------|
| **activity** | list | View audit log | --limit | N | Y | N | HTTP GET, table/json output |
| **agents** | list, get, create, update, delete, [files, instances, links, ops, wake] | Manage agents & instances | --name, --provider, --model, --type, --yes | Y (delete) | Y | N | HTTP CRUD; delete uses tui.Confirm |
| **agents.files** | list, get, create, delete | Agent context files | --data (JSON/@filepath) | Y | N | WS only | WebSocket calls, TUI-driven |
| **agents.instances** | list, get, create, delete, [trigger, reset] | Per-user agent instances | --user-id, --yes | Y | N | WS only | WebSocket calls |
| **agents.links** | list, create, delete | Delegation links | --agent, --target, --yes | Y | N | HTTP | tui.Confirm for delete |
| **agents.ops** | share, unshare, regenerate, resummon, wait | Agent operations | --user, --yes | Y | N | HTTP/WS mix | share/unshare=HTTP DELETE; wait=WS.Subscribe |
| **agents.wake** | wake | Wake sleeping agent | (none) | N | N | HTTP POST | Single action, always succeeds |
| **api-docs** | open, spec | API documentation | (none) | N | N | HTTP GET | open=browser launch, spec=JSON fetch |
| **api-keys** | list, create, reveal, revoke, extend | API key management | --name, --scopes, --expires-in, --yes | Y (revoke) | Y | HTTP | Scoped access, masked display |
| **approvals** | list, approve, deny, watch | Execution approvals | --reason, --yes | Y (deny) | N | WS + tui | watch = ws.Subscribe; approve/deny = ws.Call |
| **auth** | login, logout, whoami, use-context, list-contexts, pair | Authentication | --profile, --pair, --token | Y (logout) | N | HTTP/WS | pair flow uses device pairing code polling |
| **channels** | [instances, contacts, pending, writers] | Messaging channels | --type, --name, --agent, --yes | Y (delete) | Y | HTTP | Channel type filtering, table output |
| **channels.contacts** | list, get, create, delete | Contact management | (varies) | Y | Y | HTTP | tui.Confirm for delete |
| **channels.pending** | list | View pending messages | (none) | N | Y | HTTP GET | Filtered by channelID |
| **channels.writers** | list, add, remove | Group writer mgmt | --agent, --user, --yes | Y (remove) | Y | HTTP | tui.Confirm for remove |
| **chat** | send (primary), inject, status, abort | Interactive & one-shot chat | -m, --session, --no-stream | N | Y (--output json) | WS stream | Primary: ws.Stream (NDJSON); abort=destructive but no --yes |
| **config** | [get, set, permissions] | Server configuration | --key, --value | Y (set) | Y | WS | Uses tui.Confirm for write operations |
| **config.permissions** | list, update, grant, revoke | Config permissions | --action, --yes | Y (revoke) | N | WS | WebSocket-only, permission model |
| **contacts** | list, get, create, update, delete, verify | Manage contacts | --name, --email, --phone, --yes | Y (delete) | Y | HTTP CRUD | tui.Confirm for delete |
| **credentials** | list, create, delete, rotate | CLI credentials store | --profile, --yes | Y (delete) | Y | HTTP | Per-profile credential storage |
| **cron** | list, get, create, update, delete, trigger, history | Scheduled jobs | --agent, --schedule, --message, --yes | Y (delete) | Y | WS primary | ws.Call for CRUD; fallback to HTTP |
| **delegations** | (implied in agents.links) | Delegation links | — | — | — | — | See agents.links |
| **devices** | list, delete, approve, reject | Paired device mgmt | --device-id, --yes | Y (delete) | Y | WS | Device pairing & approval flow |
| **export** | [agent-preview, agent, team-preview, team, skills-preview, skills, mcp-preview, mcp] | Resource export | --agent, --team, --skills, --output, --yes | N | Y | HTTP | Export as JSON/YAML for import |
| **health** | (standalone cmd) | Server health check | (none) | N | N | HTTP GET | Simple status check, no args |
| **heartbeat** | [get, set, checklist, targets] | Heartbeat configuration | --agent, --interval, --yes | Y (set) | Y | WS + HTTP | Multi-part: config (WS) + targets (HTTP) |
| **heartbeat.checklist** | get, set | Checklist mgmt | --target-id, --enabled, --yes | Y (set) | Y | WS | WebSocket-driven |
| **heartbeat.targets** | list, add, remove | Checklist targets | --target-id, --name, --yes | Y (remove) | Y | HTTP | tui.Confirm for remove |
| **import** | [agent, team, skills, mcp] | Resource import | --file, --yes | Y | Y | HTTP POST | Reverse of export; requires --yes flag |
| **kg** (knowledge-graph) | [entities, traverse, graph, stats, dedup] | Knowledge graph ops | --data, --from, --yes | Y (delete entity) | Y | HTTP CRUD | Alias: kg; entities={list,get,create,delete} |
| **kg.dedup** | scan, merge-candidates, execute-merge | Entity deduplication | --agent, --yes | Y (merge) | Y | HTTP | Scan → review → execute workflow |
| **logs** | tail | Stream server logs | --agent, --level | N | Y (json NDJSON) | WS subscribe | ws.Subscribe("*") on signal.Notify; interactive |
| **mcp** | [servers, grants, requests, reconnect] | MCP server mgmt | --name, --transport, --command, --agent, --server, --yes | Y (delete server) | Y | HTTP CRUD | 3 subsystems: servers, grants, requests; reconnect=POST |
| **memory** | [list, get, store, index, delete, clear] | Agent memory docs | --user, --content, --data, --yes | Y (delete/clear) | Y | HTTP CRUD | Per-agent/user memory store |
| **packages** | list | Runtime packages | (none) | N | Y | HTTP GET | Read-only list of installed packages |
| **pending-messages** | list, create, send, delete | Pending message queue | (varies) | Y (delete) | Y | HTTP | tui.Confirm for delete |
| **providers** | [list, get, create, update, delete, verify-embedding] | LLM provider config | --name, --display-name, --api-key, --yes | Y (delete) | Y | HTTP CRUD | Verify endpoint tests embeddings |
| **sessions** | list, preview, delete, reset, label | Chat session mgmt | --agent, --user, --label, --yes | Y (delete/reset) | Y | HTTP CRUD | tui.Confirm for destructive ops |
| **skills** | [list, get, upload, download, publish, unpublish, delete, tenant-config, versions, grant, revoke] | Skill management | --search, --slug, --visibility, --yes | Y (delete/unpublish) | Y | HTTP multipart | upload=multipart form; visibility gating |
| **status** | (standalone cmd) | Server status | (none) | N | Y | WS then table | ws.Call("status") → table format fallback |
| **storage** | list, get, put, delete | Workspace files | --path, --content, --data, --yes | Y (delete) | Y | HTTP | Path-based file browsing |
| **system-config** | list, get, set, delete | Per-tenant KV config | --key, --value, --yes | Y (delete) | Y | HTTP | Tenant-scoped configuration |
| **teams** | [list, get, create, delete, members, events, tasks, workspace] | Agent team mgmt | --name, --agents, --yes | Y (delete) | N | WS only | All team ops are WebSocket; no HTTP fallback |
| **teams.members** | list, add, remove, reassign | Team membership | --user, --team-id, --yes | Y (remove) | N | WS | WebSocket-driven |
| **teams.tasks** | [list, get, create, delete, approve, reject, reassign] | Team task mgmt | --task-id, --title, --yes | Y | N | WS | ws.Call & ws.Subscribe for events |
| **teams.workspace** | list, get, put, delete | Team workspace files | --path, --content, --yes | Y (delete) | N | WS | WebSocket-based file ops |
| **tenants** | [list, get, create, update, delete, users] | Multi-tenant admin | --name, --users, --yes | Y (delete) | Y | HTTP CRUD | Admin-only; /v1/tenants endpoints |
| **tools** | [builtin, custom] | Built-in & custom tools | --search, --yes | Y (delete custom) | Y | HTTP | builtin=read-only list; custom=CRUD |
| **traces** | list, get | LLM trace viewing | --agent, --limit | N | Y | HTTP GET | OpenAI-compatible trace format |
| **tts** | status | Text-to-speech | (none) | N | Y | WS call | ws.Call("tts.status") |
| **usage** | [summary, breakdown, trends, export] | Usage analytics | --period, --agent, --format, --export | N | Y | HTTP GET | Metrics: views, completions, cost |
| **version** | (standalone cmd) | Version info | (none) | N | N | stdout | Built-in version/commit/build info |

---

## 2. PROPOSED LOGICAL CLUSTERING (12-18 reference files)

Grouping by cognitive locality and CLI skill reference organization:

| Cluster | File | Commands | Rationale |
|---------|------|----------|-----------|
| **auth-and-config** | `auth-and-config.md` | auth, api-keys, config, config.permissions, credentials | Auth/session + credential storage + server config |
| **agents-core** | `agents-core.md` | agents, agents.files, agents.instances, agents.wake | Core agent lifecycle + context + instances |
| **agents-advanced** | `agents-advanced.md` | agents.links, agents.ops, delegations | Advanced agent operations: delegation, sharing, regenerate |
| **chat-sessions** | `chat-sessions.md` | chat, sessions, approvals | Interactive chat + session mgmt + execution approval |
| **knowledge-memory** | `knowledge-memory.md` | kg, kg.dedup, memory, memory.index | Knowledge graph + entity dedup + agent memory |
| **teams-collaboration** | `teams-collaboration.md` | teams, teams.members, teams.events, teams.tasks, teams.workspace | Team-based operations + task mgmt |
| **channels-messaging** | `channels-messaging.md` | channels, channels.instances, channels.contacts, channels.pending, channels.writers | Messaging channel delivery |
| **data-movement** | `data-movement.md` | export, import, storage | Workspace file + resource export/import |
| **providers-skills** | `providers-skills.md` | providers, skills, tools, packages | LLM providers + skills + tool registry |
| **automation-scheduling** | `automation-scheduling.md` | cron, heartbeat, heartbeat.checklist, devices | Scheduled & event-driven automation |
| **mcp-integration** | `mcp-integration.md` | mcp, mcp.servers, mcp.grants, mcp.requests | MCP protocol integration + access mgmt |
| **monitoring-ops** | `monitoring-ops.md` | health, status, logs, traces, usage, version | Observability + health + analytics |
| **admin-system** | `admin-system.md` | tenants, system-config, activity | Admin operations + multi-tenancy |
| **docs-api** | `docs-api.md` | api-docs | API documentation access |

---

## 3. STREAMING COMMANDS (WebSocket-dependent, unsuitable for one-shot Bash)

Commands using `ws.Subscribe()` or long-lived `ws.Stream()` that **require persistent connection**:

| Command | Pattern | Issue | Workaround |
|---------|---------|-------|-----------|
| **chat** (interactive) | ws.Stream(chat.send) | Open-ended streaming, user-driven | Use `--no-stream` or `-m "single message"` for Bash |
| **logs tail** | ws.Subscribe("*") on signal | Real-time tail, Ctrl+C to exit | Bash integration requires monitor wrapper |
| **approvals watch** | ws.Subscribe() implied | Pending approvals stream (if implemented) | Use approvals.list + poll instead |
| **teams.events** | ws.Subscribe on team_id | Live event stream per team | No Bash-friendly alternative |
| **teams.tasks** (streaming) | ws.Subscribe tasks.* | Live task updates | List + poll for updates |
| **auth pair** | Poll-based device pairing (60× 2s sleep) | Not streaming but long-running | Bash skill should implement full flow |
| **agents.ops wait** | ws.Subscribe(agent_id) | Wait for agent idle state | Use polling with status check |
| **devices approve** | Interactive approval flow (WS) | Requires human decision | WS skill wrapper needed |

**Bash Skill Recommendation:** Wrap streaming commands in Monitor tool, not direct Bash invocation.

---

## 4. DESTRUCTIVE COMMANDS (require `--yes` flag or confirmation)

All commands with `tui.Confirm()` or that perform DELETE/destructive operations:

### Requiring `--yes` flag (automation-safe when set):
- **agents delete** — `tui.Confirm(fmt.Sprintf("Delete agent %s?", args[0]), cfg.Yes)`
- **api-keys revoke** — Permanent key revocation
- **channels delete** — Channel instance removal
- **channels.contacts delete** — Contact deletion
- **channels.pending send** — Sends pending message (irreversible)
- **channels.writers remove** — Writer access revocation
- **config set** — Modifies server config
- **config.permissions revoke** — Permission revocation
- **contacts delete** — Contact removal
- **credentials delete** — Credential store deletion
- **credentials rotate** — Key rotation (old key invalid)
- **cron delete** — Job removal
- **devices delete** — Unpair device
- **import **** — Bulk resource import (overwrites)
- **kg entities delete** — Knowledge graph entity removal
- **kg dedup execute-merge** — Permanent entity merge
- **memory delete**, **memory clear** — Memory document removal
- **pending-messages delete** — Message queue deletion
- **providers delete** — Provider config removal
- **sessions delete**, **sessions reset** — Session destruction
- **skills unpublish**, **skills delete** — Skill removal from catalog
- **storage delete** — Workspace file deletion
- **system-config delete** — KV config removal
- **teams delete** — Team removal
- **teams.members remove** — Member removal
- **teams.tasks delete** — Task deletion
- **teams.workspace delete** — Workspace file removal
- **tenants delete** — Tenant removal (admin)
- **tools delete** (custom) — Custom tool removal

### High-Risk (require skill prompt guidance):
- **export *** (with --yes) — Bulk data extraction
- **import *** (with --yes) — Bulk data overwrite
- **memory clear** — Entire agent memory wipe

**Skill Implementation:** Always prompt before operations with `--yes` flag unless user explicitly approves.

---

## 5. JSON OUTPUT SUPPORT ANALYSIS

Commands that DO vs DON'T cleanly emit JSON for programmatic parsing:

### JSON-Friendly (✅ full support):
- **agents list, get** → printer.Print(unmarshalList/Map(data))
- **api-keys list** → printer.Print(unmarshalList(data))
- **approvals list** → printer.Print(unmarshalList(data))
- **channels instances list** → printer.Print(unmarshalList(data))
- **chat (--output json)** → NDJSON event stream
- **config get** → printer.Print(unmarshalMap(data))
- **contacts list, get** → printer.Print(unmarshalList/Map(data))
- **cron list, get** → printer.Print(unmarshalList/Map(data))
- **devices list** → printer.Print(unmarshalList(data))
- **kg entities list, get** → printer.Print(unmarshalList/Map(data))
- **kg graph, stats** → printer.Print(unmarshalMap(data))
- **mcp servers list, get** → printer.Print(unmarshalList/Map(data))
- **memory list, get** → printer.Print(unmarshalList/Map(data))
- **providers list** → printer.Print(unmarshalList(data))
- **sessions list, preview** → printer.Print(unmarshalList/Map(data))
- **skills list, get** → printer.Print(unmarshalList/Map(data))
- **status** → ws.Call("status") → table fallback OR json via --output
- **storage list** → printer.Print(unmarshalList(data))
- **system-config list, get** → printer.Print(unmarshalList/Map(data))
- **tenants list, get** → printer.Print(unmarshalList/Map(data))
- **tools builtin list** → printer.Print(unmarshalList(data))
- **traces list, get** → printer.Print(unmarshalList/Map(data))
- **tts status** → ws.Call → printer.Print(unmarshalMap(data))
- **usage summary, breakdown** → printer.Print(unmarshalList/Map(data))

### JSON-Unfriendly (⚠️ table-only or TUI-driven):
- **agents.files list** — TUI output, no json flag
- **agents.instances list** — JSON response exists but unmarshalling varies
- **agents.links list** — No json output flag observed
- **agents.ops share/unshare** — Success/error only
- **agents.wake** — Single POST, no response
- **approvals approve/deny** — "Execution approved" text only
- **auth login/logout** — printer.Success() text only
- **auth whoami** — Table output, --output flag support unclear
- **chat interactive** — Real-time streaming, TUI-driven
- **chat abort** — Single action response
- **channels.contacts create** — Success message only
- **channels.pending list** — Limited JSON support
- **channels.writers add/remove** — Success text only
- **config set** — Success message only
- **config.permissions grant/revoke** — Success text only
- **contacts create, verify** — Success message or TUI prompt
- **credentials create** — Shows raw key once (not JSON)
- **cron delete, trigger** — Success message only
- **devices approve/reject** — WS call result unclear
- **export agent-preview** — JSON output but preview-only
- **heartbeat set, checklist set** — Success text only
- **import *** — "Import complete" message only
- **kg entities create, delete** — Success only
- **kg dedup execute-merge** — "Merge complete" text
- **logs tail** — Streaming plaintext logs, NDJSON in json mode
- **mcp servers create/update/delete** — Success only
- **mcp servers test** → printer.Print(unmarshalMap(data)) but test-specific
- **mcp grants grant/revoke** — "Access granted/revoked" text
- **memory store, index, delete** — Success message only
- **packages list** — Limited to simple list, no details
- **pending-messages send, delete** — Success text only
- **providers create, update, delete** — Success only
- **sessions label, delete, reset** — Success text only
- **skills upload, publish, unpublish** — Success/progress messages
- **skills tenant-config set** — Success text only
- **storage put, delete** — Success text only
- **system-config set, delete** — Success message only
- **teams create, delete** — Success text only
- **teams.members add/remove/reassign** — Success text only
- **teams.tasks create, delete, approve, reject** — Success text only
- **teams.workspace put, delete** — Success text only
- **tenants create, update, delete** — Success text only
- **tools custom create, delete** — Success text only
- **version** — Plaintext stdout (multi-line)
- **health** — "Server is healthy" text only

**Principle:** Commands marked with `printer.Success()` are **output-unfriendly** for Bash skill integration; use json-friendly alternatives where available or parse table output.

---

## 6. SUMMARY STATISTICS

| Metric | Count | Notes |
|--------|-------|-------|
| Top-level command groups | 38 | From rootCmd.AddCommand registrations |
| Subcommands (total) | ~180 | Estimated across all groups |
| HTTP-only commands | ~120 | CRUD operations, file I/O |
| WebSocket commands | ~60 | Real-time, subscribe, streaming |
| Destructive (delete/clear/reset) | ~40 | Require --yes or tui.Confirm |
| JSON-friendly | ~70 | Support `--output json` cleanly |
| Streaming (unsuitable for Bash) | ~8 | chat stream, logs tail, approvals watch, etc. |
| Interactive TUI | ~24 | Device pairing, auth, approvals |
| Global flags | 8 | --server, --token, --output, --yes, --insecure, --verbose, --profile, --tenant-id |

---

## 7. UNRESOLVED QUESTIONS & CONCERNS

1. **Approvals "watch" command:** admin.go references "approvals watch" but implementation not found in visible code. Is it implicit in approvals.list polling or WS subscribe?

2. **Teams WebSocket-only:** All teams operations use ws.Call/ws.Subscribe. What's the fallback if WS unavailable? HTTP pool?

3. **Chat interactive mode vs skill integration:** Interactive chat (readline loop) cannot be driven by Bash skill. Should skill wrap in expect/pexpect or switch to single-shot mode?

4. **TUI imports across files:** Multiple files import `internal/tui` (Confirm, Input, Password, IsInteractive). Does this pose usability challenges when --yes is set?

5. **Config set destructiveness:** config set modifies live server state. How is rollback handled? No --dry-run observed.

6. **Memory user scoping:** memory list supports --user filter but storage at /v1/memory/{agentID}/{path}. Is per-user memory automatic or manual?

7. **Sessions vs chat.session.status:** Both exist (sessions.list/delete and chat.status). Which is authoritative? Can they diverge?

8. **KG dedup workflow:** Dedup requires scan → review → merge-candidates → execute-merge. Should skill enforce this as a transaction or allow partial runs?

9. **Export without --yes:** Does export risk large data dump without automation flag? No --yes observed.

10. **Agent.files TUI mode:** agents.files list uses TUI for display. Can it be forced to json via --output flag, or is TUI hardcoded?

11. **Cron HTTP fallback:** cronListCmd has HTTP fallback if WS unavailable. Are other commands similarly protected?

12. **MCP reconnect vs servers.test:** reconnect triggers async reconnect; servers.test is sync. How are they distinguished in skill docs?

13. **Skill upload multipart:** skillsUploadCmd uses multipart/form-data. What's the MIME type, content-disposition, and maximum size?

14. **Teams namespace collision:** teams, teams.members, teams.tasks exist. Is teams.tasks a subcommand of teams or independent? AddCommand suggests independent.

15. **Heartbeat.targets vs heartbeat.checklist:** Both exist under heartbeat. Are they the same resource or separate? docs unclear.

---

## 8. SKILL DEVELOPMENT GUIDANCE

### Recommended Skill Structure (12-18 reference files):

```
~/.claude/skills/goclaw/
├── references/
│   ├── auth-and-config.md
│   ├── agents-core.md
│   ├── agents-advanced.md
│   ├── chat-sessions.md
│   ├── knowledge-memory.md
│   ├── teams-collaboration.md
│   ├── channels-messaging.md
│   ├── data-movement.md
│   ├── providers-skills.md
│   ├── automation-scheduling.md
│   ├── mcp-integration.md
│   ├── monitoring-ops.md
│   ├── admin-system.md
│   └── docs-api.md
├── SKILL.md (main entry point)
└── helpers/ (optional: example scripts)
```

### Critical Patterns for Skill:

1. **All commands accept `-o json` or `--output json` globally** (from root.go). Recommend JSON output for Bash skill parsing.

2. **--yes flag available globally.** Wrap all destructive operations with confirmation prompt in skill unless explicitly approved.

3. **WebSocket vs HTTP:** Streaming commands (chat, logs, approvals) MUST use Monitor tool, not Bash. HTTP commands safe for simple Bash invocation.

4. **Table output default.** When no --output flag set, commands emit table format. Skill should parse JSON by appending `--output json`.

5. **Profile/tenant context:** Global --profile and --tenant-id flags control context. Skill should allow setting these per invocation.

6. **Error handling:** Commands use `SilenceErrors: true` on rootCmd, so stderr is suppressed. Skill must capture exit codes for error detection.

7. **Interactive prompts disabled with --yes.** Paired device pairing, config confirmations, contact verification all respect cfg.Yes flag.

---

**Report Status: DONE**

