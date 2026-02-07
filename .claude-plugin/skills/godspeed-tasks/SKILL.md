---
name: godspeed-tasks
description: Manage Godspeed tasks. Use when user asks to view, query, create, or manage tasks in Godspeed. Supports read via CLI, write via HTTP API or SQLite direct.
---

# Godspeed Task Manager

## Quick Reference

| Operation | Method | Notes |
|-----------|--------|-------|
| Read | `go run scripts/.` | Built-in CLI |
| Write (with token) | HTTP API | Preferred |
| Write (no token) | SQLite direct | Restart app after |

### Database Path

```
~/Library/Application Support/Godspeed/godspeed-db-setapp.sqlite  # Setapp
~/Library/Application Support/Godspeed/godspeed-db.sqlite         # Standard
```

## Additional Resources

For full API details and SQLite schema, see `knowledge.md`.

## CLI Usage

Run from scripts directory:

```bash
cd ~/.claude/skills/godspeed-tasks/scripts

# Statistics
go run . stats

# List tasks
go run . list --status incomplete --limit 20
go run . list --list "home" --status incomplete
go run . list --due
go run . list --recent
go run . list --completed

# Search
go run . search "keyword"

# Get task details
go run . get <task-id>

# Show all lists
go run . lists

# JSON output
go run . list --format json
```

### CLI Flags

| Flag | Description |
|------|-------------|
| `--status <incomplete\|complete>` | Filter by status |
| `--list <name-or-id>` | Filter by list |
| `--limit <n>` | Limit results |
| `--format <table\|json>` | Output format |
| `--due` | Tasks with due dates |
| `--recent` | Recently updated |
| `--completed` | Recently completed |

## Write Operations

### Create Sub-Tasks with Notes (SQLite)

**Step 1: Get parent task info**

```bash
DB=~/Library/Application\ Support/Godspeed/godspeed-db-setapp.sqlite
sqlite3 "$DB" "SELECT id, list_id, order_index, indent_level, user_id FROM todo_items WHERE id='<parent-id>';"
```

**Step 2: Insert sub-tasks with notes**

```bash
sqlite3 "$DB" "
INSERT INTO todo_items (id, title, notes, list_id, order_index, indent_level, created_at, updated_at, user_id, timeless_due_at) VALUES
(lower(hex(randomblob(4)) || '-' || hex(randomblob(2)) || '-4' || substr(hex(randomblob(2)),2) || '-' || substr('89ab',abs(random()) % 4 + 1,1) || substr(hex(randomblob(2)),2) || '-' || hex(randomblob(6))),
'Task title',
'Task notes here
Can be multiline
Include links: https://example.com',
'<list_id>', <order> + 0.0001, <indent> + 1, datetime('now'), datetime('now'), <user_id>, '2026-02-03');
"
```

**Step 3: Restart Godspeed app**

### Update Task (SQLite)

```bash
DB=~/Library/Application\ Support/Godspeed/godspeed-db-setapp.sqlite

# Update notes
sqlite3 "$DB" "UPDATE todo_items SET notes = 'New notes' WHERE id = '<task-id>';"

# Mark complete
sqlite3 "$DB" "UPDATE todo_items SET completed_at = datetime('now') WHERE id = '<task-id>';"

# Mark incomplete
sqlite3 "$DB" "UPDATE todo_items SET completed_at = NULL WHERE id = '<task-id>';"

# Update due date
sqlite3 "$DB" "UPDATE todo_items SET timeless_due_at = '2026-02-03' WHERE id = '<task-id>';"
```

### Key SQLite Fields

| Field | Value | Effect |
|-------|-------|--------|
| `order_index` | parent + 0.0001 | Position after parent |
| `indent_level` | parent + 1 | Makes it a child |
| `indent_level` | parent + 2 | Makes it a grandchild |
| `timeless_due_at` | 'YYYY-MM-DD' | Due date (optional) |
| `notes` | 'text' | Task notes/description |
| `completed_at` | datetime('now') | Mark complete |
| `completed_at` | NULL | Mark incomplete |

### Create Task (HTTP API)

```bash
export GODSPEED_TOKEN="your-token"  # From Godspeed: Cmd+K -> "Copy API access token"

curl -X POST https://api.godspeedapp.com/tasks \
  -H "Authorization: Bearer $GODSPEED_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title": "New task", "list_id": "<list-id>"}'
```

## Notes Field (Markdown)

The `notes` field supports markdown rendering in Godspeed:

| Syntax | Example |
|--------|---------|
| Headers | `### Section` |
| Lists | `- item` or `1. item` |
| Checkboxes | `- [ ] todo` / `- [x] done` |
| Links | `[text](url)` or raw URLs |

**Example update:**
```sql
UPDATE todo_items SET notes = '### Goal
Find CPA for taxes

### Steps
- [ ] Search online
- [ ] Call for quotes
- [x] Gather documents'
WHERE id = '<task-id>';
```

## Due Date Strategy

Set task due dates **before** actual deadlines to build buffer:

| Actual Deadline | Task Due Date | Buffer |
|-----------------|---------------|--------|
| Available date (e.g., 1099s Feb 15) | +2 days after | Gives time to arrive |
| Hard deadline (e.g., Apr 15 taxes) | -2 weeks | Safe margin |
| Soft deadline | -1 week | Reasonable buffer |

**Principle:** Task due date = when YOU should complete it, not the external deadline.

## Safety Rules

1. **Confirm before delete** - Irreversible
2. **After SQLite write** - Restart Godspeed app to see changes
3. **Rate limits** - 60 writes/min, 10 reads/min
