#!/bin/bash

# Test: test_popup_exit_order.sh
# Description: Verifies that tmux switch-client happens AFTER fzf exits.

set -e

# Setup temporary directory for mocks
TEST_DIR=$(mktemp -d)
MOCK_BIN_DIR="$TEST_DIR/bin"
MOCK_LOG="$TEST_DIR/events.log"
mkdir -p "$MOCK_BIN_DIR"

cleanup() {
    rm -rf "$TEST_DIR"
}
trap cleanup EXIT

# Export Mock Log location for mocks to use
export MOCK_LOG

# 1. Mock 'tmux'
cat > "$MOCK_BIN_DIR/tmux" << 'EOF'
#!/bin/bash
if [[ "$1" == "show-option" ]]; then
    # Return 0 for @ctx_show_idle
    echo "0"
elif [[ "$1" == "switch-client" ]]; then
    echo "TMUX_SWITCH $@" >> "$MOCK_LOG"
fi
EOF
chmod +x "$MOCK_BIN_DIR/tmux"

# 2. Mock 'fzf' (and fzf-tmux to be safe)
cat > "$MOCK_BIN_DIR/fzf" << 'EOF'
#!/bin/bash
# Consume input
cat > /dev/null
# Simulate user delay slightly
sleep 0.1
# Output the simulated selection (matches ctx mock output)
echo "My Session: Window 1 ||| %100"
# Log exit event
echo "FZF_EXITED" >> "$MOCK_LOG"
EOF
chmod +x "$MOCK_BIN_DIR/fzf"
ln -s "$MOCK_BIN_DIR/fzf" "$MOCK_BIN_DIR/fzf-tmux"

# 3. Mock the 'ctx' binary (referenced as ../bin/ctx relative to script)
# The script calculates BIN based on location.
# We need to make sure the script uses our mock ctx or we just rely on PATH if we can.
# But scripts/popup.sh uses: BIN="$CURRENT_DIR/../bin/ctx"
# We cannot easily override BIN without modifying the script or placing a mock file there.
# Strategy: We will create a temporary environment where ../bin/ctx exists.

# Create project structure in temp dir
mkdir -p "$TEST_DIR/scripts"
mkdir -p "$TEST_DIR/bin"

# Copy the script under test
cp "$(dirname "$0")/../scripts/popup.sh" "$TEST_DIR/scripts/"

# Create mock ctx
cat > "$TEST_DIR/bin/ctx" << 'EOF'
#!/bin/bash
echo "My Session: Window 1 ||| %100"
EOF
chmod +x "$TEST_DIR/bin/ctx"

# 4. Run the script with mocks in PATH
export PATH="$MOCK_BIN_DIR:$PATH"

# Execute
"$TEST_DIR/scripts/popup.sh"

# 5. Assertions
if [ ! -f "$MOCK_LOG" ]; then
    echo "FAIL: No events logged."
    exit 1
fi

CONTENTS=$(cat "$MOCK_LOG")
echo "Log contents:"
echo "$CONTENTS"

# Check order
# Expected:
# FZF_EXITED
# TMUX_SWITCH -t %100

LINE1=$(sed -n '1p' "$MOCK_LOG")
LINE2=$(sed -n '2p' "$MOCK_LOG")

if [[ "$LINE1" != "FZF_EXITED" ]]; then
    echo "FAIL: First event was not FZF_EXITED. Got: '$LINE1'"
    exit 1
fi

if [[ "$LINE2" != "TMUX_SWITCH switch-client -t %100" ]]; then
    echo "FAIL: Second event was not 'TMUX_SWITCH switch-client -t %100'. Got: '$LINE2'"
    exit 1
fi

echo "PASS: fzf exited before tmux switch occurred."
