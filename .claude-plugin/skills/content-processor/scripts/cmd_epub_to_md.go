package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// epubToMd converts EPUB to markdown using pandoc.
// Usage: epub-to-md --input <file.epub> --output <file.md>
func epubToMd(args []string) {
	checkDependency("pandoc", "brew install pandoc")

	parsed := parseArgs(args)
	input := requireArg(parsed, "input", "Input EPUB file")

	if !fileExists(input) {
		fmt.Fprintf(os.Stderr, "Error: input file not found: %s\n", input)
		os.Exit(1)
	}

	// Default output: same name with .md extension
	output := getArg(parsed, "output", "")
	if output == "" {
		dir := filepath.Dir(input)
		base := getBaseName(input)
		output = filepath.Join(dir, base+".md")
	}

	// Build pandoc command
	pandocArgs := []string{
		"-f", "epub",
		"-t", "markdown",
		"--wrap=none",
		"-o", output,
		input,
	}

	fmt.Printf("Converting: %s -> %s\n", input, output)
	if err := runCommandPrint("pandoc", pandocArgs...); err != nil {
		fmt.Fprintf(os.Stderr, "Error converting EPUB: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Conversion complete: %s\n", output)
	fmt.Println("\nNext step: clean Calibre artifacts")
	fmt.Printf("  content-processor md-cleanup --input %s --output %s-cleaned.md\n",
		output, getBaseName(output))
}
