# Red Team Review: GoClaw Claude Skill Plan

**Reviewer:** adversarial / staff eng
**Target:** `plans/260417-1254-goclaw-claude-skill/` (plan.md + 4 phase files)
**Date:** 2026-04-17
**Verdict up front:** Ship AFTER critical fixes. Plan has structural soundness but several load-bearing assumptions unverified, effort underestimated, premature OSS scope, and a LICENSE blocker ignored.

---

## TL;DR — Kill / Keep / Modify

### Kill
- **"exec" hero use case as framed** (Phase 2 `exec-workflow.md`). `goclaw tools invoke exec --param command=...` is an assumption — no such hardcoded tool name in CLI source. It's `tools invoke <name>`; whether server registers `exec` as a tool is unverified. Verify BEFORE Phase 2 or the hero reference is wrong on day 1.
- **15 reference files as v1 scope.** YAGNI violation. Ship 3-5; iterate.
- **`--full-auto` install flag.** Premature. One feature flag maintaining two permission permutations for a one-user + speculative OSS case. Cut it from v1; add when a user asks.
- **"≥ 90% correct flag" success criterion.** Unfalsifiable as written. Either define measurement methodology or drop.
- **`Bash(goclaw whoami)` exact permission pattern** (Phase 4 step 1). Does NOT match `goclaw whoami --output json` — the skill's own convention is to always append `--output json`, so this rule will never fire. Same for `status`, `health`, `version`.

### Keep
- Overall structure (SKILL.md + references + install.sh + README)
- Python3 merge approach for settings.json (already-verified pattern from researcher report)
- Backup-before-patch of settings.json
- Single-repo hosting decision (D1)
- Convention "always append `--output json`" (D6)

### Modify
- Plan count: says 38 top-level groups, actually **36** registered in `rootCmd.AddCommand` (grep verified). `delegations` IS registered as top-level (admin.go:141), not "implied" as explorer claimed. `media` group exists (`admin_media.go:59`) but is missing from plan clusters entirely.
- Merge `approvals` into ONE reference (currently split across `exec-workflow.md` AND `chat-sessions.md` per explorer cluster mapping — duplicate content risk).
- Effort labels (S/M/M/M) — replace with concrete hours; see Major Issue #3.
- Permission matrix — redesign around actual wildcard semantics (see Critical #2).

---

## Critical issues (must fix before implementation)

### C1. LICENSE file does not exist at repo root — OSS publish blocked on day 1
Checked: `ls /Volumes/GOON/www/nlb/goclaw-cli/LICENSE` → **No such file or directory**. Phase 1 step 1 says "Check license file exists at repo root. If yes, copy... If no, note for Phase 4." Phase 4 does not include "write LICENSE" in its todo list — it's invisible. Result: you hit Phase 4 and realize Goal 4 ("OSS-publishable") is blocked waiting on a legal decision (MIT vs Apache-2.0?) that needs user input.

**Fix:** Before Phase 1, ask user which license. Add `write LICENSE` to Phase 1 todo list explicitly. Don't defer to Phase 4.

### C2. Permission rule patterns likely DO NOT enforce as intended
Plan's install.sh adds:
```
Bash(goclaw list:*)   Bash(goclaw get:*)
Bash(goclaw status)   Bash(goclaw whoami)
Bash(goclaw health)   Bash(goclaw version)
```

Two problems:

(a) **Skill convention says always append `--output json`** (D6 in plan). So Claude will actually run `goclaw status --output json`, which does NOT match exact pattern `Bash(goclaw status)`. Same for whoami, health, version. Four of six readonly rules are dead code. User will hit approval prompt for read-only commands — defeating the "readonly default safe" design.

(b) **`Bash(goclaw list:*)` assumes `list` is a top-level subcommand.** It is NOT. Structure is `goclaw agents list`, `goclaw sessions list`, `goclaw cron list`, etc. The pattern `goclaw list:*` matches literal string `goclaw list something` — won't match `goclaw agents list`. Plan should use patterns like `Bash(goclaw agents list:*)`, `Bash(goclaw sessions list:*)`, etc., OR use `Bash(goclaw * list *)` if wildcards in middle positions work (researcher report: "Wildcard patterns fragile with complex Bash"). More likely answer: enumerate safe leaf commands.

**Fix:**
1. Verify empirically in Claude Code which wildcard patterns actually enforce. Don't ship based on assumed syntax.
2. After `--output json` added, rules need trailing wildcard: `Bash(goclaw status *)`, `Bash(goclaw whoami *)` etc., OR use `Bash(goclaw * --output json)` (if middle wildcards work).
3. For read-only by verb, must fan out per resource group: `Bash(goclaw agents list:*)`, `Bash(goclaw agents get:*)`, `Bash(goclaw sessions list:*)`, ~20 rules minimum. Scope creep on install.sh.

### C3. "exec" tool name is unverified — hero reference may be fiction
Phase 2 centers on `goclaw tools invoke exec --param command="..."` as THE hero use case. Source code in `cmd/tools.go:111-143`:
```go
Use: "invoke <name>",
body := map[string]any{"name": args[0], "parameters": params}
data, err := c.Post("/v1/tools/invoke", body)
```
`name` is free-form. **No code anywhere hardcodes "exec" as a tool name.** It's assumed the server registers an `exec` tool. If it doesn't, or if it's named `shell`, `bash`, `run-command`, the hero reference is wrong.

Phase 2 risk section admits: "Schema response của `tools invoke exec` chưa known — phải verify bằng call thật lên dev server HOẶC đọc server code". This is fine to acknowledge — but it's BLOCKING for the phase labelled "hero". Cannot defer.

**Fix:** Before writing exec-workflow.md, run `goclaw tools builtin list --output json` against a live dev server and record the actual tool names + parameter schemas. Attach to phase 2 as prerequisite.

### C4. One-liner install via `curl | bash` has no integrity check
README (Phase 4): `curl -fsSL <raw-url>/claude-skill/install.sh | bash`.

No SHA256 verification, no pinned commit SHA, no signed release. User trusts:
1. github.com TLS not MITM'd (OK-ish)
2. The `main` branch HEAD at install time (supply-chain risk — attacker merges to main)
3. Any subsequent files the script curls from main (install.sh in Phase 4 doesn't detail this — does it curl down references from raw.githubusercontent.com one by one? If yes, that's 18+ HTTPS requests per install, any of which fails partially)

Plan doesn't specify fetch strategy. Options:
- Single tarball from GitHub releases (recommended) — one artifact, one SHA256 check.
- Script embeds a commit SHA and curls blob URLs at that SHA (reproducible).
- User runs `git clone && ./install.sh` (safest, slowest).

**Fix:** Phase 4 must specify install mode. Recommend: release tarball + SHA256 in README. `curl|bash` one-liner downloads tarball and verifies checksum. Without this, publishing to OSS is irresponsible — skill can install with modified `goclaw:*` permissions rule and user won't notice.

### C5. Phase 4 smoke test is inadequate
"Manual test 5 prompts" — what about the inverse:
- User says "give me claw feet for my couch" (woodworking) — does skill load? (false positive)
- User says "goclaw go claws" in a poem — does skill load? (over-triggering wastes context)
- User says "run a command on my server" (generic) — does skill load when it should? (under-triggering)
- Description keyword test across 10+ unrelated prompts to bound over/under-trigger rate.

Also: no regression test for "Claude hallucinates `--force` flag instead of `--yes`" — verify each reference against source code flags post-write.

**Fix:** Phase 4 test matrix expand to:
- 5 positive intents (currently in plan)
- 5 negative intents (skill should NOT load)
- 5 destructive intents (skill should prompt before running)
- Automated flag-verification pass: for each reference, grep every `--flag` name against corresponding `cmd/*.go` source.

---

## Major issues (should fix)

### M1. Effort estimate is wildly optimistic
Plan: Phase 2 = M, Phase 3 = M, Phase 4 = M. Reality:

| Phase | Claimed | Realistic (one human, focused) |
|-------|---------|--------------------------------|
| 1 (scaffold) | S | 1-2 hrs |
| 2 (5 refs, 150-250 lines each = ~1000 lines, requires reading ~10 cmd/*.go files, verifying flags, testing) | M | **8-12 hrs** |
| 3 (10 refs, ~2000 lines, same verification overhead) | M | **15-20 hrs** |
| 4 (install.sh production-grade + README + tests + journal) | M | **4-6 hrs** |

Total: **~30-40 hours** for a quality skill. Phase 3 todo list says "mỗi file 15-30 phút" — that's drastically underestimated for files requiring source verification of 5-11 subcommands, flag tables, 3-5 worked examples, and cross-links. 30 min is "skim source and write template boilerplate". You want 45-90 min per reference minimum for quality. At 10 refs × 1 hr = 10 hrs, plus 3 hrs cross-linking + testing = 13 hrs for phase 3 alone. Plan's implicit assumption of "one session" will produce fatigue-driven quality drop around file 8.

**Fix:** Either (a) cut to 5 refs for v1, iterate based on usage, or (b) split Phase 3 across multiple sessions with explicit "resume here" markers per file.

### M2. 15 references is YAGNI violation for stated scope
Goal 3 says "Cover toàn bộ 38 command groups" — but why? User's actual goal (per brainstorm citation `4/17 12:41`): "exec trên server là đủ". Covering `tenants`, `system-config`, `tts`, `media`, `api-docs open` upfront speculates need that may never materialize. Every reference written = maintenance debt (flag drift, staleness).

Compare: `gh-cli` skill in researcher report has no references/ at all — it's pure SKILL.md reference doc. Works fine.

**Fix:** v1 = 3 references:
1. `exec-workflow.md` (hero, if exec tool verified)
2. `auth-and-config.md` (prereq for any call)
3. `common-commands.md` (agents list/get, sessions list/get, status, logs — the 80% use)

Ship. Add more when user says "I want X and it's not covered". Measure actual usage before investing in 2000+ lines of maintenance surface.

### M3. Source-of-truth drift — no automated sync mechanism
Plan (Risk row 4): "goclaw CLI đổi flag → reference outdated. Mitigation: Version check ở top SKILL.md; link tới CHANGELOG." This is not mitigation — it's a README comment. When `goclaw agents create` adds `--temperature` flag in v1.3, nothing in the pipeline flags the skill as outdated.

Real mitigations (pick one):
- **CI check:** script that greps every flag in references against `cmd/*.go`, fails PR if mismatch.
- **Auto-gen with manual polish:** generate verified-flags table from `goclaw <cmd> --help`, commit alongside manual prose.
- **Version pinning in SKILL.md frontmatter:** `compatible-with: goclaw>=0.4.0,<0.5.0` and CI updates bounds on CLI release.

Plan dismisses auto-generation as "viết tay để control UX" but explicitly writing verified-flags table by hand guarantees drift. Hybrid approach (auto-gen flag table, hand-write prose) gives you both.

**Fix:** Add Phase 5 (or bake into Phase 4): write `claude-skill/check-drift.sh` that grep-validates flag names. Wire into Makefile + CI.

### M4. Approvals is placed in BOTH exec-workflow AND chat-sessions
Phase 2 says `exec-workflow.md` covers `approvals list/approve/deny`. Explorer cluster mapping puts `approvals` under `chat-sessions`. Phase 2 Priority list #3 is `chat-sessions.md` which also touches approvals. Duplicate content → maintenance burden + inconsistency risk.

**Fix:** Decide ONE canonical home. Since approvals fire from `tools invoke exec` (execution gating), keep in `exec-workflow.md`. `chat-sessions.md` cross-links only.

### M5. `media`, `packages`, `delegations` cluster assignments are wrong or missing
- `media` group (`cmd/admin_media.go:59`) exists. Not in any cluster in plan/explorer.
- `packages` (`cmd/packages.go:102`) lives where? Plan phase 3 puts it in `providers-skills-tools.md`. Packages are runtime binaries, not "skills" — confusing grouping.
- `delegations` is a root command (`cmd/admin.go:141` → `rootCmd.AddCommand(approvalsCmd, delegationsCmd)`), not "implied in agents.links" as explorer said. Plan phase 3 puts it under `agents-advanced.md`. Verify whether `delegations` is a separate resource or just an alias.

**Fix:** Run `goclaw --help` against actual built binary, extract authoritative top-level command list, re-cluster. Don't trust explorer table.

### M6. Python3 fallback path assumes `command -v python3` works — ancient macOS caveat
install.sh (Phase 4):
```bash
PY="$HOME/.claude/skills/.venv/bin/python3"
[[ ! -x "$PY" ]] && PY="$(command -v python3)"
[[ -z "$PY" ]] && { echo "python3 not found — add permissions manually"; ...}
```

Issues:
- `command -v python3` returns path but doesn't verify it's Python **3**. Some old macOS had `/usr/bin/python3` symlinked to `python2.7`. Rare in 2026 but possible.
- `[[ -z "$PY" ]]` check is too late — if `command -v` returns empty, `$PY=""` and `"$PY"` passes `-x` never because `[[ ! -x "$PY" ]]` when PY is empty: test fails falsy → fallback branch re-runs. Actual bug: if first check fails, second overwrites. Script then tries to run `"" <<EOF` which errors with "command not found" before hitting the null check.
- Script uses `${FULL_AUTO}` and `${DRY_RUN}` inside the Python heredoc — this is Bash substitution, NOT Python variable. The resulting Python becomes `rules = 1 == 1 and [...]` which works by accident but is fragile. A shellcheck lint will likely complain.

**Fix:**
- Use `python3 -c 'import sys; sys.exit(0 if sys.version_info[0]==3 else 1)'` sanity check.
- Reorder null check BEFORE invocation.
- Pass flags via env vars to Python (cleaner):
  ```bash
  FULL_AUTO=$FULL_AUTO DRY_RUN=$DRY_RUN "$PY" <<'EOF'
  import os
  full_auto = os.environ['FULL_AUTO'] == '1'
  ...
  EOF
  ```
  Note heredoc `<<'EOF'` (quoted) prevents bash substitution. Fixes injection risk too.

### M7. install.sh backup-on-failure behavior unspecified
install.sh: `[[ -f "$SETTINGS" ]] && cp "$SETTINGS" "${SETTINGS}.bak.$(date +%s)"`

If `cp` fails (disk full, permission denied on `.bak` target), `&&` short-circuits silently (exit 0 from the test, cp stderr printed but not caught). Script proceeds to mutate settings.json with NO backup. `set -e` does not save you — the `&&` chain has a non-zero exit only when cp fails, but bash with `set -e` treats `A && B` as conditional, not a failure path.

Test: `bash -euo pipefail -c '[[ -f /etc/hosts ]] && cp /etc/hosts /dev/null/x; echo reached'` — prints "reached", does not abort.

**Fix:**
```bash
if [[ -f "$SETTINGS" ]]; then
  cp "$SETTINGS" "${SETTINGS}.bak.$(date +%s)" || {
    echo "ERROR: Failed to backup settings.json, aborting" >&2
    exit 1
  }
fi
```

### M8. `--full-auto` confirmation reads from stdin — fails under `curl|bash`
Phase 4 security note: "`--full-auto` confirmation không bypass được qua pipe (read từ /dev/tty)."

Plan correctly identifies issue, but script as written uses `read -p "Continue? [y/N]" -n 1 -r` — **reads from stdin**. Under `curl|bash`, stdin is the piped script itself, so `read` consumes the next line of the script. Result: confirmation auto-approves with whatever char follows, OR script hangs if stdin already consumed.

**Fix:** `read -p "..." < /dev/tty` OR refuse to run under `--full-auto` when `[[ ! -t 0 ]]`. Plan mentions the fix in prose but not in code. Put it in code.

---

## Minor issues / nitpicks

### N1. SKILL.md name collision risk
Plan D2: skill name `goclaw`. Canonical path: `~/.claude/skills/goclaw/`. If user already has a skill named `goclaw` (from another source), install.sh overwrites silently. Add `--force` flag gating overwrite, default refuse with "existing skill found at ..., use --force to replace".

### N2. Kebab-case naming guideline vs actual file list
Plan CLAUDE.md says "Go snake_case file naming". But Phase 1 non-functional req says "Kebab-case cho tên reference files". The skill files are Markdown not Go, so kebab-case is fine — just note the mental switch to avoid accidental snake_case drift in reference filenames.

### N3. `chat-sessions.md` omits abort
Phase 2 agenda for `chat-sessions.md` lists chat send, single-shot, sessions CRUD. Explorer shows `chat abort` exists as destructive-without-`--yes`. Not in phase 2 content outline. Gap.

### N4. `auth pair` handling
Explorer: "auth pair = device pairing, poll 60× 2s, not streaming but long-running". Phase 2 outline notes "long-running polling — document 'not Bash-friendly'". Fine. But test matrix (Phase 4) doesn't include an auth-pair prompt. If Claude is asked "set up goclaw auth", what happens? Probably initiates `goclaw auth pair` and hangs the Bash tool for 120 seconds. Add explicit guidance: "skill should refuse to run pair; tell user to run manually".

### N5. Output path name convention leaks
Report naming convention in header says `code-reviewer-260417-1304-{slug}.md` but task requested `code-reviewer-260417-1254-goclaw-skill-red-team.md`. Inconsistency — the 1304 is from hook-injection time, 1254 is the plan's own timestamp. Pick one and stick. (I used 1254 per explicit user instruction.)

### N6. `effort` frontmatter field is legacy
Researcher report §1 says "`effort` — Override inference effort. Not in use for CLI wrappers." But Phase 1/2/3/4 YAML frontmatter all have `effort: S/M/M/M`. This is the plan's own tracking metadata, not SKILL.md frontmatter — so technically OK. But confusing naming overlap.

### N7. Cross-repo README link maintenance
Phase 4 step 3: "Update root `README.md` thêm link tới `claude-skill/README.md`". Fine, but this creates a doc dependency. If skill README moves or renames, root README breaks. Low-risk nitpick.

---

## Unaddressed risks

### U1. Skill "always append --output json" convention is unenforceable
D6 says Claude should always append `--output json`. But Claude's obedience to instructions varies with context pressure and conflicting signals. If user says "show me agents in a table", Claude might omit `--output json`. Downstream parsing then breaks because skill examples all assume JSON.

**Mitigation idea:** Wrap in a script. `goclaw-json` shim that forces `--output json` by default, errors on table mode. Skill instructs Claude to use `goclaw-json` not `goclaw`. Too heavy for v1; note as follow-up.

### U2. Token expiry mid-session
Plan says "Skill KHÔNG lưu token — dùng credential store của goclaw CLI". What if credential expires mid-session? `goclaw agents list --output json` exits with auth error. Skill has no reference covering "what does Claude do when goclaw returns 401?". Gap.

**Fix:** In `auth-and-config.md`, document exit codes and re-auth flow.

### U3. Streaming command over-triggering
User says "watch agent logs". Claude has a skill reference saying `logs tail` is streaming, unsupported. What does Claude actually do? Best case: explains limitation, suggests polling. Worst case: runs `goclaw logs tail`, Bash tool times out at 120s, user sees cryptic truncation. Not tested in Phase 4 matrix.

### U4. Windows users exist
Non-goal: "Windows PowerShell script — macOS/Linux only cho v1". Fine. But Claude Code runs on Windows. A Windows user who `git clone`s the repo will see install.sh and try to run it under git-bash. Most of it works, but `/dev/tty`, `date +%s`, and Python venv path all need adjustment. Plan doesn't say "Windows unsupported, print error and exit". Add that.

### U5. Multiple Claude Code installations
Some users have `~/.claude/` and `~/.config/claude/` both. install.sh hardcodes `~/.claude/`. If Claude Code changes default path in 2027, skill silently installs to stale location. Add env var override: `CLAUDE_HOME="${CLAUDE_HOME:-$HOME/.claude}"`.

### U6. "15 references" means 15× surface area for prompt injection
Not relevant today because user loads refs manually. But if references include server-returned data (e.g., example responses), a hostile gateway could inject prompts via error messages included in examples. Low-risk, but: rule "no server-rendered content in references, all examples synthetic".

---

## Concrete fix checklist (prioritized)

**BLOCKING — do before starting any phase:**
1. Decide LICENSE — write it in Phase 1.
2. Verify `exec` tool exists on server and capture its schema. Or rename hero to a verified tool.
3. Empirically verify `Bash(goclaw status)` vs `Bash(goclaw status *)` vs `Bash(goclaw status --output json)` wildcard behavior in actual Claude Code 2026.
4. Re-enumerate top-level groups from `goclaw --help` (binary output, not source grep) — fix count to 36 or whatever actual number, include `media`, confirm `delegations`.

**MUST fix in Phase 1:**
5. Fix LICENSE write into phase 1 todo.
6. Replace `effort: S/M/M/M` with hour estimates: Phase 1 = 2h, Phase 2 = 10h, Phase 3 = 18h (if kept at 10 refs), Phase 4 = 6h.

**MUST fix in Phase 2:**
7. Decide approvals canonical home (exec-workflow vs chat-sessions). Cross-link the other.
8. Cut Phase 3 to 5 refs OR explicitly budget 18 hours.

**MUST fix in Phase 4:**
9. install.sh:
   - Quote heredoc (`<<'EOF'`) and pass vars via env.
   - Robust backup (if-block, not `&&`).
   - `read < /dev/tty` for `--full-auto` confirm, or refuse when `! -t 0`.
   - Python3 sanity check (3.x verification).
   - Windows abort message.
10. README one-liner: use release tarball + SHA256, not raw main branch.
11. Test matrix: expand to 15 prompts (positive / negative / destructive).
12. Add `check-drift.sh` with flag grep validation.

**SHOULD consider:**
13. Drop `--full-auto` flag from v1 entirely. Add when demanded.
14. Cut reference count to 3-5 for v1. Mark phase 3 as "v0.2 backlog".
15. Revisit whether this should be a skill at all — a CLAUDE.md snippet + auto-approved `Bash(goclaw:*)` in project settings.json might serve one-user need. Skill is right answer only if OSS publish is confirmed priority (worth the maintenance cost).

---

## Recommendation

**Ship after critical fixes.** The plan is structurally coherent but carries ~8 unverified assumptions any one of which torpedoes UX on launch day. Specifically:
- If "exec" tool doesn't exist → hero reference is wrong.
- If wildcard permission patterns don't match `--output json` suffix → all readonly rules dead.
- If LICENSE undecided → OSS publish blocked.
- If effort underestimated 3× → Phase 3 ships as low-quality boilerplate.

Recommend downgrade to v0.1 scope: **3 references, no --full-auto flag, single tarball install with SHA256**. That's a 1-day project that proves the skill pattern works. Phase 3 + additional references become v0.2 driven by actual user demand. This is aligned with plan's stated YAGNI principle — which phase 3 violates.

**Do NOT** proceed with current plan as-is. Minimum viable correction: address C1-C5 before Phase 1.

---

## Unresolved questions

1. Does `exec` tool exist on GoClaw server, or is hero path fiction? (Blocking)
2. What is the empirical behavior of `Bash(goclaw status)` permission pattern against `goclaw status --output json`? (Blocking)
3. Which license (MIT / Apache-2.0 / proprietary)? (Blocking OSS)
4. Is OSS publish actually a goal, or is this a one-user skill? Plan treats OSS as goal; user's stated motivation was personal use.
5. Does `goclaw --help` list 36 or 38 top-level groups? Plan and explorer disagree with source code.
6. Does server register a `media` resource and is it in-scope for skill?
7. For `tools invoke`, what's the canonical JSON response schema? Referenced in Phase 2 Risk but never resolved.
8. What is expected behavior when Claude invokes a streaming command (`logs tail`) — abort at 120s Bash timeout, or should skill pre-empt?

---

**Status:** DONE_WITH_CONCERNS
**Summary:** Plan is well-structured but has 5 critical blockers (LICENSE missing, exec tool unverified, permission wildcards likely broken, no install integrity check, inadequate smoke test) and 8 major issues. Scope is 3× YAGNI — cut to 3 references for v1. Effort underestimated ~3×.
**Concerns:** Without verifying the three unverified-but-load-bearing assumptions (exec tool, wildcard semantics, license) before Phase 1, work will stall mid-Phase 2. Recommend parent agent route blocking questions back to user before dispatching implementer.
