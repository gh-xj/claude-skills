package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// mdCleanup cleans Calibre artifacts from markdown.
// Usage: md-cleanup --input <file.md> --output <file.md>
func mdCleanup(args []string) {
	parsed := parseArgs(args)
	input := requireArg(parsed, "input", "Input markdown file")

	if !fileExists(input) {
		fmt.Fprintf(os.Stderr, "Error: input file not found: %s\n", input)
		os.Exit(1)
	}

	// Default output: same name with -cleaned suffix
	output := getArg(parsed, "output", "")
	if output == "" {
		dir := filepath.Dir(input)
		base := getBaseName(input)
		ext := filepath.Ext(input)
		output = filepath.Join(dir, base+"-cleaned"+ext)
	}

	content, err := os.ReadFile(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	cleaned := cleanMarkdown(string(content))

	if err := os.WriteFile(output, []byte(cleaned), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Cleaned: %s -> %s\n", input, output)
}

// cleanMarkdown removes Calibre artifacts and normalizes markdown.
func cleanMarkdown(content string) string {
	// Remove Calibre class attributes: {.calibre1}, {.calibre2}, etc.
	calibreClass := regexp.MustCompile(`\{\.calibre\d*\}`)
	content = calibreClass.ReplaceAllString(content, "")

	// Remove Calibre ID anchors: {#part0001.xhtml}, {#part0002_split_000}, etc.
	calibreAnchor := regexp.MustCompile(`\{#[^}]+\}`)
	content = calibreAnchor.ReplaceAllString(content, "")

	// Remove standalone anchor links: [](#part0001.xhtml)
	anchorLink := regexp.MustCompile(`\[\]\(#[^)]+\)`)
	content = anchorLink.ReplaceAllString(content, "")

	// Remove image references with Calibre paths
	calibreImg := regexp.MustCompile(`!\[\]\([^)]*images/[^)]+\)`)
	content = calibreImg.ReplaceAllString(content, "")

	// Clean up multiple blank lines (more than 2)
	multiBlank := regexp.MustCompile(`\n{3,}`)
	content = multiBlank.ReplaceAllString(content, "\n\n")

	// Remove trailing whitespace from lines
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t")
	}
	content = strings.Join(lines, "\n")

	// Ensure file ends with single newline
	content = strings.TrimRight(content, "\n") + "\n"

	return content
}
