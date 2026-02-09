#!/usr/bin/env bash

CURRENT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
BIN="$CURRENT_DIR/../bin/ctx"

# Use fzf-tmux if available, otherwise fzf
FZF_CMD="fzf"
if command -v fzf-tmux >/dev/null 2>&1; then
    FZF_CMD="fzf-tmux -p 80%,80%"
fi

# Determine flags based on tmux option
SHOW_IDLE=$(tmux show-option -gqv "@ctx_show_idle")
FLAGS=""
if [ "$SHOW_IDLE" = "1" ]; then
    FLAGS="--show-idle"
fi

# Run list, pipe to fzf, capture output
# Format: Display Text ||| TargetID
SELECTED=$("$BIN" list $FLAGS | \
    $FZF_CMD --delimiter=' \|\|\| ' --with-nth=1 --reverse --header="Select Pane to Switch")

# Check if selection was made
if [ -z "$SELECTED" ]; then
    exit 0
fi

# Extract TargetID (column 2)
TARGET_ID=$(echo "$SELECTED" | awk -F ' \\|\\|\\| ' '{print $2}')

# Switch client if ID is valid
if [ -n "$TARGET_ID" ]; then
    tmux switch-client -t "$TARGET_ID"
fi
