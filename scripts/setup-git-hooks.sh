#!/bin/bash
# Copy prepare-commit-msg hook to remove Cursor co-author from commits
HOOK_SRC="$(dirname "$0")/prepare-commit-msg"
HOOK_DST="$(git rev-parse --git-dir)/hooks/prepare-commit-msg"
if [ -f "$HOOK_SRC" ]; then
  cp "$HOOK_SRC" "$HOOK_DST"
  chmod +x "$HOOK_DST"
  echo "Git hook installed: prepare-commit-msg"
else
  echo "Hook source not found: $HOOK_SRC"
  exit 1
fi
