# Configuration

## Directory Structure

```
<project>/
├── sources/
│   ├── yt-transcripts/          # YouTube transcripts
│   ├── yt-audio-downloads/      # Downloaded audio (fallback)
│   ├── yt-audio-transcripts/    # Whisper outputs
│   └── web-content/             # WebFetch saves
├── <book>-split_chapters/       # Split chapters
│   ├── 001_*.md
│   ├── README.md                # Reading times
│   └── ai-summary/
│       ├── PROGRESS.md
│       └── *_summary.md
└── summaries/                   # URL-based summaries
    └── <slug>_summary.md
```

## CLI Commands

```bash
# YouTube
content-processor yt-transcript --url "<url>" --output-dir <dir> [--language <code>]
content-processor yt-audio --url "<url>" --output-dir <dir>

# Books
content-processor epub-to-md --input "<file>.epub" --output "<file>.md"
content-processor md-cleanup --input "<file>.md" --output "<file>-cleaned.md"
content-processor md-split --input "<file>.md" --level <level> [--pattern "<regex>"]

# Whisper (audio fallback)
whisper-cli --model ~/.whisper/models/ggml-base.bin --language <code> \
  --output-txt --output-file <output> <input.wav>
```

**CLI Location:** `~/.claude/skills/content-processor/scripts/content-processor`

## Defaults

| Setting             | Default                           | Override            |
| ------------------- | --------------------------------- | ------------------- |
| Transcript language | `en`                              | `--language <code>` |
| YT output dir       | `./yt-transcripts`                | `--output-dir`      |
| Audio output dir    | `./yt-audio`                      | `--output-dir`      |
| Split heading level | `2`                               | `--level <1-6>`     |
| Whisper model       | `~/.whisper/models/ggml-base.bin` | `--model`           |

## Reading Time Formula

```
time_min = file_size_KB × 0.4 (rounded, min 1)
```
