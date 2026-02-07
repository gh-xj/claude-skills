# Go Script Patterns

## CLI Entry Pattern

```go
package main

import (
    "fmt"
    "os"
)

var commands = map[string]func([]string){
    "cmd-one": cmdOne,
    "cmd-two": cmdTwo,
}

func main() {
    if len(os.Args) < 2 {
        printUsage()
        os.Exit(1)
    }

    cmd := os.Args[1]
    args := os.Args[2:]

    if fn, ok := commands[cmd]; ok {
        fn(args)
    } else {
        fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
        printUsage()
        os.Exit(1)
    }
}

func printUsage() {
    fmt.Println(`Usage: cli <command> [args]

Commands:
  cmd-one    Description of cmd-one
  cmd-two    Description of cmd-two`)
}
```

## Helper Patterns

### osascript Execution (macOS)

```go
func runOsascript(script string) string {
    cmd := exec.Command("osascript", "-e", script)
    out, _ := cmd.Output()
    return strings.TrimSpace(string(out))
}
```

### Shell Command with Error

```go
func runCommand(name string, args ...string) (string, error) {
    cmd := exec.Command(name, args...)
    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr

    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("%w: %s", err, stderr.String())
    }
    return stdout.String(), nil
}
```

### Argument Parsing

#### Positional Args

```go
func getRepeatArg(args []string, multiplier int) int {
    if len(args) > 0 && args[0] != "" {
        if n, err := strconv.Atoi(args[0]); err == nil {
            return n * multiplier
        }
    }
    return multiplier
}

func getStringArg(args []string, index int, defaultVal string) string {
    if len(args) > index && args[index] != "" {
        return args[index]
    }
    return defaultVal
}
```

#### Flag-Style Args (--key value)

```go
func parseArgs(args []string) map[string]string {
    result := make(map[string]string)
    for i := 0; i < len(args); i++ {
        if strings.HasPrefix(args[i], "--") {
            key := strings.TrimPrefix(args[i], "--")
            if i+1 < len(args) && !strings.HasPrefix(args[i+1], "--") {
                result[key] = args[i+1]
                i++
            } else {
                result[key] = "true"  // boolean flag
            }
        }
    }
    return result
}

func requireArg(args map[string]string, key, usage string) string {
    if val, ok := args[key]; ok && val != "" {
        return val
    }
    fmt.Fprintf(os.Stderr, "Error: --%s is required. %s\n", key, usage)
    os.Exit(1)
    return ""
}

func getArg(args map[string]string, key, defaultVal string) string {
    if val, ok := args[key]; ok && val != "" {
        return val
    }
    return defaultVal
}
```

### Dependency Checking

```go
func which(cmd string) bool {
    _, err := exec.LookPath(cmd)
    return err == nil
}

func checkDependency(name, installHint string) {
    if !which(name) {
        fmt.Fprintf(os.Stderr, "Error: %s not found. Install with: %s\n", name, installHint)
        os.Exit(1)
    }
}

// Usage
checkDependency("pandoc", "brew install pandoc")
checkDependency("yt-dlp", "brew install yt-dlp")
```

### File/Path Helpers

```go
func fileExists(path string) bool {
    _, err := os.Stat(path)
    return err == nil
}

func ensureDir(dir string) error {
    return os.MkdirAll(dir, 0755)
}

func getBaseName(path string) string {
    base := filepath.Base(path)
    ext := filepath.Ext(base)
    return strings.TrimSuffix(base, ext)
}
```

### Conditional Check with Action

```go
func requireApp(appName string, fallbackAction func()) bool {
    fallbackCode := ""
    if fallbackAction != nil {
        // build fallback into script
    }
    result := runOsascript(fmt.Sprintf(`
        tell application "System Events"
            set frontApp to name of first application process whose frontmost is true
            if frontApp is not "%s" then
                %s
                return "false"
            end if
            return "true"
        end tell
    `, appName, fallbackCode))
    return result == "true"
}
```

## samber/lo Patterns

### Ternary for Conditional Values

```go
import "github.com/samber/lo"

// Instead of if/else for simple values
timeout := lo.Ternary(debug, 60, 10)
prefix := lo.Ternary(verbose, "[DEBUG] ", "")
```

### Filter and Map

```go
// Filter valid items
valid := lo.Filter(items, func(item Item, _ int) bool {
    return item.Status == "active"
})

// Transform items
names := lo.Map(items, func(item Item, _ int) string {
    return item.Name
})

// Filter then map
activeNames := lo.Map(
    lo.Filter(items, func(i Item, _ int) bool { return i.Active }),
    func(i Item, _ int) string { return i.Name },
)
```

### Find and Contains

```go
// Find first match
found, ok := lo.Find(items, func(i Item) bool { return i.ID == targetID })

// Check existence
exists := lo.Contains(names, "target")

// Index of
idx := lo.IndexOf(items, target)
```

### Grouping and Chunking

```go
// Group by key
grouped := lo.GroupBy(items, func(i Item) string { return i.Category })

// Process in batches
batches := lo.Chunk(items, 10)
for _, batch := range batches {
    processBatch(batch)
}
```

### Error Handling

```go
// Compact removes zero values
clean := lo.Compact([]string{"a", "", "b"})  // ["a", "b"]

// Uniq removes duplicates
unique := lo.Uniq([]string{"a", "b", "a"})   // ["a", "b"]

// Must panics on error (use sparingly)
value := lo.Must(strconv.Atoi("123"))
```

## Command File Template

```go
package main

import "fmt"

// myCommand does X when Y.
// Multiplier: 1 (no multiplier)
func myCommand(args []string) {
    // Optional: check preconditions
    if !requireApp("TargetApp", nil) {
        return
    }

    // Parse arguments
    repeat := getRepeatArg(args, 1)

    // Execute
    runOsascript(fmt.Sprintf(`
        tell application "System Events"
            repeat %d times
                -- actions here
            end repeat
        end tell
    `, repeat))
}
```

## Testing Pattern

```go
// In main.go or helpers.go
func runWithArgs(cmdArgs []string) {
    os.Args = append([]string{"cli"}, cmdArgs...)
    main()
}

// Manual test
// go run . my-command arg1 arg2
```

## Build Optimization

```yaml
# Taskfile.yml with checksum-based rebuild
tasks:
  build:
    cmds:
      - go build -o bin/cli .
    sources:
      - "*.go"
      - go.mod
      - go.sum
    generates:
      - bin/cli
    method: checksum
```
