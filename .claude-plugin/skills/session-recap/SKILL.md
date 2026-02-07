---
name: session-recap
description: Generate token-efficient session recap and save to file. Use when user says "recap", "summarize session", "what did we do", "session summary", or at end of work session. Opens in Zed after saving.
---

# Session Recap

Generate a structured, token-efficient summary of the current session.

## Workflow

1. **Analyze session** - Review conversation for: tasks completed, files changed, decisions made, learnings, pending items
2. **Write recap** - Save to `~/Documents/session-recaps/YYYY-MM-DD-HH-MM.md`
3. **Open in Zed** - `zed <file>`
4. **Bridge to personal-evolution** - Ask: "Any lessons or principles to capture? (y/n)"
   - If yes → Use `personal-evolution` skill Mode 3 (Capture) to record lesson

## Output Format

```markdown
# Session Recap - YYYY-MM-DD HH:MM

**Project**: [repo/project name from cwd]
**Duration**: ~[estimate based on conversation length]

## Completed
- [Task 1]: [1-line what was done]
- [Task 2]: [1-line what was done]

## Files Changed
| File | Action | Summary |
|------|--------|---------|
| `path/to/file` | created/modified/deleted | [what changed] |

## Key Decisions
- [Decision]: [rationale, 1 line]

## Learnings
- [New knowledge discovered during session]

## Pending
- [ ] [Unfinished task or follow-up]

## Commands Reference
[Only if significant commands were run that user might want to re-run]
```

## Token Efficiency Rules

- **Tables over prose** - Use tables for file changes, not paragraphs
- **1 line per item** - No multi-sentence bullets
- **Skip empty sections** - Don't include sections with no content
- **No duplication** - Each fact appears once
- **Action verbs** - "Added X" not "X was added to the project"

## Example

```markdown
# Session Recap - 2026-02-06 15:30

**Project**: url_go_scan
**Duration**: ~45 min

## Completed
- Investigated suffix intelligence: documented EnumerateSubdomains, trust_suffix vs trust_exact, filtering logic
- Investigated category update pipeline: documented async goroutine trigger, 4 gate checks, Category Intel API write

## Files Changed
| File | Action | Summary |
|------|--------|---------|
| `~/.claude/skills/url-scan-system/references/architecture/suffix-intelligence.md` | created | Domain trust hierarchy docs |
| `~/.claude/skills/url-scan-system/references/architecture/category-update-pipeline.md` | created | go_scan category write flow |
| `~/.claude/skills/url-scan-system/SKILL.md` | modified | +6 routing table entries, +2 changelog |

## Learnings
- Category update only triggers on content scan callback (not initial scan), requires severity >= 2
- Suffix intel uses filterSuffixOnlyIntel() to prevent exact-match propagation to subdomains

## Pending
- [ ] Test suffix intelligence with actual domain queries
```

## File Location

Recaps saved to: `~/Documents/session-recaps/`

Filename pattern: `YYYY-MM-DD-HH-MM.md` (allows multiple recaps per day)

Create directory if missing:
```bash
mkdir -p ~/Documents/session-recaps
```

## Cross-Skill Reference

After recap, this skill bridges to `personal-evolution` for lesson capture:
- Facts (this skill) → Growth (personal-evolution)
- See `personal-evolution` > Mode 3: Capture for lesson template
