# Timing.app SQL Queries

Database: `~/Library/Application Support/info.eurocomp.Timing2/SQLite.db`

## Daily Report (Deduplicated)

```sql
sqlite3 -header -column ~/Library/Application\ Support/info.eurocomp.Timing2/SQLite.db "
WITH RECURSIVE ProjectHierarchy AS (
    SELECT id, title, parentID, title as root_title FROM Project WHERE parentID IS NULL
    UNION ALL
    SELECT p.id, p.title, p.parentID, ph.root_title FROM Project p
    JOIN ProjectHierarchy ph ON p.parentID = ph.id
),
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
),
CategoryTotals AS (
    SELECT COALESCE(ph.root_title, 'Unassigned') as category, cast(sum(a.duration) as int) as seconds
    FROM AllActivity a LEFT JOIN ProjectHierarchy ph ON a.projectID = ph.id
    GROUP BY ph.root_title
),
GrandTotal AS (SELECT sum(seconds) as total FROM CategoryTotals)
SELECT ct.category,
    CASE WHEN ct.seconds >= 3600 THEN (ct.seconds/3600) || 'h' || ((ct.seconds%3600)/60) || 'm'
         ELSE (ct.seconds/60) || 'm' END as time,
    printf('%.1f%%', 100.0 * ct.seconds / gt.total) as pct,
    printf('%.*c', cast(30.0 * ct.seconds / gt.total as int), '█') as bar
FROM CategoryTotals ct, GrandTotal gt
ORDER BY ct.seconds DESC;
"
```

## Hourly Breakdown

```sql
sqlite3 -header -column ~/Library/Application\ Support/info.eurocomp.Timing2/SQLite.db "
WITH RECURSIVE ProjectHierarchy AS (
    SELECT id, title, parentID, title as root_title FROM Project WHERE parentID IS NULL
    UNION ALL
    SELECT p.id, p.title, p.parentID, ph.root_title FROM Project p
    JOIN ProjectHierarchy ph ON p.parentID = ph.id
),
TaskRanges AS (
    SELECT startDate as t_start, endDate as t_end
    FROM TaskActivity
    WHERE date(startDate, 'unixepoch', 'localtime') = date('now', 'localtime')
      AND isDeleted = 0
),
AppHourly AS (
    SELECT 
        strftime('%H', aa.startDate, 'unixepoch', 'localtime') as hour,
        COALESCE(ph.root_title, 'Unassigned') as category,
        COALESCE(t.stringValue, '') as detail,
        round((aa.endDate - aa.startDate) / 60.0, 0) as minutes
    FROM AppActivity aa
    LEFT JOIN ProjectHierarchy ph ON aa.projectID = ph.id
    LEFT JOIN Title t ON aa.titleID = t.id
    WHERE date(aa.startDate, 'unixepoch', 'localtime') = date('now', 'localtime')
      AND aa.isDeleted = 0
      AND NOT EXISTS (
          SELECT 1 FROM TaskRanges tr
          WHERE aa.startDate >= tr.t_start AND aa.startDate < tr.t_end
      )
),
TaskHourly AS (
    SELECT 
        strftime('%H', ta.startDate, 'unixepoch', 'localtime') as hour,
        COALESCE(ph.root_title, 'Unassigned') as category,
        COALESCE(ta.title, ph.root_title, '') as detail,
        round((ta.endDate - ta.startDate) / 60.0, 0) as minutes
    FROM TaskActivity ta
    LEFT JOIN ProjectHierarchy ph ON ta.projectID = ph.id
    WHERE date(ta.startDate, 'unixepoch', 'localtime') = date('now', 'localtime')
      AND ta.isDeleted = 0
),
AllHourly AS (
    SELECT * FROM AppHourly UNION ALL SELECT * FROM TaskHourly
)
SELECT hour, category, 
    GROUP_CONCAT(DISTINCT CASE WHEN detail != '' THEN substr(detail, 1, 30) END) as details,
    cast(sum(minutes) as int) as mins
FROM AllHourly
GROUP BY hour, category
ORDER BY hour, mins DESC;
"
```

## Top Apps This Week

```sql
sqlite3 -header -column ~/Library/Application\ Support/info.eurocomp.Timing2/SQLite.db "
SELECT a.title as app, round(sum((aa.endDate - aa.startDate) / 3600.0), 1) as hours
FROM AppActivity aa
JOIN Application a ON aa.applicationID = a.id
WHERE aa.startDate > strftime('%s', 'now', '-7 days') AND aa.isDeleted = 0
GROUP BY a.title ORDER BY hours DESC LIMIT 15;
"
```

## Manual Time Entries Today

```sql
sqlite3 -header -column ~/Library/Application\ Support/info.eurocomp.Timing2/SQLite.db "
SELECT title, 
    cast((endDate-startDate)/3600 as int) || 'h' || cast(((endDate-startDate)%3600)/60 as int) || 'm' as duration,
    strftime('%H:%M', startDate, 'unixepoch', 'localtime') as start_time
FROM TaskActivity
WHERE date(startDate, 'unixepoch', 'localtime') = date('now', 'localtime') AND isDeleted = 0
ORDER BY startDate;
"
```

## Project Hierarchy

```sql
SELECT id, title, parentID, productivityScore FROM Project WHERE parentID IS NULL ORDER BY title;
```

## Category Totals (for Daily Summary)

```sql
sqlite3 -separator '|' ~/Library/Application\ Support/info.eurocomp.Timing2/SQLite.db "
WITH RECURSIVE ProjectHierarchy AS (
    SELECT id, title, parentID, title as root_title FROM Project WHERE parentID IS NULL
    UNION ALL
    SELECT p.id, p.title, p.parentID, ph.root_title FROM Project p
    JOIN ProjectHierarchy ph ON p.parentID = ph.id
),
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
SELECT COALESCE(ph.root_title, 'Unassigned') as category,
    CASE WHEN cast(sum(a.duration) as int) >= 3600 
         THEN (cast(sum(a.duration) as int)/3600) || 'h' || ((cast(sum(a.duration) as int)%3600)/60) || 'm'
         ELSE (cast(sum(a.duration) as int)/60) || 'm' END as time
FROM AllActivity a 
LEFT JOIN ProjectHierarchy ph ON a.projectID = ph.id
GROUP BY ph.root_title
ORDER BY sum(a.duration) DESC;
"
```

## TaskActivity Notes

```sql
sqlite3 -header -column ~/Library/Application\ Support/info.eurocomp.Timing2/SQLite.db "
SELECT 
    cast(strftime('%H', startDate, 'unixepoch', 'localtime') as int) as hour,
    title,
    notes
FROM TaskActivity
WHERE date(startDate, 'unixepoch', 'localtime') = date('now', 'localtime')
  AND isDeleted = 0
  AND notes IS NOT NULL AND notes != ''
ORDER BY startDate;
"
```

## Work Tier-2 Breakdown

```sql
sqlite3 -separator '|' ~/Library/Application\ Support/info.eurocomp.Timing2/SQLite.db "
WITH RECURSIVE ProjectHierarchy AS (
    SELECT id, title, parentID, title as root_title, title as leaf_title 
    FROM Project WHERE parentID IS NULL
    UNION ALL
    SELECT p.id, p.title, p.parentID, ph.root_title, p.title
    FROM Project p JOIN ProjectHierarchy ph ON p.parentID = ph.id
),
TaskRanges AS (
    SELECT startDate as t_start, endDate as t_end FROM TaskActivity
    WHERE date(startDate, 'unixepoch', 'localtime') = date('now', 'localtime') AND isDeleted = 0
),
WorkActivity AS (
    SELECT ph.leaf_title as project, (aa.endDate - aa.startDate)/60 as mins
    FROM AppActivity aa
    JOIN ProjectHierarchy ph ON aa.projectID = ph.id
    WHERE date(aa.startDate, 'unixepoch', 'localtime') = date('now', 'localtime')
      AND aa.isDeleted = 0 AND ph.root_title = 'Work'
      AND NOT EXISTS (SELECT 1 FROM TaskRanges tr WHERE aa.startDate >= tr.t_start AND aa.startDate < tr.t_end)
    UNION ALL
    SELECT COALESCE(NULLIF(ta.title,''), ph.leaf_title) as project, (ta.endDate - ta.startDate)/60 as mins
    FROM TaskActivity ta
    JOIN ProjectHierarchy ph ON ta.projectID = ph.id
    WHERE date(ta.startDate, 'unixepoch', 'localtime') = date('now', 'localtime')
      AND ta.isDeleted = 0 AND ph.root_title = 'Work'
)
SELECT project, sum(mins) as total FROM WorkActivity GROUP BY project ORDER BY total DESC LIMIT 8;
"
```

## Hourly Summary with Project Details

Used by hourly-log-filler for granular hour-by-hour data.

```sql
sqlite3 -separator '|' ~/Library/Application\ Support/info.eurocomp.Timing2/SQLite.db "
WITH RECURSIVE ProjectHierarchy AS (
    SELECT id, title, parentID, title as root_title FROM Project WHERE parentID IS NULL
    UNION ALL
    SELECT p.id, p.title, p.parentID, ph.root_title FROM Project p
    JOIN ProjectHierarchy ph ON p.parentID = ph.id
),
TaskRanges AS (
    SELECT startDate as t_start, endDate as t_end
    FROM TaskActivity
    WHERE date(startDate, 'unixepoch', 'localtime') = date('now', 'localtime')
      AND isDeleted = 0
),
AppHourly AS (
    SELECT 
        cast(strftime('%H', aa.startDate, 'unixepoch', 'localtime') as int) as hour,
        COALESCE(ph.root_title, 'Unassigned') as category,
        COALESCE(t.stringValue, '') as detail,
        (aa.endDate - aa.startDate) / 60.0 as minutes
    FROM AppActivity aa
    LEFT JOIN ProjectHierarchy ph ON aa.projectID = ph.id
    LEFT JOIN Title t ON aa.titleID = t.id
    WHERE date(aa.startDate, 'unixepoch', 'localtime') = date('now', 'localtime')
      AND aa.isDeleted = 0
      AND NOT EXISTS (
          SELECT 1 FROM TaskRanges tr
          WHERE aa.startDate >= tr.t_start AND aa.startDate < tr.t_end
      )
),
TaskHourly AS (
    SELECT 
        cast(strftime('%H', ta.startDate, 'unixepoch', 'localtime') as int) as hour,
        COALESCE(ph.root_title, 'Maintain') as category,
        COALESCE(NULLIF(ta.title, ''), p.title, '') as detail,
        (ta.endDate - ta.startDate) / 60.0 as minutes
    FROM TaskActivity ta
    LEFT JOIN ProjectHierarchy ph ON ta.projectID = ph.id
    LEFT JOIN Project p ON ta.projectID = p.id
    WHERE date(ta.startDate, 'unixepoch', 'localtime') = date('now', 'localtime')
      AND ta.isDeleted = 0
),
AllHourly AS (
    SELECT * FROM AppHourly UNION ALL SELECT * FROM TaskHourly
),
HourSummary AS (
    SELECT hour, category, sum(minutes) as total_mins,
        GROUP_CONCAT(DISTINCT CASE WHEN detail != '' THEN substr(detail, 1, 40) END) as details
    FROM AllHourly
    GROUP BY hour, category
)
SELECT hour, category, cast(total_mins as int) as mins, details
FROM HourSummary
ORDER BY hour, total_mins DESC;
"
```

## Window Title Keywords

Common patterns for summarizing activities:

| Pattern | Summary |
|---------|---------|
| `YouTube – MIT` | MIT lecture |
| `leetcode` | leetcode |
| `project-*` | project name |
| `<app_name>` | app name |
