# Lessons Learned

## Whisper Transcription (2025-10-19)

### What Worked

- Full model path `~/.whisper/models/ggml-base.bin` required
- Background processing with `run_in_background: true`
- Apple M1 Metal: ~3 min for 47-min video
- Two-step: audio download (fast) + transcription (separate)
- Specifying `--language` upfront prevents mismatch

### Issues & Fixes

| Issue               | Fix                               |
| ------------------- | --------------------------------- |
| Directory confusion | Use `sources/yt-audio-downloads/` |
| Progress monitoring | Check `tail /tmp/whisper_log.txt` |
| Process completion  | Verify output file exists         |

## Batch Processing (2025-11-09)

**Scenario:** 24-lecture MIT playlist

**Solution:** Multiple Task tool agents in parallel

- 5 agents × ~5 lectures each
- Each agent works autonomously

**When to use:**

- 10+ items: Multiple parallel agents (5 max)
- 6-9 items: 2-3 agents
- 1-5 items: Sequential is fine

## Context Management (2025-12)

**Problem:** "Prompt too long" errors when processing many chapters.

**Root cause:** Sub-agents returning full summaries + TaskOutput retrieving content.

**Solution:**

1. Sub-agents write to files, return only "Done: \<item\>"
2. Verify via `ls` instead of TaskOutput
3. Max 3 chapters per batch for books
4. For >15 chapters, split across conversations

## Language Detection (2025-11-05)

**Process:**

1. Download with default language
2. Read first 100 lines to verify
3. If wrong, retry with explicit `--language`

**Better:** Ask user upfront, specify explicitly.

## Twitter Format (2025-11-24)

**Triggers:** "twitter format", "tweet style", "social media", "key points"

Each point under 280 chars, include hashtags and emojis.

## EPUB Artifacts (2025-12)

**Problem:** Calibre conversion leaves `{.calibre}` CSS tags.

**Solution:** Always run `markdown_cleanup` after `conv_to_md`, before splitting.

**Before:** `[[CHAPTER 1]{.calibre3}]{#part0004...}`
**After:** `CHAPTER 1`

## Summary Depth for Books (2026-01)

**Problem:** Initial book summaries were too concise, losing structure and logical flow. User: "I don't even know the structure or logic of the original book."

**Root cause:** Sub-agent prompt said "concise summary" — wrong for books.

**Solution:** For books, prompt should say "detailed summary" and emphasize preserving structure, logical flow, and key quotes. Don't just extract bullet points.

**Key insight:** Match depth to content type. Books need deep summaries that preserve how arguments build on each other.
