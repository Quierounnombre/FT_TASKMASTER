#!/bin/bash
# slow_start.sh — Takes a configurable time to "initialize"
# Used to test launch_wait behavior
STARTUP_DELAY=${STARTUP_DELAY:-5}
echo "[slow_start] PID: $$ — Initializing (will take ${STARTUP_DELAY}s)..."
sleep "$STARTUP_DELAY"
echo "[slow_start] Initialization complete!"
echo "[slow_start] Now running normally..."
COUNTER=0
while true; do
    COUNTER=$((COUNTER + 1))
    echo "[slow_start] Running #${COUNTER}"
    sleep 2
done
