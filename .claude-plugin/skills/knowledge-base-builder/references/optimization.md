# Token Optimization Guidelines

## Goal

Compress content for AI consumption while preserving:

- Core information
- Important details
- Actionable knowledge

## Transformation Principles

### Structure Compression

| Original          | Optimized                 |
| ----------------- | ------------------------- |
| Prose paragraphs  | Bullet points             |
| Repetitive lists  | Tables                    |
| Long explanations | Concise statements        |
| Repeated info     | Single source + cross-ref |

### Content Filtering

| Keep                     | Remove                 |
| ------------------------ | ---------------------- |
| Core concepts            | Filler words           |
| Key details              | Redundant explanations |
| Minimal code examples    | Verbose comments       |
| Definitions              | Marketing language     |
| 1-2 examples per concept | Excessive examples     |

### Abbreviations

Use domain-appropriate abbreviations. Common examples:

- Tech: LLM, RAG, MCP, API, CLI, DB, UI
- General: e.g., i.e., etc., vs.
- Navigation: → for "see" or "leads to"

### Code Handling

Adapt to content type:

- **API docs**: Keep signatures, minimal examples
- **Tutorials**: Condensed working examples
- **Reference**: Patterns over full implementations
- **Conceptual**: Pseudocode or omit entirely

## Cross-References

Link related content:

```markdown
→ See: related-file.md
→ Related: [concept](path/to/file.md)
```

## Metadata (Optional)

For structured knowledge bases:

```yaml
---
tags: [searchable, terms]
related: [cross-refs.md]
---
```

Skip for simple or informal content.

## Quality Check

- [ ] No unnecessary prose
- [ ] Tables for structured data
- [ ] Cross-references where helpful
- [ ] Core info preserved
- [ ] Filler removed
