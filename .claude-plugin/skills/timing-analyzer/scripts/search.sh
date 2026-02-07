#!/bin/bash
# Search window titles or paths
# Usage: ./search.sh <keyword>

DB=~/Library/Application\ Support/info.eurocomp.Timing2/SQLite.db
KEYWORD="${1:-}"

if [ -z "$KEYWORD" ]; then
    echo "Usage: $0 <keyword>"
    echo "Search window titles and file paths for a keyword"
    exit 1
fi

echo "=== Window Titles matching '$KEYWORD' ==="
sqlite3 -header -column "$DB" "
SELECT substr(t.stringValue, 1, 60) as Title,
    count(*) as Count,
    cast(sum((aa.endDate - aa.startDate))/3600 as int) || 'h' || 
    cast((sum((aa.endDate - aa.startDate))%3600)/60 as int) || 'm' as Time
FROM AppActivity aa
JOIN Title t ON aa.titleID = t.id
WHERE t.stringValue LIKE '%$KEYWORD%' AND aa.isDeleted = 0
GROUP BY t.stringValue ORDER BY sum(aa.endDate - aa.startDate) DESC LIMIT 20;
"

echo ""
echo "=== Paths/URLs matching '$KEYWORD' ==="
sqlite3 -header -column "$DB" "
SELECT substr(p.stringValue, 1, 70) as Path,
    cast(sum((aa.endDate - aa.startDate))/3600 as int) || 'h' || 
    cast((sum((aa.endDate - aa.startDate))%3600)/60 as int) || 'm' as Time
FROM AppActivity aa
JOIN Path p ON aa.pathID = p.id
WHERE p.stringValue LIKE '%$KEYWORD%' AND aa.isDeleted = 0
GROUP BY p.stringValue ORDER BY sum(aa.endDate - aa.startDate) DESC LIMIT 20;
"
