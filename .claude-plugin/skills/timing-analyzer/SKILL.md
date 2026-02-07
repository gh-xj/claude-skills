---
name: timing-analyzer
description: Personal time analyzer for Timing.app data. Use when user asks about time tracking, productivity stats, app usage, daily/weekly summaries, rule configuration, unassigned activities, or wants to explore how they spend time.
---

# Timing Analyzer

Query and analyze personal time tracking data from Timing.app.

## Additional Resources

- SQL queries: `references/queries.md`
- Rules configuration: `references/rules.md`
- Cognitive psychology analysis: `references/cognitive.md`

## Database

```
~/Library/Application Support/info.eurocomp.Timing2/SQLite.db
```

## Critical: Deduplication

**TaskActivity (manual) takes priority over AppActivity (auto-tracked).**

When they overlap, only TaskActivity counts. All queries must exclude overlapping AppActivity.

## Timestamp Handling

```sql
datetime(startDate, 'unixepoch', 'localtime')  -- convert to readable
date('now', 'localtime')                        -- current date
```

## Data Sources

| Table | Description |
|-------|-------------|
| AppActivity | Auto-tracked (exclude when overlaps TaskActivity) |
| TaskActivity | Manual entries (takes priority) |
| Application | App metadata |
| Project | Hierarchical categories |
| Title | Window titles |

## Output Format

- Visual bars: `printf('%.*c', hours/scale, '█')`
- Percentages for context
- Time as `Xh Ym` not decimals
- Tables for scanning

## Analysis Workflows

### When User Asks for Daily Report

Run deduplicated daily report query from `references/queries.md`.

### When Analyzing a Category

1. **Overview** - Monthly trend, sub-project breakdown
2. **Drill-Down** - Window titles, URLs, group by patterns
3. **Time Patterns** - By hour, by day, weekly trend
4. **Insights** - Peaks, lows, concerning patterns

### Content Extraction

```sql
CASE 
    WHEN t.stringValue LIKE '%project_a%' THEN 'Project A'
    WHEN t.stringValue LIKE '%project_b%' THEN 'Project B'
    ELSE 'Other'
END as category
```

## Cognitive Psychology Analysis

Use when user asks about:
- Personal growth, productivity patterns
- Deep work vs shallow work
- Attention, focus, distraction analysis
- Chinese: "认知心理学分析", "个人成长"

See `references/cognitive.md` for frameworks and queries.

## Helper Scripts

Located in `scripts/` directory:
- `today.sh` - Today's summary
- `week.sh` - Weekly breakdown
- `search.sh <keyword>` - Search titles/paths

## Rules & Categorization

### When User Asks About Rules

See `references/rules.md` for:
- Reading existing rules via AppleScript
- Decoding protobuf rule data
- Finding unassigned activities (rule gaps)
- Building rules programmatically

### Quick: Find Unassigned Activities

```sql
SELECT a.title as app, ROUND(SUM(aa.endDate - aa.startDate) / 60.0, 0) as mins
FROM AppActivity aa JOIN Application a ON aa.applicationID = a.id
WHERE date(aa.startDate, 'unixepoch', 'localtime') >= date('now', '-7 days')
AND aa.isDeleted = 0 AND aa.projectID IS NULL
GROUP BY app ORDER BY mins DESC LIMIT 20;
```

### Adding Rules (Manual)

Current license doesn't support AppleScript rule updates. Use:
1. Timing.app → Activities view
2. Find unassigned activity
3. **Option + drag** onto target project
