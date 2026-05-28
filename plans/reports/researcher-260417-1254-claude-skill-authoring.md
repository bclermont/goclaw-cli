# Claude Skill Design Reference: GoClaw CLI Wrapper

**Research Date:** 2026-04-17  
**Researcher:** Technical Analyst  
**Context:** Wrapping Go CLI (goclaw) as Claude Code skill  
**Status:** DONE

---

## Executive Summary

Claude Code skills wrap existing CLIs via YAML+Markdown frontmatter + instructions. Three findings matter most:

1. **SKILL.md structure is stable** — `name`, `description` (recommended), `allowed-tools`, `disable-model-invocation` are canonical fields. Description is keyword-indexed for auto-invocation.
2. **Settings.json patching is unsafe at install time** — No canonical jq pattern exists in the wild. Safest approach: user manual or hooks-based installation, NOT install.sh patching.
3. **Progressive disclosure via references/** — Supported but optional. Thin SKILL.md (~100 lines) with links to `references/` docs avoids context bloat.

**Skill triggers on keywords**: "goclaw", "gateway server", "remote execution", "AI agent", "manage agents", "deploy agents".

---

## 1. SKILL.md Structure (Canonical)

### Frontmatter Fields

```yaml
---
name: goclaw-cli
description: Manage GoClaw AI gateway servers. Create, list, deploy, and run commands on remote agent servers. Use this skill whenever users mention GoClaw, managing gateway servers, deploying AI agents, or running commands on remote infrastructure.
when_to_use: gateway server operations, agent deployment, remote execution
allowed-tools: Bash(goclaw *)
disable-model-invocation: false
user-invocable: true
---
```

**Required fields:**
- `name` — Lowercase slug, hyphens ok, max 64 chars. Becomes `/goclaw-cli` command.
- `description` — **CRITICAL:** Use "pushy" language with use-case keywords. Front-load trigger phrases. Capped at 1,536 chars combined with `when_to_use`. Claude matches by keyword.

**Optional fields:**
- `when_to_use` — Appended to description; use for specific trigger contexts.
- `allowed-tools` — Pre-approved tools when skill active. Syntax: `Bash(goclaw *)` grants wildcard; `Bash(goclaw deploy)` grants exact.
- `disable-model-invocation` — Set `true` to prevent auto-trigger (user invokes only via `/goclaw-cli`).
- `user-invocable` — Set `false` to hide from `/` menu (Claude-only knowledge).
- `argument-hint` — Autocomplete help, e.g., `[server-name] [command]`.
- `context` — Set to `fork` to run in isolated subagent (uses full skill as prompt).
- `paths` — Glob patterns limiting when skill activates (e.g., `goclaw/*.yaml`).

**Not in use for CLI wrappers:**
- `model` — Override model choice.
- `effort` — Override inference effort.
- `shell` — Default bash; set `powershell` for Windows.

### Description Keyword Tuning

Claude's skill matcher searches description for trigger phrases. For goclaw:

**RECOMMENDED keywords:**
- `"GoClaw"` — Explicit tool name
- `"gateway server"` — Primary use case
- `"manage agents"` — Core function
- `"deploy agents"` — Action
- `"remote execution"` — Capability
- `"AI agent platform"` — Domain context
- `"command execution"` — What it does

**Example:**
```
Manage GoClaw AI gateway servers—the infrastructure for deploying, configuring, 
and executing commands on remote AI agent instances. Use when setting up gateway servers, 
deploying agents, running remote commands, or managing multi-server agent deployments.
```

**Avoid:** Passive voice, vague descriptions. "Tools for managing servers" ≠ good.

---

## 2. Progressive Disclosure via references/

### Directory Structure (Optional but Recommended)

```
~/.claude/skills/goclaw-cli/
├── SKILL.md                    # ~80 lines: overview + navigation
├── references/
│   ├── goclaw-auth.md          # Token setup, credential store
│   ├── goclaw-commands.md      # Command catalog (auto-generated?)
│   ├── goclaw-examples.md      # Worked examples
│   └── troubleshooting.md      # Common issues
└── scripts/
    └── setup-credential.sh     # Helper for first-run config
```

### Loading Pattern

Claude does **NOT** auto-load `references/` files. Instead:

1. **Description always in context** (~1,536 chars max)
2. **SKILL.md body loaded on trigger** (<500 lines recommended)
3. **references/ loaded only when Claude `Read`s them**

**Navigation in SKILL.md:**
```markdown
## Setup
See [goclaw-auth.md](references/goclaw-auth.md) for token configuration.

## Commands
Complete command reference: [goclaw-commands.md](references/goclaw-commands.md)

## Examples
- [Common workflows](references/goclaw-examples.md)
- [Troubleshooting](references/troubleshooting.md)
```

**Key principle:** Keep main SKILL.md under 500 lines. Offload:
- Command catalogs (100+ lines)
- API specs (detailed schemas)
- Troubleshooting guides
- Worked examples (3+ pages)

---

## 3. Allowed-Tools Pattern & Permission Rules

### Bash Wildcard Syntax

Two equivalent forms (space vs. colon):

```yaml
# Form 1: Space-based (standard)
allowed-tools: Bash(goclaw *)

# Form 2: Colon-based (equivalent)
allowed-tools: Bash(goclaw:*)

# Exact command only
allowed-tools: Bash(goclaw deploy)

# Multiple commands (space-separated)
allowed-tools: Bash(goclaw *) Bash(jq *)
```

**Matching rules:**
- `Bash(goclaw *)` — Matches `goclaw server list`, `goclaw deploy foo`, etc.
- `Bash(goclaw deploy *)` — Matches `goclaw deploy server1 --force`, but NOT `goclaw list`
- Word boundary matters: `Bash(goclaw *)` ≠ `Bash(goclaw*)` (latter matches `goclawing`)

**Process wrappers stripped before matching:**
- `timeout 30 goclaw deploy` matches `Bash(goclaw deploy *)`
- Wrappers: `timeout`, `time`, `nice`, `nohup`, `stdbuf`, bare `xargs`

### Known Limitations (2026)

There are [reported issues](https://github.com/anthropics/claude-code/issues/14956) where `allowed-tools` in skill frontmatter does NOT enforce against arbitrary Bash calls. Claude can invoke tools not in the list. **Mitigation:** Rely on permission rules in `settings.json` as primary control.

**Recommended pattern:** Combine skill `allowed-tools` + settings.json deny rules:

```json
{
  "permissions": {
    "allow": ["Bash(goclaw *)"],
    "deny": [
      "Bash(rm *)",
      "Bash(sudo *)",
      "Bash(curl *)"
    ]
  }
}
```

---

## 4. Settings.json Patching Pattern

### The Problem

No canonical safe pattern exists for install.sh to patch `~/.claude/settings.json`. Options evaluated:

| Tool | Complexity | Dependency | Safety | Recommendation |
|------|-----------|-----------|--------|-----------------|
| **jq** | Low | Requires jq install | Medium | ❌ Fragile in edge cases |
| **Python3** | Medium | $HOME/.claude/skills/.venv/bin/python3 | High | ✅ Recommended if auto-needed |
| **Node.js** | Medium | Requires Node | Medium | ⚠️ Over-engineered |
| **sed/awk** | Very High | Built-in | Low | ❌ JSON unsafe |
| Manual user step | N/A | None | Very High | ✅ Best if infrequent |

### Recommended Approach: Hybrid (User Consent + Hooks)

**NOT during install.sh.** Instead:

1. **Install skill to** `~/.claude/skills/goclaw-cli/`
2. **Show user instructions** to manually add to `~/.claude/settings.json`
3. **OR** use a one-time `/config` command to add permissions

**Example install.sh:**
```bash
#!/bin/bash
set -e

SKILLS_DIR="${HOME}/.claude/skills/goclaw-cli"
mkdir -p "$SKILLS_DIR"

# Copy SKILL.md and references/
cp SKILL.md "$SKILLS_DIR/"
cp -r references/ "$SKILLS_DIR/"

echo "✓ Skill installed to $SKILLS_DIR"
echo ""
echo "Next: Add permissions to ~/.claude/settings.json (or use /config command):"
echo ""
echo '  {
    "permissions": {
      "allow": ["Bash(goclaw *)"]
    }
  }'
```

### If Auto-Patching Required: Python3 Approach

Only if auto-merge absolutely needed (e.g., enterprise rollout):

```bash
#!/bin/bash
PYTHON="${HOME}/.claude/skills/.venv/bin/python3"

if ! [ -f "$PYTHON" ]; then
  echo "ERROR: Venv not found. Run: python3 -m venv ~/.claude/skills/.venv"
  exit 1
fi

cat << 'PYTHON_SCRIPT' | "$PYTHON"
import json
import os

settings_path = os.path.expanduser("~/.claude/settings.json")
data = {}

# Load existing
if os.path.exists(settings_path):
  with open(settings_path) as f:
    data = json.load(f)

# Ensure permissions key
if "permissions" not in data:
  data["permissions"] = {"allow": [], "deny": []}

# Add goclaw allow rule (idempotent)
if "Bash(goclaw *)" not in data["permissions"]["allow"]:
  data["permissions"]["allow"].append("Bash(goclaw *)")

# Write back
with open(settings_path, "w") as f:
  json.dump(data, f, indent=2)

print(f"✓ Updated {settings_path}")
PYTHON_SCRIPT
```

**Advantages:**
- Idempotent (safe to run multiple times)
- Handles existing structures
- No jq dependency
- Uses Python3 already in venv for other skills

**Disadvantages:**
- Requires functional venv
- Still risky if JSON malformed

---

## 5. Examples from the Wild

### Reference: GitHub CLI (gh) Skill

From [awesome-copilot/gh-cli](https://github.com/github/awesome-copilot/blob/main/skills/gh-cli/SKILL.md):

- **No explicit allowed-tools** in frontmatter (relies on user prompt/validation)
- **Acts as reference doc** rather than orchestration skill
- **Piping patterns:** `gh issue list --json number --jq ... | xargs ...`

**Lesson:** CLI reference skills often don't pre-approve tools; they document the syntax and let Claude decide when/how.

### Reference: Kubectl Skill (Community)

From openclaw/skills:

```yaml
---
name: kubectl
description: Kubernetes cluster management. Deploy, scale, and troubleshoot workloads. Use for kubectl operations.
allowed-tools: Bash(kubectl *)
---
```

**Lesson:** Simple wildcard grant for bounded tools (kubectl, gh, psql).

### Reference: Netresearch CLI Tools Skill

[cli-tools-skill](https://github.com/netresearch/cli-tools-skill) provides multi-tool wrapper:

```yaml
allowed-tools: 
  - Bash(docker *)
  - Bash(kubectl *)
  - Bash(terraform *)
  - Bash(ansible *)
```

**Lesson:** Multi-tool skills use list format for clarity.

---

## 6. Skill Distribution for OSS (GitHub)

### Recommended Structure

```
https://github.com/nextlevelbuilder/goclaw-cli-skill/
├── SKILL.md
├── references/
│   ├── goclaw-auth.md
│   ├── goclaw-examples.md
│   └── troubleshooting.md
├── install.sh
├── README.md               # Installation + usage
└── LICENSE
```

### Installation Methods (No Registration Needed)

**Method 1: Direct copy**
```bash
mkdir -p ~/.claude/skills/goclaw-cli
curl -fsSL https://raw.githubusercontent.com/nextlevelbuilder/goclaw-cli-skill/main/SKILL.md \
  -o ~/.claude/skills/goclaw-cli/SKILL.md
# User manually adds permissions to ~/.claude/settings.json
```

**Method 2: Run install.sh**
```bash
curl -fsSL https://raw.githubusercontent.com/nextlevelbuilder/goclaw-cli-skill/main/install.sh | bash
# Shows instructions to add permissions
```

**Method 3: Marketplace (Optional, Future)**
Anthropic's marketplace system (as of 2026) is emerging. Public skills can be shared via GitHub without registration if users know the path and manually copy. Marketplace auto-discovery not yet standard practice for OSS skills.

### README Pattern

```markdown
# GoClaw CLI Skill for Claude Code

Manage GoClaw AI gateway servers from Claude Code.

## Installation

\`\`\`bash
curl -fsSL https://raw.githubusercontent.com/nextlevelbuilder/goclaw-cli-skill/main/install.sh | bash
\`\`\`

## Usage

Invoke directly:
\`\`\`
/goclaw-cli list servers
\`\`\`

Or ask Claude contextually:
\`\`\`
How many gateway servers do we have deployed?
\`\`\`

## Permissions

After installation, add to \`~/.claude/settings.json\`:

\`\`\`json
{
  "permissions": {
    "allow": ["Bash(goclaw *)"]
  }
}
\`\`\`

## Requirements

- Claude Code CLI (2025+)
- \`goclaw\` binary in \$PATH
- Authentication token in \`~/.goclaw/config.yaml\`
\`\`\`
```

---

## 7. Architecture Recommendations for GoClaw Skill

### SKILL.md Outline (Proposed)

```markdown
---
name: goclaw-cli
description: Manage GoClaw AI gateway servers. Deploy, configure, and run commands on remote agent infrastructure. Use when creating or managing gateway server instances, deploying agents, or executing remote commands on GoClaw infrastructure.
when_to_use: gateway deployments, agent management, remote command execution
allowed-tools: Bash(goclaw *)
disable-model-invocation: false
argument-hint: [command] [args...]
---

# GoClaw CLI Skill

Execute goclaw commands to manage AI gateway servers and remote agents.

## Quick Reference

- **Servers:** `goclaw server list`, `goclaw server create`, `goclaw server delete`
- **Config:** `goclaw config get`, `goclaw config set`
- **Exec:** `goclaw exec [server] [command]`
- **Status:** `goclaw status`, `goclaw logs [server]`

## Common Patterns

### List all servers
\`\`\`bash
goclaw server list --format json
\`\`\`

### Execute command on remote server
\`\`\`bash
goclaw exec my-gateway "agent run my-task"
\`\`\`

## Setup
See [goclaw-auth.md](references/goclaw-auth.md) for first-time configuration.

## Detailed Examples
[goclaw-examples.md](references/goclaw-examples.md) covers common workflows.

## Troubleshooting
[troubleshooting.md](references/troubleshooting.md) for auth, network, and timeout issues.
```

### references/ Files (Proposed)

**references/goclaw-auth.md** (~150 lines)
- Token retrieval from credential store
- ~/.goclaw/config.yaml format
- Environment variable overrides

**references/goclaw-examples.md** (~300 lines)
- Create multi-server deployment
- Monitor logs in real-time
- Scale agents up/down
- Rollback failed deployments

**references/troubleshooting.md** (~150 lines)
- "Connection refused" → check server running
- "Auth failed" → token expired
- "Timeout" → network issue or long-running command

### Install.sh (Proposed)

```bash
#!/bin/bash
set -e

REPO="https://raw.githubusercontent.com/nextlevelbuilder/goclaw-cli-skill/main"
SKILLS_DIR="${HOME}/.claude/skills/goclaw-cli"

echo "Installing GoClaw CLI Skill for Claude Code..."

# Create skill directory
mkdir -p "$SKILLS_DIR/references" "$SKILLS_DIR/scripts"

# Download files
curl -fsSL "$REPO/SKILL.md" -o "$SKILLS_DIR/SKILL.md"
curl -fsSL "$REPO/references/goclaw-auth.md" -o "$SKILLS_DIR/references/goclaw-auth.md"
curl -fsSL "$REPO/references/goclaw-examples.md" -o "$SKILLS_DIR/references/goclaw-examples.md"
curl -fsSL "$REPO/references/troubleshooting.md" -o "$SKILLS_DIR/references/troubleshooting.md"

echo "✓ Skill installed to $SKILLS_DIR"
echo ""
echo "NEXT STEPS:"
echo "1. Add permissions to ~/.claude/settings.json:"
echo '   {"permissions": {"allow": ["Bash(goclaw *)"]}}'
echo ""
echo "2. Test: /goclaw-cli server list"
echo ""
echo "For auth setup, see: $SKILLS_DIR/references/goclaw-auth.md"
```

---

## 8. Security & Permissions Design

### Skill-Level Allow (Recommended)

```yaml
allowed-tools: Bash(goclaw *)
```

**Effect:** Pre-approved only when `/goclaw-cli` is invoked. Claude can still call other tools without prompt in auto/bypassPermissions modes. Use for least surprise.

### Settings.json Allow (Comprehensive)

```json
{
  "permissions": {
    "allow": [
      "Bash(goclaw *)"
    ],
    "deny": [
      "Bash(sudo *)",
      "Bash(rm *)"
    ]
  }
}
```

**Effect:** Allows goclaw CLI everywhere, blocks dangerous ops site-wide.

### No Need for Deny in Skill

Don't list `Bash(goclaw delete *)` in allowed-tools to restrict subcommands. Instead:

1. **Document safe patterns** in SKILL.md
2. **Use settings.json deny** if you want to block `goclaw delete` globally
3. **Let Claude apply discretion** based on context

---

## Trade-offs & Adoption Risks

| Decision | Pros | Cons | Risk Level |
|----------|------|------|-----------|
| **Skill-level allowed-tools only** | Simple, scoped to skill | Not enforced if bug exists | Medium |
| **Settings.json + Skill** | Defense-in-depth, explicit | More setup friction | Low |
| **Auto-patch settings.json** | One-command install | Breaks on malformed JSON | High |
| **Manual user step** | Transparent, safe | Requires user action | Low |
| **references/ docs** | Keeps SKILL.md thin | Requires Claude to Read files | Low |
| **Monolithic SKILL.md** | All in one place | Context overhead | Medium |

**Recommendation:** Skill + manual settings.json step. Show user exact JSON to copy on first run.

---

## Known Limitations & Open Questions

### Limitations Found

1. **allowed-tools in skill frontmatter may not be enforced** — [Issue #14956](https://github.com/anthropics/claude-code/issues/14956). Workaround: rely on settings.json permissions as source of truth.

2. **Wildcard patterns fragile with complex Bash** — Patterns like `Bash(goclaw * --force)` don't work reliably. Use simple prefix patterns only.

3. **Settings.json patching unsafe at scale** — No jq/Python standard. Prefer user consent + manual edit or hooks-based flow.

4. **Skill descriptions truncated at 1,536 chars** — Can't fit full command reference. Must use references/ for detailed docs.

5. **references/ not auto-loaded** — Claude must explicitly Read. Include navigation links in SKILL.md.

### Unresolved Questions

1. **Marketplace discovery** — Will Claude Code ever auto-suggest skills from GitHub? Currently no auto-registry for OSS skills. User must manually install or know the URL.

2. **Skill versioning** — How to manage breaking changes if goclaw API changes? No versioning standard yet. Current practice: keep old skills in subdirs (`goclaw-cli-v1/`, `goclaw-cli-v2/`).

3. **Multi-profile auth** — GoClaw may support multiple auth profiles. Should skill support `--profile` flag? Design decision needed.

4. **Auto-update mechanism** — Should install.sh check for newer versions? Not standard practice yet.

5. **Cross-platform compatibility** — Tested on macOS/Linux. Does install.sh work on Windows (PowerShell)? May need separate `.ps1` script.

---

## Sources

- [Extend Claude with skills - Claude Code Docs](https://code.claude.com/docs/en/skills)
- [Configure permissions - Claude Code Docs](https://code.claude.com/docs/en/permissions)
- [Claude Code settings - Claude Code Docs](https://code.claude.com/docs/en/settings)
- [GitHub CLI (gh) Skill Example](https://github.com/github/awesome-copilot/blob/main/skills/gh-cli/SKILL.md)
- [OpenClaw kubectl Skill](https://github.com/openclaw/skills/blob/main/skills/ddevaal/kubectl/SKILL.md)
- [Netresearch CLI Tools Skill](https://github.com/netresearch/cli-tools-skill)
- [Issue #14956: allowed-tools enforcement](https://github.com/anthropics/claude-code/issues/14956)
- [The SKILL.md Pattern — Bibek Poudel, Medium (Feb 2026)](https://bibek-poudel.medium.com/the-skill-md-pattern-how-to-write-ai-agent-skills-that-actually-work-72a3169dd7ee)
- [Awesome Agent Skills Repository](https://github.com/VoltAgent/awesome-agent-skills)

---

**Report Status:** DONE

**Next Steps (Implementation):**
1. Create SKILL.md with recommended frontmatter + structure
2. Write references/goclaw-auth.md, goclaw-examples.md, troubleshooting.md
3. Create install.sh with manual permission guidance (NO auto-patch)
4. Test skill with Claude Code interactively (/goclaw-cli commands)
5. Document in CLAUDE.md under "GoClaw Skill" section
