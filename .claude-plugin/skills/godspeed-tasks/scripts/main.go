// Godspeed Task Manager CLI
package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/samber/lo"
)

type opts struct {
	status, list, format string
	limit                int
	due, recent, done    bool
}

func parseArgs(args []string) (cmd string, pos []string, o opts) {
	o.format = "table"
	if len(args) > 0 {
		cmd, args = args[0], args[1:]
	}

	flagParsers := map[string]func(string){
		"--status": func(v string) { o.status = v },
		"--list":   func(v string) { o.list = v },
		"--limit":  func(v string) { fmt.Sscanf(v, "%d", &o.limit) },
		"--format": func(v string) { o.format = v },
	}

	boolFlags := map[string]*bool{
		"--due": &o.due, "--recent": &o.recent, "--completed": &o.done,
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if fn, ok := flagParsers[arg]; ok && i+1 < len(args) {
			i++
			fn(args[i])
		} else if ptr, ok := boolFlags[arg]; ok {
			*ptr = true
		} else if !strings.HasPrefix(arg, "-") {
			pos = append(pos, arg)
		}
	}
	return
}

func die(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
	os.Exit(1)
}

func main() {
	cmd, pos, o := parseArgs(os.Args[1:])
	asJSON := o.format == "json"

	if cmd == "" || lo.Contains([]string{"help", "-h", "--help"}, cmd) {
		fmt.Println(`Usage: godspeed <command> [options]

Commands:
  stats              Task statistics
  list               List tasks (--status, --list, --due, --recent, --completed)
  search <keyword>   Search tasks
  get <id>           Get task details
  lists              Show all lists

Options:
  --status <incomplete|complete>
  --list <name|id>
  --limit <n>
  --format <table|json>
  --due, --recent, --completed`)
		return
	}

	client, err := NewClient()
	if err != nil {
		die("%v", err)
	}
	defer client.Close()

	ctx := context.Background()

	switch cmd {
	case "stats":
		s, err := client.Stats(ctx)
		if err != nil {
			die("%v", err)
		}
		printStats(s, asJSON)

	case "list":
		listID := o.list
		if listID != "" && len(listID) != 36 {
			l, err := client.ListByName(ctx, listID)
			if err != nil {
				die("%v", err)
			}
			listID = l.ID
		}

		tasks, err := client.Query(ctx, QueryOpts{
			Status: o.status, ListID: listID, Due: o.due,
			Recent: o.recent, Done: o.done, Limit: o.limit,
		})
		if err != nil {
			die("%v", err)
		}
		printTasks(tasks, asJSON)

	case "search":
		if len(pos) == 0 {
			die("search requires keyword")
		}
		tasks, err := client.Query(ctx, QueryOpts{Keyword: pos[0], Limit: o.limit})
		if err != nil {
			die("%v", err)
		}
		printTasks(tasks, asJSON)

	case "get":
		if len(pos) == 0 {
			die("get requires task id")
		}
		t, err := client.Get(ctx, pos[0])
		if err != nil {
			die("%v", err)
		}
		printTask(t, asJSON)

	case "lists":
		lists, err := client.Lists(ctx)
		if err != nil {
			die("%v", err)
		}
		printLists(lists, asJSON)

	default:
		die("unknown command: %s", cmd)
	}
}
