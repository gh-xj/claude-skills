---
name: architecture-as-code
description: Build and maintain architecture knowledge bases (Claude skills) for complex multi-repo or multi-service systems. Use when bootstrapping architecture documentation for a new system, adding cross-skill routing between architecture and operational skills, running coverage audits on existing architecture skills, refreshing source file indexes, or capturing external config snapshots. Triggers on "architecture skill", "system knowledge base", "cross-repo architecture", "source file index", "architecture coverage", "map the system".
---

# Architecture as Code

Turn living codebases into self-healing architecture knowledge bases.

## Dependencies

- **skill-builder**: Skill structure rules, quality checklist, naming, progressive disclosure.
- **skill-builder > references/multi-skill-architecture.md**: Canonical core pattern, source file indexes, versioned detection config, self-healing instructions.

This skill adds architecture-specific methodology on TOP of skill-builder's structural foundation. Do not duplicate skill-builder content.

## System Scope Questionnaire

Before starting any workflow, establish scope:

| Question | Purpose |
|----------|---------|
| How many repos/services? | Determines complexity class |
| What are the service roles? | Drives service map |
| Shared storage systems? | Identifies cross-repo contracts |
| External config systems? (rule engines, LLM prompts, feature flags) | Determines if config snapshots needed |
| Existing operational skills? (debugger, investigator) | Determines cross-skill integration needs |
| Target skill location? (`~/.claude/skills/` or `.claude/skills/`) | Global vs repo-local |

## Workflow 1: Bootstrap Architecture Skill

**When**: New system, no architecture skill exists yet.

| Step | Action | Reference |
|------|--------|-----------|
| 1 | Scope the system | Questionnaire above |
| 2 | Classify complexity & choose structure | `references/bootstrap.md` |
| 3 | Create scaffold (dirs + initial SKILL.md) | `references/bootstrap.md` + `references/templates.md` |
| 4 | Run discovery (parallel Explore agents) | `references/discovery-playbook.md` |
| 5 | Populate reference files | `references/templates.md` > Reference File Template |
| 6 | Add source file indexes | `references/templates.md` > Source Index Template |
| 7 | Add self-healing section | `references/maintenance.md` > Drift Detection |
| 8 | Validate | `skill-builder` > Quality Checklist |

## Workflow 2: Add Cross-Skill Integration

**When**: Architecture skill exists, need to connect to operational skills.

| Step | Action | Reference |
|------|--------|-----------|
| 1 | Identify operational skills | `references/cross-skill-integration.md` |
| 2 | Determine upstream/downstream | `references/cross-skill-integration.md` > Upstream/Downstream |
| 3 | Extract canonical core (if shared concepts) | `references/cross-skill-integration.md` > Canonical Core |
| 4 | Create routing tables | `references/templates.md` > Routing Table Template |
| 5 | Wire back-references in downstream skills | `references/cross-skill-integration.md` > Back-References |
| 6 | Validate cross-references | `references/cross-skill-integration.md` > Validation |

## Workflow 3: Maintain & Evolve

**When**: Architecture skill exists, codebase has changed.

| Step | Action | Reference |
|------|--------|-----------|
| 1 | Refresh source file indexes | `references/maintenance.md` > Source Index Refresh |
| 2 | Capture external config snapshots | `references/maintenance.md` > Config Snapshots |
| 3 | Detect and fix drift | `references/maintenance.md` > Drift Detection |
| 4 | Update changelog | `references/maintenance.md` > Changelog |

## Workflow 4: Assess Architecture Coverage

**When**: Need to verify completeness and accuracy.

| Step | Action | Reference |
|------|--------|-----------|
| 1 | Build coverage matrix | `references/coverage-audit.md` |
| 2 | Verify source file indexes | `references/coverage-audit.md` > Index Verification |
| 3 | Validate cross-references | `references/coverage-audit.md` > Cross-Reference Check |
| 4 | Generate gap report | `references/coverage-audit.md` > Gap Report |

## Out of Scope

| Task | Use Instead |
|------|-------------|
| Feature-level specs/plans/tasks | `speckit-autopilot` |
| Processing books/PDFs into knowledge bases | `knowledge-base-builder` |
| Generic skill creation (non-architecture) | `skill-builder` |
| Operational investigation/debugging | Generated operational skills (debugger, investigator) |

Architecture skills DESCRIBE systems. Operational skills USE that knowledge.

## Changelog

| Date | Change | File |
|------|--------|------|
| 2026-02-05 | Initial creation | all |
