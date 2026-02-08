#!/usr/bin/env bash

CURRENT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
BIN="$CURRENT_DIR/../bin/ctx"

PANE_ID="$1"
PANE_TTY="$2"

if [ -z "$PANE_ID" ] || [ -z "$PANE_TTY" ]; then
    exit 0
fi

# Detect state
STATE=$("$BIN" detect "$PANE_TTY")

# Set tmux option (user option @ctx_status)
tmux set-option -p -t "$PANE_ID" @ctx_status "$STATE"

# Update global summary
SUMMARY=$("$BIN" summary)
tmux set-option -g @ctx_summary "$SUMMARY"
