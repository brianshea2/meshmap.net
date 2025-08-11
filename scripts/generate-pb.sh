#!/bin/bash

docker run \
    --rm \
    -v "$(cd "$(dirname "$0")"; pwd)/..":/data \
    golang \
    bash -c "
        apt-get update &&
        apt-get install -y protobuf-compiler &&
        go install google.golang.org/protobuf/cmd/protoc-gen-go@latest &&
        cd /tmp &&
        git clone --progress --depth 1 https://github.com/meshtastic/protobufs.git &&
        rm -rf /data/internal/meshtastic/generated &&
        protoc \\
            -I=protobufs \\
            --go_out=/data/internal/meshtastic \\
            --go_opt=module=github.com/meshtastic/go \\
            protobufs/nanopb.proto \\
            protobufs/meshtastic/*.proto
    "
