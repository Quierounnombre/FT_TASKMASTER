#!/bin/bash
# test_random_fail.sh
while true; do
    echo "Trabajando..."
    sleep 3
    if [ $((RANDOM % 5)) -eq 0 ]; then
        echo "Error aleatorio!"
        exit 1
    fi
done