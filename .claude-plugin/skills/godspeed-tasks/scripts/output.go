package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/samber/lo"
)

func output[T any](v T, asJSON bool, tableFn func(T)) {
	if asJSON {
		e := json.NewEncoder(os.Stdout)
		e.SetIndent("", "  ")
		e.Encode(v)
		return
	}
	tableFn(v)
}

func printStats(s *TaskStats, asJSON bool) {
	output(s, asJSON, func(s *TaskStats) {
		fmt.Printf("Total: %d | Incomplete: %d | Completed: %d\n", s.Total, s.Incomplete, s.Completed)
	})
}

func printTask(t *Task, asJSON bool) {
	output(t, asJSON, func(t *Task) {
		fmt.Printf("ID:     %s\nTitle:  %s\nStatus: %s\nList:   %s\n", t.ID, t.Title, t.Status, t.Display())
		if t.Notes != "" {
			fmt.Printf("Notes:  %s\n", truncate(t.Notes, 80))
		}
		if t.TimelessDueAt != "" {
			fmt.Printf("Due:    %s\n", t.TimelessDueAt)
		}
	})
}

func printTasks(tasks []Task, asJSON bool) {
	output(tasks, asJSON, func(tasks []Task) {
		if len(tasks) == 0 {
			fmt.Println("No tasks found.")
			return
		}
		lo.ForEach(tasks, func(t Task, _ int) {
			indent := strings.Repeat("  ", t.IndentLevel)
			due := lo.Ternary(t.TimelessDueAt != "", " due:"+t.TimelessDueAt, "")
			fmt.Printf("%s%s %s (%s)%s\n", indent, t.Checkbox(), t.Title, t.Display(), due)
		})
	})
}

func printLists(lists []ListWithCount, asJSON bool) {
	output(lists, asJSON, func(lists []ListWithCount) {
		if len(lists) == 0 {
			fmt.Println("No lists found.")
			return
		}
		lo.ForEach(lists, func(l ListWithCount, _ int) {
			fmt.Printf("%-20s %3d incomplete  %s\n", l.List.Name, l.IncompleteCount, l.List.ID)
		})
	})
}

func truncate(s string, n int) string {
	return lo.Ternary(len(s) <= n, s, s[:n-3]+"...")
}
