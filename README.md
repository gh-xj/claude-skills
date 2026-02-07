# Claude Skills

A collection of productivity skills for Claude Code.

## Installation

Add this marketplace to Claude Code:

```bash
/plugin marketplace add gh-xj/claude-skills
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

### skill-builder

Interactive guide for creating and maintaining high-quality Claude skills.

**Usage:** `/skill-builder`

**Features:**
- Create new skills with proper structure
- Refactor existing skills for better organization
- Progressive disclosure patterns
- Quality checklists and best practices

## Adding More Skills

To add a new skill, create a folder in `.claude-plugin/skills/` with a `SKILL.md` file, then update `marketplace.json`.
