#!/bin/bash

#BAD PRACTICE, if something fails stop.

#RELATIVE PATH
CLI="../CLI/CLI"

#PATH launch only from tests
TEST_PATH="$(pwd)"

#LINE
LINE="--------------------------------------------------------------"

if [ ! -x "$CLI" ]; then
  echo "CLI executable not found or not executable at $CLI"
  exit 1
fi

echo "Running all tests..."

mapfile -t yamls < <(find "$TEST_PATH" -name "*.yaml" -type f | sort)

for yaml in "${yamls[@]}"; do
  echo "$LINE"
  echo "Running: $yaml"
  (
    "$CLI" load "$yaml"
  )
  if [ $? -eq 0 ]; then
    echo "Finished: $yaml ✅"
  else
    echo "FAILED: $yaml ❌"
  fi
  echo "$LINE"
done

echo "All tests executed."
