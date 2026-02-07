#!/usr/bin/env bash

CURRENT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
SCRIPTS_DIR="$CURRENT_DIR/scripts"

# Compile Go binary if not present
if [ ! -f "$CURRENT_DIR/bin/ctx" ]; then
    mkdir -p "$CURRENT_DIR/bin"
    echo "Compiling context-sentinel binary..."
    cd "$CURRENT_DIR/src" && go build -o ../bin/ctx main.go
fi

# Register hooks
# pane-focus-in: Trigger check when switching panes
tmux set-hook -g pane-focus-in "run-shell '$SCRIPTS_DIR/context_check.sh'"

# window-linked: Trigger when window created/linked (optional, but good)
# tmux set-hook -g window-linked "run-shell '$SCRIPTS_DIR/context_check.sh'"

# client-session-changed: When switching sessions
tmux set-hook -g client-session-changed "run-shell '$SCRIPTS_DIR/context_check.sh'"

# Bind key to manually refresh or toggle?
# Let's bind 'C' to check status explicitly
tmux bind-key C run-shell "$SCRIPTS_DIR/context_check.sh"
