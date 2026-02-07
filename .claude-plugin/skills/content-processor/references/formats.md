# Output Formats

## Summary Depth

| Content Type   | Depth    | Focus                         |
| -------------- | -------- | ----------------------------- |
| Blog/article   | Light    | Main thesis, takeaways        |
| YouTube video  | Medium   | Key points, main arguments    |
| Technical book | **Deep** | Structure, logic, terminology |

### Deep Summaries (books)

Preserve:

- **Logical flow** - How arguments build
- **Section structure** - Follow original organization
- **Key quotes** - Author's words for important claims
- **Relationships** - Tables/diagrams where helpful

---

## Format Templates

### Standard Summary

```markdown
---
source_type: youtube | website | epub | markdown
source: <url or file path>
language: <code>
date_processed: YYYY-MM-DD
reading_time_min: <calculated>
---

# Title

## Key Points

- Point 1
- Point 2

## Summary

<detailed content>
```

### Enriched Notes (technical content)

```markdown
# Topic

## Central Question

> [Motivating question, e.g., "Can we beat O(log n)?"]

## Quick Reference

| Operation | Time | Space |
| --------- | ---- | ----- |

## Key Formulas

| Formula | Meaning |
| ------- | ------- |

---

## The Journey

### Part 1: [First Concept]

[Problem → failed attempt → solution → worked example]

<details>
<summary>Detailed walkthrough</summary>
[Step-by-step with concrete numbers, ASCII diagrams]
</details>

---

## Flow Diagram

[ASCII art: Problem → Attempt 1 → Solution]

## Review Questions

1. ...
<details><summary>Answers</summary>...</details>
```

### Light Format (simple content)

```markdown
# Title

## Key Points

- Point 1
- Point 2

## Summary

[Content]
```

### Twitter-Style

Trigger: "twitter format", "tweet style", "key points"

```markdown
# Title - Key Takeaways

## Main Points

- [Point under 280 chars] 🎯
- [Point under 280 chars] 💡

## One-Liner

> [Tweetable summary]

#Tag1 #Tag2 #Tag3
```

---

## Quality Checklist

| Check       | Question                                           |
| ----------- | -------------------------------------------------- |
| Formulas    | Every formula has a worked example with numbers?   |
| Variables   | All symbols defined before use?                    |
| Comparisons | Tables show concrete differences, not just labels? |
| Derivations | "Why" explained, not just "what"?                  |
| Complexity  | Shows path: recurrence → tree → sum → result?      |

**Common failure**: Stating conclusions without showing the journey.

---

## Math Formatting

### Delimiters

| Type    | Syntax                     | Example       |
| ------- | -------------------------- | ------------- |
| Inline  | `$...$`                    | `$x^2 + y^2$` |
| Display | `$$...$$` with blank lines | See below     |

### Display Math

```markdown
The sum formula is:

$$
\sum_{i=1}^{n} i = \frac{n(n+1)}{2}
$$

This can be proven by induction.
```

### Common Mistakes

| Mistake             | Wrong                 | Correct                  |
| ------------------- | --------------------- | ------------------------ |
| Escaped underscores | `$x\_1$`              | `$x_1$`                  |
| No blank lines      | `text$$formula$$text` | blank lines around `$$`  |
| Backslash in text   | `O(n\log n)`          | `$O(n \log n)$`          |
| Mixed delimiters    | `\(x\)` and `$y$`     | Use `$...$` consistently |
| Bare operators      | `sum`, `log`          | `\sum`, `\log`           |

### Math Prompt Addition

For sub-agents processing math content:

```
Math formatting:
- Inline: $...$ (no escaping underscores)
- Display: $$ on own line, blank lines before/after
- Operators: \sum, \log, \lim, \frac{}{}, \sqrt{}
- Define all variables before use
```
