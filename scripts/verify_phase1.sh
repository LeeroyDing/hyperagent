#!/bin/bash
set -e

BINARY="./hyperagent"
PID_FILE="$HOME/.hyperagent/hyperagent.pid"

echo "--- Phase 1 Integration Test ---"

# 1. Start daemon
echo "Starting daemon..."
$BINARY up -d
sleep 2

# 2. Check PID file
if [ ! -f "$PID_FILE" ]; then
    echo "Error: PID file not found at $PID_FILE"
    exit 1
fi
PID=$(cat "$PID_FILE")
echo "Daemon running with PID: $PID"

# 3. Check status command
echo "Checking status..."
$BINARY status

# 4. Stop daemon
echo "Stopping daemon..."
$BINARY down
sleep 2

# 5. Verify cleanup
if [ -f "$PID_FILE" ]; then
    echo "Error: PID file still exists after shutdown"
    exit 1
fi

if ps -p $PID > /dev/null; then
    echo "Error: Process $PID still running after shutdown"
    exit 1
fi

echo "--- Phase 1 Integration Test Passed ---"
