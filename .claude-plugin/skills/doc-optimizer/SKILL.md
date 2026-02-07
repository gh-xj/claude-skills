---
name: doc-optimizer
description: Optimizes markdown documentation for token efficiency and clarity. Use when user asks to clean up, refactor, consolidate, or optimize docs/notes. Reduces duplication, improves structure.
---

# Doc Optimizer

Refactor documentation like code: no duplication, elegant structure, token-efficient.

## When to Use

- User says: "cleanup docs", "reduce duplicates", "optimize files", "refactor notes"
- Multiple files with overlapping content
- Verbose or poorly structured markdown

## Workflow

### 1. Audit

```bash
find . -name "*.md" -not -path "*/.git/*" | head -20
```

Read each file, identify:
- Duplicated content across files
- Verbose sections that could be tables
- Files that could be merged
- Outdated/stale information

### 2. Plan New Structure

Aim for:
| File Type | Purpose |
|-----------|---------|
| README.md | Overview, status, action items |
| data.md | Reference data (tables) |
| tips.md | How-to, gotchas |

### 3. Refactor

**Principles:**

| Before | After |
|--------|-------|
| Repeated info in multiple files | Single source of truth |
| Bullet lists with same structure | Tables |
| Verbose explanations | Concise bullets |
| Date-prefixed sections | Topic-organized |
| Multiple small files | Consolidated by purpose |

### 4. Delete Redundant Files

```bash
rm old-file.md
rmdir empty-dir/
```

## Token-Efficient Patterns

### Use Tables for Structured Data

```markdown
# Bad (verbose)
- **Name:** John
- **Email:** john@example.com
- **Phone:** 555-1234

# Good (compact)
| Field | Value |
|-------|-------|
| Name | John |
| Email | john@example.com |
| Phone | 555-1234 |
```

### Combine Related Info

```markdown
# Bad (2 files)
## file1.md: Passport info
## file2.md: Visa info

# Good (1 file)
## data.md: All identity/travel data
```

### Remove Redundancy

```markdown
# Bad (repeated)
## Section A
Contact: support@example.com
## Section B  
Contact: support@example.com

# Good (single reference)
## Contacts
| Purpose | Email |
|---------|-------|
| Support | support@example.com |
```

## Quality Checklist

- [ ] No duplicated information
- [ ] Tables for repetitive data
- [ ] Single source of truth per fact
- [ ] Files organized by purpose, not timeline
- [ ] README is scannable overview
- [ ] Removed empty directories
