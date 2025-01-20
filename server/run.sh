#!/bin/bash

# Get GOPATH or use default
GOPATH=${GOPATH:-$HOME/go}
AIR_PATH="$GOPATH/bin/air"

# Check if air is installed
if [ ! -f "$AIR_PATH" ]; then
    echo "Installing air..."
    go install github.com/air-verse/air@latest
fi

# Run the server with air using full path
echo "Starting server with hot reloading..."
"$AIR_PATH" 