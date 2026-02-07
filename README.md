# tmux-context-sentinel

A lightweight, offline-first Tmux plugin to guard your context.

## Features
- **Pure Tmux**: No shell config edits required.
- **Resource Light**: Core logic in a fast Go binary.
- **Context Guard**: Warns you if you enter a "BUSY" directory context.

## Installation

Add to your TPM list in `.tmux.conf`:

```tmux
set -g @plugin 'dannyh79/tmux-context-sentinel'
```

Then press `prefix` + `I` to install.

## Usage

1. **Initialize a context** in a project root:
   ```bash
   ~/path/to/plugin/bin/ctx init
   ```
   *(You might want to symlink `ctx` to your path for easier access)*

2. **Start a task**:
   ```bash
   ctx start "Refactoring API"
   ```

3. **Tmux Behavior**:
   - When you focus a pane in this directory, Tmux will display a warning: `⚠️ Context Sentinel: BUSY on 'Refactoring API'`.
   - You can add `#{E:@context_sentinel_status}` to your `status-right` or `status-left` in `.tmux.conf` to see persistent status.

4. **Stop a task**:
   ```bash
   ctx stop
   ```
