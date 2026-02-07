// Content processor CLI - standalone tools for content acquisition and processing
// Usage: content-processor <command> [args]
package main

import (
	"fmt"
	"os"
)

var commands = map[string]func([]string){
	"yt-transcript":   ytTranscript,
	"yt-audio":        ytAudio,
	"epub-to-md":      epubToMd,
	"md-cleanup":      mdCleanup,
	"md-split":        mdSplit,
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
	fmt.Println(`Usage: content-processor <command> [args]

YouTube:
  yt-transcript  --url <url> --output-dir <dir> [--language <code>]
  yt-audio       --url <url> --output-dir <dir>

Books:
  epub-to-md     --input <file.epub> [--output <file.md>]
  md-cleanup     --input <file.md> [--output <file.md>]
  md-split       --input <file.md> --level <1-6> [--pattern <regex>]`)
}
