package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// runCommand executes a command and returns stdout, stderr, error.
func runCommand(name string, args ...string) (string, string, error) {
	cmd := exec.Command(name, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

// runCommandPrint executes a command with output to terminal.
func runCommandPrint(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// parseArgs parses --key value style arguments into a map.
func parseArgs(args []string) map[string]string {
	result := make(map[string]string)
	for i := 0; i < len(args); i++ {
		if strings.HasPrefix(args[i], "--") {
			key := strings.TrimPrefix(args[i], "--")
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "--") {
				result[key] = args[i+1]
				i++
			} else {
				result[key] = "true"
			}
		}
	}
	return result
}

// requireArg checks if a required argument exists.
func requireArg(args map[string]string, key, usage string) string {
	if val, ok := args[key]; ok && val != "" {
		return val
	}
	fmt.Fprintf(os.Stderr, "Error: --%s is required. %s\n", key, usage)
	os.Exit(1)
	return ""
}

// getArg gets an optional argument with default.
func getArg(args map[string]string, key, defaultVal string) string {
	if val, ok := args[key]; ok && val != "" {
		return val
	}
	return defaultVal
}

// ensureDir creates directory if it doesn't exist.
func ensureDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}

// fileExists checks if file exists.
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// getBaseName returns filename without extension.
func getBaseName(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	return strings.TrimSuffix(base, ext)
}

// which checks if a command exists in PATH.
func which(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// checkDependency checks if required tool is installed.
func checkDependency(name, installHint string) {
	if !which(name) {
		fmt.Fprintf(os.Stderr, "Error: %s not found. Install with: %s\n", name, installHint)
		os.Exit(1)
	}
}
