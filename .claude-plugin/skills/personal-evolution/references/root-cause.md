# Root Cause Analysis - Deep Dive

## The Key Distinction

| Type | Description | Example |
|------|-------------|---------|
| **Proximate cause** | The action (verb) | "I didn't check the schedule" |
| **Root cause** | The trait/system (adjective/noun) | "I'm forgetful" or "No reminder system" |

**Rule**: Keep asking "why?" until you hit a character trait or system design.

---

## Three Root Cause Types

### 1. Character (You)

The root cause is a personal trait or tendency.

**Examples:**
- Forgetfulness
- Impatience
- Conflict avoidance
- Overconfidence
- Perfectionism

**Fix options:**
- Awareness + conscious effort (hard, slow)
- Guardrails (systems that compensate)
- Delegate to someone without this trait

**Dalio's view**: "You can't change your wiring much. Create guardrails instead."

---

### 2. Process (How)

The root cause is a missing or broken process.

**Examples:**
- No checklist for recurring task
- No verification step before sending
- No defined handoff between people
- No documentation for decisions

**Fix options:**
- Create the missing process
- Automate the check
- Add accountability step

**Signs it's a process problem:**
- "I forgot to..." (no reminder system)
- "I assumed..." (no verification step)
- "Nobody told me..." (no communication process)

---

### 3. Design (What)

The root cause is structural - the system/org/tool is wrong.

**Examples:**
- Wrong person in role
- Tool doesn't fit the job
- Incentives misaligned
- Org structure creates conflict

**Fix options:**
- Redesign the system
- Change the tool
- Move people to right roles
- Realign incentives

**Signs it's a design problem:**
- Multiple people fail at same thing
- Process exists but doesn't work
- Success requires heroic effort

---

## The "5 Whys" Method

Start with the problem, ask "why" until you hit root cause.

**Example:**
1. **Problem**: Project delivered late
2. **Why?** Testing took longer than expected
3. **Why?** Found bugs late in the process
4. **Why?** No code review before testing
5. **Why?** Team doesn't have review process
6. **Root cause**: Process gap (no code review step)

**Example 2:**
1. **Problem**: Sent wrong file to client
2. **Why?** Grabbed old version
3. **Why?** Didn't check version number
4. **Why?** Was rushing
5. **Why?** Always leave things to last minute
6. **Root cause**: Character trait (procrastination)

---

## Diagnostic Questions

| Question | Reveals |
|----------|---------|
| "Has this happened before?" | Pattern vs one-off |
| "Would this happen to anyone in this role?" | Design vs character |
| "Is there a process that should have caught this?" | Process gap |
| "Did I know better but do it anyway?" | Character trait |
| "What system would prevent this?" | Process/design fix |

---

## From Root Cause to Principle

Once you identify root cause, generalize:

| Root Cause | Principle |
|------------|-----------|
| "I'm forgetful" | "When task is important, set 2 reminders" |
| "No review process" | "Nothing ships without second pair of eyes" |
| "Wrong tool" | "Test tools on small task before committing" |
| "Overconfidence" | "When confident, seek disagreement harder" |

**Template**: "When [situation], [action], because [root cause tends to cause problem]"

---

## The Machine View

Dalio's perspective: You are both the **designer** and **worker** in your life machine.

```
Designer-You: Diagnoses problems, designs fixes
     ↓
Machine (processes, tools, habits)
     ↓
Worker-You: Executes within the machine
```

**When something fails:**
1. Don't blame worker-you
2. Ask: Is the machine designed correctly?
3. If worker-you can't do the job, find someone who can

> "The biggest mistake is not objectively seeing yourself and others as you/they really are."

---

## Red Flags: Surface-Level "Fixes"

| Surface Fix | Why It Fails | Real Fix |
|-------------|--------------|----------|
| "I'll try harder" | Doesn't address root cause | Identify what makes it hard |
| "I'll remember next time" | Memory is unreliable | Create reminder system |
| "I just need more time" | Time wasn't the issue | Identify what actually blocked |
| "Won't happen again" | No mechanism to prevent | Design prevention |

---

## Quick Reference

When capturing a lesson, always ask:

1. **What's the proximate cause?** (the action/inaction)
2. **What's the root cause?** (keep asking "why?")
3. **Is it Character, Process, or Design?**
4. **What principle prevents this next time?**
