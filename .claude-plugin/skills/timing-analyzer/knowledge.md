# Timing.app Knowledge Base

## Critical: Deduplication Logic

**The app deduplicates overlapping time between AppActivity and TaskActivity.**

When TaskActivity (manual entry) overlaps with AppActivity (auto-tracked), only TaskActivity counts.

### Deduplication CTE Pattern
```sql
TaskRanges AS (
    SELECT startDate as t_start, endDate as t_end
    FROM TaskActivity
    WHERE date(startDate, 'unixepoch', 'localtime') = date('now', 'localtime')
      AND isDeleted = 0
),
NonOverlappingApp AS (
    SELECT aa.projectID, (aa.endDate - aa.startDate) as duration
    FROM AppActivity aa
    WHERE date(aa.startDate, 'unixepoch', 'localtime') = date('now', 'localtime')
      AND aa.isDeleted = 0
      AND NOT EXISTS (
          SELECT 1 FROM TaskRanges tr
          WHERE aa.startDate >= tr.t_start AND aa.startDate < tr.t_end
      )
),
TaskDurations AS (
    SELECT projectID, (endDate - startDate) as duration
    FROM TaskActivity
    WHERE date(startDate, 'unixepoch', 'localtime') = date('now', 'localtime')
      AND isDeleted = 0
),
AllActivity AS (
    SELECT * FROM NonOverlappingApp UNION ALL SELECT * FROM TaskDurations
)
```

## Database Schema

### AppActivity (Auto-tracked)
```sql
CREATE TABLE "AppActivity" (
    "id" INTEGER PRIMARY KEY,
    "localDeviceID" INTEGER REFERENCES "Device",
    "startDate" REAL NOT NULL,      -- Unix timestamp
    "endDate" REAL NOT NULL,        -- Unix timestamp
    "applicationID" INTEGER NOT NULL REFERENCES "Application",
    "titleID" INTEGER REFERENCES "Title",
    "pathID" INTEGER REFERENCES "Path",
    "projectID" INTEGER REFERENCES "Project",
    "isDeleted" BOOLEAN DEFAULT 0
);
```

### TaskActivity (Manual entries - takes priority)
```sql
CREATE TABLE "TaskActivity" (
    "id" INTEGER PRIMARY KEY,
    "startDate" REAL NOT NULL,      -- Unix timestamp
    "endDate" REAL NOT NULL,        -- Unix timestamp
    "projectID" INTEGER REFERENCES "Project",
    "title" TEXT,
    "notes" TEXT,
    "isDeleted" BOOLEAN DEFAULT 0,
    "isRunning" BOOLEAN DEFAULT 0
);
```

### Project (Hierarchical via parentID)
```sql
CREATE TABLE "Project" (
    "id" INTEGER PRIMARY KEY,
    "title" TEXT NOT NULL,
    "parentID" INTEGER REFERENCES "Project",  -- NULL = root
    "color" TEXT NOT NULL,           -- Hex RGBA e.g. #4DC2FFFF
    "productivityScore" REAL DEFAULT 0,  -- -1.0 to +1.0
    "isArchived" BOOLEAN DEFAULT 0
);
```

## Timestamp Handling

Timestamps are **standard Unix epoch** (seconds since 1970-01-01):

```sql
-- Convert to readable datetime
datetime(startDate, 'unixepoch', 'localtime')

-- Current date
date('now', 'localtime')

-- N days ago
strftime('%s', 'now', '-7 days')
```

## Common CTEs

### Recursive Project Hierarchy
```sql
WITH RECURSIVE ProjectHierarchy AS (
    SELECT id, title, parentID, title as root_title FROM Project WHERE parentID IS NULL
    UNION ALL
    SELECT p.id, p.title, p.parentID, ph.root_title FROM Project p
    JOIN ProjectHierarchy ph ON p.parentID = ph.id
)
```

### Time Formatting
```sql
-- Hours and minutes
CASE WHEN seconds >= 3600 THEN (seconds/3600) || 'h' || ((seconds%3600)/60) || 'm'
     ELSE (seconds/60) || 'm' END
```

### Query Project Hierarchy
```sql
-- Top-level projects
SELECT title, productivityScore FROM Project WHERE parentID IS NULL ORDER BY title;

-- Full hierarchy
WITH RECURSIVE ProjectHierarchy AS (
    SELECT id, title, parentID, title as root, 0 as depth FROM Project WHERE parentID IS NULL
    UNION ALL
    SELECT p.id, p.title, p.parentID, ph.root, ph.depth + 1
    FROM Project p JOIN ProjectHierarchy ph ON p.parentID = ph.id
)
SELECT printf('%.*c', depth*2, ' ') || title as project, root FROM ProjectHierarchy ORDER BY root, depth;
```

## Data Quirks

| Issue | Cause | Solution |
|-------|-------|----------|
| Totals exceed app | Missing deduplication | Use `NOT EXISTS` to exclude AppActivity during TaskActivity |
| Screen Time gaps | iOS import limitations | Shows "(No additional information available)" |

---

## Rules & Auto-Categorization

### How Rules Work

Rules are stored in `Project.predicate` column as **binary protobuf format**. Each project can have conditions that auto-assign activities to it.

### Project Schema (Full)
```sql
CREATE TABLE "Project" (
    "id" INTEGER PRIMARY KEY NOT NULL,
    "title" TEXT NOT NULL,
    "parentID" INTEGER REFERENCES "Project",  -- NULL = root category
    "listPosition" INTEGER NOT NULL,
    "isSample" BOOLEAN NOT NULL DEFAULT 0,
    "color" TEXT NOT NULL,                    -- Hex RGBA e.g. #4DC2FFFF
    "productivityScore" REAL NOT NULL DEFAULT 0,  -- -1.0 to +1.0
    "predicate" BLOB,                         -- Rule conditions (protobuf)
    "ruleListPosition" INTEGER NOT NULL,
    "isArchived" BOOLEAN NOT NULL DEFAULT 0,
    "membershipID" INTEGER,
    "property_bag" TEXT
);
```

### Decoding Rules (Read-Only)

To inspect rule content, decode the protobuf as strings:
```bash
sqlite3 ~/Library/Application\ Support/info.eurocomp.Timing2/SQLite.db "
SELECT hex(predicate) FROM Project WHERE title = 'YourProject';" | xxd -r -p | strings
```

### Rule Condition Types

| Condition | Format in Protobuf | Example |
|-----------|-------------------|---------|
| Application list | `applicationID` + `ApplicationList.<name>` | `ApplicationList.media` |
| Domain list | `webDomain` + `DomainList.<name>` | `DomainList.media` |
| Custom domain | `webDomain"<domain>*` | `webDomain"ocw.mit.edu*` |
| File path | `filePath"<path>*` | `filePath"/Users/x/codebase*` |
| Window title | `title"<pattern>*` | `title"mit-*` |
| Title or path | `titleOrPath"<pattern>*` | `titleOrPath".go*` |

### Built-in Lists

Timing.app includes predefined lists referenced in rules:
- `ApplicationList.media` - YouTube, Netflix, Spotify, etc.
- `ApplicationList.development` - IDEs, terminals, etc.
- `ApplicationList.communication` - Slack, Mail, Messages, etc.
- `DomainList.media` - youtube.com, netflix.com, etc.

### Rule Priority

When multiple rules match, **the most specific rule wins**. Child projects override parent projects.

Example: If `Media` (under Leisure) matches `youtube.com`, but `Learn` (under Growth) has a rule for `title contains "MIT"`, the activity goes to Learn if both conditions are met.

### Query Projects with Rules
```sql
SELECT id, title, parentID, 
    CASE WHEN predicate IS NOT NULL AND length(predicate) > 0 
         THEN 'YES (' || length(predicate) || ' bytes)' 
         ELSE 'NO' END as has_rule
FROM Project 
WHERE predicate IS NOT NULL AND length(predicate) > 0
ORDER BY title;
```

### Limitations

- **Read-only**: Rules should only be modified through Timing.app UI
- Writing to `predicate` column risks data corruption (protobuf format)
- Cannot create new rules programmatically

### Common Categorization Issues

| Symptom | Cause | Fix (in Timing UI) |
|---------|-------|-------------------|
| YouTube lectures → Leisure | `ApplicationList.media` matches first | Add rule to Learn: `domain=youtube.com AND title contains "MIT"` |
| Code files → wrong project | Missing file path rule | Add rule: `filePath contains "/project-name/"` |
| Browser tabs miscategorized | Domain not in custom rules | Add specific `webDomain` rule |

---

## Cognitive Psychology Frameworks

### Deep Work Framework (Cal Newport)

| Category | Duration | Cognitive Value | Examples |
|----------|----------|-----------------|----------|
| Deep Work | ≥60 min | High output, skill growth | Coding, writing, design |
| Focused | 25-60 min | Moderate output | Research, review |
| Operational | 10-25 min | Necessary but low-growth | Email, meetings prep |
| Shallow/Micro | <10 min | Low value, interruptible | Notifications, quick checks |

**Benchmark**: 4h/day deep work = 50% of 8h workday = excellent

### Temporal Discounting

The tendency to prefer immediate rewards over delayed benefits.

| Indicator | Healthy | Warning | Critical |
|-----------|---------|---------|----------|
| Growth:Leisure Ratio | ≥1:2 | 1:3 - 1:4 | <1:4 |

**Interpretation**: High leisure relative to growth suggests instant gratification bias.

### Attention Residue (Leroy, 2009)

Each task switch incurs ~23 minutes cognitive recovery cost.

| Daily Switches | Assessment | Impact |
|----------------|------------|--------|
| <500 | Focused | Minimal residue |
| 500-1500 | Normal | Moderate impact |
| 1500-3000 | Fragmented | Significant efficiency loss |
| >3000 | Severe | 40-60% productivity reduction |

**Switch Variance**: MAX/MIN ratio >3x indicates unstable attention patterns.

### Embodied Cognition

Physical activity directly affects cognitive function via BDNF (brain-derived neurotrophic factor).

| Health Investment | Risk Level | Cognitive Impact |
|-------------------|------------|------------------|
| ≥3.5h/week (30min/day) | Optimal | Enhanced memory, focus |
| 1-3h/week | Adequate | Baseline maintenance |
| <1h/week | Warning | Gradual decline |
| Near zero | Critical | Accelerated cognitive aging |

### Cognitive Profile Scoring

| Dimension | Metric | Score Calculation |
|-----------|--------|-------------------|
| 专注力 (Focus) | Deep Work % | 51%+ = 8+, 40-50% = 6-7, <40% = 1-5 |
| 自控力 (Self-control) | G:L Ratio | 1:2 = 8+, 1:3 = 6, 1:4+ = 1-4 |
| 稳定性 (Stability) | Switch Variance | <2x = 8+, 2-4x = 5-7, >4x = 1-4 |
| 身体投资 (Body) | Health hours/week | 3.5h+ = 8+, 1-3h = 5-7, <1h = 1-4 |

### Behavioral Design Principles

When generating recommendations:

| Principle | Application |
|-----------|-------------|
| Environment Design | Modify context to reduce friction (notifications off, focus mode) |
| Habit Stacking | Attach new behavior to existing routine |
| Micro-habits | Start with tiny commitment (5 min/day) then expand |
| Trigger Binding | Link desired behavior to consistent cue |
| Substitution | Replace unwanted behavior with better alternative |

### Report Language

Support bilingual output (Chinese headers for Chinese requests):
- 时间分布总览 = Time Distribution Overview
- 深度工作分析 = Deep Work Analysis  
- 注意力碎片化 = Attention Fragmentation
- 认知模式诊断 = Cognitive Profile Diagnosis
- 行为改变建议 = Behavioral Change Recommendations
