#!/bin/bash
# Today's time summary with deduplication (matches app totals)

DB=~/Library/Application\ Support/info.eurocomp.Timing2/SQLite.db

echo "=== Today's Time Summary ==="
echo ""

sqlite3 -header -column "$DB" "
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
SELECT ct.category as Category,
    CASE WHEN ct.seconds >= 3600 THEN (ct.seconds/3600) || 'h' || ((ct.seconds%3600)/60) || 'm'
         ELSE (ct.seconds/60) || 'm' END as Time,
    printf('%.1f%%', 100.0 * ct.seconds / gt.total) as Pct,
    printf('%.*c', cast(25.0 * ct.seconds / gt.total as int), '█') as Bar
FROM CategoryTotals ct, GrandTotal gt
ORDER BY ct.seconds DESC;
"

echo ""
echo "=== Manual Entries ==="
sqlite3 -header -column "$DB" "
SELECT COALESCE(title, '(untitled)') as Task,
    cast((endDate-startDate)/3600 as int) || 'h' || cast(((endDate-startDate)%3600)/60 as int) || 'm' as Duration,
    strftime('%H:%M', startDate, 'unixepoch', 'localtime') as Start
FROM TaskActivity
WHERE date(startDate, 'unixepoch', 'localtime') = date('now', 'localtime')
  AND isDeleted = 0
ORDER BY startDate;
"
