#!/usr/bin/env bash

CURRENT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PLUGIN_ROOT="$CURRENT_DIR/.."
CTX_BIN="$PLUGIN_ROOT/bin/ctx"

# Ensure binary exists
if [ ! -f "$CTX_BIN" ]; then
    # Silent fail or try to build? For now silent to avoid spamming tmux
    exit 0
fi

# Run status check in current pane's path
# Tmux hook runs in the context of the pane, so pwd *should* be correct if we use pane_current_path
# Actually, the hook environment might not change directory automatically.
# We need to get pane_current_path from tmux.

PANE_PATH=$(tmux display-message -p "#{pane_current_path}")
cd "$PANE_PATH" || exit 1

OUTPUT=$("$CTX_BIN" status 2>/dev/null)

if [[ "$OUTPUT" == "BUSY"* ]]; then
    TASK=${OUTPUT#BUSY: }
    # Update a tmux option that can be displayed in status bar
    tmux set-option -w @context_sentinel_status " #[fg=white,bg=red] BUSY: $TASK #[default]"
    
    # Optional: Display message if it's a focus event
    # Only display message if specifically triggered, or maybe just rely on status bar for less intrusive UI
    # Requirement: "display a visible warning in the tmux message area or status bar"
    # Let's do message area for focus-in to be noticeable
    tmux display-message "⚠️  Context Sentinel: BUSY on '$TASK'"
else
    tmux set-option -w @context_sentinel_status ""
fi
