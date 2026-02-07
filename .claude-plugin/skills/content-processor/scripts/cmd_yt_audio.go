package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// ytAudio downloads YouTube audio for whisper transcription.
// Usage: yt-audio --url <url> --output-dir <dir>
func ytAudio(args []string) {
	checkDependency("yt-dlp", "brew install yt-dlp")

	parsed := parseArgs(args)
	url := requireArg(parsed, "url", "YouTube video URL")
	outputDir := getArg(parsed, "output-dir", "./yt-audio")

	if err := ensureDir(outputDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output dir: %v\n", err)
		os.Exit(1)
	}

	// Build yt-dlp command for audio download
	ytArgs := []string{
		"-x",                   // Extract audio
		"--audio-format", "wav", // WAV for whisper compatibility
		"--audio-quality", "0", // Best quality
		"--output", filepath.Join(outputDir, "%(id)s.%(ext)s"),
		url,
	}

	fmt.Printf("Downloading audio: %s\n", url)
	if err := runCommandPrint("yt-dlp", ytArgs...); err != nil {
		fmt.Fprintf(os.Stderr, "Error downloading audio: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Audio download complete")
	fmt.Println("\nTo transcribe with whisper:")
	fmt.Println("  whisper-cli --model ~/.whisper/models/ggml-base.bin \\")
	fmt.Println("    --language en --output-txt --output-file <output> <input.wav>")
}
