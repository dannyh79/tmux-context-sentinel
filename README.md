# Tmux Context Sentinel

**Context Sentinel** is a precision tool for Tmux that detects what you are working on. It identifies active AI agents and CLI tools in your panes and provides a "Jump Menu" to switch contexts instantly.

## Features

- **Surgical AI Detection**: Uses process tree analysis to identify active tools (no file-based status tracking).
- **Supports**: `gemini-cli`, `kiro-cli`, `cursor-cli`, `aider`, `opencode`.
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

### Global Status Counter

You can also display a global summary of all active AI agents across your session (e.g., `[busy AI: 1/3]`). This is stored in the global variable `@ctx_summary`.

**Example (.tmux.conf):**
```tmux
set -g status-left "#{E:@ctx_summary} "
```
This updates automatically as you switch panes.

## Requirements

- Tmux 3.2+ (for popup support)
- Go (to build the binary, happens automatically on first run)
- `fzf` (required for the menu)

## User Guide: Agent Integrations

For the most accurate status tracking, configure your AI agents to signal their state directly to Context Sentinel. This enables the `[busy AI: X/Y]` counter to work flawlessly without relying on polling.

Ensure `ctx` is in your `PATH` or use the absolute path (e.g., `~/.tmux/plugins/tmux-context-sentinel/bin/ctx`).

### 1. Gemini CLI

Add the following to your `~/.gemini/settings.json` (or project-level `.gemini/settings.json`):

```json
{
  "hooks": {
    "BeforeAgent": [
       { 
         "matcher": "*", 
         "hooks": [{ "command": "ctx signal --status running --agent Gemini", "type": "command" }] 
       }
    ],
    "AfterAgent": [
       { 
         "matcher": "*", 
         "hooks": [{ "command": "ctx signal --status waiting --agent Gemini", "type": "command" }] 
       }
    ]
  }
}
```

### 2. Claude Code

Add the following to your `~/.claude/settings.json`:

```json
{
  "hooks": {
    "UserPromptSubmit": [
      {
        "hooks": [{ "command": "ctx signal --status running --agent Claude", "type": "command" }]
      }
    ],
    "Stop": [
      {
        "hooks": [{ "command": "ctx signal --status waiting --agent Claude", "type": "command" }]
      }
    ]
  }
}
```

### 3. Kiro CLI

Refer to the Kiro CLI documentation for hook configuration. You can use the generic commands:
- **Start**: `ctx signal --status running --agent Kiro`
- **End**: `ctx signal --status waiting --agent Kiro`

### 4. OpenCode

To use Context Sentinel with OpenCode, configure the `ctx` binary as a hook or plugin in your settings.

- **Start Command**: `ctx signal --status running --agent OpenCode`
- **End Command**: `ctx signal --status waiting --agent OpenCode`

Ensure `ctx` is in your `PATH`, or use the absolute path (e.g., `~/.tmux/plugins/tmux-context-sentinel/bin/ctx`).

## Development

- **Build**: `cd src && go build -o ../bin/ctx .`
- **Test**: `go test ./src/...` and `./tests/test_integration.sh`

## License

MIT
