# Godspeed Knowledge Base

## CLI Reference

Self-contained Go program in `scripts/` directory.

### Commands

```bash
# From skill directory
cd ~/.claude/skills/godspeed-tasks

# Statistics
go run scripts/. stats
go run scripts/. stats --format json

# List tasks
go run scripts/. list
go run scripts/. list --status incomplete
go run scripts/. list --status complete
go run scripts/. list --list "home"
go run scripts/. list --limit 10
go run scripts/. list --due          # Tasks with due dates
go run scripts/. list --recent       # Recently updated incomplete
go run scripts/. list --completed    # Recently completed

# Search
go run scripts/. search "keyword"
go run scripts/. search "visa" --limit 5

# Get task details
go run scripts/. get <task-id>

# List all lists with task counts
go run scripts/. lists
```

### Output Formats

- `--format table` (default) - Human readable
- `--format json` - JSON for scripting

## API Reference

**Base URL:** `https://api.godspeedapp.com`

### Rate Limits

| Operation | Per Minute | Per Hour |
|-----------|------------|----------|
| Read (GET) | 10 | 200 |
| Write (POST/PATCH/DELETE) | 60 | 1,000 |

### Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/sessions/sign_in` | Authenticate, get token |
| GET | `/tasks` | List tasks (max 250) |
| GET | `/tasks/:id` | Get single task |
| POST | `/tasks` | Create task |
| PATCH | `/tasks/:id` | Update task |
| DELETE | `/tasks/:id` | Delete task |
| GET | `/lists` | Get all lists |

### Task Fields

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `title` | string | Yes | Task title |
| `list_id` | string | No | Defaults to Inbox |
| `location` | string | No | "start" or "end" (default) |
| `notes` | string | No | Task description (supports markdown: headers, lists, checkboxes, links) |
| `due_at` | ISO8601 | No | Due date+time |
| `timeless_due_at` | YYYY-MM-DD | No | Due date only |
| `starts_at` | ISO8601 | No | Start date+time |
| `timeless_starts_at` | YYYY-MM-DD | No | Start date only |
| `duration_minutes` | int | No | Task duration |
| `label_names` | []string | No | Label names (exact match) |
| `label_ids` | []string | No | Label IDs |
| `metadata` | object | No | Key-value pairs, max 1024 chars JSON |

### Update-Only Fields

| Field | Type | Notes |
|-------|------|-------|
| `is_complete` | bool | Mark complete/incomplete |
| `is_cleared` | bool | Clear task (must be complete first) |
| `snoozed_until` | ISO8601 | Snooze date+time |
| `timeless_snoozed_until` | YYYY-MM-DD | Snooze date only |
| `add_label_names` | []string | Labels to add |
| `remove_label_names` | []string | Labels to remove |

### Query Parameters (GET /tasks)

| Param | Values | Notes |
|-------|--------|-------|
| `status` | `incomplete`, `complete` | Filter by completion |
| `list_id` | string | Filter by list |
| `updated_before` | ISO8601 | For pagination |
| `updated_after` | ISO8601 | Filter by update time |

## SQLite Database

### Location (macOS)

```
~/Library/Application Support/Godspeed/godspeed-db-setapp.sqlite  # Setapp
~/Library/Application Support/Godspeed/godspeed-db.sqlite         # Standard
```

### Key Tables

| Table | Purpose |
|-------|---------|
| `todo_items` | Tasks |
| `lists` | Task lists |
| `labels` | Labels |

### Direct Write (when no API token)

```bash
sqlite3 ~/Library/Application\ Support/Godspeed/godspeed-db-setapp.sqlite "
INSERT INTO todo_items (id, title, list_id, order_index, indent_level, created_at, updated_at, user_id, timeless_due_at)
VALUES (
  '$(uuidgen)',
  'Sub-task title',
  '<PARENT_LIST_ID>',
  <PARENT_ORDER_INDEX + 0.0001>,
  1,
  strftime('%Y-%m-%dT%H:%M:%SZ', 'now'),
  strftime('%Y-%m-%dT%H:%M:%SZ', 'now'),
  <USER_ID>,
  'YYYY-MM-DD'
);"
```

**Note:** Direct SQLite writes bypass Godspeed sync. Restart app to see changes.

### Quick SQLite Queries (Read-Only)

```bash
DB=~/Library/Application\ Support/Godspeed/godspeed-db-setapp.sqlite

# Count tasks by status
sqlite3 -readonly "$DB" "SELECT COUNT(*), SUM(CASE WHEN completed_at IS NULL THEN 1 ELSE 0 END) as incomplete FROM todo_items;"

# Incomplete tasks by list
sqlite3 -readonly "$DB" "SELECT l.name, COUNT(*) FROM todo_items t JOIN lists l ON t.list_id = l.id WHERE t.completed_at IS NULL GROUP BY l.name ORDER BY COUNT(*) DESC;"

# Tasks due soon
sqlite3 -readonly "$DB" "SELECT title, date(due_at) FROM todo_items WHERE completed_at IS NULL AND due_at IS NOT NULL ORDER BY due_at LIMIT 10;"
```

## Getting Your Lists and Labels

Run the CLI to discover your list and label IDs:

```bash
go run scripts/. lists    # Shows all lists with IDs
```

Or query SQLite directly:

```bash
sqlite3 -readonly "$DB" "SELECT id, name FROM lists ORDER BY name;"
sqlite3 -readonly "$DB" "SELECT id, name, color FROM labels ORDER BY name;"
```

## Common Pitfalls

| Pitfall | Symptom | Fix |
|---------|---------|-----|
| No token | "authentication required" | Set `GODSPEED_TOKEN` env var |
| Rate limited | 429 response | Wait 1 min for reads |
| Label mismatch | Labels not applied | Use exact label names (case-sensitive) |
| Smart list query | No results | API cannot query smart lists |
| Clear incomplete | Error | Must complete task before clearing |
