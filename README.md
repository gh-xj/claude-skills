# Claude Skills

A collection of productivity skills for Claude Code.

## Installation

Add this marketplace to Claude Code:

```bash
/plugin marketplace add YOUR_USERNAME/claude-skills
```

## Available Skills

### content-processor

Process any content source (YouTube, websites, EPUBs, markdown) into structured summaries.

**Usage:** `/content-processor`

**Features:**
- YouTube video/playlist transcription
- Website content fetching
- EPUB to markdown conversion
- Batch content summarization

## Adding More Skills

To add a new skill, create a folder in `.claude-plugin/skills/` with a `SKILL.md` file, then update `marketplace.json`.
