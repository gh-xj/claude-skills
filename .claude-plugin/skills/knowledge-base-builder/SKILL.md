---
name: knowledge-base-builder
description: Build and maintain token-efficient AI knowledge bases from source materials. Use when processing books, PDFs, documentation, or any content into structured, AI-agent-friendly markdown files. Handles PDF conversion, chapter splitting, summarization, token optimization, and README generation.
---

# Knowledge Base Builder

Transform source materials into structured, token-optimized knowledge bases for Claude skills.

## Quick Reference

| Input Type | Workflow |
|------------|----------|
| PDF book | [PDF Workflow](#pdf-workflow) |
| EPUB book | Convert to MD first, then [Content Workflow](#content-workflow) |
| Markdown files | [Content Workflow](#content-workflow) |
| Mixed sources | Process each type, merge into single KB |

## PDF Workflow

### Step 1: Convert PDF to Text

```bash
pdftotext -layout "input.pdf" "output.txt"
```

### Step 2: Analyze & Split

Identify structure patterns:

```bash
# Find chapter/section markers
grep -nE "^(Chapter|Part|Section|\d+\.)" output.txt | head -40
```

Split approach depends on source:
- **Books**: Chapter headers, part divisions
- **Documentation**: Section headers, API modules
- **Research papers**: Abstract, sections, references
- **Mixed content**: Topic boundaries, logical breaks

Output to `chapters/` or `sections/` directory.

### Step 3: Generate Summaries (Parallel)

Use sub-agents in batches of 3:

```
Task agent → Read section → Write summary
```

Adapt summary structure to content type. See `references/formats.md`

### Step 4: Optimize for Tokens (Parallel)

Use sub-agents to optimize each file:

```
Task agent → Read source → Write optimized to knowledge-base/<category>/
```

Optimization guidelines: `references/optimization.md`

### Step 5: Create README

Generate README with:
- Structure explanation
- TOC with links
- Navigation guide
- Changelog

---

## Content Workflow

For already-split content:

1. **Analyze** → Understand content type and structure
2. **Organize** → Create appropriate category directories
3. **Optimize** → Run optimization on each file
4. **Cross-reference** → Add links between related files
5. **Document** → Generate README

---

## Directory Structure

Adapt to content type:

```
knowledge-base/
├── README.md           # Entry point, TOC, changelog
├── <category-1>/       # Primary content
├── <category-2>/       # Secondary/supporting
└── <category-3>/       # Reference material
```

**Examples:**
- Book: `patterns/`, `appendices/`, `frameworks/`
- API docs: `endpoints/`, `models/`, `examples/`
- Course: `lectures/`, `labs/`, `readings/`

**Naming:** `##-kebab-case.md` or `letter-kebab-case.md`

---

## File Format

Flexible structure based on content type. Core elements:

```markdown
---
tags: [searchable, terms]
related: [cross-refs.md]
---

# Title

## Overview
[Brief description]

## [Main Content Sections]
[Adapt to material - concepts, steps, reference, etc.]

## Key Points
[Takeaways]
```

Metadata header is recommended but optional for simple content.

---

## Parallel Processing

| Source Size | Batch Size | Agents |
|-------------|------------|--------|
| < 10 files | 3/agent | 3 |
| 10-30 files | 3/agent | 5-7 |
| > 30 files | 3/agent | 8-10 |

Sub-agents return "Done" to minimize tokens.

---

## References

- Optimization: `references/optimization.md`
- Format examples: `references/formats.md`
