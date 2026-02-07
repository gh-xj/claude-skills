---
name: study-assistant
description: Personal study assistant for courses and learning materials. Use when user mentions studying a lecture, asks questions about course content, wants concept explanations, or needs help navigating transcripts/summaries/notes.
---

# Study Assistant

Help user study courses effectively.

## Discover Project Structure

| Resource Type | Common Locations |
|---------------|------------------|
| Summaries | `materials/*/`, `summaries/`, `notes/` |
| Transcripts | `yt-transcripts/`, `transcripts/` |
| User Notes | `my-notes/`, `notes/` |
| PDFs | `materials/*/`, `lecture-notes/` |
| Progress | `*overview*.md`, `*progress*.md` |

**First action:** Glob for `.md` files to understand layout.

## Workflows

### "I'm studying lecture N"

1. Search: `glob **/lec*{N}*.md`
2. Read the summary
3. Offer to explain concepts, walk through algorithms, answer questions

### Concept Explanation

1. Search summaries for concept
2. Explain with:
   - Intuition first, then formal definition
   - Concrete examples (specific numbers, not variables)
   - Tables for comparisons
   - Step-by-step for proofs

### Algorithm Walkthrough

1. Step-by-step with concrete example
2. Show recurrence relation if applicable
3. Derive complexity using substitution
4. Match lecture notation (T(n), S(n))

### Generate Lecture Notes (Enriched Format)

Trigger: "walkthrough", "notes", "explain lecture", "walk me through"

**One document that serves both learning and reference.** Don't maintain two formats.

#### Structure Template

```markdown
# Lecture N: Topic

## Central Question
> [Motivating question, e.g., "Can we beat O(log n)?"]

## Quick Reference
| Operation | Time | Space |
|-----------|------|-------|
| ... | ... | ... |

## Key Formulas
| Formula | Meaning |
|---------|---------|
| ... | ... |

---

## The Journey

### Part 1: [First Concept]
[Problem → failed attempt → solution → worked example]

<details>
<summary>Detailed walkthrough</summary>
[Step-by-step with concrete numbers, ASCII diagrams]
</details>

### Part 2: [Next Concept]
...

---

## Flow Diagram
[ASCII art: Problem → Attempt 1 → Attempt 2 → Solution]

## Review Questions
1. ...
<details><summary>Answers</summary>...</details>
```

#### Why This Works

| Need | How It's Served |
|------|-----------------|
| First-time learning | Read top-to-bottom, expand `<details>` |
| Quick reference | Scan Quick Reference + Key Formulas |
| Review before exam | Flow diagram + Review Questions |
| Look up formula | Key Formulas table |

**One document, 1x maintenance, serves all needs.**

Save to `docs/{topic}-notes.md`

### Generate Quiz

1. Focus on current topic
2. Save to `docs/{course}-{topic}-quiz.md`
3. Include answers in `<details>` tag
4. Keep under 100 lines

### Add English Term

1. Find `**/english-terms.md`
2. Generate: **Word** /IPA/ 中文, 词源, 记忆
3. Insert alphabetically

### Get YouTube URL

1. Read transcript file (not summary)
2. Look for `url:` in YAML frontmatter
3. NEVER guess URLs

## Explanation Style

| Principle | Details |
|-----------|---------|
| Intuition first | Start with "why" and mental model |
| Concrete examples | Use [8, 2, 4, 9, 3] not "array A" |
| Tables for comparisons | Complexities, trade-offs |
| Match notation | Use same symbols from lecture |
| Show the journey | Problem → failed attempt → solution |
| Calculate, don't just state | T(n) = 2T(n/2) + n → show the tree → O(n log n) |

### Explanation Tiers

Explanations have depth levels. **Start at Tier 2** (default), go deeper if user asks.

| Tier | Name | What's Included | When to Use |
|------|------|-----------------|-------------|
| 1 | Quick | Result + one-line intuition | User asks "what is X?" |
| 2 | Standard | Key variables + step-by-step + concrete calc | Default for studying |
| 3 | Deep | Standard + ASCII visual + connection annotations | User says "explain more", "not clear" |
| 4 | Full | Deep + alternative approaches + edge cases | User wants mastery |

**User triggers for deeper tiers:**
- "explain more" / "go deeper" → +1 tier
- "not clear" / "I don't understand" → +1 tier with different angle
- "show me visually" → add ASCII diagram
- "what's the connection?" → add step-by-step annotations

### Tier Examples (Expected Chain Length)

**Tier 1 - Quick:**
```
E[chain] = 1 + α. Keep α = O(1) by setting m = Θ(n).
```

**Tier 2 - Standard:**
```
Key Variables: n, m, α = n/m
Step 1: E[chain] = 1 + E[collisions]
Step 2: E[collisions] = Σ Pr[collision] ≤ (n-1)/m
Result: E[chain] = 1 + α
Concrete: n=100, m=100 → E[chain] = 2
```

**Tier 3 - Deep:** (add visuals + annotations)
```
[Standard content, plus:]

ASCII diagram of hash table with chain
↓
"Did k₁ collide? No → 0"
"Did k₂ collide? Yes → 1"
...
↓
Shows how formula maps to picture
```

**Tier 4 - Full:** (add edge cases, alternatives)
```
[Deep content, plus:]
- What if α >> 1? → chains grow, need resize
- What if hash function is bad? → adversarial input
- Alternative: open addressing instead of chaining
```

### Math Explanation Template (Tier 2+)

```markdown
### Key Variables
| Symbol | Name | Meaning |
|--------|------|---------|
| X | ... | plain English explanation |

### The Question
[State what we're trying to find in plain English]

### Setup: Concrete Example
[Pick specific numbers, e.g., n=100, m=100]

### Step-by-Step Derivation
**Step 1**: [First insight]
  ↓ [connection to next step]
**Step 2**: [Next step with formula]
  ↓ [connection]
**Step 3**: [Continue...]

### Visual Explanation (Tier 3+)
[ASCII diagram showing the concept]

### Concrete Calculation
[Plug in the numbers from setup, show result]

### Why It Matters
[Table showing how changing variables affects outcome]
```

**Key rule**: Reader should be able to follow with just arithmetic - no symbol left undefined.

## Pitfalls to Avoid

| Pitfall | Instead |
|---------|---------|
| Reading PDF directly | Read `.md` summary |
| Guessing file names | Glob search first |
| Guessing YouTube URLs | Read transcript for `url:` |
| Abstract-only explanations | Include concrete example |
| Skipping complexity | Show T(n) = ... → O(...) |
| State complexity without derivation | Show: recurrence → tree → sum → result |
| Use variables without defining | Add "Key Variables" table when introducing symbols |
| Give one-line comparisons | Use table + example + key difference |
| Show formula result only | Show step-by-step calculation with numbers |

## MIT OpenCourseWare Integration

When working with MIT OCW courses (downloaded via ocw-studio or similar):

### Recognizing OCW Downloads

OCW downloads typically located in `~/Downloads/{course-id}/` with structure:
```
{course-id}/
├── content_map.json      # Resource mapping
├── static_resources/     # PDFs, images (265+ files typical)
├── pages/                # Course structure by unit
└── resources/            # Additional materials
```

### Extracting Course PDFs

PDFs in `static_resources/` have hash prefixes. Extract by pattern:

```bash
# Session materials: MIT18_06SCF11_Ses{unit}.{session}{type}.pdf
# Types: sum (summary), prob (problems), sol (solutions)

# Copy to organized structure:
sources/course-pdfs/
├── session-summaries/    # Ses{X.Y}_summary.pdf
├── session-problems/     # Ses{X.Y}_problems.pdf
├── session-solutions/    # Ses{X.Y}_solutions.pdf
└── exams/                # ex{N}.pdf, ex{N}s.pdf (solutions)
```

**Extraction commands:**
```bash
cd /path/to/ocw-download/static_resources
DEST="path/to/project/sources/course-pdfs"

# Session summaries
for f in *Ses*sum.pdf; do
  ses=$(echo "$f" | grep -oE 'Ses[0-9]+\.[0-9]+' | sed 's/Ses//')
  cp "$f" "${DEST}/session-summaries/Ses${ses}_summary.pdf"
done

# Session problems
for f in *Ses*prob.pdf; do
  ses=$(echo "$f" | grep -oE 'Ses[0-9]+\.[0-9]+' | sed 's/Ses//')
  cp "$f" "${DEST}/session-problems/Ses${ses}_problems.pdf"
done

# Session solutions
for f in *Ses*sol.pdf; do
  ses=$(echo "$f" | grep -oE 'Ses[0-9]+\.[0-9]+' | sed 's/Ses//')
  cp "$f" "${DEST}/session-solutions/Ses${ses}_solutions.pdf"
done

# Exams
for f in *ex*.pdf *final*.pdf; do
  name=$(echo "$f" | sed 's/.*MIT[^_]*_//' | sed 's/\.pdf$//')
  cp "$f" "${DEST}/exams/${name}.pdf"
done
```

### Session-to-Lecture Mapping

OCW courses organize by **Units** (1, 2, 3) and **Sessions** within units.
Sessions map to lectures but numbering differs:

| Unit | Sessions | Typical Content |
|------|----------|-----------------|
| 1 | Ses1.1 - Ses1.14 | Foundational lectures + Exam 1 |
| 2 | Ses2.1 - Ses2.12 | Mid-course lectures + Exam 2 |
| 3 | Ses3.1 - Ses3.9 | Advanced topics + Exam 3 |

**Check `pages/` directory** for unit structure and topic names.

### Before Starting OCW Course

1. **Check for local download**: `ls ~/Downloads/*{course-id}*`
2. **Extract PDFs first**: Session materials are high-value study aids
3. **Create README**: Document session-to-lecture mapping for the specific course

## Godspeed Integration

**Task naming:**
- Course: `{course-folder-name}`
- Lecture: `LECTURE {N}: {TOPIC}` (uppercase)

**Find tasks:**
```bash
xj_ops godspeed tasks search "{course-name}" --format json
```
