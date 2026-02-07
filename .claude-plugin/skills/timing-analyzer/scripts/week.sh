#!/bin/bash
# Weekly summary - apps and projects

DB=~/Library/Application\ Support/info.eurocomp.Timing2/SQLite.db

echo "=== This Week's Top Apps ==="
echo ""

sqlite3 -header -column "$DB" "
SELECT a.title as App,
    cast(sum((aa.endDate - aa.startDate))/3600 as int) || 'h' || 
    cast((sum((aa.endDate - aa.startDate))%3600)/60 as int) || 'm' as Time,
    round(sum((aa.endDate - aa.startDate) / 3600.0), 1) as Hours
FROM AppActivity aa
JOIN Application a ON aa.applicationID = a.id
WHERE aa.startDate > strftime('%s', 'now', '-7 days') AND aa.isDeleted = 0
GROUP BY a.title ORDER BY Hours DESC LIMIT 15;
"

echo ""
echo "=== Daily Breakdown ==="
sqlite3 -header -column "$DB" "
SELECT date(startDate, 'unixepoch', 'localtime') as Date,
    CASE cast(strftime('%w', startDate, 'unixepoch', 'localtime') as int)
        WHEN 0 THEN 'Sun' WHEN 1 THEN 'Mon' WHEN 2 THEN 'Tue'
        WHEN 3 THEN 'Wed' WHEN 4 THEN 'Thu' WHEN 5 THEN 'Fri' WHEN 6 THEN 'Sat'
    END as Day,
    cast(sum((endDate - startDate))/3600 as int) || 'h' || 
    cast((sum((endDate - startDate))%3600)/60 as int) || 'm' as Total
FROM AppActivity
WHERE startDate > strftime('%s', 'now', '-7 days') AND isDeleted = 0
GROUP BY Date ORDER BY Date DESC;
"
