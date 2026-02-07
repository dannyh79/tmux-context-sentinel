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

echo "Running ctx list with mocks..."
OUTPUT=$("$BIN_DIR/ctx" list 2>&1)
echo "Debug Output:"
echo "$OUTPUT"

# Assertions
if echo "$OUTPUT" | grep -q "Gemini (BUSY)"; then
    echo "PASS: Detected Gemini"
else
    echo "FAIL: Did not detect Gemini"
    exit 1
fi

if echo "$OUTPUT" | grep -q "Shell (IDLE)"; then
    echo "PASS: Detected Idle Shell"
else
    echo "FAIL: Did not detect Idle Shell"
    exit 1
fi

# Cleanup
rm -rf "$TEST_DIR"
