---
tags: [bootstrap, scaffold, new-system, setup]
related: [discovery-playbook.md, templates.md, maintenance.md]
---

# Bootstrap Architecture Skill

Step-by-step workflow for creating an architecture skill for a new system.

## Prerequisites

- System scope questionnaire completed (from SKILL.md)
- Read access to all repos in the system
- Understanding of which services are core vs peripheral
- `skill-builder` skill available (for structure rules and quality validation)

## Step 1: Classify System Complexity

| Repos | Services | Shared Storage | Complexity | Recommended Structure |
|-------|----------|----------------|------------|----------------------|
| 1-2 | 1-3 | 0-1 | Small | Flat `references/` |
| 3-5 | 3-8 | 2-4 | Medium | `references/core/` + `references/architecture/` |
| 6+ | 8+ | 4+ | Large | `core/` + `architecture/` + per-service subdirs |

**Small systems** often don't need a dedicated architecture skill. A single reference file in an existing skill may suffice.

**Medium systems** benefit from the canonical core pattern: shared concepts in `core/`, system-specific details in `architecture/`.

**Large systems** may need multiple architecture skills (one global, one per major subsystem) with cross-skill routing.

## Step 2: Choose Skill Location

| Location | When to Use |
|----------|-------------|
| `~/.claude/skills/{name}/` (global) | System spans multiple repos, used from any working directory |
| `{repo}/.claude/skills/{name}/` (repo-local) | System is contained within one repo or has a primary repo |

For multi-repo systems, global is usually correct. The skill describes cross-repo relationships that don't belong to any single repo.

## Step 3: Create Scaffold

Create the directory structure based on complexity class:

```
# Medium system (most common)
mkdir -p ~/.claude/skills/{system-name}/references/{core,architecture}

# Small system
mkdir -p ~/.claude/skills/{system-name}/references

# Large system
mkdir -p ~/.claude/skills/{system-name}/references/{core,architecture,services}
```

Generate initial SKILL.md using the template from `references/templates.md` > Generated Architecture SKILL.md.

Fill in immediately:
- Metadata (name, description with trigger phrases)
- Service Map (ASCII diagram — even a rough one helps)
- Repos table (name, path, role)
- Empty routing table (headers only, populated after discovery)

Leave for later:
- Critical Design Decisions (discovered during exploration)
- Changelog (add initial entry after first pass)

## Step 4: Plan Discovery

Map repos to discovery tasks. Priority order:

1. **Orchestrator/entry-point service** — Start here. It reveals the overall flow.
2. **Shared storage** — Map who reads/writes what.
3. **Core processing services** — Business logic, evaluators, renderers.
4. **Gateway/proxy services** — External interfaces, BFF layers.
5. **Peripheral services** — Utilities, monitoring, admin tools.

Parallelization guide:

| System Size | Explore Agents | Strategy |
|-------------|----------------|----------|
| Small (1-2 repos) | 1-2 sequential | Single agent per repo |
| Medium (3-5 repos) | 2-3 parallel | Orchestrator first, then parallel for others |
| Large (6+ repos) | 3-5 parallel batches | Orchestrator + storage first, then service batches |

See `references/discovery-playbook.md` for specific exploration techniques and agent prompts.

## Step 5: Run Discovery and Populate References

For each discovery task:
1. Launch Explore agent with focused prompt (see discovery-playbook.md)
2. Capture findings into the appropriate reference file
3. Use Reference File template from `references/templates.md`
4. Add Source Files table immediately (don't defer — it's harder to add later)

After each reference file is populated:
- Add a routing table entry in SKILL.md pointing to it
- Cross-link with `related:` frontmatter to sibling references

## Step 6: Add Critical Design Decisions

After discovery, identify the 3-5 most important architectural decisions:
- How do services communicate? (RPC, shared storage, events, or combination)
- What's the deployment model? (independent, coordinated, monorepo)
- What's the contract between services? (struct fields, storage keys, IDL types)
- Are there external config systems that affect behavior? (rule engines, feature flags)

Write each as a bullet in the SKILL.md Critical Design Decisions section. Focus on WHY, not just WHAT. These help agents make safe cross-repo changes.

## Step 7: Finalize

1. Add self-healing section — copy the Knowledge Freshness template from `references/maintenance.md` > Drift Detection
2. Add initial changelog entry: `| {date} | Initial creation | all |`
3. Validate using `skill-builder` > Quality Checklist:
   - SKILL.md under 500 lines (ideally under 200)
   - No duplicated information across reference files
   - All `references/` paths in routing table resolve to existing files
   - Description covers trigger phrases for common user questions
4. Verify source file indexes: do the file:line entries match actual code?

## Common Pitfalls

| Pitfall | How to Avoid |
|---------|-------------|
| Trying to document everything in one pass | Focus on core flows first. Add depth iteratively. |
| Copying code comments as architecture knowledge | Architecture knowledge is RELATIONSHIPS between components, not implementation details. |
| Huge monolithic reference files | Split by concern. If a file covers both storage and data flow, split it. |
| Skipping source file indexes | Add them during discovery. Backfilling later requires re-reading all the code. |
| No self-healing section | Every architecture skill must include drift detection instructions. |
