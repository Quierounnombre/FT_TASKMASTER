#!/bin/bash
# long_running.sh — Runs indefinitely, printing heartbeats
# Used for testing long-lived processes, stop, kill, etc.
echo "[long_running] Started PID: $$"
COUNTER=0
while true; do
    COUNTER=$((COUNTER + 1))
    echo "[long_running] Heartbeat #${COUNTER} at $(date +%T)"
    sleep ${HEARTBEAT_INTERVAL:-2}
done
