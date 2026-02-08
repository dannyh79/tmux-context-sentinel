#!/usr/bin/env bash

CURRENT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
BIN="$CURRENT_DIR/../bin/ctx"

# Use fzf-tmux if available, otherwise fzf
FZF_CMD="fzf"
if command -v fzf-tmux >/dev/null 2>&1; then
    FZF_CMD="fzf-tmux -p 80%,80%"
fi

# Run list, pipe to fzf, switch client
# Format: Display Text ||| TargetID
"$BIN" list | \
$FZF_CMD --delimiter=' \|\|\| ' --with-nth=1 --reverse --header="Select Pane to Switch" \
    --bind 'enter:execute(tmux switch-client -t {2})+accept'
