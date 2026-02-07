package main

import (
	"time"

	"github.com/samber/lo"
)

type Task struct {
	ID            string     `json:"id"`
	Title         string     `json:"title"`
	Notes         string     `json:"notes,omitempty"`
	Status        string     `json:"status"`
	ListID        string     `json:"list_id,omitempty"`
	ListName      string     `json:"list_name,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	CompletedAt   *time.Time `json:"completed_at,omitempty"`
	TimelessDueAt string     `json:"timeless_due_at,omitempty"`
	IndentLevel   int        `json:"indent_level"`
	OrderIndex    float64    `json:"order_index"`
}

func (t Task) Checkbox() string { return lo.Ternary(t.Status == "completed", "[x]", "[ ]") }
func (t Task) Display() string  { return lo.CoalesceOrEmpty(t.ListName, "Inbox") }

type TaskList struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TaskStats struct {
	Total      int `json:"total"`
	Incomplete int `json:"incomplete"`
	Completed  int `json:"completed"`
}

type ListWithCount struct {
	List            TaskList `json:"list"`
	IncompleteCount int      `json:"incomplete_count"`
}

type QueryOpts struct {
	Status  string // "incomplete" | "complete" | ""
	ListID  string
	Keyword string
	Due     bool
	Recent  bool
	Done    bool
	Limit   int
}
