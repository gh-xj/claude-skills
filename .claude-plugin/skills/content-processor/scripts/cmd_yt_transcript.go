package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ytTranscript downloads YouTube transcript using yt-dlp.
// Usage: yt-transcript --url <url> --output-dir <dir> [--language <code>]
func ytTranscript(args []string) {
	checkDependency("yt-dlp", "brew install yt-dlp")

	parsed := parseArgs(args)
	url := requireArg(parsed, "url", "YouTube video URL")
	outputDir := getArg(parsed, "output-dir", "./yt-transcripts")
	language := getArg(parsed, "language", "en")

	if err := ensureDir(outputDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output dir: %v\n", err)
		os.Exit(1)
	}

	// Build yt-dlp command for subtitle download
	ytArgs := []string{
		"--write-subs",
		"--write-auto-subs",
		"--sub-lang", language,
		"--sub-format", "vtt",
		"--skip-download",
		"--output", filepath.Join(outputDir, "%(id)s.%(ext)s"),
		url,
	}

	fmt.Printf("Downloading transcript: %s\n", url)
	if err := runCommandPrint("yt-dlp", ytArgs...); err != nil {
		fmt.Fprintf(os.Stderr, "Error downloading transcript: %v\n", err)
		os.Exit(1)
	}

	// Convert VTT to plain text markdown
	files, _ := filepath.Glob(filepath.Join(outputDir, "*.vtt"))
	for _, vttFile := range files {
		mdFile := strings.TrimSuffix(vttFile, ".vtt") + ".md"
		if err := convertVttToMd(vttFile, mdFile); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to convert %s: %v\n", vttFile, err)
		} else {
			fmt.Printf("Created: %s\n", mdFile)
			os.Remove(vttFile) // Clean up VTT
		}
	}

	fmt.Println("✅ Transcript download complete")
}

// convertVttToMd converts VTT subtitle file to clean markdown.
func convertVttToMd(vttFile, mdFile string) error {
	content, err := os.ReadFile(vttFile)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	var result []string
	var lastLine string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip VTT header, timestamps, and empty lines
		if line == "" || line == "WEBVTT" || strings.Contains(line, "-->") {
			continue
		}
		// Skip numeric cue identifiers
		if isNumeric(line) {
			continue
		}
		// Skip HTML-like tags
		if strings.HasPrefix(line, "<") && strings.HasSuffix(line, ">") {
			continue
		}

		// Remove duplicate lines (common in auto-generated subs)
		if line != lastLine {
			result = append(result, line)
			lastLine = line
		}
	}

	// Join into paragraphs
	text := strings.Join(result, " ")
	// Clean up multiple spaces
	for strings.Contains(text, "  ") {
		text = strings.ReplaceAll(text, "  ", " ")
	}

	// Write markdown
	mdContent := fmt.Sprintf("# Transcript\n\n%s\n", text)
	return os.WriteFile(mdFile, []byte(mdContent), 0644)
}

func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(s) > 0
}
