#!/usr/bin/env bash

CURRENT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
SCRIPTS_DIR="$CURRENT_DIR/scripts"

# Compile Go binary if not present (Auto-bootstrap)
if [ ! -f "$CURRENT_DIR/bin/ctx" ]; then
    mkdir -p "$CURRENT_DIR/bin"
    echo "Compiling context-sentinel binary..."
    cd "$CURRENT_DIR/src" && go build -o ../bin/ctx .
fi

# 1. Register Hooks (Requirement 1)
# Pass pane_id and pane_tty to the hook script
tmux set-hook -g pane-focus-in "run-shell '$SCRIPTS_DIR/hook.sh \"#{pane_id}\" \"#{pane_tty}\"'"
tmux set-hook -g pane-mode-changed "run-shell '$SCRIPTS_DIR/hook.sh \"#{pane_id}\" \"#{pane_tty}\"'"

# 2. Key Bindings (Requirement 4)
# Double-Prefix to open Jump Menu
get_tmux_option() {
    local option="$1"
    local default_value="$2"
    local option_value=$(tmux show-option -gqv "$option")
    if [ -z "$option_value" ]; then
        echo "$default_value"
    else
        echo "$option_value"
    fi
}

# Get current prefix
PREFIX=$(get_tmux_option "prefix" "C-b")

# The prefix key itself (e.g., C-b)
# We want to bind 'Prefix + Prefix' -> Popup
# So we bind the key '$PREFIX' inside the prefix table.
tmux bind-key "$PREFIX" run-shell "$SCRIPTS_DIR/popup.sh"

echo "Tmux Context Sentinel loaded."
