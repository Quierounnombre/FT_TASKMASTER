#!/bin/bash
# check_workdir.sh — Prints pwd and creates a marker file
# Used to verify work_dir is applied correctly
echo "=== Working Directory Check ==="
echo "PID: $$"
echo "Current directory: $(pwd)"
echo "Listing:"
ls -la
MARKER="workdir_marker_$(date +%s).txt"
echo "Creating marker file: $MARKER"
echo "Created by check_workdir.sh at $(date)" > "$MARKER"
echo "Marker created successfully: $(pwd)/$MARKER"
echo "=== Done ==="
exit 0
