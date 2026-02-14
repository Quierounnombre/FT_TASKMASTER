#!/bin/bash
# graceful_shutdown.sh — Handles signals for testing stop_signal and kill_wait
# Traps SIGTERM for graceful shutdown, SIGUSR1 as alternative stop signal

CLEANUP_TIME=${CLEANUP_TIME:-3}

cleanup() {
    echo "[graceful] Received SIGTERM — starting cleanup (${CLEANUP_TIME}s)..."
    sleep "$CLEANUP_TIME"
    echo "[graceful] Cleanup complete. Exiting gracefully."
    exit 0
}

handle_usr1() {
    echo "[graceful] Received SIGUSR1 — custom stop signal!"
    echo "[graceful] Performing quick shutdown..."
    sleep 1
    echo "[graceful] Done. Exiting."
    exit 0
}

handle_hup() {
    echo "[graceful] Received SIGHUP — ignoring (reload not supported)"
}

trap cleanup TERM
trap handle_usr1 USR1
trap handle_hup HUP

echo "[graceful] Started PID: $$"
echo "[graceful] Waiting for signals (SIGTERM, SIGUSR1, SIGHUP)..."

while true; do
    echo "[graceful] Still alive at $(date +%T)"
    sleep 2
done
