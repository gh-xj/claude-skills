---
name: personal-evolution
description: Unified system for learning from pain, making decisions, and extracting principles. Use when user is stuck, frustrated, made a mistake, facing a decision, says "lesson learned", "I'm stuck", "help me decide", "what went wrong", or reflects on experiences. Combines Dalio's Principles with structured lesson capture.
---

# Personal Evolution

**Core Formula**: Pain + Reflection = Progress → Principle → Better Decisions

```
Pain/Mistake → Reflect → Extract Lesson → Generalize to Principle → Apply Next Time
     ↑                                                                    │
     └──────────────────── (evolution loop) ──────────────────────────────┘
```

---

## Mode Detection

When triggered, identify which mode:

| User Signal | Mode | Action |
|-------------|------|--------|
| Frustrated, stuck, failed | **Pain** | Reflection prompts → capture lesson |
| "Help me decide", weighing options | **Decision** | Decision framework → believability check |
| "Lesson learned", post-mortem | **Capture** | Record lesson → extract principle |
| Similar situation to past lesson | **Recall** | Surface relevant lessons + principles |
| "Review", "what have I learned" | **Analyze** | Cross-session patterns |
| "Quiz me", "make cards", after learning | **Quiz** | Extract quiz cards from insights |

---

## Mode 1: PAIN (Reflection)

When user is frustrated, stuck, or failed:

### Step 1: One Reflection Question

Pick ONE based on situation:

| Situation | Question |
|-----------|----------|
| Frustrated | "What is this pain trying to teach you?" |
| Failed | "Is this a you-problem, process-problem, or design-problem?" |
| Stuck | "Which step are you failing at: Goal, Problem ID, Diagnosis, Plan, or Execution?" |
| Conflict | "Are you trying to be right, or trying to find what's true?" |
| Repeated issue | "What's the root cause - a character trait or a system gap?" |

### Step 2: Prompt for Lesson

After reflection, ask:
> "Want to capture this as a lesson? (y/n)"

If yes → go to **Mode 3: Capture**

---

## Mode 2: DECISION

When user is weighing options or needs to decide:

### Quick Framework

Ask ONE:

| Check | Question |
|-------|----------|
| Expected value | "What's the upside × probability vs downside × probability?" |
| Believability | "Who's done this successfully 3+ times? Are you that person?" |
| Consequences | "What happens next? And then what? And then?" |
| Reversibility | "If wrong, can you recover? What's the real cost?" |

### When Stuck Between Options

> "Step outside yourself. If you were advising a friend in this exact situation, what would you say?"

### After Decision

> "What principle guided this decision? Write it down for next time."

---

## Mode 3: CAPTURE (Lesson Recording)

### Data Locations

| Type | Path |
|------|------|
| Personal (cross-project) | `~/.lessons/personal/{category}/` |
| Project-specific | `{project}/.lessons/entries/{category}/` |
| Session reflections | `{project}/.lessons/sessions/` |

### Lesson Template

```yaml
---
date: YYYY-MM-DD
category: decisions|execution|relationships|tools|process|money
severity: critical|important|minor
root_cause_type: character|process|design
tags: []
---

# {Title}

## What Happened
[Brief description]

## Root Cause Analysis
- **Proximate cause** (action): {what you did/didn't do}
- **Root cause** (trait/system): {why you did/didn't do it}
- **Type**: Character trait | Process gap | Design flaw

## Lesson
[The takeaway - what to remember]

## Principle Extracted
[One-line rule for similar situations - write as "When X, do Y"]

## Action
[Specific change to make]
```

### Before Writing, Show Preview:

```
📝 Lesson Preview:
─────────────────
File: {path}
Title: {title}
Root Cause: {type}
Principle: {one-liner}
─────────────────
Save? (y/n)
```

---

## Mode 4: RECALL (Proactive Surfacing)

When context matches a past lesson:

1. Search `~/.lessons/personal/` and `{project}/.lessons/entries/`
2. Match by category, tags, or keywords
3. Surface briefly:

```
💡 Related lesson from {date}:
"{principle extracted}"
→ Full lesson: {filepath}
```

**Triggers for recall:**
- User starts task in category with existing lessons
- User mentions topic/tool with past lessons
- User about to repeat a past mistake pattern

---

## Mode 5: ANALYZE (Review & Patterns)

When user says "review", "what have I learned", "analyze":

### Session Review

```markdown
## Sessions Summary

| Metric | Count |
|--------|-------|
| Sessions reviewed | N |
| Lessons captured | N |
| Principles extracted | N |

### Patterns Detected
- {pattern 1}: appeared N times
- {pattern 2}: appeared N times

### 5-Step Weakness Analysis
Based on your lessons, you most often fail at:
1. {step}: N occurrences
2. {step}: N occurrences

### Unpromoted Insights
Observations marked "generalizable" but not yet principles:
- [ ] {observation} → needs principle
```

### Cross-Lesson Patterns

Look for:
| Pattern | Signal |
|---------|--------|
| Repeated root cause | Same character trait causing issues |
| Category clusters | Many lessons in same area |
| Severity trends | Are critical lessons decreasing? |

---

## Session Tracking

### Auto-Activate on First Trigger

Create/continue session file at `{project}/.lessons/sessions/{date}_{context}.md`

### During Session

Silently collect observations when:
- User corrects a mistake
- User discovers something new
- Same issue appears 3+ times
- A workaround is found

At natural breakpoints, batch-ask:
> "📝 Noticed N observations. Add to session log? (y/n/select)"

### Session Wrap-up

When user says "wrap up" or session ends:

```markdown
## End of Session

### 5-Step Audit
Which step caused issues today?
- [ ] Goals (unclear what I wanted)
- [ ] Problems (tolerated too long)  
- [ ] Diagnosis (symptoms vs root cause)
- [ ] Design (bad plan)
- [ ] Execution (didn't follow through)

### Lessons to Promote
- [ ] {observation} → permanent lesson?

### Principles Reinforced
- {principle that helped today}

### Open Questions
- {unresolved}
```

---

## Mode 6: QUIZ (Knowledge Solidification)

**Purpose**: Convert "aha moments" into retrievable knowledge through quiz cards.

### When to Trigger

- After user has a breakthrough understanding
- User says "quiz me", "make cards", "what should I remember"
- End of learning session
- When user corrects a misconception they had

### Card Types

| Type | Use When | Format |
|------|----------|--------|
| **Misconception Buster** | User had wrong intuition | Q: Why isn't [intuition]? A: Because [reality] |
| **Formula Intuition** | User learned a formula | Q: What does each part mean? A: [breakdown] |
| **Visual/Spatial** | Concept needs picture | Q: Draw/visualize X. A: [ASCII or description] |
| **Connection** | Links between concepts | Q: How does X relate to Y? A: [relationship] |
| **Principle** | General rule extracted | Q: When [situation], what to do? A: [principle] |

### Card Template

```markdown
### Card N: {Title}
**Type**: {Misconception|Formula|Visual|Connection|Principle}

Q: {Question that tests understanding}

A: {Concise answer}
   {Optional: one concrete example}
```

### Extraction Process

1. **Identify the insight**: What did user actually learn? (not what was taught)
2. **Find the friction**: Where was the confusion/misconception?
3. **Frame as question**: What question would test if they truly understood?
4. **Write minimal answer**: Just enough to trigger recall

### Example: From Session Learning

User's journey:
```
Intuition: "n=m means each slot gets 1 item"
Reality: "Random ≠ uniform, expect clustering"
```

Card generated:
```
Q: In a hash table with n=m, why is E[chain]=2, not 1?

A: Random hashing ≠ perfect distribution.
   Like birthday paradox: clustering happens.
   Some slots get 0, some get 2+.
```

### Data Location

```
{project}/docs/{topic}-quiz.md       # Project-specific
~/.lessons/quizzes/{category}.md     # Personal cross-project
```

### Quiz File Format

```markdown
---
topic: {topic}
created: YYYY-MM-DD
source_session: {session file if applicable}
---

# {Topic} Quiz Cards

## Card 1: {Title}
**Type**: {type}

Q: ...
A: ...

---

## Card 2: {Title}
...
```

### Before Saving, Show Preview

```
📝 Quiz Cards Preview:
─────────────────────
Topic: {topic}
Cards: {N}

1. {card 1 title} ({type})
2. {card 2 title} ({type})
...
─────────────────────
Save to {path}? (y/n)
```

### Integration with Other Modes

| After Mode | Quiz Action |
|------------|-------------|
| Pain → Capture | "Want quiz cards for this lesson?" |
| Analyze | "Generate cards from top patterns?" |
| Session wrap-up | "Extract quiz cards from today's insights?" |

---

## Quick Reference

| User Says | Action |
|-----------|--------|
| "I'm stuck" / frustrated | Pain mode → reflection question |
| "Help me decide" | Decision mode → framework |
| "Lesson learned" | Capture mode → record + extract principle |
| "What did I learn about X" | Recall mode → search lessons |
| "Review my lessons" | Analyze mode → patterns |
| "Quiz me" / "make cards" | Quiz mode → extract cards |
| "Wrap up" | Complete session file + offer quiz cards |

---

## Knowledge Base (Life Principles)

Token-optimized reference from Ray Dalio's Principles Part 2:

| Chapter | File | Key Concepts |
|---------|------|--------------|
| 1. Embrace Reality | `kb/01-embrace-reality.md` | Pain+Reflection=Progress, Evolution, Machine view |
| 2. Five-Step Process | `kb/02-five-step-process.md` | Goals→Problems→Diagnose→Design→Execute |
| 3. Radical Open-Minded | `kb/03-radical-open-minded.md` | Ego barrier, Blind spots, Thoughtful disagreement |
| 4. People Wired Differently | `kb/04-people-wired-differently.md` | Brain types, MBTI, Team roles |
| 5. Effective Decisions | `kb/05-effective-decisions.md` | Learning vs Deciding, Synthesis, Expected value |

When user needs deeper context, read relevant kb file.

---

## Deep Dives

- 5-Step Process details: `kb/02-five-step-process.md`
- Open-mindedness diagnostics: `kb/03-radical-open-minded.md`
- Decision frameworks: `kb/05-effective-decisions.md`
- Root cause analysis: `kb/01-embrace-reality.md` (section 1.10)
