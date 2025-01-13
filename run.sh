#!/bin/sh
set -e 
(
  cd "$(dirname "$0")" 
  go build -o /dist/app ./app
)

# Copied from .codecrafters/run.sh
#
# - Edit this to change how your program runs locally
# - Edit .codecrafters/run.sh to change how your program runs remotely
exec /tmp/codecrafters-build-http-server-go "$@"
