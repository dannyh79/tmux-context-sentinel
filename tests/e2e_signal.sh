#!/bin/bash
set -e

# Setup
TEST_DIR=$(mktemp -d)
BIN_DIR="$TEST_DIR/bin"
mkdir -p "$BIN_DIR"

# Compile main binary
cd src
go build -o "$BIN_DIR/ctx" .
cd ..

# Create Mock Tmux and TTY
cat <<EOF > "$BIN_DIR/tmux"
#!/bin/bash
if [ "\$1" == "list-panes" ]; then
    # Return a fake pane matching current TTY
    # We use a fixed TTY for testing to match the mocked 'tty' command
    FAKE_TTY="/dev/pts/fake"
    echo "%0|||\$FAKE_TTY|||$TEST_DIR|||sess|||1|||1|||win"
fi
EOF
chmod +x "$BIN_DIR/tmux"

cat <<EOF > "$BIN_DIR/tty"
#!/bin/bash
echo "/dev/pts/fake"
EOF
chmod +x "$BIN_DIR/tty"

# Export Path
export PATH="$BIN_DIR:$PATH"
export HOME="$TEST_DIR" # Isolate state.json

# Test Flow
echo "1. Run ctx signal --status running"
"$BIN_DIR/ctx" signal --status running --agent TestBot

echo "2. Run ctx summary"
SUMMARY=$("$BIN_DIR/ctx" summary)
echo "Summary output: '$SUMMARY'"

if [[ "$SUMMARY" == *"[busy AI: 1/1]"* ]]; then
    echo "SUCCESS: Summary updated correctly"
else
    echo "FAILURE: Expected '[busy AI: 1/1]', got '$SUMMARY'"
    exit 1
fi

echo "3. Run ctx signal --status waiting"
"$BIN_DIR/ctx" signal --status waiting --agent TestBot

SUMMARY=$("$BIN_DIR/ctx" summary)
echo "Summary output: '$SUMMARY'"
# Should be 0/1 (1 total, 0 busy)
# Wait, total counts agents that are in state.
if [[ "$SUMMARY" == *"[busy AI: 0/1]"* ]]; then
    echo "SUCCESS: Summary updated correctly to idle"
else
    echo "FAILURE: Expected '[busy AI: 0/1]', got '$SUMMARY'"
    exit 1
fi

rm -rf "$TEST_DIR"
