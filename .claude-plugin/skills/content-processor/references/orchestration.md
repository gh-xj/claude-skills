# Sub-Agent Orchestration

## Pre-Execution Config

**Always show and confirm before executing:**

```
📋 Processing Config
────────────────────
Files: {count}
Batch size: {calculated}
Model: sonnet
Language: {detected or ask}
Output: {path}
────────────────────
Proceed? (y/n/adjust)
```

## Batch Size Calculation

```
BATCH_SIZE = min(10, max(5, file_count / 3))
```

| File Count | Default Batch |
| ---------- | ------------- |
| < 15       | 5             |
| 15-30      | 5-10          |
| > 30       | 10            |

User can override: "batch 8", "use haiku", "language zh"

---

## Rules (CRITICAL)

| Rule           | Details                                            |
| -------------- | -------------------------------------------------- |
| Confirm first  | Always show config, wait for approval              |
| Batch size     | Default calculated, user can override              |
| Model          | Default sonnet, user can override                  |
| Agent output   | Write to file, return **only** "Done: \<item\>"    |
| Verification   | Use `ls <dir>/` — **never** TaskOutput for content |
| Progress       | Update PROGRESS.md after each batch                |
| Context safety | Sub-agents must NOT return summary content         |

---

## PROGRESS.md Format

```markdown
# Processing Progress

## Status: 5/13 completed

| Item | Status | Notes       |
| ---- | ------ | ----------- |
| 001  | ✅     | Done        |
| 002  | 🔄     | In progress |
| 003  | ⏳     | Pending     |
```

---

## Prompt Templates

### URL Summary

```
Read <full_path_to_source>.md and generate a summary.

Language: <zh/en>

Write to: <full_path>/summaries/<slug>_summary.md

IMPORTANT: Return ONLY "Done: <item>" - do not return summary content.
```

### Book Chapter

```
Read <full_path_to_source>.md and generate a detailed summary.

Language: <zh/en>

Preserve the chapter's structure and logical flow, not just key points.
Include important quotes and use tables where helpful.

Write to: <full_path>/ai-summary/<chapter>_summary.md

IMPORTANT: Return ONLY "Done: <chapter>" - do not return summary content.
```

---

## Execution Pattern

```
1. Setup: mkdir -p <summary_dir> && touch PROGRESS.md
2. For each batch:
   a. Launch N Task tools in SINGLE message (parallel)
   b. Verify: ls <summary_dir>/
   c. Update PROGRESS.md
   d. Report: "Batch N complete: X/Y total"
3. Never use TaskOutput to retrieve content
```
