#!/bin/bash

docker build \
    --no-cache \
    --pull \
    -f "$(dirname "$0")/../Dockerfile.meshobserv" \
    -t meshobserv \
    "$(dirname "$0")/.."
