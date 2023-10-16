#!/bin/bash

# Specify the path to the directory where the "gkilo" binary is located
GKILO_PATH="./bin"
go build -o "$DESTINATION_PATH/gkilo" ./cmd/main/main.go
# Check if the directory and binary exist
if [ -d "$GKILO_PATH" ]; then
  if [ -f "$GKILO_PATH/gkilo" ]; then
    # Set the environment variable
    export GKILO_HOME="$GKILO_PATH"

    # Add the directory to the PATH
    export PATH="$GKILO_PATH:$PATH"

    echo "Environment variable GKILO_HOME set to $GKILO_HOME"
    echo "Directory $GKILO_PATH added to PATH"
  else
    echo "Error: 'gkilo' binary not found in $GKILO_PATH"
  fi
else
  echo "Error: Directory $GKILO_PATH does not exist"
fi

