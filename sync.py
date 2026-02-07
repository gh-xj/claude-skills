#!/usr/bin/env python3
"""Sync skills from ~/.claude/skills/ to this repo."""

import json
import shutil
from pathlib import Path

REPO_SKILLS_DIR = Path(__file__).parent / ".claude-plugin" / "skills"
LOCAL_SKILLS_DIR = Path.home() / ".claude" / "skills"
MARKETPLACE_JSON = Path(__file__).parent / ".claude-plugin" / "marketplace.json"


def get_tracked_skills() -> list[str]:
    """Get skill names from marketplace.json."""
    with open(MARKETPLACE_JSON) as f:
        data = json.load(f)
    return [p["name"] for p in data.get("plugins", [])]


def sync_skill(name: str) -> bool:
    """Copy skill from local to repo. Returns True if synced."""
    src = LOCAL_SKILLS_DIR / name
    dst = REPO_SKILLS_DIR / name

    if not src.exists():
        print(f"  SKIP {name} (not in ~/.claude/skills/)")
        return False

    if dst.exists():
        shutil.rmtree(dst)

    shutil.copytree(src, dst, ignore=shutil.ignore_patterns(
        "__pycache__", "*.pyc", ".DS_Store", "bin", ".task"
    ))
    print(f"  SYNC {name}")
    return True


def main():
    skills = get_tracked_skills()
    print(f"Syncing {len(skills)} skills from ~/.claude/skills/\n")

    synced = sum(1 for s in skills if sync_skill(s))
    print(f"\nDone. {synced}/{len(skills)} skills synced.")
    print("\nNext: git add . && git commit && git push")


if __name__ == "__main__":
    main()
