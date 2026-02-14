#!/bin/bash
# crash_loop.sh — Always crashes after a short delay
# Used to test watcher restart behavior
echo "[crash_loop] Started PID: $$"
echo "[crash_loop] Working for a moment..."
sleep ${CRASH_DELAY:-3}
echo "[crash_loop] BOOM! Crashing with exit code ${CRASH_EXIT:-1}"
exit ${CRASH_EXIT:-1}
