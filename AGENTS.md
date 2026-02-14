# Agent Guide for Tmux Context Sentinel

This repository implements a Tmux plugin that detects and reports the status of active AI agents running in tmux panes.

## 🛠 Build & Test Commands

### Go (Core Logic)
The core logic resides in `src/`.

- **Install from Source**:
  Builds the binary and installs it to the local tmux plugin directory.
  ```bash
  ./install.sh
  ```
- **Build Binary Only**:
  ```bash
  cd src && go build -o ../bin/ctx .
  ```
- **Run All Unit Tests**:
  ```bash
  cd src && go test ./...
  ```
- **Run a Single Unit Test**:
  ```bash
  # Replace 'TestDetectStatusFromProcs' with target test name
  cd src && go test -v -run TestDetectStatusFromProcs
  ```
- **Linting**:
  Standard `go vet` and `gofmt` should be applied.
  ```bash
  cd src && go vet ./...
  ```

### Integration Tests
Integration tests are Bash scripts in `tests/` that require a running tmux environment.

- **Run Integration Tests**:
  ```bash
  # Ensure you are in the root directory
  ./tests/test_integration.sh
  ```

## 📐 Code Style & Conventions

### Go Code (`src/`)
Follow standard Go idioms (Effective Go).

- **Formatting**: Always use `gofmt`.
- **Naming**:
  - **Public**: `PascalCase` (e.g., `DetectState`, `AgentStatus`).
  - **Private**: `camelCase` (e.g., `detectStatusFromProcs`).
  - **Files**: `snake_case.go` (e.g., `detector_test.go`).
- **Error Handling**:
  - Explicit check: `if err != nil { return err }`.
  - Do not panic unless during bootstrap.
  - Wrap errors with context when helpful: `fmt.Errorf("context: %w", err)`.
- **Imports**: Group standard library imports first, separated from 3rd-party imports.

### Bash Scripts (`scripts/`, `tests/`)
- **Shebang**: Use `#!/usr/bin/env bash`.
- **Variables**:
  - Global/Env: `SCREAMING_SNAKE_CASE` (e.g., `CURRENT_DIR`).
  - Locals: `snake_case` (e.g., `local option_value`).
- **Safety**:
  - Quote all variable expansions: `"$VAR"`.
  - Use `[[ ]]` for conditions over `[ ]`.
  - Use `local` for function-scope variables.

## 📂 Project Structure

- `src/`: Go source code and unit tests.
- `scripts/`: Helper Bash scripts invoked by tmux hooks.
- `tests/`: End-to-end integration tests.
- `bin/`: Output directory for the compiled `ctx` binary.
- `context-sentinel.tmux`: Plugin entry point (handles bootstrap).

## 🤖 Common Tasks

### Adding a New Agent Detector
1. Modify `src/detector.go`.
2. Update `detectStatusFromProcs` to recognize the new process signature.
3. Add a test case in `src/detector_test.go`.

### Updating State Logic
1. `src/storage.go` handles file-based state persistence.
2. `src/main.go` wires up the CLI commands (`signal`, `list`, etc.).

## 🔄 Development Workflow for Agents

When implementing changes, follow this cycle:

1.  **Analyze**: Understand the requirement and the existing code in `src/`.
2.  **Test First**: If adding logic to `detector.go`, create a failing test case in `detector_test.go` first.
3.  **Implement**: Modify the Go code to satisfy the requirement.
4.  **Verify**: Run `go test ./...` to ensure all tests pass.
5.  **Build**: Run `go build -o ../bin/ctx .` inside `src/` to confirm compilation.

## 🐛 Debugging Tips

- **Tmux Messages**: The plugin uses `tmux display-message` for notifications.
- **Manual Execution**: You can run the binary directly to test logic without tmux hooks:
  ```bash
  ./bin/ctx detect /dev/ttys001
  ./bin/ctx list
  ```
- **Logs**: If debugging the hooks, you may need to add temporary logging to a file like `/tmp/sentinel.log` inside `scripts/hook.sh`.
