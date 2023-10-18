#!/bin/bash

# Specify the path to the directory where the "gkilo" binary is located
GKILO_PATH="./bin"
DESTINATION_PATH="$GKILO_PATH" # Assuming DESTINATION_PATH is the same as GKILO_PATH

# Check if the directory and binary exist
if [ -d "$GKILO_PATH" ]; then
  if [ -f "$GKILO_PATH/gkilo" ]; then
    # Remove the existing binary
    rm "$GKILO_PATH/gkilo"
    echo "Removed existing 'gkilo' binary from $GKILO_PATH"
  else
    echo "'gkilo' binary not found in $GKILO_PATH"
  fi
else
  # Create the directory if it doesn't exist
  mkdir -p "$GKILO_PATH"
  echo "Created directory $GKILO_PATH"
fi

# Build the new binary
go build -o "$DESTINATION_PATH/gkilo" ./cmd/main/main.go

# Set the environment variable
export GKILO_HOME="$GKILO_PATH"

# Add the directory to the PATH
export PATH="$GKILO_PATH:$PATH"

echo "Environment variable GKILO_HOME set to $GKILO_HOME"
echo "Directory $GKILO_PATH added to PATH"
