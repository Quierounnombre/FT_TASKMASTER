#!/bin/bash
# write_outputs.sh — Writes to both stdout and stderr
# Used to test stdout/stderr redirection to files
echo "[STDOUT] Process started PID: $$"
echo "[STDERR] Process started PID: $$" >&2

for i in 1 2 3 4 5; do
    echo "[STDOUT] Line $i at $(date +%T)"
    echo "[STDERR] Error line $i at $(date +%T)" >&2
    sleep 1
done

echo "[STDOUT] Process finished successfully"
echo "[STDERR] No real errors, just testing" >&2
exit 0
