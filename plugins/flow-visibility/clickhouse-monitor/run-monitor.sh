#!/usr/bin/env bash

THIS_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
DEFAULT_SKIP_ROUND_NUMBER = 3
statefile=$THIS_DIR/statefile

declare -i skip_rounds_number

if [ -f "$statefile" ]; then
  read -r skip_rounds_number <"$statefile"
else
  skip_rounds_number=0
fi

if [skip_rounds_number == 0 && $((clickhouse-monitor))]; then
  skip_rounds_number=$DEFAULT_SKIP_ROUND_NUMBER
else
  skip_rounds_number=$(($skip_rounds_number - 1))
fi

printf '%d\n' "$persistent_counter" >"$statefile"
