#!/bin/sh

#the script will exit immediately if a command returns a non-zero status
set -e

#start the app
echo "Start the app"

#take all parameters passed to the script and run it
exec "$@"