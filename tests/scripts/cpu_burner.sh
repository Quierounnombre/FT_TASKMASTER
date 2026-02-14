#!/bin/bash
# cpu_burner.sh — Burns CPU briefly then exits
# Used for stress testing with multiple instances
BURN_TIME=${BURN_TIME:-3}
echo "[cpu_burner] PID: $$ — Burning CPU for ${BURN_TIME}s..."
END=$((SECONDS + BURN_TIME))
while [ $SECONDS -lt $END ]; do
    : # busy loop
done
echo "[cpu_burner] PID: $$ — Done burning. Total time: ${SECONDS}s"
exit 0
