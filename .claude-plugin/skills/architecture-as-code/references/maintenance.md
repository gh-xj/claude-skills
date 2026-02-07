---
tags: [maintenance, drift, source-index, config-snapshot, changelog, self-healing]
related: [coverage-audit.md, templates.md, bootstrap.md]
---

# Maintenance & Evolution

Protocols for keeping architecture skills accurate as codebases change.

## Source Index Refresh

Source file indexes (`## Source Files` tables) drift as code changes. Function names are the stable anchor; line numbers are helpful hints.

### Refresh Triggers

| Trigger | Scope |
|---------|-------|
| Major refactor in a covered repo | All reference files for that repo |
| Adding a feature that changes core flows | Affected pipeline/data-flow references |
| Before making a cross-repo change | Verify assumptions in affected references |
| Investigation reveals stale info | The specific reference file |

### Refresh Process

1. Read the reference file's Source Files table
2. For each entry, verify the function still exists:
   ```bash
   grep -rn "func {FunctionName}" {repo_path}/
   ```
3. Update file:line numbers for any that shifted
4. Add new entries for functions discovered during the session
5. Remove entries for deleted functions
6. Add a changelog entry in SKILL.md

### Handling Drift

If a function moved to a different file:
- Update the file path AND line number
- Check if the surrounding architecture description is still accurate

If a function was renamed:
- Update the function name in the Source Files table
- Grep for callers to verify the rename propagated

If a function was deleted:
- Remove the entry
- Check if the concept it represented still exists (maybe in a different function)
- Update the reference file's prose if the behavior changed

## Config Snapshots

For systems with external configs (rule engines, LLM prompts, feature flags), store date-stamped snapshots in the repo for git tracking.

### Directory Convention

See `references/templates.md` > Config Snapshot Directory for the full structure.

### Capture Process

1. Export config from the external system (API call, control panel export, etc.)
2. Create a new date-stamped directory: `detection-config/{type}/{YYYY-MM-DD}/`
3. Save the raw export in its original format (JSON, YAML, markdown)
4. Update the reference file that describes this config:
   - Change `snapshot:` date in frontmatter
   - Note any significant changes in the reference file body
5. Commit with: `chore: snapshot {config-type} {date}`

### Rules

- **Never overwrite** — always create a new date directory alongside old ones
- **Keep original format** — don't transform the export
- **Git diff between dates** — `git diff detection-config/{type}/{old-date}/ detection-config/{type}/{new-date}/`
- **Reference file points to snapshot** — agents can read the actual config content during investigation

### When to Snapshot

| Trigger | Action |
|---------|--------|
| Rule engine rules changed | Snapshot the rule export |
| LLM prompts updated | Snapshot the prompt content |
| Feature flags modified | Snapshot the flag configuration |
| Before and after a config migration | Snapshot both states |
| Periodically (monthly) | Snapshot to track gradual changes |

## Drift Detection

Every generated architecture skill should include self-healing instructions. Copy this section into the skill's SKILL.md:

```markdown
## Knowledge Freshness

### Before trusting: verify critical claims
When this skill drives a decision with cross-repo impact (e.g., "safe to change
this field" or "no other service reads this prefix"), verify against live code
before committing. Use Explore agents to confirm:
- Storage key readers/writers haven't changed
- Fields referenced in the answer still exist
- Service contracts match the current IDL/Thrift definitions

### After discovering: update this skill
If during a session you discover something that contradicts, extends, or refines
this skill:
1. Read the relevant reference file
2. Update it with the new information
3. Add a changelog entry

This makes every session a self-healing opportunity. The skill improves with use.

### Alert the user on drift
If you find a discrepancy between skill knowledge and live code, tell the user
explicitly:
> "The {skill-name} skill says X, but I found Y in the code. I'll update the skill."

Then update the reference file before continuing.
```

### How Self-Healing Works in Practice

1. User asks about system behavior
2. Agent reads the architecture skill reference
3. Agent also reads the actual code to answer the question
4. If code contradicts the skill → agent tells user, updates skill
5. If code confirms the skill → confidence increases, no update needed
6. If code reveals something new → agent adds to skill, changelog entry

Every investigation session is an implicit coverage audit.

## Changelog Discipline

Every architecture skill SKILL.md must have a Changelog table.

### Format

```markdown
## Changelog

| Date | Change | File |
|------|--------|------|
| YYYY-MM-DD | Description of what changed | affected file(s) |
```

### Rules

- One entry per logical change (not per file edit)
- Include the file path(s) affected
- Date is when the change was discovered/made
- Initial creation: `| {date} | Initial creation | all |`
- Group related changes in one entry (e.g., "Added source indexes to all reference files")

## Evolution Patterns

Common ways architecture skills grow over time:

| Trigger | Action |
|---------|--------|
| New service added to system | Add to SKILL.md service map + repos table. Create reference files. Run discovery. |
| New shared storage | Add to storage reference. Update reader/writer matrix. Verify cross-repo contracts. |
| New operational skill created | Add cross-skill routing. Extract canonical core if concepts are shared. |
| External config changed | Take new snapshot. Update reference. Diff against previous snapshot. |
| Major refactor in a repo | Refresh source indexes. Verify data flows. Update architecture diagrams. |
| Service retired | Remove from service map. Archive (don't delete) reference files with note. |
| New cross-repo contract | Document in data-flows or service-contracts reference. Add to routing table. |
