#!/bin/bash
# test_signals.sh
trap 'echo "Recibí SIGTERM, limpiando..."; exit 0' TERM
trap 'echo "Recibí SIGHUP, recargando..."' HUP

while true; do
    echo "Esperando señales... PID: $$"
    sleep 5
done