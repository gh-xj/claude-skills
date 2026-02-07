package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/samber/lo"
)

type Client struct {
	db     *sql.DB
	userID int64
}

func NewClient() (*Client, error) {
	dbPath, err := findDatabase()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", dbPath+"?mode=ro")
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping db: %w", err)
	}

	userID, err := findUserID(db)
	if err != nil {
		db.Close()
		return nil, err
	}

	return &Client{db: db, userID: userID}, nil
}

func (c *Client) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

func findDatabase() (string, error) {
	home, _ := os.UserHomeDir()
	paths := []string{
		filepath.Join(home, "Library/Application Support/Godspeed/godspeed-db-setapp.sqlite"),
		filepath.Join(home, "Library/Application Support/Godspeed/godspeed-db.sqlite"),
	}
	path, found := lo.Find(paths, func(p string) bool {
		_, err := os.Stat(p)
		return err == nil
	})
	if !found {
		return "", fmt.Errorf("godspeed database not found")
	}
	return path, nil
}

func findUserID(db *sql.DB) (int64, error) {
	var id int64
	err := db.QueryRow(`
		SELECT user_id FROM todo_items WHERE user_id > 0
		GROUP BY user_id ORDER BY COUNT(*) DESC LIMIT 1
	`).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("detect user: %w", err)
	}
	return id, nil
}

// --- Queries ---

const baseSelect = `
	SELECT t.id, t.title, t.notes, t.list_id, t.completed_at,
	       t.created_at, t.updated_at, t.indent_level, t.order_index,
	       t.timeless_due_at, l.name
	FROM todo_items t LEFT JOIN lists l ON t.list_id = l.id
`

func (c *Client) Stats(ctx context.Context) (*TaskStats, error) {
	var s TaskStats
	err := c.db.QueryRowContext(ctx, `
		SELECT COUNT(*),
		       SUM(CASE WHEN completed_at IS NULL THEN 1 ELSE 0 END),
		       SUM(CASE WHEN completed_at IS NOT NULL THEN 1 ELSE 0 END)
		FROM todo_items WHERE user_id = ?
	`, c.userID).Scan(&s.Total, &s.Incomplete, &s.Completed)
	return &s, err
}

func (c *Client) Get(ctx context.Context, id string) (*Task, error) {
	row := c.db.QueryRowContext(ctx, baseSelect+`WHERE t.id = ? AND t.user_id = ?`, id, c.userID)
	return scanTask(row)
}

func (c *Client) Query(ctx context.Context, o QueryOpts) ([]Task, error) {
	where := []string{"t.user_id = ?"}
	args := []any{c.userID}

	if o.Status == "incomplete" || o.Due || o.Recent {
		where = append(where, "t.completed_at IS NULL")
	} else if o.Status == "complete" || o.Done {
		where = append(where, "t.completed_at IS NOT NULL")
	}

	if o.Due {
		where = append(where, "(t.due_at IS NOT NULL OR t.timeless_due_at IS NOT NULL)")
	}

	if o.ListID != "" {
		where = append(where, "t.list_id = ?")
		args = append(args, o.ListID)
	}

	if o.Keyword != "" {
		where = append(where, "t.title LIKE ?")
		args = append(args, "%"+o.Keyword+"%")
	}

	orderBy := lo.Switch[bool, string](true).
		Case(o.Due, "COALESCE(t.timeless_due_at, t.due_at)").
		Case(o.Recent || o.Keyword != "", "t.updated_at DESC").
		Case(o.Done, "t.completed_at DESC").
		Default("t.order_index")

	q := fmt.Sprintf("%s WHERE %s ORDER BY %s LIMIT %d",
		baseSelect, strings.Join(where, " AND "), orderBy, lo.CoalesceOrEmpty(o.Limit, 100))

	return c.queryTasks(ctx, q, args...)
}

func (c *Client) ListByName(ctx context.Context, name string) (*TaskList, error) {
	var l TaskList
	var createdAt, updatedAt string
	err := c.db.QueryRowContext(ctx, `
		SELECT id, name, created_at, updated_at FROM lists
		WHERE LOWER(name) = LOWER(?) AND archived_at IS NULL LIMIT 1
	`, name).Scan(&l.ID, &l.Name, &createdAt, &updatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("list not found: %s", name)
	}
	if err != nil {
		return nil, err
	}

	l.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	l.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	return &l, nil
}

func (c *Client) Lists(ctx context.Context) ([]ListWithCount, error) {
	rows, err := c.db.QueryContext(ctx, `
		SELECT l.id, l.name, l.created_at, l.updated_at,
		       COUNT(CASE WHEN t.completed_at IS NULL THEN 1 END)
		FROM lists l
		LEFT JOIN todo_items t ON l.id = t.list_id AND t.user_id = ?
		WHERE l.archived_at IS NULL
		GROUP BY l.id ORDER BY l.order_index
	`, c.userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []ListWithCount
	for rows.Next() {
		var lc ListWithCount
		var createdAt, updatedAt string
		if err := rows.Scan(&lc.List.ID, &lc.List.Name, &createdAt, &updatedAt, &lc.IncompleteCount); err != nil {
			return nil, err
		}
		lc.List.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		lc.List.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
		results = append(results, lc)
	}
	return results, rows.Err()
}

func (c *Client) queryTasks(ctx context.Context, q string, args ...any) ([]Task, error) {
	rows, err := c.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		t, err := scanTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, *t)
	}
	return tasks, rows.Err()
}

type scanner interface{ Scan(dest ...any) error }

func scanTask(s scanner) (*Task, error) {
	var t Task
	var notes, completedAt, dueAt, listName sql.NullString
	var createdAt, updatedAt string

	err := s.Scan(&t.ID, &t.Title, &notes, &t.ListID, &completedAt,
		&createdAt, &updatedAt, &t.IndentLevel, &t.OrderIndex, &dueAt, &listName)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("not found")
	}
	if err != nil {
		return nil, err
	}

	t.Notes, t.TimelessDueAt, t.ListName = notes.String, dueAt.String, listName.String
	t.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	t.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	t.Status = lo.Ternary(completedAt.Valid && completedAt.String != "", "completed", "pending")

	if completedAt.Valid && completedAt.String != "" {
		if ts, err := time.Parse(time.RFC3339, completedAt.String); err == nil {
			t.CompletedAt = &ts
		}
	}

	return &t, nil
}
