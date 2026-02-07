---
name: skill-builder
description: Interactive guide for creating and maintaining high-quality Claude skills. Use when user asks to create, build, update, or refactor a skill. Treats skill content like code - no duplication, elegant structure, complete but concise.
---

# Skill Builder

Build and maintain high-quality, efficient Claude skills.

## Goal

Every skill created or updated should make future work in that domain **smooth, efficient, and bug-free**.

## Core Principles

**Context window is a public good** - shared with system prompt, conversation, other skills. Only add what Claude doesn't already know.

**Treat skills like code refactoring:**

1. **No duplication** - Consolidate repeated information into single sources of truth
2. **Elegant structure** - Organize for readability and maintainability
3. **Complete but concise** - Capture all details without verbosity
4. **Incremental improvement** - Each update should make the skill better, not just bigger

---

## Three-Level Loading System

Claude loads skill content progressively:

| Level | Content                           | When Loaded         | Size Target |
| ----- | --------------------------------- | ------------------- | ----------- |
| 1     | **Metadata** (name + description) | Always in context   | ~100 words  |
| 2     | **SKILL.md body**                 | When skill triggers | <500 lines  |
| 3     | **Bundled resources**             | On-demand by Claude | Unlimited   |

**Key constraint:** SKILL.md body should stay under 500 lines. Move details to reference files.

---

## Degrees of Freedom

Choose instruction specificity based on task nature:

| Freedom    | Format                       | Use When                                     |
| ---------- | ---------------------------- | -------------------------------------------- |
| **High**   | Text instructions            | Multiple valid approaches, context-dependent |
| **Medium** | Pseudocode with parameters   | Preferred pattern exists, some variation OK  |
| **Low**    | Specific scripts, few params | Operations are fragile, consistency critical |

---

## Skill Structure

### Minimal

```
skill-name/
└── SKILL.md
```

### Standard (recommended)

```
skill-name/
├── SKILL.md              # Entry point (<500 lines)
├── references/           # Documentation loaded on-demand
│   ├── api.md
│   └── examples.md
├── scripts/              # Executable code (Python/Bash)
│   └── helper.py
└── assets/               # Files for output (templates, images)
    └── template.html
```

### Resource Types

| Folder        | Purpose                                      | When to Use                                           |
| ------------- | -------------------------------------------- | ----------------------------------------------------- |
| `references/` | Documentation Claude reads as needed         | Schemas, API docs, domain knowledge, policies         |
| `scripts/`    | Executable code                              | Repeated operations, deterministic reliability needed |
| `assets/`     | Files used in output (not read into context) | Templates, images, boilerplate code                   |

**Rules:**

- Keep references **one level deep** from SKILL.md
- For files >100 lines, include table of contents at top
- Scripts must be **tested** before finalizing

---

## Progressive Disclosure Patterns

### Pattern 1: High-level guide with references

```markdown
# PDF Processing

## Quick Start

Extract text with pdfplumber: [code example]

## Advanced Features

- **Form filling**: See `references/forms.md`
- **API reference**: See `references/api.md`
```

### Pattern 2: Domain-specific organization

```
bigquery-skill/
├── SKILL.md (overview + navigation)
└── references/
    ├── finance.md    # Revenue, billing metrics
    ├── sales.md      # Opportunities, pipeline
    └── product.md    # API usage, features
```

User asks about sales → Claude only reads `sales.md`.

### Pattern 3: Conditional details

```markdown
# DOCX Processing

## Creating Documents

Use docx-js. See `references/docx-js.md`.

## Editing Documents

For simple edits, modify XML directly.

- **Tracked changes**: See `references/redlining.md`
- **OOXML details**: See `references/ooxml.md`
```

---

## Creating a New Skill

### Step 1: Gather Requirements

Ask:

1. What problem does this skill solve?
2. When should Claude invoke it?
3. What are the key workflows?
4. Personal (`~/.claude/skills/`) or project (`.claude/skills/`)?

### Step 2: Craft the Description

CRITICAL: Description is the **primary trigger mechanism**. Include all "when to use" info here.

Format: `[What it does]. [When to use - specific triggers]. [Key capabilities].`

```
GOOD: "Comprehensive document editing with tracked changes. Use when working
      with .docx files for: creating documents, modifying content, adding
      comments, or tracking changes."

BAD:  "Helper for code" (too vague, no triggers)
```

### Step 3: Create Structure

Decide based on complexity:

- Simple skill → SKILL.md only
- Multiple workflows → Add `references/` folder
- Reusable code → Add `scripts/` folder
- Output templates → Add `assets/` folder

### Step 4: Write SKILL.md

```yaml
---
name: skill-name
description: [What]. [When - specific triggers]. [Capabilities].
---

# Skill Name

## Quick Start
[Most common workflow - brief]

## Workflows
### Workflow 1
1. Step one
2. Step two

## Additional Resources
- API details: `references/api.md`
- Examples: `references/examples.md`
```

---

## Updating an Existing Skill

### Quality-First Approach

When adding new knowledge, **refactor like code**:

1. **Read existing content first** - Understand current structure
2. **Identify duplication** - Don't add what already exists
3. **Consolidate** - Merge new info into existing sections
4. **Use tables for lists** - Convert repetitive patterns to tables
5. **Keep single source of truth** - One place for each piece of info

---

## Refactoring Skills

Apply **high cohesion, low coupling** to skills that grow too large or mix concerns.

### Separation of Concerns

| Concern | Keep In | Move To Reference |
|---------|---------|-------------------|
| **Workflows** (what to do) | SKILL.md | - |
| **Config** (CLI, defaults, paths) | - | `references/config.md` |
| **Formats** (templates, output specs) | - | `references/formats.md` |
| **Orchestration** (sub-agent rules, batching) | - | `references/orchestration.md` |
| **Domain knowledge** (API docs, schemas) | - | `references/<domain>.md` |

### SKILL.md Should Only Contain

- Quick reference / routing table
- Step-by-step workflows
- Pointers to references
- Troubleshooting (brief)

### When to Refactor

| Signal | Action |
|--------|--------|
| SKILL.md > 200 lines | Consider extraction |
| SKILL.md > 400 lines | Must extract |
| Config mixed with workflow | Extract to `config.md` |
| Format specs in workflow | Extract to `formats.md` |
| Same info in multiple sections | Consolidate to single reference |

### Refactoring Example

**Before (450 lines, mixed concerns):**
```
SKILL.md
├── Workflows (good)
├── CLI commands & defaults (config)
├── Directory structure (config)
├── Output templates (format)
├── Sub-agent batch rules (orchestration)
├── Math formatting rules (format)
└── Quality checklist (format)
```

**After (85 lines, single concern):**
```
SKILL.md           # Workflows only
references/
├── config.md      # CLI, defaults, paths
├── formats.md     # Templates, quality rules
└── orchestration.md  # Sub-agent rules
```

### Anti-Patterns to Avoid

| Anti-Pattern             | Problem                          | Solution                               |
| ------------------------ | -------------------------------- | -------------------------------------- |
| Session logs             | Grows unbounded, duplicates info | Extract facts into structured sections |
| Repeated tables          | Same info in multiple places     | Single reference table                 |
| Date-prefixed sections   | Focus on when, not what          | Organize by topic, not timeline        |
| Verbose explanations     | Hard to scan                     | Tables, bullets, code blocks           |
| Deeply nested references | Hard to navigate                 | Keep one level deep                    |

### Example: Consolidating Session Knowledge

**Before (duplicative):**

```markdown
## Session 2026-01-06

- API: py uses port 5000, go uses 6789
- Test command: curl ...

## Session 2026-01-07

- API: py uses port 5000, go uses 6789 # DUPLICATE
```

**After (consolidated):**

```markdown
## API Reference

| Service    | Port |
| ---------- | ---- |
| py-sandbox | 5000 |
| go-sandbox | 6789 |

## Verified URLs

| Date       | URL         | Result |
| ---------- | ----------- | ------ |
| 2026-01-06 | example.com | OK     |
| 2026-01-07 | other.com   | OK     |
```

---

## Naming Rules

- Lowercase, hyphens (not underscores)
- Max 64 characters
- Descriptive: `playwright-debugging`, `api-doc-generator`

---

## Quality Checklist

Before finalizing:

- [ ] SKILL.md under 500 lines (ideally <200 for workflows-only)
- [ ] No duplicated information
- [ ] Tables used for repetitive data
- [ ] Single source of truth for each fact
- [ ] Description states WHAT and WHEN (all triggers)
- [ ] References one level deep, explicitly linked
- [ ] Scripts tested and working
- [ ] Workflows are actionable and complete
- [ ] Config, formats, orchestration extracted to references (if >200 lines)

---

## Troubleshooting

| Problem                    | Cause                  | Solution                                  |
| -------------------------- | ---------------------- | ----------------------------------------- |
| Skill not activating       | Description too narrow | Add trigger phrases to description        |
| Skill activating too often | Description too broad  | Be more specific about triggers           |
| Changes not reflected      | Session cache          | Restart Claude Code                       |
| File not read              | Not referenced         | Add explicit read instruction in SKILL.md |
| SKILL.md too long          | Too much detail        | Move to `references/` files               |
