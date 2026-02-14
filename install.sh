#!/usr/bin/env bash

set -e

# Detect Tmux Plugin Manager path
detect_plugin_path() {
    if [ -n "$TMUX_PLUGIN_MANAGER_PATH" ]; then
        echo "$TMUX_PLUGIN_MANAGER_PATH"
        return
    fi

    # Check XDG location
    if [ -n "$XDG_CONFIG_HOME" ] && [ -d "$XDG_CONFIG_HOME/tmux/plugins" ]; then
        echo "$XDG_CONFIG_HOME/tmux/plugins"
        return
    fi

    # Check standard ~/.config location
    if [ -d "$HOME/.config/tmux/plugins" ]; then
        echo "$HOME/.config/tmux/plugins"
        return
    fi

    # Default to ~/.tmux/plugins
    echo "$HOME/.tmux/plugins"
}

PLUGIN_BASE_DIR=$(detect_plugin_path)
DEFAULT_INSTALL_DIR="${PLUGIN_BASE_DIR}/tmux-context-sentinel"
TARGET_DIR="${1:-$DEFAULT_INSTALL_DIR}"
CURRENT_DIR="$(pwd)"

echo "Installing to: $TARGET_DIR"

# 1. Build the Go binary
echo "Building Go binary..."
mkdir -p bin
if ! (cd src && go build -o ../bin/ctx .); then
  echo "Error: Failed to build Go binary."
  exit 1
fi

# 2. Install files if necessary
if [ "$CURRENT_DIR" != "$TARGET_DIR" ]; then
  echo "Installing files to $TARGET_DIR..."
  mkdir -p "$TARGET_DIR"

  # Use rsync if available for cleaner copy, otherwise cp
  if command -v rsync >/dev/null 2>&1; then
    rsync -av --exclude='.git' --exclude='tests' --exclude='.github' ./ "$TARGET_DIR/"
  else
    cp -R ./* "$TARGET_DIR/"
  fi
else
  echo "Running in-place installation..."
fi

# 3. Reload Tmux
echo "Reloading Tmux plugin..."
if command -v tmux >/dev/null 2>&1; then
  tmux run-shell "$TARGET_DIR/context-sentinel.tmux"
  echo "Success! Plugin installed/updated and reloaded."
else
  echo "Warning: tmux not found. Plugin installed but not reloaded."
fi
