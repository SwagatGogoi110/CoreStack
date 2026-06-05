#!/usr/bin/env bash
set -e

# Script to run unit tests for a specific service or all services

if [ $# -eq 0 ]; then
    echo "Usage: $0 <service_name> | all"
    echo "Example: $0 s3"
    echo "Example: $0 iam"
    echo "Example: $0 all"
    exit 1
fi

TARGET=$1

if [ "$TARGET" = "all" ]; then
    echo "Running unit tests for all services..."
    go test -v ./internal/services/...
else
    # Check if the service directory exists
    if [ ! -d "internal/services/$TARGET" ]; then
        echo "Error: Service '$TARGET' not found in internal/services/"
        exit 1
    fi
    
    echo "Running unit tests for service: $TARGET"
    go test -v ./internal/services/"$TARGET"/...
fi
