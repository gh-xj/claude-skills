---
tags: [audit, coverage, gaps, verification, quality]
related: [maintenance.md, bootstrap.md]
---

# Coverage Audit

Procedures for verifying architecture skill completeness and accuracy.

## Repo/Service Coverage Matrix

Build a matrix showing which parts of the system are documented:

```markdown
| Repo | In Skill? | Reference Files | Source Indexes | Last Verified |
|------|-----------|----------------|----------------|---------------|
| {repo-a} | Yes | 3 | Yes | 2026-02-04 |
| {repo-b} | Yes | 2 | Partial | 2026-01-15 |
| {repo-c} | No | - | - | Never |
```

For each gap (repo/service not covered):
- Is it core or peripheral?
- Does any existing reference partially cover it?
- Priority: **High** (core service, shared storage), **Medium** (integration point), **Low** (utility, monitoring)

## Source Index Verification

For each reference file with a `## Source Files` table:

### Automated Check

```bash
# Extract function names from Source Files tables and verify they exist
grep -rn "func {FunctionName}" {repo_path}/
```

### What to Check

- [ ] File still exists at the stated path
- [ ] Function still exists (by name — line may have shifted)
- [ ] No major functions in that area MISSING from the index
- [ ] Line numbers are within ~20 lines of actual (acceptable drift)

### Report Format

```markdown
| Reference File | Entries | Verified | Stale | Missing |
|---------------|---------|----------|-------|---------|
| core/pipeline.md | 16 | 14 | 2 | 0 |
| architecture/data-flows.md | 12 | 11 | 0 | 1 |
```

**Stale**: function exists but file/line has shifted significantly.
**Missing**: important function in the same area not yet indexed.

## Cross-Reference Validation

Verify all cross-references resolve to existing files.

### Step 1: Check internal references

```bash
# Find all references/ paths in the skill
grep -rn "references/" ~/.claude/skills/{skill-name}/ --include="*.md"
# Verify each path points to an existing file
```

### Step 2: Check cross-skill references

```bash
# Find cross-skill references (pattern: skill-name > references/)
grep -rn ">" ~/.claude/skills/{skill-name}/ --include="*.md" | grep "references/"
# Verify the referenced skill and file both exist
```

### Step 3: Check for orphaned files

```bash
# List all reference files
find ~/.claude/skills/{skill-name}/references/ -name "*.md"
# Compare against files actually referenced from SKILL.md or other references
# Any file that exists but is never linked to = orphan
```

## Gap Report

Use the template from `references/templates.md` > Coverage Audit Report.

### Prioritizing Gaps

| Priority | Criteria |
|----------|----------|
| **High** | Core service undocumented. Shared storage missing from matrix. Cross-repo contract not mapped. |
| **Medium** | Source indexes incomplete. External config not snapshotted. Peripheral service undocumented. |
| **Low** | Line numbers drifted. Changelog not up to date. Minor formatting issues. |

### Actions for Each Gap Type

| Gap Type | Recommended Action |
|----------|-------------------|
| Missing repo coverage | Run discovery (see `references/discovery-playbook.md`) and create reference files |
| Stale source indexes | Run refresh (see `references/maintenance.md` > Source Index Refresh) |
| Missing cross-references | Wire integration (see `references/cross-skill-integration.md`) |
| No config snapshots | Capture initial snapshot (see `references/maintenance.md` > Config Snapshots) |
| Outdated architecture diagram | Re-explore and update SKILL.md service map |

## Audit Cadence

| Trigger | Scope |
|---------|-------|
| After any cross-repo change | Affected reference files only |
| Monthly | Full source index verification |
| After major refactor | Full audit (all checks) |
| Before onboarding (someone new explores the system) | Full audit to ensure accuracy |
| After adding a new service | Coverage matrix + discovery for new service |
