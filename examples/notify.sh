#!/usr/bin/env bash
# Example notification script
BACKEND="$1"
STAGE="$2"
LEVEL="$3"
MESSAGE="$4"

echo "Backend: $BACKEND"
echo "Stage: $STAGE"

if [[ "$LEVEL" == "success" ]]; then
  echo "All good \o/"
else
  echo "Something went wrong: $MESSAGE"
fi

