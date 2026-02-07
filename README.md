# Tmux Context Sentinel

**Context Sentinel** is a precision tool for Tmux that detects what you are working on. It identifies active AI agents and CLI tools in your panes and provides a "Jump Menu" to switch contexts instantly.

## Features

- **Surgical AI Detection**: Uses process tree analysis to identify active tools (no file-based status tracking).
- **Supports**: `gemini-cli`, `kiro-cli`, `cursor-cli`, `aider`.
- **Zero Config**: Works out of the box with standard process names.
- **Jump Menu**: Press `Prefix` + `Prefix` to open a popup menu.
  - Lists all panes with Git Branch and Active Tool.
  - Shows `BUSY` / `IDLE` status.
  - Fuzzy search to jump to any context.

## Installation

### Manual
1. Clone the repository:
   ```bash
   git clone https://github.com/danny/tmux-context-sentinel ~/.tmux/plugins/tmux-context-sentinel
   ```
2. Add this line to your `~/.tmux.conf`:
   ```tmux
   run-shell ~/.tmux/plugins/tmux-context-sentinel/context-sentinel.tmux
   ```
3. Reload tmux (`tmux source ~/.tmux.conf`).

### TPM (Tmux Plugin Manager)
Add this to your `~/.tmux.conf`:
```tmux
set -g @plugin 'danny/tmux-context-sentinel'
```
Press `Prefix` + `I` to install.

### Binary Location
The core logic runs via a Go binary compiled automatically during installation.
- **Path**: `~/.tmux/plugins/tmux-context-sentinel/bin/ctx` (if installed via TPM)
- **Manual Path**: `~/projects/tmux-context-sentinel/bin/ctx` (if cloned manually)

You can add the `bin/` directory to your `PATH` to use the `ctx` command directly in your shell.

## Configuration

### Jump Menu Behavior
By default, the Jump Menu hides "idle" panes (panes with no active recognized tool) to reduce clutter.

- **Show Idle Panes**: To show all panes including idle ones, set this variable to `1`.
  ```tmux
  set -g @ctx_show_idle '1'  # Default: 0
  ```

## Usage

- **Jump Menu**: Press your prefix key TWICE (e.g., `C-b C-b`).
  - A popup will appear listing all panes.
  - Select a pane to switch to it.
  - The list shows: `[Branch] - [Tool Name] (BUSY/IDLE)`.

- **Automatic Detection**:
  - The plugin automatically scans panes when you focus them.
  - Status is updated live when you open the menu.

## Status Line Integration

You can display the current context directly in your tmux status bar. The plugin exposes the detected tool and branch via the pane option `@ctx_status`.

To use it, refer to the option using `#{E:@ctx_status}` in your status line configuration.

**Example (.tmux.conf):**
```tmux
set -g status-right "Context: #{E:@ctx_status}"
```

This will dynamically display the detected context (e.g., ` main ⚡ gemini-cli`) for the active pane.

## Requirements

- Tmux 3.2+ (for popup support)
- Go (to build the binary, happens automatically on first run)
- `fzf` (required for the menu)

## Development

- **Build**: `cd src && go build -o ../bin/ctx .`
- **Test**: `go test ./src/...` and `./tests/test_integration.sh`

## License

MIT
