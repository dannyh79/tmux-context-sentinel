#!/usr/bin/env bash

# Setup
TEST_DIR=$(mktemp -d)
BIN_DIR="$TEST_DIR/bin"
mkdir -p "$BIN_DIR"

# Copy ctx binary
cp ./bin/ctx "$BIN_DIR/"

# Create mock tmux
cat <<EOF > "$BIN_DIR/tmux"
#!/bin/sh
if [ "\$1" = "list-panes" ]; then
    # Return 2 panes
    # id|||tty|||path|||session|||winidx|||paneidx|||winname
    echo "%0|||mocktty1|||/tmp|||sess1|||1|||0|||win1"
    echo "%1|||mocktty2|||/tmp|||sess1|||1|||1|||win1"
fi
EOF
chmod +x "$BIN_DIR/tmux"

# Create mock ps
cat <<EOF > "$BIN_DIR/ps"
#!/bin/sh
# args: -t tty -o ...
TTY="\$2"
if [ "\$TTY" = "mocktty1" ]; then
    echo "PID PPID COMMAND"
    echo "100 1 -zsh"
    echo "101 100 gemini-cli start"
elif [ "\$TTY" = "mocktty2" ]; then
    echo "PID PPID COMMAND"
    echo "200 1 -zsh"
fi
EOF
chmod +x "$BIN_DIR/ps"

# Create mock git
cat <<EOF > "$BIN_DIR/git"
#!/bin/sh
echo "main"
EOF
chmod +x "$BIN_DIR/git"

# Run test
export PATH="$BIN_DIR:$PATH"

echo "Running ctx list (default: busy only)..."
OUTPUT=$("$BIN_DIR/ctx" list 2>&1)
echo "Debug Output (Default):"
echo "$OUTPUT"

# Assertions for Default
if echo "$OUTPUT" | grep -q "Gemini (BUSY)"; then
    echo "PASS: Detected Gemini (Default)"
else
    echo "FAIL: Did not detect Gemini (Default)"
    exit 1
fi

if echo "$OUTPUT" | grep -q "Shell (IDLE)"; then
    echo "FAIL: Detected Idle Shell (Default) - Should be hidden"
    exit 1
else
    echo "PASS: Hidden Idle Shell (Default)"
fi

echo "Running ctx list --show-idle..."
OUTPUT_ALL=$("$BIN_DIR/ctx" list --show-idle 2>&1)
echo "Debug Output (All):"
echo "$OUTPUT_ALL"

# Assertions for Show Idle
if echo "$OUTPUT_ALL" | grep -q "Gemini (BUSY)"; then
    echo "PASS: Detected Gemini (Show Idle)"
else
    echo "FAIL: Did not detect Gemini (Show Idle)"
    exit 1
fi

if echo "$OUTPUT_ALL" | grep -q "Shell (IDLE)"; then
    echo "PASS: Detected Idle Shell (Show Idle)"
else
    echo "FAIL: Did not detect Idle Shell (Show Idle)"
    exit 1
fi

# Cleanup
rm -rf "$TEST_DIR"
