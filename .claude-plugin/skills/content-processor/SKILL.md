---
name: content-processor
description: Process any content source (YouTube, websites, EPUBs, markdown) into structured summaries. Use when user provides URLs for transcription, has books/EPUBs to process, or wants batch summarization of any content type.
---

# Content Processor

Unified content acquisition and summarization.

## Quick Reference

| Input            | Detection                 | Workflow                                       |
| ---------------- | ------------------------- | ---------------------------------------------- |
| YouTube URL      | `youtube.com`, `youtu.be` | [URL Workflow](#url-workflow)                  |
| YouTube playlist | `list=` in URL            | [URL Workflow](#url-workflow) (batch)          |
| Website URL      | Any other URL             | [URL Workflow](#url-workflow) (WebFetch)       |
| `.epub` file     | File extension            | [Book Workflow](#book-workflow)                |
| `.md` file       | File extension            | [Book Workflow](#book-workflow) (skip convert) |

**First question:** What language is the content in?

**References:**

- Config & CLI: `references/config.md`
- Output formats: `references/formats.md`
- Sub-agent rules: `references/orchestration.md`

---

## URL Workflow

### YouTube Single Video

1. Ask language → use `--language` flag
2. Download transcript:
   ```bash
   content-processor yt-transcript --url "<url>" --output-dir ./sources/yt-transcripts
   ```
3. If HTTP 429 → Audio fallback:
   ```bash
   content-processor yt-audio --url "<url>" --output-dir ./sources/yt-audio-downloads
   whisper-cli --model ~/.whisper/models/ggml-base.bin --language en \
     --output-txt --output-file ./sources/yt-audio-transcripts/<id> <audio.wav>
   ```
4. Create summary → `summaries/<slug>_summary.md`

### YouTube Playlist

Download each video individually. For 10+ videos: use sub-agent orchestration.

### Website

1. WebFetch → save to `sources/web-content/<name>_<date>.md`
2. Create summary → `summaries/<slug>_summary.md`

---

## Book Workflow

### Step 1: Convert & Clean

```bash
# EPUB → Markdown
content-processor epub-to-md --input "<file>.epub" --output "<file>.md"

# Clean Calibre artifacts (required for EPUB)
content-processor md-cleanup --input "<file>.md" --output "<file>-cleaned.md"
```

### Step 2: Split Chapters

```bash
# EPUB (usually level 1)
content-processor md-split --input "<file>-cleaned.md" --level 1 --pattern "CHAPTER"

# Native markdown (usually level 2)
content-processor md-split --input "<file>.md" --level 2
```

Output: `<book>-split_chapters/` with `001_*.md`, `002_*.md`, `README.md`

### Step 3: Generate Summaries

Use sub-agent orchestration (3 chapters max per batch).

See `references/orchestration.md` for rules and templates.

---

## Troubleshooting

| Issue                  | Fix                               |
| ---------------------- | --------------------------------- |
| HTTP 429 on YouTube    | Use audio fallback                |
| Whisper stuck          | Check `tail /tmp/whisper_log.txt` |
| Prompt too long        | Ensure agents return "Done" only  |
| `{.calibre}` in output | Run `md-cleanup` before split     |
| Wrong language         | Specify `--language` explicitly   |
| Math broken            | See `references/formats.md`       |
