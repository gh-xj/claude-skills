---
tags: [cross-skill, routing, canonical-core, multi-skill, dependency]
related: [bootstrap.md, maintenance.md, templates.md]
---

# Cross-Skill Integration

How to connect architecture skills to operational skills (debuggers, investigators, deployers) and establish canonical shared concepts.

## When to Integrate

| Signal | Action |
|--------|--------|
| Operational skill asks "how does the pipeline work?" | Architecture skill provides canonical pipeline reference |
| Debugger needs to understand data flow | Architecture skill provides data-flows reference |
| Investigator needs severity/status definitions | Architecture skill provides canonical definitions |
| Two skills explain the same concept differently | Extract to canonical core in upstream skill |
| New operational skill created for a covered system | Add cross-skill routing in both directions |

## Upstream/Downstream Identification

Architecture skills are UPSTREAM (they define concepts). Operational skills are DOWNSTREAM (they use concepts for investigation, debugging, deployment).

| If... | Then... |
|-------|---------|
| Skill A explains system concepts AND skill B uses them | A is upstream |
| Both explain the same concepts from different angles | Extract shared to `core/` in the more stable skill |
| Neither depends on the other | No cross-skill integration needed |

Rules:
1. **One-way only** — downstream depends on upstream, never circular
2. **Declare explicitly** — downstream SKILL.md states: "This skill depends on `{upstream}` for X, Y, Z"
3. **Upstream never requires downstream** — it may mention downstream exists (for routing) but does not depend on it

## Extracting Canonical Core

When a concept is shared between skills, extract it to `references/core/` in the upstream skill.

### Criteria for Core Extraction

A concept belongs in `core/` when:
- Referenced by 2+ skills
- Has a single correct definition (not opinion or workflow)
- Changes should propagate to all consumers automatically

### Process

1. Identify the shared concept (e.g., pipeline phases, severity scale, storage layout)
2. Determine the "source of truth" skill (usually the architecture skill)
3. Create `references/core/{concept}.md` in the upstream skill
4. Add `canonical: true` to frontmatter
5. Move the canonical definition there from wherever it currently lives
6. In all downstream skills, replace inline definitions with back-references:
   ```
   For pipeline phase definitions, see `{upstream-skill}` > `references/core/pipeline.md`
   ```

### Common Core Concepts

| Concept | Typical Content |
|---------|----------------|
| Pipeline phases | Stage names, order, skip conditions, entry/exit criteria |
| Severity/confidence scale | Level definitions, thresholds, what determines each level |
| Storage layout | Prefixes, collections, buckets with ownership and TTL |
| Service contracts | RPC fields, shared types, implicit agreements |

## Creating Routing Tables

Three types of routing to add (see `references/templates.md` > Routing Tables for formats):

### 1. Internal routing (architecture skill SKILL.md)

Routes user questions to the correct reference file within the architecture skill. Organize by concern:
- **Core concepts** — pipeline, severity, shared definitions
- **Architecture** — data flows, storage, service contracts, external config
- **Operations** — cross-skill pointers to debugger, investigator, deployer

### 2. Cross-skill outbound (architecture -> operations)

In the architecture skill's SKILL.md, add an Operations section routing to downstream skills:
```markdown
### Operations (cross-skill)

| Question | Skill / Reference |
|----------|-------------------|
| "Debug sandbox failures" | `sandbox-debugger` skill in `{repo}` repo |
| "Investigate scan events" | `urlscan-investigator` skill in `{repo}` repo |
```

### 3. Cross-skill inbound (operations -> architecture)

In each downstream skill's SKILL.md, add a dependency section:
```markdown
## Dependencies

- **{architecture-skill}**: Pipeline context, severity definitions, storage layout
  - `{architecture-skill}` > `references/core/pipeline.md` for pipeline phases
  - `{architecture-skill}` > `references/architecture/shared-storage.md` for storage layout
```

## Wiring Back-References

In each downstream skill:

1. Add dependency declaration at top of SKILL.md (after metadata)
2. For each concept from upstream, add explicit cross-reference
3. Use convention: `See {upstream-skill} > references/{path}.md > {Section Name}`
4. Remove any inline definitions that are now canonicalized upstream

## Validation Checklist

After wiring cross-skill integration:

- [ ] Upstream skill does NOT reference downstream skill's content (only mentions it exists)
- [ ] Downstream skill declares dependency explicitly in SKILL.md
- [ ] Every `>` cross-reference resolves to an existing file
- [ ] No concept is defined in both skills (single source of truth)
- [ ] Routing tables cover the most common user questions (test 5+ queries mentally)
- [ ] Architecture skill's Operations section lists all known downstream skills
- [ ] Downstream skills' Dependencies sections list the architecture skill
