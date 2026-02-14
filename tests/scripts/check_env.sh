#!/bin/bash
# check_env.sh — Prints environment variables to verify injection
echo "=== Environment Check ==="
echo "PID: $$"
echo "APP_NAME=${APP_NAME:-NOT_SET}"
echo "APP_ENV=${APP_ENV:-NOT_SET}"
echo "APP_PORT=${APP_PORT:-NOT_SET}"
echo "APP_DEBUG=${APP_DEBUG:-NOT_SET}"
echo "CUSTOM_VAR=${CUSTOM_VAR:-NOT_SET}"
echo "=== All environment ==="
env | sort
echo "=== Done ==="
exit 0
