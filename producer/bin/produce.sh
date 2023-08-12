#!/bin/bash

fatal() {
  echo "$@" 1>&2
  exit 1
}

[ $# -eq 0 ] && fatal "missing topic"
PRODUCER_TOPIC="$1"
shift

PRODUCER_KEY=$(date +%y%m%d%H%S)
PRODUCER_BASE_URL="http://localhost:30000/produce"

PRODUCER_URL="$PRODUCER_BASE_URL/$PRODUCER_TOPIC/$PRODUCER_KEY"

PRODUCER_CMD="curl -X POST -H \"Content-Type: application/json\" \"$PRODUCER_URL\""

while true; do
  text=$(jq -R -s '.')
  echo "would pipe \"$text\" into \"$PRODUCER_CMD\""
done



