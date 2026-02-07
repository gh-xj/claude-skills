---
name: go-script-style
description: Go CLI script writing conventions and patterns. Use when creating Go automation scripts, CLI tools, or converting shell/Python scripts to Go. Covers flat CLI structure, samber/lo usage, and osascript helpers.
---

# Go Script Style Guide

> **Detailed patterns:** `references/patterns.md`

## Requirements

- **Go 1.25+** required for all new projects

## When to Use

- Creating new Go CLI tools
- Converting bash/Python scripts to Go
- Adding commands to existing Go CLI
- Reviewing Go automation code

## Project Structure

```
project/
├── main.go           # CLI entry, command routing
├── helpers.go        # Shared utilities
├── cmd_*.go          # One file per command
├── Taskfile.yml      # Build automation
├── go.mod
└── bin/              # Compiled binaries
```

## Quick Start: New Command

1. Create `cmd_<name>.go`:

```go
package main

import "fmt"

// commandName does X, Y, Z.
func commandName(args []string) {
    // implementation
}
```

2. Register in `main.go`:

```go
var commands = map[string]func([]string){
    "command-name": commandName,
    // ...
}
```

3. Build: `task build` or `go build -o bin/cli .`

## File Conventions

| File         | Purpose                                    |
| ------------ | ------------------------------------------ |
| `main.go`    | Entry point, command map, usage            |
| `helpers.go` | Shared functions (no business logic)       |
| `cmd_*.go`   | One command per file, named after function |

## Code Style

### Imports

```go
import (
    "fmt"           // stdlib first
    "os"
    "os/exec"

    "github.com/samber/lo"  // external after blank line
)
```

### Function Signatures

```go
// Commands take args slice
func myCommand(args []string) { }

// Helpers are generic
func getRepeatArg(args []string, multiplier int) int { }
```

### Common Helpers

| Helper                     | Purpose                          |
| -------------------------- | -------------------------------- |
| `parseArgs`                | Parse `--key value` flags to map |
| `requireArg` / `getArg`    | Required/optional flag access    |
| `checkDependency`          | Verify external tool exists      |
| `fileExists` / `ensureDir` | File operations                  |
| `runOsascript`             | macOS automation                 |

See `references/patterns.md` for implementations.

### Use samber/lo

```go
// Ternary
fallback := lo.Ternary(condition, "yes", "no")

// Filter
valid := lo.Filter(items, func(i Item, _ int) bool { return i.Valid })

// Map
results := lo.Map(items, func(i Item, _ int) string { return i.Name })
```

## Taskfile.yml

```yaml
version: "3"

tasks:
  build:
    desc: Build CLI binary
    cmds:
      - go build -o bin/cli .
    sources: ["*.go", "go.mod"]
    method: checksum

  run:
    desc: Run without building
    cmds:
      - go run . {{.CLI_ARGS}}
```

## Naming

| Element      | Convention      | Example          |
| ------------ | --------------- | ---------------- |
| Files        | `cmd_<name>.go` | `cmd_deploy.go`  |
| Functions    | camelCase       | `deployServer`   |
| CLI commands | kebab-case      | `deploy-server`  |
| Constants    | camelCase       | `defaultTimeout` |

## Adding a New Command

1. Create `cmd_new_feature.go`
2. Add function `newFeature(args []string)`
3. Register in `main.go` commands map
4. Update `printUsage()`
5. Rebuild: `task build`
