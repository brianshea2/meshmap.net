#!/bin/bash

bash "$(dirname "$0")/generate-pb.sh" || exit $?

docker build \
    --no-cache \
    --pull \
    -f "$(dirname "$0")/../Dockerfile.meshobserv" \
    -t meshobserv \
    "$(dirname "$0")/.."
