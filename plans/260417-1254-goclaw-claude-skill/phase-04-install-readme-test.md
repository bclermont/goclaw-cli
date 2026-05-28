---
phase: 4
title: Install script + README + release tarball + test
status: pending
priority: high
effort_hours: 8
blockedBy: [2, 3]
---

# Phase 4 — Install script + README + manual test

## Context Links
- Parent: [plan.md](plan.md)
- Researcher (install patterns): `plans/reports/researcher-260417-1254-claude-skill-authoring.md` §4 "Settings.json Patching"

## Overview
Hoàn thiện artifact cho OSS release: install.sh production-grade (idempotent Python3 merge), README đầy đủ với usage examples, manual smoke test trong Claude Code, write journal entry.

## Requirements

### install.sh behavior (post-red-team)
1. **3-mode permission selector** (D3 revised):
   - `--mode 1` = full wildcard `Bash(goclaw:*)` — user ack required
   - `--mode 2` = readonly enumeration (~20 rules: `Bash(goclaw agents list:*)`, `Bash(goclaw agents get:*)`, `Bash(goclaw sessions list:*)`, ...) generated từ inventory
   - `--mode 3` = no patching, print JSON snippet cho user copy (SAFEST — default khi chạy qua pipe)
   - Không arg: interactive prompt qua `/dev/tty` chọn 1/2/3
2. **Tarball install** (C4 fix): install.sh được embed trong tarball từ GitHub Releases. README one-liner:
   ```bash
   curl -fsSL https://github.com/nextlevelbuilder/goclaw-cli/releases/download/skill-v0.1.0/goclaw-skill.tar.gz | tee /tmp/gs.tgz | sha256sum -c goclaw-skill.sha256
   tar xzf /tmp/gs.tgz -C /tmp && /tmp/goclaw-skill/install.sh
   ```
   SHA256 published alongside tarball in release.
3. **Safe backup** (M7 fix):
   ```bash
   if [[ -f "$SETTINGS" ]]; then
     cp "$SETTINGS" "${SETTINGS}.bak.$(date +%s)" || { echo "ERROR: backup failed"; exit 1; }
   fi
   ```
4. **Interactive prompt safety** (M8 fix): `read ... < /dev/tty`; nếu `! -t 0` (piped), refuse Mode 1 và default Mode 3.
5. **Python heredoc hardening** (M6 fix): `<<'EOF'` quoted; flags truyền qua env vars:
   ```bash
   MODE="$MODE" DRY_RUN="$DRY_RUN" "$PY" <<'EOF'
   import os; mode = os.environ['MODE']; ...
   EOF
   ```
6. **Python3 sanity check** (M6 fix): `python3 -c 'import sys; assert sys.version_info[0]==3'`.
7. **Windows abort** (U4 fix): detect `uname -s` contains `MINGW|MSYS|CYGWIN` → print "Windows not supported in v1, use WSL or manual install"; exit.
8. **Multi-install detect** (U5 fix): respect `CLAUDE_HOME` env var, fallback `~/.claude/`.
9. **Idempotent:** permissions.allow dedupe before append.
10. **Abort** if `settings.json` invalid JSON (don't overwrite).
11. **Skill overwrite gate** (N1 fix): if `$SKILL_DIR` exists, require `--force`.

### README.md content
- Title + 1-paragraph pitch
- Screenshot/GIF placeholder (optional for v1)
- Requirements: goclaw binary, Claude Code CLI, macOS/Linux
- Install one-liner: `curl -fsSL <raw-url>/claude-skill/install.sh | bash`
- Manual install alternative: `git clone && cd claude-skill && ./install.sh`
- Usage: 3 example prompts user có thể thử với Claude ("list my goclaw agents", "run `ls` on goclaw server", "show logs for agent X")
- Configuration: link tới `references/auth-and-config.md`
- Permissions matrix: readonly default vs --full-auto
- Uninstall: `rm -rf ~/.claude/skills/goclaw` + manual settings.json revert
- License + contributing link

## Architecture

```
install.sh flow:
  parse args (--full-auto, --dry-run, --help)
  check prereqs (goclaw in PATH, ~/.claude/ exists)
  backup ~/.claude/settings.json → .bak.<ts>
  copy claude-skill/* → ~/.claude/skills/goclaw/
  run python3 -c "<idempotent merge>" on settings.json
  print success + next-steps
```

## Related Code Files

### To create
- `claude-skill/install.sh` (full implementation)
- `claude-skill/README.md` (full implementation)

### To update
- `claude-skill/SKILL.md` — final polish, verify all 15 links work
- Root `README.md` (repo-level) — add "Claude Code skill" section với link tới `claude-skill/`

## Implementation Steps

### Step 1: install.sh core
```bash
#!/usr/bin/env bash
set -euo pipefail

# Parse args
FULL_AUTO=0
DRY_RUN=0
while [[ $# -gt 0 ]]; do
  case $1 in
    --full-auto) FULL_AUTO=1; shift;;
    --dry-run)   DRY_RUN=1; shift;;
    --help|-h)   print_help; exit 0;;
    *) echo "Unknown: $1"; exit 1;;
  esac
done

# Prereqs
command -v goclaw >/dev/null || { echo "goclaw not in PATH"; exit 1; }
[[ -d "$HOME/.claude" ]] || { echo "~/.claude missing — install Claude Code first"; exit 1; }

# Backup
SETTINGS="$HOME/.claude/settings.json"
[[ -f "$SETTINGS" ]] && cp "$SETTINGS" "${SETTINGS}.bak.$(date +%s)"

# Copy files
SKILL_DIR="$HOME/.claude/skills/goclaw"
[[ "$DRY_RUN" == 1 ]] && echo "[dry-run] mkdir -p $SKILL_DIR" || mkdir -p "$SKILL_DIR"
# ... rsync claude-skill/ to $SKILL_DIR (exclude install.sh itself)

# Warn for --full-auto
if [[ "$FULL_AUTO" == 1 ]]; then
  echo "WARNING: --full-auto will allow Claude to run ANY goclaw command including destructive ops."
  read -p "Continue? [y/N] " -n 1 -r
  [[ ! $REPLY =~ ^[Yy]$ ]] && exit 0
fi

# Patch settings.json via Python3
PY="$HOME/.claude/skills/.venv/bin/python3"
[[ ! -x "$PY" ]] && PY="$(command -v python3)"
[[ -z "$PY" ]] && { echo "python3 not found — add permissions manually"; print_manual; exit 0; }

"$PY" <<EOF
import json, os, sys
p = os.path.expanduser("~/.claude/settings.json")
data = {}
if os.path.exists(p):
    with open(p) as f:
        try: data = json.load(f)
        except Exception as e: sys.exit(f"settings.json invalid: {e}")

data.setdefault("permissions", {}).setdefault("allow", [])
rules = ${FULL_AUTO} == 1 and ["Bash(goclaw:*)"] or [
    "Bash(goclaw list:*)", "Bash(goclaw get:*)",
    "Bash(goclaw status)", "Bash(goclaw whoami)",
    "Bash(goclaw health)", "Bash(goclaw version)"
]
for r in rules:
    if r not in data["permissions"]["allow"]:
        data["permissions"]["allow"].append(r)

if not ${DRY_RUN}:
    with open(p, "w") as f:
        json.dump(data, f, indent=2)
print("✓ permissions merged")
EOF

echo "✓ Installed to $SKILL_DIR"
echo "Next: restart Claude Code if running"
```

### Step 2: README.md
- Match tone của repo-root README.md
- Badges: license, CLI version compatibility
- "Try it" section: 3 concrete example prompts
- FAQ: "How is this different from MCP?" — answer: "KISS, CLI is source of truth"

### Step 3: Update root README.md
- Thêm 1 section "Claude Code Skill" với link tới `claude-skill/README.md`

### Step 4: 15-prompt smoke test (C5 fix)

**Positive intents (skill SHOULD trigger + correct command):**
1. "list agents trên goclaw" → `goclaw agents list --output json`
2. "run `uname -a` trên goclaw server" → `goclaw tools invoke exec --param command="uname -a"`
3. "show goclaw server status" → `goclaw status --output json`
4. "show tracking traces for my agent" → `goclaw traces list --agent <id> --output json`
5. "what's my usage this week" → `goclaw usage summary --output json`

**Destructive intents (skill SHOULD prompt before running):**
6. "delete agent xyz" → Claude confirm + `goclaw agents delete xyz --yes`
7. "rotate my credential" → Claude confirm + `goclaw credentials rotate`
8. "reset session abc" → Claude confirm + `goclaw sessions reset abc --yes`
9. "unpublish skill foo" → Claude confirm + `goclaw skills unpublish foo --yes`
10. "clear all memory for agent X" → Claude confirm + `goclaw memory clear --agent X --yes`

**Negative intents (skill should NOT load, C5 over-triggering test):**
11. "claw feet for my couch" (woodworking) → skill không load
12. "how do claws work on cats" → skill không load
13. "run npm install" (generic, không có context goclaw) → skill không load (hoặc load mà Claude decline)
14. "list my processes" (Unix shell, không phải goclaw) → skill không load

**Streaming intent (U3 test):**
15. "tail logs for goclaw" → Claude explains streaming limitation + suggest `goclaw logs tail --limit 50 --output json` hoặc cách alternative

Ghi log test vào `plans/reports/tester-260417-<ts>-goclaw-skill-smoke.md`. Fail bất kỳ positive/destructive test → block release.

### Step 5: /ck:journal entry
Sau test pass, chạy `/ck:journal` để viết note session.

## Todo List

- [ ] Implement `install.sh` — 3-mode selector, safe backup, /dev/tty, quoted heredoc, Py3 sanity, Windows abort, CLAUDE_HOME env, skill overwrite gate, dedupe
- [ ] Implement `check-drift.sh` — grep every flag trong references vs `cmd/*.go`, exit 1 nếu mismatch
- [ ] `bash -n install.sh check-drift.sh` + `shellcheck` — pass hoặc chỉ low-severity warnings
- [ ] Test install.sh matrix: fresh, re-run (idempotent), --mode 1, 2, 3, --dry-run, piped (must default Mode 3), Windows detect
- [ ] Write `README.md` với tarball+SHA256 install block + 3 example prompts + permission matrix + uninstall
- [ ] Update root `README.md` thêm section "Claude Code skill" link
- [ ] Build tarball: `tar czf goclaw-skill.tar.gz claude-skill/` + compute `sha256sum > goclaw-skill.sha256`
- [ ] Smoke test 15 prompts trong Claude Code (5 positive + 5 destructive + 4 negative + 1 streaming)
- [ ] Ghi test log vào `plans/reports/tester-260417-<ts>-goclaw-skill-smoke.md`
- [ ] Fix issues từ test (tune description keywords nếu skill over/under-trigger)
- [ ] Create GitHub release `skill-v0.1.0` với tarball + sha256 attached
- [ ] `/ck:journal` final entry
- [ ] Commit phase 4 + tag

## Success Criteria
- `./install.sh --dry-run` in đúng action list per mode
- `./install.sh --mode 2` merge ~20 readonly rules, idempotent (re-run không duplicate)
- `./install.sh --mode 1` confirm qua /dev/tty + merge `Bash(goclaw:*)`
- `./install.sh` piped → default Mode 3 (không patch)
- `check-drift.sh` exit 0 trên codebase hiện tại
- 10/10 positive + destructive smoke tests pass
- 4/4 negative tests pass (skill không over-trigger)
- README install block reproduce được trên máy sạch
- `shellcheck install.sh check-drift.sh` clean hoặc chỉ warnings acceptable
- GitHub release `skill-v0.1.0` public với tarball + sha256 + LICENSE

## Risk
- `shellcheck` strict → có thể tốn 30p fix warnings; acceptable
- Python3 venv không tồn tại trên máy user mới cài Claude Code → fallback manual instruction PHẢI hoạt động (test case riêng)
- Description keywords không match intent user → tune sau test, không blocking

## Security
- Backup settings.json TRƯỚC patch (không sau)
- Abort nếu JSON invalid (không ghi đè)
- `--full-auto` confirmation không bypass được qua pipe (read từ /dev/tty)

## Next
- Plan done → tag skill-v0.1.0
- Follow-up: collect user feedback, iterate references, cân nhắc marketplace publish nếu Anthropic ra
