#!/bin/bash
# exit_code.sh — Exits with configurable exit code via $EXIT_CODE env var
# Usage: EXIT_CODE=42 ./exit_code.sh
echo "Process starting with PID: $$"
echo "Will exit with code: ${EXIT_CODE:-0}"
sleep 2
echo "Exiting now..."
exit ${EXIT_CODE:-0}
