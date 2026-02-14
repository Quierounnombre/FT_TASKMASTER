#!/bin/bash
# check_umask.sh — Creates a file and checks its permissions
# Used to verify umask is applied correctly
echo "=== Umask Check ==="
echo "PID: $$"
echo "Current umask: $(umask)"
TESTFILE="/tmp/umask_test_$(date +%s).txt"
echo "test content" > "$TESTFILE"
echo "Created file: $TESTFILE"
echo "File permissions:"
ls -la "$TESTFILE"
PERMS=$(stat -c '%a' "$TESTFILE" 2>/dev/null || stat -f '%Lp' "$TESTFILE" 2>/dev/null)
echo "Numeric permissions: $PERMS"
rm -f "$TESTFILE"
echo "=== Done ==="
exit 0
