package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// mdSplit splits markdown into chapters based on heading level.
// Usage: md-split --input <file.md> --level <1|2> [--pattern <regex>]
func mdSplit(args []string) {
	parsed := parseArgs(args)
	input := requireArg(parsed, "input", "Input markdown file")
	levelStr := getArg(parsed, "level", "2")
	pattern := getArg(parsed, "pattern", ".*")

	if !fileExists(input) {
		fmt.Fprintf(os.Stderr, "Error: input file not found: %s\n", input)
		os.Exit(1)
	}

	level, err := strconv.Atoi(levelStr)
	if err != nil || level < 1 || level > 6 {
		fmt.Fprintf(os.Stderr, "Error: level must be 1-6\n")
		os.Exit(1)
	}

	content, err := os.ReadFile(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Create output directory
	base := getBaseName(input)
	outputDir := base + "-split_chapters"
	if err := ensureDir(outputDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output dir: %v\n", err)
		os.Exit(1)
	}

	chapters := splitByHeading(string(content), level, pattern)

	if len(chapters) == 0 {
		fmt.Fprintf(os.Stderr, "No chapters found at level %d with pattern '%s'\n", level, pattern)
		os.Exit(1)
	}

	// Write chapters
	for i, ch := range chapters {
		filename := fmt.Sprintf("%03d_%s.md", i+1, sanitizeFilename(ch.Title))
		filepath := filepath.Join(outputDir, filename)
		if err := os.WriteFile(filepath, []byte(ch.Content), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", filename, err)
			continue
		}
		fmt.Printf("  %s\n", filename)
	}

	// Create README with reading times
	createReadme(outputDir, chapters)

	fmt.Printf("\n✅ Split into %d chapters: %s/\n", len(chapters), outputDir)
}

type chapter struct {
	Title   string
	Content string
}

// splitByHeading splits content by heading level matching pattern.
func splitByHeading(content string, level int, pattern string) []chapter {
	// Build heading regex: ^#{level} (pattern)
	hashes := strings.Repeat("#", level)
	headingRe := regexp.MustCompile(fmt.Sprintf(`(?m)^%s\s+(%s.*)$`, hashes, pattern))
	patternRe := regexp.MustCompile(pattern)

	lines := strings.Split(content, "\n")
	var chapters []chapter
	var currentChapter *chapter
	var currentContent []string

	for _, line := range lines {
		if headingRe.MatchString(line) && patternRe.MatchString(line) {
			// Save previous chapter
			if currentChapter != nil {
				currentChapter.Content = strings.Join(currentContent, "\n")
				chapters = append(chapters, *currentChapter)
			}

			// Start new chapter
			title := strings.TrimLeft(line, "# ")
			currentChapter = &chapter{Title: title}
			currentContent = []string{line}
		} else if currentChapter != nil {
			currentContent = append(currentContent, line)
		}
	}

	// Save last chapter
	if currentChapter != nil {
		currentChapter.Content = strings.Join(currentContent, "\n")
		chapters = append(chapters, *currentChapter)
	}

	return chapters
}

// sanitizeFilename makes a string safe for filenames.
func sanitizeFilename(s string) string {
	// Remove or replace unsafe characters
	s = strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z':
			return r
		case r >= 'A' && r <= 'Z':
			return r
		case r >= '0' && r <= '9':
			return r
		case r == ' ' || r == '-' || r == '_':
			return '_'
		default:
			return -1
		}
	}, s)

	// Limit length
	if len(s) > 50 {
		s = s[:50]
	}
	return strings.Trim(s, "_")
}

// createReadme generates README.md with reading times.
func createReadme(outputDir string, chapters []chapter) {
	var sb strings.Builder
	sb.WriteString("# Chapters\n\n")
	sb.WriteString("| # | Chapter | Reading Time |\n")
	sb.WriteString("|---|---------|-------------|\n")

	for i, ch := range chapters {
		filename := fmt.Sprintf("%03d_%s.md", i+1, sanitizeFilename(ch.Title))
		// Estimate: ~200 words per minute, ~5 chars per word
		sizeKB := float64(len(ch.Content)) / 1024.0
		readingMin := int(sizeKB * 0.4)
		if readingMin < 1 {
			readingMin = 1
		}
		sb.WriteString(fmt.Sprintf("| %d | [%s](%s) | %d min |\n",
			i+1, ch.Title, filename, readingMin))
	}

	readme := filepath.Join(outputDir, "README.md")
	os.WriteFile(readme, []byte(sb.String()), 0644)
}
